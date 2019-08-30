package model

import (
	"testing"

	"github.com/lino-network/lino/types"

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
		GeneralMsgEMA: types.NewDecFromRat(311, 1),
		AppMsgEMA:     types.NewDecFromRat(200, 10),
		MaxMPS:        types.NewDecFromRat(12, 3),
	}
	err := bs.SetBandwidthInfo(ctx, &info)
	assert.Nil(t, err)

	resultPtr, err := bs.GetBandwidthInfo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, info, *resultPtr, "Bandwidth info should be equal")
}

func TestLastBlockInfo(t *testing.T) {
	bs := NewBandwidthStorage(TestKVStoreKey)
	ctx := getContext()

	info := LastBlockInfo{
		TotalMsgSignedByApp:  213123,
		TotalMsgSignedByUser: 0,
	}
	err := bs.SetLastBlockInfo(ctx, &info)
	assert.Nil(t, err)

	resultPtr, err := bs.GetLastBlockInfo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, info, *resultPtr, "LastBlockInfo info should be equal")
}
