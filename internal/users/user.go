package users

import "time"

type User struct {
	ID       uint64
	Login    string
	Email    string
	Password string
	Karma    int

	CreatedAt time.Time
	UpdatedAt time.Time
}
