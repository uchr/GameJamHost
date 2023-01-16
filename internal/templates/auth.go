package templates

import (
	"GameJamPlatform/internal/models/users"
)

type AuthPageData struct {
	IsAuth bool

	User *users.User
}

func NewAuthPageData(user *users.User) AuthPageData {
	return AuthPageData{
		IsAuth: user != nil,
		User:   user,
	}
}
