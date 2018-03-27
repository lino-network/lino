package event

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("event")
)

func TestPostRewardEvent(t *testing.T) {
	em := NewEventManager(TestKVStoreKey)
	ctx := getContext()

	e1 := PostRewardEvent{
		PostID: 20,
	}
	e2 := DonateRewardEvent{
		DonateID: 20,
	}

	lst := EventList{}
	lst.Events = append(lst.Events, e1)
	lst.Events = append(lst.Events, e2)

	blockHeight := types.Height(100)
	lstKey := HeightToEventListKey(blockHeight)

	err := em.SetEventList(ctx, lstKey, &lst)
	assert.Nil(t, err)

	res := em.ExecuteEvents(ctx, lstKey)
	assert.Nil(t, res)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}
