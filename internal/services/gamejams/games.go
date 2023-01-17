package gamejams

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/storages"
	"GameJamPlatform/internal/web/forms"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (sr *Service) validateGame(game gamejams.Game) forms.ValidationErrors {
	const maxTitleLength = 64
	const maxBuildLength = 1000
	const maxContentLength = 10000

	validationErrors := make(forms.ValidationErrors)

	if len(game.Title) == 0 || len(game.Title) > maxTitleLength {
		validationErrors["Title"] = fmt.Sprintf("Title must be less than %d characters and not empty", maxTitleLength)
	}
	if len(game.Build) == 0 || len(game.Build) > maxBuildLength {
		validationErrors["Build"] = fmt.Sprintf("Build must be less than %d characters and not empty", maxBuildLength)
	}
	if len(game.Content) > maxContentLength {
		validationErrors["Content"] = fmt.Sprintf("Content must be less than %d characters and not empty", maxContentLength)
	}

	return validationErrors
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
	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	result := nonAlphanumericRegex.ReplaceAllString(gameName, " ")
	spaceRegex := regexp.MustCompile(`\s+`)
	result = spaceRegex.ReplaceAllString(result, "-")
	result = strings.ToLower(result)

	return result
}

func (sr *Service) generateGameUrl(ctx context.Context, jamID int, gameName string) (string, error) {
	baseUrl := urlFromTitle(gameName)
	suffix := ""

	const maxTries = 10
	for i := 0; i < maxTries; i++ {
		resultUrl := baseUrl + "-" + suffix
		_, err := sr.repo.GetGame(ctx, jamID, resultUrl)
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

// CreateGame creates a new game in the database and returns the game's URL.
func (sr *Service) CreateGame(ctx context.Context, jamURL string, game gamejams.Game) (string, forms.ValidationErrors, error) {
	validationErrors := sr.validateGame(game)
	if len(validationErrors) > 0 {
		return "", validationErrors, nil
	}

	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return "", nil, err
	}

	game.GameJamID = jamID
	game.URL, err = sr.generateGameUrl(ctx, jamID, game.Title)
	if err != nil {
		return "", nil, err
	}

	err = sr.repo.CreateGame(ctx, game)
	return game.URL, nil, err
}

func (sr *Service) UpdateGame(ctx context.Context, jamURL string, gameURL string, game gamejams.Game) (forms.ValidationErrors, error) {
	validationErrors := sr.validateGame(game)
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	prevGame, err := sr.repo.GetGame(ctx, jamID, gameURL)
	if err != nil {
		return nil, err
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

	err = sr.repo.UpdateGame(ctx, *prevGame)
	return nil, err
}

func (sr *Service) BanGame(ctx context.Context, jamURL string, gameURL string) error {
	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return err
	}

	err = sr.repo.BanGame(ctx, jamID, gameURL)
	return err
}

func (sr *Service) GetGame(ctx context.Context, jamURL string, gameURL string) (*gamejams.Game, error) {
	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	game, err := sr.repo.GetGame(ctx, jamID, gameURL)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (sr *Service) GetGames(ctx context.Context, jamURL string) ([]gamejams.Game, error) {
	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	games, err := sr.repo.GetGames(ctx, jamID)
	return games, err
}
