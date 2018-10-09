package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/postreward/repository"
)

func NewPostRewardDB(db *sql.DB) (repository.PostRewardRepository, error) {
	PostRewardDBCleanUp(db)
	return repository.NewPostRewardDB(db)
}

func PostRewardDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM postReward")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
