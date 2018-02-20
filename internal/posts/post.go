package post

type PostType uint8

const (
	_ PostType = iota
	PostTypeText
	PostTypeLink
)

type Post struct {
	ID string

	Type PostType

	Title string
	Body  string

	Ups   uint64
	Downs int64
}
