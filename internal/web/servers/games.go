package servers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/web/pagedata"
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

func (s *server) gameNewHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "game_edit_form"

		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jamURL := chi.URLParam(r, "jamURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		pageData := pagedata.NewGameEditFormPageData(*user, *jam, gamejams.Game{}, true, nil)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) gameNewHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "game_edit_form"

		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jamURL := chi.URLParam(r, "jamURL")

		game, err := s.parseGameForm(r)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		gameURL, validationErrors, err := s.service.CreateGame(r.Context(), jamURL, *game)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		if len(validationErrors) == 0 {
			http.Redirect(w, r, "/jams/"+jamURL+"/games/"+gameURL, http.StatusSeeOther)
		}

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		pageData := pagedata.NewGameEditFormPageData(*user, *jam, *game, true, validationErrors)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) gameOverviewHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "game_overview"

		user := s.authedUser(r)

		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		game, err := s.service.GetGame(r.Context(), jamURL, gameURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		pageData := pagedata.NewGameOverviewPageData(user, *jam, *game)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) gameEditHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "game_edit_form"

		// TODO: check if user is owner of game
		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		game, err := s.service.GetGame(r.Context(), jamURL, gameURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		pageData := pagedata.NewGameEditFormPageData(*user, *jam, *game, false, nil)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) gameEditHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "game_edit_form"

		// TODO: check if user is owner of game
		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		game, err := s.parseGameForm(r)

		validationErrors, err := s.service.UpdateGame(r.Context(), jamURL, gameURL, *game)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		if len(validationErrors) == 0 {
			http.Redirect(w, r, "/jams/"+jamURL+"/games/"+gameURL, http.StatusSeeOther)
		}

		jam, err := s.service.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		pageData := pagedata.NewGameEditFormPageData(*user, *jam, *game, false, validationErrors)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) gameBanHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		err := s.service.BanGame(r.Context(), jamURL, gameURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/jams/"+jamURL+"/entries", http.StatusSeeOther)
	}
}
