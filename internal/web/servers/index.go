package servers

import (
	"net/http"

	"GameJamPlatform/internal/web/pagedata"
)

func (s *server) indexHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "index"

		user := s.authedUser(r)

		pageData := pagedata.NewAuthPageData(user)
		s.tm.Render(w, pageName, pageData)
	}
}
