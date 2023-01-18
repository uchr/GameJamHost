package usersmanager

import (
	"context"

	usersModel "GameJamPlatform/internal/models/users"
)

type UserManager interface {
	// CreateUser creates a new user in the database.
	CreateUser(ctx context.Context, user usersModel.User, password string) error

	// CheckPassword checks if the password is correct for the given username.
	CheckPassword(ctx context.Context, username, password string) (bool, error)

	// GetUserByID returns user by id or nil if user not found.
	GetUserByID(ctx context.Context, userID int) (*usersModel.User, error)

	// GetUserByUsername returns user by username or nil if user not found.
	GetUserByUsername(ctx context.Context, username string) (*usersModel.User, error)

	// GetUserByEmail returns user by email or nil if user not found.
	GetUserByEmail(ctx context.Context, email string) (*usersModel.User, error)

	// UpdateUser updates user in the database.
	UpdateUser(ctx context.Context, user usersModel.User, password string) error
}
