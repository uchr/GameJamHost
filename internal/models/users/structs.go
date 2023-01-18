package users

type User struct {
	ID int

	Email    string
	Username string
	Password []byte

	AvatarURL string
	About     string
}

type Participant struct {
	ID int

	UserID    int
	GameJamID int
	TeamID    int

	IsLookingForTeam bool
	Tags             []string

	IsAdmin bool
}
