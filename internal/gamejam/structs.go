package gamejam

import "time"

type GameJam struct {
	ID int

	Name    string
	URL     string
	Content string

	StartDate     time.Time
	EndDate       time.Time
	VotingEndDate time.Time

	HideResults     bool
	HideSubmissions bool

	// CoverImageUrl string // TODO: add cover image
}

type Game struct {
	ID        int
	GameJamID int

	Name    string
	URL     string
	Content string
	Build   string

	IsBanned bool

	// CoverImageUrl string // TODO: add cover image
	// Screenshots   []string // TODO: add screenshots
}
