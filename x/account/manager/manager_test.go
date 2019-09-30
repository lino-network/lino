package manager

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/types"
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
	registerFee           types.Coin

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
		Username:       types.AccountKey("userwithoutbalance"),
		SigningKey:     secp256k1.GenPrivKey().PubKey(),
		TransactionKey: secp256k1.GenPrivKey().PubKey(),
	}
	suite.userWithoutBalance.Address = sdk.AccAddress(suite.userWithoutBalance.TransactionKey.Address())

	suite.userWithBalance = model.AccountInfo{
		Username:       types.AccountKey("userwithbalance"),
		SigningKey:     secp256k1.GenPrivKey().PubKey(),
		TransactionKey: secp256k1.GenPrivKey().PubKey(),
	}
	suite.userWithBalance.Address = sdk.AccAddress(suite.userWithBalance.TransactionKey.Address())

	suite.unreg = model.AccountInfo{
		Username:       types.AccountKey("unreg"),
		SigningKey:     secp256k1.GenPrivKey().PubKey(),
		TransactionKey: secp256k1.GenPrivKey().PubKey(),
	}
	suite.unreg.Address = sdk.AccAddress(suite.unreg.TransactionKey.Address())

	suite.userWithBalanceSaving = types.NewCoinFromInt64(1000 * types.Decimals)
	suite.unregSaving = types.NewCoinFromInt64(1 * types.Decimals)
	suite.registerFee = types.NewCoinFromInt64(100 * types.Decimals)

	err := suite.am.CreateAccount(suite.Ctx, suite.userWithoutBalance.Username, suite.userWithoutBalance.SigningKey, suite.userWithoutBalance.TransactionKey)
	suite.NoError(err)

	err = suite.am.CreateAccount(suite.Ctx, suite.userWithBalance.Username, suite.userWithBalance.SigningKey, suite.userWithBalance.TransactionKey)
	suite.NoError(err)
	err = suite.am.AddCoinToUsername(suite.Ctx, suite.userWithBalance.Username, suite.userWithBalanceSaving)
	suite.NoError(err)

	err = suite.am.AddCoinToAddress(suite.Ctx, sdk.AccAddress(suite.unreg.TransactionKey.Address()), suite.unregSaving)
	suite.NoError(err)

	suite.ph.On("GetAccountParam", mock.Anything).Return(&parammodel.AccountParam{
		RegisterFee:       suite.registerFee,
		MinimumBalance:    types.NewCoinFromInt64(0),
		MaxNumFrozenMoney: 10,
	}, nil).Maybe()

	// // reg accounts
	// for _, v := range []types.AccountKey{suite.user1, suite.user2, suite.app1, suite.app2, suite.app3} {
	// 	suite.am.On("DoesAccountExist", mock.Anything, v).Return(true).Maybe()
	// }
	// // unreg accounts
	// for _, v := range []types.AccountKey{suite.unreg} {
	// 	suite.am.On("DoesAccountExist", mock.Anything, v).Return(false).Maybe()
	// }

	// // reg dev
	// for _, v := range []types.AccountKey{suite.app1, suite.app2, suite.app3} {
	// 	suite.dev.On("DoesDeveloperExist", mock.Anything, v).Return(true).Maybe()
	// }
	// // unreg devs
	// for _, v := range []types.AccountKey{suite.unreg, suite.user1, suite.user2} {
	// 	suite.dev.On("DoesDeveloperExist", mock.Anything, v).Return(false).Maybe()
	// }

	// rate, err := sdk.NewDecFromStr("0.099")
	// suite.Require().Nil(err)
	// suite.global.On("GetConsumptionFrictionRate", mock.Anything).Return(rate, nil).Maybe()
	// suite.rate = rate
	// // app1, app2 has issued IDA
	// suite.dev.On("GetIDAPrice", suite.Ctx, suite.app1).Return(types.NewMiniDollar(10),nil)
	// suite.dev.On("GetIDAPrice", suite.Ctx, suite.app2).Return(types.NewMiniDollar(7),nil)
}

func (suite *AccountManagerTestSuite) TestDoesAccountExist() {
	testCases := []struct {
		testName     string
		user         types.AccountKey
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
			expectErr:  acctypes.ErrAccountNotFound(unreg.Username),
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
			expectErr:  acctypes.ErrAddressAlreadyTaken(hex.EncodeToString(userWithBalance.TransactionKey.Address())),
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

		accMeta := am.storage.GetMeta(ctx, tc.username)
		if tc.JSONMeta != accMeta.JSONMeta {
			t.Errorf("%s: diff json meta, got %v, want %v", tc.testName, accMeta.JSONMeta, tc.JSONMeta)
		}
	}
}

func (suite *AccountManagerTestSuite) TestRegisterAccount() {
	txPrivKeys := []crypto.PrivKey{secp256k1.GenPrivKey(), secp256k1.GenPrivKey()}
	signingPrivKeys := []crypto.PrivKey{secp256k1.GenPrivKey(), secp256k1.GenPrivKey()}
	suite.global.On("AddToValidatorInflationPool", mock.Anything, suite.registerFee).Return(nil).Maybe()

	testCases := []struct {
		testName    string
		referrer    types.AccountKey
		registerFee types.Coin
		username    types.AccountKey
		signingKey  crypto.PubKey
		txKey       crypto.PubKey
		expectErr   sdk.Error
		accInfo     *model.AccountInfo
		accBank     *model.AccountBank
	}{
		{
			testName:    "register username already exists",
			referrer:    suite.userWithBalance.Username,
			registerFee: suite.registerFee,
			username:    suite.userWithoutBalance.Username,
			signingKey:  secp256k1.GenPrivKey().PubKey(),
			txKey:       secp256k1.GenPrivKey().PubKey(),
			expectErr:   acctypes.ErrAccountAlreadyExists(suite.userWithoutBalance.Username),
			accInfo:     &suite.userWithoutBalance,
			accBank: &model.AccountBank{
				Saving:   types.NewCoinFromInt64(0),
				Username: suite.userWithoutBalance.Username,
				PubKey:   suite.userWithoutBalance.TransactionKey,
			},
		},
		{
			testName:    "register fee not enough",
			referrer:    suite.userWithBalance.Username,
			registerFee: suite.registerFee.Minus(types.NewCoinFromInt64(1)),
			username:    "test1",
			signingKey:  secp256k1.GenPrivKey().PubKey(),
			txKey:       secp256k1.GenPrivKey().PubKey(),
			expectErr:   acctypes.ErrRegisterFeeInsufficient(),
			accInfo:     nil,
			accBank:     nil,
		},
		{
			testName:    "register success",
			referrer:    suite.userWithBalance.Username,
			registerFee: suite.registerFee,
			username:    "test1",
			signingKey:  signingPrivKeys[0].PubKey(),
			txKey:       txPrivKeys[0].PubKey(),
			expectErr:   nil,
			accInfo: &model.AccountInfo{
				Username:       "test1",
				SigningKey:     signingPrivKeys[0].PubKey(),
				TransactionKey: txPrivKeys[0].PubKey(),
				Address:        sdk.AccAddress(txPrivKeys[0].PubKey().Address()),
			},
			accBank: &model.AccountBank{
				Saving:   types.NewCoinFromInt64(0),
				Username: "test1",
				PubKey:   txPrivKeys[0].PubKey(),
			},
		},
		{
			testName:    "register with same transaction private key",
			referrer:    suite.userWithBalance.Username,
			registerFee: suite.registerFee,
			username:    "test2",
			signingKey:  signingPrivKeys[0].PubKey(),
			txKey:       txPrivKeys[0].PubKey(),
			expectErr: acctypes.ErrAddressAlreadyTaken(
				hex.EncodeToString(sdk.AccAddress(txPrivKeys[0].PubKey().Address()))),
			accInfo: nil,
			accBank: nil,
		},
		{
			testName:    "referrer is address",
			referrer:    types.AccountKey(suite.userWithBalance.TransactionKey.Address()),
			registerFee: suite.registerFee,
			username:    "test3",
			signingKey:  signingPrivKeys[1].PubKey(),
			txKey:       txPrivKeys[1].PubKey(),
			expectErr:   nil,
			accInfo: &model.AccountInfo{
				Username:       "test3",
				SigningKey:     signingPrivKeys[1].PubKey(),
				TransactionKey: txPrivKeys[1].PubKey(),
				Address:        sdk.AccAddress(txPrivKeys[1].PubKey().Address()),
			},
			accBank: &model.AccountBank{
				Saving:   types.NewCoinFromInt64(0),
				Username: "test3",
				PubKey:   txPrivKeys[1].PubKey(),
			},
		},
	}
	for _, tc := range testCases {
		err := suite.am.RegisterAccount(suite.Ctx, tc.referrer, tc.registerFee, tc.username, tc.signingKey, tc.txKey)
		suite.Equal(tc.expectErr, err)
		bank, _ := suite.am.GetBank(suite.Ctx, tc.username)
		suite.Equal(tc.accBank, bank)
		info, _ := suite.am.GetInfo(suite.Ctx, tc.username)
		suite.Equal(tc.accInfo, info)
	}
}

func (suite *AccountManagerTestSuite) TestMoveCoin() {
	testCases := []struct {
		testName              string
		sender                types.AccountKey
		amount                types.Coin
		receiver              types.AccountKey
		expectErr             sdk.Error
		expectSenderBalance   types.Coin
		expectReceiverBalance types.Coin
	}{
		{
			testName:              "sender doesnt exist",
			sender:                "movecointest",
			receiver:              suite.userWithoutBalance.Username,
			amount:                types.NewCoinFromInt64(1),
			expectErr:             acctypes.ErrAccountNotFound("movecointest"),
			expectSenderBalance:   types.Coin{},
			expectReceiverBalance: types.NewCoinFromInt64(0),
		},
		{
			testName:              "receiver doesnt exist",
			sender:                suite.userWithBalance.Username,
			receiver:              "movecointest",
			amount:                types.NewCoinFromInt64(1),
			expectErr:             acctypes.ErrAccountNotFound("movecointest"),
			expectSenderBalance:   suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(1)),
			expectReceiverBalance: types.Coin{},
		},
		{
			testName:              "send from username to username",
			sender:                suite.userWithBalance.Username,
			receiver:              suite.userWithoutBalance.Username,
			amount:                types.NewCoinFromInt64(1),
			expectErr:             nil,
			expectSenderBalance:   suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(2)),
			expectReceiverBalance: types.NewCoinFromInt64(1),
		},
		{
			testName:              "send from username to address",
			sender:                suite.userWithBalance.Username,
			receiver:              types.AccountKey(suite.userWithoutBalance.TransactionKey.Address()),
			amount:                types.NewCoinFromInt64(1),
			expectErr:             nil,
			expectSenderBalance:   suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(3)),
			expectReceiverBalance: types.NewCoinFromInt64(2),
		},
		{
			testName:              "send from address to address",
			sender:                types.AccountKey(suite.userWithBalance.TransactionKey.Address()),
			receiver:              types.AccountKey(suite.userWithoutBalance.TransactionKey.Address()),
			amount:                types.NewCoinFromInt64(1),
			expectErr:             nil,
			expectSenderBalance:   suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(4)),
			expectReceiverBalance: types.NewCoinFromInt64(3),
		},
	}
	for _, tc := range testCases {
		err := suite.am.MoveCoin(suite.Ctx, tc.sender, tc.receiver, tc.amount)
		suite.Equal(tc.expectErr, err)
		if tc.sender.IsUsername() {
			saving, _ := suite.am.GetSavingFromUsername(suite.Ctx, tc.sender)
			suite.Equal(tc.expectSenderBalance, saving)
		} else {
			saving, _ := suite.am.GetSavingFromAddress(suite.Ctx, sdk.AccAddress(tc.sender))
			suite.Equal(tc.expectSenderBalance, saving)
		}
		if tc.receiver.IsUsername() {
			saving, _ := suite.am.GetSavingFromUsername(suite.Ctx, tc.receiver)
			suite.Equal(tc.expectReceiverBalance, saving)
		} else {
			saving, _ := suite.am.GetSavingFromAddress(suite.Ctx, sdk.AccAddress(tc.receiver))
			suite.Equal(tc.expectReceiverBalance, saving)
		}
	}
}

func (suite *AccountManagerTestSuite) TestCheckSigningPubKeyOwnerByAddress() {
	txPrivKeys := []crypto.PrivKey{secp256k1.GenPrivKey(), secp256k1.GenPrivKey()}
	testCases := []struct {
		testName      string
		address       sdk.AccAddress
		signKey       crypto.PubKey
		isPaid        bool
		expectErr     sdk.Error
		expectAccBank *model.AccountBank
	}{
		{
			testName:      "bank doesn't exist",
			address:       sdk.AccAddress(txPrivKeys[0].PubKey().Address()),
			signKey:       txPrivKeys[0].PubKey(),
			isPaid:        false,
			expectErr:     model.ErrAccountBankNotFound(),
			expectAccBank: nil,
		},
		{
			testName:  "set bank to paid address",
			address:   sdk.AccAddress(txPrivKeys[0].PubKey().Address()),
			signKey:   txPrivKeys[0].PubKey(),
			isPaid:    true,
			expectErr: nil,
			expectAccBank: &model.AccountBank{
				Saving: types.NewCoinFromInt64(0),
				PubKey: txPrivKeys[0].PubKey(),
			},
		},
		{
			testName: "signing key mismatch",
			address:  sdk.AccAddress(suite.unreg.TransactionKey.Address()),
			signKey:  txPrivKeys[0].PubKey(),
			isPaid:   false,
			expectErr: sdk.ErrInvalidPubKey(
				fmt.Sprintf("PubKey does not match Signer address %s", sdk.AccAddress(suite.unreg.TransactionKey.Address()))),
			expectAccBank: &model.AccountBank{
				Saving: suite.unregSaving,
			},
		},
		{
			testName:  "set public key to bank without public key info",
			address:   sdk.AccAddress(suite.unreg.TransactionKey.Address()),
			signKey:   suite.unreg.TransactionKey,
			isPaid:    false,
			expectErr: nil,
			expectAccBank: &model.AccountBank{
				PubKey: suite.unreg.TransactionKey,
				Saving: suite.unregSaving,
			},
		},
		{
			testName:  "check public key from registered account",
			address:   sdk.AccAddress(suite.userWithoutBalance.TransactionKey.Address()),
			signKey:   suite.userWithoutBalance.TransactionKey,
			isPaid:    false,
			expectErr: nil,
			expectAccBank: &model.AccountBank{
				PubKey:   suite.userWithoutBalance.TransactionKey,
				Saving:   types.NewCoinFromInt64(0),
				Username: suite.userWithoutBalance.Username,
			},
		},
	}
	for _, tc := range testCases {
		err := suite.am.CheckSigningPubKeyOwnerByAddress(suite.Ctx, tc.address, tc.signKey, tc.isPaid)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)

		bank, _ := suite.am.storage.GetBank(suite.Ctx, tc.address)
		suite.Equal(tc.expectAccBank, bank, "%s", tc.testName)
	}
}

func (suite *AccountManagerTestSuite) TestCheckSigningPubKeyOwner() {
	txPrivKeys := []crypto.PrivKey{secp256k1.GenPrivKey(), secp256k1.GenPrivKey()}

	err := suite.am.AuthorizePermission(
		suite.Ctx, suite.userWithBalance.Username, suite.userWithoutBalance.Username,
		100, types.PreAuthorizationPermission,
		suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(1)))
	suite.Nil(err)

	testCases := []struct {
		testName     string
		username     types.AccountKey
		signKey      crypto.PubKey
		permission   types.Permission
		amount       types.Coin
		expectErr    sdk.Error
		expectSigner types.AccountKey
	}{
		{
			testName:     "account info doesn't exist",
			username:     suite.unreg.Username,
			signKey:      txPrivKeys[0].PubKey(),
			permission:   types.PreAuthorizationPermission,
			amount:       types.NewCoinFromInt64(1),
			expectErr:    model.ErrAccountInfoNotFound(),
			expectSigner: "",
		},
		{
			testName:     "public key mismatch",
			username:     suite.userWithBalance.Username,
			signKey:      txPrivKeys[0].PubKey(),
			permission:   types.PreAuthorizationPermission,
			amount:       types.NewCoinFromInt64(1),
			expectErr:    acctypes.ErrCheckAuthenticatePubKeyOwner(suite.userWithBalance.Username),
			expectSigner: "",
		},
		{
			testName:     "verify by signing key",
			username:     suite.userWithBalance.Username,
			signKey:      suite.userWithBalance.SigningKey,
			permission:   types.TransactionPermission,
			amount:       types.NewCoinFromInt64(1),
			expectErr:    nil,
			expectSigner: suite.userWithBalance.Username,
		},
		{
			testName:     "verify by transaction key",
			username:     suite.userWithBalance.Username,
			signKey:      suite.userWithBalance.SigningKey,
			permission:   types.ResetPermission,
			amount:       types.NewCoinFromInt64(1),
			expectErr:    nil,
			expectSigner: suite.userWithBalance.Username,
		},
		{
			testName:     "check preauth permission",
			username:     suite.userWithBalance.Username,
			signKey:      suite.userWithoutBalance.SigningKey,
			permission:   types.PreAuthorizationPermission,
			amount:       types.NewCoinFromInt64(1),
			expectErr:    nil,
			expectSigner: suite.userWithoutBalance.Username,
		},
		{
			testName:     "check app permission",
			username:     suite.userWithBalance.Username,
			signKey:      suite.userWithoutBalance.SigningKey,
			permission:   types.AppPermission,
			amount:       types.NewCoinFromInt64(1),
			expectErr:    acctypes.ErrCheckAuthenticatePubKeyOwner(suite.userWithBalance.Username),
			expectSigner: "",
		},
		{
			testName:   "check preauth amount is not enough",
			username:   suite.userWithBalance.Username,
			signKey:    suite.userWithoutBalance.SigningKey,
			permission: types.PreAuthorizationPermission,
			amount:     suite.userWithBalanceSaving,
			expectErr: acctypes.ErrPreAuthAmountInsufficient(
				suite.userWithoutBalance.Username,
				suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(2)),
				suite.userWithBalanceSaving),
			expectSigner: "",
		},
	}
	for _, tc := range testCases {
		signer, err := suite.am.CheckSigningPubKeyOwner(suite.Ctx, tc.username, tc.signKey, tc.permission, tc.amount)
		suite.Equal(tc.expectErr, err)
		suite.Equal(tc.expectSigner, signer)
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
				{
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
				{
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
				{
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
				{
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
				{
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
			err = am.IncreaseSequenceByOne(ctx, addr)
			if err != nil {
				panic(err)
			}
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

func (suite *AccountManagerTestSuite) TestRecoverAccount() {
	txPrivKeys := []crypto.PrivKey{secp256k1.GenPrivKey()}
	err := suite.am.AddFrozenMoney(suite.Ctx, suite.userWithBalance.Username, types.NewCoinFromInt64(1), 0, 100, 10)
	suite.Nil(err)
	testCases := []struct {
		testName         string
		username         types.AccountKey
		newTxPubKey      crypto.PubKey
		newSigningPubKey crypto.PubKey
		expectErr        sdk.Error
		oldAddr          sdk.AccAddress
		expectOldBank    *model.AccountBank
		expectNewBank    *model.AccountBank
		expectInfo       *model.AccountInfo
	}{
		{
			testName:         "username doesn't exist",
			username:         suite.unreg.Username,
			newTxPubKey:      secp256k1.GenPrivKey().PubKey(),
			newSigningPubKey: nil,
			expectErr:        model.ErrAccountInfoNotFound(),
			oldAddr:          sdk.AccAddress(suite.unreg.TransactionKey.Address()),
			expectOldBank: &model.AccountBank{
				Saving: suite.unregSaving,
			},
			expectNewBank: nil,
			expectInfo:    nil,
		},
		{
			testName:         "new bank linked to other account",
			username:         suite.userWithoutBalance.Username,
			newTxPubKey:      suite.userWithBalance.TransactionKey,
			newSigningPubKey: nil,
			expectErr: acctypes.ErrAddressAlreadyTaken(
				hex.EncodeToString(suite.userWithBalance.TransactionKey.Address())),
			oldAddr: sdk.AccAddress(suite.userWithoutBalance.TransactionKey.Address()),
			expectOldBank: &model.AccountBank{
				Username: suite.userWithoutBalance.Username,
				PubKey:   suite.userWithoutBalance.TransactionKey,
				Saving:   types.NewCoinFromInt64(0),
			},
			expectNewBank: &model.AccountBank{
				Username: suite.userWithBalance.Username,
				PubKey:   suite.userWithBalance.TransactionKey,
				Saving:   suite.userWithBalanceSaving,
				FrozenMoneyList: []model.FrozenMoney{
					{
						Amount:   types.NewCoinFromInt64(1),
						StartAt:  0,
						Interval: 100,
						Times:    10,
					},
				},
			},
			expectInfo: &suite.userWithoutBalance,
		},
		{
			testName:         "recover to empty address",
			username:         suite.userWithoutBalance.Username,
			newTxPubKey:      txPrivKeys[0].PubKey(),
			newSigningPubKey: nil,
			expectErr:        nil,
			oldAddr:          sdk.AccAddress(suite.userWithoutBalance.TransactionKey.Address()),
			expectOldBank: &model.AccountBank{
				PubKey: suite.userWithoutBalance.TransactionKey,
				Saving: types.NewCoinFromInt64(0),
			},
			expectNewBank: &model.AccountBank{
				Username: suite.userWithoutBalance.Username,
				PubKey:   txPrivKeys[0].PubKey(),
				Saving:   types.NewCoinFromInt64(0),
			},
			expectInfo: &model.AccountInfo{
				Username:       suite.userWithoutBalance.Username,
				TransactionKey: txPrivKeys[0].PubKey(),
				SigningKey:     nil,
				Address:        sdk.AccAddress(txPrivKeys[0].PubKey().Address()),
			},
		},
		{
			testName:         "recover to non empty address",
			username:         suite.userWithBalance.Username,
			newTxPubKey:      suite.unreg.TransactionKey,
			newSigningPubKey: nil,
			expectErr:        nil,
			oldAddr:          sdk.AccAddress(suite.userWithBalance.TransactionKey.Address()),
			expectOldBank: &model.AccountBank{
				PubKey: suite.userWithBalance.TransactionKey,
				Saving: types.NewCoinFromInt64(0),
			},
			expectNewBank: &model.AccountBank{
				Username: suite.userWithBalance.Username,
				PubKey:   suite.unreg.TransactionKey,
				Saving:   suite.unregSaving.Plus(suite.userWithBalanceSaving),
				FrozenMoneyList: []model.FrozenMoney{
					{
						Amount:   types.NewCoinFromInt64(1),
						StartAt:  0,
						Interval: 100,
						Times:    10,
					},
				},
			},
			expectInfo: &model.AccountInfo{
				Username:       suite.userWithBalance.Username,
				TransactionKey: suite.unreg.TransactionKey,
				SigningKey:     nil,
				Address:        sdk.AccAddress(suite.unreg.TransactionKey.Address()),
			},
		},
	}
	for _, tc := range testCases {
		err := suite.am.RecoverAccount(suite.Ctx, tc.username, tc.newTxPubKey, tc.newSigningPubKey)
		fmt.Println(err)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		oldBank, _ := suite.am.GetBankByAddress(suite.Ctx, tc.oldAddr)
		suite.Equal(tc.expectOldBank, oldBank, "%s", tc.testName)
		newBank, _ := suite.am.GetBankByAddress(suite.Ctx, sdk.AccAddress(tc.newTxPubKey.Address()))
		suite.Equal(tc.expectNewBank, newBank, "%s", tc.testName)
		info, _ := suite.am.GetInfo(suite.Ctx, tc.username)
		suite.Equal(tc.expectInfo, info, "%s", tc.testName)
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
