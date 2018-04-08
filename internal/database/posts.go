package database

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/maksadbek/dumbwall/internal/actions"
	"github.com/maksadbek/dumbwall/internal/posts"
	sq "github.com/masterminds/squirrel"
)

func (d *Database) CreatePost(userID int, post posts.Post) (posts.Post, error) {
	var id int
	var createdAt, updatedAt pq.NullTime

	err := psql.Insert("posts").
		Columns("type", "title", "body", "user_id").
		Values(post.Type, post.Title, post.Body, userID).
		Suffix("returning id,created_at,updated_at").
		RunWith(d.p.DB).
		QueryRow().
		Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return post, err
	}

	post.ID = id
	post.CreatedAt = createdAt.Time
	post.UpdatedAt = updatedAt.Time

	d.PutNew(post.ID, post.CreatedAt.Unix())
	return post, nil
}

func (d *Database) GetPost(id int) (posts.Post, error) {
	var p posts.Post

	var updatedAt pq.NullTime

	err := psql.
		Select(
			"posts.type",
			"posts.title",
			"posts.body",
			"posts.ups",
			"posts.downs",
			"posts.user_id",
			"posts.updated_at",
			"posts.created_at",
			"users.login").
		From("posts").
		Join("users on posts.user_id = users.id").
		RunWith(d.p.DB).
		Where(sq.Eq{"posts.id": id}).
		QueryRow().
		Scan(
			&p.Type,
			&p.Title,
			&p.Body,
			&p.Ups,
			&p.Downs,
			&p.UserID,
			&updatedAt,
			&p.CreatedAt,
			&p.Meta.OwnerLogin,
		)

	if err != nil {
		return p, err
	}

	p.UpdatedAt = updatedAt.Time
	p.Ups -= p.Downs

	return p, nil
}

func (d *Database) UpdatePost(id string, p posts.Post) (posts.Post, error) {
	_, err := psql.Update("posts").
		SetMap(map[string]interface{}{
			"title":      p.Title,
			"body":       p.Body,
			"updated_at": sq.Expr("now()"),
		}).
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		Exec()

	return p, err
}

func (d *Database) VotePost(userID, postID int, actionType actions.Action) error {
	tx, err := d.p.DB.Begin()
	if err != nil {
		return err
	}

	switch actionType {
	case actions.ActionUp:
		err = d.UpPost(tx, userID, postID)
	case actions.ActionDown:
		err = d.DownPost(tx, userID, postID)
	default:
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func (d *Database) UpPost(tx *sql.Tx, userID, postID int) error {
	var score int64

	done, actionID, err := d.CheckAction(tx, userID, postID, actions.ActionUp)
	if err != nil {
		return err
	}

	if done {
		err = d.DeleteAction(tx, actionID)
		if err != nil {
			return err
		}

		err = psql.Update("posts").
			Set("ups", sq.Expr("ups-1")).
			Where(sq.Eq{"id": postID}).
			Suffix("returning ups - downs").
			RunWith(tx).
			QueryRow().
			Scan(&score)

		if err != nil {
			return err
		}

		err = d.AddVote(userID, postID, actions.ActionNone)
		if err != nil {
			return err
		}
		return d.PutTop(postID, score)
	}

	err = d.CreateAction(tx, userID, postID, actions.ActionUp)
	if err != nil {
		return err
	}

	err = psql.Update("posts").
		Set("ups", sq.Expr("ups+1")).
		Where(sq.Eq{"id": postID}).
		Suffix("returning ups-downs").
		RunWith(tx).
		QueryRow().
		Scan(&score)

	if err != nil {
		return err
	}

	err = d.AddVote(userID, postID, actions.ActionUp)
	if err != nil {
		return err
	}

	return d.PutTop(postID, score)
}

func (d *Database) DownPost(tx *sql.Tx, userID, postID int) error {
	var score int64

	done, actionID, err := d.CheckAction(tx, userID, postID, actions.ActionDown)
	if err != nil {
		return err
	}

	if done {
		err = d.DeleteAction(tx, actionID)
		if err != nil {
			return err
		}

		err = psql.Update("posts").
			Set("downs", sq.Expr("downs-1")).
			Where(sq.Eq{"id": postID}).
			Suffix("returning ups - downs").
			RunWith(tx).
			QueryRow().
			Scan(&score)

		if err != nil {
			return err
		}

		err = d.AddVote(userID, postID, actions.ActionNone)
		if err != nil {
			return err
		}
		return d.PutTop(postID, score)
	}

	err = d.CreateAction(tx, userID, postID, actions.ActionDown)
	if err != nil {
		return err
	}

	err = psql.Update("posts").
		Set("downs", sq.Expr("downs+1")).
		Where(sq.Eq{"id": postID}).
		Suffix("returning ups - downs").
		RunWith(d.p.DB).
		QueryRow().
		Scan(&score)

	if err != nil {
		return err
	}

	err = d.AddVote(userID, postID, actions.ActionDown)
	if err != nil {
		return err
	}
	return d.PutTop(postID, score)
}

func (d *Database) Delete(id uint64) error {
	_, err := psql.
		Delete("posts").
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		Exec()

	return err
}

func (d *Database) GetPosts(ids []string) ([]posts.Post, []error) {
	var (
		list = []posts.Post{}
		errs = []error{}

		createdAt, updatedAt pq.NullTime
	)

	rows, err := psql.
		Select(
			"posts.id",
			"posts.type",
			"posts.title",
			"posts.body",
			"posts.ups",
			"posts.downs",
			"posts.user_id",
			"posts.created_at",
			"posts.updated_at",
			"users.login",
		).
		From("posts").
		Join("users on users.id = posts.user_id").
		Where(sq.Eq{"posts.id": ids}).
		RunWith(d.p.DB).Query()

	if err != nil {
		return list, append(errs, err)
	}

	for rows.Next() {
		var p posts.Post
		err := rows.Scan(
			&p.ID,
			&p.Type,
			&p.Title,
			&p.Body,
			&p.Ups,
			&p.Downs,
			&p.UserID,
			&createdAt,
			&updatedAt,
			&p.Meta.OwnerLogin)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		p.CreatedAt = createdAt.Time
		p.UpdatedAt = updatedAt.Time

		p.Ups -= p.Downs

		list = append(list, p)
	}

	return list, errs
}
