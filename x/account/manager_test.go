package account

import (
	"fmt"
	"testing"
	"time"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestDoesAccountExist(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	if am.DoesAccountExist(ctx, types.AccountKey("user1")) {
		t.Error("TestDoesAccountExist: user1 has already existed")
	}

	createTestAccount(ctx, am, "user1")
	if !am.DoesAccountExist(ctx, types.AccountKey("user1")) {
		t.Error("TestDoesAccountExist: user1 should exist, but not")
	}
}

// https://github.com/lino-network/lino/issues/297
func TestAddCoinBundle(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	testUser := types.AccountKey("testUser")
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	coinDayParams, _ := am.paramHolder.GetCoinDayParam(ctx)

	baseTime := time.Now()
	baseTimeSlot := baseTime.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	baseTime = time.Unix(baseTimeSlot, 0)
	d1 := time.Duration(types.CoinDayRecordIntervalSec/10) * time.Second
	baseTime1 := baseTime.Add(d1)
	baseTime2 := baseTime1.Add(d1)
	baseTime3 := baseTime2.Add(d1)

	ctx = ctx.WithBlockHeader(abci.Header{Time: baseTime, Height: types.LinoBlockchainSecondUpdateHeight + 1})
	createTestAccount(ctx, am, string(testUser))

	testCases := []struct {
		testName                  string
		amount                    types.Coin
		atWhen                    time.Time
		expectBank                model.AccountBank
		expectPendingCoinDayQueue model.PendingCoinDayQueue
	}{
		{
			testName: "add coin to account's saving",
			amount:   c100,
			atWhen:   baseTime,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c100),
				CoinDay: accParam.RegisterFee,
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     c100,
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTimeSlot,
						EndTime:   baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c100,
					},
				},
			},
		},
		{
			testName: "add coin to same bucket at time 2",
			amount:   c100,
			atWhen:   baseTime2,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c200),
				CoinDay: accParam.RegisterFee,
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     c200,
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTimeSlot,
						EndTime:   baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c200,
					},
				},
			},
		},
		{
			testName: "add coin to same bucket at time 3",
			amount:   c100,
			atWhen:   baseTime3,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c300),
				CoinDay: accParam.RegisterFee,
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     c300,
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTimeSlot,
						EndTime:   baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c300,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: types.LinoBlockchainSecondUpdateHeight + 1, Time: tc.atWhen})
		err := am.AddSavingCoin(
			ctx, testUser, tc.amount, "", "", types.TransferIn)

		if err != nil {
			t.Errorf("%s: failed to add coin, got err: %v", tc.testName, err)
			return
		}
		checkBankKVByUsername(t, ctx, tc.testName, types.AccountKey(testUser), tc.expectBank)
		checkPendingCoinDay(t, ctx, tc.testName, types.AccountKey(testUser), tc.expectPendingCoinDayQueue)
	}
}

func TestAddCoin(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		t.Error("TestAddCoin: failed to get coin day param")
	}

	fromUser1, fromUser2, testUser :=
		types.AccountKey("fromUser1"), types.AccountKey("fromuser2"), types.AccountKey("testUser")

	baseTime := time.Now()
	baseTimeSlot := baseTime.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	d1 := time.Duration(coinDayParams.SecondsToRecoverCoinDay/2) * time.Second
	baseTime1 := baseTime.Add(d1)
	baseTime1Slot := baseTime1.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	d2 := time.Duration(coinDayParams.SecondsToRecoverCoinDay+1+types.CoinDayRecordIntervalSec) * time.Second
	baseTime2 := baseTime.Add(d2)
	baseTime2Slot := baseTime2.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	d3 := time.Duration(coinDayParams.SecondsToRecoverCoinDay+1) * time.Second
	baseTime3 := baseTime2.Add(d3)

	ctx = ctx.WithBlockHeader(abci.Header{Time: baseTime})
	createTestAccount(ctx, am, string(testUser))

	testCases := []struct {
		testName                  string
		amount                    types.Coin
		from                      types.AccountKey
		detailType                types.TransferDetailType
		memo                      string
		atWhen                    time.Time
		expectBank                model.AccountBank
		expectPendingCoinDayQueue model.PendingCoinDayQueue
	}{
		{
			testName:   "add coin to account's saving",
			amount:     c100,
			from:       fromUser1,
			detailType: types.TransferIn,
			memo:       "memo",
			atWhen:     baseTime,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c100),
				CoinDay: accParam.RegisterFee,
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     c100,
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTimeSlot,
						EndTime:   baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c100,
					},
				},
			},
		},
		{
			testName:   "add coin to exist account's saving while previous tx is still in pending queue",
			amount:     c100,
			from:       fromUser2,
			detailType: types.DonationIn,
			memo:       "permlink",
			atWhen:     baseTime1,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c200),
				CoinDay: accParam.RegisterFee,
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTime1Slot,
				TotalCoinDay:  sdk.NewDec(5000000),
				TotalCoin:     c100.Plus(c100),
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTimeSlot,
						EndTime:   baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c100,
					},
					{
						StartTime: baseTime1Slot,
						EndTime:   baseTime1Slot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c100,
					},
				},
			},
		},
		{
			testName:   "add coin to exist account's saving while previous tx just finished pending",
			amount:     c100,
			from:       "",
			detailType: types.ClaimReward,
			memo:       "",
			atWhen:     baseTime2,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c300),
				CoinDay: accParam.RegisterFee.Plus(c100),
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTime2Slot,
				TotalCoinDay:  types.NewDecFromRat(316250000, 63),
				TotalCoin:     c100.Plus(c100),
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTime1Slot,
						EndTime:   baseTime1Slot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c100,
					},
					{
						StartTime: baseTime2Slot,
						EndTime:   baseTime2Slot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c100,
					},
				},
			},
		},
		{
			testName:   "add coin is zero",
			amount:     c0,
			from:       "",
			detailType: types.DelegationReturnCoin,
			memo:       "",
			atWhen:     baseTime3,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c300),
				CoinDay: accParam.RegisterFee.Plus(c100),
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTime2Slot,
				TotalCoinDay:  types.NewDecFromRat(316250000, 63),
				TotalCoin:     c100.Plus(c100),
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTime1Slot,
						EndTime:   baseTime1Slot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c100,
					},
					{
						StartTime: baseTime2Slot,
						EndTime:   baseTime2Slot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      c100,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: tc.atWhen})
		err = am.AddSavingCoin(
			ctx, testUser, tc.amount, tc.from, tc.memo, tc.detailType)

		if err != nil {
			t.Errorf("%s: failed to add coin, got err: %v", tc.testName, err)
			return
		}
		checkBankKVByUsername(t, ctx, tc.testName, types.AccountKey(testUser), tc.expectBank)
		checkPendingCoinDay(t, ctx, tc.testName, types.AccountKey(testUser), tc.expectPendingCoinDayQueue)
	}
}

func TestMinusCoin(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		t.Error("TestMinusCoin: failed to get coin day param")
	}

	userWithSufficientSaving := types.AccountKey("user1")
	userWithLimitSaving := types.AccountKey("user3")
	fromUser, toUser := types.AccountKey("fromUser"), types.AccountKey("toUser")

	// Get the minimum time of this history slot
	baseTime := time.Now()
	baseTimeSlot := baseTime.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	// baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDay + 1
	// baseTime3 := baseTime + accParam.BalanceHistoryIntervalTime + 1

	ctx = ctx.WithBlockHeader(abci.Header{Time: baseTime})
	_, _, priv1 := createTestAccount(ctx, am, string(userWithSufficientSaving))
	_, _, priv3 := createTestAccount(ctx, am, string(userWithLimitSaving))

	err = am.AddSavingCoin(
		ctx, userWithSufficientSaving, accParam.RegisterFee, fromUser, "", types.TransferIn)
	if err != nil {
		t.Errorf("TestMinusCoin: failed to add saving coin, got err %v", err)
	}

	testCases := []struct {
		testName                  string
		fromUser                  types.AccountKey
		userPriv                  crypto.PrivKey
		expectErr                 sdk.Error
		amount                    types.Coin
		atWhen                    time.Time
		to                        types.AccountKey
		memo                      string
		detailType                types.TransferDetailType
		expectBank                model.AccountBank
		expectPendingCoinDayQueue model.PendingCoinDayQueue
	}{
		{
			testName:   "minus saving coin from user with sufficient saving",
			fromUser:   userWithSufficientSaving,
			userPriv:   priv1,
			expectErr:  nil,
			amount:     coin1,
			atWhen:     baseTime,
			to:         toUser,
			memo:       "memo",
			detailType: types.TransferOut,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(accParam.RegisterFee).Minus(coin1),
				CoinDay: accParam.RegisterFee,
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     accParam.RegisterFee.Minus(coin1),
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTimeSlot,
						EndTime:   baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      accParam.RegisterFee.Minus(coin1),
					},
				},
			},
		},
		{
			testName:   "minus saving coin from user with limit saving",
			fromUser:   userWithLimitSaving,
			userPriv:   priv3,
			expectErr:  ErrAccountSavingCoinNotEnough(),
			amount:     accParam.RegisterFee.Plus(accParam.RegisterFee),
			atWhen:     baseTime,
			to:         toUser,
			memo:       "memo",
			detailType: types.TransferOut,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee,
				CoinDay: accParam.RegisterFee,
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     accParam.RegisterFee,
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTimeSlot,
						EndTime:   baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      accParam.RegisterFee,
					}},
			},
		},
		{
			testName:   "minus saving coin exceeds the coin user hold",
			fromUser:   userWithLimitSaving,
			userPriv:   priv3,
			expectErr:  ErrAccountSavingCoinNotEnough(),
			amount:     c100,
			atWhen:     baseTime,
			to:         toUser,
			memo:       "memo",
			detailType: types.TransferOut,
			expectBank: model.AccountBank{
				Saving:  accParam.RegisterFee,
				CoinDay: accParam.RegisterFee,
			},
			expectPendingCoinDayQueue: model.PendingCoinDayQueue{
				LastUpdatedAt: baseTime.Unix(),
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     accParam.RegisterFee,
				PendingCoinDays: []model.PendingCoinDay{
					{
						StartTime: baseTime.Unix(),
						EndTime:   baseTime.Unix() + coinDayParams.SecondsToRecoverCoinDay,
						Coin:      accParam.RegisterFee,
					}},
			},
		},
	}
	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: tc.atWhen})
		err = am.MinusSavingCoin(ctx, tc.fromUser, tc.amount, tc.to, tc.memo, tc.detailType)

		if !assert.Equal(t, tc.expectErr, err) {
			t.Errorf("%s: diff err, got %v, want %v", tc.testName, err, tc.expectErr)
		}
		if tc.expectErr == nil {
			checkBankKVByUsername(t, ctx, tc.testName, tc.fromUser, tc.expectBank)
			checkPendingCoinDay(t, ctx, tc.testName, tc.fromUser, tc.expectPendingCoinDayQueue)
		}
	}
}

func TestMinusCoinWithFullCoinDay(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	// accParam, _ := am.paramHolder.GetAccountParam(ctx)

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		t.Error("TestMinusCoinWithFullCoinDay: failed to get coin day param")
	}

	user := types.AccountKey("user1")

	// Get the minimum time of this history slot
	baseTime := time.Now()
	baseTimeSlot := baseTime.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	beforeFullyChargedTimeSlot := baseTimeSlot - coinDayParams.SecondsToRecoverCoinDay
	afterFullyChargedTimeSlot := baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay
	halfTimeAfterFullyChargedTimeSlot := baseTimeSlot + coinDayParams.SecondsToRecoverCoinDay/2
	// baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDay + 1
	// baseTime3 := baseTime + accParam.BalanceHistoryIntervalTime + 1

	createTestAccount(ctx, am, string(user))
	testCases := []struct {
		testName            string
		bank                *model.AccountBank
		pendingCoinDayQueue *model.PendingCoinDayQueue
		minusAmount         types.Coin
		atWhen              time.Time
		expectCoinDay       types.Coin
	}{
		{
			testName: "minus all saving coin while pending coin day is empty",
			bank: &model.AccountBank{
				Saving:  coin1,
				CoinDay: coin1,
			},
			pendingCoinDayQueue: &model.PendingCoinDayQueue{
				LastUpdatedAt:   baseTimeSlot,
				TotalCoinDay:    sdk.ZeroDec(),
				TotalCoin:       coin0,
				PendingCoinDays: []model.PendingCoinDay{},
			},
			minusAmount:   coin1,
			atWhen:        baseTime,
			expectCoinDay: coin0,
		},
		{
			testName: "minus all saving coin with pending coin day queue",
			bank: &model.AccountBank{
				Saving:  coin2,
				CoinDay: coin1,
			},
			pendingCoinDayQueue: &model.PendingCoinDayQueue{
				LastUpdatedAt: beforeFullyChargedTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     coin1,
				PendingCoinDays: []model.PendingCoinDay{
					model.PendingCoinDay{
						StartTime: beforeFullyChargedTimeSlot,
						EndTime:   baseTimeSlot,
						Coin:      coin1,
					},
				},
			},
			minusAmount:   coin2,
			atWhen:        baseTime,
			expectCoinDay: coin0,
		},
		{
			testName: "minus saving coin with full coin day",
			bank: &model.AccountBank{
				Saving:  coin2,
				CoinDay: coin1,
			},
			pendingCoinDayQueue: &model.PendingCoinDayQueue{
				LastUpdatedAt: baseTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     coin1,
				PendingCoinDays: []model.PendingCoinDay{
					model.PendingCoinDay{
						StartTime: baseTimeSlot,
						EndTime:   afterFullyChargedTimeSlot,
						Coin:      coin1,
					},
				},
			},
			minusAmount:   coin1,
			atWhen:        baseTime,
			expectCoinDay: coin0,
		},
		{
			testName: "minus coin with half charged coin day",
			bank: &model.AccountBank{
				Saving:  coin2,
				CoinDay: coin0,
			},
			pendingCoinDayQueue: &model.PendingCoinDayQueue{
				LastUpdatedAt: baseTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     coin2,
				PendingCoinDays: []model.PendingCoinDay{
					model.PendingCoinDay{
						StartTime: baseTimeSlot,
						EndTime:   afterFullyChargedTimeSlot,
						Coin:      coin2,
					},
				},
			},
			minusAmount:   coin1,
			atWhen:        time.Unix(halfTimeAfterFullyChargedTimeSlot, 0),
			expectCoinDay: coin0,
		},
		{
			testName: "minus coin with all charged coin day and half charged coin day",
			bank: &model.AccountBank{
				Saving:  coin4,
				CoinDay: coin1,
			},
			pendingCoinDayQueue: &model.PendingCoinDayQueue{
				LastUpdatedAt: baseTimeSlot,
				TotalCoinDay:  sdk.ZeroDec(),
				TotalCoin:     coin3,
				PendingCoinDays: []model.PendingCoinDay{
					model.PendingCoinDay{
						StartTime: baseTimeSlot,
						EndTime:   afterFullyChargedTimeSlot,
						Coin:      coin3,
					},
				},
			},
			minusAmount:   coin2,
			atWhen:        time.Unix(halfTimeAfterFullyChargedTimeSlot, 0),
			expectCoinDay: coin1,
		},
	}
	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: tc.atWhen})
		err := am.storage.SetBankFromAccountKey(ctx, user, tc.bank)
		if err != nil {
			t.Errorf("%s: failed to set bank", tc.testName)
		}
		am.storage.SetPendingCoinDayQueue(ctx, user, tc.pendingCoinDayQueue)
		_, err = am.MinusSavingCoinWithFullCoinDay(ctx, user, tc.minusAmount, user, "", types.TransferOut)
		if err != nil {
			t.Errorf("%s: failed to minus coin", tc.testName)
		}
		coinDay, err := am.GetCoinDay(ctx, user)
		if err != nil {
			t.Errorf("%s: failed to get coin day", tc.testName)
		}
		if !coinDay.IsEqual(tc.expectCoinDay) {
			t.Errorf("%s: diff coin day, got %v, expect %v", tc.testName, coinDay, tc.expectCoinDay)
		}
	}
}

func TestCreateAccountNormalCase(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	coinDayParam, _ := am.paramHolder.GetCoinDayParam(ctx)

	largeAmountRegisterFee := types.NewCoinFromInt64(150000 * types.Decimals)
	testCases := []struct {
		testName              string
		username              types.AccountKey
		registerFee           types.Coin
		expectFullCoinDayCoin types.Coin
		expectNumberOfTx      int64
	}{
		{
			testName:              "zero register fee",
			username:              types.AccountKey("test1"),
			registerFee:           types.NewCoinFromInt64(0),
			expectFullCoinDayCoin: types.NewCoinFromInt64(0),
			expectNumberOfTx:      0,
		},
		{
			testName:              "micro register fee",
			username:              types.AccountKey("test2"),
			registerFee:           types.NewCoinFromInt64(1),
			expectFullCoinDayCoin: types.NewCoinFromInt64(1),
			expectNumberOfTx:      1,
		},
		{
			testName:              "register fee less than full coin day coin limitation",
			username:              types.AccountKey("test3"),
			registerFee:           types.NewCoinFromInt64(1500),
			expectFullCoinDayCoin: types.NewCoinFromInt64(1500),
			expectNumberOfTx:      1,
		},
		{
			testName:              "register fee much than full coin day coin limitation",
			username:              types.AccountKey("test4"),
			registerFee:           types.NewCoinFromInt64(150000),
			expectFullCoinDayCoin: accParam.FirstDepositFullCoinDayLimit,
			expectNumberOfTx:      2,
		},
		{
			testName:              "register with large amount of coin",
			username:              types.AccountKey("test5"),
			registerFee:           largeAmountRegisterFee,
			expectFullCoinDayCoin: accParam.FirstDepositFullCoinDayLimit,
			expectNumberOfTx:      2,
		},
	}
	// normal test
	for _, tc := range testCases {
		assert.False(t, am.DoesAccountExist(ctx, tc.username))
		resetPriv := secp256k1.GenPrivKey()
		txPriv := secp256k1.GenPrivKey()
		appPriv := secp256k1.GenPrivKey()

		err := am.CreateAccount(
			ctx, accountReferrer, tc.username, resetPriv.PubKey(), txPriv.PubKey(),
			appPriv.PubKey(), tc.registerFee)
		if err != nil {
			t.Errorf("%v: failed to create account, got err %v", tc.testName, err)
		}

		assert.True(t, am.DoesAccountExist(ctx, tc.username))
		bank := model.AccountBank{
			Saving:  tc.registerFee,
			CoinDay: tc.expectFullCoinDayCoin,
		}
		checkBankKVByUsername(t, ctx, tc.testName, tc.username, bank)

		pendingCoinDayQueue :=
			model.PendingCoinDayQueue{
				TotalCoinDay: sdk.ZeroDec(),
				TotalCoin:    types.NewCoinFromInt64(0),
			}
		baseTime := ctx.BlockHeader().Time.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
		if tc.registerFee.IsGT(tc.expectFullCoinDayCoin) {
			pendingCoinDayQueue.TotalCoin = tc.registerFee.Minus(tc.expectFullCoinDayCoin)
			pendingCoinDayQueue.PendingCoinDays = []model.PendingCoinDay{{
				StartTime: baseTime,
				EndTime:   baseTime + coinDayParam.SecondsToRecoverCoinDay,
				Coin:      tc.registerFee.Minus(tc.expectFullCoinDayCoin),
			}}
			pendingCoinDayQueue.LastUpdatedAt = baseTime
		}

		checkPendingCoinDay(t, ctx, tc.testName, tc.username, pendingCoinDayQueue)
		accInfo := model.AccountInfo{
			Username:       tc.username,
			CreatedAt:      ctx.BlockHeader().Time.Unix(),
			ResetKey:       resetPriv.PubKey(),
			TransactionKey: txPriv.PubKey(),
			AppKey:         appPriv.PubKey(),
		}
		checkAccountInfo(t, ctx, tc.testName, tc.username, accInfo)
		accMeta := model.AccountMeta{
			LastActivityAt:       ctx.BlockHeader().Time.Unix(),
			LastReportOrUpvoteAt: ctx.BlockHeader().Time.Unix(),
			TransactionCapacity:  tc.expectFullCoinDayCoin,
		}
		checkAccountMeta(t, ctx, tc.testName, tc.username, accMeta)

		reward := model.Reward{
			TotalIncome:     types.NewCoinFromInt64(0),
			OriginalIncome:  types.NewCoinFromInt64(0),
			InflationIncome: types.NewCoinFromInt64(0),
			FrictionIncome:  types.NewCoinFromInt64(0),
			UnclaimReward:   types.NewCoinFromInt64(0),
		}
		checkAccountReward(t, ctx, tc.testName, tc.username, reward)
	}
}

func TestCreateAccountWithLargeRegisterFee(t *testing.T) {
	testName := "TestCreateAccountWithLargeRegisterFee"

	ctx, am, _ := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	resetPriv := secp256k1.GenPrivKey()
	txPriv := secp256k1.GenPrivKey()
	appPriv := secp256k1.GenPrivKey()

	accKey := types.AccountKey("accKey")

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		t.Errorf("%s: failed to get coin day param, got err %v", testName, err)
	}

	extraRegisterFee := types.NewCoinFromInt64(100 * types.Decimals)
	// normal test
	if am.DoesAccountExist(ctx, accKey) {
		t.Errorf("%s: account %v already exist", testName, accKey)
	}

	err = am.CreateAccount(
		ctx, accountReferrer, accKey, resetPriv.PubKey(), txPriv.PubKey(),
		appPriv.PubKey(), accParam.RegisterFee.Plus(extraRegisterFee))
	if err != nil {
		t.Errorf("%s: failed to create account, got err %v", testName, err)
	}

	assert.True(t, am.DoesAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee.Plus(extraRegisterFee),
		CoinDay: accParam.RegisterFee,
	}
	checkBankKVByUsername(t, ctx, testName, accKey, bank)

	baseTime := ctx.BlockHeader().Time.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	pendingCoinDayQueue := model.PendingCoinDayQueue{
		LastUpdatedAt: baseTime,
		TotalCoinDay:  sdk.ZeroDec(),
		TotalCoin:     extraRegisterFee,
		PendingCoinDays: []model.PendingCoinDay{
			{
				StartTime: baseTime,
				EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDay,
				Coin:      extraRegisterFee,
			},
		},
	}
	checkPendingCoinDay(t, ctx, testName, accKey, pendingCoinDayQueue)

	accInfo := model.AccountInfo{
		Username:       accKey,
		CreatedAt:      ctx.BlockHeader().Time.Unix(),
		ResetKey:       resetPriv.PubKey(),
		TransactionKey: txPriv.PubKey(),
		AppKey:         appPriv.PubKey(),
	}
	checkAccountInfo(t, ctx, testName, accKey, accInfo)

	accMeta := model.AccountMeta{
		LastActivityAt:       ctx.BlockHeader().Time.Unix(),
		LastReportOrUpvoteAt: ctx.BlockHeader().Time.Unix(),
		TransactionCapacity:  accParam.RegisterFee,
	}
	checkAccountMeta(t, ctx, testName, accKey, accMeta)

	reward := model.Reward{
		TotalIncome:     types.NewCoinFromInt64(0),
		OriginalIncome:  types.NewCoinFromInt64(0),
		InflationIncome: types.NewCoinFromInt64(0),
		FrictionIncome:  types.NewCoinFromInt64(0),
		UnclaimReward:   types.NewCoinFromInt64(0),
	}
	checkAccountReward(t, ctx, testName, accKey, reward)
}

func TestInvalidCreateAccount(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	priv1 := secp256k1.GenPrivKey()
	priv2 := secp256k1.GenPrivKey()

	accKey1 := types.AccountKey("accKey1")
	accKey2 := types.AccountKey("accKey2")

	testCases := []struct {
		testName        string
		username        types.AccountKey
		privKey         crypto.PrivKey
		registerDeposit types.Coin
		expectErr       sdk.Error
	}{
		{
			testName:        "register user with sufficient saving coin",
			username:        accKey1,
			privKey:         priv1,
			registerDeposit: accParam.RegisterFee,
			expectErr:       nil,
		},
		{
			testName:        "username already took",
			username:        accKey1,
			privKey:         priv1,
			registerDeposit: accParam.RegisterFee,
			expectErr:       ErrAccountAlreadyExists(accKey1),
		},
		{
			testName:        "username already took with different private key",
			username:        accKey1,
			privKey:         priv2,
			registerDeposit: accParam.RegisterFee,
			expectErr:       ErrAccountAlreadyExists(accKey1),
		},
		{
			testName:        "register the same private key",
			username:        accKey2,
			privKey:         priv1,
			registerDeposit: accParam.RegisterFee,
			expectErr:       nil,
		},
	}
	for _, tc := range testCases {
		err := am.CreateAccount(
			ctx, accountReferrer, tc.username, tc.privKey.PubKey(),
			secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), tc.registerDeposit)
		if !assert.Equal(t, tc.expectErr, err) {
			t.Errorf("%s: diff err, got %v, want %v", tc.testName, err, tc.expectErr)
		}
	}
}

func TestUpdateJSONMeta(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)

	accKey := types.AccountKey("accKey")
	createTestAccount(ctx, am, string(accKey))

	testCases := []struct {
		testName string
		username types.AccountKey
		JSONMeta string
	}{
		{
			testName: "normal update",
			username: accKey,
			JSONMeta: "{'link':'https://lino.network'}",
		},
	}
	for _, tc := range testCases {
		err := am.UpdateJSONMeta(ctx, tc.username, tc.JSONMeta)
		if err != nil {
			t.Errorf("%s: failed to update json meta, got err %v", tc.testName, err)
		}

		accMeta, err := am.storage.GetMeta(ctx, tc.username)
		if err != nil {
			t.Errorf("%s: failed to get meta, got err %v", tc.testName, err)
		}
		if tc.JSONMeta != accMeta.JSONMeta {
			t.Errorf("%s: diff json meta, got %v, want %v", tc.testName, accMeta.JSONMeta, tc.JSONMeta)
		}
	}
}

func TestCoinDayByAccountKey(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	accKey := types.AccountKey("accKey")

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		t.Errorf("TestCoinDayByAccountKey: failed to get coin day param, got err %v", err)
	}

	totalCoinDaysSec := coinDayParams.SecondsToRecoverCoinDay
	registerFee, _ := accParam.RegisterFee.ToInt64()
	doubleRegisterFee := types.NewCoinFromInt64(registerFee * 2)
	halfRegisterFee := types.NewCoinFromInt64(registerFee / 2)

	baseTime := ctx.BlockHeader().Time
	baseTime2 := baseTime.Add(time.Duration((totalCoinDaysSec)+types.CoinDayRecordIntervalSec) * time.Second)

	createTestAccount(ctx, am, string(accKey))

	testCases := []struct {
		testName            string
		isAdd               bool
		coin                types.Coin
		atWhen              time.Time
		expectSavingBalance types.Coin
		expectCoinDay       types.Coin
		expectCoinDayInBank types.Coin
		expectNumOfTx       int64
	}{
		{
			testName:            "add coin before charging first coin",
			isAdd:               true,
			coin:                accParam.RegisterFee,
			atWhen:              baseTime.Add(time.Duration((totalCoinDaysSec/registerFee)/2) * time.Second),
			expectSavingBalance: doubleRegisterFee,
			expectCoinDay:       accParam.RegisterFee,
			expectCoinDayInBank: accParam.RegisterFee,
			expectNumOfTx:       2,
		},
		{
			testName:            "check first coin",
			isAdd:               true,
			coin:                coin0,
			atWhen:              baseTime.Add(time.Duration((totalCoinDaysSec/registerFee)/2+1) * time.Second),
			expectSavingBalance: doubleRegisterFee,
			expectCoinDay:       accParam.RegisterFee,
			expectCoinDayInBank: accParam.RegisterFee,
			expectNumOfTx:       2,
		},
		{
			testName:            "check both transactions fully charged",
			isAdd:               true,
			coin:                coin0,
			atWhen:              baseTime2,
			expectSavingBalance: doubleRegisterFee,
			expectCoinDay:       doubleRegisterFee,
			expectCoinDayInBank: doubleRegisterFee,
			expectNumOfTx:       2,
		},
		{
			testName:            "withdraw half deposit",
			isAdd:               false,
			coin:                accParam.RegisterFee,
			atWhen:              baseTime2,
			expectSavingBalance: accParam.RegisterFee,
			expectCoinDay:       accParam.RegisterFee,
			expectCoinDayInBank: accParam.RegisterFee,
			expectNumOfTx:       3,
		},
		{
			testName:            "charge again",
			isAdd:               true,
			coin:                accParam.RegisterFee,
			atWhen:              baseTime2,
			expectSavingBalance: doubleRegisterFee,
			expectCoinDay:       accParam.RegisterFee,
			expectCoinDayInBank: accParam.RegisterFee,
			expectNumOfTx:       4,
		},
		{
			testName:            "withdraw half deposit while the last transaction is still charging",
			isAdd:               false,
			coin:                halfRegisterFee,
			atWhen:              baseTime2.Add(time.Duration(totalCoinDaysSec/2+1) * time.Second),
			expectSavingBalance: accParam.RegisterFee.Plus(halfRegisterFee),
			expectCoinDay:       accParam.RegisterFee.Plus(types.NewCoinFromInt64(registerFee / 4)),
			expectCoinDayInBank: accParam.RegisterFee,
			expectNumOfTx:       5,
		},
		{
			testName:            "withdraw last transaction which is still charging",
			isAdd:               false,
			coin:                halfRegisterFee,
			atWhen:              baseTime2.Add(time.Duration(totalCoinDaysSec/2+1) * time.Second),
			expectSavingBalance: accParam.RegisterFee,
			expectCoinDay:       accParam.RegisterFee,
			expectCoinDayInBank: accParam.RegisterFee,
			expectNumOfTx:       6,
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: tc.atWhen})
		if tc.isAdd {
			err := am.AddSavingCoin(ctx, accKey, tc.coin, "", "", types.TransferIn)
			if err != nil {
				t.Errorf("%s: failed to add saving coin, got err %v", tc.testName, err)
			}
		} else {
			err := am.MinusSavingCoin(ctx, accKey, tc.coin, "", "", types.TransferOut)
			if err != nil {
				t.Errorf("%s: failed to minus saving coin, got err %v", tc.testName, err)
			}
		}
		coin, err := am.GetCoinDay(ctx, accKey)
		if err != nil {
			t.Errorf("%s: failed to get coin day, got err %v", tc.testName, err)
		}

		if !tc.expectCoinDay.IsEqual(coin) {
			t.Errorf("%s: diff coin day, got %v, want %v", tc.testName, coin, tc.expectCoinDay)
			return
		}

		bank := model.AccountBank{
			Saving:  tc.expectSavingBalance,
			CoinDay: tc.expectCoinDayInBank,
		}
		checkBankKVByUsername(t, ctx, tc.testName, accKey, bank)
	}
}

func TestCoinDaySize(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	coinDayParam, _ := am.paramHolder.GetCoinDayParam(ctx)
	accKey := types.AccountKey("accKey")
	createTestAccount(ctx, am, string(accKey))
	baseTime := ctx.BlockHeader().Time.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	for i := baseTime; i < baseTime+coinDayParam.SecondsToRecoverCoinDay+types.CoinDayRecordIntervalSec*2; i += 40 {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(i, 0)})
		err := am.AddSavingCoin(ctx, accKey, types.NewCoinFromInt64(1), accKey, "", types.TransferIn)
		if err != nil {
			t.Errorf("%s: failed to add coin, got err %v", "TestCoinDaySize", err)
		}
	}
	pendingCoinDayQueue, err := am.storage.GetPendingCoinDayQueue(ctx, accKey)
	if err != nil {
		t.Errorf("%s: failed to add coin, got err %v", "TestCoinDaySize", err)
	}
	assert.Equal(t, 504, len(pendingCoinDayQueue.PendingCoinDays))
}

func TestAddIncomeAndReward(t *testing.T) {
	testName := "TestAddIncomeAndReward"

	ctx, am, _ := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	accKey := types.AccountKey("accKey")

	createTestAccount(ctx, am, string(accKey))

	err := am.AddIncomeAndReward(ctx, accKey, c500, c200, c300, "donor1", "postAutho1", "post1")
	if err != nil {
		t.Errorf("%s: failed to add income and reward, got err %v", testName, err)
	}

	reward := model.Reward{
		TotalIncome:     c300,
		OriginalIncome:  c200,
		FrictionIncome:  c200,
		InflationIncome: c300,
		UnclaimReward:   c300,
	}
	checkAccountReward(t, ctx, testName, accKey, reward)

	err = am.AddIncomeAndReward(ctx, accKey, c500, c300, c200, "donor2", "postAuthor1", "post1")
	if err != nil {
		t.Errorf("%s: failed to add income and reward again, got err %v", testName, err)
	}

	reward = model.Reward{
		TotalIncome:     c500,
		OriginalIncome:  c500,
		FrictionIncome:  c500,
		InflationIncome: c500,
		UnclaimReward:   c500,
	}
	checkAccountReward(t, ctx, testName, accKey, reward)

	err = am.AddDirectDeposit(ctx, accKey, c500)
	if err != nil {
		t.Errorf("%s: failed to add direct deposit, got err %v", testName, err)
	}

	reward = model.Reward{
		TotalIncome:     c1000,
		OriginalIncome:  c1000,
		FrictionIncome:  c500,
		InflationIncome: c500,
		UnclaimReward:   c500,
	}
	checkAccountReward(t, ctx, testName, accKey, reward)
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee,
		CoinDay: accParam.RegisterFee,
	}
	checkBankKVByUsername(t, ctx, testName, accKey, bank)

	err = am.ClaimReward(ctx, accKey)
	if err != nil {
		t.Errorf("%s: failed to add claim reward, got err %v", testName, err)
	}

	bank.Saving = accParam.RegisterFee.Plus(c500)
	checkBankKVByUsername(t, ctx, testName, accKey, bank)

	reward = model.Reward{
		TotalIncome:     c1000,
		OriginalIncome:  c1000,
		FrictionIncome:  c500,
		InflationIncome: c500,
		UnclaimReward:   c0,
	}
	checkAccountReward(t, ctx, testName, accKey, reward)
}

func TestCheckUserTPSCapacity(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	bandwidthParams, err := am.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		t.Errorf("TestCheckUserTPSCapacity: failed to get bandwidth param, got err %v", err)
	}
	virtualCoinAmount, _ := bandwidthParams.VirtualCoin.ToInt64()
	secondsToRecoverBandwidth := bandwidthParams.SecondsToRecoverBandwidth

	baseTime := ctx.BlockHeader().Time

	createTestAccount(ctx, am, string(accKey))
	err = am.AddSavingCoin(ctx, accKey, c100, "", "", types.TransferIn)
	if err != nil {
		t.Errorf("TestCheckUserTPSCapacity: failed to add saving coin, got err %v", err)
	}

	accStorage := model.NewAccountStorage(testAccountKVStoreKey)
	err = accStorage.SetPendingCoinDayQueue(
		ctx, accKey, &model.PendingCoinDayQueue{})
	if err != nil {
		t.Errorf("TestCheckUserTPSCapacity: failed to set pending coin day queue, got err %v", err)
	}

	testCases := []struct {
		testName             string
		tpsCapacityRatio     sdk.Dec
		userCoinDay          types.Coin
		lastActivity         int64
		lastCapacity         types.Coin
		currentTime          time.Time
		expectResult         sdk.Error
		expectRemainCapacity types.Coin
	}{
		{
			testName:             "tps capacity not enough",
			tpsCapacityRatio:     types.NewDecFromRat(1, 10),
			userCoinDay:          types.NewCoinFromInt64(10 * types.Decimals),
			lastActivity:         baseTime.Unix(),
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime,
			expectResult:         ErrAccountTPSCapacityNotEnough(accKey),
			expectRemainCapacity: types.NewCoinFromInt64(0),
		},
		{
			testName:             " 1/10 capacity ratio",
			tpsCapacityRatio:     types.NewDecFromRat(1, 10),
			userCoinDay:          types.NewCoinFromInt64(10 * types.Decimals),
			lastActivity:         baseTime.Unix(),
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth) * time.Second),
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(990000).Plus(bandwidthParams.VirtualCoin),
		},
		{
			testName:             " 1/2 capacity ratio",
			tpsCapacityRatio:     types.NewDecFromRat(1, 2),
			userCoinDay:          types.NewCoinFromInt64(10 * types.Decimals),
			lastActivity:         baseTime.Unix(),
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth) * time.Second),
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(950000).Plus(bandwidthParams.VirtualCoin),
		},
		{
			testName:             " 1/1 capacity ratio",
			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
			userCoinDay:          types.NewCoinFromInt64(10 * types.Decimals),
			lastActivity:         baseTime.Unix(),
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth) * time.Second),
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(9 * types.Decimals).Plus(bandwidthParams.VirtualCoin),
		},
		{
			testName:             " 1/1 capacity ratio with virtual coin remaining",
			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
			userCoinDay:          types.NewCoinFromInt64(1 * types.Decimals),
			lastActivity:         baseTime.Unix(),
			lastCapacity:         types.NewCoinFromInt64(10 * types.Decimals),
			currentTime:          baseTime,
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(1 * types.Decimals),
		},
		{
			testName:             " 1/1 capacity ratio with 1 coin day and 0 remaining",
			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
			userCoinDay:          types.NewCoinFromInt64(1 * types.Decimals),
			lastActivity:         baseTime.Unix(),
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth/2) * time.Second),
			expectResult:         nil,
			expectRemainCapacity: coin0,
		},
		{
			testName:             " transaction capacity not enough",
			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
			userCoinDay:          types.NewCoinFromInt64(0 * types.Decimals),
			lastActivity:         baseTime.Unix(),
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth/2) * time.Second),
			expectResult:         ErrAccountTPSCapacityNotEnough(accKey),
			expectRemainCapacity: coin0,
		},
		{
			testName:             " transaction capacity without coin day",
			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
			userCoinDay:          types.NewCoinFromInt64(0 * types.Decimals),
			lastActivity:         baseTime.Unix(),
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth) * time.Second),
			expectResult:         nil,
			expectRemainCapacity: coin0,
		},
		{
			testName:             " 1/2 capacity ratio with half virtual coin remaining",
			tpsCapacityRatio:     types.NewDecFromRat(1, 2),
			userCoinDay:          types.NewCoinFromInt64(1 * types.Decimals),
			lastActivity:         baseTime.Unix(),
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth/2) * time.Second),
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(virtualCoinAmount / 2),
		},
		{
			testName:             " 1/1 capacity ratio with virtual coin remaining and base time",
			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
			userCoinDay:          types.NewCoinFromInt64(1 * types.Decimals),
			lastActivity:         0,
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime,
			expectResult:         nil,
			expectRemainCapacity: bandwidthParams.VirtualCoin,
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: tc.currentTime})
		bank := &model.AccountBank{
			Saving:  tc.userCoinDay,
			CoinDay: tc.userCoinDay,
		}
		err := accStorage.SetBankFromAccountKey(ctx, accKey, bank)
		if err != nil {
			t.Errorf("%s: failed to set bank, got err %v", tc.testName, err)
		}

		meta := &model.AccountMeta{
			LastActivityAt:      tc.lastActivity,
			TransactionCapacity: tc.lastCapacity,
		}
		err = accStorage.SetMeta(ctx, accKey, meta)
		if err != nil {
			t.Errorf("%s: failed to set meta, got err %v", tc.testName, err)
		}

		err = am.CheckUserTPSCapacity(ctx, accKey, tc.tpsCapacityRatio)
		if !assert.Equal(t, tc.expectResult, err) {
			t.Errorf("%s: diff tps capacity, got %v, want %v", tc.testName, err, tc.expectResult)
		}

		accMeta := model.AccountMeta{
			LastActivityAt:      ctx.BlockHeader().Time.Unix(),
			TransactionCapacity: tc.expectRemainCapacity,
		}
		if tc.expectResult != nil {
			accMeta.LastActivityAt = tc.lastActivity
		}
		checkAccountMeta(t, ctx, tc.testName, accKey, accMeta)
	}
}

func TestCheckAuthenticatePubKeyOwner(t *testing.T) {
	testName := "TestCheckAuthenticatePubKeyOwner"

	ctx, am, _ := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	user1 := types.AccountKey("user1")
	appPermissionUser := types.AccountKey("user2")
	preAuthPermissionUser := types.AccountKey("user3")
	unauthUser := types.AccountKey("user4")
	resetKey := secp256k1.GenPrivKey()
	transactionKey := secp256k1.GenPrivKey()
	appKey := secp256k1.GenPrivKey()
	am.CreateAccount(
		ctx, accountReferrer, user1, resetKey.PubKey(), transactionKey.PubKey(),
		appKey.PubKey(), accParam.RegisterFee)

	_, unauthTxPriv, authAppPriv := createTestAccount(ctx, am, string(appPermissionUser))
	_, authTxPriv, unauthAppPriv := createTestAccount(ctx, am, string(preAuthPermissionUser))
	_, unauthPriv1, unauthPriv2 := createTestAccount(ctx, am, string(unauthUser))

	err := am.AuthorizePermission(ctx, user1, appPermissionUser, 100, types.AppPermission, types.NewCoinFromInt64(0))
	if err != nil {
		t.Errorf("%s: failed to authorize app permission, got err %v", testName, err)
	}

	preAuthAmount := types.NewCoinFromInt64(100)
	err = am.AuthorizePermission(ctx, user1, preAuthPermissionUser, 100, types.PreAuthorizationPermission, preAuthAmount)
	if err != nil {
		t.Errorf("%s: failed to authorize preauth permission, got err %v", testName, err)
	}
	baseTime := ctx.BlockHeader().Time

	testCases := []struct {
		testName           string
		checkUser          types.AccountKey
		checkPubKey        crypto.PubKey
		atWhen             time.Time
		amount             types.Coin
		permission         types.Permission
		expectUser         types.AccountKey
		expectResult       sdk.Error
		expectGrantPubKeys []*model.GrantPermission
	}{
		{
			testName:           "check user's reset key",
			checkUser:          user1,
			checkPubKey:        resetKey.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(0),
			permission:         types.ResetPermission,
			expectUser:         user1,
			expectResult:       nil,
			expectGrantPubKeys: nil,
		},
		{
			testName:           "check user's transaction key",
			checkUser:          user1,
			checkPubKey:        transactionKey.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(0),
			permission:         types.TransactionPermission,
			expectUser:         user1,
			expectResult:       nil,
			expectGrantPubKeys: nil,
		},
		{
			testName:           "check user's app key",
			checkUser:          user1,
			checkPubKey:        appKey.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(0),
			permission:         types.AppPermission,
			expectUser:         user1,
			expectResult:       nil,
			expectGrantPubKeys: nil,
		},
		{
			testName:           "user's transaction key can authorize grant app permission",
			checkUser:          user1,
			checkPubKey:        transactionKey.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(0),
			permission:         types.GrantAppPermission,
			expectUser:         user1,
			expectResult:       nil,
			expectGrantPubKeys: nil,
		},
		{
			testName:           "user's transaction key can authorize app permission",
			checkUser:          user1,
			checkPubKey:        transactionKey.PubKey(),
			atWhen:             baseTime,
			permission:         types.AppPermission,
			expectUser:         user1,
			expectResult:       nil,
			expectGrantPubKeys: nil,
		},
		{
			testName:           "check user's transaction key can't authorize reset permission",
			checkUser:          user1,
			checkPubKey:        transactionKey.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(0),
			permission:         types.ResetPermission,
			expectUser:         user1,
			expectResult:       ErrCheckResetKey(),
			expectGrantPubKeys: nil,
		},
		{
			testName:           "check user's app key can authorize grant app permission",
			checkUser:          user1,
			checkPubKey:        appKey.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(0),
			permission:         types.GrantAppPermission,
			expectUser:         user1,
			expectResult:       nil,
			expectGrantPubKeys: nil,
		},
		{
			testName:           "check user's app key can't authorize transaction permission",
			checkUser:          user1,
			checkPubKey:        appKey.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(0),
			permission:         types.TransactionPermission,
			expectUser:         user1,
			expectResult:       ErrCheckTransactionKey(),
			expectGrantPubKeys: nil,
		},
		{
			testName:           "check user's app key can't authorize reset permission",
			checkUser:          user1,
			checkPubKey:        appKey.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(0),
			permission:         types.ResetPermission,
			expectUser:         user1,
			expectResult:       ErrCheckResetKey(),
			expectGrantPubKeys: nil,
		},
		{
			testName:     "check app pubkey of user with app permission",
			checkUser:    user1,
			checkPubKey:  authAppPriv.PubKey(),
			atWhen:       baseTime,
			amount:       types.NewCoinFromInt64(0),
			permission:   types.AppPermission,
			expectUser:   appPermissionUser,
			expectResult: nil,
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    appPermissionUser,
					Permission: types.AppPermission,
					CreatedAt:  baseTime.Unix(),
					ExpiresAt:  baseTime.Unix() + 100,
					Amount:     types.NewCoinFromInt64(0),
				},
			},
		},
		{
			testName:           "check transaction pubkey of user with app permission",
			checkUser:          user1,
			checkPubKey:        unauthTxPriv.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(0),
			permission:         types.PreAuthorizationPermission,
			expectUser:         "",
			expectResult:       nil,
			expectGrantPubKeys: nil,
		},
		{
			testName:           "check unauthorized user app pubkey",
			checkUser:          user1,
			checkPubKey:        unauthPriv2.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(10),
			permission:         types.AppPermission,
			expectUser:         "",
			expectResult:       ErrCheckAuthenticatePubKeyOwner(user1),
			expectGrantPubKeys: nil,
		},
		{
			testName:           "check unauthorized user transaction pubkey",
			checkUser:          user1,
			checkPubKey:        unauthPriv1.PubKey(),
			atWhen:             baseTime,
			amount:             types.NewCoinFromInt64(10),
			permission:         types.PreAuthorizationPermission,
			expectUser:         "",
			expectResult:       ErrCheckAuthenticatePubKeyOwner(user1),
			expectGrantPubKeys: nil,
		},
		{
			testName:     "check transaction pubkey of user with preauthorization permission",
			checkUser:    user1,
			checkPubKey:  authTxPriv.PubKey(),
			atWhen:       baseTime,
			amount:       types.NewCoinFromInt64(10),
			permission:   types.PreAuthorizationPermission,
			expectUser:   preAuthPermissionUser,
			expectResult: nil,
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    preAuthPermissionUser,
					Permission: types.PreAuthorizationPermission,
					CreatedAt:  baseTime.Unix(),
					ExpiresAt:  baseTime.Unix() + 100,
					Amount:     preAuthAmount.Minus(types.NewCoinFromInt64(10)),
				},
			},
		},
		{
			testName:     "check app pubkey of user with preauthorization permission",
			checkUser:    user1,
			checkPubKey:  unauthAppPriv.PubKey(),
			atWhen:       baseTime,
			amount:       types.NewCoinFromInt64(10),
			permission:   types.AppPermission,
			expectUser:   preAuthPermissionUser,
			expectResult: ErrCheckAuthenticatePubKeyOwner(user1),
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    preAuthPermissionUser,
					Permission: types.PreAuthorizationPermission,
					CreatedAt:  baseTime.Unix(),
					ExpiresAt:  baseTime.Unix() + 100,
					Amount:     preAuthAmount.Minus(types.NewCoinFromInt64(10)),
				},
			},
		},
		{
			testName:    "check amount exceeds preauthorization limitation",
			checkUser:   user1,
			checkPubKey: authTxPriv.PubKey(),
			atWhen:      baseTime,
			amount:      preAuthAmount,
			permission:  types.PreAuthorizationPermission,
			expectUser:  preAuthPermissionUser,
			expectResult: ErrPreAuthAmountInsufficient(
				preAuthPermissionUser, preAuthAmount.Minus(types.NewCoinFromInt64(10)),
				preAuthAmount),
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    preAuthPermissionUser,
					Permission: types.PreAuthorizationPermission,
					CreatedAt:  baseTime.Unix(),
					ExpiresAt:  baseTime.Unix() + 100,
					Amount:     preAuthAmount.Minus(types.NewCoinFromInt64(10)),
				},
			},
		},
		{
			testName:     "check grant app key can't sign grant app msg",
			checkUser:    user1,
			checkPubKey:  authAppPriv.PubKey(),
			atWhen:       baseTime,
			permission:   types.GrantAppPermission,
			expectUser:   appPermissionUser,
			expectResult: ErrCheckGrantAppKey(),
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    appPermissionUser,
					Permission: types.AppPermission,
					CreatedAt:  baseTime.Unix(),
					ExpiresAt:  baseTime.Unix() + 100,
					Amount:     types.NewCoinFromInt64(0),
				},
			},
		},
		{
			testName:           "check expired app permission",
			checkUser:          user1,
			checkPubKey:        authAppPriv.PubKey(),
			atWhen:             baseTime.Add(time.Duration(101) * time.Second),
			permission:         types.AppPermission,
			expectUser:         "",
			expectResult:       ErrCheckAuthenticatePubKeyOwner(user1),
			expectGrantPubKeys: nil,
		},
		{
			testName:           "check expired preauth permission",
			checkUser:          user1,
			checkPubKey:        authTxPriv.PubKey(),
			atWhen:             baseTime.Add(time.Duration(101) * time.Second),
			amount:             types.NewCoinFromInt64(100),
			permission:         types.PreAuthorizationPermission,
			expectUser:         "",
			expectResult:       ErrCheckAuthenticatePubKeyOwner(user1),
			expectGrantPubKeys: nil,
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: tc.atWhen})
		keyOwner, err := am.CheckSigningPubKeyOwner(ctx, tc.checkUser, tc.checkPubKey, tc.permission, tc.amount)
		if tc.expectResult == nil {
			if tc.expectUser != keyOwner {
				t.Errorf("%s: diff key owner,  got %v, want %v", tc.testName, keyOwner, tc.expectUser)
				return
			}
		} else {
			fmt.Println(tc.testName, tc.expectResult.Result(), err)
			if !assert.Equal(t, tc.expectResult.Result(), err.Result()) {
				t.Errorf("%s: diff result,  got %v, want %v", tc.testName, err.Result(), tc.expectResult.Result())
			}
		}

		grantPubKeys, err := am.storage.GetGrantPermissions(ctx, tc.checkUser, tc.expectUser)
		if tc.expectGrantPubKeys == nil {
			if err == nil {
				t.Errorf("%s: got nil err", tc.testName)
			}
		} else {
			if err != nil {
				t.Errorf("%s: got non-empty err %v", tc.testName, err)
			}
			if len(tc.expectGrantPubKeys) != len(grantPubKeys) {
				t.Errorf("%s: expect grant pubkey length is different,  got %v, want %v", tc.testName, len(grantPubKeys), len(tc.expectGrantPubKeys))
			}
		}
	}
}

func TestRevokePermission(t *testing.T) {
	testName := "TestRevokePermission"

	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")
	userWithAppPermission := types.AccountKey("userWithAppPermission")
	userWithPreAuthPermission := types.AccountKey("userWithPreAuthPermission")

	createTestAccount(ctx, am, string(user1))
	createTestAccount(ctx, am, string(userWithAppPermission))
	createTestAccount(ctx, am, string(userWithPreAuthPermission))

	baseTime := ctx.BlockHeader().Time

	err := am.AuthorizePermission(ctx, user1, userWithAppPermission, 100, types.AppPermission, types.NewCoinFromInt64(0))
	if err != nil {
		t.Errorf("%s: failed to authorize user1 app permission to user with only app permission, got err %v", testName, err)
	}

	err = am.AuthorizePermission(ctx, user2, userWithAppPermission, 100, types.AppPermission, types.NewCoinFromInt64(0))
	if err != nil {
		t.Errorf("%s: failed to authorize user2 app permission to user with only app permission, got err %v", testName, err)
	}

	err = am.AuthorizePermission(ctx, user1, userWithPreAuthPermission, 100, types.PreAuthorizationPermission, types.NewCoinFromInt64(100))
	if err != nil {
		t.Errorf("%s: failed to authorize user1 preauth permission to user with preauth permission, got err %v", testName, err)
	}
	testCases := []struct {
		testName     string
		user         types.AccountKey
		revokeFrom   types.AccountKey
		permission   types.Permission
		atWhen       time.Time
		expectResult sdk.Error
	}{
		{
			testName:     "normal revoke app permission",
			user:         user1,
			revokeFrom:   userWithAppPermission,
			permission:   types.AppPermission,
			atWhen:       baseTime,
			expectResult: nil,
		},
		{
			testName:     "revoke non-exist permission, since it's revoked before",
			user:         user1,
			revokeFrom:   userWithAppPermission,
			permission:   types.AppPermission,
			atWhen:       baseTime,
			expectResult: model.ErrGrantPubKeyNotFound(),
		},
		{
			testName:     "normal revoke preauth permission",
			user:         user1,
			revokeFrom:   userWithPreAuthPermission,
			permission:   types.PreAuthorizationPermission,
			atWhen:       baseTime.Add(time.Duration(101) * time.Second),
			expectResult: nil,
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: tc.atWhen})
		err := am.RevokePermission(ctx, tc.user, tc.revokeFrom, tc.permission)
		if !assert.Equal(t, tc.expectResult, err) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, err, tc.expectResult)
		}
	}
}

func TestAuthorizePermission(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")
	user3 := types.AccountKey("user32")
	nonExistUser := types.AccountKey("nonExistUser")

	createTestAccount(ctx, am, string(user1))
	createTestAccount(ctx, am, string(user2))
	createTestAccount(ctx, am, string(user3))

	baseTime := ctx.BlockHeader().Time

	testCases := []struct {
		testName           string
		user               types.AccountKey
		grantTo            types.AccountKey
		level              types.Permission
		amount             types.Coin
		validityPeriod     int64
		expectResult       sdk.Error
		expectGrantPubKeys []*model.GrantPermission
	}{
		{
			testName:       "normal grant app permission",
			user:           user1,
			grantTo:        user2,
			level:          types.AppPermission,
			validityPeriod: 100,
			amount:         types.NewCoinFromInt64(0),
			expectResult:   nil,
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    user2,
					Permission: types.AppPermission,
					ExpiresAt:  baseTime.Unix() + 100,
					CreatedAt:  baseTime.Unix(),
					Amount:     types.NewCoinFromInt64(0),
				},
			},
		},
		{
			testName:       "override app permission",
			user:           user1,
			grantTo:        user2,
			level:          types.AppPermission,
			validityPeriod: 1000,
			amount:         types.NewCoinFromInt64(0),
			expectResult:   nil,
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    user2,
					Permission: types.AppPermission,
					ExpiresAt:  baseTime.Unix() + 1000,
					CreatedAt:  baseTime.Unix(),
					Amount:     types.NewCoinFromInt64(0),
				},
			},
		},
		{
			testName:       "grant app permission to non-exist user",
			user:           user1,
			grantTo:        nonExistUser,
			level:          types.AppPermission,
			validityPeriod: 1000,
			amount:         types.NewCoinFromInt64(0),
			expectResult:   ErrAccountNotFound(nonExistUser),
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    user2,
					Permission: types.AppPermission,
					ExpiresAt:  baseTime.Unix() + 1000,
					CreatedAt:  baseTime.Unix(),
					Amount:     types.NewCoinFromInt64(0),
				},
			},
		},
		{
			testName:       "grant pre authorization permission",
			user:           user1,
			grantTo:        user3,
			level:          types.PreAuthorizationPermission,
			validityPeriod: 100,
			amount:         types.NewCoinFromInt64(1000),
			expectResult:   nil,
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    user3,
					Permission: types.PreAuthorizationPermission,
					ExpiresAt:  baseTime.Unix() + 100,
					CreatedAt:  baseTime.Unix(),
					Amount:     types.NewCoinFromInt64(1000),
				},
			},
		},
		{
			testName:       "override pre authorization permission",
			user:           user1,
			grantTo:        user3,
			level:          types.PreAuthorizationPermission,
			validityPeriod: 1000,
			amount:         types.NewCoinFromInt64(10000),
			expectResult:   nil,
			expectGrantPubKeys: []*model.GrantPermission{
				&model.GrantPermission{
					GrantTo:    user3,
					Permission: types.PreAuthorizationPermission,
					ExpiresAt:  baseTime.Unix() + 1000,
					CreatedAt:  baseTime.Unix(),
					Amount:     types.NewCoinFromInt64(10000),
				},
			},
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: baseTime})
		err := am.AuthorizePermission(ctx, tc.user, tc.grantTo, tc.validityPeriod, tc.level, tc.amount)
		if !assert.Equal(t, tc.expectResult, err) {
			fmt.Println(err)
			t.Errorf("%s: failed to authorize permission, got err %v", tc.testName, err)
		}

		if tc.expectResult == nil {
			grantPubKeys, err := am.storage.GetGrantPermissions(ctx, tc.user, tc.grantTo)
			if err != nil {
				t.Errorf("%s: failed to get grant pub key, got err %v", tc.testName, err)
			}
			if !assert.Equal(t, tc.expectGrantPubKeys, grantPubKeys) {
				t.Errorf("%s: diff grant pub key, got %v, want %v", tc.testName, grantPubKeys, tc.expectGrantPubKeys)
			}
		}
	}
}

func TestAccountRecoverNormalCase(t *testing.T) {
	testName := "TestAccountRecoverNormalCase"

	ctx, am, _ := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	user1 := types.AccountKey("user1")

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		t.Errorf("%s: failed to get coin day param relationship, got err %v", testName, err)
	}

	createTestAccount(ctx, am, string(user1))

	newResetPrivKey := secp256k1.GenPrivKey()
	newTransactionPrivKey := secp256k1.GenPrivKey()
	newAppPrivKey := secp256k1.GenPrivKey()

	err = am.RecoverAccount(
		ctx, user1, newResetPrivKey.PubKey(), newTransactionPrivKey.PubKey(),
		newAppPrivKey.PubKey())
	if err != nil {
		t.Errorf("%s: failed to recover account, got err %v", testName, err)
	}

	accInfo := model.AccountInfo{
		Username:       user1,
		CreatedAt:      ctx.BlockHeader().Time.Unix(),
		ResetKey:       newResetPrivKey.PubKey(),
		TransactionKey: newTransactionPrivKey.PubKey(),
		AppKey:         newAppPrivKey.PubKey(),
	}
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee,
		CoinDay: accParam.RegisterFee,
	}

	checkAccountInfo(t, ctx, testName, user1, accInfo)
	checkBankKVByUsername(t, ctx, testName, user1, bank)

	pendingCoinDayQueue := model.PendingCoinDayQueue{
		TotalCoinDay: sdk.ZeroDec(),
		TotalCoin:    types.NewCoinFromInt64(0),
	}
	checkPendingCoinDay(t, ctx, testName, user1, pendingCoinDayQueue)

	coinDay, err := am.GetCoinDay(ctx, user1)
	if err != nil {
		t.Errorf("%s: failed to get coin day, got err %v", testName, err)
	}
	if !coinDay.IsEqual(accParam.RegisterFee) {
		t.Errorf("%s: diff coin day, got %v, want %v", testName, coinDay, accParam.RegisterFee)
	}

	ctx = ctx.WithBlockHeader(
		abci.Header{
			ChainID: "Lino", Height: 1,
			Time: ctx.BlockHeader().Time.Add(time.Duration(coinDayParams.SecondsToRecoverCoinDay) * time.Second)})
	coinDay, err = am.GetCoinDay(ctx, user1)
	if err != nil {
		t.Errorf("%s: failed to get coin day again, got err %v", testName, err)
	}
	if !coinDay.IsEqual(accParam.RegisterFee) {
		t.Errorf("%s: diff coin day again, got %v, want %v", testName, coinDay, accParam.RegisterFee)
	}
}

func TestIncreaseSequenceByOne(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	createTestAccount(ctx, am, string(user1))

	testCases := []struct {
		testName       string
		user           types.AccountKey
		increaseTimes  int
		expectSequence uint64
	}{
		{
			testName:       "increase seq once",
			user:           user1,
			increaseTimes:  1,
			expectSequence: 1,
		},
		{
			testName:       "increase seq 100 times",
			user:           user1,
			increaseTimes:  100,
			expectSequence: 101,
		},
	}

	for _, tc := range testCases {
		for i := 0; i < tc.increaseTimes; i++ {
			am.IncreaseSequenceByOne(ctx, user1)
		}
		seq, err := am.GetSequence(ctx, user1)
		if err != nil {
			t.Errorf("%s: failed to get sequence, got err %v", tc.testName, err)
		}
		if seq != tc.expectSequence {
			t.Errorf("%s: diff seq, got %v, want %v", tc.testName, seq, tc.expectSequence)
		}
	}
}

func TestLastReportOrUpvoteAt(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	createTestAccount(ctx, am, string(user1))

	testCases := []struct {
		testName             string
		lastReportOrUpvoteAt int64
	}{
		{
			testName:             "last report or upvote at current time",
			lastReportOrUpvoteAt: time.Now().Unix(),
		},
		{
			testName:             "last report or upvote at time 0",
			lastReportOrUpvoteAt: 0,
		},
	}

	for _, tc := range testCases {
		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(tc.lastReportOrUpvoteAt, 0)})
		err := am.UpdateLastReportOrUpvoteAt(newCtx, user1)
		if err != nil {
			t.Errorf("%s: failed to update last report or update at, got err %v", tc.testName, err)
		}
		lastReportOrUpdateAt, err := am.GetLastReportOrUpvoteAt(ctx, user1)
		if err != nil {
			t.Errorf("%s: failed to get last report or update at, got err %v", tc.testName, err)
		}
		assert.Equal(t, lastReportOrUpdateAt, tc.lastReportOrUpvoteAt)
	}
}

func TestLastPostAt(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	createTestAccount(ctx, am, string(user1))

	testCases := []struct {
		testName   string
		lastPostAt int64
	}{
		{
			testName:   "last post at current time",
			lastPostAt: time.Now().Unix(),
		},
		{
			testName:   "last post at time 0",
			lastPostAt: 0,
		},
	}

	for _, tc := range testCases {
		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(tc.lastPostAt, 0)})
		err := am.UpdateLastPostAt(newCtx, user1)
		if err != nil {
			t.Errorf("%s: failed to update last report or update at, got err %v", tc.testName, err)
		}
		lastPostAt, err := am.GetLastPostAt(ctx, user1)
		if err != nil {
			t.Errorf("%s: failed to get last report or update at, got err %v", tc.testName, err)
		}
		assert.Equal(t, lastPostAt, tc.lastPostAt)
	}
}
func TestAddFrozenMoney(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	createTestAccount(ctx, am, string(user1))

	testCases := []struct {
		testName                string
		frozenAmount            types.Coin
		startAt                 int64
		interval                int64
		times                   int64
		expectNumOfFrozenAmount int
	}{
		{
			testName:                "add the first 100 frozen money",
			frozenAmount:            types.NewCoinFromInt64(100),
			startAt:                 1000000,
			interval:                10,
			times:                   5,
			expectNumOfFrozenAmount: 1,
		},
		{
			testName:                "add the second 100 frozen money, clear the first one",
			frozenAmount:            types.NewCoinFromInt64(100),
			startAt:                 1200000,
			interval:                10,
			times:                   5,
			expectNumOfFrozenAmount: 1,
		},
		{
			testName:                "add the third 100 frozen money",
			frozenAmount:            types.NewCoinFromInt64(100),
			startAt:                 1300000,
			interval:                10,
			times:                   5,
			expectNumOfFrozenAmount: 2,
		},
		{
			testName:                "add the fourth 100 frozen money, clear the second one",
			frozenAmount:            types.NewCoinFromInt64(100),
			startAt:                 1400000,
			interval:                10,
			times:                   5,
			expectNumOfFrozenAmount: 2,
		},
		{
			testName:                "add the fifth 100 frozen money, clear the third and fourth ones",
			frozenAmount:            types.NewCoinFromInt64(100),
			startAt:                 1600000,
			interval:                10,
			times:                   5,
			expectNumOfFrozenAmount: 1,
		}, // this one is used to re-produce the out-of-bound bug.
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: time.Unix(tc.startAt, 0)})
		err := am.AddFrozenMoney(ctx, user1, tc.frozenAmount, tc.startAt, tc.interval, tc.times)
		if err != nil {
			t.Errorf("%s: failed to add frozen money, got err %v", tc.testName, err)
		}

		accountBank, err := am.storage.GetBankFromAccountKey(ctx, user1)
		if err != nil {
			t.Errorf("%s: failed to get bank, got err %v", tc.testName, err)
		}
		if len(accountBank.FrozenMoneyList) != tc.expectNumOfFrozenAmount {
			t.Errorf("%s: diff num of frozen money, got %v, want %v", tc.testName, len(accountBank.FrozenMoneyList), tc.expectNumOfFrozenAmount)
		}
	}
}
