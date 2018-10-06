package repository

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/stakeStat"

	_ "github.com/go-sql-driver/mysql"
)

const (
	getStakeStat       = "get-stake-stat"
	insertStakeStat    = "insert-stake-stat"
	stakeStatTableName = "stakeStat"
)

type stakeStatDB struct {
	conn  *sql.DB
	stmts map[string]*sql.Stmt
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
		conn:  conn,
		stmts: stmts,
	}, nil
}

func scanstakeStat(s dbutils.RowScanner) (*stakeStat.StakeStat, errors.Error) {
	var (
		totalConsumptionFriction int64
		unclaimedFriction        int64
		totalLinoStake           int64
		unclaimedLinoStake       int64
		timestamp                int64
	)
	if err := s.Scan(&totalConsumptionFriction, &unclaimedFriction, &totalLinoStake, &unclaimedLinoStake, &timestamp); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeUserNotFound, "user not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &stakeStat.StakeStat{
		TotalConsumptionFriction: totalConsumptionFriction,
		UnclaimedFriction:        unclaimedFriction,
		TotalLinoStake:           totalLinoStake,
		UnclaimedLinoStake:       unclaimedLinoStake,
		Timestamp:                timestamp,
	}, nil
}

func (db *stakeStatDB) Get(timestamp int64) (*stakeStat.StakeStat, errors.Error) {
	return scanstakeStat(db.stmts[getStakeStat].QueryRow(timestamp))
}

func (db *stakeStatDB) Add(stakeStat *stakeStat.StakeStat) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[insertStakeStat],
		stakeStat.TotalConsumptionFriction,
		stakeStat.UnclaimedFriction,
		stakeStat.TotalLinoStake,
		stakeStat.UnclaimedLinoStake,
		stakeStat.Timestamp,
	)
	return err
}
