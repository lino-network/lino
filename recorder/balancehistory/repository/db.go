package repository

import (
	"database/sql"
	"time"

	"github.com/lino-network/lino/recorder/balancehistory"
	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/types"

	_ "github.com/go-sql-driver/mysql"
)

const (
	insertBalanceHistory = "insert-balance-history"
	getBalanceHistory    = "get-balance-history"

	balanceHistoryTableName = "balancehistory"
)

type balanceHistoryDB struct {
	conn  *sql.DB
	stmts map[string]*sql.Stmt
}

var _ BalanceHistoryRepository = &balanceHistoryDB{}

func NewBalanceHistoryDB(conn *sql.DB) (BalanceHistoryRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("Balance history db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		insertBalanceHistory: insertBalanceHistoryStmt,
		getBalanceHistory:    getBalanceHistoryStmt,
	}
	stmts, err := dbutils.PrepareStmts(balanceHistoryTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &balanceHistoryDB{
		conn:  conn,
		stmts: stmts,
	}, nil
}

func scanBalanceHistory(s dbutils.RowScanner) (*balancehistory.BalanceHistory, errors.Error) {
	var (
		id         int64
		username   string
		fromUser   string
		toUser     string
		amount     int64
		balance    string
		detailType int64
		createdAt  time.Time
		memo       string
	)
	if err := s.Scan(&id, &username, &fromUser, &toUser, &amount, &balance, &detailType, &createdAt, &memo); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeFailedToScan, "balance history not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &balancehistory.BalanceHistory{
		ID:         id,
		Username:   username,
		FromUser:   username,
		ToUser:     toUser,
		Amount:     amount,
		Balance:    dbutils.TrimPaddedZeroFromNumber(balance),
		DetailType: types.TransferDetailType(detailType),
		CreatedAt:  createdAt,
		Memo:       memo,
	}, nil
}

func (db *balanceHistoryDB) Get(username string) (*balancehistory.BalanceHistory, errors.Error) {
	return scanBalanceHistory(db.stmts[getBalanceHistory].QueryRow(username))
}
func (db *balanceHistoryDB) Add(balanceHistory *balancehistory.BalanceHistory) errors.Error {
	paddedBalance, err := dbutils.PadNumberStrWithZero(balanceHistory.Balance)
	if err != nil {
		return err
	}
	_, err = dbutils.ExecAffectingOneRow(db.stmts[insertBalanceHistory],
		balanceHistory.Username,
		balanceHistory.FromUser,
		balanceHistory.ToUser,
		balanceHistory.Amount,
		paddedBalance,
		balanceHistory.DetailType,
		balanceHistory.CreatedAt,
		balanceHistory.Memo,
	)
	return err
}
