package account

import (
	"math/big"
	"testing"
	"time"

	"github.com/lino-network/lino/x/account/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateCoinReturnEvents(t *testing.T) {
	testCases := []struct {
		testName     string
		username     types.AccountKey
		times        int64
		interval     int64
		returnAmount types.Coin
		returnType   types.BalanceHistoryDetailType
	}{
		{"normal return coin event", "user1", 100, 100,
			types.NewCoinFromInt64(100), types.DelegationReturnCoin},
		{"return coin is insufficient for each round", "user1", 100, 100,
			types.NewCoinFromInt64(1000), types.DelegationReturnCoin},
		{"only one return event", "user1", 1, 100,
			types.NewCoinFromInt64(1000), types.DelegationReturnCoin},
		{"no return interval", "user1", 100, 0,
			types.NewCoinFromInt64(1000), types.DelegationReturnCoin},
		{"if return time is zero", "user1", 0, 0,
			types.NewCoinFromInt64(1000), types.DelegationReturnCoin},
		{"return to different user", "user2", 1, 0,
			types.NewCoinFromInt64(1000), types.DelegationReturnCoin},
		{"different return type", "user2", 1, 0,
			types.NewCoinFromInt64(1000), types.VoteReturnCoin},
	}

	for _, tc := range testCases {
		events, err := CreateCoinReturnEvents(
			tc.username, tc.times, tc.interval, tc.returnAmount, tc.returnType)
		assert.Nil(t, err)
		expectEvents := []types.Event{}
		for i := int64(0); i < tc.times; i++ {
			returnCoin, err := types.RatToCoin(big.NewRat(tc.returnAmount.ToInt64(), tc.times-i))
			assert.Nil(t, err)
			event := ReturnCoinEvent{
				Username:   tc.username,
				Amount:     returnCoin,
				ReturnType: tc.returnType,
			}
			tc.returnAmount = tc.returnAmount.Minus(returnCoin)
			expectEvents = append(expectEvents, event)
		}
		assert.Equal(t, expectEvents, events)
	}
}

func TestReturnCoinEvent(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)

	createTestAccount(ctx, am, "user1")

	// Get the minimum time of this history slot
	baseTime := time.Now().Unix()
	baseTime = baseTime / accParam.BalanceHistoryIntervalTime * accParam.BalanceHistoryIntervalTime

	testCases := []struct {
		testName             string
		event                ReturnCoinEvent
		AtWhen               int64
		expectSaving         types.Coin
		expectBalanceHistory model.BalanceHistory
	}{
		{"normal return case", ReturnCoinEvent{
			Username:   "user1",
			Amount:     types.NewCoinFromInt64(100),
			ReturnType: types.DelegationReturnCoin,
		}, baseTime, types.NewCoinFromInt64(100).Plus(accParam.RegisterFee),
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						DetailType: types.DelegationReturnCoin,
						Amount:     types.NewCoinFromInt64(100),
						CreatedAt:  baseTime,
					},
				},
			},
		},
		{"return zero coin", ReturnCoinEvent{
			Username:   "user1",
			Amount:     types.NewCoinFromInt64(0),
			ReturnType: types.VoteReturnCoin,
		}, baseTime, types.NewCoinFromInt64(100).Plus(accParam.RegisterFee),
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						DetailType: types.DelegationReturnCoin,
						Amount:     types.NewCoinFromInt64(100),
						CreatedAt:  baseTime,
					},
					model.Detail{
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
		assert.Nil(t, err)
		saving, err := am.GetSavingFromBank(ctx, tc.event.Username)
		assert.Nil(t, err)
		assert.Equal(t, saving, tc.expectSaving)
	}
}
