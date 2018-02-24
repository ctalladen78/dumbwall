package database

import (
	"github.com/maksadbek/dumbwall/internal/platform/postgres"
	"github.com/maksadbek/dumbwall/internal/platform/redis"
	sq "github.com/masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Database struct {
	r *redis.Redis
	p *postgres.Postgres
}

func New() (*Database, error) {
	return &Database{}, nil
}
