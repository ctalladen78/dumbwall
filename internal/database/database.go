package database

import (
	"github.com/maksadbek/dumbwall/internal/config"
	"github.com/maksadbek/dumbwall/internal/platform/postgres"
	"github.com/maksadbek/dumbwall/internal/platform/redis"
	sq "github.com/masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Database struct {
	r *redis.Redis
	p *postgres.Postgres
}

func New(c config.Database) (*Database, error) {
	p, err := postgres.New(c.Postgres.DSN)
	if err != nil {
		return nil, err
	}

	err = p.DB.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{
		r: redis.New(c.Redis.Addrs[0], 10, 10, 100),
		p: p,
	}, nil
}
