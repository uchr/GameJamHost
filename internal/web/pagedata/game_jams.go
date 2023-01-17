package pagedata

import (
	"html/template"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/web/forms"
)

type JamListPageData struct {
	AuthPageData

	Jams []gamejams.GameJam
}

func NewJamListPageData(user *users.User, jams []gamejams.GameJam) JamListPageData {
	return JamListPageData{
		AuthPageData: NewAuthPageData(user),

		Jams: jams,
	}
}

type JamOverviewPageData struct {
	AuthPageData

	Jam gamejams.GameJam

	RenderedContent template.HTML
}

func NewJamOverviewPageData(users *users.User, jam gamejams.GameJam) JamOverviewPageData {
	return JamOverviewPageData{
		AuthPageData: NewAuthPageData(users),

		Jam:             jam,
		RenderedContent: renderContent(jam.Content),
	}
}

type JamEditFormPageData struct {
	AuthPageData

	IsNewJam bool

	Jam    gamejams.GameJam
	Errors forms.ValidationErrors

	StartDate     string
	EndDate       string
	VotingEndDate string
}

func NewJamEditFormPageData(user users.User, jam gamejams.GameJam, isNewJam bool, validationErrors forms.ValidationErrors) JamEditFormPageData {
	pageData := JamEditFormPageData{
		AuthPageData: NewAuthPageData(&user),

		IsNewJam: isNewJam,
		Jam:      jam,
		Errors:   validationErrors,
	}

	if !jam.StartDate.IsZero() {
		pageData.StartDate = jam.StartDate.Format(forms.TimeLayout)
	}
	if !jam.EndDate.IsZero() {
		pageData.EndDate = jam.EndDate.Format(forms.TimeLayout)
	}
	if !jam.VotingEndDate.IsZero() {
		pageData.VotingEndDate = jam.VotingEndDate.Format(forms.TimeLayout)
	}

	return pageData
}

type JamEntriesPageData struct {
	AuthPageData

	Jam   gamejams.GameJam
	Games []gamejams.Game
}

func NewJamEntriesPageData(users *users.User, jam gamejams.GameJam, games []gamejams.Game) JamEntriesPageData {
	return JamEntriesPageData{
		AuthPageData: NewAuthPageData(users),

		Jam:   jam,
		Games: games,
	}
}
