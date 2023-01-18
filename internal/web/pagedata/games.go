package pagedata

import (
	"html/template"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/services/validationerr"
)

type GameEditFormPageData struct {
	AuthPageData

	Jam  gamejams.GameJam
	Game gamejams.Game

	IsNewGame bool

	Errors map[string]string
}

func NewGameEditFormPageData(user users.User, jam gamejams.GameJam, game gamejams.Game, isNewGame bool, vErr *validationerr.ValidationErrors) GameEditFormPageData {
	return GameEditFormPageData{
		AuthPageData: NewAuthPageData(&user),

		Jam:       jam,
		Game:      game,
		IsNewGame: isNewGame,
		Errors:    vErr.Errors(),
	}
}

type GameOverviewPageData struct {
	AuthPageData

	Jam  gamejams.GameJam
	Game gamejams.Game

	RenderedContent template.HTML
}

func NewGameOverviewPageData(user *users.User, jam gamejams.GameJam, game gamejams.Game) GameOverviewPageData {
	return GameOverviewPageData{
		AuthPageData: NewAuthPageData(user),

		Jam:             jam,
		Game:            game,
		RenderedContent: renderContent(game.Content),
	}
}
