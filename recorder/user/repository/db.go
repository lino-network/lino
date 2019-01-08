package repository

import (
	"database/sql"
	"time"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/user"

	_ "github.com/go-sql-driver/mysql"
)

const (
	insertUser       = "insert-user"
	getUser          = "get-user"
	increaseSeqByOne = "increase-sequence-by-one"
	updatePubKey     = "update-pub-key"
	updateBalance    = "update-balance"

	userTableName = "user"
)

type userDB struct {
	conn  *sql.DB
	stmts map[string]*sql.Stmt
}

var _ UserRepository = &userDB{}

func NewUserDB(conn *sql.DB) (UserRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("Balance history db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		insertUser:       insertUserStmt,
		getUser:          getUserStmt,
		increaseSeqByOne: increaseSeqByOneStmt,
		updatePubKey:     updatePubKeyStmt,
		updateBalance:    updateBalanceStmt,
	}
	stmts, err := dbutils.PrepareStmts(userTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &userDB{
		conn:  conn,
		stmts: stmts,
	}, nil
}

func scanUser(s dbutils.RowScanner) (*user.User, errors.Error) {
	var (
		username          string
		createdAt         time.Time
		resetPubKey       string
		transactionPubKey string
		appPubKey         string
		saving            string
		sequence          int64
	)
	if err := s.Scan(&username, &createdAt, &resetPubKey, &transactionPubKey, &appPubKey, &saving, &sequence); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeFailedToScan, "user not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &user.User{
		Username:          username,
		CreatedAt:         createdAt,
		ResetPubKey:       resetPubKey,
		TransactionPubKey: transactionPubKey,
		AppPubKey:         appPubKey,
		Saving:            dbutils.TrimPaddedZeroFromNumber(saving),
		Sequence:          sequence,
	}, nil
}

func (db *userDB) Get(username string) (*user.User, errors.Error) {
	return scanUser(db.stmts[getUser].QueryRow(username))
}
func (db *userDB) Add(user *user.User) errors.Error {
	paddingSaving, err := dbutils.PadNumberStrWithZero(user.Saving)
	if err != nil {
		return err
	}
	_, err = dbutils.ExecAffectingOneRow(db.stmts[insertUser],
		user.Username,
		user.CreatedAt,
		user.ResetPubKey,
		user.TransactionPubKey,
		user.AppPubKey,
		paddingSaving,
		user.Sequence,
	)
	return err
}

func (db *userDB) IncreaseSequenceNumber(username string) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[increaseSeqByOne],
		username,
	)
	return err
}

func (db *userDB) UpdatePubKey(username, resetPubKey, txPubKey, appPubKey string) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[updatePubKey],
		resetPubKey,
		txPubKey,
		appPubKey,
		username,
	)
	return err
}

func (db *userDB) UpdateBalance(username string, balance string) errors.Error {
	paddingBalance, err := dbutils.PadNumberStrWithZero(balance)
	if err != nil {
		return err
	}
	_, err = dbutils.ExecAffectingOneRow(db.stmts[updateBalance],
		paddingBalance,
		username,
	)
	return err
}
