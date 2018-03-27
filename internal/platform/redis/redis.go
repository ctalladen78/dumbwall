package redis

import (
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

// Redis is the pooled connection to Redis.
type Redis struct {
	pool *redigo.Pool
}

// New creates a pooled connection to Redis server.
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

// Rangess ranges over sorted set with zrevrange command of Redis
// receives key, beginning and end offset
// returns slice of values and error
func (r *Redis) Rangess(key string, b, e int) (reply []string, err error) {
	c := r.pool.Get()
	defer c.Close()

	return redigo.Strings(c.Do("ZREVRANGE", key, b, e))
}

// Addss adds record to sorted set using zadd command of Redis.
// receives key, member value and its score.
// returns an error.
func (r *Redis) Addss(key string, member int, score int64) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("ZADD", key, score, member)
	return err
}

// Do gets new connection from pool, runs given command with arguments and returns output.
func (r *Redis) Do(commandName string, args ...interface{}) (interface{}, error) {
	c := r.pool.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// Conn gets new connection from pool.
// received conn must be closed with Close() function.
func (r *Redis) Conn() redigo.Conn {
	return r.pool.Get()
}
