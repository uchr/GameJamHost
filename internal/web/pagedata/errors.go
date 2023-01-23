package pagedata

import (
	"fmt"
	"net/http"
)

type ErrorPageData struct {
	AuthPageData

	ErrorMessage      string
	DebugErrorMessage string
}

func NewErrorPageData(errorCode int, err error) ErrorPageData {
	pageData := ErrorPageData{
		ErrorMessage: fmt.Sprintf("Error %d. %s", errorCode, http.StatusText(errorCode)),
	}
	if err != nil {
		pageData.DebugErrorMessage = err.Error()
	}
	return pageData
}
