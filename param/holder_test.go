package param

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("param")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
}

func TestGlobalAllocationParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := GlobalAllocationParam{
		GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
		ContentCreatorAllocation: types.NewDecFromRat(1, 100),
		InfraAllocation:          types.NewDecFromRat(1, 100),
		DeveloperAllocation:      types.NewDecFromRat(1, 100),
		ValidatorAllocation:      types.NewDecFromRat(97, 100),
	}
	err := ph.setGlobalAllocationParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Global allocation param should be equal")
}

func TestInfraInternalAllocationParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := InfraInternalAllocationParam{
		StorageAllocation: types.NewDecFromRat(50, 100),
		CDNAllocation:     types.NewDecFromRat(50, 100),
	}
	err := ph.setInfraInternalAllocationParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetInfraInternalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Infra internal allocation param should be equal")
}

func TestDeveloperParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := DeveloperParam{
		DeveloperMinDeposit:            types.NewCoinFromInt64(100000 * types.Decimals),
		DeveloperCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DeveloperCoinReturnTimes:       int64(7),
	}
	err := ph.setDeveloperParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetDeveloperParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Developer param should be equal")
}

func TestValidatorParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := ValidatorParam{
		ValidatorMinWithdraw:           types.NewCoinFromInt64(1 * types.Decimals),
		ValidatorMinVotingDeposit:      types.NewCoinFromInt64(300000 * types.Decimals),
		ValidatorMinCommittingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissVote:                types.NewCoinFromInt64(20000 * types.Decimals),
		PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:               types.NewCoinFromInt64(1000000 * types.Decimals),
		ValidatorListSize:              int64(21),
		AbsentCommitLimitation:         int64(100),
	}
	err := ph.setValidatorParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetValidatorParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Validator param should be equal")
}

func TestVoteParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := VoteParam{
		MinStakeIn:                     types.NewCoinFromInt64(1000 * types.Decimals),
		VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
		VoterCoinReturnTimes:           int64(7),
		DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DelegatorCoinReturnTimes:       int64(7),
	}
	err := ph.setVoteParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetVoteParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Voter param should be equal")
}

func TestProposalParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := ProposalParam{
		ContentCensorshipDecideSec:  int64(7 * 24 * 3600),
		ContentCensorshipPassRatio:  types.NewDecFromRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamExecutionSec: int64(24 * 3600),
		ChangeParamDecideSec:    int64(7 * 24 * 3600),
		ChangeParamPassRatio:    types.NewDecFromRat(70, 100),
		ChangeParamPassVotes:    types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit:   types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideSec:  int64(7 * 24 * 3600),
		ProtocolUpgradePassRatio:  types.NewDecFromRat(80, 100),
		ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
		ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
	}
	err := ph.setProposalParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetProposalParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Proposal param should be equal")
}

func TestCoinDayParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := CoinDayParam{
		SecondsToRecoverCoinDay: int64(7 * 24 * 3600),
	}
	err := ph.setCoinDayParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetCoinDayParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Coin day param should be equal")
}

func TestBandwidthParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
		VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
		GeneralMsgQuotaRatio:        types.NewDecFromRat(20, 100),
		GeneralMsgEMAFactor:         types.NewDecFromRat(1, 10),
		AppMsgQuotaRatio:            types.NewDecFromRat(80, 100),
		AppMsgEMAFactor:             types.NewDecFromRat(1, 10),
		ExpectedMaxMPS:              types.NewDecFromRat(1000, 1),
		MsgFeeFactorA:               types.NewDecFromRat(6, 1),
		MsgFeeFactorB:               types.NewDecFromRat(10, 1),
		MaxMPSDecayRate:             types.NewDecFromRat(99, 100),
		AppBandwidthPoolSize:        types.NewDecFromRat(10, 1),
		AppVacancyFactor:             types.NewDecFromRat(69, 100),
		AppPunishmentFactor:          types.NewDecFromRat(14, 5),
	}
	err := ph.setBandwidthParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetBandwidthParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Bandwidth param should be equal")
}

func TestAccountParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := AccountParam{
		MinimumBalance:               types.NewCoinFromInt64(1 * types.Decimals),
		RegisterFee:                  types.NewCoinFromInt64(1 * types.Decimals),
		FirstDepositFullCoinDayLimit: types.NewCoinFromInt64(1 * types.Decimals),
		MaxNumFrozenMoney:            10,
	}
	err := ph.setAccountParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetAccountParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Account param should be equal")
}

func TestInitParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()

	ph.InitParam(ctx)

	globalAllocationParam := GlobalAllocationParam{
		GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
		InfraAllocation:          types.NewDecFromRat(20, 100),
		ContentCreatorAllocation: types.NewDecFromRat(65, 100),
		DeveloperAllocation:      types.NewDecFromRat(10, 100),
		ValidatorAllocation:      types.NewDecFromRat(5, 100),
	}

	infraInternalAllocationParam := InfraInternalAllocationParam{
		StorageAllocation: types.NewDecFromRat(50, 100),
		CDNAllocation:     types.NewDecFromRat(50, 100),
	}

	developerParam := DeveloperParam{
		DeveloperMinDeposit:            types.NewCoinFromInt64(1000000 * types.Decimals),
		DeveloperCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DeveloperCoinReturnTimes:       int64(7),
	}

	validatorParam := ValidatorParam{
		ValidatorMinWithdraw:           types.NewCoinFromInt64(1 * types.Decimals),
		ValidatorMinVotingDeposit:      types.NewCoinFromInt64(300000 * types.Decimals),
		ValidatorMinCommittingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissVote:                types.NewCoinFromInt64(20000 * types.Decimals),
		PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:               types.NewCoinFromInt64(1000000 * types.Decimals),
		ValidatorListSize:              int64(21),
		AbsentCommitLimitation:         int64(600),
	}

	voteParam := VoteParam{
		MinStakeIn:                     types.NewCoinFromInt64(1000 * types.Decimals),
		VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
		VoterCoinReturnTimes:           int64(7),
		DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DelegatorCoinReturnTimes:       int64(7),
	}
	proposalParam := ProposalParam{
		ContentCensorshipDecideSec:  int64(7 * 24 * 3600),
		ContentCensorshipPassRatio:  types.NewDecFromRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamExecutionSec: int64(24 * 3600),
		ChangeParamDecideSec:    int64(7 * 24 * 3600),
		ChangeParamPassRatio:    types.NewDecFromRat(70, 100),
		ChangeParamPassVotes:    types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit:   types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideSec:  int64(7 * 24 * 3600),
		ProtocolUpgradePassRatio:  types.NewDecFromRat(80, 100),
		ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
		ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
	}

	coinDayParam := CoinDayParam{
		SecondsToRecoverCoinDay: int64(7 * 24 * 3600),
	}
	bandwidthParam := BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
		VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
		GeneralMsgQuotaRatio:        types.NewDecFromRat(20, 100),
		GeneralMsgEMAFactor:         types.NewDecFromRat(1, 10),
		AppMsgQuotaRatio:            types.NewDecFromRat(80, 100),
		AppMsgEMAFactor:             types.NewDecFromRat(1, 10),
		ExpectedMaxMPS:              types.NewDecFromRat(1000, 1),
		MsgFeeFactorA:               types.NewDecFromRat(6, 1),
		MsgFeeFactorB:               types.NewDecFromRat(10, 1),
		MaxMPSDecayRate:             types.NewDecFromRat(99, 100),
		AppBandwidthPoolSize:        types.NewDecFromRat(10, 1),
		AppVacancyFactor:             types.NewDecFromRat(69, 100),
		AppPunishmentFactor:          types.NewDecFromRat(14, 5),
	}
	accountParam := AccountParam{
		MinimumBalance:               types.NewCoinFromInt64(0),
		RegisterFee:                  types.NewCoinFromInt64(1 * types.Decimals),
		FirstDepositFullCoinDayLimit: types.NewCoinFromInt64(1 * types.Decimals),
		MaxNumFrozenMoney:            10,
	}
	postParam := PostParam{
		ReportOrUpvoteIntervalSec: int64(24 * 3600),
		PostIntervalSec:           int64(600),
		MaxReportReputation:       types.NewCoinFromInt64(100 * types.Decimals),
	}
	checkStorage(t, ctx, ph, globalAllocationParam, infraInternalAllocationParam,
		developerParam, validatorParam, voteParam,
		proposalParam, coinDayParam, bandwidthParam, accountParam, postParam)
}

func TestInitParamFromConfig(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	globalAllocationParam := GlobalAllocationParam{
		GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
		InfraAllocation:          types.NewDecFromRat(20, 100),
		ContentCreatorAllocation: types.NewDecFromRat(65, 100),
		DeveloperAllocation:      types.NewDecFromRat(10, 100),
		ValidatorAllocation:      types.NewDecFromRat(5, 100),
	}

	infraInternalAllocationParam := InfraInternalAllocationParam{
		StorageAllocation: types.NewDecFromRat(50, 100),
		CDNAllocation:     types.NewDecFromRat(50, 100),
	}

	developerParam := DeveloperParam{
		DeveloperMinDeposit:            types.NewCoinFromInt64(1000000 * types.Decimals),
		DeveloperCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DeveloperCoinReturnTimes:       int64(7),
	}

	validatorParam := ValidatorParam{
		ValidatorMinWithdraw:           types.NewCoinFromInt64(1 * types.Decimals),
		ValidatorMinVotingDeposit:      types.NewCoinFromInt64(300000 * types.Decimals),
		ValidatorMinCommittingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissVote:                types.NewCoinFromInt64(20000 * types.Decimals),
		PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:               types.NewCoinFromInt64(1000000 * types.Decimals),
		ValidatorListSize:              int64(21),
		AbsentCommitLimitation:         int64(600),
	}

	voteParam := VoteParam{
		MinStakeIn:                     types.NewCoinFromInt64(1000 * types.Decimals),
		VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
		VoterCoinReturnTimes:           int64(7),
		DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DelegatorCoinReturnTimes:       int64(7),
	}
	proposalParam := ProposalParam{
		ContentCensorshipDecideSec:  int64(7 * 24 * 3600),
		ContentCensorshipPassRatio:  types.NewDecFromRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamExecutionSec: int64(24 * 3600),
		ChangeParamDecideSec:    int64(7 * 24 * 3600),
		ChangeParamPassRatio:    types.NewDecFromRat(70, 100),
		ChangeParamPassVotes:    types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit:   types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideSec:  int64(7 * 24 * 3600),
		ProtocolUpgradePassRatio:  types.NewDecFromRat(80, 100),
		ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
		ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
	}

	coinDayParam := CoinDayParam{
		SecondsToRecoverCoinDay: int64(7 * 24 * 3600),
	}
	bandwidthParam := BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
		VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
		GeneralMsgQuotaRatio:        types.NewDecFromRat(20, 100),
		GeneralMsgEMAFactor:         types.NewDecFromRat(1, 10),
		AppMsgQuotaRatio:            types.NewDecFromRat(80, 100),
		AppMsgEMAFactor:             types.NewDecFromRat(1, 10),
		ExpectedMaxMPS:              types.NewDecFromRat(1000, 1),
		MsgFeeFactorA:               types.NewDecFromRat(6, 1),
		MsgFeeFactorB:               types.NewDecFromRat(10, 1),
		MaxMPSDecayRate:             types.NewDecFromRat(99, 100),
		AppBandwidthPoolSize:        types.NewDecFromRat(10, 1),
		AppVacancyFactor:             types.NewDecFromRat(69, 100),
		AppPunishmentFactor:          types.NewDecFromRat(14, 5),
	}
	accountParam := AccountParam{
		MinimumBalance:               types.NewCoinFromInt64(0),
		RegisterFee:                  types.NewCoinFromInt64(1 * types.Decimals),
		FirstDepositFullCoinDayLimit: types.NewCoinFromInt64(1 * types.Decimals),
		MaxNumFrozenMoney:            10,
	}
	postParam := PostParam{
		ReportOrUpvoteIntervalSec: int64(24 * 3600),
		PostIntervalSec:           int64(600),
		MaxReportReputation:       types.NewCoinFromInt64(100 * types.Decimals),
	}
	repParam := ReputationParam{
		BestContentIndexN: 10,
	}

	err := ph.InitParamFromConfig(
		ctx, globalAllocationParam,
		infraInternalAllocationParam,
		postParam,
		developerParam,
		validatorParam,
		voteParam,
		proposalParam,
		coinDayParam,
		bandwidthParam,
		accountParam,
		repParam,
	)
	assert.Nil(t, err)

	checkStorage(t, ctx, ph, globalAllocationParam, infraInternalAllocationParam,
		developerParam, validatorParam, voteParam,
		proposalParam, coinDayParam, bandwidthParam, accountParam, postParam)
}

func checkStorage(t *testing.T, ctx sdk.Context, ph ParamHolder, expectGlobalAllocationParam GlobalAllocationParam,
	expectInfraInternalAllocationParam InfraInternalAllocationParam,
	expectDeveloperParam DeveloperParam,
	expectValidatorParam ValidatorParam, expectVoteParam VoteParam,
	expectProposalParam ProposalParam, expectCoinDayParam CoinDayParam,
	expectBandwidthParam BandwidthParam, expectAccountParam AccountParam,
	expectPostParam PostParam) {
	globalAllocationParam, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalAllocationParam, *globalAllocationParam)

	infraInternalAllocationParam, err := ph.GetInfraInternalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectInfraInternalAllocationParam, *infraInternalAllocationParam)

	developerParam, err := ph.GetDeveloperParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectDeveloperParam, *developerParam)

	validatorParam, err := ph.GetValidatorParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectValidatorParam, *validatorParam)

	voteParam, err := ph.GetVoteParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectVoteParam, *voteParam)

	proposalParam, err := ph.GetProposalParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectProposalParam, *proposalParam)

	coinDayParam, err := ph.GetCoinDayParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectCoinDayParam, *coinDayParam)

	bandwidthParam, err := ph.GetBandwidthParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectBandwidthParam, *bandwidthParam)

	accountParam, err := ph.GetAccountParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectAccountParam, *accountParam)

	postParam, err := ph.GetPostParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectPostParam, *postParam)
}

func TestUpdateGlobalGrowthRate(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()

	testCases := []struct {
		testName         string
		ceiling          sdk.Dec
		floor            sdk.Dec
		updateGrowthRate sdk.Dec
		expectGrowthRate sdk.Dec
	}{
		{
			testName:         "normal update",
			updateGrowthRate: types.NewDecFromRat(98, 1000),
			expectGrowthRate: types.NewDecFromRat(98, 1000),
		},
		{
			testName:         "update to ceiling",
			updateGrowthRate: types.NewDecFromRat(99, 1000),
			expectGrowthRate: types.NewDecFromRat(98, 1000),
		},
		{
			testName:         "update to floor",
			updateGrowthRate: types.NewDecFromRat(29, 1000),
			expectGrowthRate: types.NewDecFromRat(3, 100),
		},
	}
	for _, tc := range testCases {
		globalParam := &GlobalAllocationParam{}
		err := ph.setGlobalAllocationParam(ctx, globalParam)
		assert.Nil(t, err)
		err = ph.UpdateGlobalGrowthRate(ctx, tc.updateGrowthRate)
		assert.Nil(t, err)
		globalParam, err = ph.GetGlobalAllocationParam(ctx)
		assert.Nil(t, err)
		assert.Equal(t, globalParam.GlobalGrowthRate, tc.expectGrowthRate)
	}
}
