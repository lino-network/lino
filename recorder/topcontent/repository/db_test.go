package repository_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lino-network/lino/recorder/dbtestutil"
	"github.com/lino-network/lino/recorder/topcontent"
	"github.com/lino-network/lino/recorder/topcontent/repository"
)

func TestAddnGet(t *testing.T) {
	assert := assert.New(t)
	d1 := &topcontent.TopContent{
		Permlink:  "p1",
		Timestamp: 15535,
	}

	runTest(t, func(env TestEnv) {
		err := env.coRepo.Add(d1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", d1, err)
		}
		res, err := env.coRepo.Get("p1")

		if err != nil {
			t.Errorf("TestAddnGet: failed to get TopContent with %s, got err %v ", "p1", err)
		}
		d1.ID = res.ID
		assert.Equal(d1, res)
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	coRepo repository.TopContentRepository
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	conn, coTopContent, err := setup()
	if err != nil {
		t.Errorf("Failed to create TopContent DB : %v", err)
	}
	defer teardown(conn)

	env := TestEnv{
		coRepo: coTopContent,
	}
	fc(env)
}

func setup() (*sql.DB, repository.TopContentRepository, error) {
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		return nil, nil, err
	}
	coRepo, err := dbtestutil.NewTopContentDB(db)
	if err != nil {
		return nil, nil, err
	}

	return db, coRepo, nil
}

func teardown(db *sql.DB) {
	dbtestutil.TopContentDBCleanUp(db)
}
