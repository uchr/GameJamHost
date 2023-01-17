package gamejammanager

import (
	"context"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/web/forms"
)

type GameJamManager interface {
	GetJams(ctx context.Context) ([]gamejams.GameJam, error)
	GetJamsByUserID(ctx context.Context, userID int) ([]gamejams.GameJam, error)
	CreateJam(ctx context.Context, user users.User, jam gamejams.GameJam) (forms.ValidationErrors, error)
	DeleteJam(ctx context.Context, jamID int) error
	GetJamByURL(ctx context.Context, jamURL string) (*gamejams.GameJam, error)
	GetJamByID(ctx context.Context, jamID int) (*gamejams.GameJam, error)
	UpdateJam(ctx context.Context, jamID int, jam gamejams.GameJam) (forms.ValidationErrors, error)
	JamEntries(ctx context.Context, jamURL string) ([]gamejams.Game, error)
	IsHost(ctx context.Context, jam gamejams.GameJam, user *users.User) (bool, error)
	IsHostByID(ctx context.Context, jamID int, user *users.User) (bool, error)

	// CreateGame creates a new game in the database and returns the game's URL.
	CreateGame(ctx context.Context, jamURL string, user users.User, game gamejams.Game) (string, forms.ValidationErrors, error)
	UpdateGame(ctx context.Context, jamURL string, gameURL string, game gamejams.Game) (forms.ValidationErrors, error)
	BanGame(ctx context.Context, jamURL string, gameURL string) error
	GetGame(ctx context.Context, jamURL string, gameURL string) (*gamejams.Game, error)
	GetGames(ctx context.Context, jamURL string) ([]gamejams.Game, error)
	GetGamesByUserID(ctx context.Context, userID int) ([]gamejams.Game, error)
	IsGameOwner(ctx context.Context, game gamejams.Game, user *users.User) (bool, error)
	IsGameOwnerByURL(ctx context.Context, jamURL string, gameURL string, user *users.User) (bool, error)
}
