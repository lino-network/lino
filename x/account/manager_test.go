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

func checkAccountGrantKeyList(
	t *testing.T, ctx sdk.Context, accKey types.AccountKey, grantList model.GrantKeyList) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	grantListPtr, err := accStorage.GetGrantKeyList(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, grantList, *grantListPtr, "accout grantList should be equal")
}

func TestIsAccountExist(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	assert.False(t, am.IsAccountExist(ctx, types.AccountKey("user1")))
	createTestAccount(ctx, am, "user1")
	assert.True(t, am.IsAccountExist(ctx, types.AccountKey("user1")))
}

func TestAddCoin(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	fromUser1, fromUser2 := types.AccountKey("fromUser1"), types.AccountKey("fromuser2")
	testUser := "testUser"

	baseTime := time.Now().Unix()
	baseTime1 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake/2
	baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake + 1
	baseTime3 := baseTime2 + coinDayParams.SecondsToRecoverCoinDayStake + 1
	ctx = ctx.WithBlockHeader(abci.Header{Time: baseTime})
	createTestAccount(ctx, am, testUser)
	cases := []struct {
		testName                 string
		Amount                   types.Coin
		From                     string
		DetailType               types.TransferDetailType
		Memo                     string
		AtWhen                   int64
		ExpectBank               model.AccountBank
		ExpectPendingStakeQueue  model.PendingStakeQueue
		ExpectBalanceHistorySlot model.BalanceHistory
	}{
		{"add coin to account's saving",
			c100, string(fromUser1), types.TransferIn, "memo", baseTime,
			model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c100),
				NumOfTx: 2,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat,
				TotalCoin:        accParam.RegisterFee.Plus(c100),
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      accParam.RegisterFee,
					},
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
						From:       string(accountReferrer),
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "init account",
					},
					model.Detail{
						Amount:     c100,
						From:       string(fromUser1),
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
				},
			},
		},
		{"add coin to exist account's saving while previous tx is still in pending queue", c100,
			string(fromUser2), types.DonationIn, "permlink", baseTime1,
			model.AccountBank{
				Saving:  accParam.RegisterFee.Plus(c200),
				NumOfTx: 3,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime1,
				StakeCoinInQueue: sdk.NewRat(5050000, 1),
				TotalCoin:        accParam.RegisterFee.Plus(c100).Plus(c100),
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      accParam.RegisterFee,
					},
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
						From:       string(accountReferrer),
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "init account",
					},
					model.Detail{
						Amount:     c100,
						From:       string(fromUser1),
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
					model.Detail{
						Amount:     c100,
						From:       string(fromUser2),
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
						From:       string(accountReferrer),
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "init account",
					},
					model.Detail{
						Amount:     c100,
						From:       string(fromUser1),
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
					model.Detail{
						Amount:     c100,
						From:       string(fromUser2),
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
				Stake:   accParam.RegisterFee.Plus(c300),
				NumOfTx: 5,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime3,
				StakeCoinInQueue: sdk.ZeroRat,
				TotalCoin:        c0,
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime3,
						EndTime:   baseTime3 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c0,
					},
				},
			},
			model.BalanceHistory{
				[]model.Detail{

					model.Detail{
						Amount:     accParam.RegisterFee,
						From:       string(accountReferrer),
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "init account",
					},
					model.Detail{
						Amount:     c100,
						From:       string(fromUser1),
						To:         testUser,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "memo",
					},
					model.Detail{
						Amount:     c100,
						From:       string(fromUser2),
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
					model.Detail{
						Amount:     c0,
						From:       "",
						To:         testUser,
						CreatedAt:  baseTime3,
						DetailType: types.DelegationReturnCoin,
					},
				},
			},
		},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		err = am.AddSavingCoin(
			ctx, types.AccountKey(testUser), cs.Amount, cs.From, cs.Memo, cs.DetailType)

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
		ctx, userWithSufficientSaving, accParam.RegisterFee, string(fromUser), "", types.TransferIn)
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
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat,
				TotalCoin:        accParam.RegisterFee.Plus(accParam.RegisterFee).Minus(coin1),
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      accParam.RegisterFee,
					},
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
						From:       string(accountReferrer),
						To:         string(userWithSufficientSaving),
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "init account",
					},
					model.Detail{
						Amount:     accParam.RegisterFee,
						From:       string(fromUser),
						To:         string(userWithSufficientSaving),
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					model.Detail{
						Amount:     coin1,
						From:       string(userWithSufficientSaving),
						To:         string(toUser),
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
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat,
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
						From:       string(accountReferrer),
						To:         string(toUser),
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "init account",
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
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat,
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
						From:       string(accountReferrer),
						To:         string(userWithLimitSaving),
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
						Memo:       "init account",
					},
				},
			},
		},
	}
	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		err = am.MinusSavingCoin(ctx, cs.FromUser, cs.Amount, string(cs.To), cs.Memo, cs.DetailType)

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
	fromUser, toUser := "fromUser", "toUser"
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
	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	// normal test
	assert.False(t, am.IsAccountExist(ctx, accKey))
	err = am.CreateAccount(ctx, accountReferrer, accKey,
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), accParam.RegisterFee)
	assert.Nil(t, err)

	assert.True(t, am.IsAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee,
		NumOfTx: 1,
	}
	checkBankKVByUsername(t, ctx, accKey, bank)
	pendingStakeQueue := model.PendingStakeQueue{
		LastUpdatedAt:    ctx.BlockHeader().Time,
		StakeCoinInQueue: sdk.ZeroRat,
		TotalCoin:        accParam.RegisterFee,
		PendingStakeList: []model.PendingStake{model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake,
			Coin:      accParam.RegisterFee,
		}}}
	checkPendingStake(t, ctx, accKey, pendingStakeQueue)
	accInfo := model.AccountInfo{
		Username:       accKey,
		CreatedAt:      ctx.BlockHeader().Time,
		MasterKey:      priv.PubKey(),
		TransactionKey: priv.Generate(1).PubKey(),
		PostKey:        priv.Generate(2).PubKey(),
	}
	checkAccountInfo(t, ctx, accKey, accInfo)
	accMeta := model.AccountMeta{
		LastActivityAt: ctx.BlockHeader().Time,
	}
	checkAccountMeta(t, ctx, accKey, accMeta)

	reward := model.Reward{coin0, coin0, coin0, coin0}
	checkAccountReward(t, ctx, accKey, reward)

	var grantPubKeyList []model.GrantPubKey
	grantList := model.GrantKeyList{GrantPubKeyList: grantPubKeyList}
	checkAccountGrantKeyList(t, ctx, accKey, grantList)

	balanceHistory, err := am.storage.GetBalanceHistory(ctx, accKey, 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(balanceHistory.Details))
	assert.Equal(t, model.Detail{
		From:       string(accountReferrer),
		To:         string(accKey),
		Amount:     accParam.RegisterFee,
		CreatedAt:  ctx.BlockHeader().Time,
		DetailType: types.TransferIn,
		Memo:       "init account",
	}, balanceHistory.Details[0])
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
			crypto.GenPrivKeyEd25519().PubKey(),
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
	expectNumOfTx := int64(1)

	createTestAccount(ctx, am, string(accKey))

	cases := []struct {
		testName            string
		IsAdd               bool
		Coin                types.Coin
		AtWhen              int64
		ExpectSavingBalance types.Coin
		ExpectStake         types.Coin
		ExpectStakeInBank   types.Coin
	}{
		{"add coin before charging first coin",
			true, accParam.RegisterFee, baseTime + (totalCoinDaysSec/registerFee)/2,
			doubleRegisterFee, coin0, coin0},
		{"check first coin",
			true, coin0, baseTime + (totalCoinDaysSec/registerFee)/2 + 1,
			doubleRegisterFee, coin1, coin0},
		{"check both transactions fully charged",
			true, coin0, baseTime2, doubleRegisterFee, doubleRegisterFee, doubleRegisterFee},
		{"withdraw half deposit",
			false, accParam.RegisterFee, baseTime2,
			accParam.RegisterFee, accParam.RegisterFee, accParam.RegisterFee},
		{"charge again",
			true, accParam.RegisterFee, baseTime2,
			doubleRegisterFee, accParam.RegisterFee, accParam.RegisterFee},
		{"withdraw half deposit while the last transaction is still charging",
			false, halfRegisterFee, baseTime2 + totalCoinDaysSec/2 + 1,
			accParam.RegisterFee.Plus(halfRegisterFee),
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(registerFee / 4)), accParam.RegisterFee},
		{"withdraw last transaction which is still charging",
			false, halfRegisterFee, baseTime2 + totalCoinDaysSec/2 + 1,
			accParam.RegisterFee, accParam.RegisterFee, accParam.RegisterFee},
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

		expectNumOfTx++

		bank := model.AccountBank{
			Saving:  cs.ExpectSavingBalance,
			Stake:   cs.ExpectStakeInBank,
			NumOfTx: expectNumOfTx,
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
		Stake:   c0,
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
	user2 := types.AccountKey("user2")
	user3 := types.AccountKey("user3")

	masterKey := crypto.GenPrivKeyEd25519()
	transactionKey := crypto.GenPrivKeyEd25519()
	postKey := crypto.GenPrivKeyEd25519()
	am.CreateAccount(
		ctx, accountReferrer, user1, masterKey.PubKey(), transactionKey.PubKey(),
		postKey.PubKey(), accParam.RegisterFee)

	priv2 := createTestAccount(ctx, am, string(user2))
	priv3 := createTestAccount(ctx, am, string(user3))
	err := am.AuthorizePermission(ctx, user1, user2, 100, types.PostPermission)
	assert.Nil(t, err)

	baseTime := ctx.BlockHeader().Time

	cases := []struct {
		testName     string
		checkUser    types.AccountKey
		checkPubKey  crypto.PubKey
		atWhen       int64
		grantLevel   types.Permission
		expectUser   types.AccountKey
		expectResult sdk.Error
	}{
		{"check user's master key",
			user1, masterKey.PubKey(), baseTime, types.MasterPermission, user1, nil},
		{"check user's transaction key",
			user1, transactionKey.PubKey(), baseTime, types.TransactionPermission, user1, nil},
		{"check user's post key",
			user1, postKey.PubKey(), baseTime, types.PostPermission, user1, nil},
		{"user's transaction key can authorize post permission",
			user1, transactionKey.PubKey(), baseTime, types.PostPermission, user1, nil},
		{"check user's transaction key can't authorize master permission",
			user1, transactionKey.PubKey(), baseTime, types.MasterPermission, user1,
			ErrCheckMasterKey()},
		{"check user's post key can't authorize master permission",
			user1, postKey.PubKey(), baseTime, types.MasterPermission, user1,
			ErrCheckMasterKey()},
		{"check user's post key can't authorize transaction permission",
			user1, postKey.PubKey(), baseTime, types.TransactionPermission, user1,
			ErrCheckTransactionKey()},
		{"check user2's pubkey",
			user1, priv2.Generate(2).PubKey(), baseTime, types.PostPermission, user2, nil},
		{"check unauthorized user pubkey",
			user1, priv3.Generate(2).PubKey(), baseTime, types.PostPermission, "",
			ErrCheckAuthenticatePubKeyOwner(user1)},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: cs.atWhen})
		grantUser, err := am.CheckAuthenticatePubKeyOwner(ctx, cs.checkUser, cs.checkPubKey, cs.grantLevel)
		assert.Equal(t, cs.expectResult, err)
		if cs.expectResult == nil {
			if cs.expectUser != grantUser {
				t.Errorf(
					"%s: expect key owner incorrect, expect %v, got %v",
					cs.testName, cs.expectUser, grantUser)
				return
			}
		}
	}
}

func TestGrantPubkey(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")
	user3 := types.AccountKey("user3")

	createTestAccount(ctx, am, string(user1))
	priv2 := createTestAccount(ctx, am, string(user2))
	priv3 := createTestAccount(ctx, am, string(user3))

	baseTime := ctx.BlockHeader().Time

	cases := []struct {
		user             types.AccountKey
		grantTo          types.AccountKey
		expireTime       int64
		checkTime        int64
		checkGrantUser   types.AccountKey
		checkGrantPubKey crypto.PubKey
		expectResult     sdk.Error
	}{
		{user1, user2, 100, baseTime + 99, user2, priv2.Generate(2).PubKey(), nil},
		{user1, user3, 100, baseTime + 99, user3, priv3.Generate(2).PubKey(), nil},
		{user1, user2, 100, baseTime + 101, user2, priv2.Generate(2).PubKey(),
			ErrCheckAuthenticatePubKeyOwner(user1)},
		{user1, user2, 100, baseTime + 99, user2, priv2.Generate(2).PubKey(), nil},
		{user1, user2, 500, baseTime + 101, user2, priv2.Generate(2).PubKey(), nil},
		{user1, user2, 300, baseTime + 301, user2, priv2.Generate(2).PubKey(),
			ErrCheckAuthenticatePubKeyOwner(user1)},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: baseTime})
		err := am.AuthorizePermission(ctx, cs.user, cs.grantTo, cs.expireTime, types.PostPermission)
		assert.Nil(t, err)
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.checkTime})
		grantUser, err := am.CheckAuthenticatePubKeyOwner(ctx, cs.user, cs.checkGrantPubKey, 0)
		assert.Equal(t, err, cs.expectResult)
		if cs.expectResult == nil {
			assert.Equal(t, grantUser, cs.checkGrantUser)
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
	newTransactionPrivKey := newMasterPrivKey.Generate(1)
	newPostPrivKey := newMasterPrivKey.Generate(2)

	err = am.RecoverAccount(ctx, user1,
		newMasterPrivKey.PubKey(), newTransactionPrivKey.PubKey(), newPostPrivKey.PubKey())
	assert.Nil(t, err)
	accInfo := model.AccountInfo{
		Username:       user1,
		CreatedAt:      ctx.BlockHeader().Time,
		MasterKey:      newMasterPrivKey.PubKey(),
		TransactionKey: newTransactionPrivKey.PubKey(),
		PostKey:        newPostPrivKey.PubKey(),
	}
	bank := model.AccountBank{
		Saving:  accParam.RegisterFee,
		Stake:   coin0,
		NumOfTx: 1,
	}

	checkAccountInfo(t, ctx, user1, accInfo)
	checkBankKVByUsername(t, ctx, user1, bank)

	pendingStakeQueue := model.PendingStakeQueue{
		LastUpdatedAt:    ctx.BlockHeader().Time,
		StakeCoinInQueue: sdk.ZeroRat,
		TotalCoin:        accParam.RegisterFee,
		PendingStakeList: []model.PendingStake{
			model.PendingStake{
				StartTime: ctx.BlockHeader().Time,
				EndTime:   ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake,
				Coin:      accParam.RegisterFee,
			}},
	}
	checkPendingStake(t, ctx, user1, pendingStakeQueue)
	stake, err := am.GetStake(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, coin0, stake)
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

	cases := []struct {
		frozenAmount            types.Coin
		startAt                 int64
		interval                int64
		times                   int64
		expectNumOfFrozenAmount int
	}{
		{types.NewCoinFromInt64(100), 10000, 10, 5, 1},
		{types.NewCoinFromInt64(100), 10100, 10, 5, 1},
		{types.NewCoinFromInt64(100), 10110, 10, 5, 2},
		{types.NewCoinFromInt64(100), 10151, 10, 5, 2},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: cs.startAt})
		err := am.AddFrozenMoney(ctx, user1, cs.frozenAmount, cs.startAt, cs.interval, cs.times)
		assert.Nil(t, err)

		accountBank, err := am.storage.GetBankFromAccountKey(ctx, user1)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectNumOfFrozenAmount, len(accountBank.FrozenMoneyList))
	}
}
