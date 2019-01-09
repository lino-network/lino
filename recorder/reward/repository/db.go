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
	conn     *sql.DB
	stmts    map[string]*sql.Stmt
	EnableDB bool
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
		conn:     conn,
		stmts:    stmts,
		EnableDB: true,
	}, nil
}

func scanReward(s dbutils.RowScanner) (*reward.Reward, errors.Error) {
	var (
		id              int64
		username        string
		totalIncome     string
		originalIncome  string
		frictionIncome  string
		inflationIncome string
		unclaimReward   string
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
		TotalIncome:     dbutils.TrimPaddedZeroFromNumber(totalIncome),
		OriginalIncome:  dbutils.TrimPaddedZeroFromNumber(originalIncome),
		FrictionIncome:  dbutils.TrimPaddedZeroFromNumber(frictionIncome),
		InflationIncome: dbutils.TrimPaddedZeroFromNumber(inflationIncome),
		UnclaimReward:   dbutils.TrimPaddedZeroFromNumber(unclaimReward),
		CreatedAt:       createdAt,
	}, nil
}

func (db *rewardDB) IsEnable() bool {
	return db.EnableDB
}
func (db *rewardDB) Get(username string) (*reward.Reward, errors.Error) {
	return scanReward(db.stmts[getReward].QueryRow(username))
}
func (db *rewardDB) Add(reward *reward.Reward) errors.Error {
	totalIncome, err := dbutils.PadNumberStrWithZero(reward.TotalIncome)
	if err != nil {
		return err
	}
	originalIncome, err := dbutils.PadNumberStrWithZero(reward.OriginalIncome)
	if err != nil {
		return err
	}
	frictionIncome, err := dbutils.PadNumberStrWithZero(reward.FrictionIncome)
	if err != nil {
		return err
	}
	inflationIncome, err := dbutils.PadNumberStrWithZero(reward.InflationIncome)
	if err != nil {
		return err
	}
	unclaimReward, err := dbutils.PadNumberStrWithZero(reward.UnclaimReward)
	if err != nil {
		return err
	}
	_, err = dbutils.ExecAffectingOneRow(db.stmts[insertReward],
		reward.Username,
		totalIncome,
		originalIncome,
		frictionIncome,
		inflationIncome,
		unclaimReward,
		reward.CreatedAt,
	)
	return err
}
