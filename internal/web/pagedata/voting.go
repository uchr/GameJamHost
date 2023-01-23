package pagedata

import (
	"GameJamPlatform/internal/models/gamejams"
	"GameJamPlatform/internal/models/users"
	"GameJamPlatform/internal/models/voting"
)

type GameCriteriaPageData struct {
	Title string
	Score float32
}

type JamResultPageData struct {
	AuthPageData

	Jam     gamejams.GameJam
	Results []voting.CriteriaResult
}

func NewJamResultPageData(user *users.User, jam gamejams.GameJam, result voting.JamResult) JamResultPageData {
	return JamResultPageData{
		AuthPageData: NewAuthPageData(user),
		Jam:          jam,
		Results:      result.CriteriaResults,
	}
}
