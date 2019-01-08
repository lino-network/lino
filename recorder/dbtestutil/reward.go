package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/reward/repository"
)

func NewRewardDB(db *sql.DB) (repository.RewardRepository, error) {
	RewardDBCleanUp(db)
	return repository.NewRewardDB(db)
}

func RewardDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM reward")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
