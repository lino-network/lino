package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lino-network/lino/recorder/dbtestutil"
	"github.com/lino-network/lino/recorder/reward"
	"github.com/lino-network/lino/recorder/reward/repository"
)

func TestAddnGet(t *testing.T) {
	assert := assert.New(t)
	r1 := &reward.Reward{
		Username:        "user1",
		TotalIncome:     "10000",
		OriginalIncome:  "20000",
		FrictionIncome:  "30000",
		InflationIncome: "40000",
		UnclaimReward:   "50000",
		CreatedAt:       time.Unix(time.Now().Unix(), 0).UTC(),
	}

	runTest(t, func(env TestEnv) {
		err := env.rRepo.Add(r1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", r1, err)
		}
		res, err := env.rRepo.Get("user1")

		if err != nil {
			t.Errorf("TestAddnGet: failed to get Reward with %s, got err %v", "user1", err)
		}
		r1.ID = res.ID
		assert.Equal(r1, res)
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	rRepo repository.RewardRepository
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	conn, rRepo, err := setup()
	if err != nil {
		t.Errorf("Failed to create reward DB : %v", err)
	}
	defer teardown(conn)

	env := TestEnv{
		rRepo: rRepo,
	}
	fc(env)
}

func setup() (*sql.DB, repository.RewardRepository, error) {
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		return nil, nil, err
	}
	rRepo, err := dbtestutil.NewRewardDB(db)
	if err != nil {
		return nil, nil, err
	}

	return db, rRepo, nil
}

func teardown(db *sql.DB) {
	dbtestutil.RewardDBCleanUp(db)
}
