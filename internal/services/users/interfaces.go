package users

import (
	"context"

	usersModel "GameJamPlatform/internal/models/users"
)

type Users interface {
	CreateUser(ctx context.Context, user usersModel.User, password string) error
	CheckPassword(ctx context.Context, username, password string) (bool, error)
	GetUserByID(ctx context.Context, userID int) (*usersModel.User, error)
	GetUserByUsername(ctx context.Context, username string) (*usersModel.User, error)
	GetUserByEmail(ctx context.Context, email string) (*usersModel.User, error)
	UpdateUser(ctx context.Context, user usersModel.User) error
}
