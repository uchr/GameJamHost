package servers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/gamejam"
	"GameJamPlatform/internal/log"
	"GameJamPlatform/internal/templates"
)

func (s *server) parseGameForm(r *http.Request) (*gamejam.Game, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	game := gamejam.Game{
		Name:    r.FormValue("name"),
		Build:   r.FormValue("build"),
		Content: r.FormValue("content"),
	}

	return &game, nil
}

func (s *server) gameNewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			log.Error(err)
			return
		}

		pageData := templates.NewGameNewPageData(*jam)
		err = s.tmpl.GameNewTemplate.ExecuteTemplate(w, "base", pageData)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *server) gameOverviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			log.Error(err)
			return
		}

		game, err := s.service.GetGame(r.Context(), jamURL, gameURL)
		if err != nil {
			log.Error(err)
			return
		}

		pageData := templates.NewGameOverviewPageData(*jam, *game)
		err = s.tmpl.GameOverviewTemplate.ExecuteTemplate(w, "base", pageData)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *server) gameEditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			log.Error(err)
			return
		}

		game, err := s.service.GetGame(r.Context(), jamURL, gameURL)
		if err != nil {
			log.Error(err)
			return
		}

		pageData := templates.NewGameEditPageData(*jam, *game)
		err = s.tmpl.GameEditTemplate.ExecuteTemplate(w, "base", pageData)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *server) gameBanHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		err := s.service.BanGame(r.Context(), jamURL, gameURL)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *server) gameCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")

		game, err := s.parseGameForm(r)
		if err != nil {
			log.Error(err)
			return
		}

		gameURL, err := s.service.CreateGame(r.Context(), jamURL, *game)
		if err != nil {
			log.Error(err)
			return
		}

		http.Redirect(w, r, "/jams/"+jamURL+"/games/"+gameURL, http.StatusSeeOther)
	}
}

func (s *server) gameUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		game, err := s.parseGameForm(r)

		err = s.service.UpdateGame(r.Context(), jamURL, gameURL, *game)
		if err != nil {
			log.Error(err)
			return
		}

		http.Redirect(w, r, "/jams/"+jamURL+"/games/"+gameURL, http.StatusSeeOther)
	}
}
