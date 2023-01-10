package storages

import (
	"context"

	"GameJamPlatform/internal/gamejam"
)

type Repo interface {
	CreateJam(ctx context.Context, jam gamejam.GameJam) error
	UpdateJam(ctx context.Context, jam gamejam.GameJam) error
	DeleteJam(ctx context.Context, jamID int) error
	GetJamID(ctx context.Context, jamURL string) (int, error)
	GetJam(ctx context.Context, jamID int) (*gamejam.GameJam, error)
	GetJams(ctx context.Context) ([]gamejam.GameJam, error)
	// IsNotFoundJam(err error) bool // TODO: implement this

	CreateGame(ctx context.Context, game gamejam.Game) error
	GetGame(ctx context.Context, jamID int, gameURL string) (*gamejam.Game, error)
	UpdateGame(ctx context.Context, game gamejam.Game) error
	GetGames(ctx context.Context, gameJamID int) ([]gamejam.Game, error)
	DeleteGame(ctx context.Context, gameID int) error
	BanGame(ctx context.Context, jamID int, gameID string) error
	// IsNotFoundGame(err error) bool // TODO: implement this
}
