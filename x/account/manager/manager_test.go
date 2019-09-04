package manager

import (
	"testing"
	"time"

	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/types"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"
	acctypes "github.com/lino-network/lino/x/account/types"
	global "github.com/lino-network/lino/x/global/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type AccountManagerTestSuite struct {
	testsuites.CtxTestSuite
	am AccountManager
	ph *param.ParamKeeper
	// deps
	global *global.GlobalKeeper

	// mock data
	userWithoutBalance model.AccountInfo

	userWithBalance       model.AccountInfo
	userWithBalanceSaving types.Coin

	unreg model.AccountInfo

	unregSaving types.Coin
}

func TestAccountManagerTestSuite(t *testing.T) {
	suite.Run(t, new(AccountManagerTestSuite))
}

func (suite *AccountManagerTestSuite) SetupTest() {
	testAccountKey := sdk.NewKVStoreKey("account")
	suite.SetupCtx(0, time.Unix(0, 0), testAccountKey)
	suite.ph = &param.ParamKeeper{}
	suite.global = &global.GlobalKeeper{}
	suite.am = NewAccountManager(testAccountKey, suite.ph, suite.global)

	// background
	suite.userWithoutBalance = model.AccountInfo{
		Username:       linotypes.AccountKey("userwithoutbalance"),
		SigningKey:     secp256k1.GenPrivKey().PubKey(),
		TransactionKey: secp256k1.GenPrivKey().PubKey(),
	}
	suite.userWithoutBalance.Address = sdk.AccAddress(suite.userWithoutBalance.TransactionKey.Address())

	suite.userWithBalance = model.AccountInfo{
		Username:       linotypes.AccountKey("userwithbalance"),
		SigningKey:     secp256k1.GenPrivKey().PubKey(),
		TransactionKey: secp256k1.GenPrivKey().PubKey(),
	}
	suite.userWithBalance.Address = sdk.AccAddress(suite.userWithBalance.TransactionKey.Address())

	suite.unreg = model.AccountInfo{
		Username:       linotypes.AccountKey("unreg"),
		SigningKey:     secp256k1.GenPrivKey().PubKey(),
		TransactionKey: secp256k1.GenPrivKey().PubKey(),
	}
	suite.unreg.Address = sdk.AccAddress(suite.unreg.TransactionKey.Address())

	suite.userWithBalanceSaving = types.NewCoinFromInt64(1000 * types.Decimals)
	suite.unregSaving = types.NewCoinFromInt64(1 * types.Decimals)

	suite.am.CreateAccount(suite.Ctx, suite.userWithoutBalance.Username, suite.userWithoutBalance.SigningKey, suite.userWithoutBalance.TransactionKey)

	suite.am.CreateAccount(suite.Ctx, suite.userWithBalance.Username, suite.userWithBalance.SigningKey, suite.userWithBalance.TransactionKey)
	suite.am.AddCoinToUsername(suite.Ctx, suite.userWithBalance.Username, suite.userWithBalanceSaving)

	suite.am.AddCoinToAddress(suite.Ctx, sdk.AccAddress(suite.unreg.TransactionKey.Address()), suite.unregSaving)

	suite.ph.On("GetAccountParam", mock.Anything).Return(&parammodel.AccountParam{
		RegisterFee:    types.NewCoinFromInt64(100 * types.Decimals),
		MinimumBalance: types.NewCoinFromInt64(0),
	}, nil).Maybe()

	// // reg accounts
	// for _, v := range []linotypes.AccountKey{suite.user1, suite.user2, suite.app1, suite.app2, suite.app3} {
	// 	suite.am.On("DoesAccountExist", mock.Anything, v).Return(true).Maybe()
	// }
	// // unreg accounts
	// for _, v := range []linotypes.AccountKey{suite.unreg} {
	// 	suite.am.On("DoesAccountExist", mock.Anything, v).Return(false).Maybe()
	// }

	// // reg dev
	// for _, v := range []linotypes.AccountKey{suite.app1, suite.app2, suite.app3} {
	// 	suite.dev.On("DoesDeveloperExist", mock.Anything, v).Return(true).Maybe()
	// }
	// // unreg devs
	// for _, v := range []linotypes.AccountKey{suite.unreg, suite.user1, suite.user2} {
	// 	suite.dev.On("DoesDeveloperExist", mock.Anything, v).Return(false).Maybe()
	// }

	// rate, err := sdk.NewDecFromStr("0.099")
	// suite.Require().Nil(err)
	// suite.global.On("GetConsumptionFrictionRate", mock.Anything).Return(rate, nil).Maybe()
	// suite.rate = rate
	// // app1, app2 has issued IDA
	// suite.dev.On("GetIDAPrice", suite.Ctx, suite.app1).Return(linotypes.NewMiniDollar(10),nil)
	// suite.dev.On("GetIDAPrice", suite.Ctx, suite.app2).Return(linotypes.NewMiniDollar(7),nil)
}

func (suite *AccountManagerTestSuite) TestDoesAccountExist() {
	testCases := []struct {
		testName     string
		user         linotypes.AccountKey
		expectResult bool
	}{
		{
			testName:     "user does exists",
			user:         suite.userWithBalance.Username,
			expectResult: true,
		},
		{
			testName:     "user doesn't exists",
			user:         suite.userWithoutBalance.Username,
			expectResult: true,
		},
	}
	for _, tc := range testCases {
		res := suite.am.DoesAccountExist(suite.Ctx, tc.user)
		suite.Equal(
			tc.expectResult, res,
			"%s: does account exist for user %s, expect %t, got %t", tc.testName, tc.user, tc.expectResult, res)
	}
}

func (suite *AccountManagerTestSuite) TestAddCoinToAddress() {
	userWithBalance := suite.userWithBalance
	unreg := suite.unreg
	emptyAddress := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	testCases := []struct {
		testName   string
		address    sdk.AccAddress
		amount     types.Coin
		expectBank *model.AccountBank
	}{
		{
			testName: "add coin to bank which is linked to username",
			address:  sdk.AccAddress(userWithBalance.TransactionKey.Address()),
			amount:   c100,
			expectBank: &model.AccountBank{
				Saving:   suite.userWithBalanceSaving.Plus(c100),
				PubKey:   userWithBalance.TransactionKey,
				Username: userWithBalance.Username,
			},
		},
		{
			testName: "add coin to bank which is not linked to username",
			address:  sdk.AccAddress(unreg.TransactionKey.Address()),
			amount:   c100,
			expectBank: &model.AccountBank{
				Saving: suite.unregSaving.Plus(c100),
			},
		},
		{
			testName: "add coin to empty bank",
			address:  emptyAddress,
			amount:   c100,
			expectBank: &model.AccountBank{
				Saving: c100,
			},
		},
	}

	for _, tc := range testCases {
		err := suite.am.AddCoinToAddress(suite.Ctx, tc.address, tc.amount)
		suite.Nil(err, "%s: failed to add coin, got err: %v", tc.testName, err)
		suite.checkBankKVByAddress(tc.testName, tc.address, tc.expectBank)
	}
}

func (suite *AccountManagerTestSuite) TestAddCoinToUsername() {
	userWithBalance := suite.userWithBalance
	unreg := suite.unreg

	testCases := []struct {
		testName   string
		username   types.AccountKey
		amount     types.Coin
		expectErr  sdk.Error
		expectBank *model.AccountBank
	}{
		{
			testName:  "add coin to created username",
			username:  userWithBalance.Username,
			amount:    c100,
			expectErr: nil,
			expectBank: &model.AccountBank{
				Saving:   suite.userWithBalanceSaving.Plus(c100),
				PubKey:   userWithBalance.TransactionKey,
				Username: userWithBalance.Username,
			},
		},
		{
			testName:   "add coin to unregister username",
			username:   unreg.Username,
			amount:     c100,
			expectErr:  model.ErrAccountInfoNotFound(),
			expectBank: nil,
		},
	}

	for _, tc := range testCases {
		err := suite.am.AddCoinToUsername(suite.Ctx, tc.username, tc.amount)
		suite.Equal(
			tc.expectErr, err,
			"%s: failed to add coin to user %s, expect err %v, got %v",
			tc.testName, tc.username, tc.expectErr, err)
		if tc.expectBank != nil {
			suite.checkBankKVByUsername(tc.testName, tc.username, tc.expectBank)
		}
	}
}

func (suite *AccountManagerTestSuite) TestMinusCoinFromAddress() {
	userWithBalance := suite.userWithBalance
	userWithoutBalance := suite.userWithoutBalance
	unreg := suite.unreg
	emptyAddress := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	testCases := []struct {
		testName   string
		address    sdk.AccAddress
		amount     types.Coin
		expectErr  sdk.Error
		expectBank *model.AccountBank
	}{
		{
			testName:  "minus coin from address with sufficient balance",
			address:   userWithBalance.Address,
			expectErr: nil,
			amount:    coin100,
			expectBank: &model.AccountBank{
				Saving:   suite.userWithBalanceSaving.Minus(coin100),
				PubKey:   userWithBalance.TransactionKey,
				Username: userWithBalance.Username,
			},
		},
		{
			testName:  "minus coin from address without sufficient balance",
			address:   userWithoutBalance.Address,
			expectErr: acctypes.ErrAccountSavingCoinNotEnough(),
			amount:    coin1,
			expectBank: &model.AccountBank{
				PubKey:   userWithoutBalance.TransactionKey,
				Saving:   types.NewCoinFromInt64(0),
				Username: userWithoutBalance.Username,
			},
		},
		{
			testName:  "minus saving coin exceeds the coin address hold",
			address:   userWithBalance.Address,
			expectErr: acctypes.ErrAccountSavingCoinNotEnough(),
			amount:    suite.userWithBalanceSaving,
			expectBank: &model.AccountBank{
				Saving:   suite.userWithBalanceSaving.Minus(coin100),
				PubKey:   userWithBalance.TransactionKey,
				Username: userWithBalance.Username,
			},
		},
		{
			testName:  "minus saving coin from unregister address",
			address:   sdk.AccAddress(unreg.TransactionKey.Address()),
			expectErr: nil,
			amount:    coin100,
			expectBank: &model.AccountBank{
				Saving: suite.unregSaving.Minus(coin100),
			},
		},
		{
			testName:   "minus saving coin from empty address",
			address:    emptyAddress,
			expectErr:  model.ErrAccountBankNotFound(),
			amount:     coin1,
			expectBank: nil,
		},
	}
	for _, tc := range testCases {
		err := suite.am.MinusCoinFromAddress(suite.Ctx, tc.address, tc.amount)
		suite.Equal(
			tc.expectErr, err,
			"%s: failed to minus coin from address %s, expect err %v, got %v",
			tc.testName, tc.address, tc.expectErr, err)
		if tc.expectBank != nil {
			suite.checkBankKVByAddress(tc.testName, tc.address, tc.expectBank)
		}
	}
}

func (suite *AccountManagerTestSuite) TestMinusCoinFromUsername() {
	userWithBalance := suite.userWithBalance
	userWithoutBalance := suite.userWithoutBalance
	unreg := suite.unreg

	testCases := []struct {
		testName   string
		username   types.AccountKey
		amount     types.Coin
		expectErr  sdk.Error
		expectBank *model.AccountBank
	}{
		{
			testName:  "minus coin from user with sufficient balance",
			username:  userWithBalance.Username,
			expectErr: nil,
			amount:    coin100,
			expectBank: &model.AccountBank{
				Saving:   suite.userWithBalanceSaving.Minus(coin100),
				PubKey:   userWithBalance.TransactionKey,
				Username: userWithBalance.Username,
			},
		},
		{
			testName:  "minus coin from user without sufficient balance",
			username:  userWithoutBalance.Username,
			expectErr: acctypes.ErrAccountSavingCoinNotEnough(),
			amount:    coin1,
			expectBank: &model.AccountBank{
				PubKey:   userWithoutBalance.TransactionKey,
				Saving:   types.NewCoinFromInt64(0),
				Username: userWithoutBalance.Username,
			},
		},
		{
			testName:  "minus saving coin exceeds the coin user hold",
			username:  userWithBalance.Username,
			expectErr: acctypes.ErrAccountSavingCoinNotEnough(),
			amount:    suite.userWithBalanceSaving,
			expectBank: &model.AccountBank{
				Saving:   suite.userWithBalanceSaving.Minus(coin100),
				PubKey:   userWithBalance.TransactionKey,
				Username: userWithBalance.Username,
			},
		},
		{
			testName:   "minus saving coin from unregister account",
			username:   unreg.Username,
			expectErr:  acctypes.ErrAccountNotFound(unreg.Username),
			amount:     coin1,
			expectBank: nil,
		},
	}
	for _, tc := range testCases {
		err := suite.am.MinusCoinFromUsername(suite.Ctx, tc.username, tc.amount)
		suite.Equal(
			tc.expectErr, err,
			"%s: failed to minus coin from user %s, expect err %v, got %v",
			tc.testName, tc.username, tc.expectErr, err)
		if tc.expectBank != nil {
			suite.checkBankKVByUsername(tc.testName, tc.username, tc.expectBank)
		}
	}
}

func (suite *AccountManagerTestSuite) TestCreateAccount() {
	userWithBalance := suite.userWithBalance
	unreg := suite.unreg

	txKeyWithEmptyAddress := secp256k1.GenPrivKey().PubKey()
	signingKey := secp256k1.GenPrivKey().PubKey()
	txKey := secp256k1.GenPrivKey().PubKey()

	testCases := []struct {
		testName   string
		username   types.AccountKey
		signingKey crypto.PubKey
		txKey      crypto.PubKey
		expectErr  sdk.Error
		expectInfo *model.AccountInfo
		expectBank *model.AccountBank
	}{
		{
			testName:   "create account with registered username",
			username:   userWithBalance.Username,
			signingKey: userWithBalance.SigningKey,
			txKey:      userWithBalance.TransactionKey,
			expectErr:  acctypes.ErrAccountAlreadyExists(userWithBalance.Username),
			expectInfo: &userWithBalance,
			expectBank: &model.AccountBank{
				Saving:   suite.userWithBalanceSaving,
				PubKey:   userWithBalance.TransactionKey,
				Username: userWithBalance.Username,
			},
		},
		{
			testName:   "create account with bank linked to other username",
			username:   unreg.Username,
			signingKey: unreg.SigningKey,
			txKey:      userWithBalance.TransactionKey,
			expectErr:  acctypes.ErrAddressAlreadyTaken(sdk.AccAddress(userWithBalance.TransactionKey.Address())),
			expectInfo: nil,
			expectBank: &model.AccountBank{
				Saving:   suite.userWithBalanceSaving,
				PubKey:   userWithBalance.TransactionKey,
				Username: userWithBalance.Username,
			},
		},
		{
			testName:   "create account with exist address",
			username:   unreg.Username,
			signingKey: unreg.SigningKey,
			txKey:      unreg.TransactionKey,
			expectErr:  nil,
			expectInfo: &unreg,
			expectBank: &model.AccountBank{
				Saving:   suite.unregSaving,
				PubKey:   unreg.TransactionKey,
				Username: unreg.Username,
			},
		},
		{
			testName:   "create account with empty address",
			username:   "test1",
			signingKey: signingKey,
			txKey:      txKeyWithEmptyAddress,
			expectErr:  nil,
			expectInfo: &model.AccountInfo{
				Username:       "test1",
				SigningKey:     signingKey,
				TransactionKey: txKeyWithEmptyAddress,
				Address:        sdk.AccAddress(txKeyWithEmptyAddress.Address()),
			},
			expectBank: &model.AccountBank{
				Saving:   types.NewCoinFromInt64(0),
				PubKey:   txKeyWithEmptyAddress,
				Username: "test1",
			},
		},
		{
			testName:   "create account without signing key",
			username:   "test2",
			signingKey: nil,
			txKey:      txKey,
			expectErr:  nil,
			expectInfo: &model.AccountInfo{
				Username:       "test2",
				TransactionKey: txKey,
				Address:        sdk.AccAddress(txKey.Address()),
			},
			expectBank: &model.AccountBank{
				Saving:   types.NewCoinFromInt64(0),
				PubKey:   txKey,
				Username: "test2",
			},
		},
	}
	// normal test
	for _, tc := range testCases {
		err := suite.am.CreateAccount(suite.Ctx, tc.username, tc.signingKey, tc.txKey)
		suite.Equal(
			tc.expectErr, err,
			"%s: failed to create account for user %s, expect err %v, got %v",
			tc.testName, tc.username, tc.expectErr, err)
		if tc.expectBank != nil {
			suite.checkBankKVByAddress(tc.testName, sdk.AccAddress(tc.txKey.Address()), tc.expectBank)
		}
		if tc.expectInfo != nil {
			suite.checkInfoKVByUsername(tc.testName, tc.username, tc.expectInfo)
		}
	}
}

// func TestInvalidCreateAccount(t *testing.T) {
// 	ctx, am, _ := setupTest(t, 1)
// 	accParam, _ := am.paramHolder.GetAccountParam(ctx)
// 	priv1 := secp256k1.GenPrivKey()
// 	priv2 := secp256k1.GenPrivKey()

// 	accKey1 := types.AccountKey("accKey1")
// 	accKey2 := types.AccountKey("accKey2")

// 	testCases := []struct {
// 		testName        string
// 		username        types.AccountKey
// 		privKey         crypto.PrivKey
// 		registerDeposit types.Coin
// 		expectErr       sdk.Error
// 	}{
// 		{
// 			testName:        "register user with sufficient saving coin",
// 			username:        accKey1,
// 			privKey:         priv1,
// 			registerDeposit: accParam.RegisterFee,
// 			expectErr:       nil,
// 		},
// 		{
// 			testName:        "username already took",
// 			username:        accKey1,
// 			privKey:         priv1,
// 			registerDeposit: accParam.RegisterFee,
// 			expectErr:       ErrAccountAlreadyExists(accKey1),
// 		},
// 		{
// 			testName:        "username already took with different private key",
// 			username:        accKey1,
// 			privKey:         priv2,
// 			registerDeposit: accParam.RegisterFee,
// 			expectErr:       ErrAccountAlreadyExists(accKey1),
// 		},
// 		{
// 			testName:        "register the same private key",
// 			username:        accKey2,
// 			privKey:         priv1,
// 			registerDeposit: accParam.RegisterFee,
// 			expectErr:       nil,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		err := am.CreateAccount(
// 			ctx, accountReferrer, tc.username, tc.privKey.PubKey(),
// 			secp256k1.GenPrivKey().PubKey(),
// 			secp256k1.GenPrivKey().PubKey(), tc.registerDeposit)
// 		if !assert.Equal(t, tc.expectErr, err) {
// 			t.Errorf("%s: diff err, got %v, want %v", tc.testName, err, tc.expectErr)
// 		}
// 	}
// }

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

// func TestCheckUserTPSCapacity(t *testing.T) {
// 	ctx, am, _ := setupTest(t, 1)
// 	accKey := types.AccountKey("accKey")

// 	bandwidthParams, err := am.paramHolder.GetBandwidthParam(ctx)
// 	if err != nil {
// 		t.Errorf("TestCheckUserTPSCapacity: failed to get bandwidth param, got err %v", err)
// 	}
// 	virtualCoinAmount, _ := bandwidthParams.VirtualCoin.ToInt64()
// 	secondsToRecoverBandwidth := bandwidthParams.SecondsToRecoverBandwidth

// 	baseTime := ctx.BlockHeader().Time

// 	createTestAccount(ctx, am, string(accKey))
// 	err = am.AddSavingCoin(ctx, accKey, c100, "", "", types.TransferIn)
// 	if err != nil {
// 		t.Errorf("TestCheckUserTPSCapacity: failed to add saving coin, got err %v", err)
// 	}

// 	accStorage := model.NewAccountStorage(testAccountKVStoreKey)
// 	err = accStorage.SetPendingCoinDayQueue(
// 		ctx, accKey, &model.PendingCoinDayQueue{})
// 	if err != nil {
// 		t.Errorf("TestCheckUserTPSCapacity: failed to set pending coin day queue, got err %v", err)
// 	}

// 	testCases := []struct {
// 		testName             string
// 		tpsCapacityRatio     sdk.Dec
// 		userCoinDay          types.Coin
// 		lastActivity         int64
// 		lastCapacity         types.Coin
// 		currentTime          time.Time
// 		expectResult         sdk.Error
// 		expectRemainCapacity types.Coin
// 	}{
// 		{
// 			testName:             "tps capacity not enough",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 10),
// 			userCoinDay:          types.NewCoinFromInt64(10 * types.Decimals),
// 			lastActivity:         baseTime.Unix(),
// 			lastCapacity:         types.NewCoinFromInt64(0),
// 			currentTime:          baseTime,
// 			expectResult:         ErrAccountTPSCapacityNotEnough(accKey),
// 			expectRemainCapacity: types.NewCoinFromInt64(0),
// 		},
// 		{
// 			testName:             " 1/10 capacity ratio",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 10),
// 			userCoinDay:          types.NewCoinFromInt64(10 * types.Decimals),
// 			lastActivity:         baseTime.Unix(),
// 			lastCapacity:         types.NewCoinFromInt64(0),
// 			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth) * time.Second),
// 			expectResult:         nil,
// 			expectRemainCapacity: types.NewCoinFromInt64(990000).Plus(bandwidthParams.VirtualCoin),
// 		},
// 		{
// 			testName:             " 1/2 capacity ratio",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 2),
// 			userCoinDay:          types.NewCoinFromInt64(10 * types.Decimals),
// 			lastActivity:         baseTime.Unix(),
// 			lastCapacity:         types.NewCoinFromInt64(0),
// 			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth) * time.Second),
// 			expectResult:         nil,
// 			expectRemainCapacity: types.NewCoinFromInt64(950000).Plus(bandwidthParams.VirtualCoin),
// 		},
// 		{
// 			testName:             " 1/1 capacity ratio",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
// 			userCoinDay:          types.NewCoinFromInt64(10 * types.Decimals),
// 			lastActivity:         baseTime.Unix(),
// 			lastCapacity:         types.NewCoinFromInt64(0),
// 			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth) * time.Second),
// 			expectResult:         nil,
// 			expectRemainCapacity: types.NewCoinFromInt64(9 * types.Decimals).Plus(bandwidthParams.VirtualCoin),
// 		},
// 		{
// 			testName:             " 1/1 capacity ratio with virtual coin remaining",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
// 			userCoinDay:          types.NewCoinFromInt64(1 * types.Decimals),
// 			lastActivity:         baseTime.Unix(),
// 			lastCapacity:         types.NewCoinFromInt64(10 * types.Decimals),
// 			currentTime:          baseTime,
// 			expectResult:         nil,
// 			expectRemainCapacity: types.NewCoinFromInt64(1 * types.Decimals),
// 		},
// 		{
// 			testName:             " 1/1 capacity ratio with 1 coin day and 0 remaining",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
// 			userCoinDay:          types.NewCoinFromInt64(1 * types.Decimals),
// 			lastActivity:         baseTime.Unix(),
// 			lastCapacity:         types.NewCoinFromInt64(0),
// 			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth/2) * time.Second),
// 			expectResult:         nil,
// 			expectRemainCapacity: coin0,
// 		},
// 		{
// 			testName:             " transaction capacity not enough",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
// 			userCoinDay:          types.NewCoinFromInt64(0 * types.Decimals),
// 			lastActivity:         baseTime.Unix(),
// 			lastCapacity:         types.NewCoinFromInt64(0),
// 			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth/2) * time.Second),
// 			expectResult:         ErrAccountTPSCapacityNotEnough(accKey),
// 			expectRemainCapacity: coin0,
// 		},
// 		{
// 			testName:             " transaction capacity without coin day",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
// 			userCoinDay:          types.NewCoinFromInt64(0 * types.Decimals),
// 			lastActivity:         baseTime.Unix(),
// 			lastCapacity:         types.NewCoinFromInt64(0),
// 			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth) * time.Second),
// 			expectResult:         nil,
// 			expectRemainCapacity: coin0,
// 		},
// 		{
// 			testName:             " 1/2 capacity ratio with half virtual coin remaining",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 2),
// 			userCoinDay:          types.NewCoinFromInt64(1 * types.Decimals),
// 			lastActivity:         baseTime.Unix(),
// 			lastCapacity:         types.NewCoinFromInt64(0),
// 			currentTime:          baseTime.Add(time.Duration(secondsToRecoverBandwidth/2) * time.Second),
// 			expectResult:         nil,
// 			expectRemainCapacity: types.NewCoinFromInt64(virtualCoinAmount / 2),
// 		},
// 		{
// 			testName:             " 1/1 capacity ratio with virtual coin remaining and base time",
// 			tpsCapacityRatio:     types.NewDecFromRat(1, 1),
// 			userCoinDay:          types.NewCoinFromInt64(1 * types.Decimals),
// 			lastActivity:         0,
// 			lastCapacity:         types.NewCoinFromInt64(0),
// 			currentTime:          baseTime,
// 			expectResult:         nil,
// 			expectRemainCapacity: bandwidthParams.VirtualCoin,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: tc.currentTime})
// 		bank := &model.AccountBank{
// 			Saving:  tc.userCoinDay,
// 			CoinDay: tc.userCoinDay,
// 		}
// 		err := accStorage.SetBankFromAccountKey(ctx, accKey, bank)
// 		if err != nil {
// 			t.Errorf("%s: failed to set bank, got err %v", tc.testName, err)
// 		}

// 		meta := &model.AccountMeta{
// 			LastActivityAt:      tc.lastActivity,
// 			TransactionCapacity: tc.lastCapacity,
// 		}
// 		err = accStorage.SetMeta(ctx, accKey, meta)
// 		if err != nil {
// 			t.Errorf("%s: failed to set meta, got err %v", tc.testName, err)
// 		}

// 		err = am.CheckUserTPSCapacity(ctx, accKey, tc.tpsCapacityRatio)
// 		if !assert.Equal(t, tc.expectResult, err) {
// 			t.Errorf("%s: diff tps capacity, got %v, want %v", tc.testName, err, tc.expectResult)
// 		}

// 		accMeta := model.AccountMeta{
// 			LastActivityAt:      ctx.BlockHeader().Time.Unix(),
// 			TransactionCapacity: tc.expectRemainCapacity,
// 		}
// 		if tc.expectResult != nil {
// 			accMeta.LastActivityAt = tc.lastActivity
// 		}
// 		checkAccountMeta(t, ctx, tc.testName, accKey, accMeta)
// 	}
// }

// func TestCheckAuthenticatePubKeyOwner(t *testing.T) {
// 	testName := "TestCheckAuthenticatePubKeyOwner"

// 	ctx, am, _ := setupTest(t, 1)
// 	accParam, _ := am.paramHolder.GetAccountParam(ctx)
// 	user1 := types.AccountKey("user1")
// 	appPermissionUser := types.AccountKey("user2")
// 	preAuthPermissionUser := types.AccountKey("user3")
// 	unauthUser := types.AccountKey("user4")
// 	resetKey := secp256k1.GenPrivKey()
// 	transactionKey := secp256k1.GenPrivKey()
// 	appKey := secp256k1.GenPrivKey()
// 	am.CreateAccount(
// 		ctx, accountReferrer, user1, resetKey.PubKey(), transactionKey.PubKey(),
// 		appKey.PubKey(), accParam.RegisterFee)

// 	_, unauthTxPriv, authAppPriv := createTestAccount(ctx, am, string(appPermissionUser))
// 	_, authTxPriv, unauthAppPriv := createTestAccount(ctx, am, string(preAuthPermissionUser))
// 	_, unauthPriv1, unauthPriv2 := createTestAccount(ctx, am, string(unauthUser))

// 	err := am.AuthorizePermission(ctx, user1, appPermissionUser, 100, types.AppPermission, types.NewCoinFromInt64(0))
// 	if err != nil {
// 		t.Errorf("%s: failed to authorize app permission, got err %v", testName, err)
// 	}

// 	preAuthAmount := types.NewCoinFromInt64(100)
// 	err = am.AuthorizePermission(ctx, user1, preAuthPermissionUser, 100, types.PreAuthorizationPermission, preAuthAmount)
// 	if err != nil {
// 		t.Errorf("%s: failed to authorize preauth permission, got err %v", testName, err)
// 	}
// 	baseTime := ctx.BlockHeader().Time

// 	testCases := []struct {
// 		testName           string
// 		checkUser          types.AccountKey
// 		checkPubKey        crypto.PubKey
// 		atWhen             time.Time
// 		amount             types.Coin
// 		permission         types.Permission
// 		expectUser         types.AccountKey
// 		expectResult       sdk.Error
// 		expectGrantPubKeys []*model.GrantPermission
// 	}{
// 		{
// 			testName:           "check user's reset key",
// 			checkUser:          user1,
// 			checkPubKey:        resetKey.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(0),
// 			permission:         types.ResetPermission,
// 			expectUser:         user1,
// 			expectResult:       nil,
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "check user's transaction key",
// 			checkUser:          user1,
// 			checkPubKey:        transactionKey.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(0),
// 			permission:         types.TransactionPermission,
// 			expectUser:         user1,
// 			expectResult:       nil,
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "check user's app key",
// 			checkUser:          user1,
// 			checkPubKey:        appKey.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(0),
// 			permission:         types.AppPermission,
// 			expectUser:         user1,
// 			expectResult:       nil,
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "user's transaction key can authorize grant app permission",
// 			checkUser:          user1,
// 			checkPubKey:        transactionKey.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(0),
// 			permission:         types.GrantAppPermission,
// 			expectUser:         user1,
// 			expectResult:       nil,
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "user's transaction key can authorize app permission",
// 			checkUser:          user1,
// 			checkPubKey:        transactionKey.PubKey(),
// 			atWhen:             baseTime,
// 			permission:         types.AppPermission,
// 			expectUser:         user1,
// 			expectResult:       nil,
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "check user's transaction key can't authorize reset permission",
// 			checkUser:          user1,
// 			checkPubKey:        transactionKey.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(0),
// 			permission:         types.ResetPermission,
// 			expectUser:         user1,
// 			expectResult:       ErrCheckResetKey(),
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "check user's app key can authorize grant app permission",
// 			checkUser:          user1,
// 			checkPubKey:        appKey.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(0),
// 			permission:         types.GrantAppPermission,
// 			expectUser:         user1,
// 			expectResult:       nil,
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "check user's app key can't authorize transaction permission",
// 			checkUser:          user1,
// 			checkPubKey:        appKey.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(0),
// 			permission:         types.TransactionPermission,
// 			expectUser:         user1,
// 			expectResult:       ErrCheckTransactionKey(),
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "check user's app key can't authorize reset permission",
// 			checkUser:          user1,
// 			checkPubKey:        appKey.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(0),
// 			permission:         types.ResetPermission,
// 			expectUser:         user1,
// 			expectResult:       ErrCheckResetKey(),
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:     "check app pubkey of user with app permission",
// 			checkUser:    user1,
// 			checkPubKey:  authAppPriv.PubKey(),
// 			atWhen:       baseTime,
// 			amount:       types.NewCoinFromInt64(0),
// 			permission:   types.AppPermission,
// 			expectUser:   appPermissionUser,
// 			expectResult: nil,
// 			expectGrantPubKeys: []*model.GrantPermission{
// 				&model.GrantPermission{
// 					GrantTo:    appPermissionUser,
// 					Permission: types.AppPermission,
// 					CreatedAt:  baseTime.Unix(),
// 					ExpiresAt:  baseTime.Unix() + 100,
// 					Amount:     types.NewCoinFromInt64(0),
// 				},
// 			},
// 		},
// 		{
// 			testName:           "check transaction pubkey of user with app permission",
// 			checkUser:          user1,
// 			checkPubKey:        unauthTxPriv.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(0),
// 			permission:         types.PreAuthorizationPermission,
// 			expectUser:         "",
// 			expectResult:       nil,
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "check unauthorized user app pubkey",
// 			checkUser:          user1,
// 			checkPubKey:        unauthPriv2.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(10),
// 			permission:         types.AppPermission,
// 			expectUser:         "",
// 			expectResult:       ErrCheckAuthenticatePubKeyOwner(user1),
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "check unauthorized user transaction pubkey",
// 			checkUser:          user1,
// 			checkPubKey:        unauthPriv1.PubKey(),
// 			atWhen:             baseTime,
// 			amount:             types.NewCoinFromInt64(10),
// 			permission:         types.PreAuthorizationPermission,
// 			expectUser:         "",
// 			expectResult:       ErrCheckAuthenticatePubKeyOwner(user1),
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:     "check transaction pubkey of user with preauthorization permission",
// 			checkUser:    user1,
// 			checkPubKey:  authTxPriv.PubKey(),
// 			atWhen:       baseTime,
// 			amount:       types.NewCoinFromInt64(10),
// 			permission:   types.PreAuthorizationPermission,
// 			expectUser:   preAuthPermissionUser,
// 			expectResult: nil,
// 			expectGrantPubKeys: []*model.GrantPermission{
// 				&model.GrantPermission{
// 					GrantTo:    preAuthPermissionUser,
// 					Permission: types.PreAuthorizationPermission,
// 					CreatedAt:  baseTime.Unix(),
// 					ExpiresAt:  baseTime.Unix() + 100,
// 					Amount:     preAuthAmount.Minus(types.NewCoinFromInt64(10)),
// 				},
// 			},
// 		},
// 		{
// 			testName:     "check app pubkey of user with preauthorization permission",
// 			checkUser:    user1,
// 			checkPubKey:  unauthAppPriv.PubKey(),
// 			atWhen:       baseTime,
// 			amount:       types.NewCoinFromInt64(10),
// 			permission:   types.AppPermission,
// 			expectUser:   preAuthPermissionUser,
// 			expectResult: ErrCheckAuthenticatePubKeyOwner(user1),
// 			expectGrantPubKeys: []*model.GrantPermission{
// 				&model.GrantPermission{
// 					GrantTo:    preAuthPermissionUser,
// 					Permission: types.PreAuthorizationPermission,
// 					CreatedAt:  baseTime.Unix(),
// 					ExpiresAt:  baseTime.Unix() + 100,
// 					Amount:     preAuthAmount.Minus(types.NewCoinFromInt64(10)),
// 				},
// 			},
// 		},
// 		{
// 			testName:    "check amount exceeds preauthorization limitation",
// 			checkUser:   user1,
// 			checkPubKey: authTxPriv.PubKey(),
// 			atWhen:      baseTime,
// 			amount:      preAuthAmount,
// 			permission:  types.PreAuthorizationPermission,
// 			expectUser:  preAuthPermissionUser,
// 			expectResult: ErrPreAuthAmountInsufficient(
// 				preAuthPermissionUser, preAuthAmount.Minus(types.NewCoinFromInt64(10)),
// 				preAuthAmount),
// 			expectGrantPubKeys: []*model.GrantPermission{
// 				&model.GrantPermission{
// 					GrantTo:    preAuthPermissionUser,
// 					Permission: types.PreAuthorizationPermission,
// 					CreatedAt:  baseTime.Unix(),
// 					ExpiresAt:  baseTime.Unix() + 100,
// 					Amount:     preAuthAmount.Minus(types.NewCoinFromInt64(10)),
// 				},
// 			},
// 		},
// 		{
// 			testName:     "check grant app key can't sign grant app msg",
// 			checkUser:    user1,
// 			checkPubKey:  authAppPriv.PubKey(),
// 			atWhen:       baseTime,
// 			permission:   types.GrantAppPermission,
// 			expectUser:   appPermissionUser,
// 			expectResult: ErrCheckGrantAppKey(),
// 			expectGrantPubKeys: []*model.GrantPermission{
// 				&model.GrantPermission{
// 					GrantTo:    appPermissionUser,
// 					Permission: types.AppPermission,
// 					CreatedAt:  baseTime.Unix(),
// 					ExpiresAt:  baseTime.Unix() + 100,
// 					Amount:     types.NewCoinFromInt64(0),
// 				},
// 			},
// 		},
// 		{
// 			testName:           "check expired app permission",
// 			checkUser:          user1,
// 			checkPubKey:        authAppPriv.PubKey(),
// 			atWhen:             baseTime.Add(time.Duration(101) * time.Second),
// 			permission:         types.AppPermission,
// 			expectUser:         "",
// 			expectResult:       ErrCheckAuthenticatePubKeyOwner(user1),
// 			expectGrantPubKeys: nil,
// 		},
// 		{
// 			testName:           "check expired preauth permission",
// 			checkUser:          user1,
// 			checkPubKey:        authTxPriv.PubKey(),
// 			atWhen:             baseTime.Add(time.Duration(101) * time.Second),
// 			amount:             types.NewCoinFromInt64(100),
// 			permission:         types.PreAuthorizationPermission,
// 			expectUser:         "",
// 			expectResult:       ErrCheckAuthenticatePubKeyOwner(user1),
// 			expectGrantPubKeys: nil,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: tc.atWhen})
// 		keyOwner, err := am.CheckSigningPubKeyOwner(ctx, tc.checkUser, tc.checkPubKey, tc.permission, tc.amount)
// 		if tc.expectResult == nil {
// 			if tc.expectUser != keyOwner {
// 				t.Errorf("%s: diff key owner,  got %v, want %v", tc.testName, keyOwner, tc.expectUser)
// 				return
// 			}
// 		} else {
// 			fmt.Println(tc.testName, tc.expectResult.Result(), err)
// 			if !assert.Equal(t, tc.expectResult.Result(), err.Result()) {
// 				t.Errorf("%s: diff result,  got %v, want %v", tc.testName, err.Result(), tc.expectResult.Result())
// 			}
// 		}

// 		grantPubKeys, err := am.storage.GetGrantPermissions(ctx, tc.checkUser, tc.expectUser)
// 		if tc.expectGrantPubKeys == nil {
// 			if err == nil {
// 				t.Errorf("%s: got nil err", tc.testName)
// 			}
// 		} else {
// 			if err != nil {
// 				t.Errorf("%s: got non-empty err %v", tc.testName, err)
// 			}
// 			if len(tc.expectGrantPubKeys) != len(grantPubKeys) {
// 				t.Errorf("%s: expect grant pubkey length is different,  got %v, want %v", tc.testName, len(grantPubKeys), len(tc.expectGrantPubKeys))
// 			}
// 		}
// 	}
// }

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
			expectResult:   acctypes.ErrAccountNotFound(nonExistUser),
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

// func TestAccountRecoverNormalCase(t *testing.T) {
// 	testName := "TestAccountRecoverNormalCase"

// 	ctx, am, _ := setupTest(t, 1)
// 	accParam, _ := am.paramHolder.GetAccountParam(ctx)
// 	user1 := types.AccountKey("user1")

// 	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
// 	if err != nil {
// 		t.Errorf("%s: failed to get coin day param relationship, got err %v", testName, err)
// 	}

// 	createTestAccount(ctx, am, string(user1))

// 	newResetPrivKey := secp256k1.GenPrivKey()
// 	newTransactionPrivKey := secp256k1.GenPrivKey()
// 	newAppPrivKey := secp256k1.GenPrivKey()

// 	err = am.RecoverAccount(
// 		ctx, user1, newResetPrivKey.PubKey(), newTransactionPrivKey.PubKey(),
// 		newAppPrivKey.PubKey())
// 	if err != nil {
// 		t.Errorf("%s: failed to recover account, got err %v", testName, err)
// 	}

// 	accInfo := model.AccountInfo{
// 		Username:       user1,
// 		CreatedAt:      ctx.BlockHeader().Time.Unix(),
// 		ResetKey:       newResetPrivKey.PubKey(),
// 		TransactionKey: newTransactionPrivKey.PubKey(),
// 		AppKey:         newAppPrivKey.PubKey(),
// 	}
// 	bank := model.AccountBank{
// 		Saving:  accParam.RegisterFee,
// 		CoinDay: accParam.RegisterFee,
// 	}

// 	checkAccountInfo(t, ctx, testName, user1, accInfo)
// 	checkBankKVByUsername(t, ctx, testName, user1, bank)

// 	pendingCoinDayQueue := model.PendingCoinDayQueue{
// 		TotalCoinDay: sdk.ZeroDec(),
// 		TotalCoin:    types.NewCoinFromInt64(0),
// 	}
// 	checkPendingCoinDay(t, ctx, testName, user1, pendingCoinDayQueue)

// 	coinDay, err := am.GetCoinDay(ctx, user1)
// 	if err != nil {
// 		t.Errorf("%s: failed to get coin day, got err %v", testName, err)
// 	}
// 	if !coinDay.IsEqual(accParam.RegisterFee) {
// 		t.Errorf("%s: diff coin day, got %v, want %v", testName, coinDay, accParam.RegisterFee)
// 	}

// 	ctx = ctx.WithBlockHeader(
// 		abci.Header{
// 			ChainID: "Lino", Height: 1,
// 			Time: ctx.BlockHeader().Time.Add(time.Duration(coinDayParams.SecondsToRecoverCoinDay) * time.Second)})
// 	coinDay, err = am.GetCoinDay(ctx, user1)
// 	if err != nil {
// 		t.Errorf("%s: failed to get coin day again, got err %v", testName, err)
// 	}
// 	if !coinDay.IsEqual(accParam.RegisterFee) {
// 		t.Errorf("%s: diff coin day again, got %v, want %v", testName, coinDay, accParam.RegisterFee)
// 	}
// }

func TestIncreaseSequenceByOne(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	createTestAccount(ctx, am, string(user1))

	addr, err := am.GetAddress(ctx, user1)
	if err != nil {
		t.Errorf("TestIncreaseSequenceByOne: failed to get address, got err %v", err)
	}

	testCases := []struct {
		testName       string
		increaseTimes  int
		expectSequence uint64
	}{
		{
			testName:       "increase seq once",
			increaseTimes:  1,
			expectSequence: 1,
		},
		{
			testName:       "increase seq 100 times",
			increaseTimes:  100,
			expectSequence: 101,
		},
	}

	for _, tc := range testCases {
		for i := 0; i < tc.increaseTimes; i++ {
			am.IncreaseSequenceByOne(ctx, addr)
		}
		seq, err := am.GetSequence(ctx, addr)
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
	addr, err := am.GetAddress(ctx, user1)
	if err != nil {
		t.Errorf("TestAddFrozenMoney: failed to get address, got err %v", err)
	}

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
			interval:                10 * 3600,
			times:                   5,
			expectNumOfFrozenAmount: 1,
		},
		{
			testName:                "add the second 100 frozen money, clear the first one",
			frozenAmount:            types.NewCoinFromInt64(100),
			startAt:                 1200000,
			interval:                10 * 3600,
			times:                   5,
			expectNumOfFrozenAmount: 1,
		},
		{
			testName:                "add the third 100 frozen money",
			frozenAmount:            types.NewCoinFromInt64(100),
			startAt:                 1300000,
			interval:                10 * 3600,
			times:                   5,
			expectNumOfFrozenAmount: 2,
		},
		{
			testName:                "add the fourth 100 frozen money, clear the second one",
			frozenAmount:            types.NewCoinFromInt64(100),
			startAt:                 1400000,
			interval:                10 * 3600,
			times:                   5,
			expectNumOfFrozenAmount: 2,
		},
		{
			testName:                "add the fifth 100 frozen money, clear the third and fourth ones",
			frozenAmount:            types.NewCoinFromInt64(100),
			startAt:                 1600000,
			interval:                10 * 3600,
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

		accountBank, err := am.storage.GetBank(ctx, addr)
		if err != nil {
			t.Errorf("%s: failed to get bank, got err %v", tc.testName, err)
		}
		if len(accountBank.FrozenMoneyList) != tc.expectNumOfFrozenAmount {
			t.Errorf("%s: diff num of frozen money, got %v, want %v", tc.testName, len(accountBank.FrozenMoneyList), tc.expectNumOfFrozenAmount)
		}
	}
}

func (suite *AccountManagerTestSuite) checkInfoKVByUsername(testName string, username types.AccountKey, info *model.AccountInfo) {
	infoPtr, err := suite.am.storage.GetInfo(suite.Ctx, username)
	suite.Nil(err, "%s, failed to get info, got err %v", testName, err)
	suite.Equal(info, infoPtr, "%s: diff info, got %v, want %v", testName, *infoPtr, info)
}

func (suite *AccountManagerTestSuite) checkBankKVByAddress(testName string, address sdk.AccAddress, bank *model.AccountBank) {
	bankPtr, err := suite.am.storage.GetBank(suite.Ctx, address)
	suite.Nil(err, "%s, failed to get account bank, got err %v", testName, err)
	suite.Equal(bank, bankPtr, "%s: diff bank, got %v, want %v", testName, *bankPtr, bank)
}

func (suite *AccountManagerTestSuite) checkBankKVByUsername(testName string, username types.AccountKey, bank *model.AccountBank) {
	info, err := suite.am.storage.GetInfo(suite.Ctx, username)
	suite.Nil(err, "%s, failed to get info, got err %v", testName, err)
	suite.checkBankKVByAddress(testName, info.Address, bank)
}
