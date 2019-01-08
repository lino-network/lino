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
	postTableName = "post"
)

type postDB struct {
	conn  *sql.DB
	stmts map[string]*sql.Stmt
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
	}
	stmts, err := dbutils.PrepareStmts(postTableName, conn, unprepared)
	if err != nil {
		return nil, err
	}
	return &postDB{
		conn:  conn,
		stmts: stmts,
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
		totalReward      int64
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
		TotalReward:      totalReward,
	}, nil
}

func (db *postDB) Get(author string) (*post.Post, errors.Error) {
	return scanPost(db.stmts[getPost].QueryRow(author))
}

func (db *postDB) Add(post *post.Post) errors.Error {
	_, err := dbutils.ExecAffectingOneRow(db.stmts[insertPost],
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
		post.TotalReward,
	)
	return err
}
