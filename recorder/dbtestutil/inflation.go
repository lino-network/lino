package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/inflation/repository"
)

func NewInflationDB(db *sql.DB) (repository.InflationRepository, error) {
	InflationDBCleanUp(db)
	return repository.NewinflationDB(db)
}

func InflationDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM inflation")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
