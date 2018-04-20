package infra

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestInfraKVStoreKey = sdk.NewKVStoreKey("infra")
)

func setupTest(t *testing.T, height int64) (sdk.Context, InfraManager) {
	ctx := getContext(height)
	im := NewInfraManager(TestInfraKVStoreKey)
	return ctx, im
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestInfraKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height}, false, nil)
}
