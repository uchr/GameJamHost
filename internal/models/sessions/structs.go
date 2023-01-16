package sessions

import "time"

type Session struct {
	UID string

	UserID int

	ExpireAt time.Time
}
