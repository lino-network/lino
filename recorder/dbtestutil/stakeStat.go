package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/stakestat/repository"
)

func NewStakeStatDB(db *sql.DB) (repository.StakeStatRepository, error) {
	StakeStatDBCleanUp(db)
	return repository.NewStakeStatDB(db)
}

func StakeStatDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM stakeStat")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
