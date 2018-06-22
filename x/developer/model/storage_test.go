package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tmlibs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("developer")
)

func TestDeveloper(t *testing.T) {
	developer := Developer{
		Username:       "user1",
		Deposit:        types.NewCoinFromInt64(100),
		AppConsumption: types.NewCoinFromInt64(1000),
	}

	runTest(t, func(env TestEnv) {
		err := env.ds.SetDeveloper(env.ctx, developer.Username, &developer)
		assert.Nil(t, err)

		resultPtr, err := env.ds.GetDeveloper(env.ctx, developer.Username)
		assert.Nil(t, err)
		assert.Equal(t, developer, *resultPtr, "developer should be equal")
	})
}

func TestDeveloperList(t *testing.T) {
	lst := DeveloperList{
		AllDevelopers: []types.AccountKey{types.AccountKey("u1"), types.AccountKey("u2")},
	}

	runTest(t, func(env TestEnv) {
		err := env.ds.SetDeveloperList(env.ctx, &lst)
		assert.Nil(t, err)

		resultPtr, err := env.ds.GetDeveloperList(env.ctx)
		assert.Nil(t, err)
		assert.Equal(t, lst, *resultPtr, "developer list should be equal")
	})

}

//
// Test Environment setup
//

type TestEnv struct {
	ds  DeveloperStorage
	ctx sdk.Context
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	env := TestEnv{
		ds:  NewDeveloperStorage(TestKVStoreKey),
		ctx: getContext(),
	}
	fc(env)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
}
