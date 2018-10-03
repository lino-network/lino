package repository_test

import (
	"database/sql"
	"testing"

	"github.com/lino-network/lino/recorder/donation"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lino-network/lino/recorder/dbtestutil"
	"github.com/lino-network/lino/recorder/donation/repository"
)

func TestAddnGet(t *testing.T) {
	assert := assert.New(t)
	d1 := &donation.Donation{
		Username:       "user1",
		Seq:            0,
		Dp:             1000,
		Permlink:       "p1",
		Amount:         2000,
		FromApp:        "live",
		CoinDayDonated: 3000,
		Reputation:     4000,
		Timestamp:      1538606755,
		EvaluateResult: 5000,
	}

	runTest(t, func(env TestEnv) {
		err := env.coRepo.Add(d1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", d1, err)
		}
		res, err := env.coRepo.Get("user1")

		if err != nil {
			t.Errorf("TestAddnGet: failed to get Donation with %s, got err %v", "user1", err)
		}
		assert.Equal(d1, res)
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	coRepo repository.DonationRepository
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	conn, coDonation, err := setup()
	if err != nil {
		t.Errorf("Failed to create donation DB : %v", err)
	}
	defer teardown(conn)

	env := TestEnv{
		coRepo: coDonation,
	}
	fc(env)
}

func setup() (*sql.DB, repository.DonationRepository, error) {
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		return nil, nil, err
	}
	coRepo, err := dbtestutil.NewDonationDB(db)
	if err != nil {
		return nil, nil, err
	}

	return db, coRepo, nil
}

func teardown(db *sql.DB) {
	dbtestutil.DonationDBCleanUp(db)
}
