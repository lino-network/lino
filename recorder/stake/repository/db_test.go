package repository_test

import (
	"database/sql"
	"testing"

	"github.com/lino-network/lino/recorder/stake"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lino-network/lino/recorder/dbtestutil"
	"github.com/lino-network/lino/recorder/stake/repository"
)

func TestAddnGet(t *testing.T) {
	assert := assert.New(t)
	d1 := &stake.Stake{
		Username:  "user1",
		Amount:    2000,
		Timestamp: 1538606755,
		Op:        "IN",
	}

	runTest(t, func(env TestEnv) {
		err := env.coRepo.Add(d1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", d1, err)
		}
		res, err := env.coRepo.Get("user1")

		if err != nil {
			t.Errorf("TestAddnGet: failed to get stake with %s, got err %v", "user1", err)
		}
		assert.Equal(d1, res)
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	coRepo repository.StakeRepository
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	conn, costake, err := setup()
	if err != nil {
		t.Errorf("Failed to create stake DB : %v", err)
	}
	defer teardown(conn)

	env := TestEnv{
		coRepo: costake,
	}
	fc(env)
}

func setup() (*sql.DB, repository.StakeRepository, error) {
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		return nil, nil, err
	}
	coRepo, err := dbtestutil.NewStakeDB(db)
	if err != nil {
		return nil, nil, err
	}

	return db, coRepo, nil
}

func teardown(db *sql.DB) {
	dbtestutil.StakeDBCleanUp(db)
}
