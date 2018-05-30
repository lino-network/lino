package account

import (
	"fmt"
	"testing"
	"time"

	"github.com/lino-network/lino/tx/account/model"
	"github.com/lino-network/lino/types"

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

	// Get the minimum time of this history slot
	baseTime := time.Now().Unix()
	baseTime = baseTime / accParam.BalanceHistoryIntervalTime * accParam.BalanceHistoryIntervalTime
	baseTime1 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake/2
	baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake + 1
	baseTime3 := baseTime + accParam.BalanceHistoryIntervalTime + 1
	ctx = ctx.WithBlockHeader(abci.Header{Time: baseTime})
	priv1 := createTestAccount(ctx, am, "user1")
	cases := []struct {
		testName                 string
		ToUser                   types.AccountKey
		UserAddress              sdk.Address
		Amount                   types.Coin
		CoinType                 types.BalanceHistoryDetailType
		AtWhen                   int64
		ExpectBank               model.AccountBank
		ExpectPendingStakeQueue  model.PendingStakeQueue
		ExpectBalanceHistorySlot model.BalanceHistory
	}{
		{"add coin to account's saving",
			types.AccountKey("user1"), priv1.PubKey().Address(),
			c100, types.TransferIn, baseTime,
			model.AccountBank{
				Saving: accParam.RegisterFee.Plus(c100),
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
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					model.Detail{
						Amount:     c100,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
				},
			},
		},
		{"add coin to exist account's saving while previous tx is still in pending queue",
			types.AccountKey("user1"), priv1.PubKey().Address(), c100, types.DonationIn, baseTime1,
			model.AccountBank{
				Saving: accParam.RegisterFee.Plus(c200),
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
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					model.Detail{
						Amount:     c100,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					model.Detail{
						Amount:     c100,
						CreatedAt:  baseTime1,
						DetailType: types.DonationIn,
					},
				},
			},
		},
		{"add coin to exist account's saving while previous tx just finished pending",
			types.AccountKey("user1"), priv1.PubKey().Address(), c100, types.ClaimReward, baseTime2,
			model.AccountBank{
				Saving: accParam.RegisterFee.Plus(c300),
				Stake:  accParam.RegisterFee.Plus(c100),
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
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					model.Detail{
						Amount:     c100,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					model.Detail{
						Amount:     c100,
						CreatedAt:  baseTime1,
						DetailType: types.DonationIn,
					},
					model.Detail{
						Amount:     c100,
						CreatedAt:  baseTime2,
						DetailType: types.ClaimReward,
					},
				},
			},
		},
		{"add coin at next balance history slot",
			types.AccountKey("user1"), priv1.PubKey().Address(), c100, types.ClaimReward, baseTime3,
			model.AccountBank{
				Saving: accParam.RegisterFee.Plus(c400),
				Stake:  accParam.RegisterFee.Plus(c300),
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime3,
				StakeCoinInQueue: sdk.ZeroRat,
				TotalCoin:        c100,
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime3,
						EndTime:   baseTime3 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      c100,
					},
				},
			},
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						Amount:     c100,
						CreatedAt:  baseTime3,
						DetailType: types.ClaimReward,
					},
				},
			},
		},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		err = am.AddSavingCoin(ctx, cs.ToUser, cs.Amount, cs.CoinType)

		if err != nil {
			t.Errorf("%s: add coin failed, err: %v", cs.testName, err)
			return
		}
		checkBankKVByUsername(t, ctx, cs.ToUser, cs.ExpectBank)
		checkPendingStake(t, ctx, cs.ToUser, cs.ExpectPendingStakeQueue)
		checkBalanceHistory(
			t, ctx, cs.ToUser, cs.AtWhen/accParam.BalanceHistoryIntervalTime, cs.ExpectBalanceHistorySlot)
	}
}

func TestMinusCoin(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	userWithSufficientSaving := types.AccountKey("user1")
	userWithLimitSaving := types.AccountKey("user3")

	// Get the minimum time of this history slot
	baseTime := time.Now().Unix()
	baseTime = baseTime / accParam.BalanceHistoryIntervalTime * accParam.BalanceHistoryIntervalTime
	baseTime1 := baseTime + accParam.BalanceHistoryIntervalTime + 1
	// baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake + 1
	// baseTime3 := baseTime + accParam.BalanceHistoryIntervalTime + 1

	ctx = ctx.WithBlockHeader(abci.Header{Time: baseTime})
	priv1 := createTestAccount(ctx, am, string(userWithSufficientSaving))
	priv3 := createTestAccount(ctx, am, string(userWithLimitSaving))
	err = am.AddSavingCoin(ctx, userWithSufficientSaving, accParam.RegisterFee, types.TransferIn)
	assert.Nil(t, err)

	cases := []struct {
		TestName                string
		FromUser                types.AccountKey
		UserPriv                crypto.PrivKey
		ExpectErr               sdk.Error
		Amount                  types.Coin
		AtWhen                  int64
		ExpectBank              model.AccountBank
		ExpectPendingStakeQueue model.PendingStakeQueue
		DetailType              types.BalanceHistoryDetailType
		ExpectBalanceHistory    model.BalanceHistory
	}{
		{"minus saving coin from user with sufficient saving",
			userWithSufficientSaving, priv1, nil, coin1, baseTime,
			model.AccountBank{
				Saving: accParam.RegisterFee.Plus(accParam.RegisterFee).Minus(coin1),
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
			}, types.TransferOut,
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						Amount:     accParam.RegisterFee,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					model.Detail{
						Amount:     accParam.RegisterFee,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
					model.Detail{
						Amount:     coin1,
						CreatedAt:  baseTime,
						DetailType: types.TransferOut,
					},
				},
			},
		},
		{"minus saving coin from user with sufficient saving at next slot",
			userWithSufficientSaving, priv1, nil, coin1, baseTime1,
			model.AccountBank{
				Saving: accParam.RegisterFee.Plus(accParam.RegisterFee).Minus(coin2),
				Stake:  accParam.RegisterFee.Plus(accParam.RegisterFee).Minus(coin2),
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime1,
				StakeCoinInQueue: sdk.ZeroRat,
				TotalCoin:        types.NewCoinFromInt64(0),
				PendingStakeList: nil,
			}, types.TransferOut,
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						Amount:     coin1,
						CreatedAt:  baseTime1,
						DetailType: types.TransferOut,
					},
				},
			},
		},
		{"minus saving coin from user with limit saving",
			userWithLimitSaving, priv3, ErrAccountSavingCoinNotEnough(),
			coin1, baseTime,
			model.AccountBank{
				Saving: accParam.RegisterFee,
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
			}, types.TransferOut,
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						Amount:     accParam.RegisterFee,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
				},
			},
		},
		{"minus saving coin exceeds the coin user hold",
			userWithLimitSaving, priv3, ErrAccountSavingCoinNotEnough(),
			c100, baseTime,
			model.AccountBank{
				Saving: accParam.RegisterFee,
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
			}, types.TransferOut,
			model.BalanceHistory{
				[]model.Detail{
					model.Detail{
						Amount:     accParam.RegisterFee,
						CreatedAt:  baseTime,
						DetailType: types.TransferIn,
					},
				},
			},
		},
	}
	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		err = am.MinusSavingCoin(ctx, cs.FromUser, cs.Amount, cs.DetailType)

		assert.Equal(t, cs.ExpectErr, err, fmt.Sprintf("%s: minus coin failed, err: %v", cs.TestName, err))
		checkBankKVByUsername(t, ctx, cs.FromUser, cs.ExpectBank)
		checkPendingStake(t, ctx, cs.FromUser, cs.ExpectPendingStakeQueue)
		checkBalanceHistory(t, ctx, cs.FromUser,
			cs.AtWhen/accParam.BalanceHistoryIntervalTime, cs.ExpectBalanceHistory)
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
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey,
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), accParam.RegisterFee)
	assert.Nil(t, err)

	assert.True(t, am.IsAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Saving: accParam.RegisterFee,
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
			ctx, cs.username, cs.privkey.PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), cs.registerFee)
		assert.Equal(t, cs.expectErr, err,
			fmt.Sprintf("%s: create account failed: expect %v, got %v",
				cs.testName, cs.expectErr, err))
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
		testName              string
		IsAdd                 bool
		Coin                  types.Coin
		AtWhen                int64
		ExpectSavingBalance   types.Coin
		ExpectCheckingBalance types.Coin
		ExpectStake           types.Coin
		ExpectStakeInBank     types.Coin
	}{
		{"add coin before charging first coin",
			true, accParam.RegisterFee, baseTime + (totalCoinDaysSec/registerFee)/2,
			doubleRegisterFee, coin0, coin0, coin0},
		{"check first coin",
			true, coin0, baseTime + (totalCoinDaysSec/registerFee)/2 + 1,
			doubleRegisterFee, coin0, coin1, coin0},
		{"check both transactions fully charged",
			true, coin0, baseTime2, doubleRegisterFee, coin0, doubleRegisterFee, doubleRegisterFee},
		{"withdraw half deposit",
			false, accParam.RegisterFee, baseTime2,
			accParam.RegisterFee, coin0, accParam.RegisterFee, accParam.RegisterFee},
		{"charge again",
			true, accParam.RegisterFee, baseTime2,
			doubleRegisterFee, coin0, accParam.RegisterFee, accParam.RegisterFee},
		{"withdraw half deposit while the last transaction is still charging",
			false, halfRegisterFee, baseTime2 + totalCoinDaysSec/2 + 1,
			accParam.RegisterFee.Plus(halfRegisterFee), coin0,
			accParam.RegisterFee.Plus(types.NewCoinFromInt64(registerFee / 4)), accParam.RegisterFee},
		{"withdraw last transaction which is still charging",
			false, halfRegisterFee, baseTime2 + totalCoinDaysSec/2 + 1,
			accParam.RegisterFee, coin0, accParam.RegisterFee, accParam.RegisterFee},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		if cs.IsAdd {
			err := am.AddSavingCoin(ctx, accKey, cs.Coin, types.TransferIn)
			assert.Nil(t, err)
		} else {
			err := am.MinusSavingCoin(ctx, accKey, cs.Coin, types.TransferOut)
			assert.Nil(t, err)
		}
		coin, err := am.GetStake(ctx, accKey)
		assert.Nil(t, err)
		if !cs.ExpectStake.IsEqual(coin) {
			t.Errorf("%s: expect stake incorrect, expect %v, got %v", cs.testName, cs.ExpectStake, coin)
			return
		}

		bank := model.AccountBank{
			Saving: cs.ExpectSavingBalance,
			Stake:  cs.ExpectStakeInBank,
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
		Saving: accParam.RegisterFee,
		Stake:  c0,
	}
	checkBankKVByUsername(t, ctx, accKey, bank)

	err = am.ClaimReward(ctx, accKey)
	assert.Nil(t, err)
	bank.Saving = accParam.RegisterFee.Plus(c500)
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
	err = am.AddSavingCoin(ctx, accKey, c100, types.TransferIn)
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
		ctx, user1, masterKey.PubKey(), transactionKey.PubKey(),
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
		Saving: accParam.RegisterFee,
		Stake:  coin0,
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
