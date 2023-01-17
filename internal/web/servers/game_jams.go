package servers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/web/forms"
	"GameJamPlatform/internal/web/pagedata"
)

func (s *server) parseJamForm(r *http.Request) (*gamejams.GameJam, forms.ValidationErrors, error) {
	const maxUploadSize = 10 * 1024 * 1024 // 10 mb
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		return nil, nil, err
	}

	validationErrors := make(forms.ValidationErrors)

	jam := gamejams.GameJam{
		Title:           r.FormValue("name"),
		URL:             r.FormValue("url"),
		Content:         r.FormValue("content"),
		HideResults:     r.FormValue("hideResults") == "on",
		HideSubmissions: r.FormValue("hideSubmissions") == "on",
	}

	jam.StartDate, err = time.Parse(forms.TimeLayout, r.FormValue("startDate"))
	if err != nil {
		validationErrors["StartDate"] = "Must be a valid date"
	}
	jam.EndDate, err = time.Parse(forms.TimeLayout, r.FormValue("endDate"))
	if err != nil {
		validationErrors["EndDate"] = "Must be a valid date"
	}
	jam.VotingEndDate, err = time.Parse(forms.TimeLayout, r.FormValue("votingEndDate"))
	if err != nil {
		validationErrors["VotingEndDate"] = "Must be a valid date"
	}

	coverImageURL, err := s.uploadImage(r, "CoverImage")
	if err != nil {
		return nil, nil, err
	}
	if coverImageURL != "" {
		jam.CoverImageURL = coverImageURL
	}

	return &jam, validationErrors, nil
}

func (s *server) jamsListHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "jam_list"

		user := s.authedUser(r)

		jams, err := s.service.GetJams(r.Context())
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		pageData := pagedata.NewJamListPageData(user, jams)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) jamNewHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "jam_edit_form"

		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		pageData := pagedata.NewJamEditFormPageData(*user, gamejams.GameJam{}, true, nil)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) jamNewHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "jam_edit_form"

		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jam, validationErrors, err := s.parseJamForm(r)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}
		if len(validationErrors) > 0 {
			pageData := pagedata.NewJamEditFormPageData(*user, *jam, true, validationErrors)
			s.tm.Render(w, pageName, pageData)
			return
		}

		validationErrors, err = s.service.CreateJam(r.Context(), *jam)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		if len(validationErrors) > 0 {
			pageData := pagedata.NewJamEditFormPageData(*user, *jam, true, validationErrors)
			s.tm.Render(w, pageName, pageData)
			return
		}

		http.Redirect(w, r, "/jams/"+jam.URL, http.StatusSeeOther)
	}
}

func (s *server) jamOverviewHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "jam_overview"

		user := s.authedUser(r)

		jamURL := chi.URLParam(r, "jamURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		pageData := pagedata.NewJamOverviewPageData(user, *jam)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) jamEditHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "jam_edit_form"

		// TODO: check if user is host of jam
		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		jam, err := s.service.GetJamByID(r.Context(), jamID)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		pageData := pagedata.NewJamEditFormPageData(*user, *jam, false, nil)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) jamEditHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "jam_edit_form"

		// TODO: check if user is host of jam
		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		jam, validationErrors, err := s.parseJamForm(r)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}
		if len(validationErrors) > 0 {
			pageData := pagedata.NewJamEditFormPageData(*user, *jam, false, validationErrors)
			s.tm.Render(w, pageName, pageData)
			return
		}

		validationErrors, err = s.service.UpdateJam(r.Context(), jamID, *jam)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		if validationErrors != nil {
			pageData := pagedata.NewJamEditFormPageData(*user, *jam, false, validationErrors)
			s.tm.Render(w, pageName, pageData)
			return
		}

		http.Redirect(w, r, "/jams/"+jam.URL, http.StatusSeeOther)
	}
}

func (s *server) jamDeleteHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: check if user is host of jam
		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		err = s.service.DeleteJam(r.Context(), jamID)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/jams", http.StatusSeeOther)
	}
}

func (s *server) jamEntriesHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "jam_entries"

		user := s.authedUser(r)

		jamURL := chi.URLParam(r, "jamURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		games, err := s.service.GetGames(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		pageData := pagedata.NewJamEntriesPageData(user, *jam, games)
		s.tm.Render(w, pageName, pageData)
	}
}
