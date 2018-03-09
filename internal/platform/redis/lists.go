package redis

func (r *Redis) Hot(b, e int) ([]string, error) {
	return r.zrevrange("hot", b, e)
}

func (r *Redis) PutHot(id int, score int64) error {
	return r.zadd("hot", id, score)
}

func (r *Redis) Best(b, e int) ([]string, error) {
	return r.zrevrange("best", b, e)
}

func (r *Redis) PutBest(id int, score int64) error {
	return r.zadd("best", id, score)
}

func (r *Redis) Top(b, e int) ([]string, error) {
	return r.zrevrange("top", int(b), int(e))
}

func (r *Redis) PutTop(id int, score int64) error {
	return r.zadd("top", id, score)
}

func (r *Redis) Controversial(b, e int) ([]string, error) {
	return r.zrevrange("controversial", b, e)
}

func (r *Redis) PutControversial(id int, score int64) error {
	return r.zadd("controversial", id, score)
}

func (r *Redis) Newest(b, e int) ([]string, error) {
	return r.zrevrange("new", b, e)
}

func (r *Redis) PutNew(id int, createdAt int64) error {
	return r.zadd("new", id, createdAt)
}

func (r *Redis) Rising(b, e int) ([]string, error) {
	return r.zrevrange("rising", b, e)
}

func (r *Redis) PutRising(id int, score int64) error {
	return r.zadd("rising", id, score)
}
