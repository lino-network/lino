package global

import (
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterEventAtHeight(t *testing.T) {
	gm := NewGlobalManager(TestKVStoreKey)
	ctx := getContext()
	err := initGlobalManager(t, ctx, gm)
	assert.Nil(t, err)

	globalProxy := NewGlobalProxy(&gm)
	err = globalProxy.RegisterEventAtHeight(ctx, types.Height(2), PostRewardEvent{})
	assert.Nil(t, err)
	heightEventList, err := gm.GetHeightEventList(ctx, HeightToEventListKey(types.Height(2)))
	assert.Nil(t, err)
	assert.Equal(t, HeightEventList{Events: []Event{PostRewardEvent{}}}, *heightEventList)
}
