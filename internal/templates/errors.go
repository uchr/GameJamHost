package templates

import (
	"fmt"
	"net/http"

	"GameJamPlatform/internal/models/users"
)

type ErrorPageData struct {
	AuthPageData

	ErrorMessage      string
	DebugErrorMessage string
}

func NewErrorPageData(user *users.User, errorCode int, err error) ErrorPageData {
	return ErrorPageData{
		AuthPageData:      NewAuthPageData(user),
		ErrorMessage:      fmt.Sprintf("Error %d. %s", errorCode, http.StatusText(errorCode)),
		DebugErrorMessage: err.Error(),
	}
}
