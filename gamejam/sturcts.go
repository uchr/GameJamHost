package gamejam

import "time"

type GameJam struct {
	UID string

	Name  string
	Theme string

	StartDate   time.Time
	EndDate     time.Time
	VoteEndDate time.Time

	Url         string
	Description string

	HideResults     bool
	HideSubmissions bool
}

type Game struct {
	UID string

	Name        string
	Description string
	Url         string

	Author string
}
