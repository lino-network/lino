package repository_test

import (
	"database/sql"
	"testing"

	"github.com/lino-network/lino/recorder/inflation"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lino-network/lino/recorder/dbtestutil"
	"github.com/lino-network/lino/recorder/inflation/repository"
)

func TestAddnGet(t *testing.T) {
	assert := assert.New(t)
	d1 := &inflation.Inflation{
		InfraPool:     1,
		DevPool:       2,
		CreatorPool:   3,
		ValidatorPool: 4,
		Timestamp:     1538606755,
	}

	runTest(t, func(env TestEnv) {
		err := env.coRepo.Add(d1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", d1, err)
		}
		res, err := env.coRepo.Get(1538606755)

		if err != nil {
			t.Errorf("TestAddnGet: failed to get inflation with %s, got err %v", "user1", err)
		}
		assert.Equal(d1, res)
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	coRepo repository.InflationRepository
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	conn, coinflation, err := setup()
	if err != nil {
		t.Errorf("Failed to create inflation DB : %v", err)
	}
	defer teardown(conn)

	env := TestEnv{
		coRepo: coinflation,
	}
	fc(env)
}

func setup() (*sql.DB, repository.InflationRepository, error) {
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		return nil, nil, err
	}
	coRepo, err := dbtestutil.NewInflationDB(db)
	if err != nil {
		return nil, nil, err
	}

	return db, coRepo, nil
}

func teardown(db *sql.DB) {
	// dbtestutil.InflationDBCleanUp(db)
}
