package repository

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/dbutils"
	"github.com/lino-network/lino/recorder/donation"
	errors "github.com/lino-network/lino/recorder/errors"

	_ "github.com/go-sql-driver/mysql"
)

const (
	getDonation       = "get-donation"
	insertDonation    = "insert-donation"
	donationTableName = "donation"
)

type donationDB struct {
	conn  *sql.DB
	stmts map[string]*sql.Stmt
}

var _ DonationRepository = &donationDB{}

func NewDonationDB(conn *sql.DB) (DonationRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("Donation db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		getDonation:    getDonationStmt,
		insertDonation: insertDonationStmt,
	}
	stmts, err := dbutils.PrepareStmts(donationTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &donationDB{
		conn:  conn,
		stmts: stmts,
	}, nil
}

func scanDonation(s dbutils.RowScanner) (*donation.Donation, errors.Error) {
	var (
		username       string
		seq            int64
		dp             int64
		permlink       string
		amount         int64
		fromApp        string
		coinDayDonated int64
	)
	if err := s.Scan(&username, &seq, &dp, &permlink, &amount, &fromApp, &coinDayDonated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeUserNotFound, "user not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &donation.Donation{
		Username:       username,
		Seq:            seq,
		Dp:             dp,
		Permlink:       permlink,
		Amount:         amount,
		FromApp:        fromApp,
		CoinDayDonated: coinDayDonated,
	}, nil
}

func (db *donationDB) Get(username string) (*donation.Donation, errors.Error) {
	return scanDonation(db.stmts[getDonation].QueryRow(username))
}

func (db *donationDB) Add(donation *donation.Donation) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[insertDonation],
		donation.Username,
		donation.Seq,
		donation.Dp,
		donation.Permlink,
		donation.Amount,
		donation.FromApp,
		donation.CoinDayDonated,
	)
	return err
}
