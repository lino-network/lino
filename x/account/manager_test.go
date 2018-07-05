package account

import (
	"fmt"
	"testing"
	"time"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
)

func checkBankKVByUsername(
	t *testing.T, ctx sdk.Context, username types.AccountKey, bank model.AccountBank) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	bankPtr, err := accStorage.GetBankFromAccountKey(ctx, username)
	assert.Nil(t, err)
	assert.Equal(t, bank, *bankPtr, "bank should be equal")
}

func checkBalanceHistory(
	t *testing.T, ctx sdk.Context, username types.AccountKey,
	timeSlot int64, balanceHistory model.BalanceHistory) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	balanceHistoryPtr, err := accStorage.GetBalanceHistory(ctx, username, timeSlot)
	assert.Nil(t, err)
	assert.Equal(t, balanceHistory, *balanceHistoryPtr, "balance history should be equal")
}

func checkPendingStake(
	t *testing.T, ctx sdk.Context, username types.AccountKey, pendingStakeQueue model.PendingStakeQueue) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	pendingStakeQueuePtr, err := accStorage.GetPendingStakeQueue(ctx, username)
	assert.Nil(t, err)
	assert.Equal(t, pendingStakeQueue, *pendingStakeQueuePtr, "pending stake should be equal")
}

func checkAccountInfo(
	t *testing.T, ctx sdk.Context, accKey types.AccountKey, accInfo model.AccountInfo) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	info, err := accStorage.GetInfo(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *info, "accout info should be equal")
}

func checkAccountMeta(
	t *testing.T, ctx sdk.Context, accKey types.AccountKey, accMeta model.AccountMeta) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	metaPtr, err := accStorage.GetMeta(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *metaPtr, "accout meta should be equal")
}

func checkAccountReward(
	t *testing.T, ctx sdk.Context, accKey types.AccountKey, reward model.Reward) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	rewardPtr, err := accStorage.GetReward(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, reward, *rewardPtr, "accout reward should be equal")
}

func TestDoesAccountExist(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	assert.False(t, am.DoesAccountExist(ctx, types.AccountKey("user1")))
	createTestAccount(ctx, am, "user1")
	assert.True(t, am.DoesAccountExist(ctx, types.AccountKey("user1")))
}

func TestAddCoin(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	fromUser1, fromUser2, testUser :=
		types.AccountKey("fromUser1"), types.AccountKey("fromuser2"), types.AccountKey("testUser")

	baseTime := time.Now().Unix()
	baseTime1 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake/2
	baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake + 1
	baseTime3 := baseTime2 + coinDayParams.SecondsToRecoverCoinDayStake + 1
	ctx = ctx.WithBlockHeader(abci.Header{Time: baseTime})
	createTestAccount(ctx, am, string(testUser))
	cases := []struct {
		testName                 string
		Amount                   types.Coin
		From                     types.AccountKey
		DetailType               types.TransferDetailType
		Memo                     string
		AtWhen                   int64
		ExpectBank               model.AccountBank
		ExpectPendingStakeQueue  model.PendingStakeQueue
		ExpectBalanceHistorySlot model.BalanceHistory
	}{
		{"add coin to account's saving",
			c100, fromUser1, types.TransferIn, "memo", baseTime,
			model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c100),
				Stake:   accParam.RegisterFee,
				NumOfTx: 2,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat(),
				TotalCoin:        c100,
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
				},
			},
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					model.Detail{
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
		{"add coin to exist account's saving while previous tx is still in pending queue", c100,
			fromUser2, types.DonationIn, "permlink", baseTime1,
			model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c200),
				Stake:   accParam.RegisterFee,
				NumOfTx: 3,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime1,
				StakeCoinInQueue: sdk.NewRat(5000000, 1),
				TotalCoin:        c100.Plus(c100),
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
					model.PendingStake{
						StartTime: baseTime1,
						EndTime:   baseTime1 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
				},
			},
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					model.Detail{
						Amount:     c100,
						From:       fromUser1,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
					model.Detail{
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
		{"add coin to exist account's saving while previous tx just finished pending", c100, "",
			types.ClaimReward, "", baseTime2,
			model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c300),
				Stake:   accParam.RegisterFee.Plus(c100),
				NumOfTx: 4,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime2,
				StakeCoinInQueue: sdk.NewRat(945003125, 189),
				TotalCoin:        c100.Plus(c100),
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime1,
						EndTime:   baseTime1 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
					model.PendingStake{
						StartTime: baseTime2,
						EndTime:   baseTime2 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
				},
			},
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					model.Detail{
						Amount:     c100,
						From:       fromUser1,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
					model.Detail{
						Amount:     c100,
						From:       fromUser2,
						To:         testUser,
						CreatedAt:  baseTime1,
						DetailType: types.DonationIn,
						Memo:       "permlink",
					},
					model.Detail{
						Amount:     c100,
						From:       "",
						To:         testUser,
						CreatedAt:  baseTime2,
						DetailType: types.ClaimReward,
					},
				},
			},
		},
		{"add coin is zero", c0, "",
			types.DelegationReturnCoin, "", baseTime3,
			model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c300),
				Stake:   accParam.RegisterFee.Plus(c100),
				NumOfTx: 4,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime2,
				StakeCoinInQueue: sdk.NewRat(945003125, 189),
				TotalCoin:        c100.Plus(c100),
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime1,
						EndTime:   baseTime1 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
					model.PendingStake{
						StartTime: baseTime2,
						EndTime:   baseTime2 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
				},
			},
			model.BalanceHistory{
				[]model.Detail{

					model.Detail{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					model.Detail{
						Amount:     c100,
						From:       fromUser1,
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
					model.Detail{
						Amount:     c100,
						From:       fromUser2,
						To:         testUser,
						CreatedAt:  baseTime1,
						DetailType: types.DonationIn,
						Memo:       "permlink",
					},
					model.Detail{
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

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		err = am.AddSavingCoin(
			ctx, testUser, cs.Amount, cs.From, cs.Memo, cs.DetailType)

		if err != nil {
			t.Errorf("%s: add coin failed, err: %v", cs.testName, err)
			return
		}
		checkBankKVByUsername(t, ctx, types.AccountKey(testUser), cs.ExpectBank)
		checkPendingStake(t, ctx, types.AccountKey(testUser), cs.ExpectPendingStakeQueue)
		checkBalanceHistory(
			t, ctx, types.AccountKey(testUser), 0, cs.ExpectBalanceHistorySlot)
	}
}

func TestMinusCoin(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	userWithSufficientSaving := types.AccountKey("user1")
	userWithLimitSaving := types.AccountKey("user3")
	fromUser, toUser := types.AccountKey("fromUser"), types.AccountKey("toUser")

	// Get the minimum time of this history slot
	baseTime := time.Now().Unix()
	// baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake + 1
	// baseTime3 := baseTime + accParam.BalanceHistoryIntervalTime + 1

	ctx = ctx.WithBlockHeader(abci.Header{Time: baseTime})
	priv1 := createTestAccount(ctx, am, string(userWithSufficientSaving))
	priv3 := createTestAccount(ctx, am, string(userWithLimitSaving))
	err = am.AddSavingCoin(
		ctx, userWithSufficientSaving, accParam.RegisterFee, fromUser, "", types.TransferIn)
	assert.Nil(t, err)

	cases := []struct {
		TestName                string
		FromUser                types.AccountKey
		UserPriv                crypto.PrivKey
		ExpectErr               sdk.Error
		Amount                  types.Coin
		AtWhen                  int64
		To                      types.AccountKey
		Memo                    string
		DetailType              types.TransferDetailType
		ExpectBank              model.AccountBank
		ExpectPendingStakeQueue model.PendingStakeQueue
		ExpectBalanceHistory    model.BalanceHistory
	}{
		{"minus saving coin from user with sufficient saving",
			userWithSufficientSaving, priv1, nil, coin1, baseTime, toUser, "memo", types.TransferOut,
			model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(accParam.RegisterFee).Minus(coin1),
				NumOfTx: 3,
				Stake:   accParam.RegisterFee,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat(),
				TotalCoin:        accParam.RegisterFee.Minus(coin1),
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      accParam.RegisterFee.Minus(coin1),
					}},
			},
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						Amount:     accParam.RegisterFee,
						From:       accountReferrer,
						To:         userWithSufficientSaving,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       types.InitAccountWithFullStakeMemo,
					},
					model.Detail{
						Amount:     accParam.RegisterFee,
						From:       fromUser,
						To:         userWithSufficientSaving,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					model.Detail{
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
		{"minus saving coin from user with limit saving",
			userWithLimitSaving, priv3, ErrAccountSavingCoinNotEnough(),
			coin1, baseTime, toUser, "memo", types.TransferOut,
			model.AccountBank{
				Saving:  accParam.RegisterFee,
				NumOfTx: 1,
				Stake:   accParam.RegisterFee,
			},
			model.PendingStakeQueue{
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
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
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
		{"minus saving coin exceeds the coin user hold",
			userWithLimitSaving, priv3, ErrAccountSavingCoinNotEnough(),
			c100, baseTime, toUser, "memo", types.TransferOut,
			model.AccountBank{
				Saving:  accParam.RegisterFee,
				NumOfTx: 1,
				Stake:   accParam.RegisterFee,
			},
			model.PendingStakeQueue{
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
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
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
	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		err = am.MinusSavingCoin(ctx, cs.FromUser, cs.Amount, cs.To, cs.Memo, cs.DetailType)

		assert.Equal(t, cs.ExpectErr, err, fmt.Sprintf("%s: minus coin failed, err: %v", cs.TestName, err))
		if cs.ExpectErr == nil {
			checkBankKVByUsername(t, ctx, cs.FromUser, cs.ExpectBank)
			checkPendingStake(t, ctx, cs.FromUser, cs.ExpectPendingStakeQueue)
			checkBalanceHistory(t, ctx, cs.FromUser,
				0, cs.ExpectBalanceHistory)
		}
	}
}

func TestBalanceHistory(t *testing.T) {
	fromUser, toUser := types.AccountKey("fromUser"), types.AccountKey("toUser")
	cases := []struct {
		TestName        string
		NumOfAdding     int
		NumOfMinus      int
		expectTotalSlot int64
	}{
		{"test only one adding", 1, 0, 1},
		{"test 99 adding, which fullfills 1 bundles", 99, 0, 1},
		{"test adding and minus, which results in 2 bundles", 50, 50, 2},
	}
	for _, cs := range cases {
		ctx, am, accParam := setupTest(t, 1)

		user1 := types.AccountKey("user1")
		createTestAccount(ctx, am, string(user1))

		for i := 0; i < cs.NumOfAdding; i++ {
			err := am.AddSavingCoin(ctx, user1, coin1, fromUser, "", types.TransferIn)
			assert.Nil(t, err)
		}
		for i := 0; i < cs.NumOfMinus; i++ {
			err := am.MinusSavingCoin(ctx, user1, coin1, toUser, "", types.TransferOut)
			assert.Nil(t, err)
		}
		bank, err := am.storage.GetBankFromAccountKey(ctx, user1)
		assert.Nil(t, err)
		// add one init transfer in
		expectNumOfTx := int64(cs.NumOfAdding + cs.NumOfMinus + 1)
		assert.Equal(t, expectNumOfTx, bank.NumOfTx)

		// total slot should use previous states to get expected slots
		actualTotalSlot := (expectNumOfTx-1)/accParam.BalanceHistoryBundleSize + 1
		assert.Equal(t, cs.expectTotalSlot, actualTotalSlot)
		actualNumOfAdding, actualNumOfMinus := 0, 0
		for slot := int64(0); slot < actualTotalSlot; slot++ {
			balanceHistory, err := am.storage.GetBalanceHistory(ctx, user1, slot)
			assert.Nil(t, err)
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
		assert.Equal(t, cs.NumOfAdding+1, actualNumOfAdding)
		assert.Equal(t, cs.NumOfMinus, actualNumOfMinus)
	}
}

func TestAddBalanceHistory(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)

	cases := []struct {
		testName              string
		numOfTx               int64
		detail                model.Detail
		expectNumOfTxInBundle int
	}{
		{"try first transaction in first slot",
			0, model.Detail{
				From:       "test1",
				To:         "test2",
				Amount:     types.NewCoinFromInt64(1),
				DetailType: types.TransferIn,
				CreatedAt:  time.Now().Unix(),
			}, 1,
		},
		{"try second transaction in first slot",
			1, model.Detail{
				From:       "test2",
				To:         "test1",
				Amount:     types.NewCoinFromInt64(1 * types.Decimals),
				DetailType: types.TransferOut,
				CreatedAt:  time.Now().Unix(),
			}, 2,
		},
		{"add transaction to the end of the first slot limitation",
			99, model.Detail{
				From:       "test1",
				To:         "post",
				Amount:     types.NewCoinFromInt64(1 * types.Decimals),
				DetailType: types.DonationOut,
				CreatedAt:  time.Now().Unix(),
				Memo:       "",
			}, 3,
		},
		{"add transaction to next slot",
			100, model.Detail{
				From:       "",
				To:         "test1",
				Amount:     types.NewCoinFromInt64(1 * types.Decimals),
				DetailType: types.DeveloperDeposit,
				CreatedAt:  time.Now().Unix(),
			}, 1,
		},
	}

	for _, cs := range cases {
		err := am.AddBalanceHistory(ctx, user1, cs.numOfTx, cs.detail)
		assert.Nil(t, err)
		balanceHistory, err :=
			am.storage.GetBalanceHistory(
				ctx, user1, cs.numOfTx/accParam.BalanceHistoryBundleSize)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectNumOfTxInBundle, len(balanceHistory.Details))
		assert.Equal(t, cs.detail, balanceHistory.Details[cs.expectNumOfTxInBundle-1])
	}
}

func TestCreateAccountNormalCase(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	priv := crypto.GenPrivKeyEd25519()
	accKey := types.AccountKey("accKey")

	// normal test
	assert.False(t, am.DoesAccountExist(ctx, accKey))
	err := am.CreateAccount(
		ctx, accountReferrer, accKey, priv.PubKey(), priv.Generate(0).PubKey(),
		priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), accParam.RegisterFee)
	assert.Nil(t, err)

	assert.True(t, am.DoesAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee,
		NumOfTx: 1,
		Stake:   accParam.RegisterFee,
	}
	checkBankKVByUsername(t, ctx, accKey, bank)
	pendingStakeQueue := model.PendingStakeQueue{StakeCoinInQueue: sdk.ZeroRat()}
	checkPendingStake(t, ctx, accKey, pendingStakeQueue)
	accInfo := model.AccountInfo{
		Username:        accKey,
		CreatedAt:       ctx.BlockHeader().Time,
		MasterKey:       priv.PubKey(),
		TransactionKey:  priv.Generate(0).PubKey(),
		MicropaymentKey: priv.Generate(1).PubKey(),
		PostKey:         priv.Generate(2).PubKey(),
	}
	checkAccountInfo(t, ctx, accKey, accInfo)
	accMeta := model.AccountMeta{
		LastActivityAt:       ctx.BlockHeader().Time,
		LastReportOrUpvoteAt: ctx.BlockHeader().Time,
	}
	checkAccountMeta(t, ctx, accKey, accMeta)

	reward := model.Reward{coin0, coin0, coin0, coin0}
	checkAccountReward(t, ctx, accKey, reward)

	balanceHistory, err := am.storage.GetBalanceHistory(ctx, accKey, 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(balanceHistory.Details))
	assert.Equal(t, model.Detail{
		From:       accountReferrer,
		To:         accKey,
		Amount:     accParam.RegisterFee,
		CreatedAt:  ctx.BlockHeader().Time,
		DetailType: types.TransferIn,
		Memo:       types.InitAccountWithFullStakeMemo,
	}, balanceHistory.Details[0])
}

func TestCreateAccountWithLargeRegisterFee(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	priv := crypto.GenPrivKeyEd25519()
	accKey := types.AccountKey("accKey")

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	extraRegisterFee := types.NewCoinFromInt64(100 * types.Decimals)
	// normal test
	assert.False(t, am.DoesAccountExist(ctx, accKey))
	err = am.CreateAccount(
		ctx, accountReferrer, accKey, priv.PubKey(), priv.Generate(0).PubKey(),
		priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), accParam.RegisterFee.Plus(extraRegisterFee))
	assert.Nil(t, err)

	assert.True(t, am.DoesAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee.Plus(extraRegisterFee),
		NumOfTx: 2,
		Stake:   accParam.RegisterFee,
	}
	checkBankKVByUsername(t, ctx, accKey, bank)
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
	checkPendingStake(t, ctx, accKey, pendingStakeQueue)
	accInfo := model.AccountInfo{
		Username:        accKey,
		CreatedAt:       ctx.BlockHeader().Time,
		MasterKey:       priv.PubKey(),
		TransactionKey:  priv.Generate(0).PubKey(),
		MicropaymentKey: priv.Generate(1).PubKey(),
		PostKey:         priv.Generate(2).PubKey(),
	}
	checkAccountInfo(t, ctx, accKey, accInfo)
	accMeta := model.AccountMeta{
		LastActivityAt:       ctx.BlockHeader().Time,
		LastReportOrUpvoteAt: ctx.BlockHeader().Time,
	}
	checkAccountMeta(t, ctx, accKey, accMeta)

	reward := model.Reward{coin0, coin0, coin0, coin0}
	checkAccountReward(t, ctx, accKey, reward)

	balanceHistory, err := am.storage.GetBalanceHistory(ctx, accKey, 0)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(balanceHistory.Details))
	assert.Equal(t, model.Detail{
		From:       accountReferrer,
		To:         accKey,
		Amount:     accParam.RegisterFee,
		CreatedAt:  ctx.BlockHeader().Time,
		DetailType: types.TransferIn,
		Memo:       types.InitAccountWithFullStakeMemo,
	}, balanceHistory.Details[0])
	assert.Equal(t, model.Detail{
		From:       accountReferrer,
		To:         accKey,
		Amount:     extraRegisterFee,
		CreatedAt:  ctx.BlockHeader().Time,
		DetailType: types.TransferIn,
		Memo:       types.InitAccountRegisterDepositMemo,
	}, balanceHistory.Details[1])
}

func TestInvalidCreateAccount(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	priv1 := crypto.GenPrivKeyEd25519()
	priv2 := crypto.GenPrivKeyEd25519()

	accKey1 := types.AccountKey("accKey1")
	accKey2 := types.AccountKey("accKey2")
	accKey3 := types.AccountKey("accKey3")

	cases := []struct {
		testName    string
		username    types.AccountKey
		privkey     crypto.PrivKey
		registerFee types.Coin
		expectErr   sdk.Error
	}{
		{"register user with sufficient saving coin",
			accKey1, priv1, accParam.RegisterFee, nil,
		},
		{"username already took",
			accKey1, priv1, accParam.RegisterFee, ErrAccountAlreadyExists(accKey1),
		},
		{"username already took with different private key",
			accKey1, priv2, accParam.RegisterFee, ErrAccountAlreadyExists(accKey1),
		},
		{"register the same private key",
			accKey2, priv1, accParam.RegisterFee, nil,
		},
		{"insufficient register fee",
			accKey3, priv1, types.NewCoinFromInt64(1), ErrRegisterFeeInsufficient(),
		},
	}
	for _, cs := range cases {
		err := am.CreateAccount(
			ctx, accountReferrer, cs.username, cs.privkey.PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), cs.registerFee)
		assert.Equal(t, cs.expectErr, err,
			fmt.Sprintf("%s: create account failed: expect %v, got %v",
				cs.testName, cs.expectErr, err))
	}
}

func TestUpdateJSONMeta(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)

	accKey := types.AccountKey("accKey")
	createTestAccount(ctx, am, string(accKey))

	cases := []struct {
		testName string
		username types.AccountKey
		JSONMeta string
	}{
		{"normal update",
			accKey, "{'link':'https://lino.network'}",
		},
	}
	for _, cs := range cases {
		err := am.UpdateJSONMeta(ctx, cs.username, cs.JSONMeta)
		assert.Nil(t, err)
		accMeta, err := am.storage.GetMeta(ctx, cs.username)
		assert.Nil(t, err)
		assert.Equal(t, cs.JSONMeta, accMeta.JSONMeta)
	}
}

func TestCoinDayByAccountKey(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)
	totalCoinDaysSec := coinDayParams.SecondsToRecoverCoinDayStake
	registerFee := accParam.RegisterFee.ToInt64()
	doubleRegisterFee := types.NewCoinFromInt64(registerFee * 2)
	halfRegisterFee := types.NewCoinFromInt64(registerFee / 2)

	baseTime := ctx.BlockHeader().Time
	baseTime2 := baseTime + totalCoinDaysSec + (totalCoinDaysSec/registerFee)/2 + 1

	createTestAccount(ctx, am, string(accKey))

	cases := []struct {
		testName            string
		IsAdd               bool
		Coin                types.Coin
		AtWhen              int64
		ExpectSavingBalance types.Coin
		ExpectStake         types.Coin
		ExpectStakeInBank   types.Coin
		ExpectNumOfTx       int64
	}{
		{"add coin before charging first coin",
			true, accParam.RegisterFee, baseTime + (totalCoinDaysSec/registerFee)/2,
			doubleRegisterFee, accParam.RegisterFee, accParam.RegisterFee, 2},
		{"check first coin",
			true, coin0, baseTime + (totalCoinDaysSec/registerFee)/2 + 1,
			doubleRegisterFee, accParam.RegisterFee, accParam.RegisterFee, 2},
		{"check both transactions fully charged",
			true, coin0, baseTime2, doubleRegisterFee, doubleRegisterFee, doubleRegisterFee, 2},
		{"withdraw half deposit",
			false, accParam.RegisterFee, baseTime2,
			accParam.RegisterFee, accParam.RegisterFee, accParam.RegisterFee, 3},
		{"charge again",
			true, accParam.RegisterFee, baseTime2,
			doubleRegisterFee, accParam.RegisterFee, accParam.RegisterFee, 4},
		{"withdraw half deposit while the last transaction is still charging",
			false, halfRegisterFee, baseTime2 + totalCoinDaysSec/2 + 1,
			accParam.RegisterFee.Plus(halfRegisterFee),
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(registerFee / 4)), accParam.RegisterFee, 5},
		{"withdraw last transaction which is still charging",
			false, halfRegisterFee, baseTime2 + totalCoinDaysSec/2 + 1,
			accParam.RegisterFee, accParam.RegisterFee, accParam.RegisterFee, 6},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		if cs.IsAdd {
			err := am.AddSavingCoin(ctx, accKey, cs.Coin, "", "", types.TransferIn)
			assert.Nil(t, err)
		} else {
			err := am.MinusSavingCoin(ctx, accKey, cs.Coin, "", "", types.TransferOut)
			assert.Nil(t, err)
		}
		coin, err := am.GetStake(ctx, accKey)
		assert.Nil(t, err)
		if !cs.ExpectStake.IsEqual(coin) {
			t.Errorf("%s: expect stake incorrect, expect %v, got %v", cs.testName, cs.ExpectStake, coin)
			return
		}

		bank := model.AccountBank{
			Saving:  cs.ExpectSavingBalance,
			Stake:   cs.ExpectStakeInBank,
			NumOfTx: cs.ExpectNumOfTx,
		}
		checkBankKVByUsername(t, ctx, accKey, bank)
	}
}

func TestAccountReward(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	createTestAccount(ctx, am, string(accKey))

	err := am.AddIncomeAndReward(ctx, accKey, c500, c200, c300)
	assert.Nil(t, err)
	reward := model.Reward{c500, c200, c300, c300}
	checkAccountReward(t, ctx, accKey, reward)
	err = am.AddIncomeAndReward(ctx, accKey, c500, c300, c200)
	assert.Nil(t, err)
	reward = model.Reward{c1000, c500, c500, c500}
	checkAccountReward(t, ctx, accKey, reward)

	bank := model.AccountBank{
		Saving:  accParam.RegisterFee,
		NumOfTx: 1,
		Stake:   accParam.RegisterFee,
	}
	checkBankKVByUsername(t, ctx, accKey, bank)

	err = am.ClaimReward(ctx, accKey)
	assert.Nil(t, err)
	bank.Saving = accParam.RegisterFee.Plus(c500)
	bank.NumOfTx = 2
	checkBankKVByUsername(t, ctx, accKey, bank)
	reward = model.Reward{c1000, c500, c500, c0}
	checkAccountReward(t, ctx, accKey, reward)
}

func TestCheckUserTPSCapacity(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	bandwidthParams, err := am.paramHolder.GetBandwidthParam(ctx)
	assert.Nil(t, err)
	secondsToRecoverBandwidth := bandwidthParams.SecondsToRecoverBandwidth

	baseTime := ctx.BlockHeader().Time

	createTestAccount(ctx, am, string(accKey))
	err = am.AddSavingCoin(ctx, accKey, c100, "", "", types.TransferIn)
	assert.Nil(t, err)

	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	err = accStorage.SetPendingStakeQueue(
		ctx, accKey, &model.PendingStakeQueue{})
	assert.Nil(t, err)

	cases := []struct {
		TPSCapacityRatio     sdk.Rat
		UserStake            types.Coin
		LastActivity         int64
		LastCapacity         types.Coin
		CurrentTime          int64
		ExpectResult         sdk.Error
		ExpectRemainCapacity types.Coin
	}{
		{sdk.NewRat(1, 10), types.NewCoinFromInt64(10 * types.Decimals), baseTime, types.NewCoinFromInt64(0),
			baseTime, ErrAccountTPSCapacityNotEnough(accKey), types.NewCoinFromInt64(0)},
		{sdk.NewRat(1, 10), types.NewCoinFromInt64(10 * types.Decimals), baseTime, types.NewCoinFromInt64(0),
			baseTime + secondsToRecoverBandwidth, nil, types.NewCoinFromInt64(990000)},
		{sdk.NewRat(1, 2), types.NewCoinFromInt64(10 * types.Decimals), baseTime, types.NewCoinFromInt64(0),
			baseTime + secondsToRecoverBandwidth, nil, types.NewCoinFromInt64(950000)},
		{sdk.NewRat(1, 1), types.NewCoinFromInt64(10 * types.Decimals), baseTime, types.NewCoinFromInt64(0),
			baseTime + secondsToRecoverBandwidth, nil, types.NewCoinFromInt64(9 * types.Decimals)},
		{sdk.NewRat(1, 1), types.NewCoinFromInt64(1 * types.Decimals), baseTime,
			types.NewCoinFromInt64(10 * types.Decimals), baseTime, nil, types.NewCoinFromInt64(0)},
		{sdk.NewRat(1, 1), types.NewCoinFromInt64(10), baseTime, types.NewCoinFromInt64(1 * types.Decimals),
			baseTime, ErrAccountTPSCapacityNotEnough(accKey), types.NewCoinFromInt64(1 * types.Decimals)},
		{sdk.NewRat(1, 1), types.NewCoinFromInt64(1 * types.Decimals), baseTime, types.NewCoinFromInt64(0),
			baseTime + secondsToRecoverBandwidth/2,
			ErrAccountTPSCapacityNotEnough(accKey), types.NewCoinFromInt64(0)},
		{sdk.NewRat(1, 2), types.NewCoinFromInt64(1 * types.Decimals), baseTime, types.NewCoinFromInt64(0),
			baseTime + secondsToRecoverBandwidth/2, nil, types.NewCoinFromInt64(0)},
		{sdk.NewRat(1, 1), types.NewCoinFromInt64(1 * types.Decimals), 0, types.NewCoinFromInt64(0),
			baseTime, nil, types.NewCoinFromInt64(0)},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: cs.CurrentTime})
		bank := &model.AccountBank{
			Saving: cs.UserStake,
			Stake:  cs.UserStake,
		}
		err = accStorage.SetBankFromAccountKey(ctx, accKey, bank)
		assert.Nil(t, err)
		meta := &model.AccountMeta{
			LastActivityAt:      cs.LastActivity,
			TransactionCapacity: cs.LastCapacity,
		}
		err = accStorage.SetMeta(ctx, accKey, meta)
		assert.Nil(t, err)

		err = am.CheckUserTPSCapacity(ctx, accKey, cs.TPSCapacityRatio)
		assert.Equal(t, cs.ExpectResult, err)

		accMeta := model.AccountMeta{
			LastActivityAt:      ctx.BlockHeader().Time,
			TransactionCapacity: cs.ExpectRemainCapacity,
		}
		if cs.ExpectResult != nil {
			accMeta.LastActivityAt = cs.LastActivity
		}
		checkAccountMeta(t, ctx, accKey, accMeta)
	}
}

func TestCheckAuthenticatePubKeyOwner(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	postPermissionUser := types.AccountKey("user2")
	micropaymentPermissionUser := types.AccountKey("user3")
	multiTimesUser := types.AccountKey("user4")
	unauthUser := types.AccountKey("user5")
	fullyAuthUser := types.AccountKey("user6")

	masterKey := crypto.GenPrivKeyEd25519()
	transactionKey := crypto.GenPrivKeyEd25519()
	micropaymentKey := crypto.GenPrivKeyEd25519()
	postKey := crypto.GenPrivKeyEd25519()
	am.CreateAccount(
		ctx, accountReferrer, user1, masterKey.PubKey(), transactionKey.PubKey(),
		micropaymentKey.PubKey(), postKey.PubKey(), accParam.RegisterFee)

	postPriv := createTestAccount(ctx, am, string(postPermissionUser))
	microPriv := createTestAccount(ctx, am, string(micropaymentPermissionUser))
	unauthPriv := createTestAccount(ctx, am, string(unauthUser))
	fullyAuthPriv := createTestAccount(ctx, am, string(fullyAuthUser))
	multiTimesPriv := createTestAccount(ctx, am, string(multiTimesUser))
	defaultGrantTimes := int64(1)
	err := am.AuthorizePermission(ctx, user1, postPermissionUser, 100, defaultGrantTimes, types.PostPermission)
	assert.Nil(t, err)
	err = am.AuthorizePermission(
		ctx, user1, micropaymentPermissionUser, 100, defaultGrantTimes, types.MicropaymentPermission)
	assert.Nil(t, err)
	err = am.AuthorizePermission(ctx, user1, fullyAuthUser, 100, defaultGrantTimes, types.PostPermission)
	assert.Nil(t, err)
	err = am.AuthorizePermission(ctx, user1, fullyAuthUser, 100, defaultGrantTimes, types.MicropaymentPermission)
	assert.Nil(t, err)
	err = am.AuthorizePermission(ctx, user1, multiTimesUser, 100, defaultGrantTimes, types.MicropaymentPermission)
	assert.Nil(t, err)

	baseTime := ctx.BlockHeader().Time

	cases := []struct {
		testName          string
		checkUser         types.AccountKey
		checkPubKey       crypto.PubKey
		atWhen            int64
		permission        types.Permission
		expectUser        types.AccountKey
		expectResult      sdk.Error
		expectGrantPubKey *model.GrantPubKey
	}{
		{"check user's master key",
			user1, masterKey.PubKey(), baseTime, types.MasterPermission, user1, nil, nil},
		{"check user's transaction key",
			user1, transactionKey.PubKey(), baseTime, types.TransactionPermission, user1, nil, nil},
		{"check user's micropayment key",
			user1, micropaymentKey.PubKey(), baseTime, types.MicropaymentPermission, user1, nil, nil},
		{"check user's post key",
			user1, postKey.PubKey(), baseTime, types.PostPermission, user1, nil, nil},
		{"user's transaction key can authorize micropayment permission",
			user1, transactionKey.PubKey(), baseTime, types.MicropaymentPermission, user1, nil, nil},
		{"user's transaction key can authorize grant micropayment permission",
			user1, transactionKey.PubKey(), baseTime, types.GrantMicropaymentPermission, user1, nil, nil},
		{"user's transaction key can authorize grant post permission",
			user1, transactionKey.PubKey(), baseTime, types.GrantPostPermission, user1, nil, nil},
		{"user's transaction key can authorize post permission",
			user1, transactionKey.PubKey(), baseTime, types.PostPermission, user1, nil, nil},
		{"check user's transaction key can't authorize master permission",
			user1, transactionKey.PubKey(), baseTime, types.MasterPermission, user1,
			ErrCheckMasterKey(), nil},
		{"user's micropayment key can authorize post permission",
			user1, micropaymentKey.PubKey(), baseTime, types.PostPermission, user1, nil, nil},
		{"user's micropayment key can authorize grant micropayment permission",
			user1, micropaymentKey.PubKey(), baseTime, types.GrantMicropaymentPermission, user1, nil, nil},
		{"user's micropayment key can authorize grant post permission",
			user1, micropaymentKey.PubKey(), baseTime, types.GrantPostPermission, user1, nil, nil},
		{"user's micropayment key can't authorize master permission",
			user1, micropaymentKey.PubKey(), baseTime, types.MasterPermission, user1, ErrCheckMasterKey(), nil},
		{"user's micropayment key can't authorize transaction permission",
			user1, micropaymentKey.PubKey(), baseTime, types.TransactionPermission, user1, ErrCheckTransactionKey(), nil},
		{"check user's post key can authorize grant post permission",
			user1, postKey.PubKey(), baseTime, types.GrantPostPermission, user1, nil, nil},
		{"check user's post key can't authorize master permission",
			user1, postKey.PubKey(), baseTime, types.MasterPermission, user1,
			ErrCheckMasterKey(), nil},
		{"check user's post key can't authorize transaction permission",
			user1, postKey.PubKey(), baseTime, types.TransactionPermission, user1,
			ErrCheckTransactionKey(), nil},
		{"check user's post key can't authorize micropayment permission",
			user1, postKey.PubKey(), baseTime, types.MicropaymentPermission, user1,
			ErrCheckAuthenticatePubKeyOwner(user1), nil},
		{"check post pubkey of user with post permission",
			user1, postPriv.Generate(2).PubKey(), baseTime, types.PostPermission, postPermissionUser, nil,
			&model.GrantPubKey{
				Username:   postPermissionUser,
				Permission: types.PostPermission,
				LeftTimes:  defaultGrantTimes,
				CreatedAt:  baseTime,
				ExpiresAt:  baseTime + 100,
			}},
		{"check micropayment pubkey of user with post permission",
			user1, postPriv.Generate(1).PubKey(), baseTime, types.PostPermission,
			postPermissionUser, ErrCheckAuthenticatePubKeyOwner(user1), nil},
		{"check micropayment pubkey of user with micropayment permission",
			user1, multiTimesPriv.Generate(1).PubKey(), baseTime, types.MicropaymentPermission, multiTimesUser, nil,
			&model.GrantPubKey{
				Username:   multiTimesUser,
				Permission: types.MicropaymentPermission,
				LeftTimes:  defaultGrantTimes - 1,
				CreatedAt:  baseTime,
				ExpiresAt:  baseTime + 100,
			}},
		{"check post pubkey of user with micropayment permission",
			user1, microPriv.Generate(2).PubKey(), baseTime, types.MicropaymentPermission,
			micropaymentPermissionUser, ErrCheckAuthenticatePubKeyOwner(user1), nil},
		{"check unauthorized user post pubkey",
			user1, unauthPriv.Generate(2).PubKey(), baseTime, types.PostPermission, "",
			ErrCheckAuthenticatePubKeyOwner(user1), nil},
		{"check unauthorized user micropayment pubkey",
			user1, unauthPriv.Generate(1).PubKey(), baseTime, types.MicropaymentPermission, "",
			ErrCheckAuthenticatePubKeyOwner(user1), nil},
		{"check fully authed user micropayment pubkey but post permission",
			user1, fullyAuthPriv.Generate(1).PubKey(), baseTime, types.PostPermission, "",
			ErrGrantKeyMismatch(fullyAuthUser), nil},
		{"check fully authed user post pubkey but micropayment permission",
			user1, fullyAuthPriv.Generate(2).PubKey(), baseTime, types.MicropaymentPermission, "",
			ErrGrantKeyMismatch(fullyAuthUser), nil},
		{"check expired micropayment permission",
			user1, microPriv.Generate(1).PubKey(), baseTime + 101,
			types.MicropaymentPermission, "", ErrGrantKeyExpired(user1), nil},
		{"check expired post permission",
			user1, postPriv.Generate(2).PubKey(), baseTime + 101, types.PostPermission,
			"", ErrGrantKeyExpired(user1), nil},
		{"check micropayment pubkey exceeds limitation",
			user1, multiTimesPriv.Generate(1).PubKey(), baseTime,
			types.MicropaymentPermission, multiTimesUser, ErrGrantKeyExpired(user1), nil},
		{"check grant micropayment key can't sign grant permission msg",
			user1, micropaymentKey.Generate(1).PubKey(), baseTime,
			types.GrantMicropaymentPermission, micropaymentPermissionUser, ErrCheckGrantMicropaymentKey(), nil},
		{"check grant micropayment key can't sign grant post msg",
			user1, micropaymentKey.Generate(1).PubKey(), baseTime,
			types.GrantPostPermission, micropaymentPermissionUser, ErrCheckGrantPostKey(), nil},
		{"check grant post key can't sign grant post msg",
			user1, postKey.Generate(1).PubKey(), baseTime,
			types.GrantPostPermission, postPermissionUser, ErrCheckGrantPostKey(), nil},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: cs.atWhen})
		grantPubKey, err := am.CheckSigningPubKeyOwner(ctx, cs.checkUser, cs.checkPubKey, cs.permission)
		if cs.expectResult == nil {
			if cs.expectUser != grantPubKey {
				t.Errorf(
					"%s: expect key owner incorrect, expect %v, got %v",
					cs.testName, cs.expectUser, grantPubKey)
				return
			}
		} else {
			assert.Equal(t, cs.expectResult.Result(), err.Result())
		}
		grantPubKeyInfo, err := am.storage.GetGrantPubKey(ctx, cs.checkUser, cs.checkPubKey)
		if cs.expectGrantPubKey == nil {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, *cs.expectGrantPubKey, *grantPubKeyInfo)
		}
	}
}

func TestRevokePermission(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	userWithMicropaymentPermission := types.AccountKey("userWithMicropaymentPermission")
	userWithBothPermission := types.AccountKey("userWithBothPermission")

	createTestAccount(ctx, am, string(user1))
	priv2 := createTestAccount(ctx, am, string(userWithMicropaymentPermission))
	priv3 := createTestAccount(ctx, am, string(userWithBothPermission))

	baseTime := ctx.BlockHeader().Time

	err := am.AuthorizePermission(ctx, user1, userWithMicropaymentPermission, 100, 10, types.MicropaymentPermission)
	assert.Nil(t, err)

	err = am.AuthorizePermission(ctx, user1, userWithBothPermission, 100, 10, types.MicropaymentPermission)
	assert.Nil(t, err)
	err = am.AuthorizePermission(ctx, user1, userWithBothPermission, 100, 10, types.PostPermission)
	assert.Nil(t, err)

	cases := []struct {
		testName     string
		user         types.AccountKey
		revokePubkey crypto.PubKey
		atWhen       int64
		level        types.Permission
		expectResult sdk.Error
	}{
		{"normal revoke post permission", user1, priv3.Generate(2).PubKey(), baseTime, types.PostPermission, nil},
		{"normal revoke micropayment permission", user1, priv2.Generate(1).PubKey(), baseTime, types.MicropaymentPermission, nil},
		{"revoke permission mismatch", user1, priv3.Generate(1).PubKey(),
			baseTime, types.PostPermission, ErrRevokePermissionLevelMismatch(types.PostPermission, types.MicropaymentPermission)},
		{"revoke non-exist pubkey", user1, priv3.Generate(2).PubKey(),
			baseTime, types.PostPermission, model.ErrGetGrantPubKeyFailed()},
		{"revoke expired pubkey", user1, priv3.Generate(1).PubKey(),
			baseTime + 101, types.PostPermission, nil},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: cs.atWhen})
		err := am.RevokePermission(ctx, cs.user, cs.revokePubkey, cs.level)
		assert.Equal(t, cs.expectResult, err, cs.testName)
	}
}

func TestAuthorizePermission(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")
	user3 := types.AccountKey("user3")

	createTestAccount(ctx, am, string(user1))
	priv2 := createTestAccount(ctx, am, string(user2))
	priv3 := createTestAccount(ctx, am, string(user3))

	baseTime := ctx.BlockHeader().Time

	cases := []struct {
		testName       string
		user           types.AccountKey
		grantTo        types.AccountKey
		level          types.Permission
		validityPeriod int64
		allowTimes     int64
		expectResult   sdk.Error
		expectPubKey   crypto.PubKey
	}{
		{"normal grant post permission", user1, user2, types.PostPermission,
			100, 10, nil, priv2.Generate(2).PubKey()},
		{"normal grant micropayment permission", user1, user3, types.MicropaymentPermission,
			100, 10, nil, priv3.Generate(1).PubKey()},
		{"override permission", user1, user3, types.MicropaymentPermission,
			1000, 10, nil, priv3.Generate(1).PubKey()},
		{"micropayment authorization exceeds maximum requirement", user1, user3, types.MicropaymentPermission,
			1000, accParam.MaximumMicropaymentGrantTimes + 1,
			ErrGrantTimesExceedsLimitation(accParam.MaximumMicropaymentGrantTimes), priv3.Generate(1).PubKey()},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: baseTime})
		err := am.AuthorizePermission(ctx, cs.user, cs.grantTo, cs.validityPeriod, cs.allowTimes, cs.level)
		assert.Equal(t, cs.expectResult, err, cs.testName)
		if cs.expectResult == nil {
			grantPubKey, err := am.storage.GetGrantPubKey(ctx, cs.user, cs.expectPubKey)
			assert.Nil(t, err)
			expectGrantPubKey := model.GrantPubKey{
				Username:   cs.grantTo,
				ExpiresAt:  baseTime + cs.validityPeriod,
				CreatedAt:  baseTime,
				LeftTimes:  cs.allowTimes,
				Permission: cs.level,
			}
			assert.Equal(t, expectGrantPubKey, *grantPubKey)
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

	cases := []struct {
		user             types.AccountKey
		donateTo         types.AccountKey
		expectDonateTime int64
	}{
		{user1, user2, 1},
		{user1, user2, 2},
		{user1, user3, 1},
		{user3, user1, 1},
		{user2, user1, 1},
	}

	for _, cs := range cases {
		err := am.UpdateDonationRelationship(ctx, cs.user, cs.donateTo)
		assert.Nil(t, err)
		donateTime, err := am.GetDonationRelationship(ctx, cs.user, cs.donateTo)
		assert.Nil(t, err)
		assert.Equal(t, donateTime, cs.expectDonateTime)
	}
}

func TestAccountRecoverNormalCase(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	createTestAccount(ctx, am, string(user1))

	newMasterPrivKey := crypto.GenPrivKeyEd25519()
	newTransactionPrivKey := newMasterPrivKey.Generate(0)
	newMicropaymentPrivKey := newMasterPrivKey.Generate(1)
	newPostPrivKey := newMasterPrivKey.Generate(2)

	err = am.RecoverAccount(
		ctx, user1, newMasterPrivKey.PubKey(), newTransactionPrivKey.PubKey(),
		newMicropaymentPrivKey.PubKey(), newPostPrivKey.PubKey())
	assert.Nil(t, err)
	accInfo := model.AccountInfo{
		Username:        user1,
		CreatedAt:       ctx.BlockHeader().Time,
		MasterKey:       newMasterPrivKey.PubKey(),
		TransactionKey:  newTransactionPrivKey.PubKey(),
		MicropaymentKey: newMicropaymentPrivKey.PubKey(),
		PostKey:         newPostPrivKey.PubKey(),
	}
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee,
		Stake:   accParam.RegisterFee,
		NumOfTx: 1,
	}

	checkAccountInfo(t, ctx, user1, accInfo)
	checkBankKVByUsername(t, ctx, user1, bank)

	pendingStakeQueue := model.PendingStakeQueue{
		StakeCoinInQueue: sdk.ZeroRat(),
	}
	checkPendingStake(t, ctx, user1, pendingStakeQueue)
	stake, err := am.GetStake(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, accParam.RegisterFee, stake)
	ctx = ctx.WithBlockHeader(
		abci.Header{
			ChainID: "Lino", Height: 1,
			Time: ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake})
	stake, err = am.GetStake(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, accParam.RegisterFee, stake)
}

func TestIncreaseSequenceByOne(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	createTestAccount(ctx, am, string(user1))

	cases := []struct {
		user           types.AccountKey
		increaseTimes  int
		expectSequence int64
	}{
		{user1, 1, 1},
		{user1, 100, 101},
	}

	for _, cs := range cases {

		for i := 0; i < cs.increaseTimes; i++ {
			am.IncreaseSequenceByOne(ctx, user1)
		}
		seq, err := am.GetSequence(ctx, user1)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectSequence, seq)
	}
}

func TestAddFrozenMoney(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	createTestAccount(ctx, am, string(user1))

	testCases := []struct {
		frozenAmount            types.Coin
		startAt                 int64
		interval                int64
		times                   int64
		expectNumOfFrozenAmount int
	}{
		{types.NewCoinFromInt64(100), 1000000, 10, 5, 1},
		{types.NewCoinFromInt64(100), 1200000, 10, 5, 1},
		{types.NewCoinFromInt64(100), 1300000, 10, 5, 2},
		{types.NewCoinFromInt64(100), 1400000, 10, 5, 2},
		{types.NewCoinFromInt64(100), 1600000, 10, 5, 1}, // this one is used to re-produce the out-of-bound bug.
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: tc.startAt})
		err := am.AddFrozenMoney(ctx, user1, tc.frozenAmount, tc.startAt, tc.interval, tc.times)
		assert.Nil(t, err)

		accountBank, err := am.storage.GetBankFromAccountKey(ctx, user1)
		assert.Nil(t, err)
		assert.Equal(t, tc.expectNumOfFrozenAmount, len(accountBank.FrozenMoneyList))
	}
}
