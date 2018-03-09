package redis

import (
	"fmt"

	redigo "github.com/garyburd/redigo/redis"
)

var uvotes = "uvotes_tmp:%d:%d"

func (r *Redis) AddVote(userID, postID, voteType int) error {
	c := r.pool.Get()
	defer c.Close()

	if voteType == 0 {
		println("deleting votes")
		c.Do("SREM", fmt.Sprintf(uvotes, userID, 1), postID)
		c.Do("SREM", fmt.Sprintf(uvotes, userID, 2), postID)
		return nil
	}

	key := fmt.Sprintf(uvotes, userID, voteType)

	_, err := c.Do("SADD", key, postID)

	return err
}

func (r *Redis) CheckVotes(userID int, postIDs ...int) ([]int, error) {
	c := r.pool.Get()
	defer c.Close()

	var voteTypes = make([]int, 0, len(postIDs))

	for i := range postIDs {
		err := c.Send("SISMEMBER", fmt.Sprintf(uvotes, userID, 1), postIDs[i])
		if err != nil {
			return voteTypes, err
		}
		err = c.Send("SISMEMBER", fmt.Sprintf(uvotes, userID, 2), postIDs[i])
		if err != nil {
			return voteTypes, err
		}
	}

	err := c.Flush()
	if err != nil {
		panic(err)
		return voteTypes, err
	}

	for i := 0; i < len(postIDs); i++ {
		voteType1, err := redigo.Int(c.Receive())
		if err != nil {
			return voteTypes, err
		}

		voteType2, err := redigo.Int(c.Receive())
		if err != nil {
			return voteTypes, err
		}

		if voteType1 == 1 {
			voteTypes = append(voteTypes, 1)
		} else if voteType2 == 1 {
			voteTypes = append(voteTypes, 2)
		} else {
			voteTypes = append(voteTypes, 0)
		}
	}

	return voteTypes, nil
}
