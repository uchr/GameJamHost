package servers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/services/validationerr"
	"GameJamPlatform/internal/web/defs"
	"GameJamPlatform/internal/web/pagedata"
)

func (s *server) parseJamForm(r *http.Request) (*gamejams.GameJam, error) {
	const maxUploadSize = 10 * 1024 * 1024 // 10 mb
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		return nil, err
	}

	vErr := validationerr.New()

	jam := gamejams.GameJam{
		Title:           r.FormValue("name"),
		URL:             r.FormValue("url"),
		Content:         r.FormValue("content"),
		HideResults:     r.FormValue("hide_results") == "on",
		HideSubmissions: r.FormValue("hide_submissions") == "on",
	}

	timezoneValue := r.FormValue("timezone")
	userTimezone, err := time.LoadLocation(timezoneValue)
	if err != nil {
		return nil, err
	}

	jam.StartDate, err = time.ParseInLocation(defs.TimeLayout, r.FormValue("start_date"), userTimezone)
	if err != nil {
		vErr.Add("StartDate", "Must be a valid date")
	}
	jam.EndDate, err = time.ParseInLocation(defs.TimeLayout, r.FormValue("end_date"), userTimezone)
	if err != nil {
		vErr.Add("EndDate", "Must be a valid date")
	}
	jam.VotingEndDate, err = time.ParseInLocation(defs.TimeLayout, r.FormValue("voting_end_date"), userTimezone)
	if err != nil {
		vErr.Add("VotingEndDate", "Must be a valid date")
	}

	jam.StartDate = jam.StartDate.In(time.UTC)
	jam.EndDate = jam.EndDate.In(time.UTC)
	jam.VotingEndDate = jam.VotingEndDate.In(time.UTC)

	coverImageURL, err := s.uploadImage(r, "cover_image")
	if err != nil {
		return nil, err
	}
	if coverImageURL != "" {
		jam.CoverImageURL = coverImageURL
	}

	criteriaTitleValues := r.Form["criteria_title[]"]
	criteriaDescValues := r.Form["criteria_desc[]"]
	if len(criteriaTitleValues) != len(criteriaDescValues) {
		return nil, errors.New("criteria and criteria_desc must be the same length")
	}
	for i := range criteriaTitleValues {
		jam.Criteria = append(jam.Criteria, gamejams.Criteria{
			Title:       criteriaTitleValues[i],
			Description: criteriaDescValues[i]},
		)
	}

	questionTitleValues := r.Form["question_title[]"]
	questionDescValues := r.Form["question_desc[]"]
	questionCriteriaValues := r.Form["question_criteria[]"]
	if len(questionTitleValues) != len(questionDescValues) || len(questionTitleValues) != len(questionCriteriaValues) {
		return nil, errors.New("question_title, question_desc, and question_criteria must be the same length")
	}
	for i := range questionTitleValues {
		jam.Questions = append(jam.Questions, gamejams.JamQuestion{
			Title:          questionTitleValues[i],
			Description:    questionDescValues[i],
			HiddenCriteria: questionCriteriaValues[i]},
		)
	}

	if vErr.HasErrors() {
		return &jam, vErr
	}

	return &jam, nil
}

func (s *server) jamsListHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const pageName = "jam_list"

		user := s.authedUser(r)

		jams, err := s.gameJams.GetJams(r.Context())
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
		user := s.authedUser(r)
		if user == nil {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jam, err := s.parseJamForm(r)
		if err != nil {
			var vErr validationerr.ValidationErrors
			if errors.As(err, &vErr) {
				s.redirectToValidatedJamForm(w, *user, *jam, true, &vErr)
				return
			}
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		err = s.gameJams.CreateJam(r.Context(), *user, *jam)
		if err != nil {
			var vErr validationerr.ValidationErrors
			if errors.As(err, &vErr) {
				s.redirectToValidatedJamForm(w, *user, *jam, true, &vErr)
				return
			}
			s.tm.RenderError(w, http.StatusInternalServerError, err)
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

		jam, err := s.gameJams.GetJamByURL(r.Context(), jamURL)
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

		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		jam, err := s.gameJams.GetJamByID(r.Context(), jamID)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		user := s.authedUser(r)
		isHost, err := s.gameJams.IsHost(r.Context(), *jam, user)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		if !isHost {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		pageData := pagedata.NewJamEditFormPageData(*user, *jam, false, nil)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) jamEditHandlerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		user := s.authedUser(r)
		isHost, err := s.gameJams.IsHostByID(r.Context(), jamID, user)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		if !isHost {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		jam, err := s.parseJamForm(r)
		if err != nil {
			var vErr validationerr.ValidationErrors
			if errors.As(err, &vErr) {
				s.redirectToValidatedJamForm(w, *user, *jam, false, &vErr)
				return
			}
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		err = s.gameJams.UpdateJam(r.Context(), jamID, *jam)
		if err != nil {
			var vErr validationerr.ValidationErrors
			if errors.As(err, &vErr) {
				s.redirectToValidatedJamForm(w, *user, *jam, true, &vErr)
				return
			}
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/jams/"+jam.URL, http.StatusSeeOther)
	}
}

func (s *server) jamDeleteHandlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jamIDText := chi.URLParam(r, "jamID")
		jamID, err := strconv.Atoi(jamIDText)
		if err != nil {
			s.tm.RenderError(w, http.StatusBadRequest, err)
			return
		}

		user := s.authedUser(r)
		isHost, err := s.gameJams.IsHostByID(r.Context(), jamID, user)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}
		if !isHost {
			s.tm.RenderError(w, http.StatusUnauthorized, nil)
			return
		}

		err = s.gameJams.DeleteJam(r.Context(), jamID)
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

		jam, err := s.gameJams.GetJamByURL(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusNotFound, err)
			return
		}

		games, err := s.gameJams.GetGames(r.Context(), jamURL)
		if err != nil {
			s.tm.RenderError(w, http.StatusInternalServerError, err)
			return
		}

		pageData := pagedata.NewJamEntriesPageData(user, *jam, games)
		s.tm.Render(w, pageName, pageData)
	}
}

func (s *server) redirectToValidatedJamForm(w http.ResponseWriter, user users.User, jam gamejams.GameJam, isNewJam bool, vErr *validationerr.ValidationErrors) {
	const pageName = "jam_edit_form"

	pageData := pagedata.NewJamEditFormPageData(user, jam, isNewJam, vErr)
	s.tm.Render(w, pageName, pageData)
}
