package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lino-network/lino/recorder/dbtestutil"
	"github.com/lino-network/lino/recorder/post"
	"github.com/lino-network/lino/recorder/post/repository"
)

func TestAddnGet(t *testing.T) {
	assert := assert.New(t)
	p1 := &post.Post{
		Author:           "author",
		PostID:           "postID",
		Title:            "title",
		Content:          "content",
		ParentAuthor:     "parentauthor",
		ParentPostID:     "parentPostID",
		SourceAuthor:     "sourceAuthor",
		SourcePostID:     "sourcePostID",
		Links:            "links",
		CreatedAt:        time.Unix(time.Now().Unix(), 0).UTC(),
		TotalDonateCount: 1,
		TotalReward:      "100",
	}

	runTest(t, func(env TestEnv) {
		err := env.pRepo.Add(p1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", p1, err)
		}
		res, err := env.pRepo.Get("author")

		if err != nil {
			t.Errorf("TestAddnGet: failed to get Donation with %s, got err %v", "user1", err)
		}
		assert.Equal(p1, res)
	})
}

func TestSetReward(t *testing.T) {
	assert := assert.New(t)
	p1 := &post.Post{
		Author:           "author",
		PostID:           "postID",
		Title:            "title",
		Content:          "content",
		ParentAuthor:     "parentauthor",
		ParentPostID:     "parentPostID",
		SourceAuthor:     "sourceAuthor",
		SourcePostID:     "sourcePostID",
		Links:            "links",
		CreatedAt:        time.Unix(time.Now().Unix(), 0).UTC(),
		TotalDonateCount: 1,
		TotalReward:      "100",
	}

	runTest(t, func(env TestEnv) {
		err := env.pRepo.Add(p1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", p1, err)
		}
		p1.TotalReward = "10000"
		err = env.pRepo.SetReward(p1.Author, p1.PostID, "10000")
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", p1, err)
		}
		res, err := env.pRepo.Get("author")

		if err != nil {
			t.Errorf("TestAddnGet: failed to get Donation with %s, got err %v", "user1", err)
		}
		assert.Equal(p1, res)
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	pRepo repository.PostRepository
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	conn, coPost, err := setup()
	if err != nil {
		t.Errorf("Failed to create donation DB : %v", err)
	}
	defer teardown(conn)

	env := TestEnv{
		pRepo: coPost,
	}
	fc(env)
}

func setup() (*sql.DB, repository.PostRepository, error) {
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		return nil, nil, err
	}
	pRepo, err := dbtestutil.NewPostDB(db)
	if err != nil {
		return nil, nil, err
	}

	return db, pRepo, nil
}

func teardown(db *sql.DB) {
	dbtestutil.DonationDBCleanUp(db)
}
