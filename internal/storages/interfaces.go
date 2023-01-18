package storages

import (
	"context"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/sessions"
	"GameJamPlatform/internal/models/users"
)

type Repo interface {
	CreateJam(ctx context.Context, jam gamejams.GameJam) error
	UpdateJam(ctx context.Context, jam gamejams.GameJam) error
	DeleteJam(ctx context.Context, jamID int) error
	GetJamID(ctx context.Context, jamURL string) (int, error)
	GetJam(ctx context.Context, jamID int) (*gamejams.GameJam, error)
	GetJams(ctx context.Context) ([]gamejams.GameJam, error)
	GetJamsByUserID(ctx context.Context, userID int) ([]gamejams.GameJam, error)

	CreateGame(ctx context.Context, game gamejams.Game) error
	GetGame(ctx context.Context, jamID int, gameURL string) (*gamejams.Game, error)
	UpdateGame(ctx context.Context, game gamejams.Game) error
	GetGames(ctx context.Context, gameJamID int) ([]gamejams.Game, error)
	GetGamesByUserID(ctx context.Context, userID int) ([]gamejams.Game, error)
	DeleteGame(ctx context.Context, gameID int) error
	BanGame(ctx context.Context, jamID int, gameID string) error

	CreateUser(ctx context.Context, user users.User) error
	GetUserByID(ctx context.Context, userID int) (*users.User, error)
	GetUserByEmail(ctx context.Context, email string) (*users.User, error)
	GetUserByUsername(ctx context.Context, username string) (*users.User, error)
	UpdateUser(ctx context.Context, user users.User) error

	CreateSession(ctx context.Context, session sessions.Session) error
	GetSession(ctx context.Context, sessionID string) (*sessions.Session, error)
	UpdateSession(ctx context.Context, session sessions.Session) error
	DeleteSession(ctx context.Context, sessionID string) error
}
