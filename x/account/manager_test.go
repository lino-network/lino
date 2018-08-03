package account

import (
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

func checkBankKVByUsername(
	t *testing.T, ctx sdk.Context, testName string, username types.AccountKey, bank model.AccountBank) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	bankPtr, err := accStorage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		t.Errorf("%s, failed to get bank, got err %v", testName, err)
	}
	if !assert.Equal(t, bank, *bankPtr) {
		t.Errorf("%s: diff bank, got %v, want %v", testName, *bankPtr, bank)
	}
}

func checkBalanceHistory(
	t *testing.T, ctx sdk.Context, testName string, username types.AccountKey,
	timeSlot int64, balanceHistory model.BalanceHistory) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	balanceHistoryPtr, err := accStorage.GetBalanceHistory(ctx, username, timeSlot)
	if err != nil {
		t.Errorf("%s, failed to get balance history, got err %v", testName, err)
	}
	if !assert.Equal(t, balanceHistory, *balanceHistoryPtr) {
		t.Errorf("%s: diff balance history, got %v, want %v", testName, *balanceHistoryPtr, balanceHistory)
	}
}

func checkPendingStake(
	t *testing.T, ctx sdk.Context, testName string, username types.AccountKey, pendingStakeQueue model.PendingStakeQueue) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	pendingStakeQueuePtr, err := accStorage.GetPendingStakeQueue(ctx, username)
	if err != nil {
		t.Errorf("%s, failed to get pending stake queue, got err %v", testName, err)
	}
	if !assert.Equal(t, pendingStakeQueue, *pendingStakeQueuePtr) {
		t.Errorf("%s: diff pending stake queue, got %v, want %v", testName, *pendingStakeQueuePtr, pendingStakeQueue)
	}
}

func checkAccountInfo(
	t *testing.T, ctx sdk.Context, testName string, accKey types.AccountKey, accInfo model.AccountInfo) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	info, err := accStorage.GetInfo(ctx, accKey)
	if err != nil {
		t.Errorf("%s, failed to get account info, got err %v", testName, err)
	}
	if !assert.Equal(t, accInfo, *info) {
		t.Errorf("%s: diff account info, got %v, want %v", testName, *info, accInfo)
	}
}

func checkAccountMeta(
	t *testing.T, ctx sdk.Context, testName string, accKey types.AccountKey, accMeta model.AccountMeta) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	metaPtr, err := accStorage.GetMeta(ctx, accKey)
	if err != nil {
		t.Errorf("%s, failed to get account meta, got err %v", testName, err)
	}
	if !assert.Equal(t, accMeta, *metaPtr) {
		t.Errorf("%s: diff account meta, got %v, want %v", testName, *metaPtr, accMeta)
	}
}

func checkAccountReward(
	t *testing.T, ctx sdk.Context, testName string, accKey types.AccountKey, reward model.Reward) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	rewardPtr, err := accStorage.GetReward(ctx, accKey)
	if err != nil {
		t.Errorf("%s, failed to get reward, got err %v", testName, err)
	}
	if !assert.Equal(t, reward, *rewardPtr) {
		t.Errorf("%s: diff reward, got %v, want %v", testName, *rewardPtr, reward)
	}
}

func checkRewardHistory(
	t *testing.T, ctx sdk.Context, testName string, accKey types.AccountKey, bucketSlot int64, wantNumOfReward int) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	rewardHistoryPtr, err := accStorage.GetRewardHistory(ctx, accKey, bucketSlot)
	if err != nil {
		t.Errorf("%s, failed to get reward history, got err %v", testName, err)
	}
	if wantNumOfReward != len(rewardHistoryPtr.Details) {
		t.Errorf("%s: diff account rewards, got %v, want %v", testName, len(rewardHistoryPtr.Details), wantNumOfReward)
	}
}

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

func TestAddCoin(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		t.Error("TestAddCoin: failed to get coin day param")
	}

	fromUser1, fromUser2, testUser :=
		types.AccountKey("fromUser1"), types.AccountKey("fromuser2"), types.AccountKey("testUser")

	baseTime := time.Now().Unix()
	baseTime1 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake/2
	baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake + 1
	baseTime3 := baseTime2 + coinDayParams.SecondsToRecoverCoinDayStake + 1

	ctx = ctx.WithBlockHeader(abci.Header{Time: baseTime})
	createTestAccount(ctx, am, string(testUser))

	testCases := []struct {
		testName                 string
		amount                   types.Coin
		from                     types.AccountKey
		detailType               types.TransferDetailType
		memo                     string
		atWhen                   int64
		expectBank               model.AccountBank
		expectPendingStakeQueue  model.PendingStakeQueue
		expectBalanceHistorySlot model.BalanceHistory
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
				Stake:   accParam.RegisterFee,
				NumOfTx: 2,
			},
			expectPendingStakeQueue: model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat(),
				TotalCoin:        c100,
				PendingStakeList: []model.PendingStake{
					{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
				},
			},
			expectBalanceHistorySlot: model.BalanceHistory{
				Details: []model.Detail{
					{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					{
						Amount:     c100,
						From:       fromUser1,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
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
				Stake:   accParam.RegisterFee,
				NumOfTx: 3,
			},
			expectPendingStakeQueue: model.PendingStakeQueue{
				LastUpdatedAt:    baseTime1,
				StakeCoinInQueue: sdk.NewRat(5000000, 1),
				TotalCoin:        c100.Plus(c100),
				PendingStakeList: []model.PendingStake{
					{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
					{
						StartTime: baseTime1,
						EndTime:   baseTime1 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
				},
			},
			expectBalanceHistorySlot: model.BalanceHistory{
				Details: []model.Detail{
					{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					{
						Amount:     c100,
						From:       fromUser1,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
					{
						Amount:     c100,
						From:       fromUser2,
						To:         testUser,
						CreatedAt:  baseTime1,
						DetailType: types.DonationIn,
						Memo:       "permlink",
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
				Stake:   accParam.RegisterFee.Plus(c100),
				NumOfTx: 4,
			},
			expectPendingStakeQueue: model.PendingStakeQueue{
				LastUpdatedAt:    baseTime2,
				StakeCoinInQueue: sdk.NewRat(945003125, 189),
				TotalCoin:        c100.Plus(c100),
				PendingStakeList: []model.PendingStake{
					{
						StartTime: baseTime1,
						EndTime:   baseTime1 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
					{
						StartTime: baseTime2,
						EndTime:   baseTime2 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
				},
			},
			expectBalanceHistorySlot: model.BalanceHistory{
				Details: []model.Detail{
					{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					{
						Amount:     c100,
						From:       fromUser1,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
					{
						Amount:     c100,
						From:       fromUser2,
						To:         testUser,
						CreatedAt:  baseTime1,
						DetailType: types.DonationIn,
						Memo:       "permlink",
					},
					{
						Amount:     c100,
						From:       "",
						To:         testUser,
						CreatedAt:  baseTime2,
						DetailType: types.ClaimReward,
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
				Stake:   accParam.RegisterFee.Plus(c100),
				NumOfTx: 4,
			},
			expectPendingStakeQueue: model.PendingStakeQueue{
				LastUpdatedAt:    baseTime2,
				StakeCoinInQueue: sdk.NewRat(945003125, 189),
				TotalCoin:        c100.Plus(c100),
				PendingStakeList: []model.PendingStake{
					{
						StartTime: baseTime1,
						EndTime:   baseTime1 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
					{
						StartTime: baseTime2,
						EndTime:   baseTime2 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
				},
			},
			expectBalanceHistorySlot: model.BalanceHistory{
				Details: []model.Detail{
					{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					{
						Amount:     c100,
						From:       fromUser1,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
					{
						Amount:     c100,
						From:       fromUser2,
						To:         testUser,
						CreatedAt:  baseTime1,
						DetailType: types.DonationIn,
						Memo:       "permlink",
					},
					{
						Amount:     c100,
						From:       "",
						To:         testUser,
						CreatedAt:  baseTime2,
						DetailType: types.ClaimReward,
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
		checkPendingStake(t, ctx, tc.testName, types.AccountKey(testUser), tc.expectPendingStakeQueue)
		checkBalanceHistory(
			t, ctx, tc.testName, types.AccountKey(testUser), 0, tc.expectBalanceHistorySlot)
	}
}

func TestMinusCoin(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		t.Error("TestMinusCoin: failed to get coin day param")
	}

	userWithSufficientSaving := types.AccountKey("user1")
	userWithLimitSaving := types.AccountKey("user3")
	fromUser, toUser := types.AccountKey("fromUser"), types.AccountKey("toUser")

	// Get the minimum time of this history slot
	baseTime := time.Now().Unix()
	// baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake + 1
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
		testName                string
		fromUser                types.AccountKey
		userPriv                crypto.PrivKey
		expectErr               sdk.Error
		amount                  types.Coin
		atWhen                  int64
		to                      types.AccountKey
		memo                    string
		detailType              types.TransferDetailType
		expectBank              model.AccountBank
		expectPendingStakeQueue model.PendingStakeQueue
		expectBalanceHistory    model.BalanceHistory
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
				NumOfTx: 3,
				Stake:   accParam.RegisterFee,
			},
			expectPendingStakeQueue: model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat(),
				TotalCoin:        accParam.RegisterFee.Minus(coin1),
				PendingStakeList: []model.PendingStake{
					{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      accParam.RegisterFee.Minus(coin1),
					}},
			},
			expectBalanceHistory: model.BalanceHistory{
				Details: []model.Detail{
					{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         userWithSufficientSaving,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					{
						Amount:     accParam.RegisterFee,
						From:       fromUser,
						To:         userWithSufficientSaving,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					{
						Amount:     coin1,
						From:       userWithSufficientSaving,
						To:         toUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferOut,
						Memo:       "memo",
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
				NumOfTx: 1,
				Stake:   accParam.RegisterFee,
			},
			expectPendingStakeQueue: model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat(),
				TotalCoin:        accParam.RegisterFee,
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      accParam.RegisterFee,
					}},
			},
			expectBalanceHistory: model.BalanceHistory{
				Details: []model.Detail{
					{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         toUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
				},
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
				NumOfTx: 1,
				Stake:   accParam.RegisterFee,
			},
			expectPendingStakeQueue: model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat(),
				TotalCoin:        accParam.RegisterFee,
				PendingStakeList: []model.PendingStake{
					{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      accParam.RegisterFee,
					}},
			},
			expectBalanceHistory: model.BalanceHistory{
				Details: []model.Detail{
					{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         userWithLimitSaving,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountRegisterDepositMemo,
					},
				},
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
			checkPendingStake(t, ctx, tc.testName, tc.fromUser, tc.expectPendingStakeQueue)
			checkBalanceHistory(t, ctx, tc.testName, tc.fromUser, 0, tc.expectBalanceHistory)
		}
	}
}

func TestBalanceHistory(t *testing.T) {
	fromUser, toUser := types.AccountKey("fromUser"), types.AccountKey("toUser")

	testCases := []struct {
		testName        string
		numOfAdding     int
		numOfMinus      int
		expectTotalSlot int64
	}{
		{
			testName:        "test only one adding",
			numOfAdding:     1,
			numOfMinus:      0,
			expectTotalSlot: 1,
		},
		{
			testName:        "test 99 adding, which fullfills 1 bundles",
			numOfAdding:     99,
			numOfMinus:      0,
			expectTotalSlot: 1,
		},
		{
			testName:        "test adding and minus, which results in 2 bundles",
			numOfAdding:     50,
			numOfMinus:      50,
			expectTotalSlot: 2,
		},
	}
	for _, tc := range testCases {
		ctx, am, accParam := setupTest(t, 1)

		user1 := types.AccountKey("user1")
		createTestAccount(ctx, am, string(user1))

		for i := 0; i < tc.numOfAdding; i++ {
			err := am.AddSavingCoin(ctx, user1, coin1, fromUser, "", types.TransferIn)
			if err != nil {
				t.Errorf("%s: failed to add saving coin, got err %v", tc.testName, err)
			}
		}
		for i := 0; i < tc.numOfMinus; i++ {
			err := am.MinusSavingCoin(ctx, user1, coin1, toUser, "", types.TransferOut)
			if err != nil {
				t.Errorf("%s: failed to minus saving coin, got err %v", tc.testName, err)
			}
		}

		bank, err := am.storage.GetBankFromAccountKey(ctx, user1)
		if err != nil {
			t.Errorf("%s: failed to get bank, got err %v", tc.testName, err)
		}

		// add one init transfer in
		expectNumOfTx := int64(tc.numOfAdding + tc.numOfMinus + 1)
		if expectNumOfTx != bank.NumOfTx {
			t.Errorf("%s: diff num of tx, got %v, want %v", tc.testName, bank.NumOfTx, expectNumOfTx)
		}

		// total slot should use previous states to get expected slots
		actualTotalSlot := (expectNumOfTx-1)/accParam.BalanceHistoryBundleSize + 1
		if tc.expectTotalSlot != actualTotalSlot {
			t.Errorf("%s: diff total slot, got %v, want %v", tc.testName, actualTotalSlot, tc.expectTotalSlot)
		}

		actualNumOfAdding, actualNumOfMinus := 0, 0
		for slot := int64(0); slot < actualTotalSlot; slot++ {
			balanceHistory, err := am.storage.GetBalanceHistory(ctx, user1, slot)
			if err != nil {
				t.Errorf("%s: failed to get balance history, got err %v", tc.testName, err)
			}

			for _, tx := range balanceHistory.Details {
				if tx.DetailType == types.TransferIn {
					actualNumOfAdding++
				}
				if tx.DetailType == types.TransferOut {
					actualNumOfMinus++
				}
			}
		}
		// include create account init transaction
		if tc.numOfAdding+1 != actualNumOfAdding {
			t.Errorf("%s: diff num of adding, got %v, want %v", tc.testName, actualNumOfAdding, tc.numOfAdding+1)
		}
		if tc.numOfMinus != actualNumOfMinus {
			t.Errorf("%s: diff num of minus, got %v, want %v", tc.testName, actualNumOfMinus, tc.numOfMinus)
		}
	}
}

func TestAddBalanceHistory(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)

	testCases := []struct {
		testName              string
		numOfTx               int64
		detail                model.Detail
		expectNumOfTxInBundle int
	}{
		{
			testName: "try first transaction in first slot",
			numOfTx:  0,
			detail: model.Detail{
				From:       "test1",
				To:         "test2",
				Amount:     types.NewCoinFromInt64(1),
				DetailType: types.TransferIn,
				CreatedAt:  time.Now().Unix(),
			},
			expectNumOfTxInBundle: 1,
		},
		{
			testName: "try second transaction in first slot",
			numOfTx:  1,
			detail: model.Detail{
				From:       "test2",
				To:         "test1",
				Amount:     types.NewCoinFromInt64(1 * types.Decimals),
				DetailType: types.TransferOut,
				CreatedAt:  time.Now().Unix(),
			},
			expectNumOfTxInBundle: 2,
		},
		{
			testName: "add transaction to the end of the first slot limitation",
			numOfTx:  99,
			detail: model.Detail{
				From:       "test1",
				To:         "post",
				Amount:     types.NewCoinFromInt64(1 * types.Decimals),
				DetailType: types.DonationOut,
				CreatedAt:  time.Now().Unix(),
				Memo:       "",
			},
			expectNumOfTxInBundle: 3,
		},
		{
			testName: "add transaction to next slot",
			numOfTx:  100,
			detail: model.Detail{
				From:       "",
				To:         "test1",
				Amount:     types.NewCoinFromInt64(1 * types.Decimals),
				DetailType: types.DeveloperDeposit,
				CreatedAt:  time.Now().Unix(),
			},
			expectNumOfTxInBundle: 1,
		},
	}

	for _, tc := range testCases {
		err := am.AddBalanceHistory(ctx, user1, tc.numOfTx, tc.detail)
		if err != nil {
			t.Errorf("%s: failed to add balance history, got err %v", tc.testName, err)
		}

		balanceHistory, err :=
			am.storage.GetBalanceHistory(
				ctx, user1, tc.numOfTx/accParam.BalanceHistoryBundleSize)
		if err != nil {
			t.Errorf("%s: failed to get balance history, got err %v", tc.testName, err)
		}

		if tc.expectNumOfTxInBundle != len(balanceHistory.Details) {
			t.Errorf("%s: diff num of tx in bunlde, got %v, want %v", tc.testName, len(balanceHistory.Details), tc.expectNumOfTxInBundle)
		}
		if !assert.Equal(t, tc.detail, balanceHistory.Details[tc.expectNumOfTxInBundle-1]) {
			t.Errorf("%s: diff detail, got %v, want %v", tc.testName, balanceHistory.Details[tc.expectNumOfTxInBundle-1], tc.detail)
		}
	}
}

func TestCreateAccountNormalCase(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)

	resetPriv := secp256k1.GenPrivKey()
	txPriv := secp256k1.GenPrivKey()
	appPriv := secp256k1.GenPrivKey()

	accKey := types.AccountKey("accKey")

	// normal test
	assert.False(t, am.DoesAccountExist(ctx, accKey))
	err := am.CreateAccount(
		ctx, accountReferrer, accKey, resetPriv.PubKey(), txPriv.PubKey(),
		appPriv.PubKey(), accParam.RegisterFee)
	if err != nil {
		t.Errorf("TestCreateAccountNormalCase: failed to create account, got err %v", err)
	}

	assert.True(t, am.DoesAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee,
		NumOfTx: 1,
		Stake:   accParam.RegisterFee,
	}
	checkBankKVByUsername(t, ctx, "TestCreateAccountNormalCase", accKey, bank)

	pendingStakeQueue := model.PendingStakeQueue{StakeCoinInQueue: sdk.ZeroRat(), TotalCoin: types.NewCoinFromInt64(0)}
	checkPendingStake(t, ctx, "TestCreateAccountNormalCase", accKey, pendingStakeQueue)

	accInfo := model.AccountInfo{
		Username:       accKey,
		CreatedAt:      ctx.BlockHeader().Time,
		ResetKey:       resetPriv.PubKey(),
		TransactionKey: txPriv.PubKey(),
		AppKey:         appPriv.PubKey(),
	}
	checkAccountInfo(t, ctx, "TestCreateAccountNormalCase", accKey, accInfo)
	accMeta := model.AccountMeta{
		LastActivityAt:       ctx.BlockHeader().Time,
		LastReportOrUpvoteAt: ctx.BlockHeader().Time,
		TransactionCapacity:  accParam.RegisterFee,
	}
	checkAccountMeta(t, ctx, "TestCreateAccountNormalCase", accKey, accMeta)

	reward := model.Reward{coin0, coin0, coin0, coin0, coin0}
	checkAccountReward(t, ctx, "TestCreateAccountNormalCase", accKey, reward)

	balanceHistory, err := am.storage.GetBalanceHistory(ctx, accKey, 0)
	if err != nil {
		t.Errorf("TestCreateAccountNormalCase: failed to get balance history, got err %v", err)
	}
	if len(balanceHistory.Details) != 1 {
		t.Errorf("TestCreateAccountNormalCase: diff num of balance, got %v, want %v", len(balanceHistory.Details), 1)
	}

	wantDetail := model.Detail{
		From:       accountReferrer,
		To:         accKey,
		Amount:     accParam.RegisterFee,
		CreatedAt:  ctx.BlockHeader().Time,
		DetailType: types.TransferIn,
		Memo:       types.InitAccountWithFullStakeMemo,
	}
	if !assert.Equal(t, wantDetail, balanceHistory.Details[0]) {
		t.Errorf("TestCreateAccountNormalCase: diff detail, got %v, want %v", balanceHistory.Details[0], wantDetail)
	}
}

func TestCreateAccountWithLargeRegisterFee(t *testing.T) {
	testName := "TestCreateAccountWithLargeRegisterFee"

	ctx, am, accParam := setupTest(t, 1)

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
		NumOfTx: 2,
		Stake:   accParam.RegisterFee,
	}
	checkBankKVByUsername(t, ctx, testName, accKey, bank)

	pendingStakeQueue := model.PendingStakeQueue{
		LastUpdatedAt:    ctx.BlockHeader().Time,
		StakeCoinInQueue: sdk.ZeroRat(),
		TotalCoin:        extraRegisterFee,
		PendingStakeList: []model.PendingStake{
			model.PendingStake{
				StartTime: ctx.BlockHeader().Time,
				EndTime:   ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake,
				Coin:      extraRegisterFee,
			},
		},
	}
	checkPendingStake(t, ctx, testName, accKey, pendingStakeQueue)

	accInfo := model.AccountInfo{
		Username:       accKey,
		CreatedAt:      ctx.BlockHeader().Time,
		ResetKey:       resetPriv.PubKey(),
		TransactionKey: txPriv.PubKey(),
		AppKey:         appPriv.PubKey(),
	}
	checkAccountInfo(t, ctx, testName, accKey, accInfo)

	accMeta := model.AccountMeta{
		LastActivityAt:       ctx.BlockHeader().Time,
		LastReportOrUpvoteAt: ctx.BlockHeader().Time,
		TransactionCapacity:  accParam.RegisterFee,
	}
	checkAccountMeta(t, ctx, testName, accKey, accMeta)

	reward := model.Reward{coin0, coin0, coin0, coin0, coin0}
	checkAccountReward(t, ctx, testName, accKey, reward)

	balanceHistory, err := am.storage.GetBalanceHistory(ctx, accKey, 0)
	if err != nil {
		t.Errorf("%s: failed to get balance history, got err %v", testName, err)
	}
	if len(balanceHistory.Details) != 2 {
		t.Errorf("%s: diff num of balance history, got %v, want %v", testName, len(balanceHistory.Details), 2)
	}

	wantDetails := []model.Detail{
		{
			From:       accountReferrer,
			To:         accKey,
			Amount:     accParam.RegisterFee,
			CreatedAt:  ctx.BlockHeader().Time,
			DetailType: types.TransferIn,
			Memo:       types.InitAccountWithFullStakeMemo,
		},
		{
			From:       accountReferrer,
			To:         accKey,
			Amount:     extraRegisterFee,
			CreatedAt:  ctx.BlockHeader().Time,
			DetailType: types.TransferIn,
			Memo:       types.InitAccountRegisterDepositMemo,
		},
	}
	if !assert.Equal(t, wantDetails, balanceHistory.Details) {
		t.Errorf("%s: diff details, got %v, want %v", testName, balanceHistory.Details, wantDetails)
	}
}

func TestInvalidCreateAccount(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	priv1 := secp256k1.GenPrivKey()
	priv2 := secp256k1.GenPrivKey()

	accKey1 := types.AccountKey("accKey1")
	accKey2 := types.AccountKey("accKey2")
	accKey3 := types.AccountKey("accKey3")

	testCases := []struct {
		testName    string
		username    types.AccountKey
		privKey     crypto.PrivKey
		registerFee types.Coin
		expectErr   sdk.Error
	}{
		{
			testName:    "register user with sufficient saving coin",
			username:    accKey1,
			privKey:     priv1,
			registerFee: accParam.RegisterFee,
			expectErr:   nil,
		},
		{
			testName:    "username already took",
			username:    accKey1,
			privKey:     priv1,
			registerFee: accParam.RegisterFee,
			expectErr:   ErrAccountAlreadyExists(accKey1),
		},
		{
			testName:    "username already took with different private key",
			username:    accKey1,
			privKey:     priv2,
			registerFee: accParam.RegisterFee,
			expectErr:   ErrAccountAlreadyExists(accKey1),
		},
		{
			testName:    "register the same private key",
			username:    accKey2,
			privKey:     priv1,
			registerFee: accParam.RegisterFee,
			expectErr:   nil,
		},
		{
			testName:    "insufficient register fee",
			username:    accKey3,
			privKey:     priv1,
			registerFee: types.NewCoinFromInt64(1),
			expectErr:   ErrRegisterFeeInsufficient(),
		},
	}
	for _, tc := range testCases {
		err := am.CreateAccount(
			ctx, accountReferrer, tc.username, tc.privKey.PubKey(),
			secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), tc.registerFee)
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
	ctx, am, accParam := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		t.Errorf("TestCoinDayByAccountKey: failed to get coin day param, got err %v", err)
	}

	totalCoinDaysSec := coinDayParams.SecondsToRecoverCoinDayStake
	registerFee := accParam.RegisterFee.ToInt64()
	doubleRegisterFee := types.NewCoinFromInt64(registerFee * 2)
	halfRegisterFee := types.NewCoinFromInt64(registerFee / 2)

	baseTime := ctx.BlockHeader().Time
	baseTime2 := baseTime + totalCoinDaysSec + (totalCoinDaysSec/registerFee)/2 + 1

	createTestAccount(ctx, am, string(accKey))

	testCases := []struct {
		testName            string
		isAdd               bool
		coin                types.Coin
		atWhen              int64
		expectSavingBalance types.Coin
		expectStake         types.Coin
		expectStakeInBank   types.Coin
		expectNumOfTx       int64
	}{
		{
			testName:            "add coin before charging first coin",
			isAdd:               true,
			coin:                accParam.RegisterFee,
			atWhen:              baseTime + (totalCoinDaysSec/registerFee)/2,
			expectSavingBalance: doubleRegisterFee,
			expectStake:         accParam.RegisterFee,
			expectStakeInBank:   accParam.RegisterFee,
			expectNumOfTx:       2,
		},
		{
			testName:            "check first coin",
			isAdd:               true,
			coin:                coin0,
			atWhen:              baseTime + (totalCoinDaysSec/registerFee)/2 + 1,
			expectSavingBalance: doubleRegisterFee,
			expectStake:         accParam.RegisterFee,
			expectStakeInBank:   accParam.RegisterFee,
			expectNumOfTx:       2,
		},
		{
			testName:            "check both transactions fully charged",
			isAdd:               true,
			coin:                coin0,
			atWhen:              baseTime2,
			expectSavingBalance: doubleRegisterFee,
			expectStake:         doubleRegisterFee,
			expectStakeInBank:   doubleRegisterFee,
			expectNumOfTx:       2,
		},
		{
			testName:            "withdraw half deposit",
			isAdd:               false,
			coin:                accParam.RegisterFee,
			atWhen:              baseTime2,
			expectSavingBalance: accParam.RegisterFee,
			expectStake:         accParam.RegisterFee,
			expectStakeInBank:   accParam.RegisterFee,
			expectNumOfTx:       3,
		},
		{
			testName:            "charge again",
			isAdd:               true,
			coin:                accParam.RegisterFee,
			atWhen:              baseTime2,
			expectSavingBalance: doubleRegisterFee,
			expectStake:         accParam.RegisterFee,
			expectStakeInBank:   accParam.RegisterFee,
			expectNumOfTx:       4,
		},
		{
			testName:            "withdraw half deposit while the last transaction is still charging",
			isAdd:               false,
			coin:                halfRegisterFee,
			atWhen:              baseTime2 + totalCoinDaysSec/2 + 1,
			expectSavingBalance: accParam.RegisterFee.Plus(halfRegisterFee),
			expectStake:         accParam.RegisterFee.Plus(types.NewCoinFromInt64(registerFee / 4)),
			expectStakeInBank:   accParam.RegisterFee,
			expectNumOfTx:       5,
		},
		{
			testName:            "withdraw last transaction which is still charging",
			isAdd:               false,
			coin:                halfRegisterFee,
			atWhen:              baseTime2 + totalCoinDaysSec/2 + 1,
			expectSavingBalance: accParam.RegisterFee,
			expectStake:         accParam.RegisterFee,
			expectStakeInBank:   accParam.RegisterFee,
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
		coin, err := am.GetStake(ctx, accKey)
		if err != nil {
			t.Errorf("%s: failed to get stake, got err %v", tc.testName, err)
		}

		if !tc.expectStake.IsEqual(coin) {
			t.Errorf("%s: diff stake, got %v, want %v", tc.testName, coin, tc.expectStake)
			return
		}

		bank := model.AccountBank{
			Saving:  tc.expectSavingBalance,
			Stake:   tc.expectStakeInBank,
			NumOfTx: tc.expectNumOfTx,
		}
		checkBankKVByUsername(t, ctx, tc.testName, accKey, bank)
	}
}

func TestAddIncomeAndReward(t *testing.T) {
	testName := "TestAddIncomeAndReward"

	ctx, am, accParam := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	createTestAccount(ctx, am, string(accKey))

	err := am.AddIncomeAndReward(ctx, accKey, c500, c200, c300, "donor1", "postAutho1", "post1")
	if err != nil {
		t.Errorf("%s: failed to add income and reward, got err %v", testName, err)
	}

	reward := model.Reward{c300, c200, c200, c300, c300}
	checkAccountReward(t, ctx, testName, accKey, reward)
	checkRewardHistory(t, ctx, testName, accKey, 0, 1)

	err = am.AddIncomeAndReward(ctx, accKey, c500, c300, c200, "donor2", "postAuthor1", "post1")
	if err != nil {
		t.Errorf("%s: failed to add income and reward again, got err %v", testName, err)
	}

	reward = model.Reward{c500, c500, c500, c500, c500}
	checkAccountReward(t, ctx, testName, accKey, reward)
	checkRewardHistory(t, ctx, testName, accKey, 0, 2)

	err = am.AddDirectDeposit(ctx, accKey, c500)
	if err != nil {
		t.Errorf("%s: failed to add direct deposit, got err %v", testName, err)
	}

	reward = model.Reward{c1000, c1000, c500, c500, c500}
	checkAccountReward(t, ctx, testName, accKey, reward)
	checkRewardHistory(t, ctx, testName, accKey, 0, 2)
	bank := model.AccountBank{
		Saving:      accParam.RegisterFee,
		NumOfTx:     1,
		NumOfReward: 2,
		Stake:       accParam.RegisterFee,
	}
	checkBankKVByUsername(t, ctx, testName, accKey, bank)

	err = am.ClaimReward(ctx, accKey)
	if err != nil {
		t.Errorf("%s: failed to add claim reward, got err %v", testName, err)
	}

	bank.Saving = accParam.RegisterFee.Plus(c500)
	bank.NumOfTx = 2
	bank.NumOfReward = 0
	checkBankKVByUsername(t, ctx, testName, accKey, bank)

	reward = model.Reward{c1000, c1000, c500, c500, c0}
	checkAccountReward(t, ctx, testName, accKey, reward)
}

func TestCheckUserTPSCapacity(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	bandwidthParams, err := am.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		t.Errorf("TestCheckUserTPSCapacity: failed to get bandwidth param, got err %v", err)
	}
	secondsToRecoverBandwidth := bandwidthParams.SecondsToRecoverBandwidth

	baseTime := ctx.BlockHeader().Time

	createTestAccount(ctx, am, string(accKey))
	err = am.AddSavingCoin(ctx, accKey, c100, "", "", types.TransferIn)
	if err != nil {
		t.Errorf("TestCheckUserTPSCapacity: failed to add saving coin, got err %v", err)
	}

	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	err = accStorage.SetPendingStakeQueue(
		ctx, accKey, &model.PendingStakeQueue{})
	if err != nil {
		t.Errorf("TestCheckUserTPSCapacity: failed to set pending stake queue, got err %v", err)
	}

	testCases := []struct {
		testName             string
		tpsCapacityRatio     sdk.Rat
		userStake            types.Coin
		lastActivity         int64
		lastCapacity         types.Coin
		currentTime          int64
		expectResult         sdk.Error
		expectRemainCapacity types.Coin
	}{
		{
			testName:             "tps capacity not enough",
			tpsCapacityRatio:     sdk.NewRat(1, 10),
			userStake:            types.NewCoinFromInt64(10 * types.Decimals),
			lastActivity:         baseTime,
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime,
			expectResult:         ErrAccountTPSCapacityNotEnough(accKey),
			expectRemainCapacity: types.NewCoinFromInt64(0)},
		{
			testName:             " 1/10 capacity ratio",
			tpsCapacityRatio:     sdk.NewRat(1, 10),
			userStake:            types.NewCoinFromInt64(10 * types.Decimals),
			lastActivity:         baseTime,
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime + secondsToRecoverBandwidth,
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(990000),
		},
		{
			testName:             " 1/2 capacity ratio",
			tpsCapacityRatio:     sdk.NewRat(1, 2),
			userStake:            types.NewCoinFromInt64(10 * types.Decimals),
			lastActivity:         baseTime,
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime + secondsToRecoverBandwidth,
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(950000),
		},
		{
			testName:             " 1/1 capacity ratio",
			tpsCapacityRatio:     sdk.NewRat(1, 1),
			userStake:            types.NewCoinFromInt64(10 * types.Decimals),
			lastActivity:         baseTime,
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime + secondsToRecoverBandwidth,
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(9 * types.Decimals),
		},
		{
			testName:             " 1/1 capacity ratio with 0 remaining",
			tpsCapacityRatio:     sdk.NewRat(1, 1),
			userStake:            types.NewCoinFromInt64(1 * types.Decimals),
			lastActivity:         baseTime,
			lastCapacity:         types.NewCoinFromInt64(10 * types.Decimals),
			currentTime:          baseTime,
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(0),
		},
		{
			testName:             " 1/1 capacity ratio with 1 remaining",
			tpsCapacityRatio:     sdk.NewRat(1, 1),
			userStake:            types.NewCoinFromInt64(10),
			lastActivity:         baseTime,
			lastCapacity:         types.NewCoinFromInt64(1 * types.Decimals),
			currentTime:          baseTime,
			expectResult:         ErrAccountTPSCapacityNotEnough(accKey),
			expectRemainCapacity: types.NewCoinFromInt64(1 * types.Decimals),
		},
		{
			testName:             " 1/1 capacity ratio with 1 stake and 0 remaining",
			tpsCapacityRatio:     sdk.NewRat(1, 1),
			userStake:            types.NewCoinFromInt64(1 * types.Decimals),
			lastActivity:         baseTime,
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime + secondsToRecoverBandwidth/2,
			expectResult:         ErrAccountTPSCapacityNotEnough(accKey),
			expectRemainCapacity: types.NewCoinFromInt64(0),
		},
		{
			testName:             " 1/2 capacity ratio with 0 remaining",
			tpsCapacityRatio:     sdk.NewRat(1, 2),
			userStake:            types.NewCoinFromInt64(1 * types.Decimals),
			lastActivity:         baseTime,
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime + secondsToRecoverBandwidth/2,
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(0),
		},
		{
			testName:             " 1/1 capacity ratio with 0 remaining and base time",
			tpsCapacityRatio:     sdk.NewRat(1, 1),
			userStake:            types.NewCoinFromInt64(1 * types.Decimals),
			lastActivity:         0,
			lastCapacity:         types.NewCoinFromInt64(0),
			currentTime:          baseTime,
			expectResult:         nil,
			expectRemainCapacity: types.NewCoinFromInt64(0),
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: tc.currentTime})
		bank := &model.AccountBank{
			Saving: tc.userStake,
			Stake:  tc.userStake,
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
			LastActivityAt:      ctx.BlockHeader().Time,
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

	ctx, am, accParam := setupTest(t, 1)
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
		testName          string
		checkUser         types.AccountKey
		checkPubKey       crypto.PubKey
		atWhen            int64
		amount            types.Coin
		permission        types.Permission
		expectUser        types.AccountKey
		expectResult      sdk.Error
		expectGrantPubKey *model.GrantPubKey
	}{
		{
			testName:          "check user's reset key",
			checkUser:         user1,
			checkPubKey:       resetKey.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(0),
			permission:        types.ResetPermission,
			expectUser:        user1,
			expectResult:      nil,
			expectGrantPubKey: nil,
		},
		{
			testName:          "check user's transaction key",
			checkUser:         user1,
			checkPubKey:       transactionKey.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(0),
			permission:        types.TransactionPermission,
			expectUser:        user1,
			expectResult:      nil,
			expectGrantPubKey: nil,
		},
		{
			testName:          "check user's app key",
			checkUser:         user1,
			checkPubKey:       appKey.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(0),
			permission:        types.AppPermission,
			expectUser:        user1,
			expectResult:      nil,
			expectGrantPubKey: nil,
		},
		{
			testName:          "user's transaction key can authorize grant app permission",
			checkUser:         user1,
			checkPubKey:       transactionKey.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(0),
			permission:        types.GrantAppPermission,
			expectUser:        user1,
			expectResult:      nil,
			expectGrantPubKey: nil,
		},
		{
			testName:          "user's transaction key can authorize app permission",
			checkUser:         user1,
			checkPubKey:       transactionKey.PubKey(),
			atWhen:            baseTime,
			permission:        types.AppPermission,
			expectUser:        user1,
			expectResult:      nil,
			expectGrantPubKey: nil,
		},
		{
			testName:          "check user's transaction key can't authorize reset permission",
			checkUser:         user1,
			checkPubKey:       transactionKey.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(0),
			permission:        types.ResetPermission,
			expectUser:        user1,
			expectResult:      ErrCheckResetKey(),
			expectGrantPubKey: nil,
		},
		{
			testName:          "check user's app key can authorize grant app permission",
			checkUser:         user1,
			checkPubKey:       appKey.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(0),
			permission:        types.GrantAppPermission,
			expectUser:        user1,
			expectResult:      nil,
			expectGrantPubKey: nil,
		},
		{
			testName:          "check user's app key can't authorize transaction permission",
			checkUser:         user1,
			checkPubKey:       appKey.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(0),
			permission:        types.TransactionPermission,
			expectUser:        user1,
			expectResult:      ErrCheckTransactionKey(),
			expectGrantPubKey: nil,
		},
		{
			testName:          "check user's app key can't authorize reset permission",
			checkUser:         user1,
			checkPubKey:       appKey.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(0),
			permission:        types.ResetPermission,
			expectUser:        user1,
			expectResult:      ErrCheckResetKey(),
			expectGrantPubKey: nil,
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
			expectGrantPubKey: &model.GrantPubKey{
				Username:   appPermissionUser,
				Permission: types.AppPermission,
				CreatedAt:  baseTime,
				ExpiresAt:  baseTime + 100,
				Amount:     types.NewCoinFromInt64(0),
			},
		},
		{
			testName:          "check transaction pubkey of user with app permission",
			checkUser:         user1,
			checkPubKey:       unauthTxPriv.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(0),
			permission:        types.PreAuthorizationPermission,
			expectUser:        "",
			expectResult:      nil,
			expectGrantPubKey: nil,
		},
		{
			testName:          "check unauthorized user app pubkey",
			checkUser:         user1,
			checkPubKey:       unauthPriv2.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(10),
			permission:        types.AppPermission,
			expectUser:        "",
			expectResult:      model.ErrGrantPubKeyNotFound(),
			expectGrantPubKey: nil,
		},
		{
			testName:          "check unauthorized user transaction pubkey",
			checkUser:         user1,
			checkPubKey:       unauthPriv1.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(10),
			permission:        types.PreAuthorizationPermission,
			expectUser:        "",
			expectResult:      model.ErrGrantPubKeyNotFound(),
			expectGrantPubKey: nil,
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
			expectGrantPubKey: &model.GrantPubKey{
				Username:   preAuthPermissionUser,
				Permission: types.PreAuthorizationPermission,
				CreatedAt:  baseTime,
				ExpiresAt:  baseTime + 100,
				Amount:     preAuthAmount.Minus(types.NewCoinFromInt64(10)),
			},
		},
		{
			testName:          "check app pubkey of user with preauthorization permission",
			checkUser:         user1,
			checkPubKey:       unauthAppPriv.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(10),
			permission:        types.AppPermission,
			expectUser:        preAuthPermissionUser,
			expectResult:      model.ErrGrantPubKeyNotFound(),
			expectGrantPubKey: nil,
		},
		{
			testName:          "check app pubkey of user with preauthorization permission",
			checkUser:         user1,
			checkPubKey:       unauthAppPriv.PubKey(),
			atWhen:            baseTime,
			amount:            types.NewCoinFromInt64(10),
			permission:        types.AppPermission,
			expectUser:        preAuthPermissionUser,
			expectResult:      model.ErrGrantPubKeyNotFound(),
			expectGrantPubKey: nil,
		},
		{
			testName:    "check amount exceeds preauthorization limitation",
			checkUser:   user1,
			checkPubKey: authTxPriv.PubKey(),
			atWhen:      baseTime,
			amount:      preAuthAmount,
			permission:  types.PreAuthorizationPermission,
			expectUser:  "",
			expectResult: ErrPreAuthAmountInsufficient(
				preAuthPermissionUser, preAuthAmount.Minus(types.NewCoinFromInt64(10)),
				preAuthAmount),
			expectGrantPubKey: &model.GrantPubKey{
				Username:   preAuthPermissionUser,
				Permission: types.PreAuthorizationPermission,
				CreatedAt:  baseTime,
				ExpiresAt:  baseTime + 100,
				Amount:     preAuthAmount.Minus(types.NewCoinFromInt64(10)),
			},
		},
		{
			testName:     "check grant app key can't sign grant app msg",
			checkUser:    user1,
			checkPubKey:  authAppPriv.PubKey(),
			atWhen:       baseTime,
			permission:   types.GrantAppPermission,
			expectUser:   "",
			expectResult: nil,
			expectGrantPubKey: &model.GrantPubKey{
				Username:   appPermissionUser,
				Permission: types.AppPermission,
				CreatedAt:  baseTime,
				ExpiresAt:  baseTime + 100,
				Amount:     types.NewCoinFromInt64(0),
			},
		},
		{
			testName:          "check expired app permission",
			checkUser:         user1,
			checkPubKey:       authAppPriv.PubKey(),
			atWhen:            baseTime + 101,
			permission:        types.AppPermission,
			expectUser:        "",
			expectResult:      ErrGrantKeyExpired(user1),
			expectGrantPubKey: nil,
		},
		{
			testName:          "check expired preauth permission",
			checkUser:         user1,
			checkPubKey:       authTxPriv.PubKey(),
			atWhen:            baseTime + 101,
			amount:            types.NewCoinFromInt64(100),
			permission:        types.PreAuthorizationPermission,
			expectUser:        "",
			expectResult:      ErrGrantKeyExpired(user1),
			expectGrantPubKey: nil,
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: tc.atWhen})
		grantPubKey, err := am.CheckSigningPubKeyOwner(ctx, tc.checkUser, tc.checkPubKey, tc.permission, tc.amount)
		if tc.expectResult == nil {
			if tc.expectUser != grantPubKey {
				t.Errorf("%s: diff key owner,  got %v, want %v", tc.testName, grantPubKey, tc.expectUser)
				return
			}
		} else {
			if !assert.Equal(t, tc.expectResult.Result(), err.Result()) {
				t.Errorf("%s: diff result,  got %v, want %v", tc.testName, err.Result(), tc.expectResult.Result())
			}
		}

		grantPubKeyInfo, err := am.storage.GetGrantPubKey(ctx, tc.checkUser, tc.checkPubKey)
		if tc.expectGrantPubKey == nil {
			if err == nil {
				t.Errorf("%s: got nil err", tc.testName)
			}
		} else {
			if err != nil {
				t.Errorf("%s: got non-empty err %v", tc.testName, err)
			}
			if !assert.Equal(t, *tc.expectGrantPubKey, *grantPubKeyInfo) {
				t.Errorf("%s: diff grant key,  got %v, want %v", tc.testName, *grantPubKeyInfo, *tc.expectGrantPubKey)
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
	_, _, appPriv2 := createTestAccount(ctx, am, string(userWithAppPermission))
	_, txPriv, _ := createTestAccount(ctx, am, string(userWithPreAuthPermission))

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
		revokePubKey crypto.PubKey
		atWhen       int64
		expectResult sdk.Error
	}{
		{
			testName:     "normal revoke app permission",
			user:         user1,
			revokePubKey: appPriv2.PubKey(),
			atWhen:       baseTime,
			expectResult: nil,
		},
		{
			testName:     "revoke non-exist pubkey, since it's revoked before",
			user:         user1,
			revokePubKey: appPriv2.PubKey(),
			atWhen:       baseTime,
			expectResult: model.ErrGrantPubKeyNotFound(),
		},
		{
			testName:     "revoke expired pubkey",
			user:         user2,
			revokePubKey: appPriv2.PubKey(),
			atWhen:       baseTime + 101,
			expectResult: nil,
		},
		{
			testName:     "normal revoke preauth permission",
			user:         user1,
			revokePubKey: txPriv.PubKey(),
			atWhen:       baseTime + 101,
			expectResult: nil,
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: tc.atWhen})
		err := am.RevokePermission(ctx, tc.user, tc.revokePubKey)
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
	_, _, appPriv1 := createTestAccount(ctx, am, string(user2))
	_, txPriv1, _ := createTestAccount(ctx, am, string(user3))

	baseTime := ctx.BlockHeader().Time

	testCases := []struct {
		testName       string
		user           types.AccountKey
		grantTo        types.AccountKey
		level          types.Permission
		amount         types.Coin
		validityPeriod int64
		expectResult   sdk.Error
		expectPubKey   crypto.PubKey
	}{
		{
			testName:       "normal grant app permission",
			user:           user1,
			grantTo:        user2,
			level:          types.AppPermission,
			validityPeriod: 100,
			amount:         types.NewCoinFromInt64(0),
			expectResult:   nil,
			expectPubKey:   appPriv1.PubKey(),
		},
		{
			testName:       "override app permission",
			user:           user1,
			grantTo:        user2,
			level:          types.AppPermission,
			validityPeriod: 1000,
			amount:         types.NewCoinFromInt64(0),
			expectResult:   nil,
			expectPubKey:   appPriv1.PubKey(),
		},
		{
			testName:       "grant app permission to non-exist user",
			user:           user1,
			grantTo:        nonExistUser,
			level:          types.AppPermission,
			validityPeriod: 1000,
			amount:         types.NewCoinFromInt64(0),
			expectResult:   ErrGetAppKey(nonExistUser),
			expectPubKey:   appPriv1.PubKey(),
		},
		{
			testName:       "grant pre authorization permission",
			user:           user1,
			grantTo:        user3,
			level:          types.PreAuthorizationPermission,
			validityPeriod: 100,
			amount:         types.NewCoinFromInt64(1000),
			expectResult:   nil,
			expectPubKey:   txPriv1.PubKey(),
		},
		{
			testName:       "override pre authorization permission",
			user:           user1,
			grantTo:        user3,
			level:          types.PreAuthorizationPermission,
			validityPeriod: 1000,
			amount:         types.NewCoinFromInt64(10000),
			expectResult:   nil,
			expectPubKey:   txPriv1.PubKey(),
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: baseTime})
		err := am.AuthorizePermission(ctx, tc.user, tc.grantTo, tc.validityPeriod, tc.level, tc.amount)
		if !assert.Equal(t, tc.expectResult, err) {
			t.Errorf("%s: failed to authorize permission, got err %v", tc.testName, err)
		}

		if tc.expectResult == nil {
			grantPubKey, err := am.storage.GetGrantPubKey(ctx, tc.user, tc.expectPubKey)
			if err != nil {
				t.Errorf("%s: failed to get grant pub key, got err %v", tc.testName, err)
			}
			expectGrantPubKey := model.GrantPubKey{
				Username:   tc.grantTo,
				ExpiresAt:  baseTime + tc.validityPeriod,
				CreatedAt:  baseTime,
				Permission: tc.level,
				Amount:     tc.amount,
			}
			if !assert.Equal(t, expectGrantPubKey, *grantPubKey) {
				t.Errorf("%s: diff grant pub key, got %v, want %v", tc.testName, *grantPubKey, expectGrantPubKey)
			}
		}
	}
}
func TestDonationRelationship(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")
	user3 := types.AccountKey("user3")

	createTestAccount(ctx, am, string(user1))
	createTestAccount(ctx, am, string(user2))
	createTestAccount(ctx, am, string(user3))

	testCases := []struct {
		testName         string
		user             types.AccountKey
		donateTo         types.AccountKey
		expectDonateTime int64
	}{
		{
			testName:         "user1 donates to user2",
			user:             user1,
			donateTo:         user2,
			expectDonateTime: 1,
		},
		{
			testName:         "user1 donates to user2 again",
			user:             user1,
			donateTo:         user2,
			expectDonateTime: 2,
		},
		{
			testName:         "user1 donates to user3",
			user:             user1,
			donateTo:         user3,
			expectDonateTime: 1,
		},
		{
			testName:         "user3 donates to user1",
			user:             user3,
			donateTo:         user1,
			expectDonateTime: 1,
		},
		{
			testName:         "user2 donates to user1",
			user:             user2,
			donateTo:         user1,
			expectDonateTime: 1,
		},
	}

	for _, tc := range testCases {
		err := am.UpdateDonationRelationship(ctx, tc.user, tc.donateTo)
		if err != nil {
			t.Errorf("%s: failed to update donation relationship, got err %v", tc.testName, err)
		}

		donateTime, err := am.GetDonationRelationship(ctx, tc.user, tc.donateTo)
		if err != nil {
			t.Errorf("%s: failed to get donation relationship, got err %v", tc.testName, err)
		}
		if donateTime != tc.expectDonateTime {
			t.Errorf("%s: diff donate time, got %v, want %v", tc.testName, donateTime, tc.expectDonateTime)
		}
	}
}

func TestAccountRecoverNormalCase(t *testing.T) {
	testName := "TestAccountRecoverNormalCase"

	ctx, am, accParam := setupTest(t, 1)
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
		CreatedAt:      ctx.BlockHeader().Time,
		ResetKey:       newResetPrivKey.PubKey(),
		TransactionKey: newTransactionPrivKey.PubKey(),
		AppKey:         newAppPrivKey.PubKey(),
	}
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee,
		Stake:   accParam.RegisterFee,
		NumOfTx: 1,
	}

	checkAccountInfo(t, ctx, testName, user1, accInfo)
	checkBankKVByUsername(t, ctx, testName, user1, bank)

	pendingStakeQueue := model.PendingStakeQueue{
		StakeCoinInQueue: sdk.ZeroRat(),
		TotalCoin:        types.NewCoinFromInt64(0),
	}
	checkPendingStake(t, ctx, testName, user1, pendingStakeQueue)

	stake, err := am.GetStake(ctx, user1)
	if err != nil {
		t.Errorf("%s: failed to get stake, got err %v", testName, err)
	}
	if !stake.IsEqual(accParam.RegisterFee) {
		t.Errorf("%s: diff stake, got %v, want %v", testName, stake, accParam.RegisterFee)
	}

	ctx = ctx.WithBlockHeader(
		abci.Header{
			ChainID: "Lino", Height: 1,
			Time: ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake})
	stake, err = am.GetStake(ctx, user1)
	if err != nil {
		t.Errorf("%s: failed to get stake again, got err %v", testName, err)
	}
	if !stake.IsEqual(accParam.RegisterFee) {
		t.Errorf("%s: diff stake again, got %v, want %v", testName, stake, accParam.RegisterFee)
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
		expectSequence int64
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
			testName:     "add the first 100 frozen money",
			frozenAmount: types.NewCoinFromInt64(100),
			startAt:      1000000,
			interval:     10,
			times:        5,
			expectNumOfFrozenAmount: 1,
		},
		{
			testName:     "add the second 100 frozen money, clear the first one",
			frozenAmount: types.NewCoinFromInt64(100),
			startAt:      1200000,
			interval:     10,
			times:        5,
			expectNumOfFrozenAmount: 1,
		},
		{
			testName:     "add the third 100 frozen money",
			frozenAmount: types.NewCoinFromInt64(100),
			startAt:      1300000,
			interval:     10,
			times:        5,
			expectNumOfFrozenAmount: 2,
		},
		{
			testName:     "add the fourth 100 frozen money, clear the second one",
			frozenAmount: types.NewCoinFromInt64(100),
			startAt:      1400000,
			interval:     10,
			times:        5,
			expectNumOfFrozenAmount: 2,
		},
		{
			testName:     "add the fifth 100 frozen money, clear the third and fourth ones",
			frozenAmount: types.NewCoinFromInt64(100),
			startAt:      1600000,
			interval:     10,
			times:        5,
			expectNumOfFrozenAmount: 1,
		}, // this one is used to re-produce the out-of-bound bug.
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: tc.startAt})
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
