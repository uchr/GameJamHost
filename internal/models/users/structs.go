package users

type User struct {
	ID int

	Email    string
	Username string
	Password []byte

	AvatarURL string
	About     string
}
