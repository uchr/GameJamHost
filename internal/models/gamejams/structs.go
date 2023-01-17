package gamejams

import "time"

type GameJam struct {
	ID     int
	UserID int // TODO: Multiple hosts

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
	UserID    int // TODO: Multiple authors (team)

	Title   string
	URL     string
	Content string
	Build   string

	IsBanned bool

	CoverImageURL  string
	ScreenshotURLs []string
}

type JamHost struct {
	ID        int
	GameJamID int
	UserID    int
}

type JamHostInvite struct {
	UID       string
	GameJamID int
}
