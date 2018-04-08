package database

import (
	"sort"

	"github.com/maksadbek/dumbwall/internal/posts"
)

// keys for sorted sets that keep listings in Redis
const (
	keyHot           = "hot"
	keyNewest        = "newest"
	keyTop           = "top"
	keyBest          = "best"
	keyControversial = "controversial"
	keyRising        = "rising"
)

func (d *Database) Hot(b, e int) ([]string, error) {
	return d.r.Rangess(keyHot, b, e)
}

func (d *Database) PutHot(id int, score int64) error {
	return d.r.Addss(keyHot, id, score)
}

func (d *Database) Best(b, e int) ([]string, error) {
	return d.r.Rangess(keyBest, b, e)
}

func (d *Database) PutBest(id int, score int64) error {
	return d.r.Addss(keyBest, id, score)
}

func (d *Database) Top(b, e int) ([]posts.Post, []error) {
	ids, err := d.r.Rangess(keyTop, int(b), int(e))
	if err != nil {
		return nil, []error{err}
	}

	posts, errs := d.GetPosts(ids)
	if len(errs) > 0 {
		return posts, errs
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Ups > posts[j].Ups
	})

	return posts, errs
}

func (d *Database) PutTop(id int, score int64) error {
	return d.r.Addss(keyTop, id, score)
}

func (d *Database) Controversial(b, e int) ([]string, error) {
	return d.r.Rangess(keyControversial, b, e)
}

func (d *Database) PutControversial(id int, score int64) error {
	return d.r.Addss(keyControversial, id, score)
}

func (d *Database) Newest(b, e int) ([]posts.Post, []error) {
	ids, err := d.r.Rangess(keyNewest, b, e)
	if err != nil {
		return nil, []error{err}
	}

	posts, errs := d.GetPosts(ids)
	if len(errs) > 0 {
		return posts, errs
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Unix() > posts[j].CreatedAt.Unix()
	})

	return posts, errs
}

func (d *Database) PutNew(id int, createdAt int64) error {
	return d.r.Addss(keyNewest, id, createdAt)
}

func (d *Database) Rising(b, e int) ([]string, error) {
	return d.r.Rangess(keyRising, b, e)
}

func (d *Database) PutRising(id int, score int64) error {
	return d.r.Addss(keyRising, id, score)
}
