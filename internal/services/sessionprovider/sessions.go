package sessionprovider

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"GameJamPlatform/internal/models/sessions"
	"GameJamPlatform/internal/storages"
)

type sessionsProvider struct {
	repo storages.Repo
}

var _ = SessionProvider(&sessionsProvider{})

func NewProvider(repo storages.Repo) *sessionsProvider {
	return &sessionsProvider{repo: repo}
}

func (sp *sessionsProvider) Create(ctx context.Context, userID int) (*sessions.Session, error) {
	session := sessions.Session{
		UID:      uuid.New().String(),
		UserID:   userID,
		ExpireAt: time.Now().Add(7 * 24 * time.Hour),
	}

	err := sp.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (sp *sessionsProvider) GetAndUpdate(ctx context.Context, sessionID string) (*sessions.Session, error) {
	session, err := sp.repo.GetSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, storages.ErrNotFound) {
			return nil, ErrSessionNotAuthenticated
		}
		return nil, err
	}

	if session.ExpireAt.Before(time.Now()) {
		err = sp.repo.DeleteSession(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		return nil, ErrSessionNotAuthenticated
	}

	session.ExpireAt = time.Now().Add(7 * 24 * time.Hour)
	err = sp.repo.UpdateSession(ctx, *session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (sp *sessionsProvider) CheckAndUpdate(ctx context.Context, sessionID string) (*sessions.Session, error) {
	session, err := sp.repo.GetSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, storages.ErrNotFound) {
			return nil, ErrSessionNotAuthenticated
		}
		return nil, err
	}

	if session.ExpireAt.Before(time.Now()) {
		err = sp.repo.DeleteSession(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		return nil, ErrSessionNotAuthenticated
	}

	session.ExpireAt = time.Now().Add(7 * 24 * time.Hour)
	err = sp.repo.UpdateSession(ctx, *session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (sp *sessionsProvider) Delete(ctx context.Context, sessionID string) error {
	return sp.repo.DeleteSession(ctx, sessionID)
}
