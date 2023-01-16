package servers

import (
	"net/http"

	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/templates"
)

func (s *server) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.isAuthorized(r)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		var user *users.User
		if session != nil {
			user, err = s.users.GetUserByID(r.Context(), session.UserID)
			if err != nil {
				s.executeErrorPage(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		indexPageData := templates.NewIndexPageData(user)
		if err := s.tmpl.IndexTemplate.ExecuteTemplate(w, "base", indexPageData); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}
