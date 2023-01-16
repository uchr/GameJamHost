package servers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/forms"
	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/templates"
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

func (s *server) jamsListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jams, err := s.service.GetJams(r.Context())
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		pageData := templates.NewJamListPageData(jams)
		if err := s.tmpl.JamListTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) jamNewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageData := templates.NewJamEditFormPageData(true, gamejams.GameJam{}, nil)
		if err := s.tmpl.JamEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) jamCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jam, validationErrors, err := s.parseJamForm(r)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusBadRequest, err)
			return
		}
		if len(validationErrors) > 0 {
			pageData := templates.NewJamEditFormPageData(true, *jam, validationErrors)
			if err := s.tmpl.JamEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
				s.executeErrorPage(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}

		validationErrors, err = s.service.CreateJam(r.Context(), *jam)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
		if len(validationErrors) > 0 {
			pageData := templates.NewJamEditFormPageData(true, *jam, validationErrors)
			if err := s.tmpl.JamEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
				s.executeErrorPage(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}

		http.Redirect(w, r, "/jams/"+jam.URL, http.StatusSeeOther)
	}
}

func (s *server) jamOverviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusNotFound, err)
			return
		}

		pageData := templates.NewJamOverviewPageData(*jam)
		if err := s.tmpl.JamOverviewTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) jamEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusBadRequest, err)
			return
		}

		jam, err := s.service.GetJamByID(r.Context(), jamID)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusNotFound, err)
			return
		}

		pageData := templates.NewJamEditFormPageData(false, *jam, nil)
		if err := s.tmpl.JamEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/jams/"+jam.URL, http.StatusSeeOther)
	}
}

func (s *server) jamUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusBadRequest, err)
			return
		}

		jam, validationErrors, err := s.parseJamForm(r)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusBadRequest, err)
			return
		}
		if len(validationErrors) > 0 {
			pageData := templates.NewJamEditFormPageData(false, *jam, validationErrors)
			if err := s.tmpl.JamEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
				s.executeErrorPage(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}

		validationErrors, err = s.service.UpdateJam(r.Context(), jamID, *jam)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
		if validationErrors != nil {
			pageData := templates.NewJamEditFormPageData(false, *jam, validationErrors)
			if err := s.tmpl.JamEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
				s.executeErrorPage(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}

		http.Redirect(w, r, "/jams/"+jam.URL, http.StatusSeeOther)
	}
}

func (s *server) jamDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.service.DeleteJam(r.Context(), jamID)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/jams", http.StatusSeeOther)
	}
}

func (s *server) jamEntriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusNotFound, err)
			return
		}

		games, err := s.service.GetGames(r.Context(), jamURL)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}

		pageData := templates.NewJamEntriesPageData(*jam, games)
		if err := s.tmpl.JamEntriesTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}
