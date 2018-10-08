package repository

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/stake"

	_ "github.com/go-sql-driver/mysql"
)

const (
	getStake       = "get-stake"
	insertStake    = "insert-stake"
	stakeTableName = "stake"
)

type stakeDB struct {
	conn  *sql.DB
	stmts map[string]*sql.Stmt
}

var _ StakeRepository = &stakeDB{}

func NewStakeDB(conn *sql.DB) (StakeRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("stake db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		getStake:    getStakeStmt,
		insertStake: insertStakeStmt,
	}
	stmts, err := dbutils.PrepareStmts(stakeTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &stakeDB{
		conn:  conn,
		stmts: stmts,
	}, nil
}

func scanstake(s dbutils.RowScanner) (*stake.Stake, errors.Error) {
	var (
		username  string
		amount    int64
		timestamp int64
		op        string
	)
	if err := s.Scan(&username, &amount, &timestamp, &op); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeUserNotFound, "stake not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &stake.Stake{
		Username:  username,
		Amount:    amount,
		Timestamp: timestamp,
		Op:        op,
	}, nil
}

func (db *stakeDB) Get(timestamp string) (*stake.Stake, errors.Error) {
	return scanstake(db.stmts[getStake].QueryRow(timestamp))
}

func (db *stakeDB) Add(stake *stake.Stake) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[insertStake],
		stake.Username,
		stake.Amount,
		stake.Timestamp,
		stake.Op,
	)
	return err
}
