package sessionprovider

import (
	"context"

	"GameJamPlatform/internal/models/sessions"
)

type SessionProvider interface {
	// Create creates a new session for the given user ID. It returns the session ID.
	Create(ctx context.Context, userID int) (*sessions.Session, error)

	// Check checks if the given session ID is valid. If it is, it returns the session.
	Check(ctx context.Context, sessionID string) (*sessions.Session, error)

	// Delete deletes the session with the given UID.
	Delete(ctx context.Context, sessionID string) error
}
