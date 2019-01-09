package repository

import (
	"database/sql"
	"time"

	"github.com/lino-network/lino/recorder/dbutils"
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/post"

	_ "github.com/go-sql-driver/mysql"
)

const (
	getPost       = "get-post"
	insertPost    = "insert-post"
	setReward     = "set-reward"
	postTableName = "post"
)

type postDB struct {
	conn     *sql.DB
	stmts    map[string]*sql.Stmt
	EnableDB bool
}

var _ PostRepository = &postDB{}

func NewPostDB(conn *sql.DB) (PostRepository, errors.Error) {
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, errors.Unavailable("Donation db conn is unavaiable").TraceCause(err, "")
	}
	unprepared := map[string]string{
		getPost:    getPostStmt,
		insertPost: insertPostStmt,
		setReward:  setRewardStmt,
	}
	stmts, err := dbutils.PrepareStmts(postTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &postDB{
		conn:     conn,
		stmts:    stmts,
		EnableDB: true,
	}, nil
}

func scanPost(s dbutils.RowScanner) (*post.Post, errors.Error) {
	var (
		author           string
		postID           string
		title            string
		content          string
		parentAuthor     string
		parentPostID     string
		sourceAuthor     string
		sourcePostID     string
		links            string
		createdAt        time.Time
		totalDonateCount int64
		totalReward      string
	)
	if err := s.Scan(&author, &postID, &title, &content, &parentAuthor, &parentPostID, &sourceAuthor, &sourcePostID, &links, &createdAt, &totalDonateCount, &totalReward); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewErrorf(errors.CodeUserNotFound, "post not found: %s", err)
		}
		return nil, errors.NewErrorf(errors.CodeFailedToScan, "failed to scan %s", err)
	}

	return &post.Post{
		Author:           author,
		PostID:           postID,
		Title:            title,
		Content:          content,
		ParentAuthor:     parentAuthor,
		ParentPostID:     parentPostID,
		SourceAuthor:     sourceAuthor,
		SourcePostID:     sourcePostID,
		Links:            links,
		CreatedAt:        createdAt,
		TotalDonateCount: totalDonateCount,
		TotalReward:      dbutils.TrimPaddedZeroFromNumber(totalReward),
	}, nil
}

func (db *postDB) IsEnable() bool {
	return db.EnableDB
}
func (db *postDB) Get(author string) (*post.Post, errors.Error) {
	return scanPost(db.stmts[getPost].QueryRow(author))
}

func (db *postDB) Add(post *post.Post) errors.Error {
	totalReward, err := dbutils.PadNumberStrWithZero(post.TotalReward)
	if err != nil {
		return err
	}
	_, err = dbutils.ExecAffectingOneRow(db.stmts[insertPost],
		post.Author,
		post.PostID,
		post.Title,
		post.Content,
		post.ParentAuthor,
		post.ParentPostID,
		post.SourceAuthor,
		post.SourcePostID,
		post.Links,
		post.CreatedAt,
		post.TotalDonateCount,
		totalReward,
	)
	return err
}

func (db *postDB) SetReward(author, postID string, amount string) errors.Error {
	paddingAmount, err := dbutils.PadNumberStrWithZero(amount)
	if err != nil {
		return err
	}
	_, err = dbutils.Exec(db.stmts[setReward],
		paddingAmount, author, postID,
	)
	return err
}
