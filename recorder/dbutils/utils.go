package dbutils

import (
	"database/sql"
	"fmt"

	"github.com/lino-network/lino/recorder/errors"
)

// ExecAffectingOneRow executes a given statement, expecting one row to be affected.
func ExecAffectingOneRow(stmt *sql.Stmt, args ...interface{}) (sql.Result, errors.Error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, errors.Internalf("ExecAffectingOneRow: failed to execute statement [%v]", stmt).TraceCause(err, "")
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return r, errors.Internalf("ExecAffectingOneRow: can't get rows affected for statement [%v]", stmt).TraceCause(err, "")
	} else if rowsAffected != 1 {
		return r, errors.Internalf("ExecAffectingOneRow: expect 1, but got [%d] row affected for statement [%v]", rowsAffected, stmt)
	}
	return r, nil
}

// Exec executes a given statement
func Exec(stmt *sql.Stmt, args ...interface{}) (sql.Result, errors.Error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, errors.Internalf("Exec: failed to execute statement [%v]", stmt).TraceCause(err, "")
	}
	return r, nil
}

// RowScanner is implemented by sql.Row and sql.Rows
type RowScanner interface {
	Scan(dest ...interface{}) error
}

// PrepareStmts will attempt to prepare each unprepared
// query on the database. If one fails, the function returns
// with an error.
func PrepareStmts(service string, db *sql.DB, unprepared map[string]string) (map[string]*sql.Stmt, errors.Error) {
	prepared := map[string]*sql.Stmt{}
	for k, v := range unprepared {
		stmt, err := db.Prepare(v)
		if err != nil {
			return nil, errors.UnablePrepareStatement(fmt.Sprintf("service: %s can't prepare %v statement", service, stmt)).TraceCause(err, "")
		}
		prepared[k] = stmt
	}

	return prepared, nil
}
