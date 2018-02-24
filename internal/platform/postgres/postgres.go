package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func New(dsn string) (*Postgres, error) {
	db, err := sql.Open("postgresql", dsn)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		DB: db,
	}, nil
}
