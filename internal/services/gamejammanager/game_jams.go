package gamejammanager

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/web/forms"
)

func (jm *gameJamManager) validateJam(ctx context.Context, jam gamejams.GameJam) (forms.ValidationErrors, error) {
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

	prevGameJam, err := jm.repo.GetJamID(ctx, jam.URL)
	if err == nil && prevGameJam != jam.ID {
		validationErrors["URL"] = "URL already exists"
	}

	return validationErrors, nil
}

func (jm *gameJamManager) GetJams(ctx context.Context) ([]gamejams.GameJam, error) {
	gameJams, err := jm.repo.GetJams(ctx)
	return gameJams, err
}

func (jm *gameJamManager) GetJamsByUserID(ctx context.Context, userID int) ([]gamejams.GameJam, error) {
	jams, err := jm.repo.GetJamsByUserID(ctx, userID)
	return jams, err
}

func (jm *gameJamManager) CreateJam(ctx context.Context, user users.User, jam gamejams.GameJam) (forms.ValidationErrors, error) {
	jam.Title = strings.TrimSpace(jam.Title)
	validationErrors, err := jm.validateJam(ctx, jam)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	jam.UserID = user.ID
	err = jm.repo.CreateJam(ctx, jam)
	return nil, err
}

func (jm *gameJamManager) DeleteJam(ctx context.Context, jamID int) error {
	err := jm.repo.DeleteJam(ctx, jamID)
	return err
}

func (jm *gameJamManager) GetJamByURL(ctx context.Context, jamURL string) (*gamejams.GameJam, error) {
	jamID, err := jm.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	jam, err := jm.repo.GetJam(ctx, jamID)
	return jam, err
}

func (jm *gameJamManager) GetJamByID(ctx context.Context, jamID int) (*gamejams.GameJam, error) {
	jam, err := jm.repo.GetJam(ctx, jamID)
	return jam, err
}

func (jm *gameJamManager) UpdateJam(ctx context.Context, jamID int, jam gamejams.GameJam) (forms.ValidationErrors, error) {
	jam.Title = strings.TrimSpace(jam.Title)
	jam.ID = jamID

	validationErrors, err := jm.validateJam(ctx, jam)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	err = jm.repo.UpdateJam(ctx, jam)
	return nil, err
}

func (jm *gameJamManager) JamEntries(ctx context.Context, jamURL string) ([]gamejams.Game, error) {
	jamID, err := jm.repo.GetJamID(ctx, jamURL)
	if err != nil {
		return nil, err
	}

	jam, err := jm.repo.GetJam(ctx, jamID)
	if err != nil {
		return nil, err
	}

	games, err := jm.repo.GetGames(ctx, jam.ID)
	return games, err
}

func (jm *gameJamManager) IsHost(_ context.Context, jam gamejams.GameJam, user *users.User) (bool, error) {
	if user == nil {
		return false, nil
	}

	return jam.UserID == user.ID, nil
}

func (jm *gameJamManager) IsHostByID(ctx context.Context, jamID int, user *users.User) (bool, error) {
	if user == nil {
		return false, nil
	}

	jam, err := jm.repo.GetJam(ctx, jamID)
	if err != nil {
		return false, err
	}

	return jam.UserID == user.ID, nil
}
