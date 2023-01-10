package servers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/gamejam"
	"GameJamPlatform/internal/log"
	"GameJamPlatform/internal/templates"
)

func (s *server) parseJamForm(r *http.Request) (*gamejam.GameJam, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	jam := gamejam.GameJam{
		Name:            r.FormValue("name"),
		URL:             r.FormValue("url"),
		Content:         r.FormValue("content"),
		HideResults:     r.FormValue("hideResults") == "on",
		HideSubmissions: r.FormValue("hideSubmissions") == "on",
	}

	jam.StartDate, err = time.Parse(templates.TimeLayout, r.FormValue("startDate"))
	if err != nil {
		return nil, err
	}
	jam.EndDate, err = time.Parse(templates.TimeLayout, r.FormValue("endDate"))
	if err != nil {
		return nil, err
	}
	jam.VotingEndDate, err = time.Parse(templates.TimeLayout, r.FormValue("votingEndDate"))
	if err != nil {
		return nil, err
	}

	return &jam, nil
}

func (s *server) jamsListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jams, err := s.service.GetJams(r.Context())
		if err != nil {
			log.Error(err)

			err = s.tmpl.ErrorTemplate.Execute(w, templates.NewErrorPageData(http.StatusInternalServerError))
			if err != nil {
				log.Error(err)
			}
			return
		}

		pageData := templates.NewJamListPageData(jams)
		err = s.tmpl.JamListTemplate.ExecuteTemplate(w, "base", pageData)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *server) jamNewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.tmpl.JamNewTemplate.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *server) jamCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jam, err := s.parseJamForm(r)
		if err != nil {
			log.Error(err)
			return
		}

		err = s.service.CreateJam(r.Context(), *jam)
		if err != nil {
			log.Error(err)
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
			log.Error(err)
			return
		}

		pageData := templates.NewJamOverviewPageData(*jam)
		err = s.tmpl.JamOverviewTemplate.ExecuteTemplate(w, "base", pageData)
		if err != nil {
			log.Error(err)
			return
		}

	}
}

func (s *server) jamEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			log.Error(err)
			return
		}

		jam, err := s.service.GetJamByID(r.Context(), jamID)
		if err != nil {
			log.Error(err)
			return
		}

		pageData := templates.NewJamEditPageData(*jam)
		err = s.tmpl.JamEditTemplate.ExecuteTemplate(w, "base", pageData)
		if err != nil {
			log.Error(err)
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
			log.Error(err)
			return
		}

		jam, err := s.parseJamForm(r)
		if err != nil {
			log.Error(err)
			return
		}

		err = s.service.UpdateJam(r.Context(), jamID, *jam)
		if err != nil {
			log.Error(err)
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
			log.Error(err)
			return
		}

		err = s.service.DeleteJam(r.Context(), jamID)
		if err != nil {
			log.Error(err)
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
			log.Error(err)
			return
		}

		games, err := s.service.GetGames(r.Context(), jamURL)
		if err != nil {
			log.Error(err)
			return
		}

		pageData := templates.NewJamEntriesPageData(*jam, games)
		err = s.tmpl.JamEntriesTemplate.ExecuteTemplate(w, "base", pageData)
		if err != nil {
			log.Error(err)
			return
		}
	}
}
