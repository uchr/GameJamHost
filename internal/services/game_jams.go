package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"GameJamPlatform/internal/forms"
	"GameJamPlatform/internal/gamejam"
)

func (sr *Service) validateJam(ctx context.Context, jam gamejam.GameJam) (forms.ValidationErrors, error) {
	const maxTitleLength = 64
	const maxURLLength = 64
	const maxContentLength = 10000

	validationErrors := make(forms.ValidationErrors)

	urlFormat := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !urlFormat.MatchString(jam.URL) {
		validationErrors["URL"] = "URL must only contain lowercase letters, numbers and dashes"
	}

	if len(jam.URL) == 0 || len(jam.URL) > maxURLLength {
		validationErrors["URL"] = fmt.Sprintf("URL length must be less than %d and not empty", maxURLLength)
	}
	if len(jam.Title) == 0 || len(jam.Title) > maxTitleLength {
		validationErrors["Title"] = fmt.Sprintf("Name length must be less than %d and not empty", maxTitleLength)
	}
	if len(jam.Content) > maxContentLength {
		validationErrors["Content"] = fmt.Sprintf("Content length must be less than %d", maxContentLength)
	}

	// TODO: Validate dates

	prevGameJam, err := sr.repo.GetJamID(ctx, jam.URL)
	if err == nil && prevGameJam != jam.ID {
		validationErrors["URL"] = "URL already exists"
	}

	return validationErrors, nil
}

func (sr *Service) GetJams(ctx context.Context) ([]gamejam.GameJam, error) {
	gameJams, err := sr.repo.GetJams(ctx)
	return gameJams, err
}

func (sr *Service) CreateJam(ctx context.Context, jam gamejam.GameJam) (forms.ValidationErrors, error) {
	jam.Title = strings.TrimSpace(jam.Title)
	validationErrors, err := sr.validateJam(ctx, jam)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	err = sr.repo.CreateJam(ctx, jam)
	return nil, err
}

func (sr *Service) DeleteJam(ctx context.Context, jamID int) error {
	err := sr.repo.DeleteJam(ctx, jamID)
	return err
}

func (sr *Service) GetJamByURL(ctx context.Context, jamURL string) (*gamejam.GameJam, error) {
	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	jam, err := sr.repo.GetJam(ctx, jamID)
	return jam, err
}

func (sr *Service) GetJamByID(ctx context.Context, jamID int) (*gamejam.GameJam, error) {
	jam, err := sr.repo.GetJam(ctx, jamID)
	return jam, err
}

func (sr *Service) UpdateJam(ctx context.Context, jamID int, jam gamejam.GameJam) (forms.ValidationErrors, error) {
	jam.Title = strings.TrimSpace(jam.Title)
	jam.ID = jamID

	validationErrors, err := sr.validateJam(ctx, jam)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	err = sr.repo.UpdateJam(ctx, jam)
	return nil, err
}

func (sr *Service) JamEntries(ctx context.Context, jamURL string) ([]gamejam.Game, error) {
	jamID, err := sr.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	jam, err := sr.repo.GetJam(ctx, jamID)
	if err != nil {
		return nil, err
	}

	games, err := sr.repo.GetGames(ctx, jam.ID)
	return games, err
}
