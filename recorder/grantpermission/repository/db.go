package repository

import (
	"database/sql"
	"time"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/grantpermission"
	"github.com/lino-network/lino/types"

	_ "github.com/go-sql-driver/mysql"
)

const (
	insertPermission = "insert-permission"
	getPermission    = "get-permission"
	updateAmount     = "update-amount"
	updatePermission = "update-permission"
	deletePermission = "delete-permission"

	grantPermissionTableName = "grantpermission"
)

type grantPermissionDB struct {
	conn     *sql.DB
	stmts    map[string]*sql.Stmt
	EnableDB bool
}

var _ GrantPermissionRepository = &grantPermissionDB{}

func NewGrantPermissionDB(conn *sql.DB) (GrantPermissionRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("grant permission db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		insertPermission: insertGrantPermissionStmt,
		getPermission:    getGrantPermissionStmt,
		updateAmount:     updateAmountStmt,
		updatePermission: updateGrantPermissionStmt,
		deletePermission: deletePermissionStmt,
	}
	stmts, err := dbutils.PrepareStmts(grantPermissionTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &grantPermissionDB{
		EnableDB: true,
		conn:     conn,
		stmts:    stmts,
	}, nil
}

func scanPermission(s dbutils.RowScanner) (*grantpermission.GrantPubKey, errors.Error) {
	var (
		username   string
		authTo     string
		permission int
		createdAt  time.Time
		expiresAt  time.Time
		amount     string
	)
	if err := s.Scan(&username, &authTo, &permission, &createdAt, &expiresAt, &amount); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeFailedToScan, "permission not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &grantpermission.GrantPubKey{
		Username:   username,
		AuthTo:     authTo,
		Permission: types.Permission(permission),
		CreatedAt:  createdAt,
		ExpiresAt:  expiresAt,
		Amount:     dbutils.TrimPaddedZeroFromNumber(amount),
	}, nil
}

func (db *grantPermissionDB) IsEnable() bool {
	return db.EnableDB
}

func (db *grantPermissionDB) Get(username, authTo string) (*grantpermission.GrantPubKey, errors.Error) {
	return scanPermission(db.stmts[getPermission].QueryRow(username, authTo))
}

func (db *grantPermissionDB) Add(grantPubKey *grantpermission.GrantPubKey) errors.Error {
	paddingAmount, err := dbutils.PadNumberStrWithZero(grantPubKey.Amount)
	if err != nil {
		return err
	}
	_, err = dbutils.ExecAffectingOneRow(db.stmts[insertPermission],
		grantPubKey.Username,
		grantPubKey.AuthTo,
		grantPubKey.Permission,
		grantPubKey.CreatedAt,
		grantPubKey.ExpiresAt,
		paddingAmount,
	)
	return err
}

func (db *grantPermissionDB) SetAmount(username, authTo, amount string) errors.Error {
	paddingAmount, err := dbutils.PadNumberStrWithZero(amount)
	if err != nil {
		return err
	}
	_, err = dbutils.ExecAffectingOneRow(db.stmts[updateAmount],
		paddingAmount,
		username,
		authTo,
	)
	return err
}

func (db *grantPermissionDB) Update(grantPubKey *grantpermission.GrantPubKey) errors.Error {
	paddingAmount, err := dbutils.PadNumberStrWithZero(grantPubKey.Amount)
	if err != nil {
		return err
	}
	_, err = dbutils.Exec(db.stmts[updatePermission],
		grantPubKey.Permission,
		grantPubKey.CreatedAt,
		grantPubKey.ExpiresAt,
		paddingAmount,
		grantPubKey.Username,
		grantPubKey.AuthTo,
	)
	return err
}

func (db *grantPermissionDB) Delete(username, authTo string) errors.Error {
	_, err := dbutils.Exec(db.stmts[deletePermission],
		username, authTo,
	)
	return err
}
