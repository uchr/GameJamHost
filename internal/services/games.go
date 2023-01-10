package services

import (
	"context"
	"errors"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"GameJamPlatform/internal/gamejam"
	"GameJamPlatform/internal/storages"
)

func init() {
	rand.Seed(time.Now().UnixNano())
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
		if errors.Is(err, storages.ErrGameNotFound) {
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
func (sr *Service) CreateGame(ctx context.Context, jamURL string, game gamejam.Game) (string, error) {
	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return "", err
	}

	game.GameJamID = jamID
	game.URL, err = sr.generateGameUrl(ctx, jamID, game.Name)
	if err != nil {
		return "", err
	}

	err = sr.repo.CreateGame(ctx, game)
	return game.URL, err
}

func (sr *Service) UpdateGame(ctx context.Context, jamURL string, gameURL string, game gamejam.Game) error {
	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return err
	}

	prevGame, err := sr.repo.GetGame(ctx, jamID, gameURL)
	if err != nil {
		return err
	}

	prevGame.Name = game.Name
	prevGame.Content = game.Content
	prevGame.Build = game.Build

	err = sr.repo.UpdateGame(ctx, *prevGame)
	return err
}

func (sr *Service) BanGame(ctx context.Context, jamURL string, gameURL string) error {
	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return err
	}

	err = sr.repo.BanGame(ctx, jamID, gameURL)
	return err
}

func (sr *Service) GetGame(ctx context.Context, jamURL string, gameURL string) (*gamejam.Game, error) {
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

func (sr *Service) GetGames(ctx context.Context, jamURL string) ([]gamejam.Game, error) {
	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	games, err := sr.repo.GetGames(ctx, jamID)
	return games, err
}