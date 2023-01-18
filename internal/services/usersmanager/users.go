package usersmanager

import (
	"context"
	"errors"

	usersModel "GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/storages"
)

type users struct {
	repo storages.Repo
}

func NewUsers(repo storages.Repo) *users {
	return &users{repo: repo}
}

func (u *users) CreateUser(ctx context.Context, user usersModel.User, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hash
	return u.repo.CreateUser(ctx, user)
}

func (u *users) GetUserByID(ctx context.Context, userID int) (*usersModel.User, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storages.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *users) GetUserByUsername(ctx context.Context, username string) (*usersModel.User, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, storages.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *users) GetUserByEmail(ctx context.Context, email string) (*usersModel.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storages.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *users) UpdateUser(ctx context.Context, user usersModel.User, password string) error {
	prevUser, err := u.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil
	}

	if password != "" {
		hash, err := hashPassword(password)
		if err != nil {
			return err
		}
		prevUser.Password = hash
	}
	if user.About != "" {
		prevUser.About = user.About
	}
	if user.AvatarURL != "" {
		prevUser.AvatarURL = user.AvatarURL
	}

	return u.repo.UpdateUser(ctx, *prevUser)
}

func (u *users) CheckPassword(ctx context.Context, username, password string) (bool, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	return checkPassword(password, user.Password), nil
}
