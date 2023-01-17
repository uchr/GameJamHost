package servers

import (
	"net/http"

	"GameJamPlatform/internal/models/users"
)

func (s *server) parseRegistrationForm(r *http.Request) (*users.User, string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, "", err
	}

	user := users.User{
		Username: r.Form.Get("Username"),
		Email:    r.Form.Get("Email"),
		About:    r.Form.Get("About"),
	}

	password := r.Form.Get("Password")
	confirmPassword := r.Form.Get("ConfirmPassword")
	if password != confirmPassword {
		return nil, "", ErrPasswordsNotMatch
	}

	return &user, password, nil
}

func (s *server) userNewHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "user_registration"

		user := s.authedUser(r)
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		s.tm.Render(w, pageName, nil)
	}
}

func (s *server) userNewHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, password, err := s.parseRegistrationForm(r)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		if err := s.users.CreateUser(r.Context(), *user, password); err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		user, err = s.users.GetUserByUsername(r.Context(), user.Username)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		err = s.login(w, r, user.ID)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *server) loginHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "user_login"

		user := s.authedUser(r)
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		s.tm.Render(w, pageName, nil)
	}
}

func (s *server) logoutHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.isAuthorized(r)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		if session != nil {
			err := s.logout(w, r, *session)
			if err != nil {
				s.tm.RenderError(w, http.StatusInternalServerError, err)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *server) loginHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		username := r.Form.Get("Username")
		password := r.Form.Get("Password")

		user, err := s.users.GetUserByUsername(r.Context(), username)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, ErrFailedAuth)
			return
		}

		ok, err := s.users.CheckPassword(r.Context(), username, password)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		if !ok {
			s.tm.RenderError(w, http.StatusUnauthorized, ErrFailedAuth)
			return
		}

		err = s.login(w, r, user.ID)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
