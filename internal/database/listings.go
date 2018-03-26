package database

func (d *Database) Hot(b, e int) ([]string, error) {
	return d.r.Zrevrange("hot", b, e)
}

func (d *Database) PutHot(id int, score int64) error {
	return d.r.Zadd("hot", id, score)
}

func (d *Database) Best(b, e int) ([]string, error) {
	return d.r.Zrevrange("best", b, e)
}

func (d *Database) PutBest(id int, score int64) error {
	return d.r.Zadd("best", id, score)
}

func (d *Database) Top(b, e int) ([]string, error) {
	return d.r.Zrevrange("top", int(b), int(e))
}

func (d *Database) PutTop(id int, score int64) error {
	return d.r.Zadd("top", id, score)
}

func (d *Database) Controversial(b, e int) ([]string, error) {
	return d.r.Zrevrange("controversial", b, e)
}

func (d *Database) PutControversial(id int, score int64) error {
	return d.r.Zadd("controversial", id, score)
}

func (d *Database) Newest(b, e int) ([]posts.Post, []error) {
	ids, err := d.r.Zrevrange("new", b, e)
	if err != nil {
		return nil, []error{err}
	}

	posts, errs := d.GetPosts(ids)
	if len(errs) > 0 {
		return posts, errs
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Unix() > posts[j].CreatedAt.Unix()
	})

	return posts, errs
}

func (d *Database) PutNew(id int, createdAt int64) error {
	return d.r.Zadd("new", id, createdAt)
}

func (d *Database) Rising(b, e int) ([]string, error) {
	return d.r.Zrevrange("rising", b, e)
}

func (d *Database) PutRising(id int, score int64) error {
	return d.r.Zadd("rising", id, score)
}
