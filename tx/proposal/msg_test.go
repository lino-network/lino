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
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
	}

	p2 := param.GlobalAllocationParam{
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{25, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
	}

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
		CDNAllocation:     sdk.Rat{20, 100},
		StorageAllocation: sdk.Rat{80, 100},
	}

	p2 := param.InfraInternalAllocationParam{
		CDNAllocation:     sdk.ZeroRat,
		StorageAllocation: sdk.Rat{101, 100},
	}

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
		VoterMinDeposit:               types.NewCoin(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoin(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}

	p2 := param.VoteParam{
		VoterMinDeposit:               types.NewCoin(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoin(0 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}

	p3 := param.VoteParam{
		VoterMinDeposit:               types.NewCoin(-1 * types.Decimals),
		VoterMinWithdraw:              types.NewCoin(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}

	p4 := param.VoteParam{
		VoterMinDeposit:               types.NewCoin(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoin(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoin(0 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}
	p5 := param.VoteParam{
		VoterMinDeposit:               types.NewCoin(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoin(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(0),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}
	p6 := param.VoteParam{
		VoterMinDeposit:               types.NewCoin(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoin(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(0),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(7),
	}
	p7 := param.VoteParam{
		VoterMinDeposit:               types.NewCoin(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoin(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(-1),
		DelegatorCoinReturnTimes:      int64(7),
	}
	p8 := param.VoteParam{
		VoterMinDeposit:               types.NewCoin(1000 * types.Decimals),
		VoterMinWithdraw:              types.NewCoin(1 * types.Decimals),
		DelegatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		VoterCoinReturnIntervalHr:     int64(7 * 24),
		VoterCoinReturnTimes:          int64(7),
		DelegatorCoinReturnIntervalHr: int64(7 * 24),
		DelegatorCoinReturnTimes:      int64(0),
	}

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
		DeveloperMinDeposit:           types.NewCoin(1 * types.Decimals),
	}

	p2 := param.DeveloperParam{
		DeveloperCoinReturnIntervalHr: int64(7 * 24),
		DeveloperCoinReturnTimes:      int64(-7),
		DeveloperMinDeposit:           types.NewCoin(1 * types.Decimals),
	}

	p3 := param.DeveloperParam{
		DeveloperCoinReturnIntervalHr: int64(7 * 24),
		DeveloperCoinReturnTimes:      int64(0),
		DeveloperMinDeposit:           types.NewCoin(1 * types.Decimals),
	}

	p4 := param.DeveloperParam{
		DeveloperCoinReturnIntervalHr: int64(7 * 24),
		DeveloperCoinReturnTimes:      int64(7),
		DeveloperMinDeposit:           types.NewCoin(-1 * types.Decimals),
	}

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
		ValidatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoin(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(1000 * types.Decimals),
	}

	p2 := param.ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoin(-1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoin(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(1000 * types.Decimals),
	}

	p3 := param.ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(0 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoin(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(1000 * types.Decimals),
	}

	p4 := param.ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(-1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoin(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(1000 * types.Decimals),
	}

	p5 := param.ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(-7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoin(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(1000 * types.Decimals),
	}

	p6 := param.ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(23),
		PenaltyMissVote:               types.NewCoin(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(-200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(1000 * types.Decimals),
	}

	p7 := param.ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(0),
		PenaltyMissVote:               types.NewCoin(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(1000 * types.Decimals),
	}

	p8 := param.ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoin(-200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(1000 * types.Decimals),
	}

	p9 := param.ValidatorParam{
		ValidatorMinWithdraw:          types.NewCoin(1 * types.Decimals),
		ValidatorMinVotingDeposit:     types.NewCoin(3000 * types.Decimals),
		ValidatorMinCommitingDeposit:  types.NewCoin(1000 * types.Decimals),
		ValidatorCoinReturnIntervalHr: int64(7 * 24),
		ValidatorCoinReturnTimes:      int64(7),
		PenaltyMissVote:               types.NewCoin(200 * types.Decimals),
		PenaltyMissCommit:             types.NewCoin(200 * types.Decimals),
		PenaltyByzantine:              types.NewCoin(-10233232300 * types.Decimals),
	}

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
		{NewChangeValidatorParamMsg("", p1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.ChangeValidatorParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)

	}
}
