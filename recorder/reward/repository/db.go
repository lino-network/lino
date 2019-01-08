package repository

import (
	"database/sql"
	"time"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/reward"

	_ "github.com/go-sql-driver/mysql"
)

const (
	insertReward = "insert-reward"
	getReward    = "get-reward"

	rewardTableName = "reward"
)

type rewardDB struct {
	conn  *sql.DB
	stmts map[string]*sql.Stmt
}

var _ RewardRepository = &rewardDB{}

func NewRewardDB(conn *sql.DB) (RewardRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("Balance history db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		insertReward: insertRewardStmt,
		getReward:    getRewardStmt,
	}
	stmts, err := dbutils.PrepareStmts(rewardTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &rewardDB{
		conn:  conn,
		stmts: stmts,
	}, nil
}

func scanReward(s dbutils.RowScanner) (*reward.Reward, errors.Error) {
	var (
		id              int64
		username        string
		totalIncome     int64
		originalIncome  int64
		frictionIncome  int64
		inflationIncome int64
		unclaimReward   int64
		createdAt       time.Time
	)
	if err := s.Scan(&id, &username, &totalIncome, &originalIncome, &frictionIncome, &inflationIncome, &unclaimReward, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeFailedToScan, "balance history not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &reward.Reward{
		ID:              id,
		Username:        username,
		TotalIncome:     totalIncome,
		OriginalIncome:  originalIncome,
		FrictionIncome:  frictionIncome,
		InflationIncome: inflationIncome,
		UnclaimReward:   unclaimReward,
		CreatedAt:       createdAt,
	}, nil
}

func (db *rewardDB) Get(username string) (*reward.Reward, errors.Error) {
	return scanReward(db.stmts[getReward].QueryRow(username))
}
func (db *rewardDB) Add(reward *reward.Reward) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[insertReward],
		reward.Username,
		reward.TotalIncome,
		reward.OriginalIncome,
		reward.FrictionIncome,
		reward.InflationIncome,
		reward.UnclaimReward,
		reward.CreatedAt,
	)
	return err
}
