package repository_test

import (
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lino-network/lino/recorder/dbtestutil"
	"github.com/lino-network/lino/recorder/user"
	"github.com/lino-network/lino/recorder/user/repository"
)

func TestAddnGet(t *testing.T) {
	assert := assert.New(t)
	r1 := &user.User{
		Username:          "user1",
		CreatedAt:         time.Unix(time.Now().Unix(), 0).UTC(),
		ResetPubKey:       hex.EncodeToString(secp256k1.GenPrivKey().PubKey().Bytes()),
		TransactionPubKey: hex.EncodeToString(secp256k1.GenPrivKey().PubKey().Bytes()),
		AppPubKey:         hex.EncodeToString(secp256k1.GenPrivKey().PubKey().Bytes()),
		Saving:            "100000",
		Sequence:          1,
	}

	runTest(t, func(env TestEnv) {
		err := env.uRepo.Add(r1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", r1, err)
		}
		res, err := env.uRepo.Get("user1")

		if err != nil {
			t.Errorf("TestAddnGet: failed to get User with %s, got err %v", "user1", err)
		}
		assert.Equal(r1, res)
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	uRepo repository.UserRepository
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	conn, uRepo, err := setup()
	if err != nil {
		t.Errorf("Failed to create donation DB : %v", err)
	}
	defer teardown(conn)

	env := TestEnv{
		uRepo: uRepo,
	}
	fc(env)
}

func setup() (*sql.DB, repository.UserRepository, error) {
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		return nil, nil, err
	}
	uRepo, err := dbtestutil.NewUserDB(db)
	if err != nil {
		return nil, nil, err
	}

	return db, uRepo, nil
}

func teardown(db *sql.DB) {
	dbtestutil.UserDBCleanUp(db)
}
