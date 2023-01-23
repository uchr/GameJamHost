package servers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/services/validationerr"
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

	coverImageURL, err := s.uploadImage(r, "cover_image")
	if err != nil {
		return nil, err
	}
	if coverImageURL != "" {
		game.CoverImageURL = coverImageURL
	}

	for i := 0; i < 3; i++ {
		imageURL, err := s.uploadImage(r, fmt.Sprintf("screenshot_%d", i))
		if err != nil {
			return nil, err
		}
		if imageURL != "" {
			game.ScreenshotURLs = append(game.ScreenshotURLs, imageURL)
		}
	}

	answerValues := r.Form["answers[]"]
	for _, a := range answerValues {
		game.Answers = append(game.Answers, gamejams.GameAnswer{Answer: a == "Yes"})
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

		jam, err := s.gameJams.GetJamByURL(r.Context(), jamURL)
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

		gameURL, err := s.gameJams.CreateGame(r.Context(), jamURL, *user, *game)
		if err == nil {
			http.Redirect(w, r, "/jams/"+jamURL+"/games/"+gameURL, http.StatusSeeOther)
			return
		}

		var vErr validationerr.ValidationErrors
		if !errors.As(err, &vErr) {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		s.redirectToValidatedGameForm(w, r, jamURL, *user, *game, true, &vErr)
		return
	}
}

func (s *server) gameOverviewHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "game_overview"

		user := s.authedUser(r)

		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		jam, err := s.gameJams.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		game, err := s.gameJams.GetGame(r.Context(), jamURL, gameURL)
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

		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		jam, err := s.gameJams.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		game, err := s.gameJams.GetGame(r.Context(), jamURL, gameURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		user := s.authedUser(r)
		isOwner, err := s.gameJams.IsGameOwner(r.Context(), *game, user)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		if !isOwner {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		pageData := pagedata.NewGameEditFormPageData(*user, *jam, *game, false, nil)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) gameEditHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		user := s.authedUser(r)
		isOwner, err := s.gameJams.IsGameOwnerByURL(r.Context(), jamURL, gameURL, user)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		if !isOwner {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		game, err := s.parseGameForm(r)

		err = s.gameJams.UpdateGame(r.Context(), jamURL, gameURL, *game)
		if err == nil {
			http.Redirect(w, r, "/jams/"+jamURL+"/games/"+gameURL, http.StatusSeeOther)
			return
		}

		var vErr validationerr.ValidationErrors
		if !errors.As(err, &vErr) {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		s.redirectToValidatedGameForm(w, r, jamURL, *user, *game, true, &vErr)
		return
	}
}

func (s *server) gameBanHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		err := s.gameJams.BanGame(r.Context(), jamURL, gameURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/jams/"+jamURL+"/entries", http.StatusSeeOther)
	}
}

func (s *server) redirectToValidatedGameForm(w http.ResponseWriter, r *http.Request, jamURL string, user users.User, game gamejams.Game, isNewGame bool, vErr *validationerr.ValidationErrors) {
	const pageName = "game_edit_form"

	jam, err := s.gameJams.GetJamByURL(r.Context(), jamURL)
	if err != nil {
		s.tm.RenderError(w, http.StatusNotFound, err)
		return
	}

	pageData := pagedata.NewGameEditFormPageData(user, *jam, game, isNewGame, vErr)
	s.tm.Render(w, pageName, pageData)
}
