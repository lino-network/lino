//nolint:unused
package param

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/lino-network/lino/types"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("param")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

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
		ValidatorMinDeposit:            types.NewCoinFromInt64(200000 * types.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:               types.NewCoinFromInt64(1000 * types.Decimals),
		AbsentCommitLimitation:         int64(600), // 30min
		OncallSize:                     int64(22),
		StandbySize:                    int64(7),
		ValidatorRevokePendingSec:      int64(7 * 24 * 3600),
		OncallInflationWeight:          int64(2),
		StandbyInflationWeight:         int64(1),
		MaxVotedValidators:             int64(3),
	}
	err := ph.setValidatorParam(ctx, &parameter)
	assert.Nil(t, err)

	resultPtr := ph.GetValidatorParam(ctx)
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
		ExpectedMaxMPS:              types.NewDecFromRat(300, 1),
		MsgFeeFactorA:               types.NewDecFromRat(6, 1),
		MsgFeeFactorB:               types.NewDecFromRat(10, 1),
		MaxMPSDecayRate:             types.NewDecFromRat(99, 100),
		AppBandwidthPoolSize:        types.NewDecFromRat(10, 1),
		AppVacancyFactor:            types.NewDecFromRat(69, 100),
		AppPunishmentFactor:         types.NewDecFromRat(14, 5),
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

	err := ph.InitParam(ctx)
	if err != nil {
		panic(err)
	}

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
		ValidatorMinDeposit:            types.NewCoinFromInt64(200000 * types.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:               types.NewCoinFromInt64(1000 * types.Decimals),
		AbsentCommitLimitation:         int64(600), // 30min
		OncallSize:                     int64(22),
		StandbySize:                    int64(7),
		ValidatorRevokePendingSec:      int64(7 * 24 * 3600),
		OncallInflationWeight:          int64(2),
		StandbyInflationWeight:         int64(1),
		MaxVotedValidators:             int64(3),
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
		ExpectedMaxMPS:              types.NewDecFromRat(300, 1),
		MsgFeeFactorA:               types.NewDecFromRat(6, 1),
		MsgFeeFactorB:               types.NewDecFromRat(10, 1),
		MaxMPSDecayRate:             types.NewDecFromRat(99, 100),
		AppBandwidthPoolSize:        types.NewDecFromRat(10, 1),
		AppVacancyFactor:            types.NewDecFromRat(69, 100),
		AppPunishmentFactor:         types.NewDecFromRat(14, 5),
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
		BestContentIndexN: 200,
		UserMaxN:          50,
	}
	priceParam := PriceParam{
		TestnetMode:     true,
		UpdateEverySec:  int64(time.Hour.Seconds()),
		FeedEverySec:    int64((10 * time.Minute).Seconds()),
		HistoryMaxLen:   71,
		PenaltyMissFeed: types.NewCoinFromInt64(10000 * types.Decimals),
	}

	checkStorage(t, ctx, ph, globalAllocationParam, infraInternalAllocationParam,
		developerParam, validatorParam, voteParam,
		proposalParam, coinDayParam, bandwidthParam, accountParam, postParam, repParam, priceParam)
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
		ValidatorMinDeposit:            types.NewCoinFromInt64(200000 * types.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:               types.NewCoinFromInt64(1000 * types.Decimals),
		AbsentCommitLimitation:         int64(600), // 30min
		OncallSize:                     int64(22),
		StandbySize:                    int64(7),
		ValidatorRevokePendingSec:      int64(7 * 24 * 3600),
		OncallInflationWeight:          int64(2),
		StandbyInflationWeight:         int64(1),
		MaxVotedValidators:             int64(3),
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
		AppVacancyFactor:            types.NewDecFromRat(69, 100),
		AppPunishmentFactor:         types.NewDecFromRat(14, 5),
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
		BestContentIndexN: 200,
		UserMaxN:          40,
	}
	priceParam := PriceParam{
		UpdateEverySec:  int64(time.Hour.Seconds()),
		FeedEverySec:    int64((10 * time.Minute).Seconds()),
		HistoryMaxLen:   123,
		PenaltyMissFeed: types.NewCoinFromInt64(10000 * types.Decimals),
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
		priceParam,
	)
	assert.Nil(t, err)

	checkStorage(t, ctx, ph, globalAllocationParam, infraInternalAllocationParam,
		developerParam, validatorParam, voteParam,
		proposalParam, coinDayParam, bandwidthParam, accountParam, postParam, repParam, priceParam)
}

func checkStorage(t *testing.T, ctx sdk.Context, ph ParamHolder, expectGlobalAllocationParam GlobalAllocationParam,
	expectInfraInternalAllocationParam InfraInternalAllocationParam,
	expectDeveloperParam DeveloperParam,
	expectValidatorParam ValidatorParam, expectVoteParam VoteParam,
	expectProposalParam ProposalParam, expectCoinDayParam CoinDayParam,
	expectBandwidthParam BandwidthParam, expectAccountParam AccountParam,
	expectPostParam PostParam,
	expectedRepParam ReputationParam,
	expectedPriceParam PriceParam,
) {
	globalAllocationParam, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalAllocationParam, *globalAllocationParam)

	infraInternalAllocationParam, err := ph.GetInfraInternalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectInfraInternalAllocationParam, *infraInternalAllocationParam)

	developerParam, err := ph.GetDeveloperParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectDeveloperParam, *developerParam)

	validatorParam := ph.GetValidatorParam(ctx)
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

	repParam := ph.GetReputationParam(ctx)
	assert.Equal(t, expectedRepParam, *repParam)

	priceParam := ph.GetPriceParam(ctx)
	assert.Equal(t, expectedPriceParam, *priceParam)
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
