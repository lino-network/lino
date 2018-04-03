package global

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("global")
)

func TestPostRewardEvent(t *testing.T) {
	gm := NewGlobalManager(TestKVStoreKey)
	ctx := getContext()

	lst := HeightEventList{}

	lstKey := HeightToEventListKey(100)

	err := gm.SetHeightEventList(ctx, lstKey, &lst)
	assert.Nil(t, err)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}
