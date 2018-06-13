package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestChangeGlobalAllocationParamMsg(t *testing.T) {
	p1 := param.GlobalAllocationParam{
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(55, 100),
		DeveloperAllocation:      sdk.NewRat(20, 100),
		ValidatorAllocation:      sdk.NewRat(5, 100),
	}
	p2 := p1
	p2.DeveloperAllocation = sdk.NewRat(25, 100)

	cases := []struct {
		ChangeGlobalAllocationParamMsg ChangeGlobalAllocationParamMsg
		expectError                    sdk.Error
	}{
		{NewChangeGlobalAllocationParamMsg("user1", p1), nil},
		{NewChangeGlobalAllocationParamMsg("user1", p2), ErrIllegalParameter()},
		{NewChangeGlobalAllocationParamMsg("", p1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.ChangeGlobalAllocationParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              sdk.Msg
		expectPermission types.Permission
	}{
		"change evaluate of content value param": {
			NewChangeEvaluateOfContentValueParamMsg("creator",
				param.EvaluateOfContentValueParam{}),
			types.TransactionPermission},
		"change global allocation param": {
			NewChangeGlobalAllocationParamMsg("creator",
				param.GlobalAllocationParam{}),
			types.TransactionPermission},
		"change infra internal allocation param": {
			NewChangeInfraInternalAllocationParamMsg("creator",
				param.InfraInternalAllocationParam{}),
			types.TransactionPermission},
		"change vote param": {
			NewChangeVoteParamMsg("creator", param.VoteParam{}),
			types.TransactionPermission},
		"change proposal param": {
			NewChangeProposalParamMsg("creator", param.ProposalParam{}),
			types.TransactionPermission},
		"change developer param": {
			NewChangeDeveloperParamMsg("creator", param.DeveloperParam{}),
			types.TransactionPermission},
		"change validator param": {
			NewChangeValidatorParamMsg("creator", param.ValidatorParam{}),
			types.TransactionPermission},
		"change coinday param": {
			NewChangeCoinDayParamMsg("creator", param.CoinDayParam{}),
			types.TransactionPermission},
		"change bandwidth param": {
			NewChangeBandwidthParamMsg("creator", param.BandwidthParam{}),
			types.TransactionPermission},
		"change account param": {
			NewChangeAccountParamMsg("creator", param.AccountParam{}),
			types.TransactionPermission},
	}

	for testName, cs := range cases {
		permissionLevel := cs.msg.Get(types.PermissionLevel)
		if permissionLevel == nil {
			if cs.expectPermission != types.PostPermission {
				t.Errorf(
					"%s: expect permission incorrect, expect %v, got %v",
					testName, cs.expectPermission, types.PostPermission)
				return
			} else {
				continue
			}
		}
		permission, ok := permissionLevel.(types.Permission)
		assert.Equal(t, ok, true)
		if cs.expectPermission != permission {
			t.Errorf(
				"%s: expect permission incorrect, expect %v, got %v",
				testName, cs.expectPermission, permission)
			return
		}
	}
}

func TestChangeInfraInternalAllocationParamMsg(t *testing.T) {
	p1 := param.InfraInternalAllocationParam{
		CDNAllocation:     sdk.NewRat(20, 100),
		StorageAllocation: sdk.NewRat(80, 100),
	}

	p2 := p1
	p2.StorageAllocation = sdk.NewRat(101, 100)
	cases := []struct {
		ChangeInfraInternalAllocationParamMsg ChangeInfraInternalAllocationParamMsg
		expectError                           sdk.Error
	}{
		{NewChangeInfraInternalAllocationParamMsg("user1", p1), nil},
		{NewChangeInfraInternalAllocationParamMsg("user1", p2), ErrIllegalParameter()},
		{NewChangeInfraInternalAllocationParamMsg("", p1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.ChangeInfraInternalAllocationParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestChangeVoteParamMsg(t *testing.T) {
	p1 := param.VoteParam{
		VoterMinDeposit:               types.NewCoinFromInt64(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoinFromInt64(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoinFromInt64(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}

	p2 := p1
	p2.VoterMinDeposit = types.NewCoinFromInt64(-1 * types.Decimals)

	p3 := p1
	p3.VoterMinWithdraw = types.NewCoinFromInt64(0 * types.Decimals)

	p4 := p1
	p4.DelegatorMinWithdraw = types.NewCoinFromInt64(0 * types.Decimals)

	p5 := p1
	p5.VoterCoinReturnIntervalHr = int64(0)

	p6 := p1
	p6.VoterCoinReturnTimes = int64(0)

	p7 := p1
	p7.DelegatorCoinReturnIntervalHr = int64(-1)

	p8 := p1
	p8.DelegatorCoinReturnTimes = int64(0)

	cases := []struct {
		ChangeVoteParamMsg ChangeVoteParamMsg
		expectError        sdk.Error
	}{
		{NewChangeVoteParamMsg("user1", p1), nil},
		{NewChangeVoteParamMsg("user1", p2), ErrIllegalParameter()},
		{NewChangeVoteParamMsg("user1", p3), ErrIllegalParameter()},
		{NewChangeVoteParamMsg("user1", p4), ErrIllegalParameter()},
		{NewChangeVoteParamMsg("user1", p5), ErrIllegalParameter()},
		{NewChangeVoteParamMsg("user1", p6), ErrIllegalParameter()},
		{NewChangeVoteParamMsg("user1", p7), ErrIllegalParameter()},
		{NewChangeVoteParamMsg("user1", p8), ErrIllegalParameter()},
		{NewChangeVoteParamMsg("", p1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.ChangeVoteParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestChangeDeveloperParamMsg(t *testing.T) {
	p1 := param.DeveloperParam{
		DeveloperCoinReturnIntervalHr: int64(7 * 24),
		DeveloperCoinReturnTimes:      int64(7),
		DeveloperMinDeposit:           types.NewCoinFromInt64(1 * types.Decimals),
	}

	p2 := p1
	p2.DeveloperCoinReturnTimes = int64(-7)

	p3 := p1
	p3.DeveloperCoinReturnIntervalHr = int64(0)

	p4 := p1
	p4.DeveloperMinDeposit = types.NewCoinFromInt64(-1 * types.Decimals)

	cases := []struct {
		ChangeDeveloperParamMsg ChangeDeveloperParamMsg
		expectError             sdk.Error
	}{
		{NewChangeDeveloperParamMsg("user1", p1), nil},
		{NewChangeDeveloperParamMsg("user1", p2), ErrIllegalParameter()},
		{NewChangeDeveloperParamMsg("user1", p3), ErrIllegalParameter()},
		{NewChangeDeveloperParamMsg("user1", p4), ErrIllegalParameter()},
		{NewChangeDeveloperParamMsg("", p1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.ChangeDeveloperParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestChangeValidatorParamMsg(t *testing.T) {
	p1 := param.ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoinFromInt64(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoinFromInt64(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoinFromInt64(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoinFromInt64(1000 * types.Decimals),
		ValidatorListSize:             int64(21),
		AbsentCommitLimitation:        int64(100),
	}

	p2 := p1
	p2.ValidatorMinWithdraw = types.NewCoinFromInt64(-1 * types.Decimals)

	p3 := p1
	p3.ValidatorMinVotingDeposit = types.NewCoinFromInt64(0 * types.Decimals)

	p4 := p1
	p4.ValidatorMinCommitingDeposit = types.NewCoinFromInt64(-1000 * types.Decimals)

	p5 := p1
	p5.ValidatorCoinReturnIntervalHr = int64(-7 * 24)

	p6 := p1
	p6.ValidatorCoinReturnTimes = int64(0)

	p7 := p1
	p7.PenaltyMissVote = types.NewCoinFromInt64(-200 * types.Decimals)

	p8 := p1
	p8.PenaltyByzantine = types.NewCoinFromInt64(-10233232300 * types.Decimals)

	p9 := p1
	p9.PenaltyMissCommit = types.NewCoinFromInt64(0 * types.Decimals)

	p10 := p1
	p10.AbsentCommitLimitation = int64(0)

	p11 := p1
	p11.ValidatorListSize = int64(-1)

	cases := []struct {
		ChangeValidatorParamMsg ChangeValidatorParamMsg
		expectError             sdk.Error
	}{
		{NewChangeValidatorParamMsg("user1", p1), nil},
		{NewChangeValidatorParamMsg("user1", p2), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("user1", p3), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("user1", p4), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("user1", p5), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("user1", p6), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("user1", p7), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("user1", p8), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("user1", p9), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("user1", p10), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("user1", p11), ErrIllegalParameter()},
		{NewChangeValidatorParamMsg("", p1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.ChangeValidatorParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestChangeProposalParamMsg(t *testing.T) {
	p1 := param.ProposalParam{
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

	p2 := p1
	p2.ContentCensorshipDecideHr = int64(-24 * 7)

	p3 := p1
	p3.ContentCensorshipPassRatio = sdk.NewRat(150, 100)

	p4 := p1
	p4.ContentCensorshipPassVotes = types.NewCoinFromInt64(-10000 * types.Decimals)

	p5 := p1
	p5.ContentCensorshipMinDeposit = types.NewCoinFromInt64(-100 * types.Decimals)

	p6 := p1
	p6.ChangeParamDecideHr = int64(-24 * 7)

	p7 := p1
	p7.ChangeParamPassRatio = sdk.NewRat(0, 8)

	p8 := p1
	p8.ChangeParamPassVotes = types.NewCoinFromInt64(0 * types.Decimals)

	p9 := p1
	p9.ChangeParamMinDeposit = types.NewCoinFromInt64(-100000 * types.Decimals)

	p10 := p1
	p10.ProtocolUpgradeDecideHr = int64(0)

	p11 := p1
	p11.ProtocolUpgradePassRatio = sdk.NewRat(0, 100)

	p12 := p1
	p12.ProtocolUpgradePassVotes = types.NewCoinFromInt64(-10000000 * types.Decimals)

	p13 := p1
	p13.ProtocolUpgradeMinDeposit = types.NewCoinFromInt64(-1000000 * types.Decimals)

	cases := []struct {
		ChangeProposalParamMsg ChangeProposalParamMsg
		expectError            sdk.Error
	}{
		{NewChangeProposalParamMsg("user1", p1), nil},
		{NewChangeProposalParamMsg("", p1), ErrInvalidUsername()},
		{NewChangeProposalParamMsg("user1", p2), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p3), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p4), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p5), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p6), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p7), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p8), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p9), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p10), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p11), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p12), ErrIllegalParameter()},
		{NewChangeProposalParamMsg("user1", p13), ErrIllegalParameter()},
	}

	for _, cs := range cases {
		result := cs.ChangeProposalParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestChangeAccountParamMsg(t *testing.T) {
	p1 := param.AccountParam{
		MinimumBalance: types.NewCoinFromInt64(1 * types.Decimals),
		RegisterFee:    types.NewCoinFromInt64(1 * types.Decimals),
	}

	p2 := p1
	p2.MinimumBalance = types.NewCoinFromInt64(0)

	p3 := p1
	p3.RegisterFee = types.NewCoinFromInt64(0)

	p4 := p1
	p4.RegisterFee = types.NewCoinFromInt64(-1)

	p5 := p1
	p5.RegisterFee = types.NewCoinFromInt64(-1)

	cases := []struct {
		changeAccountParamMsg ChangeAccountParamMsg
		expectError           sdk.Error
	}{
		{NewChangeAccountParamMsg("user1", p1), nil},
		{NewChangeAccountParamMsg("us", p1), ErrInvalidUsername()},
		{NewChangeAccountParamMsg("user1user1user1user1user1user1", p1), ErrInvalidUsername()},
		{NewChangeAccountParamMsg("user1", p2), nil},
		{NewChangeAccountParamMsg("user1", p3), nil},
		{NewChangeAccountParamMsg("user1", p4), ErrIllegalParameter()},
		{NewChangeAccountParamMsg("user1", p5), ErrIllegalParameter()},
	}

	for _, cs := range cases {
		result := cs.changeAccountParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestChangeBandwidthParamMsg(t *testing.T) {
	p1 := param.BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
	}

	p2 := p1
	p2.SecondsToRecoverBandwidth = int64(-1)

	p3 := p1
	p3.CapacityUsagePerTransaction = types.NewCoinFromInt64(-1)

	cases := []struct {
		changeBandwidthParamMsg ChangeBandwidthParamMsg
		expectError             sdk.Error
	}{
		{NewChangeBandwidthParamMsg("user1", p1), nil},
		{NewChangeBandwidthParamMsg("us", p1), ErrInvalidUsername()},
		{NewChangeBandwidthParamMsg("user1user1user1user1user1user1", p1), ErrInvalidUsername()},
		{NewChangeBandwidthParamMsg("user1", p2), ErrIllegalParameter()},
		{NewChangeBandwidthParamMsg("user1", p3), ErrIllegalParameter()},
	}

	for _, cs := range cases {
		result := cs.changeBandwidthParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestChangeEvaluateOfContentValueParamMsg(t *testing.T) {
	p1 := param.EvaluateOfContentValueParam{
		ConsumptionTimeAdjustBase:      3153600,
		ConsumptionTimeAdjustOffset:    5,
		NumOfConsumptionOnAuthorOffset: 7,
		TotalAmountOfConsumptionBase:   1000 * types.Decimals,
		TotalAmountOfConsumptionOffset: 5,
		AmountOfConsumptionExponent:    sdk.NewRat(8, 10),
	}

	p2 := p1
	p2.ConsumptionTimeAdjustBase = 0

	p3 := p1
	p3.TotalAmountOfConsumptionBase = 0

	cases := []struct {
		changeAccountParamMsg ChangeEvaluateOfContentValueParamMsg
		expectError           sdk.Error
	}{
		{NewChangeEvaluateOfContentValueParamMsg("user1", p1), nil},
		{NewChangeEvaluateOfContentValueParamMsg("user1", p2), ErrIllegalParameter()},
		{NewChangeEvaluateOfContentValueParamMsg("user1", p3), ErrIllegalParameter()},
		{NewChangeEvaluateOfContentValueParamMsg("us", p1), ErrInvalidUsername()},
		{NewChangeEvaluateOfContentValueParamMsg("user1user1user1user1user1", p1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.changeAccountParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestChangeCoinDayParamMsg(t *testing.T) {
	p1 := param.CoinDayParam{
		DaysToRecoverCoinDayStake:    7,
		SecondsToRecoverCoinDayStake: 7 * 24 * 3600,
	}

	p2 := p1
	p2.DaysToRecoverCoinDayStake = 0

	p3 := p1
	p3.SecondsToRecoverCoinDayStake = 0

	p4 := p1
	p2.DaysToRecoverCoinDayStake = 1
	p4.SecondsToRecoverCoinDayStake = 3600

	cases := []struct {
		changeCoinDayParamMsg ChangeCoinDayParamMsg
		expectError           sdk.Error
	}{
		{NewChangeCoinDayParamMsg("user1", p1), nil},
		{NewChangeCoinDayParamMsg("us", p1), ErrInvalidUsername()},
		{NewChangeCoinDayParamMsg("user1user1user1user1user1user1", p1), ErrInvalidUsername()},
		{NewChangeCoinDayParamMsg("user1", p2), ErrIllegalParameter()},
		{NewChangeCoinDayParamMsg("user1", p3), ErrIllegalParameter()},
		{NewChangeCoinDayParamMsg("user1", p4), ErrIllegalParameter()},
	}

	for _, cs := range cases {
		result := cs.changeCoinDayParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestDeletePostContentMsg(t *testing.T) {
	cases := []struct {
		deletePostContentMsg DeletePostContentMsg
		expectError          sdk.Error
	}{
		{NewDeletePostContentMsg("user1", "permLink", "reason"), nil},
		{NewDeletePostContentMsg("us", "permLink", "reason"), ErrInvalidUsername()},
		{NewDeletePostContentMsg("user1user1user1user1user1user1", "permLink", "reason"), ErrInvalidUsername()},
		{NewDeletePostContentMsg("user1", "", "reason"), ErrInvalidPermLink()},
	}

	for _, cs := range cases {
		result := cs.deletePostContentMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}
