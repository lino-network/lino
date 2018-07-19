package account

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateCoinReturnEvents(t *testing.T) {
	testCases := []struct {
		testName     string
		username     types.AccountKey
		times        int64
		interval     int64
		returnAmount types.Coin
		returnType   types.TransferDetailType
	}{
		{
			testName:     "normal return coin event",
			username:     "user1",
			times:        100,
			interval:     100,
			returnAmount: types.NewCoinFromInt64(100),
			returnType:   types.DelegationReturnCoin,
		},
		{
			testName:     "return coin is insufficient for each round",
			username:     "user1",
			times:        100,
			interval:     100,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.DelegationReturnCoin,
		},
		{
			testName:     "only one return event",
			username:     "user1",
			times:        1,
			interval:     100,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.DelegationReturnCoin,
		},
		{
			testName:     "no return interval",
			username:     "user1",
			times:        100,
			interval:     0,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.DelegationReturnCoin,
		},
		{
			testName:     "if return time is zero",
			username:     "user1",
			times:        0,
			interval:     0,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.DelegationReturnCoin,
		},
		{
			testName:     "return to different user",
			username:     "user2",
			times:        1,
			interval:     0,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.DelegationReturnCoin,
		},
		{
			testName:     "different return type",
			username:     "user2",
			times:        1,
			interval:     0,
			returnAmount: types.NewCoinFromInt64(1000),
			returnType:   types.VoteReturnCoin,
		},
	}

	for _, tc := range testCases {
		events, err := CreateCoinReturnEvents(
			tc.username, tc.times, tc.interval, tc.returnAmount, tc.returnType)
		if err != nil {
			t.Errorf("%s: failed to create coin return events, got err %v", tc.testName, err)
		}

		expectEvents := []types.Event{}
		for i := int64(0); i < tc.times; i++ {
			returnCoin := types.RatToCoin(sdk.NewRat(tc.returnAmount.ToInt64(), tc.times-i))

			event := ReturnCoinEvent{
				Username:   tc.username,
				Amount:     returnCoin,
				ReturnType: tc.returnType,
			}
			tc.returnAmount = tc.returnAmount.Minus(returnCoin)
			expectEvents = append(expectEvents, event)
		}

		if !assert.Equal(t, expectEvents, events) {
			t.Errorf("%s: diff events, got %v, want %v", tc.testName, events, expectEvents)
		}
	}
}

func TestReturnCoinEvent(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)

	createTestAccount(ctx, am, "user1")

	// Get the minimum time of this history slot
	baseTime := time.Now().Unix()
	testCases := []struct {
		testName             string
		event                ReturnCoinEvent
		atWhen               int64
		expectSaving         types.Coin
		expectBalanceHistory model.BalanceHistory
	}{
		{
			testName: "normal return case",
			event: ReturnCoinEvent{
				Username:   "user1",
				Amount:     types.NewCoinFromInt64(100),
				ReturnType: types.DelegationReturnCoin,
			},
			atWhen:       baseTime,
			expectSaving: types.NewCoinFromInt64(100).Plus(accParam.RegisterFee),
			expectBalanceHistory: model.BalanceHistory{
				[]model.Detail{
					{
						From:       "",
						DetailType: types.DelegationReturnCoin,
						Amount:     types.NewCoinFromInt64(100),
						CreatedAt:  baseTime,
					},
				},
			},
		},
		{
			testName: "return zero coin",
			event: ReturnCoinEvent{
				Username:   "user1",
				Amount:     types.NewCoinFromInt64(0),
				ReturnType: types.VoteReturnCoin,
			},
			atWhen:       baseTime,
			expectSaving: types.NewCoinFromInt64(100).Plus(accParam.RegisterFee),
			expectBalanceHistory: model.BalanceHistory{
				[]model.Detail{
					{
						From:       "",
						DetailType: types.DelegationReturnCoin,
						Amount:     types.NewCoinFromInt64(100),
						CreatedAt:  baseTime,
					},
					{
						From:       "",
						DetailType: types.VoteReturnCoin,
						Amount:     types.NewCoinFromInt64(0),
						CreatedAt:  baseTime,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		err := tc.event.Execute(ctx, am)
		if err != nil {
			t.Errorf("%s: failed to execute event, got err %v", tc.testName, err)
		}
		saving, err := am.GetSavingFromBank(ctx, tc.event.Username)
		if err != nil {
			t.Errorf("%s: failed to get saving from bank, got err %v", tc.testName, err)
		}
		if !saving.IsEqual(tc.expectSaving) {
			t.Errorf("%s: diff saving, got %v, want %v", tc.testName, saving, tc.expectSaving)
		}
	}
}
