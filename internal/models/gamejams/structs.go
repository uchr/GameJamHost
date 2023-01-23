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

	Criteria  []Criteria
	Questions []JamQuestion
}

type Game struct {
	ID     int
	JamID  int
	UserID int // TODO: Multiple authors (team)

	Title   string
	URL     string
	Content string
	Build   string

	IsBanned bool

	CoverImageURL  string
	ScreenshotURLs []string

	Answers []GameAnswer
}

type Criteria struct {
	ID    int
	JamID int

	Title       string
	Description string
}

type Vote struct {
	ID          int
	GameID      int
	UserID      int
	CriteriaUID string

	Value int
}

type JamQuestion struct {
	ID    int
	JamID int

	Title       string
	Description string

	HiddenCriteria string
}

type GameAnswer struct {
	ID         int
	GameID     int
	QuestionID int

	Answer bool
}

type JamState int

const (
	JamStateNotStarted JamState = iota
	JamStateStarted
	JamStateEnded
	JamStateVotingEnded
)
