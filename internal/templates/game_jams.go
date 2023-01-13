package templates

import (
	"html/template"

	"GameJamPlatform/internal/forms"
	"GameJamPlatform/internal/models"
)

type JamListPageData struct {
	Jams []models.GameJam
}

func NewJamListPageData(jams []models.GameJam) JamListPageData {
	return JamListPageData{Jams: jams}
}

type JamOverviewPageData struct {
	Jam models.GameJam

	RenderedContent template.HTML
}

func NewJamOverviewPageData(jam models.GameJam) JamOverviewPageData {
	return JamOverviewPageData{Jam: jam, RenderedContent: renderContent(jam.Content)}
}

type JamEditFormPageData struct {
	IsNewJam bool

	Jam    models.GameJam
	Errors forms.ValidationErrors

	StartDate     string
	EndDate       string
	VotingEndDate string
}

func NewJamEditFormPageData(isNewJam bool, jam models.GameJam, validationErrors forms.ValidationErrors) JamEditFormPageData {
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
	Jam   models.GameJam
	Games []models.Game
}

func NewJamEntriesPageData(jam models.GameJam, games []models.Game) JamEntriesPageData {
	return JamEntriesPageData{Jam: jam, Games: games}
}
