package infra

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	testInfraKVStoreKey = sdk.NewKVStoreKey("infra")
	testParamKVStoreKey = sdk.NewKVStoreKey("param")
)

func setupTest(t *testing.T, height int64) (sdk.Context, InfraManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(testParamKVStoreKey)
	ph.InitParam(ctx)
	im := NewInfraManager(testInfraKVStoreKey, ph)
	return ctx, im
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(testInfraKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height}, false, log.NewNopLogger())
}
