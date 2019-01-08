package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/balancehistory/repository"
)

func NewBalanceHistoryDB(db *sql.DB) (repository.BalanceHistoryRepository, error) {
	BalanceHistoryDBCleanUp(db)
	return repository.NewBalanceHistoryDB(db)
}

func BalanceHistoryDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM balancehistory")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
