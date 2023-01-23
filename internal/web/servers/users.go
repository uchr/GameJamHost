package servers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/web/pagedata"
)

// TODO: add validation
func (s *server) parseRegistrationForm(r *http.Request) (*users.User, string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, "", err
	}

	user := users.User{
		Username: r.Form.Get("username"),
		Email:    r.Form.Get("email"),
		About:    r.Form.Get("about"),
	}

	password := r.Form.Get("password")
	confirmPassword := r.Form.Get("confirm_password")
	if password != confirmPassword {
		return nil, "", ErrPasswordsNotMatch
	}

	return &user, password, nil
}

// TODO: add validation
func (s *server) parseEditForm(r *http.Request) (*users.User, string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, "", err
	}

	user := users.User{
		About: r.Form.Get("about"),
	}

	password := r.Form.Get("password")
	confirmPassword := r.Form.Get("confirm_password")
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

		username := r.Form.Get("username")
		password := r.Form.Get("password")

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

func (s *server) userProfileHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "user_profile"

		authedUser := s.authedUser(r)

		username := chi.URLParam(r, "username")
		profileUser, err := s.users.GetUserByUsername(r.Context(), username)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		if profileUser == nil {
			s.tm.RenderError(w, http.StatusNotFound, fmt.Errorf("user %s not found", username))
			return
		}

		games, err := s.gameJams.GetGamesByUserID(r.Context(), profileUser.ID)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		jams, err := s.gameJams.GetJamsByUserID(r.Context(), profileUser.ID)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		jamURLs := make(map[int]string)
		for _, game := range games {
			jam, err := s.gameJams.GetJamByID(r.Context(), game.JamID)
			if err != nil {
				s.tm.RenderError(w, http.StatusInternalServerError, err)
				return
			}
			jamURLs[jam.ID] = jam.URL
		}

		pageData := pagedata.NewUserProfilePageData(authedUser, *profileUser, jams, games, jamURLs)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) userEditHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "user_edit_form"

		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		pageData := pagedata.NewUserEditFormPageData(*user, nil)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) userEditHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		userForm, password, err := s.parseEditForm(r)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		userForm.ID = user.ID
		err = s.users.UpdateUser(r.Context(), *userForm, password)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/users/%s", user.Username), http.StatusSeeOther)
	}
}
