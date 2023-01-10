package services

import (
	"context"

	"GameJamPlatform/internal/gamejam"
	"GameJamPlatform/internal/storages"
)

type Service struct {
	repo storages.Repo
}

func NewService(repo storages.Repo) *Service {
	return &Service{repo: repo}
}

func (sr *Service) GetJams(ctx context.Context) ([]gamejam.GameJam, error) {
	gameJams, err := sr.repo.GetJams(ctx)
	return gameJams, err
}

func (sr *Service) CreateJam(ctx context.Context, jam gamejam.GameJam) error {
	err := sr.repo.CreateJam(ctx, jam)
	return err
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

func (sr *Service) UpdateJam(ctx context.Context, jamID int, jam gamejam.GameJam) error {
	jam.ID = jamID
	err := sr.repo.UpdateJam(ctx, jam)
	return err
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
