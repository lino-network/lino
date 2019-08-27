package bandwidth

import (
	"testing"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/x/bandwidth/model"
	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var (
	testBandwidthKVStoreKey = sdk.NewKVStoreKey("bandwidth")
	testParamKVStoreKey     = sdk.NewKVStoreKey("param")
)

func setupTest(t *testing.T, height int64) (sdk.Context, BandwidthManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(testParamKVStoreKey)
	ph.InitParam(ctx)
	bandwidthManager := NewBandwidthManager(testBandwidthKVStoreKey, ph)
	bandwidthManager.InitGenesis(ctx)
	return ctx, bandwidthManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(testBandwidthKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Now()},
		false, log.NewNopLogger())
}

func checkCurBlockInfo(
	t *testing.T, ctx sdk.Context, testName string, bm BandwidthManager, curBlockInfo model.CurBlockInfo) {
	curBlockInfoPtr, err := bm.storage.GetCurBlockInfo(ctx)
	assert.Nil(t, err, "%s, failed to get cur block info, got err %v", testName, err)
	assert.Equal(t, curBlockInfo, *curBlockInfoPtr, "%s: diff cur block info, got %v, want %v", testName, *curBlockInfoPtr, curBlockInfo)
}
