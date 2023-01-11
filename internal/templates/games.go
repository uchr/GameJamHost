package templates

import (
	"html/template"

	"GameJamPlatform/internal/forms"
	"GameJamPlatform/internal/gamejam"
)

type GameEditFormPageData struct {
	IsNewGame bool

	Jam  gamejam.GameJam
	Game gamejam.Game

	Errors forms.ValidationErrors
}

func NewGameEditFormPageData(isNewGame bool, jam gamejam.GameJam, game gamejam.Game, validationErrors forms.ValidationErrors) GameEditFormPageData {
	return GameEditFormPageData{
		IsNewGame: isNewGame,
		Jam:       jam,
		Game:      game,
		Errors:    validationErrors}
}

type GameOverviewPageData struct {
	Jam  gamejam.GameJam
	Game gamejam.Game

	RenderedContent template.HTML
}

func NewGameOverviewPageData(jam gamejam.GameJam, game gamejam.Game) GameOverviewPageData {
	return GameOverviewPageData{Jam: jam, Game: game, RenderedContent: renderContent(game.Content)}
}
