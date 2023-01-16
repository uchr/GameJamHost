package servers

import (
	"errors"
	"net/http"

	"GameJamPlatform/internal/models/sessions"
	"GameJamPlatform/internal/services/sessionprovider"
)

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
