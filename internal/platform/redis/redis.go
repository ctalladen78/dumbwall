package redis

import (
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

type Redis struct {
	pool *redigo.Pool
}

func New(addr string, maxIdle, maxActive, idleTimeout int) *Redis {
	return &Redis{
		pool: &redigo.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: time.Second * time.Duration(idleTimeout),
			Dial:        func() (redigo.Conn, error) { return redigo.Dial("tcp", addr) },
		},
	}
}

func (r *Redis) Hot(b, e uint64) ([]string, error) {
	return r.zrevrange("hot", b, e)
}

func (r *Redis) PutHot(id string, score uint64) error {
	return r.zadd("hot", id, score)
}

func (r *Redis) Best(b, e uint64) ([]string, error) {
	return r.zrevrange("best", b, e)
}

func (r *Redis) PutBest(id string, score uint64) error {
	return r.zadd("best", id, score)
}

func (r *Redis) Top(b, e uint64) ([]string, error) {
	return r.zrevrange("top", b, e)
}

func (r *Redis) PutTop(id string, score uint64) error {
	return r.zadd("put", id, score)
}

func (r *Redis) Controversial(b, e uint64) ([]string, error) {
	return r.zrevrange("controversial", b, e)
}

func (r *Redis) PutControversial(id string, score uint64) error {
	return r.zadd("controversial", id, score)
}

func (r *Redis) New(b, e uint64) ([]string, error) {
	return r.zrevrange("new", b, e)
}

func (r *Redis) PutNew(id string, createdAt int64) error {
	return r.zadd("new", id, uint64(createdAt))
}

func (r *Redis) Rising(b, e uint64) ([]string, error) {
	return r.zrevrange("rising", b, e)
}

func (r *Redis) PutRising(id string, score uint64) error {
	return r.zadd("rising", id, score)
}

func (r *Redis) zrevrange(key string, b, e uint64) (reply []string, err error) {
	c := r.pool.Get()
	defer c.Close()

	return redigo.Strings(c.Do("ZREVRANGE", key, b, e))
}

func (r *Redis) zadd(key string, member string, score uint64) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("ZADD", key, score, member)
	return err
}
