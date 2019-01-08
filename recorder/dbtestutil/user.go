package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/user/repository"
)

func NewUserDB(db *sql.DB) (repository.UserRepository, error) {
	UserDBCleanUp(db)
	return repository.NewUserDB(db)
}

func UserDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM user")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
