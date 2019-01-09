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
	conn     *sql.DB
	stmts    map[string]*sql.Stmt
	EnableDB bool
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
		conn:     conn,
		stmts:    stmts,
		EnableDB: true,
	}, nil
}

func scaninflation(s dbutils.RowScanner) (*inflation.Inflation, errors.Error) {
	var (
		id                 int64
		infraPool          string
		devPool            string
		creatorPool        string
		validatorPool      string
		infraInflation     string
		devInflation       string
		creatorInflation   string
		validatorInflation string
		timestamp          int64
	)
	if err := s.Scan(&id, &infraPool, &devPool, &creatorPool, &validatorPool, &infraInflation, &devInflation, &creatorInflation, &validatorInflation, &timestamp); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeUserNotFound, "user not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &inflation.Inflation{
		InfraPool:          dbutils.TrimPaddedZeroFromNumber(infraPool),
		DevPool:            dbutils.TrimPaddedZeroFromNumber(devPool),
		CreatorPool:        dbutils.TrimPaddedZeroFromNumber(creatorPool),
		ValidatorPool:      dbutils.TrimPaddedZeroFromNumber(validatorPool),
		InfraInflation:     dbutils.TrimPaddedZeroFromNumber(infraInflation),
		DevInflation:       dbutils.TrimPaddedZeroFromNumber(devInflation),
		CreatorInflation:   dbutils.TrimPaddedZeroFromNumber(creatorInflation),
		ValidatorInflation: dbutils.TrimPaddedZeroFromNumber(validatorInflation),
		Timestamp:          timestamp,
	}, nil
}

func (db *inflationDB) IsEnable() bool {
	return db.EnableDB
}
func (db *inflationDB) Get(timestamp int64) (*inflation.Inflation, errors.Error) {
	return scaninflation(db.stmts[getInflation].QueryRow(timestamp))
}

func (db *inflationDB) Add(inflation *inflation.Inflation) errors.Error {
	infraPool, err := dbutils.PadNumberStrWithZero(inflation.InfraPool)
	if err != nil {
		return err
	}
	devPool, err := dbutils.PadNumberStrWithZero(inflation.DevPool)
	if err != nil {
		return err
	}
	creatorPool, err := dbutils.PadNumberStrWithZero(inflation.CreatorPool)
	if err != nil {
		return err
	}
	validatorPool, err := dbutils.PadNumberStrWithZero(inflation.ValidatorPool)
	if err != nil {
		return err
	}
	infraInflation, err := dbutils.PadNumberStrWithZero(inflation.InfraInflation)
	if err != nil {
		return err
	}
	devInflation, err := dbutils.PadNumberStrWithZero(inflation.DevInflation)
	if err != nil {
		return err
	}
	creatorInflation, err := dbutils.PadNumberStrWithZero(inflation.CreatorInflation)
	if err != nil {
		return err
	}
	validatorInflation, err := dbutils.PadNumberStrWithZero(inflation.ValidatorInflation)
	if err != nil {
		return err
	}
	_, err = dbutils.ExecAffectingOneRow(db.stmts[insertInflation],
		infraPool,
		devPool,
		creatorPool,
		validatorPool,
		infraInflation,
		devInflation,
		creatorInflation,
		validatorInflation,
		inflation.Timestamp,
	)
	return err
}
