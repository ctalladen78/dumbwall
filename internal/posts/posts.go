package posts

import (
	"time"

	"github.com/maksadbek/dumbwall/internal/actions"
)

type PostType uint8

const (
	_ PostType = iota
	PostTypeText
	PostTypeLink
)

type Post struct {
	ID int

	Type PostType

	Title string
	Body  string

	Ups   int64
	Downs int64

	CreatedAt time.Time
	UpdatedAt time.Time

	UserID int

	Meta Meta
}

type Meta struct {
	OwnerLogin string
	Action     actions.Action
}

func (p Post) Validate() []error {
	return nil
}
