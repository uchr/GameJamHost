package templates

import (
	"html/template"

	"GameJamPlatform/internal/forms"
	"GameJamPlatform/internal/gamejam"
)

type JamListPageData struct {
	Jams []gamejam.GameJam
}

func NewJamListPageData(jams []gamejam.GameJam) JamListPageData {
	return JamListPageData{Jams: jams}
}

type JamOverviewPageData struct {
	Jam gamejam.GameJam

	RenderedContent template.HTML
}

func NewJamOverviewPageData(jam gamejam.GameJam) JamOverviewPageData {
	return JamOverviewPageData{Jam: jam, RenderedContent: renderContent(jam.Content)}
}

type JamEditFormPageData struct {
	IsNewJam bool

	Jam    gamejam.GameJam
	Errors forms.ValidationErrors

	StartDate     string
	EndDate       string
	VotingEndDate string
}

func NewJamEditFormPageData(isNewJam bool, jam gamejam.GameJam, validationErrors forms.ValidationErrors) JamEditFormPageData {
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
	Jam   gamejam.GameJam
	Games []gamejam.Game
}

func NewJamEntriesPageData(jam gamejam.GameJam, games []gamejam.Game) JamEntriesPageData {
	return JamEntriesPageData{Jam: jam, Games: games}
}
