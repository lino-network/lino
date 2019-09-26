package manager

import (
	"testing"
	"time"

	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	linotypes "github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account/mocks"
	global "github.com/lino-network/lino/x/global/mocks"
	hk "github.com/lino-network/lino/x/vote/manager/mocks"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/lino-network/lino/x/vote/types"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type VoteManagerTestSuite struct {
	testsuites.CtxTestSuite
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

	stakeInAmount linotypes.Coin
	interest      linotypes.Coin
}

func TestPostManagerTestSuite(t *testing.T) {
	suite.Run(t, new(VoteManagerTestSuite))
}

func (suite *VoteManagerTestSuite) SetupTest() {
	testVoteKey := sdk.NewKVStoreKey("vote")
	suite.SetupCtx(0, time.Unix(0, 0), testVoteKey)
	suite.am = &acc.AccountKeeper{}
	suite.ph = &param.ParamKeeper{}
	suite.global = &global.GlobalKeeper{}
	suite.hooks = &hk.StakingHooks{}
	suite.vm = NewVoteManager(testVoteKey, suite.ph, suite.am, suite.global)
	suite.vm = *suite.vm.SetHooks(suite.hooks)
	suite.user1 = linotypes.AccountKey("user1")
	suite.user2 = linotypes.AccountKey("user2")
	suite.user3 = linotypes.AccountKey("user3")
	suite.stakeInAmount = linotypes.NewCoinFromInt64(1000 * linotypes.Decimals)
	suite.interest = linotypes.NewCoinFromInt64(20 * linotypes.Decimals)

	suite.voter2 = model.Voter{
		Username:  suite.user2,
		LinoStake: suite.stakeInAmount,
		Duty:      types.DutyVoter,
	}
	suite.voter3 = model.Voter{
		Username:     suite.user3,
		LinoStake:    suite.stakeInAmount,
		FrozenAmount: suite.stakeInAmount,
		Duty:         types.DutyValidator,
	}

	suite.ph.On("GetVoteParam", mock.Anything).Return(&parammodel.VoteParam{
		MinStakeIn:                 suite.stakeInAmount,
		VoterCoinReturnIntervalSec: 100,
		VoterCoinReturnTimes:       1,
	}, nil).Maybe()
	suite.global.On(
		"GetInterestSince", mock.Anything, int64(0),
		linotypes.NewCoinFromInt64(0)).Return(linotypes.NewCoinFromInt64(0), nil).Maybe()
	suite.global.On(
		"GetInterestSince", mock.Anything, int64(0),
		suite.stakeInAmount).Return(linotypes.NewCoinFromInt64(0), nil).Maybe()
	suite.global.On("AddLinoStakeToStat", mock.Anything, suite.stakeInAmount).Return(nil).Maybe()
	suite.global.On("MinusLinoStakeFromStat", mock.Anything, suite.stakeInAmount).Return(nil).Maybe()
	suite.global.On(
		"RegisterCoinReturnEvent", mock.Anything, mock.Anything, int64(1), int64(100)).Return(nil).Maybe()

	suite.hooks.On("AfterAddingStake", mock.Anything, suite.user1).Return(nil).Maybe()
}

func (suite *VoteManagerTestSuite) TestStakeIn() {
	suite.am.On("MinusCoinFromUsername", mock.Anything, suite.user1, suite.stakeInAmount).Return(nil).Maybe()
	suite.global.On(
		"GetInterestSince", mock.Anything, int64(100),
		linotypes.NewCoinFromInt64(0)).Return(linotypes.NewCoinFromInt64(0), nil).Maybe()
	suite.global.On(
		"GetInterestSince", mock.Anything, int64(100), suite.stakeInAmount).Return(suite.interest, nil).Maybe()
	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		amount      linotypes.Coin
		atWhen      time.Time
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:    "stake in amount less than minimum requirement",
			username:    suite.user1,
			amount:      suite.stakeInAmount.Minus(linotypes.NewCoinFromInt64(1)),
			atWhen:      time.Unix(0, 0),
			expectErr:   types.ErrInsufficientDeposit(),
			expectVoter: nil,
		},
		{
			testName:  "stake in minimum requirement",
			username:  suite.user1,
			amount:    suite.stakeInAmount,
			atWhen:    time.Unix(100, 0),
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:          suite.user1,
				LinoStake:         suite.stakeInAmount,
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 100,
			},
		},
		{
			testName:  "stake in again minimum requirement with interest",
			username:  suite.user1,
			amount:    suite.stakeInAmount,
			atWhen:    time.Unix(200, 0),
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:          suite.user1,
				LinoStake:         suite.stakeInAmount.Plus(suite.stakeInAmount),
				Interest:          suite.interest,
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 200,
			},
		},
	}

	for _, tc := range testCases {
		ctx := suite.Ctx.WithBlockHeader(abci.Header{Time: tc.atWhen})
		err := suite.vm.StakeIn(ctx, tc.username, tc.amount)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		voter, _ := suite.vm.GetVoter(ctx, tc.username)
		suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
	}
}

func (suite *VoteManagerTestSuite) TestStakeOut() {
	suite.hooks.On("AfterSubtractingStake", mock.Anything, suite.user2).Return(nil).Maybe()
	suite.am.On(
		"AddFrozenMoney", mock.Anything, suite.user2, suite.stakeInAmount,
		int64(300), int64(100), int64(1)).Return(nil).Maybe()

	// add stake to user2
	e := suite.vm.storage.SetVoter(suite.Ctx, suite.user2, &suite.voter2)
	suite.Nil(e)

	// add stake and frozon amount to user3
	e = suite.vm.storage.SetVoter(suite.Ctx, suite.user3, &suite.voter3)
	suite.Nil(e)

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
			amount:      suite.stakeInAmount,
			atWhen:      time.Unix(0, 0),
			expectErr:   model.ErrVoterNotFound(),
			expectVoter: nil,
		},
		{
			testName:  "stake out amount more than user has",
			username:  suite.user2,
			amount:    suite.stakeInAmount.Plus(linotypes.NewCoinFromInt64(1)),
			atWhen:    time.Unix(100, 0),
			expectErr: types.ErrInsufficientStake(),
			expectVoter: &model.Voter{
				Username:          suite.user2,
				LinoStake:         suite.stakeInAmount,
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 0,
			},
		},
		{
			testName:  "stake out from user with all stake frozen",
			username:  suite.user3,
			amount:    suite.stakeInAmount,
			atWhen:    time.Unix(200, 0),
			expectErr: types.ErrInsufficientStake(),
			expectVoter: &model.Voter{
				Username:     suite.user3,
				LinoStake:    suite.stakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: suite.stakeInAmount,
			},
		},
		{
			testName:  "stake out from user with sufficient stake",
			username:  suite.user2,
			amount:    suite.stakeInAmount,
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

func (suite *VoteManagerTestSuite) TestClaimInterest() {
	suite.global.On(
		"GetInterestSince", mock.Anything, int64(500),
		suite.stakeInAmount).Return(suite.interest, nil).Twice()
	suite.am.On(
		"AddCoinToUsername", mock.Anything, suite.user2,
		suite.interest).Return(nil).Once()
	suite.am.On(
		"AddCoinToUsername", mock.Anything, suite.user3,
		suite.interest.Plus(suite.interest)).Return(nil).Once()

	// add stake to user2
	suite.voter2.LastPowerChangeAt = 500
	e := suite.vm.storage.SetVoter(suite.Ctx, suite.user2, &suite.voter2)
	suite.Nil(e)

	// add stake and interest to user3
	suite.voter3.Interest = suite.interest
	suite.voter3.LastPowerChangeAt = 500
	e = suite.vm.storage.SetVoter(suite.Ctx, suite.user3, &suite.voter3)
	suite.Nil(e)

	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		atWhen      time.Time
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:  "claim interest from user without interest in voter struct",
			username:  suite.user2,
			atWhen:    time.Unix(600, 0),
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:          suite.user2,
				LinoStake:         suite.stakeInAmount,
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyVoter,
				FrozenAmount:      linotypes.NewCoinFromInt64(0),
				LastPowerChangeAt: 600,
			},
		},
		{
			testName:  "claim interest from user with interest in voter struct",
			username:  suite.user3,
			atWhen:    time.Unix(600, 0),
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:          suite.user3,
				LinoStake:         suite.stakeInAmount,
				Interest:          linotypes.NewCoinFromInt64(0),
				Duty:              types.DutyValidator,
				FrozenAmount:      suite.stakeInAmount,
				LastPowerChangeAt: 600,
			},
		},
	}

	for _, tc := range testCases {
		ctx := suite.Ctx.WithBlockHeader(abci.Header{Time: tc.atWhen})
		err := suite.vm.ClaimInterest(ctx, tc.username)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		voter, _ := suite.vm.GetVoter(ctx, tc.username)
		suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
	}
}

func (suite *VoteManagerTestSuite) TestAssignDuty() {
	// add stake to user2
	e := suite.vm.storage.SetVoter(suite.Ctx, suite.user2, &suite.voter2)
	suite.Nil(e)

	// add stake and interest to user3
	e = suite.vm.storage.SetVoter(suite.Ctx, suite.user3, &suite.voter3)
	suite.Nil(e)

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
			expectErr:    model.ErrVoterNotFound(),
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
				LinoStake:    suite.stakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: suite.stakeInAmount,
			},
		},
		{
			testName:     "frozen money larger than stake",
			username:     suite.user2,
			duty:         types.DutyValidator,
			frozenAmount: suite.stakeInAmount.Plus(linotypes.NewCoinFromInt64(1)),
			expectErr:    types.ErrInsufficientStake(),
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    suite.stakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyVoter,
				FrozenAmount: linotypes.NewCoinFromInt64(0),
			},
		},
		{
			testName:     "assign duty successfully",
			username:     suite.user2,
			duty:         types.DutyValidator,
			frozenAmount: suite.stakeInAmount,
			expectErr:    nil,
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    suite.stakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: suite.stakeInAmount,
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
	e := suite.vm.storage.SetVoter(suite.Ctx, suite.user2, &suite.voter2)
	suite.Nil(e)

	// add stake and interest to user3
	e = suite.vm.storage.SetVoter(suite.Ctx, suite.user3, &suite.voter3)
	suite.Nil(e)

	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:    "unassign duty from user without stake",
			username:    suite.user1,
			expectErr:   model.ErrVoterNotFound(),
			expectVoter: nil,
		},
		{
			testName:  "unassign duty from user doesn't have duty",
			username:  suite.user2,
			expectErr: types.ErrNoDuty(),
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    suite.stakeInAmount,
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
				LinoStake:    suite.stakeInAmount,
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: suite.stakeInAmount,
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
	suite.hooks.On("AfterSubtractingStake", mock.Anything, suite.user2).Return(nil).Maybe()
	suite.hooks.On("AfterSubtractingStake", mock.Anything, suite.user3).Return(nil).Maybe()
	// add stake to user2
	e := suite.vm.storage.SetVoter(suite.Ctx, suite.user2, &suite.voter2)
	suite.Nil(e)

	// add stake and interest to user3
	e = suite.vm.storage.SetVoter(suite.Ctx, suite.user3, &suite.voter3)
	suite.Nil(e)

	testCases := []struct {
		testName    string
		username    linotypes.AccountKey
		amount      linotypes.Coin
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:    "slash stake from user without stake",
			username:    suite.user1,
			amount:      linotypes.NewCoinFromInt64(1),
			expectErr:   model.ErrVoterNotFound(),
			expectVoter: nil,
		},
		{
			testName:  "slash more than user's stake",
			username:  suite.user2,
			amount:    suite.stakeInAmount.Plus(linotypes.NewCoinFromInt64(1)),
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    linotypes.NewCoinFromInt64(0),
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyVoter,
				FrozenAmount: linotypes.NewCoinFromInt64(0),
			},
		},
		{
			testName:  "slash user's stake with frozen",
			username:  suite.user3,
			amount:    suite.stakeInAmount,
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:     suite.user3,
				LinoStake:    linotypes.NewCoinFromInt64(0),
				Interest:     linotypes.NewCoinFromInt64(0),
				Duty:         types.DutyValidator,
				FrozenAmount: suite.stakeInAmount,
			},
		},
	}

	for _, tc := range testCases {
		err := suite.vm.SlashStake(suite.Ctx, tc.username, tc.amount)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
		voter, _ := suite.vm.GetVoter(suite.Ctx, tc.username)
		suite.Equal(tc.expectVoter, voter, "%s", tc.testName)
	}
}

func (suite *VoteManagerTestSuite) TestExecUnassignDutyEvent() {
	e := suite.vm.storage.SetVoter(suite.Ctx, suite.user2, &suite.voter2)
	suite.Nil(e)

	// add stake and interest to user3
	e = suite.vm.storage.SetVoter(suite.Ctx, suite.user3, &suite.voter3)
	suite.Nil(e)

	testCases := []struct {
		testName    string
		event       types.UnassignDutyEvent
		expectErr   sdk.Error
		expectVoter *model.Voter
	}{
		{
			testName:    "execute event on non exist voter",
			event:       types.UnassignDutyEvent{Username: suite.user1},
			expectErr:   model.ErrVoterNotFound(),
			expectVoter: nil,
		},
		{
			testName:  "execute event on voter without duty",
			event:     types.UnassignDutyEvent{Username: suite.user2},
			expectErr: nil,
			expectVoter: &model.Voter{
				Username:     suite.user2,
				LinoStake:    suite.stakeInAmount,
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
				LinoStake:    suite.stakeInAmount,
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
	e := suite.vm.storage.SetVoter(suite.Ctx, suite.user2, &suite.voter2)
	suite.Nil(e)

	// add stake and interest to user3
	e = suite.vm.storage.SetVoter(suite.Ctx, suite.user3, &suite.voter3)
	suite.Nil(e)

	testCases := []struct {
		username linotypes.AccountKey
		stake    linotypes.Coin
		duty     types.VoterDuty
	}{
		{
			username: suite.user2,
			stake:    suite.stakeInAmount,
			duty:     types.DutyVoter,
		},
		{
			username: suite.user3,
			stake:    suite.stakeInAmount,
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
