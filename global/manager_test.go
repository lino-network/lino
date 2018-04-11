package global

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/global/model"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestGlobalKVStoreKey = sdk.NewKVStoreKey("global")
)

func InitGlobalManager(ctx sdk.Context, gm *GlobalManager) error {
	globalState := genesis.GlobalState{
		TotalLino:                10000,
		GrowthRate:               sdk.Rat{98, 1000},
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
		ConsumptionFrictionRate:  sdk.Rat{1, 100},
		FreezingPeriodHr:         24 * 7,
	}
	return gm.InitGlobalManager(ctx, globalState)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func setupTest(t *testing.T) (sdk.Context, *GlobalManager) {
	ctx := getContext()
	globalManager := NewGlobalManager(TestGlobalKVStoreKey)
	err := InitGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, globalManager
}

func TestRegisterEventAtHeight(t *testing.T) {
}

func TestTPS(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := time.Now().Unix()
	var initMaxTPS = sdk.NewRat(1000)

	cases := []struct {
		BaseTime            int64
		NextTime            int64
		NumOfTx             int32
		ExpectCurrentTPS    sdk.Rat
		ExpectMaxTPS        sdk.Rat
		ExpectCapacityRatio sdk.Rat
	}{
		{BaseTime: baseTime, NextTime: baseTime, NumOfTx: 0, ExpectCurrentTPS: initMaxTPS,
			ExpectMaxTPS: initMaxTPS, ExpectCapacityRatio: sdk.OneRat},
		{BaseTime: baseTime, NextTime: baseTime + 2, NumOfTx: 2, ExpectCurrentTPS: sdk.OneRat,
			ExpectMaxTPS: initMaxTPS, ExpectCapacityRatio: sdk.NewRat(1, 1000)},
		{BaseTime: baseTime, NextTime: baseTime + 1, NumOfTx: 1000, ExpectCurrentTPS: initMaxTPS,
			ExpectMaxTPS: initMaxTPS, ExpectCapacityRatio: sdk.OneRat},
		{BaseTime: baseTime, NextTime: baseTime + 2, NumOfTx: 2000, ExpectCurrentTPS: initMaxTPS,
			ExpectMaxTPS: initMaxTPS, ExpectCapacityRatio: sdk.OneRat},
		{BaseTime: baseTime, NextTime: baseTime + 2, NumOfTx: 3000, ExpectCurrentTPS: sdk.NewRat(1500),
			ExpectMaxTPS: sdk.NewRat(1500), ExpectCapacityRatio: sdk.OneRat},
		{BaseTime: baseTime, NextTime: baseTime + 2, NumOfTx: 2000, ExpectCurrentTPS: sdk.NewRat(1000),
			ExpectMaxTPS: sdk.NewRat(1500), ExpectCapacityRatio: sdk.NewRat(2, 3)},
	}
	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: cs.NextTime, NumTxs: cs.NumOfTx})
		err := gm.UpdateTPS(ctx, cs.BaseTime)
		assert.Nil(t, err)
		storage := model.NewGlobalStorage(TestGlobalKVStoreKey)
		tps, err := storage.GetTPS(ctx)
		assert.Equal(t, true, cs.ExpectCurrentTPS.Equal(tps.CurrentTPS))
		assert.Equal(t, true, cs.ExpectMaxTPS.Equal(tps.MaxTPS))
		ratio, err := gm.GetTPSCapacityRatio(ctx)
		assert.Nil(t, err)
		assert.Equal(t, true, cs.ExpectCapacityRatio.Equal(ratio))
	}
}
