package manager

import (
	"math"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"

	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	linotypes "github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account/mocks"
	global "github.com/lino-network/lino/x/global/mocks"
	"github.com/lino-network/lino/x/validator/model"
	"github.com/lino-network/lino/x/validator/types"
	vote "github.com/lino-network/lino/x/vote/mocks"
	votetypes "github.com/lino-network/lino/x/vote/types"
)

type ValidatorManagerTestSuite struct {
	testsuites.CtxTestSuite
	vm       ValidatorManager
	baseTime time.Time
	// deps
	ph     *param.ParamKeeper
	global *global.GlobalKeeper
	vote   *vote.VoteKeeper
	acc    *acc.AccountKeeper
}

func TestValidatorManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorManagerTestSuite))
}

func (suite *ValidatorManagerTestSuite) SetupTest() {
	suite.baseTime = time.Now()
	testValidatorKey := sdk.NewKVStoreKey("validator")
	suite.SetupCtx(0, suite.baseTime.Add(3*time.Second), testValidatorKey)
	suite.global = &global.GlobalKeeper{}
	suite.ph = &param.ParamKeeper{}
	suite.vote = &vote.VoteKeeper{}
	suite.acc = &acc.AccountKeeper{}

	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("user1")).Return(linotypes.NewCoinFromInt64(300), nil).Maybe()
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("val")).Return(linotypes.NewCoinFromInt64(300), nil).Maybe()
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("jail1")).Return(linotypes.NewCoinFromInt64(200000*linotypes.Decimals), nil).Maybe()
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("jail2")).Return(linotypes.NewCoinFromInt64(200), nil).Maybe()
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("abs")).Return(linotypes.NewCoinFromInt64(200), nil).Maybe()
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("byz")).Return(linotypes.NewCoinFromInt64(2000000*linotypes.Decimals), nil).Maybe()
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("changedVoter")).Return(linotypes.NewCoinFromInt64(600), nil).Maybe()
	suite.vote.On("GetVoterDuty", suite.Ctx, linotypes.AccountKey("val")).Return(votetypes.DutyVoter, nil).Maybe()
	suite.vote.On("AssignDuty", suite.Ctx, linotypes.AccountKey("val"), votetypes.DutyValidator,
		linotypes.NewCoinFromInt64(200000*linotypes.Decimals)).Return(nil).Maybe()
	suite.vote.On("UnassignDuty", suite.Ctx, linotypes.AccountKey("val"), mock.Anything).Return(nil).Maybe()
	suite.vote.On("SlashStake", suite.Ctx, linotypes.AccountKey("abs"),
		linotypes.NewCoinFromInt64(200*linotypes.Decimals), linotypes.InflationValidatorPool).Return(linotypes.NewCoinFromInt64(200*linotypes.Decimals), nil).Maybe()
	suite.vote.On("SlashStake", suite.Ctx, linotypes.AccountKey("byz"),
		linotypes.NewCoinFromInt64(1000*linotypes.Decimals), linotypes.InflationValidatorPool).Return(linotypes.NewCoinFromInt64(200*linotypes.Decimals), nil).Maybe()

	suite.vote.On("ClaimInterest", suite.Ctx, mock.Anything).Return(nil).Maybe()

	suite.vm = NewValidatorManager(testValidatorKey, suite.ph, suite.vote, suite.global, suite.acc)
	suite.vm.InitGenesis(suite.Ctx)
	suite.ph.On("GetValidatorParam", mock.Anything).Return(&parammodel.ValidatorParam{
		ValidatorMinDeposit:            linotypes.NewCoinFromInt64(200000 * linotypes.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissCommit:              linotypes.NewCoinFromInt64(200 * linotypes.Decimals),
		PenaltyByzantine:               linotypes.NewCoinFromInt64(1000 * linotypes.Decimals),
		AbsentCommitLimitation:         int64(600), // 30min
		OncallSize:                     int64(3),
		StandbySize:                    int64(3),
		ValidatorRevokePendingSec:      int64(7 * 24 * 3600),
		OncallInflationWeight:          int64(2),
		StandbyInflationWeight:         int64(1),
		MaxVotedValidators:             int64(3),
		SlashLimitation:                int64(5),
	}, nil).Maybe()

}

func (suite *ValidatorManagerTestSuite) SetupValidatorAndVotes(m map[linotypes.AccountKey]linotypes.Coin) {
	for name, votes := range m {
		val := model.Validator{
			ABCIValidator: abci.Validator{
				Address: secp256k1.GenPrivKey().PubKey().Address(),
				Power:   0},
			Username:      name,
			ReceivedVotes: votes,
		}
		suite.vm.storage.SetValidator(suite.Ctx, name, &val)
	}
}
func (suite *ValidatorManagerTestSuite) TestAddValidatortToOncallList() {
	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		prevList    model.ValidatorList
		prevVal     model.Validator
		expectList  model.ValidatorList
		expectPower int64
	}{
		{
			testName: "add user to oncall",
			username: linotypes.AccountKey("test"),
			prevVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: secp256k1.GenPrivKey().PubKey().Address(),
					Power:   0},
				Username:      linotypes.AccountKey("test"),
				ReceivedVotes: linotypes.NewCoinFromInt64(100000000),
			},
			prevList: model.ValidatorList{
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Oncall:             []linotypes.AccountKey{linotypes.AccountKey("test")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectPower: 1000,
		},
		{
			testName: "add user to oncall2",
			username: linotypes.AccountKey("test"),
			prevVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: secp256k1.GenPrivKey().PubKey().Address(),
					Power:   0},
				Username:      linotypes.AccountKey("test"),
				ReceivedVotes: linotypes.NewCoinFromInt64(900000000000 * linotypes.Decimals),
			},
			prevList: model.ValidatorList{
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Oncall:             []linotypes.AccountKey{linotypes.AccountKey("test")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectPower: 100000000000,
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidator(suite.Ctx, tc.username, &tc.prevVal)
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.addValidatortToOncallList(suite.Ctx, tc.username)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
		val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.username)
		suite.NoError(err)
		suite.Equal(tc.expectPower, val.ABCIValidator.Power, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestAddValidatortToStandbyList() {
	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		prevList    model.ValidatorList
		prevVal     model.Validator
		expectList  model.ValidatorList
		expectPower int64
	}{
		{
			testName: "add user to standby",
			username: linotypes.AccountKey("test"),
			prevVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: secp256k1.GenPrivKey().PubKey().Address(),
					Power:   1000},
				Username:      linotypes.AccountKey("test"),
				ReceivedVotes: linotypes.NewCoinFromInt64(100000000),
			},
			prevList: model.ValidatorList{
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Standby:            []linotypes.AccountKey{linotypes.AccountKey("test")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectPower: 1,
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidator(suite.Ctx, tc.username, &tc.prevVal)
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.addValidatortToStandbyList(suite.Ctx, tc.username)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
		val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.username)
		suite.NoError(err)
		suite.Equal(tc.expectPower, val.ABCIValidator.Power, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestAddValidatortToCandidateList() {
	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		prevList    model.ValidatorList
		prevVal     model.Validator
		expectList  model.ValidatorList
		expectPower int64
	}{
		{
			testName: "add user to candidate",
			username: linotypes.AccountKey("test"),
			prevVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: secp256k1.GenPrivKey().PubKey().Address(),
					Power:   1000},
				Username:      linotypes.AccountKey("test"),
				ReceivedVotes: linotypes.NewCoinFromInt64(100000000),
			},
			prevList: model.ValidatorList{
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Candidates:         []linotypes.AccountKey{linotypes.AccountKey("test")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectPower: 0,
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidator(suite.Ctx, tc.username, &tc.prevVal)
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.addValidatortToCandidateList(suite.Ctx, tc.username)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
		val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.username)
		suite.NoError(err)
		suite.Equal(tc.expectPower, val.ABCIValidator.Power, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestRmValidatortFromCandidateList() {
	testCases := []struct {
		testName   string
		username   linotypes.AccountKey
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "rm user from candidate",
			username: linotypes.AccountKey("test1"),
			prevList: model.ValidatorList{
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Candidates:         []linotypes.AccountKey{linotypes.AccountKey("test2")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		suite.vm.removeValidatorFromCandidateList(suite.Ctx, tc.username)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestRmValidatortFromOncallList() {
	testCases := []struct {
		testName   string
		username   linotypes.AccountKey
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "rm user from oncall",
			username: linotypes.AccountKey("test1"),
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Oncall:             []linotypes.AccountKey{linotypes.AccountKey("test2")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		suite.vm.removeValidatorFromOncallList(suite.Ctx, tc.username)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestRmValidatortFromStandbyList() {
	testCases := []struct {
		testName   string
		username   linotypes.AccountKey
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "rm user from standby",
			username: linotypes.AccountKey("test1"),
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Standby:            []linotypes.AccountKey{linotypes.AccountKey("test2")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		suite.vm.removeValidatorFromStandbyList(suite.Ctx, tc.username)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestRmValidatortFromAllList() {
	testCases := []struct {
		testName   string
		username   linotypes.AccountKey
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "rm user from all list",
			username: linotypes.AccountKey("test1"),
			prevList: model.ValidatorList{
				Standby:            []linotypes.AccountKey{linotypes.AccountKey("test1")},
				Oncall:             []linotypes.AccountKey{linotypes.AccountKey("test2")},
				Candidates:         []linotypes.AccountKey{linotypes.AccountKey("test3")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Oncall:             []linotypes.AccountKey{linotypes.AccountKey("test2")},
				Candidates:         []linotypes.AccountKey{linotypes.AccountKey("test3")},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.removeValidatorFromAllLists(suite.Ctx, tc.username)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestUpdateLowestOncall() {
	testCases := []struct {
		testName   string
		validators map[linotypes.AccountKey]linotypes.Coin
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "update lowest oncall",
			validators: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(1000),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(100),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(2000),
			},
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test2"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
		},
		{
			testName: "update lowest oncall2",
			validators: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(10000000),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(100),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(1),
			},
			prevList: model.ValidatorList{
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test2"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		suite.SetupValidatorAndVotes(tc.validators)
		err := suite.vm.updateLowestOncall(suite.Ctx)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestUpdateLowestStandby() {
	testCases := []struct {
		testName   string
		validators map[linotypes.AccountKey]linotypes.Coin
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "update lowest standby",
			validators: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(10000000),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(100),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(1),
			},
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(1),
				LowestStandby:      linotypes.AccountKey("test3"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
		{
			testName: "update lowest standby2",
			validators: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(10000000),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(100),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(1),
			},
			prevList: model.ValidatorList{
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(100),
				LowestStandby:      linotypes.AccountKey("test2"),
			},
			expectList: model.ValidatorList{
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		suite.SetupValidatorAndVotes(tc.validators)
		err := suite.vm.updateLowestStandby(suite.Ctx)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.NoError(err)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestGetHighestVotesAndValidator() {
	testCases := []struct {
		testName      string
		validators    map[linotypes.AccountKey]linotypes.Coin
		lst           []linotypes.AccountKey
		expectValName linotypes.AccountKey
		expectVotes   linotypes.Coin
	}{
		{
			testName: "get highest votes and val",
			validators: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(10000000),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(100),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(1),
			},
			lst: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test2"),
				linotypes.AccountKey("test3"),
			},
			expectValName: linotypes.AccountKey("test1"),
			expectVotes:   linotypes.NewCoinFromInt64(10000000),
		},
		{
			testName: "get highest votes and val2",
			validators: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(0),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(0),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(0),
			},
			lst: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test2"),
				linotypes.AccountKey("test3"),
			},
			expectValName: linotypes.AccountKey("test3"),
			expectVotes:   linotypes.NewCoinFromInt64(0),
		},
	}

	for _, tc := range testCases {
		suite.SetupValidatorAndVotes(tc.validators)
		name, votes, err := suite.vm.getHighestVotesAndValidator(suite.Ctx, tc.lst)
		suite.Require().Nil(err)
		suite.Equal(tc.expectValName, name, "%s", tc.testName)
		suite.Equal(tc.expectVotes, votes, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestGetLowestVotesAndValidator() {
	testCases := []struct {
		testName      string
		validators    map[linotypes.AccountKey]linotypes.Coin
		lst           []linotypes.AccountKey
		expectValName linotypes.AccountKey
		expectVotes   linotypes.Coin
	}{
		{
			testName: "get lowest votes and val",
			validators: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(10000000),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(100),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(1),
			},
			lst: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test2"),
				linotypes.AccountKey("test3"),
			},
			expectValName: linotypes.AccountKey("test3"),
			expectVotes:   linotypes.NewCoinFromInt64(1),
		},
		{
			testName: "get lowest votes and val2",
			validators: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(math.MaxInt64),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(math.MaxInt64),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(math.MaxInt64),
			},
			lst: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test2"),
				linotypes.AccountKey("test3"),
			},
			expectValName: linotypes.AccountKey("test3"),
			expectVotes:   linotypes.NewCoinFromInt64(math.MaxInt64),
		},
	}

	for _, tc := range testCases {
		suite.SetupValidatorAndVotes(tc.validators)
		name, votes, err := suite.vm.getLowestVotesAndValidator(suite.Ctx, tc.lst)
		suite.Require().Nil(err)
		suite.Equal(tc.expectValName, name, "%s", tc.testName)
		suite.Equal(tc.expectVotes, votes, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestRemoveExtraOncall() {
	testCases := []struct {
		testName   string
		validators map[linotypes.AccountKey]linotypes.Coin
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "rm extra oncall",
			validators: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(6),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(5),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(3),
				linotypes.AccountKey("test4"): linotypes.NewCoinFromInt64(4),
				linotypes.AccountKey("test5"): linotypes.NewCoinFromInt64(2),
				linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(1),
			},
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
				},
				Standby: []linotypes.AccountKey{

					linotypes.AccountKey("test6"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test4"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test3"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		suite.SetupValidatorAndVotes(tc.validators)
		err := suite.vm.removeExtraOncall(suite.Ctx)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.NoError(err)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestRemoveExtraStandby() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(6),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(5),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(3),
		linotypes.AccountKey("test4"): linotypes.NewCoinFromInt64(4),
		linotypes.AccountKey("test5"): linotypes.NewCoinFromInt64(2),
		linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(1),
	}
	suite.SetupValidatorAndVotes(validators)
	testCases := []struct {
		testName   string
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "rm extra standby",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test4"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test3"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.removeExtraStandby(suite.Ctx)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.NoError(err)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestFillEmptyStandby() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(6),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(5),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(3),
		linotypes.AccountKey("test4"): linotypes.NewCoinFromInt64(4),
		linotypes.AccountKey("test5"): linotypes.NewCoinFromInt64(2),
		linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(1),
	}
	suite.SetupValidatorAndVotes(validators)
	testCases := []struct {
		testName   string
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "fill empty standby",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test4"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
		{
			testName: "fill empty standby2",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.fillEmptyStandby(suite.Ctx)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.NoError(err)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestFillEmptyOncall() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(6),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(5),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(3),
		linotypes.AccountKey("test4"): linotypes.NewCoinFromInt64(4),
		linotypes.AccountKey("test5"): linotypes.NewCoinFromInt64(2),
		linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(1),
	}
	suite.SetupValidatorAndVotes(validators)
	testCases := []struct {
		testName   string
		prevList   model.ValidatorList
		expectList model.ValidatorList
	}{
		{
			testName: "fill empty standby",
			prevList: model.ValidatorList{
				Oncall:  []linotypes.AccountKey{},
				Standby: []linotypes.AccountKey{linotypes.AccountKey("test1")},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test4"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
		{
			testName: "fill empty oncall2",
			prevList: model.ValidatorList{
				Oncall:  []linotypes.AccountKey{},
				Standby: []linotypes.AccountKey{},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
		{
			testName: "fill empty oncall3",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test4"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.fillEmptyOncall(suite.Ctx)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestGetAllValidators() {
	testCases := []struct {
		testName  string
		prevList  model.ValidatorList
		expectRes []linotypes.AccountKey
	}{
		{
			testName: "get all validators",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectRes: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test2"),
				linotypes.AccountKey("test3"),
				linotypes.AccountKey("test4"),
				linotypes.AccountKey("test5"),
				linotypes.AccountKey("test6"),
				linotypes.AccountKey("test7"),
			},
		},
		{
			testName: "get all validators2",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectRes: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test4"),
			},
		},
		{
			testName: "get all validators3",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectRes: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test4"),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		lst := suite.vm.GetAllValidators(suite.Ctx)
		suite.Equal(tc.expectRes, lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestGetCommittingValidators() {
	testCases := []struct {
		testName  string
		prevList  model.ValidatorList
		expectRes []linotypes.AccountKey
	}{
		{
			testName: "get committing validators",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectRes: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test2"),
				linotypes.AccountKey("test3"),
				linotypes.AccountKey("test4"),
				linotypes.AccountKey("test5"),
				linotypes.AccountKey("test6"),
			},
		},
		{
			testName: "get committing validators2",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectRes: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		lst := suite.vm.GetCommittingValidators(suite.Ctx)
		suite.Equal(tc.expectRes, lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestOnCandidateVotesInc() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(100),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(200),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(300),
		linotypes.AccountKey("test4"): linotypes.NewCoinFromInt64(400),
		linotypes.AccountKey("test5"): linotypes.NewCoinFromInt64(500),
		linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(600),
		linotypes.AccountKey("test7"): linotypes.NewCoinFromInt64(700),
	}
	suite.SetupValidatorAndVotes(validators)

	testCases := []struct {
		testName      string
		prevList      model.ValidatorList
		increasedUser linotypes.AccountKey
		expectList    model.ValidatorList
	}{
		{
			testName: "on candidate votes inc",
			prevList: model.ValidatorList{
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			increasedUser: linotypes.AccountKey("test1"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test1"),
			},
		},
		{
			testName: "on candidate votes inc2",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(300),
				LowestOncall:       linotypes.AccountKey("test3"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(100),
				LowestStandby:      linotypes.AccountKey("test1"),
			},
			increasedUser: linotypes.AccountKey("test6"),
			expectList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(100),
				LowestStandby:      linotypes.AccountKey("test1"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(400),
				LowestOncall:       linotypes.AccountKey("test4"),
			},
		},
		{
			testName: "on candidate votes inc3",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(400),
				LowestOncall:       linotypes.AccountKey("test4"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(100),
				LowestStandby:      linotypes.AccountKey("test1"),
			},
			increasedUser: linotypes.AccountKey("test7"),
			expectList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
				},
				Oncall: []linotypes.AccountKey{

					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test7"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200),
				LowestStandby:      linotypes.AccountKey("test2"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
			},
		},
		{
			testName: "on candidate votes inc4",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(100),
				LowestStandby:      linotypes.AccountKey("test1"),
			},
			increasedUser: linotypes.AccountKey("test2"),
			expectList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test2"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200),
				LowestStandby:      linotypes.AccountKey("test2"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
			},
		},
		{
			testName: "on candidate votes inc5",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(600),
				LowestOncall:       linotypes.AccountKey("test6"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(100),
				LowestStandby:      linotypes.AccountKey("test1"),
			},
			increasedUser: linotypes.AccountKey("test3"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test1"),
			},
		},
		{
			testName: "on candidate votes inc6",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test7"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200),
				LowestStandby:      linotypes.AccountKey("test2"),
			},
			increasedUser: linotypes.AccountKey("test1"),
			expectList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test7"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200),
				LowestStandby:      linotypes.AccountKey("test2"),
			},
		},
		{
			testName: "on candidate votes inc7",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test5"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(600),
				LowestOncall:       linotypes.AccountKey("test6"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(500),
				LowestStandby:      linotypes.AccountKey("test5"),
			},
			increasedUser: linotypes.AccountKey("test3"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test3"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(300),
				LowestOncall:       linotypes.AccountKey("test3"),
			},
		},
		{
			testName: "on candidate votes inc8",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test3"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(300),
				LowestStandby:      linotypes.AccountKey("test3"),
			},
			increasedUser: linotypes.AccountKey("test1"),
			expectList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test1"),
				},
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(100),
				LowestStandby:      linotypes.AccountKey("test1"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
			},
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.onCandidateVotesInc(suite.Ctx, tc.increasedUser)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestOnStandbyVotesInc() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(100),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(200),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(300),
		linotypes.AccountKey("test4"): linotypes.NewCoinFromInt64(400),
		linotypes.AccountKey("test5"): linotypes.NewCoinFromInt64(500),
		linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(600),
		linotypes.AccountKey("test7"): linotypes.NewCoinFromInt64(700),
	}
	suite.SetupValidatorAndVotes(validators)

	testCases := []struct {
		testName      string
		prevList      model.ValidatorList
		increasedUser linotypes.AccountKey
		expectList    model.ValidatorList
	}{
		{
			testName: "on standby votes inc",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			increasedUser: linotypes.AccountKey("test1"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test1"),
			},
		},
		{
			testName: "on standby votes inc2",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test4"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(400),
				LowestOncall:       linotypes.AccountKey("test4"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			increasedUser: linotypes.AccountKey("test7"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test7"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(400),
				LowestStandby:      linotypes.AccountKey("test4"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
			},
		},
		{
			testName: "on standby votes inc3",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test4"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(400),
				LowestOncall:       linotypes.AccountKey("test4"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			increasedUser: linotypes.AccountKey("test3"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test4"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(300),
				LowestStandby:      linotypes.AccountKey("test3"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(400),
				LowestOncall:       linotypes.AccountKey("test4"),
			},
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.onStandbyVotesInc(suite.Ctx, tc.increasedUser)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestOnOncallVotesInc() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(600 * linotypes.Decimals),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(300),
	}
	suite.SetupValidatorAndVotes(validators)

	testCases := []struct {
		testName      string
		prevList      model.ValidatorList
		increasedUser linotypes.AccountKey
		expectList    model.ValidatorList
		expectPower   int64
	}{
		{
			testName: "on oncall votes inc",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			increasedUser: linotypes.AccountKey("test6"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(600 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test6"),
			},
			expectPower: 600,
		},
		{
			testName: "on oncall votes inc2",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(600 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test6"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(300),
				LowestStandby:      linotypes.AccountKey("test3"),
			},
			increasedUser: linotypes.AccountKey("test6"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test3"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(300),
				LowestOncall:       linotypes.AccountKey("test3"),
			},
			expectPower: 600,
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.onOncallVotesInc(suite.Ctx, tc.increasedUser)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.NoError(err)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
		val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.increasedUser)
		suite.NoError(err)
		suite.Equal(tc.expectPower, val.ABCIValidator.Power, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestCheckDupPubKey() {
	key1 := secp256k1.GenPrivKey().PubKey()
	key2 := secp256k1.GenPrivKey().PubKey()
	testCases := []struct {
		testName    string
		newKey      crypto.PubKey
		existVal    model.Validator
		prevList    model.ValidatorList
		expectedRes sdk.Error
	}{
		{
			testName: "check dup pubkey",
			newKey:   key1,
			existVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: key2.Address(),
					Power:   0},
				Username: linotypes.AccountKey("test1"),
			},
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectedRes: nil,
		},
		{
			testName: "check dup pubkey2",
			newKey:   key2,
			existVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: key2.Address(),
					Power:   0},
				Username: linotypes.AccountKey("test1"),
			},
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectedRes: types.ErrValidatorPubKeyAlreadyExist(),
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		suite.vm.storage.SetValidator(suite.Ctx, tc.existVal.Username, &tc.existVal)
		err := suite.vm.checkDupPubKey(suite.Ctx, tc.newKey)
		suite.Equal(tc.expectedRes, err, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestGetElectionVoteListUpdates() {
	testCases := []struct {
		testName        string
		username        linotypes.AccountKey
		votedValidators []linotypes.AccountKey
		prevList        model.ElectionVoteList
		expectedUpdates []*model.ElectionVote
	}{
		{
			testName: "get election vote list updates",
			username: linotypes.AccountKey("user1"),
			votedValidators: []linotypes.AccountKey{
				linotypes.AccountKey("val4"),
				linotypes.AccountKey("val5"),
				linotypes.AccountKey("val6"),
			},
			prevList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("val1"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
					{
						ValidatorName: linotypes.AccountKey("val2"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
					{
						ValidatorName: linotypes.AccountKey("val3"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
				},
			},
			expectedUpdates: []*model.ElectionVote{
				{
					ValidatorName: linotypes.AccountKey("val1"),
					Vote:          linotypes.NewCoinFromInt64(-100),
				},
				{
					ValidatorName: linotypes.AccountKey("val2"),
					Vote:          linotypes.NewCoinFromInt64(-100),
				},
				{
					ValidatorName: linotypes.AccountKey("val3"),
					Vote:          linotypes.NewCoinFromInt64(-100),
				},
				{
					ValidatorName: linotypes.AccountKey("val4"),
					Vote:          linotypes.NewCoinFromInt64(100),
				},
				{
					ValidatorName: linotypes.AccountKey("val5"),
					Vote:          linotypes.NewCoinFromInt64(100),
				},
				{
					ValidatorName: linotypes.AccountKey("val6"),
					Vote:          linotypes.NewCoinFromInt64(100),
				},
			},
		},
		{
			testName: "get election vote list updates2",
			username: linotypes.AccountKey("user1"),
			votedValidators: []linotypes.AccountKey{
				linotypes.AccountKey("val1"),
			},
			prevList: model.ElectionVoteList{},
			expectedUpdates: []*model.ElectionVote{
				{
					ValidatorName: linotypes.AccountKey("val1"),
					Vote:          linotypes.NewCoinFromInt64(300),
				},
			},
		},
		{
			testName: "get election vote list updates3",
			username: linotypes.AccountKey("user1"),
			votedValidators: []linotypes.AccountKey{
				linotypes.AccountKey("val1"),
				linotypes.AccountKey("val2"),
			},
			prevList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("val1"),
						Vote:          linotypes.NewCoinFromInt64(300),
					},
				},
			},
			expectedUpdates: []*model.ElectionVote{
				{
					ValidatorName: linotypes.AccountKey("val1"),
					Vote:          linotypes.NewCoinFromInt64(-150),
				},
				{
					ValidatorName: linotypes.AccountKey("val2"),
					Vote:          linotypes.NewCoinFromInt64(150),
				},
			},
		},
		{
			testName: "get election vote list updates4",
			username: linotypes.AccountKey("user1"),
			votedValidators: []linotypes.AccountKey{
				linotypes.AccountKey("val1"),
				linotypes.AccountKey("val2"),
			},
			prevList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("val1"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
					{
						ValidatorName: linotypes.AccountKey("val2"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
				},
			},
			expectedUpdates: []*model.ElectionVote{
				{
					ValidatorName: linotypes.AccountKey("val1"),
					Vote:          linotypes.NewCoinFromInt64(50),
				},
				{
					ValidatorName: linotypes.AccountKey("val2"),
					Vote:          linotypes.NewCoinFromInt64(50),
				},
			},
		},
		{
			testName: "get election vote list updates5",
			username: linotypes.AccountKey("user1"),
			votedValidators: []linotypes.AccountKey{
				linotypes.AccountKey("val1"),
			},
			prevList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("val1"),
						Vote:          linotypes.NewCoinFromInt64(300),
					},
				},
			},
			expectedUpdates: []*model.ElectionVote{
				{
					ValidatorName: linotypes.AccountKey("val1"),
					Vote: linotypes.NewCoinFromInt64(300).Plus(
						linotypes.NewCoinFromInt64(300).Neg()),
				},
			},
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetElectionVoteList(suite.Ctx, tc.username, &tc.prevList)
		updates, err := suite.vm.getElectionVoteListUpdates(suite.Ctx, tc.username, tc.votedValidators)
		suite.NoError(err)
		suite.Equal(tc.expectedUpdates, updates, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestSetNewElectionVoteList() {
	testCases := []struct {
		testName        string
		username        linotypes.AccountKey
		votedValidators []linotypes.AccountKey
		prevList        model.ElectionVoteList
		expectedList    model.ElectionVoteList
	}{
		{
			testName: "set new election vote list",
			username: linotypes.AccountKey("user1"),
			votedValidators: []linotypes.AccountKey{
				linotypes.AccountKey("val1"),
				linotypes.AccountKey("val2"),
			},
			prevList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("val6"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
				},
			},
			expectedList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("val1"),
						Vote:          linotypes.NewCoinFromInt64(150),
					},
					{
						ValidatorName: linotypes.AccountKey("val2"),
						Vote:          linotypes.NewCoinFromInt64(150),
					},
				},
			},
		},
		{
			testName:        "set new election vote list2",
			username:        linotypes.AccountKey("user1"),
			votedValidators: []linotypes.AccountKey{},
			prevList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("val6"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
				},
			},
			expectedList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("val6"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetElectionVoteList(suite.Ctx, tc.username, &tc.prevList)
		err := suite.vm.setNewElectionVoteList(suite.Ctx, tc.username, tc.votedValidators)
		suite.NoError(err)
		lst := suite.vm.storage.GetElectionVoteList(suite.Ctx, tc.username)
		suite.Equal(tc.expectedList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestOnStandbyVotesDec() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(100),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(200),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(300),
		linotypes.AccountKey("test4"): linotypes.NewCoinFromInt64(400),
		linotypes.AccountKey("test5"): linotypes.NewCoinFromInt64(500),
		linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(600),
		linotypes.AccountKey("test7"): linotypes.NewCoinFromInt64(700),
	}
	suite.SetupValidatorAndVotes(validators)

	testCases := []struct {
		testName      string
		prevList      model.ValidatorList
		decreasedUser linotypes.AccountKey
		expectList    model.ValidatorList
	}{
		{
			testName: "on standby votes dec1",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(300),
				LowestStandby:      linotypes.AccountKey("test3"),
			},
			decreasedUser: linotypes.AccountKey("test2"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test2"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200),
				LowestOncall:       linotypes.AccountKey("test2"),
			},
		},
		{
			testName: "on standby votes dec2",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test1"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(300),
				LowestStandby:      linotypes.AccountKey("test3"),
			},
			decreasedUser: linotypes.AccountKey("test1"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200),
				LowestStandby:      linotypes.AccountKey("test2"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
			},
		},
		{
			testName: "on standby votes dec3",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(300),
				LowestStandby:      linotypes.AccountKey("test3"),
			},
			decreasedUser: linotypes.AccountKey("test2"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200),
				LowestStandby:      linotypes.AccountKey("test2"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500),
				LowestOncall:       linotypes.AccountKey("test5"),
			},
		},
		{
			testName: "on standby votes dec4",
			prevList: model.ValidatorList{
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(300),
				LowestStandby:      linotypes.AccountKey("test3"),
			},
			decreasedUser: linotypes.AccountKey("test3"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test4"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200),
				LowestOncall:       linotypes.AccountKey("test2"),
			},
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.onStandbyVotesDec(suite.Ctx, tc.decreasedUser)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestOnOncallVotesDec() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(100 * linotypes.Decimals),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(200 * linotypes.Decimals),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(300 * linotypes.Decimals),
		linotypes.AccountKey("test4"): linotypes.NewCoinFromInt64(400 * linotypes.Decimals),
		linotypes.AccountKey("test5"): linotypes.NewCoinFromInt64(500 * linotypes.Decimals),
		linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(600 * linotypes.Decimals),
		linotypes.AccountKey("test7"): linotypes.NewCoinFromInt64(700 * linotypes.Decimals),
	}
	suite.SetupValidatorAndVotes(validators)

	testCases := []struct {
		testName      string
		prevList      model.ValidatorList
		decreasedUser linotypes.AccountKey
		expectList    model.ValidatorList
		expectPower   int64
	}{
		{
			testName: "on oncall votes dec",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test1"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test3"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(600 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test6"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200 * linotypes.Decimals),
				LowestStandby:      linotypes.AccountKey("test2"),
			},
			decreasedUser: linotypes.AccountKey("test1"),
			expectPower:   0,
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200 * linotypes.Decimals),
				LowestStandby:      linotypes.AccountKey("test2"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test5"),
			},
		},
		{
			testName: "on oncall votes dec2",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(300 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test2"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			decreasedUser: linotypes.AccountKey("test2"),
			expectPower:   200,
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test2"),
			},
		},
		{
			testName: "on oncall votes dec3",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test3"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test2"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(600 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test6"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200 * linotypes.Decimals),
				LowestStandby:      linotypes.AccountKey("test2"),
			},
			decreasedUser: linotypes.AccountKey("test3"),
			expectPower:   1,
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
					linotypes.AccountKey("test6"),
					linotypes.AccountKey("test5"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(200 * linotypes.Decimals),
				LowestStandby:      linotypes.AccountKey("test2"),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(500 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test5"),
			},
		},
		{
			testName: "on oncall votes dec4",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test1"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			decreasedUser: linotypes.AccountKey("test2"),
			expectPower:   200,
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test1"),
			},
		},
		{
			testName: "on oncall votes dec5",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test2"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			decreasedUser: linotypes.AccountKey("test2"),
			expectPower:   200,
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100 * linotypes.Decimals),
				LowestOncall:       linotypes.AccountKey("test1"),
			},
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.onOncallVotesDec(suite.Ctx, tc.decreasedUser)
		suite.Require().Nil(err)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
		val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.decreasedUser)
		suite.NoError(err)
		suite.Equal(tc.expectPower, val.ABCIValidator.Power, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestGetValidatorUpdates() {
	valKey1 := secp256k1.GenPrivKey().PubKey()
	valKey2 := secp256k1.GenPrivKey().PubKey()

	user1 := linotypes.AccountKey("user1")
	user2 := linotypes.AccountKey("user2")
	validator1 := model.Validator{
		ABCIValidator: abci.Validator{
			Address: valKey1.Address(),
			Power:   linotypes.TendermintValidatorPower,
		},
		PubKey:        valKey1,
		Username:      user1,
		ReceivedVotes: linotypes.NewCoinFromInt64(0),
	}
	validator2 := model.Validator{
		ABCIValidator: abci.Validator{
			Address: valKey2.Address(),
			Power:   linotypes.TendermintValidatorPower,
		},
		PubKey:        valKey2,
		Username:      user2,
		ReceivedVotes: linotypes.NewCoinFromInt64(0),
	}
	suite.vm.storage.SetValidator(suite.Ctx, user1, &validator1)
	suite.vm.storage.SetValidator(suite.Ctx, user2, &validator2)

	val1 := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey1),
		Power:  linotypes.TendermintValidatorPower,
	}

	val2 := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey2),
		Power:  linotypes.TendermintValidatorPower,
	}

	val1NoPower := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey1),
		Power:  0,
	}

	val2NoPower := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey2),
		Power:  0,
	}

	testCases := []struct {
		testName            string
		oncallValidators    []linotypes.AccountKey
		preBlockValidators  []linotypes.AccountKey
		expectedUpdatedList []abci.ValidatorUpdate
	}{
		{
			testName:            "only one oncall validator",
			oncallValidators:    []linotypes.AccountKey{user1},
			preBlockValidators:  []linotypes.AccountKey{},
			expectedUpdatedList: []abci.ValidatorUpdate{val1},
		},
		{
			testName:            "two oncall validators and one pre block validator",
			oncallValidators:    []linotypes.AccountKey{user1, user2},
			preBlockValidators:  []linotypes.AccountKey{user1},
			expectedUpdatedList: []abci.ValidatorUpdate{val1, val2},
		},
		{
			testName:            "two oncall validatos and two pre block validators",
			oncallValidators:    []linotypes.AccountKey{user1, user2},
			preBlockValidators:  []linotypes.AccountKey{user1, user2},
			expectedUpdatedList: []abci.ValidatorUpdate{val1, val2},
		},
		{
			testName:            "one oncall validator and two pre block validators",
			oncallValidators:    []linotypes.AccountKey{user2},
			preBlockValidators:  []linotypes.AccountKey{user1, user2},
			expectedUpdatedList: []abci.ValidatorUpdate{val1NoPower, val2},
		},
		{
			testName:            "only one pre block validator",
			oncallValidators:    []linotypes.AccountKey{},
			preBlockValidators:  []linotypes.AccountKey{user2},
			expectedUpdatedList: []abci.ValidatorUpdate{val2NoPower},
		},
	}

	for _, tc := range testCases {
		lst := &model.ValidatorList{
			Oncall:             tc.oncallValidators,
			PreBlockValidators: tc.preBlockValidators,
		}
		suite.vm.storage.SetValidatorList(suite.Ctx, lst)

		actualList, err := suite.vm.GetValidatorUpdates(suite.Ctx)
		suite.NoError(err)
		suite.Equal(tc.expectedUpdatedList, actualList, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestRejoinFromJail() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("jail1"): linotypes.NewCoinFromInt64(100),
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(200),
	}
	suite.SetupValidatorAndVotes(validators)

	testCases := []struct {
		testName   string
		prevList   model.ValidatorList
		rejoinUser linotypes.AccountKey
		expectList model.ValidatorList
		expectRes  sdk.Error
	}{
		{
			testName: "rejoin from jail",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				Jail: []linotypes.AccountKey{
					linotypes.AccountKey("jail1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200),
				LowestOncall:       linotypes.AccountKey("test1"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			rejoinUser: linotypes.AccountKey("jail1"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("jail1"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("jail1"),
			},
			expectRes: nil,
		},
		{
			testName: "rejoin from jail2",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				Jail: []linotypes.AccountKey{
					linotypes.AccountKey("jail2"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200),
				LowestOncall:       linotypes.AccountKey("test1"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			rejoinUser: linotypes.AccountKey("jail2"),
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				Jail: []linotypes.AccountKey{
					linotypes.AccountKey("jail2"),
				},
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200),
				LowestOncall:       linotypes.AccountKey("test1"),
			},
			expectRes: types.ErrInsufficientDeposit(),
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.rejoinFromJail(suite.Ctx, tc.rejoinUser)
		suite.Equal(tc.expectRes, err, "%s", tc.testName)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestVoteValidator() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(100),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(200),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(300),
	}
	suite.SetupValidatorAndVotes(validators)

	testCases := []struct {
		testName           string
		prevValList        model.ValidatorList
		voter              linotypes.AccountKey
		votedValidators    []linotypes.AccountKey
		expectElectionList model.ElectionVoteList
		expectValAndVotes  map[linotypes.AccountKey]linotypes.Coin
		expectRes          sdk.Error
	}{
		{
			testName: "vote validator",
			prevValList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test1"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			voter: linotypes.AccountKey("user1"),
			votedValidators: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test2"),
				linotypes.AccountKey("test3"),
			},
			expectElectionList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("test1"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
					{
						ValidatorName: linotypes.AccountKey("test2"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
					{
						ValidatorName: linotypes.AccountKey("test3"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
				},
			},
			expectValAndVotes: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(200),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(300),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(400),
			},
			expectRes: nil,
		},
		{
			testName: "vote validator2",
			prevValList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test1"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			voter: linotypes.AccountKey("user1"),
			votedValidators: []linotypes.AccountKey{
				linotypes.AccountKey("test1"),
				linotypes.AccountKey("test2"),
				linotypes.AccountKey("dummy"),
			},
			expectElectionList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("test1"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
					{
						ValidatorName: linotypes.AccountKey("test2"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
					{
						ValidatorName: linotypes.AccountKey("test3"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
				},
			},
			expectValAndVotes: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(200),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(300),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(400),
			},
			expectRes: types.ErrValidatorNotFound(linotypes.AccountKey("dummy")),
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevValList)
		err := suite.vm.VoteValidator(suite.Ctx, tc.voter, tc.votedValidators)
		suite.Equal(tc.expectRes, err, "%s", tc.testName)
		lst := suite.vm.storage.GetElectionVoteList(suite.Ctx, tc.voter)
		suite.Equal(tc.expectElectionList, *lst, "%s", tc.testName)
		for k, v := range tc.expectValAndVotes {
			val, _ := suite.vm.storage.GetValidator(suite.Ctx, k)
			suite.Equal(v, val.ReceivedVotes, "%s", tc.testName)
		}
	}
}

func (suite *ValidatorManagerTestSuite) TestGetInitValidators() {
	valKey1 := secp256k1.GenPrivKey().PubKey()
	valKey2 := secp256k1.GenPrivKey().PubKey()

	user1 := linotypes.AccountKey("user1")
	user2 := linotypes.AccountKey("user2")
	validator1 := model.Validator{
		ABCIValidator: abci.Validator{
			Address: valKey1.Address(),
			Power:   linotypes.TendermintValidatorPower,
		},
		PubKey:        valKey1,
		Username:      user1,
		ReceivedVotes: linotypes.NewCoinFromInt64(0),
	}
	validator2 := model.Validator{
		ABCIValidator: abci.Validator{
			Address: valKey2.Address(),
			Power:   linotypes.TendermintValidatorPower,
		},
		PubKey:        valKey2,
		Username:      user2,
		ReceivedVotes: linotypes.NewCoinFromInt64(0),
	}
	suite.vm.storage.SetValidator(suite.Ctx, user1, &validator1)
	suite.vm.storage.SetValidator(suite.Ctx, user2, &validator2)

	val1 := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey1),
		Power:  linotypes.TendermintValidatorPower,
	}

	val2 := abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(valKey2),
		Power:  linotypes.TendermintValidatorPower,
	}

	testCases := []struct {
		testName            string
		oncallValidators    []linotypes.AccountKey
		expectedUpdatedList []abci.ValidatorUpdate
	}{
		{
			testName:            "only one oncall validator",
			oncallValidators:    []linotypes.AccountKey{user1},
			expectedUpdatedList: []abci.ValidatorUpdate{val1},
		},
		{
			testName:            "two oncall validators",
			oncallValidators:    []linotypes.AccountKey{user1, user2},
			expectedUpdatedList: []abci.ValidatorUpdate{val1, val2},
		},
		{
			testName:            "no validators exists",
			oncallValidators:    []linotypes.AccountKey{},
			expectedUpdatedList: []abci.ValidatorUpdate{},
		},
	}

	for _, tc := range testCases {
		lst := &model.ValidatorList{
			Oncall: tc.oncallValidators,
		}
		suite.vm.storage.SetValidatorList(suite.Ctx, lst)

		actualList, err := suite.vm.GetInitValidators(suite.Ctx)
		suite.NoError(err)
		suite.Equal(tc.expectedUpdatedList, actualList, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestFireIncompetentValidator() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(100),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(200),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(300),
	}
	suite.SetupValidatorAndVotes(validators)

	byzKey := secp256k1.GenPrivKey().PubKey()
	byz := linotypes.AccountKey("byz")
	byzVal := model.Validator{
		ABCIValidator: abci.Validator{
			Address: byzKey.Address(),
			Power:   linotypes.TendermintValidatorPower,
		},
		PubKey:        byzKey,
		Username:      byz,
		ReceivedVotes: linotypes.NewCoinFromInt64(2000),
	}

	absKey := secp256k1.GenPrivKey().PubKey()
	abs := linotypes.AccountKey("abs")
	absVal := model.Validator{
		ABCIValidator: abci.Validator{
			Address: absKey.Address(),
			Power:   linotypes.TendermintValidatorPower,
		},
		PubKey:        absKey,
		Username:      abs,
		ReceivedVotes: linotypes.NewCoinFromInt64(2000),
		AbsentCommit:  20000,
	}
	suite.vm.storage.SetValidator(suite.Ctx, byz, &byzVal)
	suite.vm.storage.SetValidator(suite.Ctx, abs, &absVal)

	testCases := []struct {
		testName            string
		prevList            model.ValidatorList
		expectedList        model.ValidatorList
		byzantineValidators []abci.Evidence
	}{
		{
			testName: "fire validator",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("abs"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200),
				LowestOncall:       linotypes.AccountKey("test2"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(100),
				LowestStandby:      linotypes.AccountKey("test1"),
			},
			expectedList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test1"),
				},
				Jail: []linotypes.AccountKey{
					linotypes.AccountKey("abs"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test1"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			byzantineValidators: []abci.Evidence{},
		},
		{
			testName: "fire validator2",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("byz"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200),
				LowestOncall:       linotypes.AccountKey("test2"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(100),
				LowestStandby:      linotypes.AccountKey("test1"),
			},
			expectedList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test3"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test1"),
				},
				Jail: []linotypes.AccountKey{
					linotypes.AccountKey("byz"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test1"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			byzantineValidators: []abci.Evidence{
				{
					Validator: abci.Validator{
						Address: byzKey.Address(),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.fireIncompetentValidator(suite.Ctx, tc.byzantineValidators)
		suite.NoError(err)
		actualList := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectedList, *actualList, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestOnStakeChange() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(100),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(200),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(300),
	}
	suite.SetupValidatorAndVotes(validators)

	testCases := []struct {
		testName           string
		prevValList        model.ValidatorList
		voter              linotypes.AccountKey
		prevElectionList   model.ElectionVoteList
		expectElectionList model.ElectionVoteList
		expectValAndVotes  map[linotypes.AccountKey]linotypes.Coin
		expectRes          sdk.Error
	}{
		{
			testName: "vote validator",
			prevValList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(100),
				LowestOncall:       linotypes.AccountKey("test1"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			voter: linotypes.AccountKey("changedVoter"),
			prevElectionList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("test1"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
					{
						ValidatorName: linotypes.AccountKey("test2"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
					{
						ValidatorName: linotypes.AccountKey("test3"),
						Vote:          linotypes.NewCoinFromInt64(100),
					},
				},
			},
			expectElectionList: model.ElectionVoteList{
				ElectionVotes: []model.ElectionVote{
					{
						ValidatorName: linotypes.AccountKey("test1"),
						Vote:          linotypes.NewCoinFromInt64(200),
					},
					{
						ValidatorName: linotypes.AccountKey("test2"),
						Vote:          linotypes.NewCoinFromInt64(200),
					},
					{
						ValidatorName: linotypes.AccountKey("test3"),
						Vote:          linotypes.NewCoinFromInt64(200),
					},
				},
			},
			expectValAndVotes: map[linotypes.AccountKey]linotypes.Coin{
				linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(200),
				linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(300),
				linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(400),
			},
			expectRes: nil,
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevValList)
		suite.vm.storage.SetElectionVoteList(suite.Ctx, tc.voter, &tc.prevElectionList)
		err := suite.vm.onStakeChange(suite.Ctx, tc.voter)
		suite.Equal(tc.expectRes, err, "%s", tc.testName)
		lst := suite.vm.storage.GetElectionVoteList(suite.Ctx, tc.voter)
		suite.Equal(tc.expectElectionList, *lst, "%s", tc.testName)
		for k, v := range tc.expectValAndVotes {
			val, _ := suite.vm.storage.GetValidator(suite.Ctx, k)
			suite.Equal(v, val.ReceivedVotes, "%s", tc.testName)
		}
	}
}

func (suite *ValidatorManagerTestSuite) TestRegisterValidator() {
	valKey := secp256k1.GenPrivKey().PubKey()
	val := linotypes.AccountKey("val")

	testCases := []struct {
		testName   string
		username   linotypes.AccountKey
		link       string
		expectList model.ValidatorList
		expectVal  model.Validator
		expectRes  sdk.Error
	}{
		{
			testName: "vote validator",
			link:     "web1",
			expectList: model.ValidatorList{
				Oncall:             []linotypes.AccountKey{val},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(300),
				LowestOncall:       linotypes.AccountKey("val"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			username: val,
			expectVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: valKey.Address(),
					Power:   1,
				},
				Link:          "web1",
				PubKey:        valKey,
				Username:      val,
				ReceivedVotes: linotypes.NewCoinFromInt64(300),
			},
			expectRes: nil,
		},
	}
	for _, tc := range testCases {
		err := suite.vm.RegisterValidator(suite.Ctx, tc.username, valKey, tc.link)
		suite.Equal(tc.expectRes, err, "%s", tc.testName)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
		val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.username)
		suite.NoError(err)
		suite.Equal(tc.expectVal, *val, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestRegisterFromRevoked() {
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("valx")).Return(linotypes.NewCoinFromInt64(1), nil).Maybe()
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("valy")).Return(linotypes.NewCoinFromInt64(2), nil).Maybe()
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("valz")).Return(linotypes.NewCoinFromInt64(3), nil).Maybe()
	suite.vote.On("GetVoterDuty", suite.Ctx, linotypes.AccountKey("valx")).Return(votetypes.DutyVoter, nil).Maybe()
	suite.vote.On("GetVoterDuty", suite.Ctx, linotypes.AccountKey("valy")).Return(votetypes.DutyVoter, nil).Maybe()
	suite.vote.On("GetVoterDuty", suite.Ctx, linotypes.AccountKey("valz")).Return(votetypes.DutyVoter, nil).Maybe()
	suite.vote.On("AssignDuty", suite.Ctx, linotypes.AccountKey("valx"), votetypes.DutyValidator,
		linotypes.NewCoinFromInt64(200000*linotypes.Decimals)).Return(nil).Maybe()
	suite.vote.On("AssignDuty", suite.Ctx, linotypes.AccountKey("valy"), votetypes.DutyValidator,
		linotypes.NewCoinFromInt64(200000*linotypes.Decimals)).Return(nil).Maybe()
	suite.vote.On("AssignDuty", suite.Ctx, linotypes.AccountKey("valz"), votetypes.DutyValidator,
		linotypes.NewCoinFromInt64(200000*linotypes.Decimals)).Return(nil).Maybe()

	err := suite.vm.RegisterValidator(suite.Ctx, linotypes.AccountKey("valx"), secp256k1.GenPrivKey().PubKey(), "link")
	suite.NoError(err)
	err = suite.vm.RegisterValidator(suite.Ctx, linotypes.AccountKey("valy"), secp256k1.GenPrivKey().PubKey(), "link")
	suite.NoError(err)
	err = suite.vm.RegisterValidator(suite.Ctx, linotypes.AccountKey("valz"), secp256k1.GenPrivKey().PubKey(), "link")
	suite.NoError(err)

	valKey := secp256k1.GenPrivKey().PubKey()
	valName := linotypes.AccountKey("val")
	val := model.Validator{
		ABCIValidator: abci.Validator{
			Address: valKey.Address(),
			Power:   linotypes.TendermintValidatorPower,
		},
		PubKey:        valKey,
		Username:      valName,
		ReceivedVotes: linotypes.NewCoinFromInt64(300),
		HasRevoked:    true,
	}
	suite.vm.storage.SetValidator(suite.Ctx, valName, &val)
	suite.vm.storage.SetElectionVoteList(suite.Ctx, valName, &model.ElectionVoteList{
		ElectionVotes: []model.ElectionVote{
			{
				ValidatorName: linotypes.AccountKey("val"),
				Vote:          linotypes.NewCoinFromInt64(300),
			},
		},
	})

	testCases := []struct {
		testName   string
		username   linotypes.AccountKey
		link       string
		expectList model.ValidatorList
		expectVal  model.Validator
		expectRes  sdk.Error
	}{
		{
			testName: "register a revoked one",
			link:     "web1",
			expectList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("valy"),
					linotypes.AccountKey("valz"),
					linotypes.AccountKey("val"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("valx"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(2),
				LowestOncall:       linotypes.AccountKey("valy"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(1),
				LowestStandby:      linotypes.AccountKey("valx"),
			},
			username: valName,
			expectVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: valKey.Address(),
					Power:   1,
				},
				Link:          "web1",
				PubKey:        valKey,
				Username:      valName,
				ReceivedVotes: linotypes.NewCoinFromInt64(300),
			},
			expectRes: nil,
		},
	}
	for _, tc := range testCases {
		err := suite.vm.RegisterValidator(suite.Ctx, tc.username, valKey, tc.link)
		suite.Equal(tc.expectRes, err, "%s", tc.testName)
		lst := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectList, *lst, "%s", tc.testName)
		val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.username)
		suite.NoError(err)
		suite.Equal(tc.expectVal, *val, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestDistributeInflationToValidator() {
	suite.acc.On("GetPool", mock.Anything, linotypes.InflationValidatorPool).Return(
		(linotypes.NewCoinFromInt64(6)), nil).Once()
	for _, v := range []struct {
		validator linotypes.AccountKey
		amount    linotypes.Coin
	}{
		{
			linotypes.AccountKey("oncall1"),
			linotypes.NewCoinFromInt64(2),
		},
		{
			linotypes.AccountKey("oncall2"),
			linotypes.NewCoinFromInt64(2),
		},
		{
			linotypes.AccountKey("standby1"),
			linotypes.NewCoinFromInt64(1),
		},
		{
			linotypes.AccountKey("standby2"),
			linotypes.NewCoinFromInt64(1),
		},
	} {
		suite.acc.On(
			"MoveFromPool", mock.Anything, linotypes.InflationValidatorPool,
			linotypes.NewAccOrAddrFromAcc(v.validator), v.amount).Return(nil).Once()
	}

	testCases := []struct {
		testName string
		prevList model.ValidatorList
	}{
		{
			testName: "distribute inflation",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("oncall1"),
					linotypes.AccountKey("oncall2"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("standby1"),
					linotypes.AccountKey("standby2"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
		},
	}
	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.DistributeInflationToValidator(suite.Ctx)
		suite.NoError(err)
		suite.acc.AssertExpectations(suite.T())
	}
}

func (suite *ValidatorManagerTestSuite) TestRevokeValidator() {
	valKey := secp256k1.GenPrivKey().PubKey()
	val := linotypes.AccountKey("val")
	err := suite.vm.RegisterValidator(suite.Ctx, val, valKey, "link")
	suite.NoError(err)

	testCases := []struct {
		testName   string
		username   linotypes.AccountKey
		expectList model.ValidatorList
		expectVal  model.Validator
		expectRes  sdk.Error
	}{
		{
			testName: "revoke validator",
			expectList: model.ValidatorList{
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			username: val,
			expectVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: valKey.Address(),
					Power:   1,
				},
				Link:          "link",
				PubKey:        valKey,
				Username:      val,
				ReceivedVotes: linotypes.NewCoinFromInt64(300),
				HasRevoked:    true,
			},
			expectRes: nil,
		},
		{
			testName: "revoke validator2",
			expectList: model.ValidatorList{
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			username: val,
			expectVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: valKey.Address(),
					Power:   1,
				},
				Link:          "link",
				PubKey:        valKey,
				Username:      val,
				ReceivedVotes: linotypes.NewCoinFromInt64(300),
				HasRevoked:    true,
			},
			expectRes: types.ErrInvalidValidator(),
		},
	}
	for _, tc := range testCases {
		err := suite.vm.RevokeValidator(suite.Ctx, tc.username)
		suite.Equal(tc.expectRes, err, "%s", tc.testName)

		if tc.expectRes == nil {
			lst := suite.vm.storage.GetValidatorList(suite.Ctx)
			suite.Equal(tc.expectList, *lst, "%s", tc.testName)
			val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.username)
			suite.NoError(err)
			suite.Equal(tc.expectVal, *val, "%s", tc.testName)
		}

	}
}

func (suite *ValidatorManagerTestSuite) TestGetCommittingValidatorsVotes() {
	validators := map[linotypes.AccountKey]linotypes.Coin{
		linotypes.AccountKey("test1"): linotypes.NewCoinFromInt64(100),
		linotypes.AccountKey("test2"): linotypes.NewCoinFromInt64(200),
		linotypes.AccountKey("test3"): linotypes.NewCoinFromInt64(300),
		linotypes.AccountKey("test4"): linotypes.NewCoinFromInt64(400),
		linotypes.AccountKey("test5"): linotypes.NewCoinFromInt64(500),
		linotypes.AccountKey("test6"): linotypes.NewCoinFromInt64(600),
		linotypes.AccountKey("test7"): linotypes.NewCoinFromInt64(700),
	}
	suite.SetupValidatorAndVotes(validators)

	testCases := []struct {
		testName  string
		prevList  model.ValidatorList
		expectRes []model.ReceivedVotesStatus
	}{
		{
			testName: "get committing validators votes status",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("test1"),
					linotypes.AccountKey("test2"),
					linotypes.AccountKey("test3"),
				},
				Standby: []linotypes.AccountKey{
					linotypes.AccountKey("test4"),
					linotypes.AccountKey("test5"),
					linotypes.AccountKey("test6"),
				},
				Candidates: []linotypes.AccountKey{
					linotypes.AccountKey("test7"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectRes: []model.ReceivedVotesStatus{
				{
					ValidatorName: linotypes.AccountKey("test1"),
					ReceivedVotes: linotypes.NewCoinFromInt64(100),
				},
				{
					ValidatorName: linotypes.AccountKey("test2"),
					ReceivedVotes: linotypes.NewCoinFromInt64(200),
				},
				{
					ValidatorName: linotypes.AccountKey("test3"),
					ReceivedVotes: linotypes.NewCoinFromInt64(300),
				},
				{
					ValidatorName: linotypes.AccountKey("test4"),
					ReceivedVotes: linotypes.NewCoinFromInt64(400),
				},
				{
					ValidatorName: linotypes.AccountKey("test5"),
					ReceivedVotes: linotypes.NewCoinFromInt64(500),
				},
				{
					ValidatorName: linotypes.AccountKey("test6"),
					ReceivedVotes: linotypes.NewCoinFromInt64(600),
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		lst := suite.vm.GetCommittingValidatorVoteStatus(suite.Ctx)
		suite.Equal(tc.expectRes, lst, "%s", tc.testName)
	}
}

func (suite *ValidatorManagerTestSuite) TestUpdateValidator() {
	valKey := secp256k1.GenPrivKey().PubKey()
	val := linotypes.AccountKey("val")

	testCases := []struct {
		testName  string
		username  linotypes.AccountKey
		link      string
		expectVal model.Validator
	}{
		{
			testName: "update validator",
			link:     "web1111111",
			username: val,
			expectVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: valKey.Address(),
					Power:   1,
				},
				Link:          "web1111111",
				PubKey:        valKey,
				Username:      val,
				ReceivedVotes: linotypes.NewCoinFromInt64(300),
			},
		},
	}
	for _, tc := range testCases {
		err := suite.vm.RegisterValidator(suite.Ctx, tc.username, valKey, tc.link)
		suite.NoError(err)

		err = suite.vm.UpdateValidator(suite.Ctx, tc.username, tc.link)
		suite.NoError(err)

		val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.username)
		suite.NoError(err)
		suite.Equal(tc.expectVal, *val, "%s", tc.testName)
	}
}
func (suite *ValidatorManagerTestSuite) TestPunishCommittingValidator() {
	suite.vote.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("abs2")).Return(
		linotypes.NewCoinFromInt64(200000*linotypes.Decimals), nil).Maybe()
	suite.vote.On("SlashStake", suite.Ctx, linotypes.AccountKey("abs2"),
		linotypes.NewCoinFromInt64(200*linotypes.Decimals),
		linotypes.InflationValidatorPool).Return(
		linotypes.NewCoinFromInt64(200*linotypes.Decimals), nil).Maybe()
	absKey := secp256k1.GenPrivKey().PubKey()
	abs := linotypes.AccountKey("abs2")
	absVal := model.Validator{
		ABCIValidator: abci.Validator{
			Address: absKey.Address(),
			Power:   linotypes.TendermintValidatorPower,
		},
		PubKey:        absKey,
		Username:      abs,
		ReceivedVotes: linotypes.NewCoinFromInt64(2000),
		AbsentCommit:  1,
		NumSlash:      5,
	}
	suite.vm.storage.SetValidator(suite.Ctx, abs, &absVal)

	testCases := []struct {
		testName     string
		prevList     model.ValidatorList
		expectedList model.ValidatorList
		expectedVal  model.Validator
		username     linotypes.AccountKey
	}{
		{
			testName: "punish validator",
			prevList: model.ValidatorList{
				Oncall: []linotypes.AccountKey{
					linotypes.AccountKey("abs2"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(200),
				LowestOncall:       linotypes.AccountKey("abs2"),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectedList: model.ValidatorList{
				Jail: []linotypes.AccountKey{
					linotypes.AccountKey("abs2"),
				},
				LowestOncallVotes:  linotypes.NewCoinFromInt64(0),
				LowestOncall:       linotypes.AccountKey(""),
				LowestStandbyVotes: linotypes.NewCoinFromInt64(0),
				LowestStandby:      linotypes.AccountKey(""),
			},
			expectedVal: model.Validator{
				ABCIValidator: abci.Validator{
					Address: absKey.Address(),
					Power:   0,
				},
				PubKey:        absKey,
				Username:      abs,
				ReceivedVotes: linotypes.NewCoinFromInt64(2000),
				AbsentCommit:  0,
				NumSlash:      0,
			},
			username: linotypes.AccountKey("abs2"),
		},
	}

	for _, tc := range testCases {
		suite.vm.storage.SetValidatorList(suite.Ctx, &tc.prevList)
		err := suite.vm.PunishCommittingValidator(suite.Ctx, tc.username, linotypes.NewCoinFromInt64(200*linotypes.Decimals), linotypes.PunishNoPriceFed)
		suite.NoError(err)
		actualList := suite.vm.storage.GetValidatorList(suite.Ctx)
		suite.Equal(tc.expectedList, *actualList, "%s", tc.testName)
		val, err := suite.vm.storage.GetValidator(suite.Ctx, tc.username)
		suite.NoError(err)
		suite.Equal(tc.expectedVal, *val, "%s", tc.testName)
	}
}
