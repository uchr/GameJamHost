package servers

import (
	"net/http"

	"GameJamPlatform/internal/log"
	"GameJamPlatform/internal/templates"
)

func (s *server) executeErrorPage(w http.ResponseWriter, r *http.Request, errorCode int, err error) {
	log.Error(err)
	err = s.tmpl.ErrorTemplate.ExecuteTemplate(w, "base", templates.NewErrorPageData(errorCode, err))
	if err != nil {
		log.Error(err)
	}
}
