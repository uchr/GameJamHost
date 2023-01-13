package storages

import (
	"context"

	"GameJamPlatform/internal/models"
)

type Repo interface {
	CreateJam(ctx context.Context, jam models.GameJam) error
	UpdateJam(ctx context.Context, jam models.GameJam) error
	DeleteJam(ctx context.Context, jamID int) error
	GetJamID(ctx context.Context, jamURL string) (int, error)
	GetJam(ctx context.Context, jamID int) (*models.GameJam, error)
	GetJams(ctx context.Context) ([]models.GameJam, error)
	// IsNotFoundJam(err error) bool // TODO: implement this

	CreateGame(ctx context.Context, game models.Game) error
	GetGame(ctx context.Context, jamID int, gameURL string) (*models.Game, error)
	UpdateGame(ctx context.Context, game models.Game) error
	GetGames(ctx context.Context, gameJamID int) ([]models.Game, error)
	DeleteGame(ctx context.Context, gameID int) error
	BanGame(ctx context.Context, jamID int, gameID string) error
	// IsNotFoundGame(err error) bool // TODO: implement this
}
