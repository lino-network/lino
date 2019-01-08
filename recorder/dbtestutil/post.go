package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/post/repository"
)

func NewPostDB(db *sql.DB) (repository.PostRepository, error) {
	PostDBCleanUp(db)
	return repository.NewPostDB(db)
}

func PostDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM post")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
