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

func (s *server) userNewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.tmpl.UserRegistrationTemplate.ExecuteTemplate(w, "base", nil); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) userCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, password, err := s.parseRegistrationForm(r)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.users.CreateUser(r.Context(), *user, password); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.authorize(w, r, user.ID)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *server) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.tmpl.UserLoginTemplate.ExecuteTemplate(w, "base", nil); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.isAuthorized(r)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		if session != nil {
			if err := s.sessionProvider.Delete(r.Context(), session.UID); err != nil {
				s.executeErrorPage(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *server) authHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			s.executeErrorPage(w, r, http.StatusBadRequest, err)
			return
		}

		username := r.Form.Get("Username")
		password := r.Form.Get("Password")

		user, err := s.users.GetUserByUsername(r.Context(), username)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		if user == nil {
			s.executeErrorPage(w, r, http.StatusUnauthorized, ErrFailedAuth)
			return
		}

		ok, err := s.users.CheckPassword(r.Context(), username, password)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		if !ok {
			s.executeErrorPage(w, r, http.StatusUnauthorized, ErrFailedAuth)
			return
		}

		err = s.authorize(w, r, user.ID)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
