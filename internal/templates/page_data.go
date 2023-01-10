package templates

import (
	"fmt"
	"html/template"
	"net/http"

	"GameJamPlatform/internal/gamejam"
)

type ErrorPageData struct {
	ErrorMessage string
}

func NewErrorPageData(errorCode int) ErrorPageData {
	pageData := ErrorPageData{}

	pageData.ErrorMessage = fmt.Sprintf("Error %d. %s", errorCode, http.StatusText(errorCode))

	return pageData
}

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

type JamEditPageData struct {
	Jam gamejam.GameJam

	StartDate     string
	EndDate       string
	VotingEndDate string
}

func NewJamEditPageData(jam gamejam.GameJam) JamEditPageData {
	return JamEditPageData{
		Jam:           jam,
		StartDate:     jam.StartDate.Format(TimeLayout),
		EndDate:       jam.EndDate.Format(TimeLayout),
		VotingEndDate: jam.VotingEndDate.Format(TimeLayout),
	}
}

type JamEntriesPageData struct {
	Jam   gamejam.GameJam
	Games []gamejam.Game
}

func NewJamEntriesPageData(jam gamejam.GameJam, games []gamejam.Game) JamEntriesPageData {
	return JamEntriesPageData{Jam: jam, Games: games}
}

type GameNewPageData struct {
	Jam gamejam.GameJam
}

func NewGameNewPageData(jam gamejam.GameJam) GameNewPageData {
	return GameNewPageData{Jam: jam}
}

type GameOverviewPageData struct {
	Jam  gamejam.GameJam
	Game gamejam.Game

	RenderedContent template.HTML
}

func NewGameOverviewPageData(jam gamejam.GameJam, game gamejam.Game) GameOverviewPageData {
	return GameOverviewPageData{Jam: jam, Game: game, RenderedContent: renderContent(game.Content)}
}

type GameEditPageData struct {
	Jam  gamejam.GameJam
	Game gamejam.Game
}

func NewGameEditPageData(jam gamejam.GameJam, game gamejam.Game) GameEditPageData {
	return GameEditPageData{Jam: jam, Game: game}
}
