package repository

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/topcontent"

	_ "github.com/go-sql-driver/mysql"
)

const (
	getTopContent       = "get-top-content"
	insertTopContent    = "insert-top-content"
	topContentTableName = "topContent"
)

type topContentDB struct {
	conn  *sql.DB
	stmts map[string]*sql.Stmt
}

var _ TopContentRepository = &topContentDB{}

func NewTopContentDB(conn *sql.DB) (TopContentRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("TopContent db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		getTopContent:    getTopContentStmt,
		insertTopContent: insertTopContentStmt,
	}
	stmts, err := dbutils.PrepareStmts(topContentTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &topContentDB{
		conn:  conn,
		stmts: stmts,
	}, nil
}

func scanTopContent(s dbutils.RowScanner) (*topcontent.TopContent, errors.Error) {
	var (
		id        int64
		permlink  string
		timestamp int64
	)
	if err := s.Scan(&id, &permlink, &timestamp); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeUserNotFound, "TopContent not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &topcontent.TopContent{
		ID:        id,
		Permlink:  permlink,
		Timestamp: timestamp,
	}, nil
}

func (db *topContentDB) Get(permlink string) (*topcontent.TopContent, errors.Error) {
	return scanTopContent(db.stmts[getTopContent].QueryRow(permlink))
}

func (db *topContentDB) Add(topContent *topcontent.TopContent) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[insertTopContent],
		topContent.Permlink,
		topContent.Timestamp,
	)
	return err
}
