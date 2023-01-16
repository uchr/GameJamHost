package users

import (
	"context"

	users2 "GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/storages"
)

type users struct {
	repo storages.Repo
}

func NewUsers(repo storages.Repo) *users {
	return &users{repo: repo}
}

func (u *users) CreateUser(ctx context.Context, user users2.User, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hash
	return u.repo.CreateUser(ctx, user)
}

func (u *users) GetUserByID(ctx context.Context, userID int) (*users2.User, error) {
	return u.repo.GetUserByID(ctx, userID)
}

func (u *users) GetUserByUsername(ctx context.Context, username string) (*users2.User, error) {
	return u.repo.GetUserByUsername(ctx, username)
}

func (u *users) GetUserByEmail(ctx context.Context, email string) (*users2.User, error) {
	return u.repo.GetUserByEmail(ctx, email)
}

func (u *users) UpdateUser(ctx context.Context, user users2.User) error {
	if user.Password != nil {
		hash, err := hashPassword(string(user.Password))
		if err != nil {
			return err
		}
		user.Password = hash
	}

	return u.repo.UpdateUser(ctx, user)
}

func (u *users) CheckPassword(ctx context.Context, username, password string) (bool, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	return checkPassword(password, user.Password), nil
}
