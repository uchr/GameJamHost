package servers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/templates"
)

func (s *server) parseGameForm(r *http.Request) (*gamejams.Game, error) {
	const maxUploadSize = 10 * 1024 * 1024 // 10 mb
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		return nil, err
	}

	game := gamejams.Game{
		Title:   r.FormValue("name"),
		Build:   r.FormValue("build"),
		Content: r.FormValue("content"),
	}

	coverImageURL, err := s.uploadImage(r, "CoverImage")
	if err != nil {
		return nil, err
	}
	if coverImageURL != "" {
		game.CoverImageURL = coverImageURL
	}

	for i := 0; i < 3; i++ {
		imageURL, err := s.uploadImage(r, fmt.Sprintf("Screenshot-%d", i))
		if err != nil {
			return nil, err
		}
		if imageURL != "" {
			game.ScreenshotURLs = append(game.ScreenshotURLs, imageURL)
		}
	}

	return &game, nil
}

func (s *server) gameNewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusNotFound, err)
			return
		}

		pageData := templates.NewGameEditFormPageData(true, *jam, gamejams.Game{}, nil)
		if err := s.tmpl.GameEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) gameCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")

		game, err := s.parseGameForm(r)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusBadRequest, err)
			return
		}

		gameURL, validationErrors, err := s.service.CreateGame(r.Context(), jamURL, *game)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
		if len(validationErrors) > 0 {
			jam, err := s.service.GetJamByURL(r.Context(), jamURL)
			if err != nil {
				s.executeErrorPage(w, r, http.StatusNotFound, err)
				return
			}

			pageData := templates.NewGameEditFormPageData(true, *jam, *game, validationErrors)
			if err := s.tmpl.JamEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
				s.executeErrorPage(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}

		http.Redirect(w, r, "/jams/"+jamURL+"/games/"+gameURL, http.StatusSeeOther)
	}
}

func (s *server) gameOverviewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusNotFound, err)
			return
		}

		game, err := s.service.GetGame(r.Context(), jamURL, gameURL)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusNotFound, err)
			return
		}

		pageData := templates.NewGameOverviewPageData(*jam, *game)
		if err := s.tmpl.GameOverviewTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
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
			s.executeErrorPage(w, r, http.StatusNotFound, err)
			return
		}

		game, err := s.service.GetGame(r.Context(), jamURL, gameURL)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusNotFound, err)
			return
		}

		pageData := templates.NewGameEditFormPageData(false, *jam, *game, nil)
		if err := s.tmpl.GameEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) gameUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		game, err := s.parseGameForm(r)

		validationErrors, err := s.service.UpdateGame(r.Context(), jamURL, gameURL, *game)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
		if len(validationErrors) > 0 {
			jam, err := s.service.GetJamByURL(r.Context(), jamURL)
			if err != nil {
				s.executeErrorPage(w, r, http.StatusNotFound, err)
				return
			}

			pageData := templates.NewGameEditFormPageData(false, *jam, *game, validationErrors)
			if err := s.tmpl.GameEditFormTemplate.ExecuteTemplate(w, "base", pageData); err != nil {
				s.executeErrorPage(w, r, http.StatusInternalServerError, err)
				return
			}
			return
		}

		http.Redirect(w, r, "/jams/"+jamURL+"/games/"+gameURL, http.StatusSeeOther)
	}
}

func (s *server) gameBanHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		err := s.service.BanGame(r.Context(), jamURL, gameURL)
		if err != nil {
			s.executeErrorPage(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}
