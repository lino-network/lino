package manager

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"
)

func TestCreateCoinReturnEvents(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		testName     string
		username     types.AccountKey
		start        int64
		times        int64
		interval     int64
		returnAmount types.Coin
		returnType   types.TransferDetailType
		pool         types.PoolName
	}{
		{
			testName:     "normal return coin event",
			username:     "user1",
			start:        0,
			times:        100,
			interval:     100,
			returnAmount: types.NewCoinFromInt64(100),
			returnType:   types.DelegationReturnCoin,
			pool:         types.AccountVestingPool,
		},
		{
			testName:     "return coin is insufficient for each round",
			username:     "user1",
			start:        0,
			times:        100,
			interval:     100,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.DelegationReturnCoin,
			pool:         types.AccountVestingPool,
		},
		{
			testName:     "only one return event",
			username:     "user1",
			start:        0,
			times:        1,
			interval:     100,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.DelegationReturnCoin,
			pool:         types.AccountVestingPool,
		},
		{
			testName:     "no return interval",
			username:     "user1",
			start:        0,
			times:        100,
			interval:     0,
			returnAmount: types.NewCoinFromInt64(1000),
			pool:         types.AccountVestingPool,
		},
		{
			testName:     "if return time is zero",
			username:     "user1",
			start:        0,
			times:        0,
			interval:     0,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.DelegationReturnCoin,
			pool:         types.AccountVestingPool,
		},
		{
			testName:     "return to different user",
			username:     "user2",
			start:        0,
			times:        1,
			interval:     0,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.DelegationReturnCoin,
			pool:         types.AccountVestingPool,
		},
		{
			testName:     "different return type",
			username:     "user2",
			start:        0,
			times:        1,
			interval:     0,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.VoteReturnCoin,
			pool:         types.InflationValidatorPool,
		},
	}

	for _, tc := range testCases {
		events := CreateCoinReturnEvents(
			tc.username, tc.start, tc.interval, tc.times, tc.returnAmount, tc.returnType, tc.pool)
		expectEvents := []ReturnCoinEvent{}
		for i := int64(0); i < tc.times; i++ {
			returnAmount, _ := tc.returnAmount.ToInt64()
			returnCoin := types.DecToCoin(types.NewDecFromRat(returnAmount, tc.times-i))

			event := ReturnCoinEvent{
				Username:   tc.username,
				Amount:     returnCoin,
				ReturnType: tc.returnType,
				FromPool:   tc.pool,
				At:         tc.start + (i+1)*tc.interval,
			}
			tc.returnAmount = tc.returnAmount.Minus(returnCoin)
			expectEvents = append(expectEvents, event)
		}

		if !assert.Equal(expectEvents, events) {
			t.Errorf("%s: diff events, got %v, want %v", tc.testName, events, expectEvents)
		}
	}
}

func TestReturnCoinEvent(t *testing.T) {
	ctx, am := setupTest(t, 1)
	am.storage.SetPool(ctx, &model.Pool{
		Name:    types.AccountVestingPool,
		Balance: types.MustLinoToCoin("10000000000"),
	})
	accParam := am.paramHolder.GetAccountParam(ctx)

	createTestAccount(ctx, am, "user1")

	// Get the minimum time of this history slot
	baseTime := time.Now().Unix()
	testCases := []struct {
		testName     string
		event        ReturnCoinEvent
		atWhen       int64
		expectSaving types.Coin
	}{
		{
			testName: "normal return case",
			event: ReturnCoinEvent{
				Username:   "user1",
				Amount:     types.NewCoinFromInt64(100),
				ReturnType: types.DelegationReturnCoin,
				FromPool:   types.AccountVestingPool,
			},
			atWhen:       baseTime,
			expectSaving: types.NewCoinFromInt64(100).Plus(accParam.RegisterFee),
		},
		{
			testName: "return zero coin",
			event: ReturnCoinEvent{
				Username:   "user1",
				Amount:     types.NewCoinFromInt64(0),
				ReturnType: types.VoteReturnCoin,
				FromPool: types.AccountVestingPool,
			},
			atWhen:       baseTime,
			expectSaving: types.NewCoinFromInt64(100).Plus(accParam.RegisterFee),
		},
	}

	for _, tc := range testCases {
		err := tc.event.Execute(ctx, am)
		if err != nil {
			t.Errorf("%s: failed to execute event, got err %v", tc.testName, err)
		}
		saving, err := am.GetSavingFromUsername(ctx, tc.event.Username)
		if err != nil {
			t.Errorf("%s: failed to get saving from bank, got err %v", tc.testName, err)
		}
		if !saving.IsEqual(tc.expectSaving) {
			t.Errorf("%s: diff saving, got %v, want %v", tc.testName, saving, tc.expectSaving)
		}
	}
}
