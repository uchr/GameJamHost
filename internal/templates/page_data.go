package templates

import (
	"fmt"
	"net/http"
)

type ErrorPageData struct {
	ErrorMessage      string
	DebugErrorMessage string
}

func NewErrorPageData(errorCode int, err error) ErrorPageData {
	return ErrorPageData{
		ErrorMessage:      fmt.Sprintf("Error %d. %s", errorCode, http.StatusText(errorCode)),
		DebugErrorMessage: err.Error(),
	}
}
