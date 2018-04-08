package posts

import (
	"fmt"
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

func (p Post) FormattedCreatedAt() string {
	sub := time.Now().Sub(p.CreatedAt)

	if sub.Minutes() < 1 {
		if sub.Seconds() < 1 {
			return "a moment ago"
		}
		return fmt.Sprintf("%.0f seconds ago", sub.Seconds())
	}

	if sub.Minutes() < 60 {
		return fmt.Sprintf("%.0f minutes ago", sub.Minutes())
	}

	if sub.Hours() < 24 {
		return fmt.Sprintf("%.0f hours ago", sub.Hours())
	}

	if sub.Hours() < 48 {
		return "yesterday"
	}

	return p.CreatedAt.Format("Mar 7 2015")
}
