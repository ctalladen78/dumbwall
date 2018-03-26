package database

var uvotes = "uvotes_tmp:%d:%d"

func (d *Database) AddVote(userID, postID, voteType int) error {
	if voteType == 0 {
		d.r.Do("SREM", fmt.Sprintf(uvotes, userID, 1), postID)
		d.r.Do("SREM", fmt.Sprintf(uvotes, userID, 2), postID)
		return nil
	}

	key := fmt.Sprintf(uvotes, userID, voteType)

	_, err := c.Do("SADD", key, postID)

	return err
}

func (d *Database) CheckVotes(userID int, postIDs ...int) ([]int, error) {
	var voteTypes = make([]int, 0, len(postIDs))

	conn := r.r.Conn()
	defer conn.Close()

	for i := range postIDs {
		err := conn.Send("SISMEMBER", fmt.Sprintf(uvotes, userID, 1), postIDs[i])
		if err != nil {
			return voteTypes, err
		}

		err = conn.Send("SISMEMBER", fmt.Sprintf(uvotes, userID, 2), postIDs[i])
		if err != nil {
			return voteTypes, err
		}
	}

	err := conn.Flush()
	if err != nil {
		return voteTypes, err
	}

	for i := 0; i < len(postIDs); i++ {
		voteType1, err := redigo.Int(conn.Receive())
		if err != nil {
			return voteTypes, err
		}

		voteType2, err := redigo.Int(conn.Receive())
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
