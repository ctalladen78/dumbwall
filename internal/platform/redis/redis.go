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

func (r *Redis) Zrevrange(key string, b, e int) (reply []string, err error) {
	c := r.pool.Get()
	defer c.Close()

	return redigo.Strings(c.Do("ZREVRANGE", key, b, e))
}

func (r *Redis) Zadd(key string, member int, score int64) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("ZADD", key, score, member)
	return err
}

func (r *Redis) Do(commandName string, args ...interface{}) (interface{}, error) {
	c := r.pool.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

func (r *Redis) Conn() redigo.Conn {
	return r.pool.Get()
}
