package gamejams

import "time"

type GameJam struct {
	ID int

	Title   string
	URL     string
	Content string

	StartDate     time.Time
	EndDate       time.Time
	VotingEndDate time.Time

	HideResults     bool
	HideSubmissions bool

	CoverImageURL string
}

type Game struct {
	ID        int
	GameJamID int

	Title   string
	URL     string
	Content string
	Build   string

	IsBanned bool

	CoverImageURL  string
	ScreenshotURLs []string
}
