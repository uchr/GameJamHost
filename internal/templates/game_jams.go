package templates

import (
	"html/template"

	"GameJamPlatform/internal/forms"
	"GameJamPlatform/internal/models/gamejams"
)

type JamListPageData struct {
	AuthPageData

	Jams []gamejams.GameJam
}

func NewJamListPageData(jams []gamejams.GameJam) JamListPageData {
	return JamListPageData{Jams: jams}
}

type JamOverviewPageData struct {
	AuthPageData

	Jam gamejams.GameJam

	RenderedContent template.HTML
}

func NewJamOverviewPageData(jam gamejams.GameJam) JamOverviewPageData {
	return JamOverviewPageData{Jam: jam, RenderedContent: renderContent(jam.Content)}
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

func NewJamEditFormPageData(isNewJam bool, jam gamejams.GameJam, validationErrors forms.ValidationErrors) JamEditFormPageData {
	pageData := JamEditFormPageData{
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

func NewJamEntriesPageData(jam gamejams.GameJam, games []gamejams.Game) JamEntriesPageData {
	return JamEntriesPageData{Jam: jam, Games: games}
}
