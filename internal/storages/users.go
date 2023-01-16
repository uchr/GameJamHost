package storages

import (
	"context"

	"GameJamPlatform/internal/models/users"
)

func (st *storage) CreateUser(ctx context.Context, user users.User) error {
	_, err := st.db.Exec(ctx, "INSERT INTO users (username, email, password, avatar, about) VALUES ($1, $2, $3, $4, $5)",
		user.Username, user.Email, user.Password, user.AvatarURL, user.About)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) GetUserByID(ctx context.Context, userID int) (*users.User, error) {
	row := st.db.QueryRow(ctx, "SELECT user_id, username, email, password, avatar, about FROM users WHERE user_id = $1", userID)

	var user users.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AvatarURL, &user.About)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (st *storage) GetUserByUsername(ctx context.Context, username string) (*users.User, error) {
	row := st.db.QueryRow(ctx, "SELECT user_id, username, email, password, avatar, about FROM users WHERE username = $1", username)

	var user users.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AvatarURL, &user.About)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (st *storage) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	row := st.db.QueryRow(ctx, "SELECT user_id, username, email, password, avatar, about FROM users WHERE email = $1", email)

	var user users.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AvatarURL, &user.About)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (st *storage) UpdateUser(ctx context.Context, user users.User) error {
	_, err := st.db.Exec(ctx, "UPDATE users SET username = $1, email = $2, password = $3, avatar = $4, about = $5 WHERE user_id = $6",
		user.Username, user.Email, user.Password, user.AvatarURL, user.About, user.ID)
	if err != nil {
		return err
	}
	return nil
}
