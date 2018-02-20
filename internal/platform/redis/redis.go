package redis

import (
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

type Redis struct {
	pool *redigo.Pool
}

func New() *Redis {
	return &Redis{}
}

func (r *Redis) Hot(b, e uint64) ([]string, error) {
	c := r.pool.Get()
	defer c.Close()

	c.

	return []string{}, nil
}

func (r *Redis) PutHot(id string, score int) error {
	return nil
}

func (r *Redis) Best(b, e uint64) ([]string, error) {
	return []string{}, nil
}

func (r *Redis) PutBest(id string, score int) error {
	return nil
}

func (r *Redis) Top(b, e uint64) ([]string, error) {
	return []string{}, nil
}

func (r *Redis) PutTop(id string, score int) error {
	return nil
}

func (r *Redis) Controversial(b, e uint64) ([]string, error) {
	return []string{}, nil
}

func (r *Redis) PutControversial(id string, score int) error {
	return nil
}

func (r *Redis) New(b, e uint64) ([]string, error) {
	return []string{}, nil
}

func (r *Redis) PutNew(id string, createdAt int32) error {
	return nil
}

func (r *Redis) Rising(b, e uint64) ([]string, error) {
	return []string{}, nil
}

func (r *Redis) PutRising(id string, score int32) error {
	return nil
}
