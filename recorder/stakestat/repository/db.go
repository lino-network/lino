package repository

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/stakestat"

	_ "github.com/go-sql-driver/mysql"
)

const (
	getStakeStat       = "get-stake-stat"
	insertStakeStat    = "insert-stake-stat"
	stakeStatTableName = "stakeStat"
)

type stakeStatDB struct {
	conn     *sql.DB
	stmts    map[string]*sql.Stmt
	EnableDB bool
}

var _ StakeStatRepository = &stakeStatDB{}

func NewStakeStatDB(conn *sql.DB) (StakeStatRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("stakeStat db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		getStakeStat:    getStakeStatStmt,
		insertStakeStat: insertStakeStatStmt,
	}
	stmts, err := dbutils.PrepareStmts(stakeStatTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &stakeStatDB{
		conn:     conn,
		stmts:    stmts,
		EnableDB: true,
	}, nil
}

func scanstakeStat(s dbutils.RowScanner) (*stakestat.StakeStat, errors.Error) {
	var (
		totalConsumptionFriction int64
		unclaimedFriction        int64
		totalLinoStake           string
		unclaimedLinoStake       string
		timestamp                int64
	)
	if err := s.Scan(&totalConsumptionFriction, &unclaimedFriction, &totalLinoStake, &unclaimedLinoStake, &timestamp); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeUserNotFound, "user not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &stakestat.StakeStat{
		TotalConsumptionFriction: totalConsumptionFriction,
		UnclaimedFriction:        unclaimedFriction,
		TotalLinoStake:           dbutils.TrimPaddedZeroFromNumber(totalLinoStake),
		UnclaimedLinoStake:       dbutils.TrimPaddedZeroFromNumber(unclaimedLinoStake),
		Timestamp:                timestamp,
	}, nil
}

func (db *stakeStatDB) IsEnable() bool {
	return db.EnableDB
}
func (db *stakeStatDB) Get(timestamp int64) (*stakestat.StakeStat, errors.Error) {
	return scanstakeStat(db.stmts[getStakeStat].QueryRow(timestamp))
}

func (db *stakeStatDB) Add(stakeStat *stakestat.StakeStat) errors.Error {
	totalLinoStake, err := dbutils.PadNumberStrWithZero(stakeStat.TotalLinoStake)
	if err != nil {
		return err
	}
	unclaimedLinoStake, err := dbutils.PadNumberStrWithZero(stakeStat.UnclaimedLinoStake)
	if err != nil {
		return err
	}
	_, err = dbutils.ExecAffectingOneRow(db.stmts[insertStakeStat],
		stakeStat.TotalConsumptionFriction,
		stakeStat.UnclaimedFriction,
		totalLinoStake,
		unclaimedLinoStake,
		stakeStat.Timestamp,
	)
	return err
}
