package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/stake/repository"
)

func NewStakeDB(db *sql.DB) (repository.StakeRepository, error) {
	StakeDBCleanUp(db)
	return repository.NewStakeDB(db)
}

func StakeDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM stake")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
