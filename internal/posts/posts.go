package posts

import "time"

type PostType uint8

const (
	_ PostType = iota
	PostTypeText
	PostTypeLink
)

type Post struct {
	ID uint64

	Type PostType

	Title string
	Body  string

	Ups   int64
	Downs int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Post) Validate() []error {
	return nil
}
