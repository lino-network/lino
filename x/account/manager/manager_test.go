package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	"github.com/lino-network/lino/types"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"
	acctypes "github.com/lino-network/lino/x/account/types"
)

var (
	storeKeyStr = "testAccountStore"
	kvStoreKey  = sdk.NewKVStoreKey(storeKeyStr)
)

type AccountStoreDumper struct{}

func (dumper AccountStoreDumper) NewDumper() *testutils.Dumper {
	return model.NewAccountDumper(model.NewAccountStorage(kvStoreKey))
}

type AccountManagerTestSuite struct {
	testsuites.GoldenTestSuite
	am AccountManager
	ph *param.ParamKeeper

	// mock data
	userWithoutBalance model.AccountInfo

	userWithBalance       model.AccountInfo
	userWithBalanceSaving types.Coin
	registerFee           types.Coin

	unreg model.AccountInfo

	unregSaving types.Coin
}

func TestAccountManagerTestSuite(t *testing.T) {
	suite.Run(t, &AccountManagerTestSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(AccountStoreDumper{}, kvStoreKey),
	})
}

func (suite *AccountManagerTestSuite) SetupTest() {
	suite.SetupCtx(0, time.Unix(0, 0), kvStoreKey)
	suite.ph = &param.ParamKeeper{}
	suite.am = NewAccountManager(kvStoreKey, suite.ph)

	// background
	suite.userWithoutBalance = model.AccountInfo{
		Username:       types.AccountKey("userwithoutbalance"),
		SigningKey:     sampleKeys()[0],
		TransactionKey: sampleKeys()[1],
	}
	suite.userWithoutBalance.Address = sdk.AccAddress(suite.userWithoutBalance.TransactionKey.Address())

	suite.userWithBalance = model.AccountInfo{
		Username:       types.AccountKey("userwithbalance"),
		SigningKey:     sampleKeys()[2],
		TransactionKey: sampleKeys()[3],
	}
	suite.userWithBalance.Address = sdk.AccAddress(suite.userWithBalance.TransactionKey.Address())

	suite.unreg = model.AccountInfo{
		Username:       types.AccountKey("unreg"),
		SigningKey:     sampleKeys()[4],
		TransactionKey: sampleKeys()[5],
	}
	suite.unreg.Address = sdk.AccAddress(suite.unreg.TransactionKey.Address())

	suite.userWithBalanceSaving = types.NewCoinFromInt64(1000 * types.Decimals)
	suite.unregSaving = types.NewCoinFromInt64(1 * types.Decimals)
	suite.registerFee = types.NewCoinFromInt64(100 * types.Decimals)

	err := suite.am.GenesisAccount(suite.Ctx, suite.userWithoutBalance.Username, suite.userWithoutBalance.SigningKey, suite.userWithoutBalance.TransactionKey)
	suite.NoError(err)

	err = suite.am.GenesisAccount(suite.Ctx, suite.userWithBalance.Username, suite.userWithBalance.SigningKey, suite.userWithBalance.TransactionKey)
	suite.NoError(err)
	err = suite.am.addCoinToUsername(suite.Ctx, suite.userWithBalance.Username, suite.userWithBalanceSaving)
	suite.NoError(err)

	suite.am.addCoinToAddress(suite.Ctx, sdk.AccAddress(suite.unreg.TransactionKey.Address()), suite.unregSaving)

	suite.ph.On("GetAccountParam", mock.Anything).Return(&parammodel.AccountParam{
		RegisterFee:       suite.registerFee,
		MinimumBalance:    types.NewCoinFromInt64(0),
		MaxNumFrozenMoney: 10,
	}, nil).Maybe()
}

func (suite *AccountManagerTestSuite) TestInitGenesis() {
	suite.NextBlock(time.Unix(123, 0))
	am := suite.am
	ctx := suite.Ctx

	total := linotypes.NewCoinFromInt64(2000000)

	am.InitGenesis(ctx, total, []model.Pool{
		{
			Name:    linotypes.InflationValidatorPool,
			Balance: linotypes.NewCoinFromInt64(123),
		},
		{
			Name:    linotypes.AccountVestingPool,
			Balance: linotypes.NewCoinFromInt64(1000000),
		},
	})

	supply := am.GetSupply(ctx)
	suite.Equal(model.Supply{
		LastYearTotal:     total,
		Total:             total,
		ChainStartTime:    ctx.BlockTime().Unix(),
		LastInflationTime: ctx.BlockTime().Unix(),
	}, supply)

	pool1, err := am.GetPool(ctx, linotypes.InflationValidatorPool)
	suite.Nil(err)
	suite.Equal(linotypes.NewCoinFromInt64(123), pool1)

	pool2, err := am.GetPool(ctx, linotypes.AccountVestingPool)
	suite.Nil(err)
	suite.Equal(linotypes.NewCoinFromInt64(1000000), pool2)

	_, err = am.GetPool(ctx, "not-a-pool")
	suite.NotNil(err)

	suite.Panics(func() {
		am.InitGenesis(ctx, total, nil)
	})

	suite.Golden()
}

func (suite *AccountManagerTestSuite) TestMoveFromPools() {
	initBackground := func() {
		suite.NextBlock(time.Unix(123, 0))
		am := suite.am
		ctx := suite.Ctx

		total := linotypes.NewCoinFromInt64(2000000)

		am.InitGenesis(ctx, total, []model.Pool{
			{
				Name:    linotypes.InflationValidatorPool,
				Balance: linotypes.NewCoinFromInt64(123),
			},
			{
				Name:    linotypes.AccountVestingPool,
				Balance: linotypes.NewCoinFromInt64(1000000),
			},
		})
	}
	cases := []struct {
		name             string
		pool             linotypes.PoolName
		to               linotypes.AccOrAddr
		amount           linotypes.Coin
		expectedErr      sdk.Error
		expectedBalance  linotypes.Coin
		expectedPoolLeft linotypes.Coin
	}{
		{
			name:        "move negative amount",
			pool:        linotypes.AccountVestingPool,
			to:          linotypes.NewAccOrAddrFromAcc(types.AccountKey("userwithoutbalance")),
			amount:      linotypes.NewCoinFromInt64(-1),
			expectedErr: acctypes.ErrNegativeMoveAmount(linotypes.NewCoinFromInt64(-1)),
		},
		{
			name:        "pool not enough",
			pool:        linotypes.InflationValidatorPool,
			to:          linotypes.NewAccOrAddrFromAcc(types.AccountKey("userwithoutbalance")),
			amount:      linotypes.NewCoinFromInt64(124),
			expectedErr: acctypes.ErrPoolNotEnough(linotypes.InflationValidatorPool),
		},
		{
			name:        "pool not exists",
			pool:        "poolnotexists",
			to:          linotypes.NewAccOrAddrFromAcc(types.AccountKey("userwithoutbalance")),
			amount:      linotypes.NewCoinFromInt64(124),
			expectedErr: acctypes.ErrPoolNotFound("poolnotexists"),
		},
		{
			name:             "succ move to account",
			pool:             linotypes.InflationValidatorPool,
			to:               linotypes.NewAccOrAddrFromAcc(types.AccountKey("userwithoutbalance")),
			amount:           linotypes.NewCoinFromInt64(100),
			expectedBalance:  linotypes.NewCoinFromInt64(100),
			expectedPoolLeft: linotypes.NewCoinFromInt64(23),
		},
		{
			name:             "succ move to addr",
			pool:             linotypes.AccountVestingPool,
			to:               linotypes.NewAccOrAddrFromAddr(suite.userWithoutBalance.Address),
			amount:           linotypes.NewCoinFromInt64(1000000),
			expectedBalance:  linotypes.NewCoinFromInt64(1000000),
			expectedPoolLeft: linotypes.NewCoinFromInt64(0),
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			initBackground()
			err := suite.am.MoveFromPool(suite.Ctx, tc.pool, tc.to, tc.amount)
			suite.Equal(tc.expectedErr, err)
			if tc.expectedErr == nil {
				if tc.to.IsAddr {
					bank, err := suite.am.GetBankByAddress(suite.Ctx, tc.to.Addr)
					suite.Nil(err)
					suite.Equal(tc.expectedBalance, bank.Saving)
				} else {
					bank, err := suite.am.GetBank(suite.Ctx, tc.to.AccountKey)
					suite.Nil(err)
					suite.Equal(tc.expectedBalance, bank.Saving)
				}
				pool, err := suite.am.GetPool(suite.Ctx, tc.pool)
				suite.Nil(err)
				suite.Equal(tc.expectedPoolLeft, pool)
			}
			suite.Golden()
		})
	}
}

func (suite *AccountManagerTestSuite) TestMoveToPools() {
	initBackground := func() {
		suite.NextBlock(time.Unix(123, 0))
		am := suite.am
		ctx := suite.Ctx

		total := linotypes.NewCoinFromInt64(2000000)

		am.InitGenesis(ctx, total, []model.Pool{
			{
				Name:    linotypes.InflationValidatorPool,
				Balance: linotypes.NewCoinFromInt64(123),
			},
			{
				Name:    linotypes.AccountVestingPool,
				Balance: linotypes.NewCoinFromInt64(1000000),
			},
		})
	}
	cases := []struct {
		name             string
		pool             linotypes.PoolName
		from             linotypes.AccOrAddr
		amount           linotypes.Coin
		expectedErr      sdk.Error
		expectedBalance  linotypes.Coin
		expectedPoolLeft linotypes.Coin
	}{
		{
			name:        "move negative amount",
			pool:        linotypes.AccountVestingPool,
			from:        linotypes.NewAccOrAddrFromAcc(types.AccountKey("userwithbalance")),
			amount:      linotypes.NewCoinFromInt64(-1),
			expectedErr: acctypes.ErrNegativeMoveAmount(linotypes.NewCoinFromInt64(-1)),
		},
		{
			name:        "balance not enough",
			pool:        linotypes.InflationValidatorPool,
			from:        linotypes.NewAccOrAddrFromAcc(suite.userWithBalance.Username),
			amount:      suite.userWithBalanceSaving.Plus(linotypes.NewCoinFromInt64(1)),
			expectedErr: acctypes.ErrAccountSavingCoinNotEnough(),
		},
		{
			name:        "pool not exists",
			pool:        "poolnotexists",
			from:        linotypes.NewAccOrAddrFromAcc(types.AccountKey("userwithbalance")),
			amount:      linotypes.NewCoinFromInt64(1),
			expectedErr: acctypes.ErrPoolNotFound("poolnotexists"),
		},
		{
			name:             "succ move from account",
			pool:             linotypes.InflationValidatorPool,
			from:             linotypes.NewAccOrAddrFromAcc(types.AccountKey("userwithbalance")),
			amount:           suite.userWithBalanceSaving,
			expectedBalance:  linotypes.NewCoinFromInt64(0),
			expectedPoolLeft: suite.userWithBalanceSaving.Plus(linotypes.NewCoinFromInt64(123)),
		},
		{
			name:             "succ move from addr",
			pool:             linotypes.AccountVestingPool,
			from:             linotypes.NewAccOrAddrFromAddr(suite.userWithBalance.Address),
			amount:           suite.userWithBalanceSaving.Minus(linotypes.NewCoinFromInt64(1)),
			expectedBalance:  linotypes.NewCoinFromInt64(1),
			expectedPoolLeft: linotypes.NewCoinFromInt64(1000000).Plus(suite.userWithBalanceSaving.Minus(linotypes.NewCoinFromInt64(1))),
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			initBackground()
			err := suite.am.MoveToPool(suite.Ctx, tc.pool, tc.from, tc.amount)
			suite.Equal(tc.expectedErr, err)
			if tc.expectedErr == nil {
				if tc.from.IsAddr {
					bank, err := suite.am.GetBankByAddress(suite.Ctx, tc.from.Addr)
					suite.Nil(err)
					suite.Equal(tc.expectedBalance, bank.Saving)
				} else {
					bank, err := suite.am.GetBank(suite.Ctx, tc.from.AccountKey)
					suite.Nil(err)
					suite.Equal(tc.expectedBalance, bank.Saving)
				}
				pool, err := suite.am.GetPool(suite.Ctx, tc.pool)
				suite.Nil(err)
				suite.Equal(tc.expectedPoolLeft, pool)
			}
			suite.Golden()
		})
	}
}

func (suite *AccountManagerTestSuite) TestBetweenPools() {
	initBackground := func() {
		suite.NextBlock(time.Unix(123, 0))
		am := suite.am
		ctx := suite.Ctx

		total := linotypes.NewCoinFromInt64(2000000)

		am.InitGenesis(ctx, total, []model.Pool{
			{
				Name:    linotypes.InflationValidatorPool,
				Balance: linotypes.NewCoinFromInt64(123),
			},
			{
				Name:    linotypes.AccountVestingPool,
				Balance: linotypes.NewCoinFromInt64(1000000),
			},
		})
	}
	cases := []struct {
		name         string
		from         linotypes.PoolName
		to           linotypes.PoolName
		amount       linotypes.Coin
		expectedErr  sdk.Error
		expectedFrom linotypes.Coin
		expectedTo   linotypes.Coin
	}{
		{
			name:        "move negative amount",
			from:        linotypes.AccountVestingPool,
			to:          linotypes.AccountVestingPool,
			amount:      linotypes.NewCoinFromInt64(-1),
			expectedErr: acctypes.ErrNegativeMoveAmount(linotypes.NewCoinFromInt64(-1)),
		},
		{
			name:        "from pool not exists",
			from:        "poolnotexists",
			to:          linotypes.AccountVestingPool,
			amount:      linotypes.NewCoinFromInt64(1),
			expectedErr: acctypes.ErrPoolNotFound("poolnotexists"),
		},
		{
			name:        "to pool not exists",
			from:        linotypes.AccountVestingPool,
			to:          "poolnotexists",
			amount:      linotypes.NewCoinFromInt64(1),
			expectedErr: acctypes.ErrPoolNotFound("poolnotexists"),
		},
		{
			name:        "balance not enough",
			from:        linotypes.InflationValidatorPool,
			to:          linotypes.AccountVestingPool,
			amount:      linotypes.NewCoinFromInt64(124),
			expectedErr: acctypes.ErrPoolNotEnough(linotypes.InflationValidatorPool),
		},
		{
			name:         "succ",
			from:         linotypes.InflationValidatorPool,
			to:           linotypes.AccountVestingPool,
			amount:       linotypes.NewCoinFromInt64(1),
			expectedFrom: linotypes.NewCoinFromInt64(122),
			expectedTo:   linotypes.NewCoinFromInt64(1000001),
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			initBackground()
			err := suite.am.MoveBetweenPools(suite.Ctx, tc.from, tc.to, tc.amount)
			suite.Equal(tc.expectedErr, err)
			if tc.expectedErr == nil {
				poolFrom, _ := suite.am.GetPool(suite.Ctx, tc.from)
				suite.Equal(tc.expectedFrom, poolFrom)
				poolTo, _ := suite.am.GetPool(suite.Ctx, tc.to)
				suite.Equal(tc.expectedTo, poolTo)
			}
			suite.Golden()
		})
	}
}

// test mint schedule
func (suite *AccountManagerTestSuite) TestMint() {
	// Genesis
	init := int64(123)
	suite.NextBlock(time.Unix(init, 0))
	am := suite.am

	total := linotypes.MustLinoToCoin("10000000000")
	am.InitGenesis(suite.Ctx, total, []model.Pool{
		{
			Name:    linotypes.InflationValidatorPool,
			Balance: linotypes.NewCoinFromInt64(0),
		},
		{
			Name:    linotypes.InflationDeveloperPool,
			Balance: linotypes.NewCoinFromInt64(0),
		},
		{
			Name:    linotypes.InflationConsumptionPool,
			Balance: linotypes.NewCoinFromInt64(0),
		},
		{
			Name:    linotypes.AccountVestingPool,
			Balance: total,
		},
	})

	// param
	rate := sdk.MustNewDecFromStr("0.065")
	cc := sdk.MustNewDecFromStr("0.10")
	dev := sdk.MustNewDecFromStr("0.75")
	val := sdk.MustNewDecFromStr("0.15")
	suite.ph.On("GetGlobalAllocationParam", mock.Anything).Return(
		&parammodel.GlobalAllocationParam{
			GlobalGrowthRate:         rate,
			ContentCreatorAllocation: cc,
			DeveloperAllocation:      dev,
			ValidatorAllocation:      val,
		})

	computeHourly := func(total linotypes.Coin, growth sdk.Dec) (linotypes.Coin, linotypes.Coin, linotypes.Coin) {
		amount := linotypes.DecToCoin(total.ToDec().Mul(growth).Mul(
			linotypes.NewDecFromRat(1, nHourOfOneYear)))
		ccAmount := linotypes.DecToCoin(amount.ToDec().Mul(cc))
		valAmount := linotypes.DecToCoin(amount.ToDec().Mul(val))
		devAmount := amount.Minus(ccAmount).Minus(valAmount)
		return ccAmount, valAmount, devAmount
	}

	getPools := func(ctx sdk.Context) (linotypes.Coin, linotypes.Coin, linotypes.Coin) {
		cpool, _ := suite.am.GetPool(ctx, linotypes.InflationConsumptionPool)
		vpool, _ := suite.am.GetPool(ctx, linotypes.InflationValidatorPool)
		dpool, _ := suite.am.GetPool(ctx, linotypes.InflationDeveloperPool)
		return cpool, vpool, dpool
	}

	checkPool := func(ctx sdk.Context, cc, val, dev linotypes.Coin) {
		cpool, vpool, dpool := getPools(ctx)
		suite.Equal(cc, cpool)
		suite.Equal(val, vpool)
		suite.Equal(dev, dpool)
	}

	// test first hour
	firstYearOneHourCC := linotypes.MustLinoToCoin("7415.01255")
	firstYearOneHourVal := linotypes.MustLinoToCoin("11122.51882")
	firstYearOneHourDev := linotypes.MustLinoToCoin("55612.59411")
	base := total
	t := init + nSecOfOneHour
	suite.NextBlock(time.Unix(t, 0))
	err := am.Mint(suite.Ctx)
	suite.Nil(err)
	checkPool(suite.Ctx, firstYearOneHourCC, firstYearOneHourVal, firstYearOneHourDev)

	// same time again, won't mint
	err = am.Mint(suite.Ctx)
	suite.Nil(err)
	checkPool(suite.Ctx, firstYearOneHourCC, firstYearOneHourVal, firstYearOneHourDev)

	// test first 50 hours
	for ; t <= init+50*nSecOfOneHour; t += 5 {
		suite.NextBlock(time.Unix(t, 0))
		err := am.Mint(suite.Ctx)
		suite.Nil(err)
		if (t-init)%nSecOfOneHour == 0 {
			n := (t - init) / nSecOfOneHour
			checkPool(suite.Ctx,
				linotypes.DecToCoin(firstYearOneHourCC.ToDec().Mul(sdk.NewDec(n))),
				linotypes.DecToCoin(firstYearOneHourVal.ToDec().Mul(sdk.NewDec(n))),
				linotypes.DecToCoin(firstYearOneHourDev.ToDec().Mul(sdk.NewDec(n))),
			)
		}
	}

	// first year, 123 + nSecOfOneHour * nHourOfOneYear, is the first year
	// math check
	hourSum := firstYearOneHourCC.Plus(firstYearOneHourVal).Plus(firstYearOneHourDev)
	oneYearComputedMint := linotypes.DecToCoin(hourSum.ToDec().Mul(sdk.NewDec(nHourOfOneYear)))
	oneYearTotal := linotypes.MustLinoToCoin("10649999999.95768")
	suite.Equal(oneYearTotal.Minus(base), oneYearComputedMint)

	t = init + nSecOfOneHour*nHourOfOneYear
	suite.NextBlock(time.Unix(t, 0))
	err = suite.am.Mint(suite.Ctx)
	suite.Nil(err)
	supply := suite.am.GetSupply(suite.Ctx)
	suite.Equal(oneYearTotal, supply.LastYearTotal)
	suite.Equal(oneYearTotal, supply.Total)
	// cc1, val1, dev1 := getPools(suite.Ctx)
	// fmt.Println(cc1, val1, dev1)
	// fmt.Println(cc1.Plus(val1).Plus(dev1))

	// second year, 123 + 2 * (nSecOfOneHour * nHourOfOneYear), is the second year
	// next year first hour
	base = oneYearTotal
	lastYearCC, lastYearVal, lastYearDev := getPools(suite.Ctx)
	ccAmount, valAmount, devAmount := computeHourly(base, rate)
	t += nSecOfOneHour
	suite.NextBlock(time.Unix(t, 0))
	err = suite.am.Mint(suite.Ctx)
	suite.Nil(err)
	checkPool(suite.Ctx,
		lastYearCC.Plus(ccAmount),
		lastYearVal.Plus(valAmount),
		lastYearDev.Plus(devAmount),
	)
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
		suite.am.addCoinToAddress(suite.Ctx, tc.address, tc.amount)
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
		err := suite.am.addCoinToUsername(suite.Ctx, tc.username, tc.amount)
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
			expectErr:  acctypes.ErrAccountBankNotFound(emptyAddress),
			amount:     coin1,
			expectBank: nil,
		},
	}
	for _, tc := range testCases {
		err := suite.am.minusCoinFromAddress(suite.Ctx, tc.address, tc.amount)
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
		err := suite.am.minusCoinFromUsername(suite.Ctx, tc.username, tc.amount)
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
			expectErr:  acctypes.ErrAddressAlreadyTaken(sdk.AccAddress(userWithBalance.TransactionKey.Address()).String()),
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
		err := suite.am.GenesisAccount(suite.Ctx, tc.username, tc.signingKey, tc.txKey)
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
	ctx, am := setupTest(t, 1)

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
	suite.am.storage.SetPool(suite.Ctx, &model.Pool{
		Name:    types.InflationValidatorPool,
		Balance: types.MustLinoToCoin("10000000000"),
	})

	txPrivKeys := []crypto.PrivKey{secp256k1.GenPrivKey(), secp256k1.GenPrivKey()}
	signingPrivKeys := []crypto.PrivKey{secp256k1.GenPrivKey(), secp256k1.GenPrivKey()}

	testCases := []struct {
		testName    string
		referrer    types.AccOrAddr
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
			referrer:    types.NewAccOrAddrFromAcc(suite.userWithBalance.Username),
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
			referrer:    types.NewAccOrAddrFromAcc(suite.userWithBalance.Username),
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
			referrer:    types.NewAccOrAddrFromAcc(suite.userWithBalance.Username),
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
			referrer:    types.NewAccOrAddrFromAcc(suite.userWithBalance.Username),
			registerFee: suite.registerFee,
			username:    "test2",
			signingKey:  signingPrivKeys[0].PubKey(),
			txKey:       txPrivKeys[0].PubKey(),
			expectErr: acctypes.ErrAddressAlreadyTaken(
				sdk.AccAddress(txPrivKeys[0].PubKey().Address()).String()),
			accInfo: nil,
			accBank: nil,
		},
		{
			testName: "referrer is address",
			referrer: types.NewAccOrAddrFromAddr(
				sdk.AccAddress(suite.userWithBalance.TransactionKey.Address())),
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

func (suite *AccountManagerTestSuite) TestMoveCoinAccOrAddr() {
	testCases := []struct {
		testName              string
		sender                types.AccOrAddr
		amount                types.Coin
		receiver              types.AccOrAddr
		expectErr             sdk.Error
		expectSenderBalance   types.Coin
		expectReceiverBalance types.Coin
	}{
		{
			testName:              "negative amount",
			sender:                types.NewAccOrAddrFromAcc("movecointest"),
			receiver:              types.NewAccOrAddrFromAcc(suite.userWithoutBalance.Username),
			amount:                types.NewCoinFromInt64(-1),
			expectErr:             acctypes.ErrNegativeMoveAmount(types.NewCoinFromInt64(-1)),
			expectSenderBalance:   types.Coin{},
			expectReceiverBalance: types.NewCoinFromInt64(0),
		},
		{
			testName:              "sender doesnt exist",
			sender:                types.NewAccOrAddrFromAcc("movecointest"),
			receiver:              types.NewAccOrAddrFromAcc(suite.userWithoutBalance.Username),
			amount:                types.NewCoinFromInt64(1),
			expectErr:             acctypes.ErrAccountNotFound("movecointest"),
			expectSenderBalance:   types.Coin{},
			expectReceiverBalance: types.NewCoinFromInt64(0),
		},
		{
			testName:              "receiver doesnt exist",
			sender:                types.NewAccOrAddrFromAcc(suite.userWithBalance.Username),
			receiver:              types.NewAccOrAddrFromAcc("movecointest"),
			amount:                types.NewCoinFromInt64(1),
			expectErr:             acctypes.ErrAccountNotFound("movecointest"),
			expectSenderBalance:   suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(1)),
			expectReceiverBalance: types.Coin{},
		},
		{
			testName:              "send from username to username",
			sender:                types.NewAccOrAddrFromAcc(suite.userWithBalance.Username),
			receiver:              types.NewAccOrAddrFromAcc(suite.userWithoutBalance.Username),
			amount:                types.NewCoinFromInt64(1),
			expectErr:             nil,
			expectSenderBalance:   suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(2)),
			expectReceiverBalance: types.NewCoinFromInt64(1),
		},
		{
			testName: "send from username to address",
			sender:   types.NewAccOrAddrFromAcc(suite.userWithBalance.Username),
			receiver: types.NewAccOrAddrFromAddr(
				sdk.AccAddress(suite.userWithoutBalance.TransactionKey.Address())),
			amount:                types.NewCoinFromInt64(1),
			expectErr:             nil,
			expectSenderBalance:   suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(3)),
			expectReceiverBalance: types.NewCoinFromInt64(2),
		},
		{
			testName: "send from address to address",
			sender: types.NewAccOrAddrFromAddr(
				sdk.AccAddress(suite.userWithBalance.TransactionKey.Address())),
			receiver: types.NewAccOrAddrFromAddr(
				sdk.AccAddress(suite.userWithoutBalance.TransactionKey.Address())),
			amount:                types.NewCoinFromInt64(1),
			expectErr:             nil,
			expectSenderBalance:   suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(4)),
			expectReceiverBalance: types.NewCoinFromInt64(3),
		},
		{
			testName: "send from address to user",
			sender: types.NewAccOrAddrFromAddr(
				sdk.AccAddress(suite.userWithBalance.TransactionKey.Address())),
			receiver:              types.NewAccOrAddrFromAcc(suite.userWithoutBalance.Username),
			amount:                types.NewCoinFromInt64(1),
			expectErr:             nil,
			expectSenderBalance:   suite.userWithBalanceSaving.Minus(types.NewCoinFromInt64(5)),
			expectReceiverBalance: types.NewCoinFromInt64(4),
		},
	}
	for _, tc := range testCases {
		err := suite.am.MoveCoin(suite.Ctx, tc.sender, tc.receiver, tc.amount)
		suite.Equal(tc.expectErr, err)
		if !tc.sender.IsAddr {
			saving, _ := suite.am.GetSavingFromUsername(suite.Ctx, tc.sender.AccountKey)
			suite.Equal(tc.expectSenderBalance, saving)
		} else {
			saving, _ := suite.am.GetSavingFromAddress(suite.Ctx, tc.sender.Addr)
			suite.Equal(tc.expectSenderBalance, saving)
		}
		if !tc.receiver.IsAddr {
			saving, _ := suite.am.GetSavingFromUsername(suite.Ctx, tc.receiver.AccountKey)
			suite.Equal(tc.expectReceiverBalance, saving)
		} else {
			saving, _ := suite.am.GetSavingFromAddress(suite.Ctx, tc.receiver.Addr)
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
			testName: "bank doesn't exist",
			address:  sdk.AccAddress(txPrivKeys[0].PubKey().Address()),
			signKey:  txPrivKeys[0].PubKey(),
			isPaid:   false,
			expectErr: acctypes.ErrAccountBankNotFound(
				sdk.AccAddress(txPrivKeys[0].PubKey().Address())),
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
			expectErr:    acctypes.ErrAccountNotFound(suite.unreg.Username),
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

	ctx, am := setupTest(t, 1)
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
			expectResult: acctypes.ErrGrantPubKeyNotFound(),
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
	ctx, am := setupTest(t, 1)
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
	ctx, am := setupTest(t, 1)
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
	ctx, am := setupTest(t, 1)
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
			expectErr:        acctypes.ErrAccountNotFound(suite.unreg.Username),
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
				sdk.AccAddress(suite.userWithBalance.TransactionKey.Address()).String()),
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

func (suite *AccountManagerTestSuite) TestImportExport() {
	// background data
	suite.NextBlock(time.Unix(123, 0))
	am := suite.am
	ctx := suite.Ctx
	total := linotypes.NewCoinFromInt64(2000000)
	am.InitGenesis(ctx, total, []model.Pool{
		{
			Name:    linotypes.InflationValidatorPool,
			Balance: linotypes.NewCoinFromInt64(123),
		},
		{
			Name:    linotypes.AccountVestingPool,
			Balance: linotypes.NewCoinFromInt64(1000000),
		},
	})
	err := am.UpdateJSONMeta(ctx, suite.userWithoutBalance.Username, `{"key":"value"}`)
	suite.Nil(err)

	cdc := wire.New()
	wire.RegisterCrypto(cdc)

	dir, err2 := ioutil.TempDir("", "test")
	suite.Require().Nil(err2)
	defer os.RemoveAll(dir) // clean up

	tmpfn := filepath.Join(dir, "tmpfile")
	err2 = suite.am.ExportToFile(suite.Ctx, cdc, tmpfn)
	suite.Nil(err2)

	// reset state
	suite.SetupCtx(0, time.Unix(0, 0), kvStoreKey)
	suite.ph = &param.ParamKeeper{}
	suite.am = NewAccountManager(kvStoreKey, suite.ph)
	err2 = suite.am.ImportFromFile(suite.Ctx, cdc, tmpfn)
	suite.Nil(err2)

	suite.Golden()
}

// cdc := wire.New()
// wire.RegisterCrypto(cdc)
// keys := make([]crypto.PubKey, 0)
// for i := 0 ; i < 10; i++ {
// keys = append(keys, secp256k1.GenPrivKey().PubKey())
// }
// fmt.Print(string(cdc.MustMarshalJSON(keys)))
func sampleKeys() []crypto.PubKey {
	json := `
[{"type":"tendermint/PubKeySecp256k1","value":"Aot3u5m7vuxUOszkS6IZW5XYVu6ATvZsfSQIjtQo9tML"},{"type":"tendermint/PubKeySecp256k1","value":"AoFqbXKmblwKVggqb8Cqo30gRKs9EfqwhOhuyOKlGCuD"},{"type":"tendermint/PubKeySecp256k1","value":"Aj/1EOLKUKUPhp+mx3fLNoZOEEsY+tjPeTW4nOPbqwwq"},{"type":"tendermint/PubKeySecp256k1","value":"A1SxTVyDiXljmHeimniCQiNZQ3dcDsgppP0gDCMgJtdp"},{"type":"tendermint/PubKeySecp256k1","value":"Ax8b6HzTh9el9/NfE8fI4awCvMZWGQkjl+rYOGWeGJc9"},{"type":"tendermint/PubKeySecp256k1","value":"A4r+RjYEc2V9p43J4CovoktRTXNY9vvcQbx0aOW9bhoq"},{"type":"tendermint/PubKeySecp256k1","value":"AwFSpofxlQGAQv167WveHyeUvTh/3fukkJU7gkEW+iMm"},{"type":"tendermint/PubKeySecp256k1","value":"AjglddkWGGlMZck7uvWMDCtyqpNWSBy9HmnJV9vPnu2k"},{"type":"tendermint/PubKeySecp256k1","value":"A+KW7obJ0BpKqUWmY33svTBxGdTfRhmOym7A5imWWwGm"},{"type":"tendermint/PubKeySecp256k1","value":"A6P8IUdt9DKrYCe3/Tflt7DBdgFokRcCKkixt+UbhjZ8"}]
`

	keys := make([]crypto.PubKey, 0)
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	cdc.MustUnmarshalJSON([]byte(json), &keys)
	return keys
}
