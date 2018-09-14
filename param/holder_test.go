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
		GlobalGrowthRate:         sdk.NewRat(98, 1000),
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
		VoterMinWithdraw:               types.NewCoinFromInt64(1 * types.Decimals),
		DelegatorMinWithdraw:           types.NewCoinFromInt64(1 * types.Decimals),
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
		ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamExecutionSec: int64(24 * 3600),
		ChangeParamDecideSec:    int64(7 * 24 * 3600),
		ChangeParamPassRatio:    sdk.NewRat(70, 100),
		ChangeParamPassVotes:    types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit:   types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideSec:  int64(7 * 24 * 3600),
		ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
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
		GlobalGrowthRate:         sdk.NewRat(98, 1000),
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
		VoterMinWithdraw:               types.NewCoinFromInt64(2 * types.Decimals),
		DelegatorMinWithdraw:           types.NewCoinFromInt64(2 * types.Decimals),
		VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
		VoterCoinReturnTimes:           int64(7),
		DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DelegatorCoinReturnTimes:       int64(7),
	}
	proposalParam := ProposalParam{
		ContentCensorshipDecideSec:  int64(7 * 24 * 3600),
		ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamExecutionSec: int64(24 * 3600),
		ChangeParamDecideSec:    int64(7 * 24 * 3600),
		ChangeParamPassRatio:    sdk.NewRat(70, 100),
		ChangeParamPassVotes:    types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit:   types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideSec:  int64(7 * 24 * 3600),
		ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
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
	}
	checkStorage(t, ctx, ph, globalAllocationParam, infraInternalAllocationParam,
		evaluateOfContentValueParam, developerParam, validatorParam, voteParam,
		proposalParam, coinDayParam, bandwidthParam, accountParam, postParam)
}

func TestInitParamFromConfig(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()
	globalAllocationParam := GlobalAllocationParam{
		GlobalGrowthRate:         sdk.NewRat(98, 1000),
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
		VoterMinWithdraw:               types.NewCoinFromInt64(2 * types.Decimals),
		DelegatorMinWithdraw:           types.NewCoinFromInt64(2 * types.Decimals),
		VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
		VoterCoinReturnTimes:           int64(7),
		DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DelegatorCoinReturnTimes:       int64(7),
	}
	proposalParam := ProposalParam{
		ContentCensorshipDecideSec:  int64(7 * 24 * 3600),
		ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamExecutionSec: int64(24 * 3600),
		ChangeParamDecideSec:    int64(7 * 24 * 3600),
		ChangeParamPassRatio:    sdk.NewRat(70, 100),
		ChangeParamPassVotes:    types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit:   types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideSec:  int64(7 * 24 * 3600),
		ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
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
	}

	err := ph.InitParamFromConfig(
		ctx, globalAllocationParam,
		infraInternalAllocationParam,
		postParam,
		evaluateOfContentValueParam,
		developerParam,
		validatorParam,
		voteParam,
		proposalParam,
		coinDayParam,
		bandwidthParam,
		accountParam,
	)
	assert.Nil(t, err)

	checkStorage(t, ctx, ph, globalAllocationParam, infraInternalAllocationParam,
		evaluateOfContentValueParam, developerParam, validatorParam, voteParam,
		proposalParam, coinDayParam, bandwidthParam, accountParam, postParam)
}

func checkStorage(t *testing.T, ctx sdk.Context, ph ParamHolder, expectGlobalAllocationParam GlobalAllocationParam,
	expectInfraInternalAllocationParam InfraInternalAllocationParam,
	expectEvaluateOfContentValueParam EvaluateOfContentValueParam, expectDeveloperParam DeveloperParam,
	expectValidatorParam ValidatorParam, expectVoteParam VoteParam,
	expectProposalParam ProposalParam, expectCoinDayParam CoinDayParam,
	expectBandwidthParam BandwidthParam, expectAccountParam AccountParam,
	expectPostParam PostParam) {
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

	postParam, err := ph.GetPostParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectPostParam, *postParam)
}

func TestUpdateGlobalGrowthRate(t *testing.T) {
	ph := NewParamHolder(TestKVStoreKey)
	ctx := getContext()

	testCases := []struct {
		testName         string
		ceiling          sdk.Rat
		floor            sdk.Rat
		updateGrowthRate sdk.Rat
		expectGrowthRate sdk.Rat
	}{
		{
			testName:         "normal update",
			updateGrowthRate: sdk.NewRat(98, 1000),
			expectGrowthRate: sdk.NewRat(98, 1000),
		},
		{
			testName:         "update to ceiling",
			updateGrowthRate: sdk.NewRat(99, 1000),
			expectGrowthRate: sdk.NewRat(98, 1000),
		},
		{
			testName:         "update to floor",
			updateGrowthRate: sdk.NewRat(29, 1000),
			expectGrowthRate: sdk.NewRat(3, 100),
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
