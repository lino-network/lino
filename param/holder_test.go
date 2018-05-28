package param

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("param")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func TestGlobalAllocationParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := GlobalAllocationParam{
		ContentCreatorAllocation: sdk.NewRat(1, 100),
		InfraAllocation:          sdk.NewRat(1, 100),
		DeveloperAllocation:      sdk.NewRat(1, 100),
		ValidatorAllocation:      sdk.NewRat(97, 100),
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
		StorageAllocation: sdk.NewRat(50, 100),
		CDNAllocation:     sdk.NewRat(50, 100),
	}
	err := ph.setInfraInternalAllocationParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetInfraInternalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Infra internal allocation param should be equal")
}

func TestEvaluateOfContenValueParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := EvaluateOfContentValueParam{
		ConsumptionTimeAdjustBase:      3153600,
		ConsumptionTimeAdjustOffset:    5,
		NumOfConsumptionOnAuthorOffset: 7,
		TotalAmountOfConsumptionBase:   1000 * types.Decimals,
		TotalAmountOfConsumptionOffset: 5,
		AmountOfConsumptionExponent:    sdk.NewRat(8, 10),
	}
	err := ph.setEvaluateOfContentValueParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr, err := ph.GetEvaluateOfContentValueParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, parameter, *resultPtr, "Evaluate of content value param should be equal")
}

func TestDeveloperParam(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	parameter := DeveloperParam{
		DeveloperMinDeposit:           types.NewCoinFromInt64(100000 * types.Decimals),
		DeveloperCoinReturnIntervalHr: int64(7 * 24),
		DeveloperCoinReturnTimes:      int64(7),
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
		ValidatorMinWithdraw:          types.NewCoinFromInt64(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoinFromInt64(300000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoinFromInt64(20000 * types.Decimals),
		PenaltyMissCommit:             types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoinFromInt64(1000000 * types.Decimals),
		ValidatorListSize:             int64(21),
		AbsentCommitLimitation:        int64(100),
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
		VoterMinDeposit:               types.NewCoinFromInt64(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoinFromInt64(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoinFromInt64(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
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
		ContentCensorshipDecideHr:   int64(24 * 7),
		ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamDecideHr:   int64(24 * 7),
		ChangeParamPassRatio:  sdk.NewRat(70, 100),
		ChangeParamPassVotes:  types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit: types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideHr:   int64(24 * 7),
		ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
		ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
		ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),

		NextProposalID: int64(0),
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
		DaysToRecoverCoinDayStake:    int64(7),
		SecondsToRecoverCoinDayStake: int64(7 * 24 * 3600),
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
		MinimumBalance: types.NewCoinFromInt64(1 * types.Decimals),
		RegisterFee:    types.NewCoinFromInt64(1 * types.Decimals),
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
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(65, 100),
		DeveloperAllocation:      sdk.NewRat(10, 100),
		ValidatorAllocation:      sdk.NewRat(5, 100),
	}

	infraInternalAllocationParam := InfraInternalAllocationParam{
		StorageAllocation: sdk.NewRat(50, 100),
		CDNAllocation:     sdk.NewRat(50, 100),
	}

	evaluateOfContentValueParam := EvaluateOfContentValueParam{
		ConsumptionTimeAdjustBase:      3153600,
		ConsumptionTimeAdjustOffset:    5,
		NumOfConsumptionOnAuthorOffset: 7,
		TotalAmountOfConsumptionBase:   1000 * types.Decimals,
		TotalAmountOfConsumptionOffset: 5,
		AmountOfConsumptionExponent:    sdk.NewRat(8, 10),
	}

	developerParam := DeveloperParam{
		DeveloperMinDeposit:           types.NewCoinFromInt64(1000000 * types.Decimals),
		DeveloperCoinReturnIntervalHr: int64(7 * 24),
		DeveloperCoinReturnTimes:      int64(7),
	}

	validatorParam := ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoinFromInt64(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoinFromInt64(300000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoinFromInt64(20000 * types.Decimals),
		PenaltyMissCommit:             types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoinFromInt64(1000000 * types.Decimals),
		ValidatorListSize:             int64(21),
		AbsentCommitLimitation:        int64(100),
	}

	voteParam := VoteParam{
		VoterMinDeposit:               types.NewCoinFromInt64(2000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoinFromInt64(2 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoinFromInt64(2 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}
	proposalParam := ProposalParam{
		ContentCensorshipDecideHr:   int64(24 * 7),
		ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamDecideHr:   int64(24 * 7),
		ChangeParamPassRatio:  sdk.NewRat(70, 100),
		ChangeParamPassVotes:  types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit: types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideHr:   int64(24 * 7),
		ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
		ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
		ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),

		NextProposalID: int64(0),
	}

	coinDayParam := CoinDayParam{
		DaysToRecoverCoinDayStake:    int64(7),
		SecondsToRecoverCoinDayStake: int64(7 * 24 * 3600),
	}
	bandwidthParam := BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
	}
	accountParam := AccountParam{
		MinimumBalance: types.NewCoinFromInt64(1 * types.Decimals),
		RegisterFee:    types.NewCoinFromInt64(1 * types.Decimals),
	}
	checkStorage(t, ctx, ph, globalAllocationParam, infraInternalAllocationParam, evaluateOfContentValueParam,
		developerParam, validatorParam, voteParam, proposalParam, coinDayParam, bandwidthParam, accountParam)
}

func checkStorage(t *testing.T, ctx sdk.Context, ph ParamHolder, expectGlobalAllocationParam GlobalAllocationParam,
	expectInfraInternalAllocationParam InfraInternalAllocationParam,
	expectEvaluateOfContentValueParam EvaluateOfContentValueParam, expectDeveloperParam DeveloperParam,
	expectValidatorParam ValidatorParam, expectVoteParam VoteParam,
	expectProposalParam ProposalParam, expectCoinDayParam CoinDayParam,
	expectBandwidthParam BandwidthParam, expectAccountParam AccountParam) {
	evaluateOfContentValueParam, err := ph.GetEvaluateOfContentValueParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectEvaluateOfContentValueParam, *evaluateOfContentValueParam)

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
}
