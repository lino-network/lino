package manager

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
	accmn "github.com/lino-network/lino/x/account/manager"
	acc "github.com/lino-network/lino/x/account/mocks"
	global "github.com/lino-network/lino/x/global/mocks"
	hk "github.com/lino-network/lino/x/vote/manager/mocks"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/lino-network/lino/x/vote/types"
)

// background:
// 3 voters plus a userPendingDuty(with pending duty).
// user1 has 2000, staked in on day 0
// user2 has 1000, staked in on day 1, 100 interests unsettled.
// user3 has 1000 and 1000 frozen, validator
// friction day0 888,
// friction day1 999.
// all units in LINO.

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

	// for common/3voters.input
	user1           linotypes.AccountKey
	user2           linotypes.AccountKey
	user3           linotypes.AccountKey
	userNotVoter    linotypes.AccountKey
	userPendingDuty linotypes.AccountKey

	// exmaple data
	minStakeInAmount  linotypes.Coin
	returnIntervalSec int64
	returnTimes       int64
}

func TestVoteManagerTestSuite(t *testing.T) {
	suite.Run(t, &VoteManagerTestSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(VoteStoreDumper{}, kvStoreKey),
	})
}

func (suite *VoteManagerTestSuite) SetupTest() {
	suite.SetupCtx(0, time.Unix(0, 0), kvStoreKey)
	suite.user1 = linotypes.AccountKey("user1")
	suite.user2 = linotypes.AccountKey("user2")
	suite.user3 = linotypes.AccountKey("user3")
	suite.userNotVoter = linotypes.AccountKey("notavoter")
	suite.userPendingDuty = linotypes.AccountKey("pendingdutyuser")
	suite.am = &acc.AccountKeeper{}
	suite.ph = &param.ParamKeeper{}
	suite.global = &global.GlobalKeeper{}
	suite.hooks = &hk.StakingHooks{}
	suite.vm = NewVoteManager(kvStoreKey, suite.ph, suite.am, suite.global)
	suite.vm = *suite.vm.SetHooks(suite.hooks)

	suite.minStakeInAmount = linotypes.NewCoinFromInt64(1000 * linotypes.Decimals)
	suite.returnIntervalSec = 100
	suite.returnTimes = 1
	suite.ph.On("GetVoteParam", mock.Anything).Return(&parammodel.VoteParam{
		MinStakeIn:                 suite.minStakeInAmount,
		VoterCoinReturnIntervalSec: suite.returnIntervalSec,
		VoterCoinReturnTimes:       suite.returnTimes,
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
			suite.global.On("GetPastDay", mock.Anything, tc.atWhen.Unix()).Return(int64(0)).Maybe()
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
	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		amount      linotypes.Coin
		atWhen      time.Time
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:  "stake out from user without stake",
			username:  suite.userNotVoter,
			amount:    suite.minStakeInAmount,
			atWhen:    time.Unix(1, 0),
			expectErr: types.ErrVoterNotFound(),
		},
		{
			testName:  "stake out amount more than user has",
			username:  suite.user2,
			amount:    linotypes.NewCoinFromInt64(1000*linotypes.Decimals + 1),
			atWhen:    time.Unix(1, 0),
			expectErr: types.ErrInsufficientStake(),
		},
		{
			testName:  "stake out from user with stakes not enough due to fronzen",
			username:  suite.user3,
			amount:    suite.minStakeInAmount.Plus(linotypes.NewCoinFromInt64(1)),
			atWhen:    time.Unix(1, 0),
			expectErr: types.ErrInsufficientStake(),
		},
		{
			testName:  "stake out from user with sufficient stake",
			username:  suite.user1,
			amount:    suite.minStakeInAmount,
			atWhen:    time.Unix(1, 0),
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:          suite.user1,
				LinoStake:         linotypes.NewCoinFromInt64(1000 * linotypes.Decimals),
				Interest:          linotypes.NewCoinFromInt64(888 * linotypes.Decimals),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 1,
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.testName, func() {
			suite.SetupTest()
			suite.LoadState(false, "3voters")
			suite.hooks.On("AfterSubtractingStake",
				mock.Anything, mock.Anything).Return(nil).Maybe()
			for i := int64(0); i <= tc.atWhen.Unix(); i++ {
				suite.global.On("GetPastDay", mock.Anything, i).Return(i).Maybe()
			}
			suite.NextBlock(tc.atWhen)

			if tc.expectErr == nil {
				suite.am.On("MoveBetweenPools", mock.Anything,
					linotypes.VoteStakeInPool, linotypes.VoteStakeReturnPool,
					tc.amount).Return(nil).Once()

				suite.am.On("AddFrozenMoney", mock.Anything,
					tc.username, tc.amount, tc.atWhen.Unix(),
					suite.returnIntervalSec, suite.returnTimes).Return(nil).Once()
				suite.global.On(
					"RegisterEventAtTime", mock.Anything,
					tc.atWhen.Unix()+suite.returnIntervalSec,
					accmn.ReturnCoinEvent{
						Username:   tc.username,
						Amount:     tc.amount,
						ReturnType: linotypes.VoteReturnCoin,
						FromPool:   linotypes.VoteStakeReturnPool,
						At:         tc.atWhen.Unix() + suite.returnIntervalSec,
					}).Return(nil).Once()
			}
			err := suite.vm.StakeOut(suite.Ctx, tc.username, tc.amount)
			suite.Equal(tc.expectErr, err, "%s", tc.testName)
			if tc.expectVoter != nil {
				voter, err := suite.vm.GetVoter(suite.Ctx, tc.username)
				suite.Nil(err)
				suite.Equal(tc.expectVoter, voter)
			}
			suite.am.AssertExpectations(suite.T())
			suite.global.AssertExpectations(suite.T())
		})
	}
}

type claim struct {
	username     linotypes.AccountKey
	atWhen       int64
	expectErr    sdk.Error
	expectAmount *linotypes.Coin
	expectVoter  *model.Voter
}

func (suite *VoteManagerTestSuite) TestClaimInterest() {
	testCases := []struct {
		testName string
		maxDay   int64
		claims   []claim
	}{
		{
			testName: "voter not exists",
			maxDay:   1,
			claims: []claim{
				{
					username:  suite.userNotVoter,
					atWhen:    1,
					expectErr: types.ErrVoterNotFound(),
				},
			},
		},
		{
			testName: "claim interest for day0",
			maxDay:   1,
			claims: []claim{
				{
					username:     suite.user1,
					atWhen:       1,
					expectAmount: newCoin(888 * linotypes.Decimals),
					expectVoter: &model.Voter{
						Username:          suite.user1,
						LinoStake:         linotypes.NewCoinFromInt64(2000 * linotypes.Decimals),
						Interest:          linotypes.NewCoinFromInt64(0),
						Duty:              types.DutyVoter,
						FrozenAmount:      linotypes.NewCoinFromInt64(0),
						LastPowerChangeAt: 1,
					},
				},
				{ // claim again, no interest.
					username:     suite.user1,
					atWhen:       1,
					expectAmount: newCoin(0),
					expectVoter: &model.Voter{
						Username:          suite.user1,
						LinoStake:         linotypes.NewCoinFromInt64(2000 * linotypes.Decimals),
						Interest:          linotypes.NewCoinFromInt64(0),
						Duty:              types.DutyVoter,
						FrozenAmount:      linotypes.NewCoinFromInt64(0),
						LastPowerChangeAt: 1,
					},
				},
			},
		},
		{
			testName: "claim interest for all past days",
			maxDay:   2,
			claims: []claim{
				{
					username: suite.user1,
					atWhen:   2,
					expectAmount: newCoin(
						888*linotypes.Decimals + (999*linotypes.Decimals)/5*2),
					expectVoter: &model.Voter{
						Username:          suite.user1,
						LinoStake:         linotypes.NewCoinFromInt64(2000 * linotypes.Decimals),
						Interest:          linotypes.NewCoinFromInt64(0),
						Duty:              types.DutyVoter,
						FrozenAmount:      linotypes.NewCoinFromInt64(0),
						LastPowerChangeAt: 2,
					},
				},
			},
		},
		{
			testName: "claim interest from user with interest in voter struct",
			maxDay:   2,
			claims: []claim{
				{
					username:     suite.user2,
					atWhen:       2,
					expectAmount: newCoin((999*linotypes.Decimals)/5 + 100*linotypes.Decimals),
					expectVoter: &model.Voter{
						Username:          suite.user2,
						LinoStake:         linotypes.NewCoinFromInt64(1000 * linotypes.Decimals),
						Interest:          linotypes.NewCoinFromInt64(0),
						Duty:              types.DutyVoter,
						FrozenAmount:      linotypes.NewCoinFromInt64(0),
						LastPowerChangeAt: 2,
					},
				},
			},
		},
		{
			testName: "all claimed",
			maxDay:   2,
			claims: []claim{
				{
					username: suite.user1,
					atWhen:   2,
					expectAmount: newCoin(
						888*linotypes.Decimals + (999*linotypes.Decimals)/5*2),
					expectVoter: &model.Voter{
						Username:          suite.user1,
						LinoStake:         linotypes.NewCoinFromInt64(2000 * linotypes.Decimals),
						Interest:          linotypes.NewCoinFromInt64(0),
						Duty:              types.DutyVoter,
						FrozenAmount:      linotypes.NewCoinFromInt64(0),
						LastPowerChangeAt: 2,
					},
				},
				{
					username:     suite.user2,
					atWhen:       2,
					expectAmount: newCoin((999*linotypes.Decimals)/5 + 100*linotypes.Decimals),
					expectVoter: &model.Voter{
						Username:          suite.user2,
						LinoStake:         linotypes.NewCoinFromInt64(1000 * linotypes.Decimals),
						Interest:          linotypes.NewCoinFromInt64(0),
						Duty:              types.DutyVoter,
						FrozenAmount:      linotypes.NewCoinFromInt64(0),
						LastPowerChangeAt: 2,
					},
				},
				{
					username:     suite.user3,
					atWhen:       2,
					expectAmount: newCoin((999 * linotypes.Decimals) / 5 * 2),
					expectVoter: &model.Voter{
						Username:          suite.user3,
						LinoStake:         linotypes.NewCoinFromInt64(2000 * linotypes.Decimals),
						Interest:          linotypes.NewCoinFromInt64(0),
						Duty:              types.DutyValidator,
						FrozenAmount:      linotypes.NewCoinFromInt64(1000 * linotypes.Decimals),
						LastPowerChangeAt: 2,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.testName, func() {
			suite.SetupTest()
			suite.LoadState(false, "3voters")
			for i := int64(0); i <= tc.maxDay; i++ {
				suite.global.On("GetPastDay", mock.Anything, i).Return(i).Maybe()
			}
			for _, claim := range tc.claims {
				suite.NextBlock(time.Unix(claim.atWhen, 0))
				if claim.expectErr == nil {
					suite.am.On("MoveFromPool", mock.Anything, linotypes.VoteFrictionPool,
						linotypes.NewAccOrAddrFromAcc(claim.username),
						*claim.expectAmount).Return(nil).Once()
				}
				err := suite.vm.ClaimInterest(suite.Ctx, claim.username)
				suite.Equal(claim.expectErr, err)
				if claim.expectVoter != nil {
					voter, err := suite.vm.GetVoter(suite.Ctx, claim.username)
					suite.Nil(err)
					suite.Equal(claim.expectVoter, voter)
				}
			}
			suite.am.AssertExpectations(suite.T())
			suite.global.AssertExpectations(suite.T())
			suite.Golden() //ensures that stake-stats are all correct.
		})
	}
}

func (suite *VoteManagerTestSuite) TestAssignDuty() {
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
			username:     suite.userNotVoter,
			duty:         types.DutyValidator,
			frozenAmount: linotypes.NewCoinFromInt64(1),
			expectErr:    types.ErrVoterNotFound(),
		},
		{
			testName:     "assign duty to user with other duty",
			username:     suite.user3,
			duty:         types.DutyValidator,
			frozenAmount: linotypes.NewCoinFromInt64(1),
			expectErr:    types.ErrNotAVoterOrHasDuty(),
		},
		{
			testName:     "negative frozen amount",
			username:     suite.user1,
			duty:         types.DutyValidator,
			frozenAmount: linotypes.NewCoinFromInt64(-1),
			expectErr:    types.ErrNegativeFrozenAmount(),
		},
		{
			testName:     "frozen money larger than stake",
			username:     suite.user1,
			duty:         types.DutyValidator,
			frozenAmount: *newCoin(2000*linotypes.Decimals + 1),
			expectErr:    types.ErrInsufficientStake(),
		},
		{
			testName:     "assign duty successfully",
			username:     suite.user1,
			duty:         types.DutyValidator,
			frozenAmount: *newCoin(1000 * linotypes.Decimals),
			expectErr:    nil,
			expectVoter: &model.Voter{
				Username:     suite.user1,
				LinoStake:    *newCoin(2000 * linotypes.Decimals),
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: *newCoin(1000 * linotypes.Decimals),
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.testName, func() {
			suite.SetupTest()
			suite.LoadState(false, "3voters")
			err := suite.vm.AssignDuty(suite.Ctx, tc.username, tc.duty, tc.frozenAmount)
			suite.Equal(tc.expectErr, err, "%s", tc.testName)
			if tc.expectVoter != nil {
				voter, _ := suite.vm.GetVoter(suite.Ctx, tc.username)
				suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
			}
			suite.Golden()
		})
	}
}

func (suite *VoteManagerTestSuite) TestUnassignDuty() {
	testCases := []struct {
		testName          string
		username          linotypes.AccountKey
		expectErr         sdk.Error
		testDoubleUnassin bool
		expectVoter       *model.Voter
	}{
		{
			testName:  "unassign duty from user without stake",
			username:  suite.userNotVoter,
			expectErr: types.ErrVoterNotFound(),
		},
		{
			testName:  "unassign duty from user doesnt have duty",
			username:  suite.user2,
			expectErr: types.ErrNoDuty(),
		},
		{
			testName:  "unassign duty from user has pending duty",
			username:  suite.userPendingDuty,
			expectErr: types.ErrNoDuty(),
		},
		{
			testName: "unassign duty from user who has validator duty",
			username: suite.user3,
			expectVoter: &model.Voter{
				Username:          suite.user3,
				LinoStake:         linotypes.NewCoinFromInt64(2000 * linotypes.Decimals),
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyPending,
				FrozenAmount:      linotypes.NewCoinFromInt64(1000 * linotypes.Decimals),
				LastPowerChangeAt: 1,
			},
		},
	}

	waitingPeriodSec := int64(100)
	for _, tc := range testCases {
		suite.Run(tc.testName, func() {
			suite.LoadState(false, "3voters")
			suite.NextBlock(time.Unix(1, 0))
			if tc.expectErr == nil {
				suite.global.On("RegisterEventAtTime", mock.Anything,
					1+waitingPeriodSec,
					types.UnassignDutyEvent{Username: tc.username}).Return(nil).Once()
			}
			err := suite.vm.UnassignDuty(suite.Ctx, tc.username, waitingPeriodSec)
			suite.Equal(tc.expectErr, err)
			if tc.expectVoter != nil {
				voter, err := suite.vm.GetVoter(suite.Ctx, tc.username)
				suite.Nil(err)
				suite.Equal(tc.expectVoter, voter)
			}
			suite.global.AssertExpectations(suite.T())
			suite.Golden()
		})
	}
}

func (suite *VoteManagerTestSuite) TestSlashStake() {
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
			username:            suite.userNotVoter,
			amount:              linotypes.NewCoinFromInt64(1),
			expectErr:           types.ErrVoterNotFound(),
			expectSlashedAmount: linotypes.NewCoinFromInt64(0),
		},
		{
			testName:            "slash more than user stake",
			username:            suite.user2,
			amount:              *newCoin(1000*linotypes.Decimals + 1),
			expectSlashedAmount: *newCoin(1000 * linotypes.Decimals),
			expectVoter: &model.Voter{
				Username:  suite.user2,
				LinoStake: linotypes.NewCoinFromInt64(0),
				Interest: linotypes.NewCoinFromInt64(
					(999*linotypes.Decimals)/5 + 100*linotypes.Decimals),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 2,
			},
		},
		{
			testName:            "slash users stake with frozen",
			username:            suite.user3,
			amount:              *newCoin(1500 * linotypes.Decimals),
			expectSlashedAmount: *newCoin(1500 * linotypes.Decimals),
			expectVoter: &model.Voter{
				Username:          suite.user3,
				LinoStake:         linotypes.NewCoinFromInt64(500 * linotypes.Decimals),
				Interest:          linotypes.NewCoinFromInt64(39960000),
				Duty:              types.DutyValidator,
				FrozenAmount:      suite.minStakeInAmount,
				LastPowerChangeAt: 2,
			},
		},
	}

	// all cases are assumed to happen at day2 to test poping interests upon slash.
	var destPool linotypes.PoolName = "dest"
	for _, tc := range testCases {
		suite.Run(tc.testName, func() {
			suite.SetupTest()
			suite.LoadState(false, "3voters")
			for i := int64(0); i <= 2; i++ {
				suite.global.On("GetPastDay", mock.Anything, i).Return(i).Maybe()
			}
			suite.NextBlock(time.Unix(2, 0))
			suite.vm.DailyAdvanceLinoStakeStats(suite.Ctx)
			if tc.expectErr == nil {
				suite.hooks.On("AfterSlashing", mock.Anything, tc.username).Return(nil).Once()
				suite.am.On("MoveBetweenPools", mock.Anything,
					linotypes.VoteStakeInPool, destPool, tc.expectSlashedAmount).Return(nil).Once()
			}

			amount, err := suite.vm.SlashStake(suite.Ctx, tc.username, tc.amount, destPool)
			suite.Equal(tc.expectErr, err)
			suite.Equal(tc.expectSlashedAmount, amount, "%s vs %s", tc.expectSlashedAmount, amount)
			if tc.expectVoter != nil {
				voter, _ := suite.vm.GetVoter(suite.Ctx, tc.username)
				suite.Equal(tc.expectVoter, voter)
			}

			suite.am.AssertExpectations(suite.T())
			suite.hooks.AssertExpectations(suite.T())
			suite.Golden() // ensure stake-stats are correct.
		})
	}
}

func (suite *VoteManagerTestSuite) TestExecUnassignDutyEvent() {
	testCases := []struct {
		testName    string
		event       types.UnassignDutyEvent
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:    "execute event on non exist voter",
			event:       types.UnassignDutyEvent{Username: suite.userNotVoter},
			expectErr:   types.ErrVoterNotFound(),
			expectVoter: nil,
		},
		{
			testName:  "execute event on voter with validator duty",
			event:     types.UnassignDutyEvent{Username: suite.user3},
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:          suite.user3,
				LinoStake:         linotypes.NewCoinFromInt64(2000 * linotypes.Decimals),
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 1,
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.testName, func() {
			suite.LoadState(false, "3voters")
			err := suite.vm.ExecUnassignDutyEvent(suite.Ctx, tc.event)
			suite.Equal(tc.expectErr, err, "%s", tc.testName)
			voter, _ := suite.vm.GetVoter(suite.Ctx, tc.event.Username)
			suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
		})
	}
}

func (suite *VoteManagerTestSuite) TestGetLinoStakeAndDuty() {
	testCases := []struct {
		username linotypes.AccountKey
		stake    linotypes.Coin
		duty     types.VoterDuty
	}{
		{
			username: suite.user2,
			stake:    linotypes.NewCoinFromInt64(1000 * linotypes.Decimals),
			duty:     types.DutyVoter,
		},
		{
			username: suite.user3,
			stake:    linotypes.NewCoinFromInt64(2000 * linotypes.Decimals),
			duty:     types.DutyValidator,
		},
	}

	suite.LoadState(false, "3voters")
	for _, tc := range testCases {
		stake, err := suite.vm.GetLinoStake(suite.Ctx, tc.username)
		suite.Nil(err)
		suite.Equal(stake, tc.stake)
		duty, err := suite.vm.GetVoterDuty(suite.Ctx, tc.username)
		suite.Nil(err)
		suite.Equal(duty, tc.duty)
	}
}

func newCoin(n int64) *linotypes.Coin {
	coin := linotypes.NewCoinFromInt64(n)
	return &coin
}
