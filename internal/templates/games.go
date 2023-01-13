package templates

import (
	"html/template"

	"GameJamPlatform/internal/forms"
	"GameJamPlatform/internal/models"
)

type GameEditFormPageData struct {
	IsNewGame bool

	Jam  models.GameJam
	Game models.Game

	Errors forms.ValidationErrors
}

func NewGameEditFormPageData(isNewGame bool, jam models.GameJam, game models.Game, validationErrors forms.ValidationErrors) GameEditFormPageData {
	return GameEditFormPageData{
		IsNewGame: isNewGame,
		Jam:       jam,
		Game:      game,
		Errors:    validationErrors}
}

type GameOverviewPageData struct {
	Jam  models.GameJam
	Game models.Game

	RenderedContent template.HTML
}

func NewGameOverviewPageData(jam models.GameJam, game models.Game) GameOverviewPageData {
	return GameOverviewPageData{Jam: jam, Game: game, RenderedContent: renderContent(game.Content)}
}
