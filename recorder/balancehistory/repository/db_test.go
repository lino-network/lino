package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lino-network/lino/recorder/balancehistory"
	"github.com/lino-network/lino/recorder/balancehistory/repository"
	"github.com/lino-network/lino/recorder/dbtestutil"
	"github.com/lino-network/lino/types"
)

func TestAddnGet(t *testing.T) {
	assert := assert.New(t)
	b1 := &balancehistory.BalanceHistory{
		Username:   "user1",
		FromUser:   "user1",
		ToUser:     "user2",
		Amount:     100,
		Balance:    10000,
		DetailType: types.TransferIn,
		CreatedAt:  time.Unix(time.Now().Unix(), 0).UTC(),
		Memo:       "transfer",
	}

	runTest(t, func(env TestEnv) {
		err := env.bhRepo.Add(b1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", b1, err)
		}
		res, err := env.bhRepo.Get("user1")

		if err != nil {
			t.Errorf("TestAddnGet: failed to get BalanceHistory with %s, got err %v", "user1", err)
		}
		b1.ID = res.ID
		assert.Equal(b1, res)
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	bhRepo repository.BalanceHistoryRepository
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	_, bhRepo, err := setup()
	if err != nil {
		t.Errorf("Failed to create balancehistory DB : %v", err)
	}
	// defer teardown(conn)

	env := TestEnv{
		bhRepo: bhRepo,
	}
	fc(env)
}

func setup() (*sql.DB, repository.BalanceHistoryRepository, error) {
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		return nil, nil, err
	}
	bhRepo, err := dbtestutil.NewBalanceHistoryDB(db)
	if err != nil {
		return nil, nil, err
	}

	return db, bhRepo, nil
}

// func teardown(db *sql.DB) {
// 	dbtestutil.BalanceHistoryDBCleanUp(db)
// }
