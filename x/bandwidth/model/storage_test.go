package model

import (
	"testing"

	linotypes "github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("bandwidth")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
}

func TestBandwidthInfo(t *testing.T) {
	bs := NewBandwidthStorage(TestKVStoreKey)
	ctx := getContext()

	info := BandwidthInfo{
		GeneralMsgEMA: linotypes.NewDecFromRat(311, 1),
		AppMsgEMA:     linotypes.NewDecFromRat(200, 10),
		MaxMPS:        linotypes.NewDecFromRat(12, 3),
	}
	err := bs.SetBandwidthInfo(ctx, &info)
	assert.Nil(t, err)

	resultPtr, err := bs.GetBandwidthInfo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, info, *resultPtr, "Bandwidth info should be equal")
}

func TestBlockInfo(t *testing.T) {
	bs := NewBandwidthStorage(TestKVStoreKey)
	ctx := getContext()

	info := BlockInfo{
		TotalMsgSignedByApp:  213123,
		TotalMsgSignedByUser: 0,
		CurMsgFee:            linotypes.NewCoinFromInt64(int64(123)),
		CurU:                 sdk.NewDec(1),
	}
	err := bs.SetBlockInfo(ctx, &info)
	assert.Nil(t, err)

	resultPtr, err := bs.GetBlockInfo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, info, *resultPtr, "BlockInfo info should be equal")
}

func TestAppBandwidthInfo(t *testing.T) {
	bs := NewBandwidthStorage(TestKVStoreKey)
	ctx := getContext()

	info := AppBandwidthInfo{
		MaxBandwidthCredit: sdk.NewDec(1000),
		CurBandwidthCredit: sdk.NewDec(1000),
		MessagesInCurBlock: 100,
		LastRefilledAt:     1230,
		ExpectedMPS:        sdk.NewDec(200),
	}
	accName := linotypes.AccountKey("test")
	err := bs.SetAppBandwidthInfo(ctx, accName, &info)
	assert.Nil(t, err)

	resultPtr, err := bs.GetAppBandwidthInfo(ctx, accName)
	assert.Nil(t, err)
	assert.Equal(t, info, *resultPtr, "App bandwidth info should be equal")
}

func TestGetAllAppBandwidthInfo(t *testing.T) {
	bs := NewBandwidthStorage(TestKVStoreKey)
	ctx := getContext()
	appName1 := linotypes.AccountKey("app1")
	appName2 := linotypes.AccountKey("app2")
	info1 := AppBandwidthInfo{
		Username:           appName1,
		MaxBandwidthCredit: sdk.NewDec(1000),
		CurBandwidthCredit: sdk.NewDec(1000),
		MessagesInCurBlock: 100,
		LastRefilledAt:     1230,
		ExpectedMPS:        sdk.NewDec(200),
	}
	info2 := AppBandwidthInfo{
		Username:           appName2,
		MaxBandwidthCredit: sdk.NewDec(1000),
		CurBandwidthCredit: sdk.NewDec(2000),
		MessagesInCurBlock: 1300,
		LastRefilledAt:     1130,
		ExpectedMPS:        sdk.NewDec(100),
	}

	err := bs.SetAppBandwidthInfo(ctx, appName1, &info1)
	assert.Nil(t, err)
	err = bs.SetAppBandwidthInfo(ctx, appName2, &info2)
	assert.Nil(t, err)

	resultList, err := bs.GetAllAppBandwidthInfo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, []*AppBandwidthInfo{&info1, &info2}, resultList, "App bandwidth info should be equal")
}
