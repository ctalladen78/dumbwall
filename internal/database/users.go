package database

import (
	"github.com/maksadbek/dumbwall/internal/users"
	sq "github.com/masterminds/squirrel"
)

func (d *Database) CreateUser(u users.User) (users.User, error) {
	var id uint64

	err := psql.Insert("users").
		Columns("login", "email", "password").
		Values(u.Login, u.Email, u.Password).
		Suffix("returning id").
		RunWith(d.p.DB).
		QueryRow().
		Scan(&id)

	if err != nil {
		return u, err
	}

	u.ID = id

	return u, nil
}

func (d *Database) GetUser(id uint64) (users.User, error) {
	var u users.User

	err := psql.
		Select(
			"login",
			"email",
			"created_at",
			"updated_at",
			"karma",
		).
		From("users").
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		QueryRow().
		Scan(
			&u.Login,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.Karma,
		)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (d *Database) UpdateUser(id uint64, u users.User) error {
	_, err := psql.Update("users").
		SetMap(map[string]interface{}{
			"login":      u.Login,
			"updated_at": "now()",
		}).
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		Exec()

	return err
}

func (d *Database) ChangeKarma(id uint64, delta int) error {
	_, err := psql.Update("users").
		Set("karma", delta).
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		Exec()

	return err
}

func (d *Database) ConfirmUserEmail(id uint32) error {
	_, err := psql.Update("users").
		Set("email_verified", true).
		Where(sq.Eq{"id": id}).
		RunWith(d.p.DB).
		Exec()

	return err
}

func (d *Database) CheckLogin(login string) error {
	return nil
}

func (d *Database) CheckEmail(email string) error {
	return nil
}
