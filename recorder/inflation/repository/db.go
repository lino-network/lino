package repository

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/inflation"

	_ "github.com/go-sql-driver/mysql"
)

const (
	getInflation       = "get-inflation"
	insertInflation    = "insert-inflation"
	InflationTableName = "inflation"
)

type inflationDB struct {
	conn  *sql.DB
	stmts map[string]*sql.Stmt
}

var _ InflationRepository = &inflationDB{}

func NewinflationDB(conn *sql.DB) (InflationRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("inflation db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		getInflation:    getInflationStmt,
		insertInflation: insertInflationStmt,
	}
	stmts, err := dbutils.PrepareStmts(InflationTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &inflationDB{
		conn:  conn,
		stmts: stmts,
	}, nil
}

func scaninflation(s dbutils.RowScanner) (*inflation.Inflation, errors.Error) {
	var (
		infraPool     int64
		devPool       int64
		creatorPool   int64
		validatorPool int64
		timestamp     int64
	)
	if err := s.Scan(&infraPool, &devPool, &creatorPool, &validatorPool, &timestamp); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeUserNotFound, "user not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &inflation.Inflation{
		InfraPool:     infraPool,
		DevPool:       devPool,
		CreatorPool:   creatorPool,
		ValidatorPool: validatorPool,
		Timestamp:     timestamp,
	}, nil
}

func (db *inflationDB) Get(timestamp int64) (*inflation.Inflation, errors.Error) {
	return scaninflation(db.stmts[getInflation].QueryRow(timestamp))
}

func (db *inflationDB) Add(inflation *inflation.Inflation) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[insertInflation],
		inflation.InfraPool,
		inflation.DevPool,
		inflation.CreatorPool,
		inflation.ValidatorPool,
		inflation.Timestamp,
	)
	return err
}
