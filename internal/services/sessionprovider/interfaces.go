package sessionprovider

import (
	"context"

	"GameJamPlatform/internal/models/sessions"
)

type SessionProvider interface {
	// Create creates a new session for the given user ID. It returns the session ID.
	Create(ctx context.Context, userID int) (*sessions.Session, error)

	// GetAndUpdate returns the session for the given UID, and updates the session's expiration time.
	GetAndUpdate(ctx context.Context, sessionID string) (*sessions.Session, error)

	// CheckAndUpdate checks if the given session ID is valid, and updates the session's expiration time.
	// If the session is valid, it returns the user ID.
	CheckAndUpdate(ctx context.Context, sessionID string) (*sessions.Session, error)

	// Delete deletes the session with the given UID.
	Delete(ctx context.Context, sessionID string) error
}
