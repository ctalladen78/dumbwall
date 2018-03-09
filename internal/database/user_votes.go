package database

import (
	"github.com/maksadbek/dumbwall/internal/actions"
)

func (db *Database) AddVote(userID, postID int, voteType actions.Action) error {
	println("addVote", userID, postID, voteType)
	return db.r.AddVote(userID, postID, int(voteType))
}

func (db *Database) CheckVotes(userID int, postIDs ...int) ([]actions.Action, error) {
	var usersActions = []actions.Action{}

	voteTypes, err := db.r.CheckVotes(userID, postIDs...)
	if err != nil {
		return usersActions, err
	}

	for _, v := range voteTypes {
		usersActions = append(usersActions, actions.Action(v))

	}

	return usersActions, nil
}
