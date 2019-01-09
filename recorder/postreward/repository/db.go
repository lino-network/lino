package repository

import (
	"database/sql"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/postreward"

	_ "github.com/go-sql-driver/mysql"
)

const (
	getPostReward       = "get-post-reward"
	insertPostReward    = "insert-post-reward"
	PostRewardTableName = "postReward"
)

type postRewardDB struct {
	conn     *sql.DB
	stmts    map[string]*sql.Stmt
	EnableDB bool
}

var _ PostRewardRepository = &postRewardDB{}

func NewPostRewardDB(conn *sql.DB) (PostRewardRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("PostReward db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		getPostReward:    getPostRewardStmt,
		insertPostReward: insertPostRewardStmt,
	}
	stmts, err := dbutils.PrepareStmts(PostRewardTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &postRewardDB{
		conn:     conn,
		stmts:    stmts,
		EnableDB: true,
	}, nil
}

func scanPostReward(s dbutils.RowScanner) (*postreward.PostReward, errors.Error) {
	var (
		id           int64
		permlink     string
		reward       int64
		penaltyScore string
		timestamp    int64
		evaluate     int64
		original     int64
		consumer     string
	)
	if err := s.Scan(&id, &permlink, &reward, &penaltyScore, &timestamp, &evaluate, &original, &consumer); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeUserNotFound, "post not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &postreward.PostReward{
		ID:           id,
		Permlink:     permlink,
		Reward:       reward,
		PenaltyScore: penaltyScore,
		Timestamp:    timestamp,
		Evaluate:     evaluate,
		Original:     original,
		Consumer:     consumer,
	}, nil
}

func (db *postRewardDB) IsEnable() bool {
	return db.EnableDB
}
func (db *postRewardDB) Get(permlink string) (*postreward.PostReward, errors.Error) {
	return scanPostReward(db.stmts[getPostReward].QueryRow(permlink))
}

func (db *postRewardDB) Add(postReward *postreward.PostReward) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[insertPostReward],
		postReward.Permlink,
		postReward.Reward,
		postReward.PenaltyScore,
		postReward.Timestamp,
		postReward.Evaluate,
		postReward.Original,
		postReward.Consumer,
	)
	return err
}
