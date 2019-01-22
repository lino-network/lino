package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/grantpermission/repository"
)

func NewPermissionDB(db *sql.DB) (repository.GrantPermissionRepository, error) {
	PermissionDBCleanUp(db)
	return repository.NewGrantPermissionDB(db)
}

func PermissionDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM grantpermission")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
