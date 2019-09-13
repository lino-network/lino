package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("infra")
)

func TestInfraProvider(t *testing.T) {
	provider := InfraProvider{
		Username: "user1",
		Usage:    int64(1000),
	}

	runTest(t, func(env TestEnv) {
		err := env.is.SetInfraProvider(env.ctx, provider.Username, &provider)
		assert.Nil(t, err)

		resultPtr, err := env.is.GetInfraProvider(env.ctx, provider.Username)
		assert.Nil(t, err)
		assert.Equal(t, provider, *resultPtr, "infra provider should be equal")
	})

}

func TestInfraProviderList(t *testing.T) {
	lst := InfraProviderList{
		AllInfraProviders: []types.AccountKey{types.AccountKey("u1"), types.AccountKey("u2")},
	}

	runTest(t, func(env TestEnv) {
		err := env.is.SetInfraProviderList(env.ctx, &lst)
		assert.Nil(t, err)

		resultPtr, err := env.is.GetInfraProviderList(env.ctx)
		assert.Nil(t, err)
		assert.Equal(t, lst, *resultPtr, "infra provider list should be equal")
	})

}

//
// Test Environment setup
//

type TestEnv struct {
	is  InfraProviderStorage
	ctx sdk.Context
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	env := TestEnv{
		is:  NewInfraProviderStorage(TestKVStoreKey),
		ctx: getContext(),
	}
	fc(env)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	return sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
}
