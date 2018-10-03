package dbtestutil

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/donation/repository"
)

func NewDonationDB(db *sql.DB) (repository.DonationRepository, error) {
	DonationDBCleanUp(db)
	return repository.NewDonationDB(db)
}

func DonationDBCleanUp(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM donation")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
