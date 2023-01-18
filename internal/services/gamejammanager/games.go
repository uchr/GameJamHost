package gamejammanager

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/services/validationerr"
	"GameJamPlatform/internal/storages"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (jm *gameJamManager) validateGame(game gamejams.Game) error {
	const maxTitleLength = 64
	const maxBuildLength = 1000
	const maxContentLength = 10000

	vErr := validationerr.New()

	if len(game.Title) == 0 || len(game.Title) > maxTitleLength {
		vErr.Add("Title", fmt.Sprintf("Title must be less than %d characters and not empty", maxTitleLength))
	}
	if len(game.Build) == 0 || len(game.Build) > maxBuildLength {
		vErr.Add("Build", fmt.Sprintf("Build must be less than %d characters and not empty", maxBuildLength))
	}
	if len(game.Content) > maxContentLength {
		vErr.Add("Content", fmt.Sprintf("Content must be less than %d characters and not empty", maxContentLength))
	}

	if vErr.HasErrors() {
		return vErr
	}

	return nil
}

func generateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz"

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = letters[b[i]%byte(len(letters))]
	}
	return string(b), nil
}

func urlFromTitle(gameName string) string {
	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9- ]+`)
	result := nonAlphanumericRegex.ReplaceAllString(gameName, " ")
	spaceRegex := regexp.MustCompile(`\s+`)
	result = spaceRegex.ReplaceAllString(result, "-")
	result = strings.ToLower(result)

	return result
}

func (jm *gameJamManager) generateGameUrl(ctx context.Context, jamID int, gameName string) (string, error) {
	baseUrl := urlFromTitle(gameName)
	suffix := ""

	const maxTries = 10
	for i := 0; i < maxTries; i++ {
		resultUrl := baseUrl
		if suffix != "" {
			resultUrl += "-" + suffix
		}
		_, err := jm.repo.GetGame(ctx, jamID, resultUrl)
		if errors.Is(err, storages.ErrNotFound) {
			return resultUrl, nil
		}

		suffix, err = generateRandomString(5)
		if err != nil {
			return "", err
		}
	}

	return "", errors.New("failed to generate unique game URL")
}

func (jm *gameJamManager) CreateGame(ctx context.Context, jamURL string, user users.User, game gamejams.Game) (string, error) {
	err := jm.validateGame(game)
	if err != nil {
		return "", err
	}

	jamID, err := jm.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return "", err
	}

	game.Title = strings.TrimSpace(game.Title)
	game.JamID = jamID
	game.URL, err = jm.generateGameUrl(ctx, jamID, game.Title)
	game.UserID = user.ID
	if err != nil {
		return "", err
	}

	err = jm.repo.CreateGame(ctx, game)
	if err != nil {
		return "", err
	}

	return game.URL, nil
}

func (jm *gameJamManager) UpdateGame(ctx context.Context, jamURL string, gameURL string, game gamejams.Game) error {
	err := jm.validateGame(game)
	if err != nil {
		return err
	}

	jamID, err := jm.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return err
	}

	prevGame, err := jm.repo.GetGame(ctx, jamID, gameURL)
	if err != nil {
		return err
	}

	prevGame.Title = game.Title
	prevGame.Content = game.Content
	prevGame.Build = game.Build
	if game.CoverImageURL != "" {
		prevGame.CoverImageURL = game.CoverImageURL
	}
	if game.ScreenshotURLs != nil {
		prevGame.ScreenshotURLs = game.ScreenshotURLs
	}

	err = jm.repo.UpdateGame(ctx, *prevGame)
	return err
}

func (jm *gameJamManager) BanGame(ctx context.Context, jamURL string, gameURL string) error {
	jamID, err := jm.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return err
	}

	err = jm.repo.BanGame(ctx, jamID, gameURL)
	return err
}

func (jm *gameJamManager) GetGame(ctx context.Context, jamURL string, gameURL string) (*gamejams.Game, error) {
	jamID, err := jm.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	game, err := jm.repo.GetGame(ctx, jamID, gameURL)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (jm *gameJamManager) GetGames(ctx context.Context, jamURL string) ([]gamejams.Game, error) {
	jamID, err := jm.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	games, err := jm.repo.GetGames(ctx, jamID)
	return games, err
}

func (jm *gameJamManager) GetGamesByUserID(ctx context.Context, userID int) ([]gamejams.Game, error) {
	games, err := jm.repo.GetGamesByUserID(ctx, userID)
	return games, err
}

func (jm *gameJamManager) IsGameOwner(_ context.Context, game gamejams.Game, user *users.User) (bool, error) {
	if user == nil {
		return false, nil
	}

	return game.UserID == user.ID, nil
}

func (jm *gameJamManager) IsGameOwnerByURL(ctx context.Context, jamURL string, gameURL string, user *users.User) (bool, error) {
	if user == nil {
		return false, nil
	}

	jamID, err := jm.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return false, err
	}

	game, err := jm.repo.GetGame(ctx, jamID, gameURL)
	if err != nil {
		return false, err
	}

	return game.UserID == user.ID, nil
}
