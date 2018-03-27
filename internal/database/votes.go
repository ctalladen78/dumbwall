package database

import (
	"fmt"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/maksadbek/dumbwall/internal/actions"
)

var uvotes = "uvotes_tmp:%d:%d"

func (d *Database) AddVote(userID, postID int, action actions.Action) error {
	if action == actions.ActionNone {
		d.r.Do("SREM", fmt.Sprintf(uvotes, userID, actions.ActionUp), postID)
		d.r.Do("SREM", fmt.Sprintf(uvotes, userID, actions.ActionDown), postID)
		return nil
	}

	key := fmt.Sprintf(uvotes, userID, action)

	_, err := d.r.Do("SADD", key, postID)

	return err
}

func (d *Database) CheckVotes(userID int, postIDs ...int) ([]actions.Action, error) {
	var usersActions = make([]actions.Action, 0, len(postIDs))

	conn := d.r.Conn()
	defer conn.Close()

	for i := range postIDs {
		err := conn.Send("SISMEMBER", fmt.Sprintf(uvotes, userID, actions.ActionUp), postIDs[i])
		if err != nil {
			return usersActions, err
		}

		err = conn.Send("SISMEMBER", fmt.Sprintf(uvotes, userID, actions.ActionDown), postIDs[i])
		if err != nil {
			return usersActions, err
		}
	}

	err := conn.Flush()
	if err != nil {
		return usersActions, err
	}

	for i := 0; i < len(postIDs); i++ {
		voteType1, err := redigo.Int(conn.Receive())
		if err != nil {
			return usersActions, err
		}

		voteType2, err := redigo.Int(conn.Receive())
		if err != nil {
			return usersActions, err
		}

		if voteType1 == 1 {
			usersActions = append(usersActions, actions.ActionUp)
		} else if voteType2 == 1 {
			usersActions = append(usersActions, actions.ActionDown)
		} else {
			usersActions = append(usersActions, actions.ActionNone)
		}
	}

	return usersActions, nil
}
