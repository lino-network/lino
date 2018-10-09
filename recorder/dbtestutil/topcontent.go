package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/topcontent/repository"
)

func NewTopContentDB(db *sql.DB) (repository.TopContentRepository, error) {
	TopContentDBCleanUp(db)
	return repository.NewTopContentDB(db)
}

func TopContentDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM topContent")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
