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

func checkBankKVByAddress(
	t *testing.T, ctx sdk.Context, addr sdk.Address, bank model.AccountBank) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	bankPtr, err := accStorage.GetBankFromAddress(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, bank, *bankPtr, "bank should be equal")
}

func checkPendingStake(
	t *testing.T, ctx sdk.Context, addr sdk.Address, pendingStakeQueue model.PendingStakeQueue) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	pendingStakeQueuePtr, err := accStorage.GetPendingStakeQueue(ctx, addr)
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

func TestAddCoinFromAddress(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	baseTime := time.Now().Unix()
	baseTime1 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake/2
	baseTime2 := baseTime + coinDayParams.SecondsToRecoverCoinDayStake + 1

	cases := []struct {
		testName                string
		AddToSaving             bool
		ToAddress               sdk.Address
		Amount                  types.Coin
		AtWhen                  int64
		ExpectBank              model.AccountBank
		ExpectPendingStakeQueue model.PendingStakeQueue
	}{
		{"add coin to non-exist account's saving",
			true, sdk.Address("test"), coin1, baseTime,
			model.AccountBank{
				Address: sdk.Address("test"),
				Saving:  coin1,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat,
				TotalCoin:        coin1,
				PendingStakeList: []model.PendingStake{model.PendingStake{
					StartTime: baseTime,
					EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
					Coin:      coin1,
				}},
			},
		},
		{"add coin to non-exist account's checking",
			false, sdk.Address("test1"), coin1, baseTime,
			model.AccountBank{
				Address:  sdk.Address("test1"),
				Checking: coin1,
			},
			model.PendingStakeQueue{},
		},
		{"add coin to exist account's saving",
			true, sdk.Address("test1"), coin1, baseTime,
			model.AccountBank{
				Address:  sdk.Address("test1"),
				Saving:   coin1,
				Checking: coin1,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat,
				TotalCoin:        coin1,
				PendingStakeList: []model.PendingStake{model.PendingStake{
					StartTime: baseTime,
					EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
					Coin:      coin1,
				}},
			},
		},
		{"add coin to exist account's checking",
			false, sdk.Address("test"), coin1, baseTime,
			model.AccountBank{
				Address:  sdk.Address("test"),
				Saving:   coin1,
				Checking: coin1,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat,
				TotalCoin:        coin1,
				PendingStakeList: []model.PendingStake{model.PendingStake{
					StartTime: baseTime,
					EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
					Coin:      coin1,
				}},
			},
		},
		{"add coin to exist account's saving while previous tx is still in pending queue",
			true, sdk.Address("test"), coin1, baseTime1,
			model.AccountBank{
				Address:  sdk.Address("test"),
				Saving:   coin2,
				Checking: coin1,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime1,
				StakeCoinInQueue: sdk.NewRat(1, 2),
				TotalCoin:        coin2,
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      coin1,
					},
					model.PendingStake{
						StartTime: baseTime1,
						EndTime:   baseTime1 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      coin1,
					}},
			},
		},
		{"add coin to exist account's checking while previous tx is still in pending queue",
			false, sdk.Address("test1"), coin1, baseTime1,
			model.AccountBank{
				Address:  sdk.Address("test1"),
				Saving:   coin1,
				Checking: coin2,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime,
				StakeCoinInQueue: sdk.ZeroRat,
				TotalCoin:        coin1,
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime,
						EndTime:   baseTime + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      coin1,
					},
				},
			},
		},
		{"add coin to exist account's saving while previous tx just finished pending",
			true, sdk.Address("test1"), coin1, baseTime2,
			model.AccountBank{
				Address:  sdk.Address("test1"),
				Saving:   coin2,
				Checking: coin2,
				Stake:    coin1,
			},
			model.PendingStakeQueue{
				LastUpdatedAt:    baseTime2,
				TotalCoin:        coin1,
				StakeCoinInQueue: sdk.ZeroRat,
				PendingStakeList: []model.PendingStake{
					model.PendingStake{
						StartTime: baseTime2,
						EndTime:   baseTime2 + coinDayParams.SecondsToRecoverCoinDayStake,
						Coin:      coin1,
					},
				},
			},
		},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		if cs.AddToSaving {
			err = am.AddSavingCoinToAddress(ctx, cs.ToAddress, cs.Amount)
		} else {
			err = am.AddCheckingCoinToAddress(ctx, cs.ToAddress, cs.Amount)
		}
		if err != nil {
			t.Errorf("%s: add coin failed, err: %v", cs.testName, err)
			return
		}
		fmt.Println(cs.testName, cs.ExpectBank)
		checkBankKVByAddress(t, ctx, cs.ToAddress, cs.ExpectBank)
		checkPendingStake(t, ctx, cs.ToAddress, cs.ExpectPendingStakeQueue)
	}
}

func TestMinusCoin(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	userWithSufficientSaving := types.AccountKey("user1")
	userWithSufficientChecking := types.AccountKey("user2")
	userWithLimitSaving := types.AccountKey("user3")

	priv1 := createTestAccount(ctx, am, string(userWithSufficientSaving))
	priv2 := createTestAccount(ctx, am, string(userWithSufficientChecking))
	priv3 := createTestAccount(ctx, am, string(userWithLimitSaving))
	err = am.AddSavingCoin(ctx, userWithSufficientSaving, accParam.RegisterFee)
	assert.Nil(t, err)
	err = am.AddCheckingCoin(ctx, userWithSufficientChecking, accParam.RegisterFee)
	assert.Nil(t, err)

	baseTime := ctx.BlockHeader().Time

	cases := []struct {
		TestName                string
		MinusFromSaving         bool
		FromUser                types.AccountKey
		UserPriv                crypto.PrivKey
		ExpectErr               sdk.Error
		Amount                  types.Coin
		AtWhen                  int64
		ExpectBank              model.AccountBank
		ExpectPendingStakeQueue model.PendingStakeQueue
	}{
		{"minus saving coin from user with sufficient saving",
			true, userWithSufficientSaving, priv1, nil, coin1, baseTime,
			model.AccountBank{
				Username: userWithSufficientSaving,
				Address:  priv1.PubKey().Address(),
				Saving:   accParam.RegisterFee.Plus(accParam.RegisterFee).Minus(coin1),
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
		},
		{"minus checking coin from user with sufficient checking",
			false, userWithSufficientChecking, priv2, nil, coin1, baseTime,
			model.AccountBank{
				Username: userWithSufficientChecking,
				Address:  priv2.PubKey().Address(),
				Saving:   accParam.RegisterFee,
				Checking: accParam.RegisterFee.Minus(coin1),
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
		},
		{"minus too much checking coin from user with sufficient checking",
			false, userWithSufficientChecking, priv2, ErrAccountCheckingCoinNotEnough(),
			accParam.RegisterFee, baseTime,
			model.AccountBank{
				Username: userWithSufficientChecking,
				Address:  priv2.PubKey().Address(),
				Saving:   accParam.RegisterFee,
				Checking: accParam.RegisterFee.Minus(coin1),
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
		},
		{"minus too much checking coin from user without checking",
			false, userWithLimitSaving, priv3, ErrAccountCheckingCoinNotEnough(),
			accParam.RegisterFee, baseTime,
			model.AccountBank{
				Username: userWithLimitSaving,
				Address:  priv3.PubKey().Address(),
				Saving:   accParam.RegisterFee,
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
		},
		{"minus saving coin from user with limit saving",
			true, userWithLimitSaving, priv3, ErrAccountSavingCoinNotEnough(),
			coin1, baseTime,
			model.AccountBank{
				Username: userWithLimitSaving,
				Address:  priv3.PubKey().Address(),
				Saving:   accParam.RegisterFee,
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
		},
	}
	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		if cs.MinusFromSaving {
			err = am.MinusSavingCoin(ctx, cs.FromUser, cs.Amount)
		} else {
			err = am.MinusCheckingCoin(ctx, cs.FromUser, cs.Amount)
		}

		assert.Equal(t, cs.ExpectErr, err, fmt.Sprintf("%s: minus coin failed, err: %v", cs.TestName, err))
		checkBankKVByAddress(t, ctx, cs.UserPriv.PubKey().Address(), cs.ExpectBank)
		checkPendingStake(t, ctx, cs.UserPriv.PubKey().Address(), cs.ExpectPendingStakeQueue)
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
	err = am.AddSavingCoinToAddress(ctx, priv.PubKey().Address(), accParam.RegisterFee)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey,
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	assert.Nil(t, err)

	assert.True(t, am.IsAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Address:  priv.PubKey().Address(),
		Saving:   accParam.RegisterFee,
		Username: accKey,
	}
	checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
	pendingStakeQueue := model.PendingStakeQueue{
		LastUpdatedAt:    ctx.BlockHeader().Time,
		StakeCoinInQueue: sdk.ZeroRat,
		TotalCoin:        accParam.RegisterFee,
		PendingStakeList: []model.PendingStake{model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake,
			Coin:      accParam.RegisterFee,
		}}}
	checkPendingStake(t, ctx, priv.PubKey().Address(), pendingStakeQueue)
	accInfo := model.AccountInfo{
		Username:       accKey,
		CreatedAt:      ctx.BlockHeader().Time,
		MasterKey:      priv.PubKey(),
		TransactionKey: priv.Generate(1).PubKey(),
		PostKey:        priv.Generate(2).PubKey(),
		Address:        priv.PubKey().Address(),
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
	enoughSavingCoinPriv := crypto.GenPrivKeyEd25519()
	enoughCheckingCoinPriv := crypto.GenPrivKeyEd25519()
	insufficientSavingCoinPriv := crypto.GenPrivKeyEd25519()
	insufficientCheckingCoinPriv := crypto.GenPrivKeyEd25519()
	newPrivWithoutBank := crypto.GenPrivKeyEd25519()

	accKey1 := types.AccountKey("accKey1")
	accKey2 := types.AccountKey("accKey2")

	err := am.AddSavingCoinToAddress(
		ctx, enoughSavingCoinPriv.PubKey().Address(), accParam.RegisterFee)
	err = am.AddSavingCoinToAddress(
		ctx, insufficientSavingCoinPriv.PubKey().Address(), types.NewCoin(1))
	err = am.AddCheckingCoinToAddress(
		ctx, enoughCheckingCoinPriv.PubKey().Address(), accParam.RegisterFee)
	err = am.AddSavingCoinToAddress(
		ctx, insufficientCheckingCoinPriv.PubKey().Address(), types.NewCoin(1))

	cases := []struct {
		testName  string
		username  types.AccountKey
		privkey   crypto.PrivKey
		expectErr sdk.Error
	}{
		{"register user with sufficient saving coin",
			accKey1, enoughSavingCoinPriv, nil,
		},
		{"username already took",
			accKey1, enoughSavingCoinPriv, ErrAccountAlreadyExists(accKey1),
		},
		{"bank already registered",
			accKey2, enoughSavingCoinPriv, ErrBankAlreadyRegistered(),
		},
		{"bank doesn't have coin",
			accKey2, newPrivWithoutBank, ErrAccountCreateFailed(accKey2),
		},
		{"bank has insufficient saving",
			accKey2, insufficientSavingCoinPriv, ErrRegisterFeeInsufficient(),
		},
		{"bank has sufficient checking",
			accKey2, insufficientSavingCoinPriv, ErrRegisterFeeInsufficient(),
		},
		{"bank has insufficient checking",
			accKey2, insufficientCheckingCoinPriv, ErrRegisterFeeInsufficient(),
		},
	}
	for _, cs := range cases {
		err = am.CreateAccount(
			ctx, cs.username, cs.privkey.PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey())
		assert.Equal(t, cs.expectErr, err,
			fmt.Sprintf("%s: create account failed: expect %v, got %v",
				cs.testName, cs.expectErr, err))
	}
}

func TestCoinDayByAddress(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)
	assert.Nil(t, err)
	totalCoinDaysSec := coinDayParams.SecondsToRecoverCoinDayStake
	registerFee := accParam.RegisterFee.ToInt64()
	doubleRegisterFee := types.NewCoin(registerFee * 2)
	halfRegisterFee := types.NewCoin(registerFee / 2)

	// create bank and account
	priv := createTestAccount(ctx, am, string(accKey))

	baseTime1 := ctx.BlockHeader().Time
	baseTime2 := baseTime1 + totalCoinDaysSec*2
	testCases := []struct {
		testName          string
		AddCoin           types.Coin
		AtWhen            int64
		ExpectBalance     types.Coin
		ExpectStake       types.Coin
		ExpectStakeInBank types.Coin
	}{
		{"before charge first coin",
			coin0, baseTime1 + (totalCoinDaysSec/registerFee)/2,
			accParam.RegisterFee, coin0, coin0},
		{"after charge first coin",
			coin0, baseTime1 + (totalCoinDaysSec/registerFee)/2 + 1,
			accParam.RegisterFee, coin1, coin0},
		{"charge half coin",
			coin0, baseTime1 + totalCoinDaysSec/2, accParam.RegisterFee,
			halfRegisterFee, coin0},
		{"transfer new coin",
			accParam.RegisterFee, baseTime1 + totalCoinDaysSec/2,
			doubleRegisterFee, halfRegisterFee, coin0},
		{"first transaction charge finished",
			coin0, baseTime1 + totalCoinDaysSec + 1, doubleRegisterFee,
			accParam.RegisterFee.Plus(halfRegisterFee), accParam.RegisterFee},
		{"all transaction charge finished",
			coin0, baseTime1 + totalCoinDaysSec*2 + 1,
			doubleRegisterFee, doubleRegisterFee, doubleRegisterFee},
		{"transaction with only one coin",
			coin1, baseTime2, types.NewCoin(registerFee*2 + 1), doubleRegisterFee,
			doubleRegisterFee},
		{"transaction with one coin charge ongoing",
			coin0, baseTime2 + totalCoinDaysSec/2, types.NewCoin(registerFee*2 + 1),
			doubleRegisterFee, doubleRegisterFee},
		{"transaction with one coin charge finished",
			coin0, baseTime2 + totalCoinDaysSec/2 + 1,
			types.NewCoin(registerFee*2 + 1), types.NewCoin(registerFee*2 + 1), doubleRegisterFee},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: tc.AtWhen})
		err := am.AddSavingCoinToAddress(ctx, priv.PubKey().Address(), tc.AddCoin)
		if err != nil {
			t.Errorf("%s: add coin failed, expect %v, got %v", tc.testName, "nil", err)
			return
		}
		coin, err := am.GetStake(ctx, accKey)
		if err != nil {
			t.Errorf("%s: get stake failed, expect %v, got %v", tc.testName, "nil", err)
			return
		}
		if !tc.ExpectStake.IsEqual(coin) {
			t.Errorf("%s: expect stake incorrect, expect %v, got %v", tc.testName, tc.ExpectStake, coin)
			return
		}
		bank := model.AccountBank{
			Address:  priv.PubKey().Address(),
			Saving:   tc.ExpectBalance,
			Stake:    tc.ExpectStakeInBank,
			Username: accKey,
		}
		checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
	}
}

func TestCoinDayByAccountKey(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)
	totalCoinDaysSec := coinDayParams.SecondsToRecoverCoinDayStake
	registerFee := accParam.RegisterFee.ToInt64()
	doubleRegisterFee := types.NewCoin(registerFee * 2)
	halfRegisterFee := types.NewCoin(registerFee / 2)

	baseTime := ctx.BlockHeader().Time
	baseTime2 := baseTime + totalCoinDaysSec + (totalCoinDaysSec/registerFee)/2 + 1

	priv := createTestAccount(ctx, am, string(accKey))

	cases := []struct {
		testName          string
		IsAdd             bool
		Coin              types.Coin
		AtWhen            int64
		ExpectBalance     types.Coin
		ExpectStake       types.Coin
		ExpectStakeInBank types.Coin
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
			accParam.RegisterFee.Plus(types.NewCoin(registerFee / 4)), accParam.RegisterFee},
		{"withdraw last transaction which is still charging",
			false, halfRegisterFee, baseTime2 + totalCoinDaysSec/2 + 1,
			accParam.RegisterFee, accParam.RegisterFee, accParam.RegisterFee},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		if cs.IsAdd {
			err := am.AddSavingCoinToAddress(ctx, priv.PubKey().Address(), cs.Coin)
			assert.Nil(t, err)
		} else {
			err := am.MinusSavingCoin(ctx, accKey, cs.Coin)
			assert.Nil(t, err)
		}
		coin, err := am.GetStake(ctx, accKey)
		assert.Nil(t, err)
		if !cs.ExpectStake.IsEqual(coin) {
			t.Errorf("%s: expect stake incorrect, expect %v, got %v", cs.testName, cs.ExpectStake, coin)
			return
		}

		bank := model.AccountBank{
			Address:  priv.PubKey().Address(),
			Saving:   cs.ExpectBalance,
			Stake:    cs.ExpectStakeInBank,
			Username: accKey,
		}
		checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
	}
}

func TestAccountReward(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	accKey := types.AccountKey("accKey")
	priv := crypto.GenPrivKeyEd25519()

	err := am.AddSavingCoinToAddress(ctx, priv.PubKey().Address(), accParam.RegisterFee)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey,
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	assert.Nil(t, err)

	err = am.AddIncomeAndReward(ctx, accKey, c500, c200, c300)
	assert.Nil(t, err)
	reward := model.Reward{c500, c200, c300, c300}
	checkAccountReward(t, ctx, accKey, reward)
	err = am.AddIncomeAndReward(ctx, accKey, c500, c300, c200)
	assert.Nil(t, err)
	reward = model.Reward{c1000, c500, c500, c500}
	checkAccountReward(t, ctx, accKey, reward)

	bank := model.AccountBank{
		Address:  priv.PubKey().Address(),
		Saving:   accParam.RegisterFee,
		Stake:    c0,
		Username: accKey,
	}
	checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)

	err = am.ClaimReward(ctx, accKey)
	assert.Nil(t, err)
	bank.Saving = accParam.RegisterFee.Plus(c500)
	checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
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

	priv := createTestAccount(ctx, am, string(accKey))
	err = am.AddSavingCoinToAddress(ctx, priv.PubKey().Address(), c100)
	assert.Nil(t, err)

	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	err = accStorage.SetPendingStakeQueue(
		ctx, priv.PubKey().Address(), &model.PendingStakeQueue{})
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
		{sdk.NewRat(1, 10), types.NewCoin(10 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime, ErrAccountTPSCapacityNotEnough(accKey), types.NewCoin(0)},
		{sdk.NewRat(1, 10), types.NewCoin(10 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth, nil, types.NewCoin(990000)},
		{sdk.NewRat(1, 2), types.NewCoin(10 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth, nil, types.NewCoin(950000)},
		{sdk.NewRat(1), types.NewCoin(10 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth, nil, types.NewCoin(9 * types.Decimals)},
		{sdk.NewRat(1), types.NewCoin(1 * types.Decimals), baseTime,
			types.NewCoin(10 * types.Decimals), baseTime, nil, types.NewCoin(0)},
		{sdk.NewRat(1), types.NewCoin(10), baseTime, types.NewCoin(1 * types.Decimals),
			baseTime, ErrAccountTPSCapacityNotEnough(accKey), types.NewCoin(1 * types.Decimals)},
		{sdk.NewRat(1), types.NewCoin(1 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth/2,
			ErrAccountTPSCapacityNotEnough(accKey), types.NewCoin(0)},
		{sdk.NewRat(1, 2), types.NewCoin(1 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth/2, nil, types.NewCoin(0)},
		{sdk.OneRat, types.NewCoin(1 * types.Decimals), 0, types.NewCoin(0),
			baseTime, nil, types.NewCoin(0)},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: cs.CurrentTime})
		bank := &model.AccountBank{
			Address: priv.PubKey().Address(),
			Saving:  cs.UserStake,
			Stake:   cs.UserStake,
		}
		err = accStorage.SetBankFromAddress(ctx, priv.PubKey().Address(), bank)
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
		err := am.AuthorizePermission(ctx, cs.user, cs.grantTo, cs.expireTime, 0)
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

func TestAccountRecover(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	priv := createTestAccount(ctx, am, string(user1))

	cases := []struct {
		user              types.AccountKey
		newPostKey        crypto.PubKey
		newTransactionKey crypto.PubKey
	}{
		{user1, crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey()},
	}

	for _, cs := range cases {
		err := am.RecoverAccount(ctx, cs.user, cs.newPostKey, cs.newTransactionKey)
		assert.Nil(t, err)
		accInfo := model.AccountInfo{
			Username:       cs.user,
			CreatedAt:      ctx.BlockHeader().Time,
			MasterKey:      priv.PubKey(),
			TransactionKey: cs.newTransactionKey,
			PostKey:        cs.newPostKey,
			Address:        priv.PubKey().Address(),
		}
		checkAccountInfo(t, ctx, cs.user, accInfo)
	}
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
		{types.NewCoin(100), 10000, 10, 5, 1},
		{types.NewCoin(100), 10100, 10, 5, 1},
		{types.NewCoin(100), 10110, 10, 5, 2},
		{types.NewCoin(100), 10151, 10, 5, 2},
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
