package servers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/web/pagedata"
)

func (s *server) parseVoteForm(r *http.Request, jam gamejams.GameJam, gameID, userID int) ([]gamejams.Vote, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	votes := make([]gamejams.Vote, 0)

	voteValues := r.Form["vote[]"]
	for i, v := range voteValues {
		if v == "" {
			return nil, errors.New("empty vote")
		}

		vote, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		votes = append(votes, gamejams.Vote{
			GameID:     gameID,
			UserID:     userID,
			CriteriaID: jam.Criteria[i].ID,
			Value:      vote,
		})
	}

	return votes, nil
}

func (s *server) gameVoteHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamURL := chi.URLParam(r, "jamURL")
		gameURL := chi.URLParam(r, "gameURL")

		user := s.authedUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		jam, err := s.gameJams.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		game, err := s.gameJams.GetGame(r.Context(), jamURL, gameURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		votes, err := s.parseVoteForm(r, *jam, game.ID, user.ID)

		err = s.gameJams.VoteGame(r.Context(), votes)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/jams/"+jamURL+"/games/"+gameURL, http.StatusSeeOther)
	}
}

func (s *server) jamResultsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "jam_results"

		user := s.authedUser(r)

		jamURL := chi.URLParam(r, "jamURL")

		jam, err := s.gameJams.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		jamResults, err := s.gameJams.GetJamResult(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		pageData := pagedata.NewJamResultPageData(user, *jam, *jamResults)
		s.tm.Render(w, pageName, pageData)
	}
}
