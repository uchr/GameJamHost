package pagedata

import (
	"html/template"

	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/web/forms"
)

type UserProfilePageData struct {
	AuthPageData

	ProfileUser users.User
	Jams        []gamejams.GameJam
	Games       []gamejams.Game

	RenderedAbout template.HTML
}

func NewUserProfilePageData(authedUser *users.User, profileUser users.User, jams []gamejams.GameJam, games []gamejams.Game) *UserProfilePageData {
	return &UserProfilePageData{
		AuthPageData: NewAuthPageData(authedUser),

		ProfileUser: profileUser,
		Jams:        jams,
		Games:       games,

		RenderedAbout: renderContent(profileUser.About),
	}
}

type UserEditFormPageData struct {
	AuthPageData

	Errors forms.ValidationErrors
}

func NewUserEditFormPageData(user users.User, validationErrors forms.ValidationErrors) UserEditFormPageData {
	return UserEditFormPageData{
		AuthPageData: NewAuthPageData(&user),

		Errors: validationErrors,
	}
}
