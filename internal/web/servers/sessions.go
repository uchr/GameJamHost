package servers

import (
	"context"
	"errors"
	"net/http"

	"GameJamPlatform/internal/models/sessions"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/services/sessionprovider"
)

const authedUserKey = "authed_user"

// isAuthorized checks if user is authorized and returns session if it is.
func (s *server) isAuthorized(r *http.Request) (*sessions.Session, error) {
	sessionCookie, err := s.cookieStore.Get(r, "session")
	if err != nil {
		return nil, err
	}

	token, ok := sessionCookie.Values["token"]
	if !ok {
		return nil, nil
	}

	session, err := s.sessionProvider.CheckAndUpdate(r.Context(), token.(string))
	if err != nil {
		if errors.Is(err, sessionprovider.ErrSessionNotAuthenticated) {
			return nil, nil
		}
		return nil, err
	}
	return session, nil
}

func (s *server) authorize(w http.ResponseWriter, r *http.Request, userID int) error {
	session, err := s.isAuthorized(r)
	if err != nil {
		return err
	}
	if session != nil {
		return nil
	}

	session, err = s.sessionProvider.Create(r.Context(), userID)
	if err != nil {
		return err
	}

	sessionCookie, err := s.cookieStore.Get(r, "session")
	if err != nil {
		return err
	}

	sessionCookie.Values["token"] = session.UID
	if err := sessionCookie.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.isAuthorized(r)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		if session == nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := s.users.GetUserByID(r.Context(), session.UserID)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		ctx := context.WithValue(r.Context(), authedUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) authedUser(r *http.Request) *users.User {
	value := r.Context().Value(authedUserKey)
	if value == nil {
		return nil
	}
	return value.(*users.User)
}
