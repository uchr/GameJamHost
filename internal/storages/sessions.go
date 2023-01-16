package storages

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"GameJamPlatform/internal/models/sessions"
)

func (st *storage) CreateSession(ctx context.Context, session sessions.Session) error {
	_, err := st.db.Exec(ctx, "INSERT INTO sessions (session_id, user_id, expire_at) VALUES ($1, $2, $3)", session.UID, session.UserID, session.ExpireAt)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) GetSession(ctx context.Context, sessionID string) (*sessions.Session, error) {
	row := st.db.QueryRow(ctx, "SELECT session_id, user_id, expire_at FROM sessions WHERE session_id = $1", sessionID)

	var session sessions.Session
	err := row.Scan(&session.UID, &session.UserID, &session.ExpireAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (st *storage) UpdateSession(ctx context.Context, session sessions.Session) error {
	_, err := st.db.Exec(ctx, "UPDATE sessions SET user_id = $1, expire_at = $2 WHERE session_id = $3", session.UserID, session.ExpireAt, session.UID)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := st.db.Exec(ctx, "DELETE FROM sessions WHERE session_id = $1", sessionID)
	if err != nil {
		return err
	}
	return nil
}
