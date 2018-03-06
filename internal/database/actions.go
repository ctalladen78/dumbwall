package database

import (
	"database/sql"
	"errors"

	"github.com/maksadbek/dumbwall/internal/actions"
	sq "github.com/masterminds/squirrel"
)

func (d *Database) CreateAction(tx *sql.Tx, userID, postID int, actionType actions.Action) error {
	res, err := psql.
		Insert("actions").
		Columns("user_id", "post_id", "action_type").
		Values(userID, postID, actionType).
		RunWith(tx).
		Exec()

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected <= 0 {
		return errors.New("no rows were affected")
	}

	return nil
}

func (d *Database) DeleteAction(tx *sql.Tx, id int) error {
	_, err := psql.
		Delete("actions").
		Where(sq.Eq{"id": id}).
		RunWith(tx).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

func (d *Database) CheckAction(tx *sql.Tx, userID, postID int, actionType actions.Action) (bool, int, error) {
	var id int

	err := psql.
		Select("id").
		From("actions").
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.Eq{"post_id": postID},
			sq.Eq{"action_type": actionType},
		}).
		RunWith(tx).
		QueryRow().
		Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, id, nil
		}

		return false, id, err
	}

	return true, id, err
}
