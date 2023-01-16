package templates

import (
	"GameJamPlatform/internal/models/users"
)

func NewIndexPageData(user *users.User) AuthPageData {
	return NewAuthPageData(user)
}
