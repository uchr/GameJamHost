package templates

import (
	"html/template"

	"GameJamPlatform/internal/forms"
	"GameJamPlatform/internal/models/gamejams"
)

type GameEditFormPageData struct {
	AuthPageData
	IsNewGame bool

	Jam  gamejams.GameJam
	Game gamejams.Game

	Errors forms.ValidationErrors
}

func NewGameEditFormPageData(isNewGame bool, jam gamejams.GameJam, game gamejams.Game, validationErrors forms.ValidationErrors) GameEditFormPageData {
	return GameEditFormPageData{
		IsNewGame: isNewGame,
		Jam:       jam,
		Game:      game,
		Errors:    validationErrors}
}

type GameOverviewPageData struct {
	AuthPageData
	Jam  gamejams.GameJam
	Game gamejams.Game

	RenderedContent template.HTML
}

func NewGameOverviewPageData(jam gamejams.GameJam, game gamejams.Game) GameOverviewPageData {
	return GameOverviewPageData{Jam: jam, Game: game, RenderedContent: renderContent(game.Content)}
}
