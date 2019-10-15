package manager

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account/mocks"
	global "github.com/lino-network/lino/x/global/mocks"
	hk "github.com/lino-network/lino/x/vote/manager/mocks"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/lino-network/lino/x/vote/types"
)

var (
	storeKeyStr = "testVoterStore"
	kvStoreKey  = sdk.NewKVStoreKey(storeKeyStr)
)

type VoteStoreDumper struct{}

func (dumper VoteStoreDumper) NewDumper() *testutils.Dumper {
	return model.NewVoteDumper(model.NewVoteStorage(kvStoreKey))
}

type VoteManagerTestSuite struct {
	testsuites.GoldenTestSuite
	vm     VoteManager
	ph     *param.ParamKeeper
	am     *acc.AccountKeeper
	global *global.GlobalKeeper
	hooks  *hk.StakingHooks

	// mock data
	user1 linotypes.AccountKey
	user2 linotypes.AccountKey
	user3 linotypes.AccountKey

	voter2 model.Voter
	voter3 model.Voter

	minStakeInAmount linotypes.Coin
}

func TestVoteManagerTestSuite(t *testing.T) {
	suite.Run(t, &VoteManagerTestSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(VoteStoreDumper{}, kvStoreKey),
	})
}

func (suite *VoteManagerTestSuite) SetupTest() {
	suite.SetupCtx(0, time.Unix(0, 0), kvStoreKey)
	suite.am = &acc.AccountKeeper{}
	suite.ph = &param.ParamKeeper{}
	suite.global = &global.GlobalKeeper{}
	suite.hooks = &hk.StakingHooks{}
	suite.vm = NewVoteManager(kvStoreKey, suite.ph, suite.am, suite.global)
	suite.vm = *suite.vm.SetHooks(suite.hooks)
	suite.minStakeInAmount = linotypes.NewCoinFromInt64(1000 * linotypes.Decimals)

	// suite.voter2 = model.Voter{
	// 	Username:  suite.user2,
	// 	LinoStake: suite.stakeInAmount,
	// 	Duty:      types.DutyVoter,
	// }
	// suite.voter3 = model.Voter{
	// 	Username:     suite.user3,
	// 	LinoStake:    suite.stakeInAmount,
	// 	FrozenAmount: suite.stakeInAmount,
	// 	Duty:         types.DutyValidator,
	// }

	suite.ph.On("GetVoteParam", mock.Anything).Return(&parammodel.VoteParam{
		MinStakeIn:                 suite.minStakeInAmount,
		VoterCoinReturnIntervalSec: 100,
		VoterCoinReturnTimes:       1,
	}).Maybe()
	// set initial stake stats for day 0.
	suite.vm.InitGenesis(suite.Ctx)
}

func (suite *VoteManagerTestSuite) ResetGlobal() {
	suite.global = &global.GlobalKeeper{}
	suite.vm = NewVoteManager(kvStoreKey, suite.ph, suite.am, suite.global)
	suite.vm = *suite.vm.SetHooks(suite.hooks)
}

func (suite *VoteManagerTestSuite) ResetParam() {
	suite.ph = &param.ParamKeeper{}
	suite.vm = NewVoteManager(kvStoreKey, suite.ph, suite.am, suite.global)
	suite.vm = *suite.vm.SetHooks(suite.hooks)
}

func (suite *VoteManagerTestSuite) TestStakeIn() {
	user1 := linotypes.AccountKey("user1")

	testCases := []struct {
		testName     string
		username     linotypes.AccountKey
		amount       linotypes.Coin
		lessThanMin  bool
		moveErr      sdk.Error
		atWhen       time.Time
		expectErr    sdk.Error
		expectVoter  *model.Voter
		expetecStats *model.LinoStakeStat
	}{
		{
			testName:    "stake in amount less than minimum requirement",
			username:    user1,
			amount:      suite.minStakeInAmount.Minus(linotypes.NewCoinFromInt64(1)),
			lessThanMin: true,
			expectErr:   types.ErrInsufficientDeposit(),
		},
		{
			testName:  "stake in with insufficient balance",
			username:  user1,
			amount:    suite.minStakeInAmount,
			moveErr:   linotypes.ErrTestDummyError(), // just a mock, any error is fine
			expectErr: linotypes.ErrTestDummyError(),
		},
		{
			testName:  "stake in minimum requirement",
			username:  user1,
			amount:    suite.minStakeInAmount,
			atWhen:    time.Unix(100, 0),
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:          user1,
				LinoStake:         suite.minStakeInAmount,
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 100,
			},
			expetecStats: &model.LinoStakeStat{
				TotalConsumptionFriction: linotypes.NewCoinFromInt64(0),
				UnclaimedFriction:        linotypes.NewCoinFromInt64(0),
				TotalLinoStake:           suite.minStakeInAmount,
				UnclaimedLinoStake:       suite.minStakeInAmount,
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.testName, func() {
			suite.hooks.On("AfterAddingStake", mock.Anything, tc.username).Return(nil).Maybe()
			suite.NextBlock(tc.atWhen)
			if !tc.lessThanMin {
				suite.am.On("MoveToPool", mock.Anything, linotypes.VoteStakeInPool,
					linotypes.NewAccOrAddrFromAcc(tc.username), tc.amount).Return(tc.moveErr).Once()
			}
			suite.global.On("GetPastDay", mock.Anything, int64(100)).Return(int64(0)).Maybe()
			err := suite.vm.StakeIn(suite.Ctx, tc.username, tc.amount)
			suite.Equal(tc.expectErr, err)
			if tc.expectErr == nil {
				voter, err := suite.vm.GetVoter(suite.Ctx, tc.username)
				suite.Nil(err)
				suite.Equal(tc.expectVoter, voter)
				stats, err := suite.vm.storage.GetLinoStakeStat(suite.Ctx, 0)
				suite.Nil(err)
				suite.Equal(tc.expetecStats, stats)
			}

			suite.global.AssertExpectations(suite.T())
			suite.am.AssertExpectations(suite.T())

			suite.Golden()
		})
	}
}

func (suite *VoteManagerTestSuite) TestStakeInFor() {
	user1 := linotypes.AccountKey("user1")
	user2 := linotypes.AccountKey("user2")

	testCases := []struct {
		testName     string
		username     linotypes.AccountKey
		stakeInFor   linotypes.AccountKey
		amount       linotypes.Coin
		lessThanMin  bool
		moveErr      sdk.Error
		atWhen       time.Time
		expectErr    sdk.Error
		expectVoter  *model.Voter
		expetecStats *model.LinoStakeStat
	}{
		{
			testName:    "stake in amount less than minimum requirement",
			username:    user1,
			stakeInFor:  user2,
			amount:      suite.minStakeInAmount.Minus(linotypes.NewCoinFromInt64(1)),
			lessThanMin: true,
			expectErr:   types.ErrInsufficientDeposit(),
		},
		{
			testName:   "stake in with insufficient balance",
			username:   user1,
			stakeInFor: user2,
			amount:     suite.minStakeInAmount,
			moveErr:    linotypes.ErrTestDummyError(), // just a mock, any error is fine
			expectErr:  linotypes.ErrTestDummyError(),
		},
		{
			testName:   "stake in minimum requirement",
			username:   user1,
			stakeInFor: user2,
			amount:     suite.minStakeInAmount,
			atWhen:     time.Unix(100, 0),
			expectErr:  nil,
			expectVoter: &model.Voter{
				Username:          user2,
				LinoStake:         suite.minStakeInAmount,
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 100,
			},
			expetecStats: &model.LinoStakeStat{
				TotalConsumptionFriction: linotypes.NewCoinFromInt64(0),
				UnclaimedFriction:        linotypes.NewCoinFromInt64(0),
				TotalLinoStake:           suite.minStakeInAmount,
				UnclaimedLinoStake:       suite.minStakeInAmount,
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.testName, func() {
			suite.hooks.On("AfterAddingStake", mock.Anything, tc.stakeInFor).Return(nil).Maybe()
			suite.NextBlock(tc.atWhen)
			if !tc.lessThanMin {
				suite.am.On("MoveToPool", mock.Anything, linotypes.VoteStakeInPool,
					linotypes.NewAccOrAddrFromAcc(tc.username), tc.amount).Return(tc.moveErr).Once()
			}
			suite.global.On("GetPastDay", mock.Anything, int64(100)).Return(int64(0)).Maybe()
			err := suite.vm.StakeInFor(suite.Ctx, tc.username, tc.stakeInFor, tc.amount)
			suite.Equal(tc.expectErr, err)
			if tc.expectErr == nil {
				_, err := suite.vm.GetVoter(suite.Ctx, tc.username)
				suite.NotNil(err)
				voter, err := suite.vm.GetVoter(suite.Ctx, tc.stakeInFor)
				suite.Nil(err)
				suite.Equal(tc.expectVoter, voter)
				stats, err := suite.vm.storage.GetLinoStakeStat(suite.Ctx, 0)
				suite.Nil(err)
				suite.Equal(tc.expetecStats, stats)
			}

			suite.global.AssertExpectations(suite.T())
			suite.am.AssertExpectations(suite.T())

			suite.Golden()
		})
	}

}

// TestMultipleStakeInWithConsumption
// script: (0.claim. 1.stake-in  2.consumption 3.end of day)
//       user1    user2    consumption   user1-calim  user2-claim
// day0  1000     0        222           x            x
// day1  2000     2000     444           222          x
// day2  -2000    -1000    888           x            x
// day3  0        0        100           x            x
// day4  2000     2000     500           760          672
// day5  0        1000     300           250          250         // all claimed
// day6  2000     2000     0             x            171         // 300 * (4/7)
// day7  0        0        300           129          0           // 1 - 300 * (4/7)
func (suite *VoteManagerTestSuite) TestMultipleStakeInWithConsumption() {
	// setup global for PastDay calculation
	suite.ResetGlobal()
	for i := 0; i <= 100; i++ {
		suite.global.On("GetPastDay", mock.Anything, int64(i)).Return(int64(i))
	}
	suite.global.On("RegisterEventAtTime",
		mock.Anything, mock.Anything, mock.Anything).Return(nil)

	suite.ResetParam()
	suite.ph.On("GetVoteParam", mock.Anything).Return(&parammodel.VoteParam{
		MinStakeIn:                 linotypes.NewCoinFromInt64(1),
		VoterCoinReturnIntervalSec: 100,
		VoterCoinReturnTimes:       1,
	}).Maybe()

	user1 := linotypes.AccountKey("user1")
	user2 := linotypes.AccountKey("user2")
	suite.hooks.On("AfterAddingStake", mock.Anything, mock.Anything).Return(nil).Maybe()
	suite.hooks.On("AfterSubtractingStake", mock.Anything, mock.Anything).Return(nil).Maybe()
	suite.am.On("AddFrozenMoney", mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, int64(100), int64(1)).Return(nil)
	// linotypes.VoteReturnCoin, linotypes.VoteStakeReturnPool
	suite.am.On("MoveToPool", mock.Anything, linotypes.VoteStakeInPool,
		mock.Anything, mock.Anything).Return(nil)
	suite.am.On("MoveBetweenPools", mock.Anything,
		linotypes.VoteStakeInPool, linotypes.VoteStakeReturnPool, mock.Anything).Return(nil)

	newCoin := func(n int64) *linotypes.Coin {
		coin := linotypes.NewCoinFromInt64(n)
		return &coin
	}

	for i, tc := range []struct {
		user1claim   *linotypes.Coin
		user2claim   *linotypes.Coin
		user1stakein *linotypes.Coin
		user2stakein *linotypes.Coin
		consumption  linotypes.Coin
	}{
		{
			// day0
			user1stakein: newCoin(1000),
			consumption:  *newCoin(222),
		},
		{
			// day1
			user1claim:   newCoin(222),
			user1stakein: newCoin(2000),
			user2stakein: newCoin(2000),
			consumption:  *newCoin(444),
		},
		{
			// day2
			user1stakein: newCoin(-2000),
			user2stakein: newCoin(-1000),
			consumption:  *newCoin(888),
		},
		{
			// day3
			consumption: *newCoin(100),
		},
		{

			// day4
			user1claim:   newCoin(760),
			user2claim:   newCoin(672),
			user1stakein: newCoin(2000),
			user2stakein: newCoin(2000),
			consumption:  *newCoin(500),
		},
		{
			// day5
			user1claim:   newCoin(250),
			user2claim:   newCoin(250),
			user2stakein: newCoin(1000),
			consumption:  *newCoin(300),
		},
		{
			// day6
			user2claim:   newCoin(171),
			user1stakein: newCoin(2000),
			user2stakein: newCoin(2000),
			consumption:  *newCoin(0),
		},
		{
			// day7
			user1claim:  newCoin(129),
			consumption: *newCoin(300),
		},
	} {
		suite.NextBlock(time.Unix(int64(i), 0))
		if i != 0 {
			err := suite.vm.DailyAdvanceLinoStakeStats(suite.Ctx)
			suite.Nil(err)
		}
		if tc.user1claim != nil {
			suite.am.On("MoveFromPool", mock.Anything, linotypes.VoteFrictionPool,
				linotypes.NewAccOrAddrFromAcc(user1), *tc.user1claim).Return(nil).Once()
			err := suite.vm.ClaimInterest(suite.Ctx, user1)
			suite.Nil(err)
		}
		if tc.user2claim != nil {
			suite.am.On("MoveFromPool", mock.Anything, linotypes.VoteFrictionPool,
				linotypes.NewAccOrAddrFromAcc(user2), *tc.user2claim).Return(nil).Once()
			err := suite.vm.ClaimInterest(suite.Ctx, user2)
			suite.Nil(err)
		}
		if tc.user1stakein != nil {
			if tc.user1stakein.IsPositive() {
				err := suite.vm.StakeIn(suite.Ctx, user1, *tc.user1stakein)
				suite.Nil(err)
			} else {
				err := suite.vm.StakeOut(suite.Ctx, user1, tc.user1stakein.Neg())
				suite.Nil(err)
			}
		}
		if tc.user2stakein != nil {
			if tc.user2stakein.IsPositive() {
				err := suite.vm.StakeIn(suite.Ctx, user2, *tc.user2stakein)
				suite.Nil(err)
			} else {
				err := suite.vm.StakeOut(suite.Ctx, user2, tc.user2stakein.Neg())
				suite.Nil(err)
			}
		}
		err := suite.vm.RecordFriction(suite.Ctx, tc.consumption)
		suite.Nil(err)
	}
	suite.Golden()
}

func (suite *VoteManagerTestSuite) TestStakeOut() {
	suite.hooks.On("AfterSubtractingStake", mock.Anything, suite.user2).Return(nil).Maybe()
	suite.am.On(
		"AddFrozenMoney", mock.Anything, suite.user2, suite.minStakeInAmount,
		int64(300), int64(100), int64(1)).Return(nil).Maybe()

	// add stake to user2
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter2)

	// add stake and frozon amount to user3
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter3)

	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		amount      linotypes.Coin
		atWhen      time.Time
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:    "stake out from user without stake",
			username:    suite.user1,
			amount:      suite.minStakeInAmount,
			atWhen:      time.Unix(0, 0),
			expectErr:   types.ErrVoterNotFound(),
			expectVoter: nil,
		},
		{
			testName:  "stake out amount more than user has",
			username:  suite.user2,
			amount:    suite.minStakeInAmount.Plus(linotypes.NewCoinFromInt64(1)),
			atWhen:    time.Unix(100, 0),
			expectErr: types.ErrInsufficientStake(),
			expectVoter: &model.Voter{
				Username:          suite.user2,
				LinoStake:         suite.minStakeInAmount,
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 0,
			},
		},
		{
			testName:  "stake out from user with all stake frozen",
			username:  suite.user3,
			amount:    suite.minStakeInAmount,
			atWhen:    time.Unix(200, 0),
			expectErr: types.ErrInsufficientStake(),
			expectVoter: &model.Voter{
				Username:     suite.user3,
				LinoStake:    suite.minStakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: suite.minStakeInAmount,
			},
		},
		{
			testName:  "stake out from user with sufficient stake",
			username:  suite.user2,
			amount:    suite.minStakeInAmount,
			atWhen:    time.Unix(300, 0),
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:          suite.user2,
				LinoStake:         linotypes.NewCoinFromInt64(0),
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 300,
			},
		},
	}

	for _, tc := range testCases {
		ctx := suite.Ctx.WithBlockHeader(abci.Header{Time: tc.atWhen})
		err := suite.vm.StakeOut(ctx, tc.username, tc.amount)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		voter, _ := suite.vm.GetVoter(ctx, tc.username)
		suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
	}
}

// func (suite *VoteManagerTestSuite) TestClaimInterest() {
// 	suite.global.On(
// 		"GetInterestSince", mock.Anything, int64(500),
// 		suite.minStakeInAmount).Return(suite.interest, nil).Twice()
// 	suite.am.On(
// 		"AddCoinToUsername", mock.Anything, suite.user2,
// 		suite.interest).Return(nil).Once()
// 	suite.am.On(
// 		"AddCoinToUsername", mock.Anything, suite.user3,
// 		suite.interest.Plus(suite.interest)).Return(nil).Once()

// 	// add stake to user2
// 	suite.voter2.LastPowerChangeAt = 500
// 	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter2)

// 	// add stake and interest to user3
// 	suite.voter3.Interest = suite.interest
// 	suite.voter3.LastPowerChangeAt = 500
// 	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter3)

// 	testCases := []struct {
// 		testName    string
// 		username    linotypes.AccountKey
// 		atWhen      time.Time
// 		expectErr   sdk.Error
// 		expectVoter *model.Voter
// 	}{
// 		{
// 			testName:  "claim interest from user without interest in voter struct",
// 			username:  suite.user2,
// 			atWhen:    time.Unix(600, 0),
// 			expectErr: nil,
// 			expectVoter: &model.Voter{
// 				Username:          suite.user2,
// 				LinoStake:         suite.minStakeInAmount,
// 				Interest:          linotypes.NewCoinFromInt64(0),
// 				Duty:              types.DutyVoter,
// 				FrozenAmount:      linotypes.NewCoinFromInt64(0),
// 				LastPowerChangeAt: 600,
// 			},
// 		},
// 		{
// 			testName:  "claim interest from user with interest in voter struct",
// 			username:  suite.user3,
// 			atWhen:    time.Unix(600, 0),
// 			expectErr: nil,
// 			expectVoter: &model.Voter{
// 				Username:          suite.user3,
// 				LinoStake:         suite.minStakeInAmount,
// 				Interest:          linotypes.NewCoinFromInt64(0),
// 				Duty:              types.DutyValidator,
// 				FrozenAmount:      suite.minStakeInAmount,
// 				LastPowerChangeAt: 600,
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		ctx := suite.Ctx.WithBlockHeader(abci.Header{Time: tc.atWhen})
// 		err := suite.vm.ClaimInterest(ctx, tc.username)
// 		suite.Equal(tc.expectErr, err, "%s", tc.testName)
// 		voter, _ := suite.vm.GetVoter(ctx, tc.username)
// 		suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
// 	}
// }

func (suite *VoteManagerTestSuite) TestAssignDuty() {
	// add stake to user2
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter2)

	// add stake and interest to user3
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter3)

	testCases := []struct {
		testName     string
		username     linotypes.AccountKey
		duty         types.VoterDuty
		frozenAmount linotypes.Coin
		expectErr    sdk.Error
		expectVoter  *model.Voter
	}{
		{
			testName:     "assign duty to user without stake",
			username:     suite.user1,
			duty:         types.DutyValidator,
			frozenAmount: linotypes.NewCoinFromInt64(1),
			expectErr:    types.ErrVoterNotFound(),
			expectVoter:  nil,
		},
		{
			testName:     "assign duty to user with other duty",
			username:     suite.user3,
			duty:         types.DutyValidator,
			frozenAmount: linotypes.NewCoinFromInt64(1),
			expectErr:    types.ErrNotAVoterOrHasDuty(),
			expectVoter: &model.Voter{
				Username:     suite.user3,
				LinoStake:    suite.minStakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: suite.minStakeInAmount,
			},
		},
		{
			testName:     "frozen money larger than stake",
			username:     suite.user2,
			duty:         types.DutyValidator,
			frozenAmount: suite.minStakeInAmount.Plus(linotypes.NewCoinFromInt64(1)),
			expectErr:    types.ErrInsufficientStake(),
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    suite.minStakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyVoter,
				FrozenAmount: linotypes.NewCoinFromInt64(0),
			},
		},
		{
			testName:     "assign duty successfully",
			username:     suite.user2,
			duty:         types.DutyValidator,
			frozenAmount: suite.minStakeInAmount,
			expectErr:    nil,
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    suite.minStakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: suite.minStakeInAmount,
			},
		},
	}

	for _, tc := range testCases {
		err := suite.vm.AssignDuty(suite.Ctx, tc.username, tc.duty, tc.frozenAmount)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		voter, _ := suite.vm.GetVoter(suite.Ctx, tc.username)
		suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
	}
}

func (suite *VoteManagerTestSuite) TestUnassignDuty() {
	var waitingPeriodSec int64 = 100
	suite.global.On(
		"RegisterEventAtTime", mock.Anything, waitingPeriodSec,
		types.UnassignDutyEvent{Username: suite.user3}).Return(nil).Once()

	// add stake to user2
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter2)

	// add stake and interest to user3
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter3)

	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:    "unassign duty from user without stake",
			username:    suite.user1,
			expectErr:   types.ErrVoterNotFound(),
			expectVoter: nil,
		},
		{
			testName:  "unassign duty from user doesn't have duty",
			username:  suite.user2,
			expectErr: types.ErrNoDuty(),
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    suite.minStakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyVoter,
				FrozenAmount: linotypes.NewCoinFromInt64(0),
			},
		},
		{
			testName:  "unassign duty from user who has validator duty",
			username:  suite.user3,
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:     suite.user3,
				LinoStake:    suite.minStakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyPending,
				FrozenAmount: suite.minStakeInAmount,
			},
		},
		{
			testName:  "unassign duty again",
			username:  suite.user3,
			expectErr: types.ErrNoDuty(),
			expectVoter: &model.Voter{
				Username:     suite.user3,
				LinoStake:    suite.minStakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyPending,
				FrozenAmount: suite.minStakeInAmount,
			},
		},
	}

	for _, tc := range testCases {
		err := suite.vm.UnassignDuty(suite.Ctx, tc.username, waitingPeriodSec)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		voter, _ := suite.vm.GetVoter(suite.Ctx, tc.username)
		suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
	}
}

func (suite *VoteManagerTestSuite) TestSlashStake() {
	suite.hooks.On("AfterSlashing", mock.Anything, suite.user2).Return(nil).Maybe()
	suite.hooks.On("AfterSlashing", mock.Anything, suite.user3).Return(nil).Maybe()
	// add stake to user2
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter2)

	// add stake and interest to user3
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter3)

	testCases := []struct {
		testName            string
		username            linotypes.AccountKey
		amount              linotypes.Coin
		expectErr           sdk.Error
		expectSlashedAmount linotypes.Coin
		expectVoter         *model.Voter
	}{
		{
			testName:            "slash stake from user without stake",
			username:            suite.user1,
			amount:              linotypes.NewCoinFromInt64(1),
			expectErr:           types.ErrVoterNotFound(),
			expectSlashedAmount: linotypes.NewCoinFromInt64(0),
			expectVoter:         nil,
		},
		{
			testName:            "slash more than user's stake",
			username:            suite.user2,
			amount:              suite.minStakeInAmount.Plus(linotypes.NewCoinFromInt64(1)),
			expectErr:           nil,
			expectSlashedAmount: suite.minStakeInAmount,
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    linotypes.NewCoinFromInt64(0),
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyVoter,
				FrozenAmount: linotypes.NewCoinFromInt64(0),
			},
		},
		{
			testName:            "slash user's stake with frozen",
			username:            suite.user3,
			amount:              suite.minStakeInAmount,
			expectErr:           nil,
			expectSlashedAmount: suite.minStakeInAmount,
			expectVoter: &model.Voter{
				Username:     suite.user3,
				LinoStake:    linotypes.NewCoinFromInt64(0),
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: suite.minStakeInAmount,
			},
		},
	}

	for _, tc := range testCases {
		amount, err := suite.vm.SlashStake(suite.Ctx, tc.username, tc.amount)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		suite.Equal(tc.expectSlashedAmount, amount, "%s", tc.testName)
		voter, _ := suite.vm.GetVoter(suite.Ctx, tc.username)
		suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
	}
}

func (suite *VoteManagerTestSuite) TestExecUnassignDutyEvent() {
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter2)

	// add stake and interest to user3
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter3)

	testCases := []struct {
		testName    string
		event       types.UnassignDutyEvent
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:    "execute event on non exist voter",
			event:       types.UnassignDutyEvent{Username: suite.user1},
			expectErr:   types.ErrVoterNotFound(),
			expectVoter: nil,
		},
		{
			testName:  "execute event on voter without duty",
			event:     types.UnassignDutyEvent{Username: suite.user2},
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    suite.minStakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyVoter,
				FrozenAmount: linotypes.NewCoinFromInt64(0),
			},
		},
		{
			testName:  "execute event on voter with validator duty",
			event:     types.UnassignDutyEvent{Username: suite.user3},
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:     suite.user3,
				LinoStake:    suite.minStakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyVoter,
				FrozenAmount: linotypes.NewCoinFromInt64(0),
			},
		},
	}

	for _, tc := range testCases {
		err := suite.vm.ExecUnassignDutyEvent(suite.Ctx, tc.event)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		voter, _ := suite.vm.GetVoter(suite.Ctx, tc.event.Username)
		suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
	}
}

func (suite *VoteManagerTestSuite) TestGetLinoStakeAndDuty() {
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter2)

	// add stake and interest to user3
	suite.vm.storage.SetVoter(suite.Ctx, &suite.voter3)

	testCases := []struct {
		username linotypes.AccountKey
		stake    linotypes.Coin
		duty     types.VoterDuty
	}{
		{
			username: suite.user2,
			stake:    suite.minStakeInAmount,
			duty:     types.DutyVoter,
		},
		{
			username: suite.user3,
			stake:    suite.minStakeInAmount,
			duty:     types.DutyValidator,
		},
	}

	for _, tc := range testCases {
		stake, err := suite.vm.GetLinoStake(suite.Ctx, tc.username)
		suite.Nil(err)
		suite.Equal(stake, tc.stake)
		duty, err := suite.vm.GetVoterDuty(suite.Ctx, tc.username)
		suite.Nil(err)
		suite.Equal(duty, tc.duty)
	}
}
