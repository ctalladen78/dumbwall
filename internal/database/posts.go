package database

import (
	"github.com/maksadbek/dumbwall/internal/posts"
	sq "github.com/masterminds/squirrel"
)

func (d *Database) CreatePost(userID int64, post posts.Post) (posts.Post, error) {
	var id uint64

	err := psql.Insert("posts").
		Columns("type", "title", "body", "user_id").
		Values(post.Type, post.Title, post.Body, userID).
		RunWith(d.p.DB).
		QueryRow().
		Scan(&id)

	if err != nil {
		return post, err
	}

	post.ID = id

	return post, nil
}

func (d *Database) GetPost(id string) (posts.Post, error) {
	var p posts.Post

	err := psql.Select("type", "title", "body", "ups", "downs").
		From("posts").
		RunWith(d.p.DB).
		QueryRow().
		Scan(&p.Type, &p.Title, &p.Body, &p.Ups, &p.Downs, &p.CreatedAt)

	if err != nil {
		return p, err
	}

	p.Ups -= p.Downs

	return p, nil
}

func (d *Database) UpdatePost(id string, p posts.Post) (posts.Post, error) {
	_, err := psql.Update("posts").
		SetMap(map[string]interface{}{
			"title":      p.Title,
			"body":       p.Body,
			"updated_at": "now()",
		}).
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		Exec()

	return p, err
}

func (d *Database) UpPost(id uint64) error {
	_, err := psql.Update("posts").
		Set("ups", "ups + 1").
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		Exec()

	return err
}

func (d *Database) DownPost(id uint64) error {
	_, err := psql.Update("posts").
		Set("downs", "down+1").
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		Exec()

	return err
}

func (d *Database) Delete(id uint64) error {
	_, err := psql.
		Delete("posts").
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		Exec()

	return err
}

func (d *Database) GetPosts(ids []uint64) ([]posts.Post, []error) {
	var (
		list = []posts.Post{}
		errs = []error{}
	)

	rows, err := psql.Select("type", "title", "body", "ups", "downs").
		From("posts").
		Where(sq.Eq{"id": ids}).
		RunWith(d.p.DB).Query()

	if err != nil {
		return list, append(errs, err)
	}

	for rows.Next() {
		var p posts.Post
		err := rows.Scan(&p.Type, &p.Title, &p.Body, &p.Ups, &p.Downs)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		p.Ups -= p.Downs

		list = append(list, p)
	}

	return list, errs
}
