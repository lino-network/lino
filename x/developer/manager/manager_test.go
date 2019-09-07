package developer // To test private filed `storage`

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	mparam "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
	maccount "github.com/lino-network/lino/x/account/mocks"
	"github.com/lino-network/lino/x/developer/model"
	"github.com/lino-network/lino/x/developer/types"
	mglobal "github.com/lino-network/lino/x/global/mocks"
	mprice "github.com/lino-network/lino/x/price/mocks"
	mvote "github.com/lino-network/lino/x/vote/mocks"
	votetypes "github.com/lino-network/lino/x/vote/types"
)

var (
	storeKeyStr      = "testStoreKey"
	storeKey         = sdk.NewKVStoreKey(storeKeyStr)
	appTest          = linotypes.AccountKey("testapp")
	appDoesNotExists = linotypes.AccountKey("testapp-does-not-exist")
	appWithoutIDA    = linotypes.AccountKey("testapp-without-ida")
	appHasRevokedIDA = linotypes.AccountKey("testapp-revoked-ida")
)

type DeveloperDumper struct{}

func (dumper DeveloperDumper) NewDumper() *testutils.Dumper {
	return model.NewDeveloperDumper(model.NewDeveloperStorage(storeKey))
}

type DeveloperManagerSuite struct {
	testsuites.GoldenTestSuite
	manager        DeveloperManager
	mParamKeeper   *mparam.ParamKeeper
	mVoteKeeper    *mvote.VoteKeeper
	mAccountKeeper *maccount.AccountKeeper
	mPriceKeeper   *mprice.PriceKeeper
	mGlobalKeeper  *mglobal.GlobalKeeper
}

func NewDeveloperManagerSuite() *DeveloperManagerSuite {
	return &DeveloperManagerSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(DeveloperDumper{}, storeKey),
	}
}

func (suite *DeveloperManagerSuite) SetupTest() {
	suite.mParamKeeper = new(mparam.ParamKeeper)
	suite.mVoteKeeper = new(mvote.VoteKeeper)
	suite.mAccountKeeper = new(maccount.AccountKeeper)
	suite.mPriceKeeper = new(mprice.PriceKeeper)
	suite.mGlobalKeeper = new(mglobal.GlobalKeeper)
	suite.manager = NewDeveloperManager(storeKey, suite.mParamKeeper, suite.mVoteKeeper, suite.mAccountKeeper, suite.mPriceKeeper, suite.mGlobalKeeper)
	suite.SetupCtx(0, time.Unix(0, 0), storeKey)
}

func TestDeveloperManagerSuite(t *testing.T) {
	suite.Run(t, NewDeveloperManagerSuite())
}

func (suite *DeveloperManagerSuite) TestInitGenesis() {
	testCases := []struct {
		reservePoolAmount linotypes.Coin
		name              string
		expected          sdk.Error
		expectedStore     *model.ReservePool
	}{
		{
			name:              "Success Valid Genesis",
			reservePoolAmount: linotypes.NewCoin(sdk.NewInt(1)),
			expected:          nil,
			expectedStore:     &model.ReservePool{Total: linotypes.NewCoin(sdk.NewInt(1))},
		},
		{
			name:              "Fail Invalid Genesis negative reservePoolAmount",
			reservePoolAmount: linotypes.NewCoin(sdk.NewInt(-1)),
			expected:          types.ErrInvalidReserveAmount(linotypes.NewCoin(sdk.NewInt(-1))),
			expectedStore:     nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.SetupCtx(0, time.Unix(0, 0), storeKey)
			suite.Equal(c.expected, suite.manager.InitGenesis(suite.Ctx, c.reservePoolAmount))
			suite.Golden()
		})
	}
}

func (suite *DeveloperManagerSuite) TestDoesDeveloperExist() {
	devDeleted := linotypes.AccountKey("testapp-deleted")
	testCases := []struct {
		name     string
		username linotypes.AccountKey
		expected bool
	}{
		{
			name:     "Developer exists not deleted",
			username: appTest,
			expected: true,
		},
		{
			name:     "Developer exist deleted",
			username: devDeleted,
			expected: false,
		},
		{
			name:     "Developer does not exist",
			username: appDoesNotExists,
			expected: false,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false)
			suite.Equal(c.expected, suite.manager.DoesDeveloperExist(suite.Ctx, c.username))
			suite.Golden()
			suite.AssertStateUnchanged(false)
		})
	}
}

func (suite *DeveloperManagerSuite) TestGetDeveloper() {
	testCases := []struct {
		name        string
		username    linotypes.AccountKey
		expected    model.Developer
		expectedErr sdk.Error
	}{
		{
			name:     "Developer exists",
			username: appTest,
			expected: model.Developer{
				Username:       appTest,
				IsDeleted:      false,
				Deposit:        linotypes.NewCoinFromInt64(0),
				AppConsumption: linotypes.NewMiniDollar(0),
			},
			expectedErr: nil,
		},
		{
			name:        "Developer does not exist",
			username:    appDoesNotExists,
			expected:    model.Developer{},
			expectedErr: types.ErrDeveloperNotFound(),
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false)
			got, err := suite.manager.GetDeveloper(suite.Ctx, c.username)
			suite.Equal(c.expectedErr, err)
			suite.Equal(c.expected, got)
			suite.Golden()
			suite.AssertStateUnchanged(false)
		})
	}
}

func (suite *DeveloperManagerSuite) TestGetLiveDeveloper() {
	testCases := []struct {
		name     string
		expected []model.Developer
	}{
		{
			name: "All developers",
			expected: []model.Developer{
				{
					Username:       "test",
					Deposit:        linotypes.NewCoinFromInt64(0),
					AppConsumption: linotypes.NewMiniDollar(0),
					IsDeleted:      false,
				},
			},
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(true)
			suite.Equal(c.expected, suite.manager.GetLiveDevelopers(suite.Ctx))
			suite.Golden()
			suite.AssertStateUnchanged(true)
		})
	}
}

func (suite *DeveloperManagerSuite) TestRegisterDeveloper() {
	username := linotypes.AccountKey("test_username")
	duplicateUsername := linotypes.AccountKey("test")
	userRoleUsername := linotypes.AccountKey("test2")
	voterDuty := votetypes.DutyVoter
	invalidVoterDuty := votetypes.DutyApp

	minDeposit := linotypes.NewCoinFromInt64(50)
	params := param.DeveloperParam{
		DeveloperMinDeposit: minDeposit,
	}
	stake := linotypes.NewCoinFromInt64(100)
	noEnoughStake := linotypes.NewCoinFromInt64(10)
	testCases := []struct {
		name                 string
		username             linotypes.AccountKey
		website              string
		description          string
		appMetaData          string
		accountExist         bool
		voterDuty            *votetypes.VoterDuty
		developerParam       *param.DeveloperParam
		developerParamError  sdk.Error
		voteLinoStake        *linotypes.Coin
		voteLinoStakeError   sdk.Error
		voteAssignDutyCalled bool
		voteAssignDutyError  sdk.Error
		expected             sdk.Error
	}{
		{
			name:         "Fail Account does not exist",
			accountExist: false,
			username:     username,
			expected:     types.ErrAccountNotFound(),
		},
		{
			name:         "Fail Developer already exists",
			accountExist: true,
			username:     duplicateUsername,
			expected:     types.ErrDeveloperAlreadyExist(duplicateUsername),
		},
		{
			name:         "Fail Account is not a Duty Voter",
			accountExist: true,
			voterDuty:    &invalidVoterDuty,
			username:     username,
			expected:     types.ErrInvalidVoterDuty(),
		},
		{
			name:         "Fail Account has user role",
			accountExist: true,
			voterDuty:    &voterDuty,
			username:     userRoleUsername,
			expected:     types.ErrInvalidUserRole(),
		},
		{
			name:                "Fail Error from paramHolder",
			accountExist:        true,
			voterDuty:           &voterDuty,
			developerParam:      &params,
			developerParamError: sdk.ErrInternal("test"),
			username:            username,
			expected:            sdk.ErrInternal("test"),
		},
		{
			name:                "Fail Error from vote.GetLinoStake",
			accountExist:        true,
			voterDuty:           &voterDuty,
			developerParam:      &params,
			developerParamError: nil,
			voteLinoStake:       &stake,
			voteLinoStakeError:  sdk.ErrInternal("test linostake"),
			username:            username,
			expected:            sdk.ErrInternal("test linostake"),
		},
		{
			name:                "Fail not enough stake",
			accountExist:        true,
			voterDuty:           &voterDuty,
			developerParam:      &params,
			developerParamError: nil,
			voteLinoStake:       &noEnoughStake,
			voteLinoStakeError:  nil,
			username:            username,
			expected:            types.ErrInsufficientDeveloperDeposit(),
		},
		{
			name:                 "Fail Error from vote.AssignDuty",
			accountExist:         true,
			voterDuty:            &voterDuty,
			developerParam:       &params,
			developerParamError:  nil,
			voteLinoStake:        &stake,
			voteLinoStakeError:   nil,
			voteAssignDutyCalled: true,
			voteAssignDutyError:  sdk.ErrInternal("test assign duty"),
			username:             username,
			expected:             sdk.ErrInternal("test assign duty"),
		},
		{
			name:                 "Success",
			accountExist:         true,
			voterDuty:            &voterDuty,
			developerParam:       &params,
			developerParamError:  nil,
			voteLinoStake:        &stake,
			voteAssignDutyCalled: true,
			voteLinoStakeError:   nil,
			voteAssignDutyError:  nil,
			username:             username,
			website:              "test website",
			description:          "test description",
			appMetaData:          "test meta",
			expected:             nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false)
			suite.mAccountKeeper.On("DoesAccountExist", mock.Anything, c.username).Return(c.accountExist).Once()
			if c.voterDuty != nil {
				suite.mVoteKeeper.On("GetVoterDuty", mock.Anything, c.username).Return(*c.voterDuty).Once()
			}
			if c.developerParam != nil {
				suite.mParamKeeper.On("GetDeveloperParam", mock.Anything).Return(c.developerParam, c.developerParamError).Once()
			}
			if c.voteLinoStake != nil {
				suite.mVoteKeeper.On("GetLinoStake", mock.Anything, c.username).Return(*c.voteLinoStake, c.voteLinoStakeError).Once()
			}
			if c.voteAssignDutyCalled {
				suite.mVoteKeeper.On("AssignDuty", mock.Anything, c.username, votetypes.DutyApp, params.DeveloperMinDeposit).Return(c.voteAssignDutyError).Once()
			}
			suite.Equal(c.expected, suite.manager.RegisterDeveloper(suite.Ctx, c.username, c.website, c.description, c.appMetaData))
			suite.mAccountKeeper.AssertExpectations(suite.T())
			suite.mVoteKeeper.AssertExpectations(suite.T())
			suite.mParamKeeper.AssertExpectations(suite.T())
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false)
			}
		})
	}
}

func (suite *DeveloperManagerSuite) TestUpdateDeveloper() {
	username := linotypes.AccountKey("test")
	usernameDoesNotExist := linotypes.AccountKey("test-no")
	usernameDeleted := linotypes.AccountKey("test-deleted")
	website := "test website"
	description := "test description"
	meta := "test meta"
	testCases := []struct {
		name     string
		username linotypes.AccountKey
		expected sdk.Error
	}{
		{
			name:     "Fail developer doesnt exist",
			username: usernameDoesNotExist,
			expected: types.ErrDeveloperNotFound(),
		}, {
			name:     "Fail developer deleted",
			username: usernameDeleted,
			expected: types.ErrDeveloperNotFound(),
		}, {
			name:     "Success",
			username: username,
			expected: nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false)
			suite.Equal(c.expected, suite.manager.UpdateDeveloper(suite.Ctx, c.username, website, description, meta))
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false)
			}
		})
	}
}

func (suite *DeveloperManagerSuite) TestIssueIDA() {
	app := linotypes.AccountKey("testapp")
	appDoesNotExists := linotypes.AccountKey("testapp-does-not-exist")
	appHasIDA := linotypes.AccountKey("testapp-has-ida")
	idaName := "test-lemon"
	var idaPrice int64 = 100
	testCases := []struct {
		name     string
		appName  linotypes.AccountKey
		expected sdk.Error
	}{
		{
			name:     "Fail Developer doesnt exist",
			appName:  appDoesNotExists,
			expected: types.ErrDeveloperNotFound(),
		},
		{
			name:     "Fail Developer has already issued IDA",
			appName:  appHasIDA,
			expected: types.ErrIDAIssuedBefore(),
		},
		{
			name:     "Sucess",
			appName:  app,
			expected: nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false)
			suite.Equal(c.expected, suite.manager.IssueIDA(suite.Ctx, c.appName, idaName, idaPrice))
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false)
			}
		})
	}
}

func (suite *DeveloperManagerSuite) TestMintIDA() {
	amount := linotypes.NewCoinFromInt64(1)
	zeroMiniDollar := linotypes.NewMiniDollar(0)
	validMiniDollar := linotypes.NewMiniDollar(1)
	testCases := []struct {
		name                    string
		appName                 linotypes.AccountKey
		amount                  linotypes.Coin
		expected                sdk.Error
		coinToMiniDollar        *linotypes.MiniDollar
		minusCoinFromUserCalled bool
		minusCoinFromUser       sdk.Error
	}{
		{
			name:     "Fail Developer doesnt exist",
			appName:  appDoesNotExists,
			amount:   amount,
			expected: types.ErrDeveloperNotFound(),
		},
		{
			name:     "Fail App doesnt have IDA",
			appName:  appWithoutIDA,
			amount:   amount,
			expected: types.ErrIDANotFound(),
		},
		{
			name:     "Fail App has revoked IDA",
			appName:  appHasRevokedIDA,
			amount:   amount,
			expected: types.ErrIDARevoked(),
		},
		{
			name:             "Fail priceCoinToMiniDollar returns 0",
			appName:          appTest,
			amount:           amount,
			coinToMiniDollar: &zeroMiniDollar,
			expected:         types.ErrExchangeMiniDollarZeroAmount(),
		},
		{
			name:                    "Fail accMinusCoinFromUsername returns error",
			appName:                 appTest,
			amount:                  amount,
			coinToMiniDollar:        &validMiniDollar,
			minusCoinFromUserCalled: true,
			minusCoinFromUser:       sdk.ErrInternal("minus coin from username failed"),
			expected:                sdk.ErrInternal("minus coin from username failed"),
		},
		{
			name:                    "Success",
			appName:                 appTest,
			amount:                  amount,
			coinToMiniDollar:        &validMiniDollar,
			minusCoinFromUserCalled: true,
			minusCoinFromUser:       nil,
			expected:                nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			if c.coinToMiniDollar != nil {
				suite.mPriceKeeper.On("CoinToMiniDollar", c.amount).Return(*c.coinToMiniDollar).Once()
			}
			if c.minusCoinFromUserCalled {
				suite.mAccountKeeper.On("MinusCoinFromUsername", mock.Anything, c.appName, c.amount).Return(c.minusCoinFromUser).Once()
			}
			suite.LoadState(false)
			suite.Equal(c.expected, suite.manager.MintIDA(suite.Ctx, c.appName, c.amount))
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false)
			}
			suite.mPriceKeeper.AssertExpectations(suite.T())
			suite.mAccountKeeper.AssertExpectations(suite.T())
		})
	}
}

func (suite *DeveloperManagerSuite) TestPrivateAppIDAMove() {
	to := linotypes.AccountKey("to")
	from := linotypes.AccountKey("from")
	fromNotEnough := linotypes.AccountKey("from-not-enough")
	fromUnauthed := linotypes.AccountKey("from-unauthed")
	fromNotFound := linotypes.AccountKey("from-not-found")
	toNotFound := linotypes.AccountKey("to-not-found")
	amount := linotypes.NewMiniDollar(100)
	testCases := []struct {
		name     string
		app      linotypes.AccountKey
		from     linotypes.AccountKey
		to       linotypes.AccountKey
		amount   linotypes.MiniDollar
		expected sdk.Error
	}{
		{
			name:     "Fail negative amount",
			app:      appTest,
			from:     from,
			to:       to,
			amount:   linotypes.NewMiniDollar(-1),
			expected: linotypes.ErrInvalidIDAAmount(),
		},
		{
			name:     "Fail from not found",
			app:      appTest,
			from:     fromNotFound,
			to:       to,
			amount:   amount,
			expected: types.ErrNotEnoughIDA(),
		},
		{
			name:     "Fail from unauthed",
			app:      appTest,
			from:     fromUnauthed,
			to:       to,
			amount:   amount,
			expected: types.ErrIDAUnauthed(),
		},
		{
			name:     "Fail from balance not enough",
			app:      appTest,
			from:     fromNotEnough,
			to:       to,
			amount:   amount,
			expected: types.ErrNotEnoughIDA(),
		},
		{
			name:     "Succ should add to existing to account bank",
			app:      appTest,
			from:     from,
			to:       to,
			amount:   amount,
			expected: nil,
		},
		{
			name:     "Succ should create new to account bank",
			app:      appTest,
			from:     from,
			to:       toNotFound,
			amount:   amount,
			expected: nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false, "IDABasic")
			suite.Equal(c.expected, suite.manager.appIDAMove(suite.Ctx, c.app, c.from, c.to, c.amount))
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false, "IDABasic")
			}
		})
	}
}

func (suite *DeveloperManagerSuite) TestAppTransferIDA() {
	to := linotypes.AccountKey("to")
	from := linotypes.AccountKey("from")
	app1affiliated := linotypes.AccountKey("testapp-affiliated")
	exists := true
	doesntExists := false
	amount := sdk.NewInt(1)
	testCases := []struct {
		name       string
		appName    linotypes.AccountKey
		signer     linotypes.AccountKey
		from       linotypes.AccountKey
		to         linotypes.AccountKey
		fromExists *bool
		toExists   *bool
		expected   sdk.Error
	}{
		{
			name:     "Fail from and to both not sender",
			appName:  appTest,
			signer:   appTest,
			from:     from,
			to:       to,
			expected: types.ErrInvalidTransferTarget(),
		},
		{
			name:     "Fail App does not exist",
			appName:  appDoesNotExists,
			signer:   appDoesNotExists,
			from:     from,
			to:       appDoesNotExists,
			expected: types.ErrDeveloperNotFound(),
		},
		{
			name:     "Fail App doesnt have IDA",
			appName:  appWithoutIDA,
			signer:   appWithoutIDA,
			from:     from,
			to:       appWithoutIDA,
			expected: types.ErrIDANotFound(),
		},
		{
			name:     "Fail App has revoked IDA",
			appName:  appHasRevokedIDA,
			signer:   appHasRevokedIDA,
			from:     from,
			to:       appHasRevokedIDA,
			expected: types.ErrIDARevoked(),
		},
		{
			name:       "Fail from account doesnt exist",
			appName:    appTest,
			signer:     appTest,
			from:       from,
			to:         appTest,
			fromExists: &doesntExists,
			expected:   types.ErrAccountNotFound(),
		},
		{
			name:       "Fail to account doesnt exist",
			appName:    appTest,
			signer:     appTest,
			from:       from,
			to:         appTest,
			fromExists: &exists,
			toExists:   &doesntExists,
			expected:   types.ErrAccountNotFound(),
		},
		{
			name:       "Fail signer does not match",
			appName:    appTest,
			signer:     from,
			from:       from,
			to:         appTest,
			fromExists: &exists,
			toExists:   &exists,
			expected:   types.ErrInvalidSigner(),
		},
		{
			name:       "Success Transfer from App",
			appName:    appTest,
			signer:     appTest,
			from:       appTest,
			to:         to,
			fromExists: &exists,
			toExists:   &exists,
			expected:   nil,
		},
		{
			name:       "Success Transfer to App",
			appName:    appTest,
			signer:     appTest,
			from:       from,
			to:         appTest,
			fromExists: &exists,
			toExists:   &exists,
			expected:   nil,
		},
		{
			name:       "Success Transfer to App by affiliated",
			appName:    appTest,
			signer:     app1affiliated,
			from:       from,
			to:         appTest,
			fromExists: &exists,
			toExists:   &exists,
			expected:   nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			if c.fromExists != nil {
				suite.mAccountKeeper.On("DoesAccountExist", mock.Anything, c.from).Return(*c.fromExists).Once()
			}
			if c.toExists != nil {
				suite.mAccountKeeper.On("DoesAccountExist", mock.Anything, c.to).Return(*c.toExists).Once()
			}
			suite.LoadState(false, "IDABasic")
			suite.Equal(c.expected, suite.manager.AppTransferIDA(
				suite.Ctx, c.appName, c.signer, amount, c.from, c.to))
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false, "IDABasic")
			}
			suite.mAccountKeeper.AssertExpectations(suite.T())
		})
	}
}

func (suite *DeveloperManagerSuite) TestMoveIDA() {
	to := linotypes.AccountKey("to")
	from := linotypes.AccountKey("from")
	fromNotEnough := linotypes.AccountKey("from-not-enough")
	fromUnauthed := linotypes.AccountKey("from-unauthed")
	fromNotFound := linotypes.AccountKey("from-not-found")
	toNotFound := linotypes.AccountKey("to-not-found")
	amount := linotypes.NewMiniDollar(100)
	exists := true
	doesntExists := false
	testCases := []struct {
		name       string
		app        linotypes.AccountKey
		from       linotypes.AccountKey
		to         linotypes.AccountKey
		amount     linotypes.MiniDollar
		fromExists *bool
		toExists   *bool
		expected   sdk.Error
	}{
		{
			name:     "Fail Developer doesnt exist",
			app:      appDoesNotExists,
			amount:   amount,
			expected: types.ErrDeveloperNotFound(),
		},
		{
			name:     "Fail App doesnt have IDA",
			app:      appWithoutIDA,
			amount:   amount,
			expected: types.ErrIDANotFound(),
		},
		{
			name:     "Fail App has revoked IDA",
			app:      appHasRevokedIDA,
			amount:   amount,
			expected: types.ErrIDARevoked(),
		},
		{
			name:       "Fail from account doesnt exist",
			app:        appTest,
			from:       from,
			to:         appTest,
			fromExists: &doesntExists,
			expected:   types.ErrAccountNotFound(),
		},
		{
			name:       "Fail to account doesnt exist",
			app:        appTest,
			from:       from,
			to:         appTest,
			fromExists: &exists,
			toExists:   &doesntExists,
			expected:   types.ErrAccountNotFound(),
		},
		{
			name:       "Fail negative amount",
			app:        appTest,
			from:       from,
			to:         to,
			fromExists: &exists,
			toExists:   &exists,
			amount:     linotypes.NewMiniDollar(-1),
			expected:   linotypes.ErrInvalidIDAAmount(),
		},
		{
			name:       "Fail from bank not found",
			app:        appTest,
			from:       fromNotFound,
			to:         to,
			amount:     amount,
			fromExists: &exists,
			toExists:   &exists,
			expected:   types.ErrNotEnoughIDA(),
		},
		{
			name:       "Fail from bank unauthed",
			app:        appTest,
			from:       fromUnauthed,
			to:         to,
			amount:     amount,
			fromExists: &exists,
			toExists:   &exists,
			expected:   types.ErrIDAUnauthed(),
		},
		{
			name:       "Fail from bank balance not enough",
			app:        appTest,
			from:       fromNotEnough,
			to:         to,
			amount:     amount,
			fromExists: &exists,
			toExists:   &exists,
			expected:   types.ErrNotEnoughIDA(),
		},
		{
			name:       "Succes should add to existing to account bank",
			app:        appTest,
			from:       from,
			to:         to,
			amount:     amount,
			fromExists: &exists,
			toExists:   &exists,
			expected:   nil,
		},
		{
			name:       "Succes should create new to account bank",
			app:        appTest,
			from:       from,
			to:         toNotFound,
			amount:     amount,
			fromExists: &exists,
			toExists:   &exists,
			expected:   nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			if c.fromExists != nil {
				suite.mAccountKeeper.On("DoesAccountExist", mock.Anything, c.from).Return(*c.fromExists).Once()
			}
			if c.toExists != nil {
				suite.mAccountKeeper.On("DoesAccountExist", mock.Anything, c.to).Return(*c.toExists).Once()
			}
			suite.LoadState(false, "IDABasic")
			suite.Equal(c.expected, suite.manager.MoveIDA(suite.Ctx, c.app, c.from, c.to, c.amount))
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false, "IDABasic")
			}
			suite.mAccountKeeper.AssertExpectations(suite.T())
		})
	}
}

// This also includes GetMiniIDAPrice and GetIDA since they are really simple
func (suite *DeveloperManagerSuite) TestPrivateValidAppIDA() {
	zeroPrices := linotypes.NewMiniDollar(0)
	zeroIDA := model.AppIDA{}
	price := linotypes.NewMiniDollar(10000)
	ida := model.AppIDA{
		App:             appTest,
		Name:            "test-lemon",
		MiniIDAPrice:    price,
		RevokeCoinPrice: zeroPrices,
	}
	testCases := []struct {
		name          string
		app           linotypes.AccountKey
		expectedError sdk.Error
		expectedPrice linotypes.MiniDollar
		expectedIDA   model.AppIDA
	}{
		{
			name:          "Fail Developer doesnt exist",
			app:           appDoesNotExists,
			expectedError: types.ErrDeveloperNotFound(),
			expectedPrice: zeroPrices,
			expectedIDA:   zeroIDA,
		},
		{
			name:          "Fail App doesnt have IDA",
			app:           appWithoutIDA,
			expectedError: types.ErrIDANotFound(),
			expectedPrice: zeroPrices,
			expectedIDA:   zeroIDA,
		},
		{
			name:          "Fail App has revoked IDA",
			app:           appHasRevokedIDA,
			expectedError: types.ErrIDARevoked(),
			expectedPrice: zeroPrices,
			expectedIDA:   zeroIDA,
		},
		{
			name:          "Success",
			app:           appTest,
			expectedError: nil,
			expectedPrice: price,
			expectedIDA:   ida,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false, "IDABasic")
			var ida *model.AppIDA = nil
			if c.expectedError == nil {
				ida = &c.expectedIDA
			}
			i, e := suite.manager.validAppIDA(suite.Ctx, c.app)
			suite.Equal(ida, i)
			suite.Equal(c.expectedError, e)
			price, e := suite.manager.GetMiniIDAPrice(suite.Ctx, c.app)
			suite.Equal(c.expectedPrice, price)
			suite.Equal(c.expectedError, e)
			id, e := suite.manager.GetIDA(suite.Ctx, c.app)
			suite.Equal(c.expectedIDA, id)
			suite.Equal(c.expectedError, e)
			suite.Golden()
			suite.AssertStateUnchanged(false, "IDABasic")
		})
	}
}

// This also tests GetIDABank
func (suite *DeveloperManagerSuite) TestBurnIDA() {
	zeroCoin := linotypes.NewCoinFromInt64(0)
	coin := linotypes.NewCoinFromInt64(1)
	coinAmountMoreThanPool := linotypes.NewCoinFromInt64(11)
	used := linotypes.NewMiniDollar(100)
	userNotEnough := linotypes.AccountKey("from-not-enough")
	user := linotypes.AccountKey("from")
	amount := linotypes.NewMiniDollar(10000)
	exists := true
	noExists := false
	testCases := []struct {
		name          string
		app           linotypes.AccountKey
		user          linotypes.AccountKey
		amount        linotypes.MiniDollar
		expectedError sdk.Error
		expectedCoin  linotypes.Coin
		accountExists *bool
		bought        *linotypes.Coin
		used          *linotypes.MiniDollar
	}{
		{
			name:          "Fail app does not exist",
			app:           appDoesNotExists,
			user:          user,
			amount:        amount,
			expectedCoin:  zeroCoin,
			expectedError: types.ErrDeveloperNotFound(),
		},
		{
			name:          "Fail user does not exist",
			app:           appTest,
			user:          user,
			amount:        amount,
			expectedCoin:  zeroCoin,
			accountExists: &noExists,
			expectedError: types.ErrAccountNotFound(),
		},
		{
			name:          "Fail user does not have enough mini dollar",
			app:           appTest,
			user:          userNotEnough,
			amount:        amount,
			expectedCoin:  zeroCoin,
			accountExists: &exists,
			expectedError: types.ErrNotEnoughIDA(),
		},
		{
			name:          "Fail cannot burn 0 coin",
			app:           appTest,
			user:          user,
			amount:        amount,
			expectedCoin:  zeroCoin,
			accountExists: &exists,
			expectedError: types.ErrBurnZeroIDA(),
			bought:        &zeroCoin,
			used:          &used,
		},
		{
			name:          "Fail burn amount more than reserve pool",
			app:           appTest,
			user:          user,
			amount:        amount,
			expectedCoin:  zeroCoin,
			accountExists: &exists,
			expectedError: types.ErrInsuffientReservePool(),
			bought:        &coinAmountMoreThanPool,
			used:          &used,
		},
		{
			name:          "Success",
			app:           appTest,
			user:          user,
			amount:        amount,
			expectedCoin:  coin,
			accountExists: &exists,
			expectedError: nil,
			bought:        &coin,
			used:          &used,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			if c.accountExists != nil {
				suite.mAccountKeeper.On("DoesAccountExist", mock.Anything, c.user).Return(*c.accountExists).Once()
			}
			if c.bought != nil {
				suite.mPriceKeeper.On("MiniDollarToCoin", c.amount).Return(*c.bought, *c.used).Once()
			}
			suite.LoadState(false, "IDABasic")
			coin, err := suite.manager.BurnIDA(suite.Ctx, c.app, c.user, c.amount)
			suite.Equal(c.expectedCoin, coin)
			suite.Equal(c.expectedError, err)
			suite.Golden()
			if c.expectedError != nil {
				suite.AssertStateUnchanged(false, "IDABasic")
			}
			suite.mAccountKeeper.AssertExpectations(suite.T())
			suite.mPriceKeeper.AssertExpectations(suite.T())
		})
	}
}

func (suite *DeveloperManagerSuite) TestGetIDABank() {
	zeroBank := model.IDABank{}
	user := linotypes.AccountKey("from")
	bank := model.IDABank{
		Balance: linotypes.NewMiniDollar(100000),
	}
	exists := true
	noExists := false
	testCases := []struct {
		name          string
		app           linotypes.AccountKey
		user          linotypes.AccountKey
		expectedError sdk.Error
		expectedBank  model.IDABank
		accountExists *bool
	}{
		{
			name:          "Fail app does not exist",
			app:           appDoesNotExists,
			user:          user,
			expectedError: types.ErrDeveloperNotFound(),
			expectedBank:  zeroBank,
		},
		{
			name:          "Fail user does not exist",
			app:           appTest,
			user:          user,
			accountExists: &noExists,
			expectedBank:  zeroBank,
			expectedError: types.ErrAccountNotFound(),
		},
		{
			name:          "Success",
			app:           appTest,
			user:          user,
			accountExists: &exists,
			expectedBank:  bank,
			expectedError: nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			if c.accountExists != nil {
				suite.mAccountKeeper.On("DoesAccountExist", mock.Anything, c.user).Return(*c.accountExists).Once()
			}
			suite.LoadState(false, "IDABasic")
			bank, err := suite.manager.GetIDABank(suite.Ctx, c.app, c.user)
			suite.Equal(c.expectedBank, bank)
			suite.Equal(c.expectedError, err)
			suite.Golden()
			suite.AssertStateUnchanged(false, "IDABasic")
			suite.mAccountKeeper.AssertExpectations(suite.T())
		})
	}
}

func (suite *DeveloperManagerSuite) TestUpdateAffiliated() {
	t := true
	f := false
	username := linotypes.AccountKey("testuser")
	appMaxAffiliated := linotypes.AccountKey("testapp-max-affiliated")
	userAffiliated := linotypes.AccountKey("testuser-affiliated")
	userDeveloper := linotypes.AccountKey("testuser-developer")
	userDeactivate := linotypes.AccountKey("testuser-deactivate")
	voterDuty := votetypes.DutyVoter
	invalidVoterDuty := votetypes.DutyApp
	testCases := []struct {
		name          string
		appName       linotypes.AccountKey
		username      linotypes.AccountKey
		activate      bool
		expected      sdk.Error
		accountExists *bool
		vote          *votetypes.VoterDuty
	}{
		{
			name:     "Fail app doesnt exist",
			appName:  appDoesNotExists,
			username: username,
			expected: types.ErrDeveloperNotFound(),
		},
		{
			name:          "Fail user doesnt exist",
			appName:       appTest,
			username:      username,
			accountExists: &f,
			expected:      types.ErrAccountNotFound(),
		},
		{
			name:          "Fail max affiliated account reached",
			appName:       appMaxAffiliated,
			username:      username,
			accountExists: &t,
			expected:      types.ErrMaxAffiliatedExceeded(),
		},
		{
			name:          "Fail activate user is already affiliated with some app",
			appName:       appTest,
			username:      userAffiliated,
			accountExists: &t,
			activate:      t,
			expected:      types.ErrInvalidAffiliatedAccount("is affiliated already"),
		},
		{
			name:          "Fail activate user is already a developer",
			appName:       appTest,
			username:      userDeveloper,
			accountExists: &t,
			activate:      t,
			expected:      types.ErrInvalidAffiliatedAccount("is/was developer"),
		},
		{
			name:          "Fail activate user has some other duty",
			appName:       appTest,
			username:      username,
			accountExists: &t,
			vote:          &invalidVoterDuty,
			activate:      t,
			expected:      types.ErrInvalidAffiliatedAccount("on duty of something else"),
		},
		{
			name:          "Success activate",
			appName:       appTest,
			username:      username,
			accountExists: &t,
			vote:          &voterDuty,
			activate:      t,
			expected:      nil,
		},
		{
			name:          "Fail deactivate user isn't affiliated with any app",
			appName:       appTest,
			username:      username,
			accountExists: &t,
			vote:          &voterDuty,
			activate:      f,
			expected:      types.ErrInvalidUserRole(),
		},
		{
			name:          "Fail deactivate user has different affiliated account",
			appName:       appTest,
			username:      userAffiliated,
			accountExists: &t,
			vote:          &voterDuty,
			activate:      f,
			expected:      types.ErrInvalidAffiliatedAccount("not affiliated account of provided app"),
		},
		{
			name:          "Success deactivate",
			appName:       appTest,
			username:      userDeactivate,
			accountExists: &t,
			vote:          &voterDuty,
			activate:      f,
			expected:      nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			if c.accountExists != nil {
				suite.mAccountKeeper.On("DoesAccountExist", mock.Anything, c.username).Return(*c.accountExists).Once()
			}
			if c.vote != nil && c.activate {
				suite.mVoteKeeper.On("GetVoterDuty", mock.Anything, c.username).Return(*c.vote).Once()
			}
			suite.LoadState(false, "AffiliatedBasic")
			suite.Equal(c.expected, suite.manager.UpdateAffiliated(suite.Ctx, c.appName, c.username, c.activate))
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false, "AffiliatedBasic")
			}
			suite.mAccountKeeper.AssertExpectations(suite.T())
			suite.mVoteKeeper.AssertExpectations(suite.T())
		})
	}
}

func (suite *DeveloperManagerSuite) TestGetAffiliatingApp() {
	userDev := linotypes.AccountKey("testuser-developer")
	username := linotypes.AccountKey("testuser-affiliated")
	app := linotypes.AccountKey("testapp-a")
	usernameNotAf := linotypes.AccountKey("testuser")
	testCases := []struct {
		name          string
		username      linotypes.AccountKey
		expectedApp   linotypes.AccountKey
		expectedError sdk.Error
	}{
		{
			name:          "Fail no affiliation found",
			username:      usernameNotAf,
			expectedApp:   "",
			expectedError: types.ErrInvalidUserRole(),
		},
		{
			name:          "Success is developer",
			username:      userDev,
			expectedApp:   userDev,
			expectedError: nil,
		},
		{
			name:          "Success found affiliated app",
			username:      username,
			expectedApp:   app,
			expectedError: nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false, "AffiliatedBasic")
			app, err := suite.manager.GetAffiliatingApp(suite.Ctx, c.username)
			suite.Equal(c.expectedApp, app)
			suite.Equal(c.expectedError, err)
			suite.Golden()
			suite.AssertStateUnchanged(false, "AffiliatedBasic")
		})
	}
}

func (suite *DeveloperManagerSuite) TestGetAffiliated() {
	testCases := []struct {
		name     string
		app      linotypes.AccountKey
		expected []linotypes.AccountKey
	}{
		{
			name:     "Success developer doesnt exist",
			app:      appDoesNotExists,
			expected: nil,
		},
		{
			name: "Success get all affiliated account",
			app:  appTest,
			expected: []linotypes.AccountKey{
				linotypes.AccountKey("testuser-deactivate"),
			},
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false, "AffiliatedBasic")
			suite.Equal(c.expected, suite.manager.GetAffiliated(suite.Ctx, c.app))
			suite.Golden()
			suite.AssertStateUnchanged(false, "AffiliatedBasic")
		})
	}
}

func (suite *DeveloperManagerSuite) TestUpdateIDAAuth() {
	t := true
	f := false
	userAf := linotypes.AccountKey("testuser-affiliated")
	user := linotypes.AccountKey("testuser")
	testCases := []struct {
		name     string
		app      linotypes.AccountKey
		username linotypes.AccountKey
		active   bool
		expected sdk.Error
		aExists  *bool
	}{
		{
			name:     "Fail developer doesnt exist",
			app:      appDoesNotExists,
			expected: types.ErrDeveloperNotFound(),
		},
		{
			name:     "Fail account doesnt exist",
			app:      appTest,
			aExists:  &f,
			expected: types.ErrAccountNotFound(),
		},
		{
			name:     "Fail user is affiliated account",
			app:      appTest,
			username: userAf,
			aExists:  &t,
			expected: types.ErrInvalidIDAAuth(),
		},
		{
			name:     "Fail bank already has the target active state",
			app:      appTest,
			username: user,
			aExists:  &t,
			active:   true,
			expected: types.ErrInvalidIDAAuth(),
		},
		{
			name:     "Success",
			app:      appTest,
			username: user,
			aExists:  &t,
			active:   false,
			expected: nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			if c.aExists != nil {
				suite.mAccountKeeper.On("DoesAccountExist", mock.Anything, c.username).Return(*c.aExists).Once()
			}
			suite.LoadState(false, "IDABasic")
			suite.Equal(c.expected, suite.manager.UpdateIDAAuth(suite.Ctx, c.app, c.username, c.active))
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false, "IDABasic")
			}
			suite.mAccountKeeper.AssertExpectations(suite.T())
		})
	}
}

func (suite *DeveloperManagerSuite) TestReportConsumption() {
	testCases := []struct {
		name        string
		app         linotypes.AccountKey
		consumption linotypes.MiniDollar
		expected    sdk.Error
	}{
		{
			name:     "Fail developer doesnt exist",
			app:      appDoesNotExists,
			expected: types.ErrDeveloperNotFound(),
		},
		{
			name:        "Success",
			app:         appTest,
			consumption: linotypes.NewMiniDollar(10),
			expected:    nil,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			suite.LoadState(false, "DeveloperBasic")
			suite.Equal(c.expected, suite.manager.ReportConsumption(suite.Ctx, c.app, c.consumption))
			suite.Golden()
			if c.expected != nil {
				suite.AssertStateUnchanged(false, "DeveloperBasic")
			}
		})
	}
}

func (suite *DeveloperManagerSuite) TestDistributeDevInflation() {
	zeroCoin := linotypes.NewCoinFromInt64(0)
	inflation := linotypes.NewCoinFromInt64(100)
	testCases := []struct {
		name                string
		expected            sdk.Error
		totalInflation      *linotypes.Coin
		totalInflationError sdk.Error
		shares              []linotypes.Coin
		stateChange         bool
		addCoinError        sdk.Error
	}{
		{
			name:                "Fail error from global.PopDeveloperMonthlyInflation",
			expected:            sdk.ErrInternal(""),
			totalInflation:      &zeroCoin,
			totalInflationError: sdk.ErrInternal(""),
		},
		{
			name:                "Fail error from acc.AddCoinToUsername",
			expected:            sdk.ErrInternal(""),
			totalInflation:      &inflation,
			totalInflationError: nil,
			shares: []linotypes.Coin{
				linotypes.NewCoinFromInt64(50),
			},
			addCoinError: sdk.ErrInternal(""),
		},
		{
			name:     "Succ no developers",
			expected: nil,
		},
		{
			name:                "Success even distribution",
			expected:            nil,
			totalInflation:      &inflation,
			totalInflationError: nil,
			shares: []linotypes.Coin{
				linotypes.NewCoinFromInt64(50),
				linotypes.NewCoinFromInt64(50),
			},
		},
		{
			name:                "Success even distribution with remainder",
			expected:            nil,
			totalInflation:      &inflation,
			totalInflationError: nil,
			shares: []linotypes.Coin{
				linotypes.NewCoinFromInt64(33),
				linotypes.NewCoinFromInt64(33),
				linotypes.NewCoinFromInt64(34),
			},
			stateChange: true,
		},
		{
			name:                "Success distribute according to consumption",
			expected:            nil,
			totalInflation:      &inflation,
			totalInflationError: nil,
			shares: []linotypes.Coin{
				linotypes.NewCoinFromInt64(14),
				linotypes.NewCoinFromInt64(29),
				linotypes.NewCoinFromInt64(57),
			},
			stateChange: true,
		},
	}
	for _, c := range testCases {
		suite.Run(c.name, func() {
			if c.totalInflation != nil {
				suite.mGlobalKeeper.On("PopDeveloperMonthlyInflation", mock.Anything).Return(*c.totalInflation, c.totalInflationError).Once()
			}
			for i, share := range c.shares {
				suite.mAccountKeeper.On("AddCoinToUsername", mock.Anything, linotypes.AccountKey(fmt.Sprintf("testapp-%d", i)), share).Return(c.addCoinError).Once()
			}
			suite.LoadState(true)
			suite.Equal(c.expected, suite.manager.DistributeDevInflation(suite.Ctx))
			suite.Golden()
			if !c.stateChange {
				suite.AssertStateUnchanged(true)
			}
			suite.mGlobalKeeper.AssertExpectations(suite.T())
			suite.mAccountKeeper.AssertExpectations(suite.T())
		})
	}
}

// import (
// 	"testing"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/lino-network/lino/types"
// 	"github.com/stretchr/testify/assert"
// )

// func TestReportConsumption(t *testing.T) {
// 	ctx, _, dm, _ := setupTest(t, 0)
// 	dm.InitGenesis(ctx)

// 	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
// 	dm.RegisterDeveloper(ctx, "developer1", devParam.DeveloperMinDeposit, "", "", "")
// 	dm.RegisterDeveloper(ctx, "developer2", devParam.DeveloperMinDeposit, "", "", "")

// 	con1 := types.NewCoinFromInt64(100)
// 	dm.ReportConsumption(ctx, "developer1", con1)
// 	p1, _ := dm.GetConsumptionWeight(ctx, "developer1")
// 	assert.True(t, p1.Equal(types.NewDecFromRat(1, 1)))

// 	con2 := types.NewCoinFromInt64(100)
// 	dm.ReportConsumption(ctx, "developer2", con2)
// 	p2, _ := dm.GetConsumptionWeight(ctx, "developer1")
// 	assert.True(t, p2.Equal(types.NewDecFromRat(1, 2)))

// 	dm.ClearConsumption(ctx)
// 	p3, _ := dm.GetConsumptionWeight(ctx, "developer1")
// 	assert.True(t, p3.Equal(types.NewDecFromRat(1, 2)))

// 	testCases := map[string]struct {
// 		developer1Consumption             types.Coin
// 		developer2Consumption             types.Coin
// 		expectDeveloper1ConsumptionWeight sdk.Dec
// 		expectDeveloper2ConsumptionWeight sdk.Dec
// 	}{
// 		"test normal consumption": {
// 			developer1Consumption:             types.NewCoinFromInt64(2500 * types.Decimals),
// 			developer2Consumption:             types.NewCoinFromInt64(7500 * types.Decimals),
// 			expectDeveloper1ConsumptionWeight: types.NewDecFromRat(1, 4),
// 			expectDeveloper2ConsumptionWeight: types.NewDecFromRat(3, 4),
// 		},
// 		"test empty consumption": {
// 			developer1Consumption:             types.NewCoinFromInt64(0),
// 			developer2Consumption:             types.NewCoinFromInt64(0),
// 			expectDeveloper1ConsumptionWeight: types.NewDecFromRat(1, 2),
// 			expectDeveloper2ConsumptionWeight: types.NewDecFromRat(1, 2),
// 		},
// 		"large numbers": {
// 			developer1Consumption:             types.NewCoinFromInt64(3333333),
// 			developer2Consumption:             types.NewCoinFromInt64(4444444),
// 			expectDeveloper1ConsumptionWeight: types.NewDecFromRat(3333333, 7777777),
// 			expectDeveloper2ConsumptionWeight: types.NewDecFromRat(4444444, 7777777),
// 		},
// 	}
// 	for testName, tc := range testCases {
// 		dm.ReportConsumption(ctx, "developer1", tc.developer1Consumption)
// 		dm.ReportConsumption(ctx, "developer2", tc.developer2Consumption)

// 		p1, _ := dm.GetConsumptionWeight(ctx, "developer1")
// 		if !tc.expectDeveloper1ConsumptionWeight.Equal(p1) {
// 			t.Errorf("%s: diff developer1 usage weight, got %v, want %v",
// 				testName, p1, tc.expectDeveloper1ConsumptionWeight)
// 			return
// 		}

// 		p2, _ := dm.GetConsumptionWeight(ctx, "developer2")
// 		if !tc.expectDeveloper2ConsumptionWeight.Equal(p2) {
// 			t.Errorf("%s: diff developer2 usage weight, got %v, want %v",
// 				testName, p2, tc.expectDeveloper2ConsumptionWeight)
// 			return
// 		}
// 		dm.ClearConsumption(ctx)
// 	}
// }

// func TestDistributeInflationToDevelopers(t *testing.T) {
// 	cases := map[string]struct {
// 		beforeDistributionInflationPool types.Coin
// 		pastMinutes                     int64
// 		numberOfDevelopers              int
// 		consumptionList                 []types.MiniDollar
// 	}{
// 		"distribute to one developer with zero consumption": {
// 			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
// 			numberOfDevelopers:              1,
// 			pastMinutes:                     types.MinutesPerMonth,
// 			consumptionList:                 []types.MiniDollar{types.NewMiniDollar(0)},
// 		},
// 		"distribute to five developers with zero consumption": {
// 			beforeDistributionInflationPool: types.NewCoinFromInt64(1000 * types.Decimals),
// 			numberOfDevelopers:              5,
// 			pastMinutes:                     types.MinutesPerMonth,
// 			consumptionList: []types.MiniDollar{
// 				types.NewMiniDollar(0),
// 				types.NewMiniDollar(0),
// 				types.NewMiniDollar(0),
// 				types.NewMiniDollar(0),
// 				types.NewMiniDollar(0)},
// 		},
// 		"test inflation need to be rounded case": {
// 			beforeDistributionInflationPool: types.NewCoinFromInt64(100 * types.Decimals),
// 			numberOfDevelopers:              3,
// 			pastMinutes:                     types.MinutesPerMonth,
// 			consumptionList: []types.MiniDollar{
// 				types.NewMiniDollar(0),
// 				types.NewMiniDollar(0),
// 				types.NewMiniDollar(0),
// 			},
// 		},
// 		"test different consumption case": {
// 			beforeDistributionInflationPool: types.NewCoinFromInt64(100 * types.Decimals),
// 			numberOfDevelopers:              3,
// 			pastMinutes:                     types.MinutesPerMonth,
// 			consumptionList: []types.MiniDollar{
// 				types.NewMiniDollar(1000 * types.Decimals),
// 				types.NewMiniDollar(2000 * types.Decimals),
// 				types.NewMiniDollar(20),
// 			},
// 		},
// 	}
// 	for testName, cs := range cases {
// 		lb := newLinoBlockchain(t, 21)
// 		ctx := lb.BaseApp.NewContext(true, abci.Header{})
// 		devStorage := devModel.NewDeveloperStorage(lb.CapKeyDeveloperStore)
// 		totalConsumption := types.NewMiniDollar(0)
// 		for i := 0; i < cs.numberOfDevelopers; i++ {
// 			err := lb.accountManager.CreateAccount(
// 				ctx, types.AccountKey("dev"+strconv.Itoa(i)),
// 				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey())
// 			if err != nil {
// 				t.Errorf("%s: failed to register account, got err %v", testName, err)
// 			}
// 			err = lb.developerManager.RegisterDeveloper(
// 				ctx, types.AccountKey("dev"+strconv.Itoa(i)), "", "", "")
// 			if err != nil {
// 				t.Errorf("%s: failed to register developer, got err %v", testName, err)
// 			}
// 			developer, _ := devStorage.GetDeveloper(ctx, types.AccountKey("dev"+strconv.Itoa(i)))
// 			developer.AppConsumption = cs.consumptionList[i]
// 			devStorage.SetDeveloper(ctx, *developer)
// 			totalConsumption = totalConsumption.Plus(cs.consumptionList[i])
// 		}
// 		globalStore := globalModel.NewGlobalStorage(lb.CapKeyGlobalStore)
// 		err := globalStore.SetInflationPool(ctx, &globalModel.InflationPool{
// 			DeveloperInflationPool: cs.beforeDistributionInflationPool,
// 		})
// 		if err != nil {
// 			t.Errorf("%s: failed to set inflation pool, got err %v", testName, err)
// 		}

// 		lb.distributeInflationToDeveloper(ctx)
// 		inflationPool, err := globalStore.GetInflationPool(ctx)
// 		if err != nil {
// 			t.Errorf("%s: failed to get inflation pool, got err %v", testName, err)
// 		}

// 		if !inflationPool.DeveloperInflationPool.IsZero() {
// 			t.Errorf(
// 				"%s: diff developer inflation pool, got %v, want %v",
// 				testName, inflationPool.DeveloperInflationPool,
// 				types.NewCoinFromInt64(0))
// 			return
// 		}

// 		actualInflation := types.NewCoinFromInt64(0)
// 		for i := 0; i < cs.numberOfDevelopers; i++ {
// 			saving, err :=
// 				lb.accountManager.GetSavingFromUsername(
// 					ctx, types.AccountKey("dev"+strconv.Itoa(i)))
// 			assert.Nil(t, err)
// 			var inflation types.Coin
// 			if totalConsumption.IsZero() {
// 				inflation =
// 					types.DecToCoin(
// 						types.NewDecFromRat(1, int64(len(cs.consumptionList))).
// 							Mul(cs.beforeDistributionInflationPool.ToDec()))
// 			} else {
// 				inflation =
// 					types.DecToCoin(
// 						cs.consumptionList[i].ToDec().
// 							Quo(totalConsumption.ToDec()).
// 							Mul(cs.beforeDistributionInflationPool.ToDec()))
// 			}
// 			if i == (cs.numberOfDevelopers - 1) {
// 				inflation = cs.beforeDistributionInflationPool.Minus(actualInflation)
// 			}
// 			actualInflation = actualInflation.Plus(inflation)
// 			if !saving.IsEqual(inflation) {
// 				t.Errorf(
// 					"%s: diff inflation for %v, got %v, want %v",
// 					testName, "dev"+strconv.Itoa(i), inflation,
// 					saving)
// 				return
// 			}
// 			developer, err := devStorage.GetDeveloper(ctx, types.AccountKey("dev"+strconv.Itoa(i)))
// 			assert.Nil(t, err)
// 			assert.True(t, developer.AppConsumption.IsZero())
// 		}
// 	}
// }
