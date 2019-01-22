package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lino-network/lino/recorder/dbtestutil"
	"github.com/lino-network/lino/recorder/grantpermission"
	"github.com/lino-network/lino/recorder/grantpermission/repository"
	"github.com/lino-network/lino/types"
)

func TestAddnGet(t *testing.T) {
	assert := assert.New(t)
	g1 := &grantpermission.GrantPubKey{
		Username:   "user1",
		AuthTo:     "authto",
		Permission: types.Permission(0),
		CreatedAt:  time.Unix(time.Now().Unix(), 0).UTC(),
		ExpiresAt:  time.Unix(time.Now().Unix(), 0).UTC(),
		Amount:     "1",
	}

	runTest(t, func(env TestEnv) {
		err := env.gRepo.Add(g1)
		if err != nil {
			t.Errorf("TestAddnGet: failed to add %v, got err %v", g1, err)
		}
		res, err := env.gRepo.Get("user1", "authto")

		if err != nil {
			t.Errorf("TestAddnGet: failed to get permission, got err %v", err)
		}
		assert.Equal(g1, res)
	})
}

func TestSetAmount(t *testing.T) {
	assert := assert.New(t)
	g1 := &grantpermission.GrantPubKey{
		Username:   "user1",
		AuthTo:     "authto",
		Permission: types.Permission(0),
		CreatedAt:  time.Unix(time.Now().Unix(), 0).UTC(),
		ExpiresAt:  time.Unix(time.Now().Unix(), 0).UTC(),
		Amount:     "1",
	}

	runTest(t, func(env TestEnv) {
		err := env.gRepo.Add(g1)
		if err != nil {
			t.Errorf("TestSetAmount: failed to set %v, got err %v", g1, err)
		}

		err = env.gRepo.SetAmount("user1", "authto", "100000000")
		if err != nil {
			t.Errorf("TestSetAmount: failed to update amount %v, got err %v", g1, err)
		}
		g1.Amount = "100000000"
		res, err := env.gRepo.Get("user1", "authto")

		if err != nil {
			t.Errorf("TestSetAmount: failed to get permission, got err %v", err)
		}
		assert.Equal(g1, res)
	})
}

func TestUpdatePermission(t *testing.T) {
	assert := assert.New(t)
	g1 := &grantpermission.GrantPubKey{
		Username:   "user1",
		AuthTo:     "authto",
		Permission: types.Permission(0),
		CreatedAt:  time.Unix(time.Now().Unix(), 0).UTC(),
		ExpiresAt:  time.Unix(time.Now().Unix(), 0).UTC(),
		Amount:     "1",
	}

	runTest(t, func(env TestEnv) {
		err := env.gRepo.Add(g1)
		if err != nil {
			t.Errorf("TestSetAmount: failed to set %v, got err %v", g1, err)
		}
		g1.ExpiresAt = time.Unix(time.Now().Add(10*time.Hour).Unix(), 0).UTC()
		g1.CreatedAt = time.Unix(time.Now().Add(time.Hour).Unix(), 0).UTC()
		g1.Amount = "1000000"
		err = env.gRepo.Update(g1)
		if err != nil {
			t.Errorf("TestSetAmount: failed to update amount %v, got err %v", g1, err)
		}
		res, err := env.gRepo.Get("user1", "authto")

		if err != nil {
			t.Errorf("TestSetAmount: failed to get permission, got err %v", err)
		}
		assert.Equal(g1, res)
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	gRepo repository.GrantPermissionRepository
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	conn, gRepo, err := setup()
	if err != nil {
		t.Errorf("Failed to create donation DB : %v", err)
	}
	defer teardown(conn)

	env := TestEnv{
		gRepo: gRepo,
	}
	fc(env)
}

func setup() (*sql.DB, repository.GrantPermissionRepository, error) {
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		return nil, nil, err
	}
	gRepo, err := dbtestutil.NewPermissionDB(db)
	if err != nil {
		return nil, nil, err
	}

	return db, gRepo, nil
}

func teardown(db *sql.DB) {
	dbtestutil.PermissionDBCleanUp(db)
}
