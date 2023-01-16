package servers

import (
	"context"

	"GameJamPlatform/internal/models/users"
)

type Server interface {
	Run() error
}

type Users interface {
	CreateUser(ctx context.Context, user users.User, password string) error
	CheckPassword(ctx context.Context, username, password string) (bool, error)
	GetUserByID(ctx context.Context, userID int) (*users.User, error)
	GetUserByUsername(ctx context.Context, username string) (*users.User, error)
	GetUserByEmail(ctx context.Context, email string) (*users.User, error)
	UpdateUser(ctx context.Context, user users.User) error
}
