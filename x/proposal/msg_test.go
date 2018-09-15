package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// len of 1000
	maxLenOfUTF8Reason = `
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌123`

	// len of 1001
	tooLongOfUTF8Reason = `
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌1234`
)

func TestVoteProposalMsg(t *testing.T) {
	testCases := []struct {
		testName        string
		voteProposalMsg VoteProposalMsg
		expectedError   sdk.Error
	}{
		{
			testName:        "normal case",
			voteProposalMsg: NewVoteProposalMsg("user1", 1, true),
			expectedError:   nil,
		},
		{
			testName:        "empty username is illegal",
			voteProposalMsg: NewVoteProposalMsg("", 1, true),
			expectedError:   ErrInvalidUsername(),
		},
	}

	for _, tc := range testCases {
		result := tc.voteProposalMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestChangeGlobalAllocationParamMsg(t *testing.T) {
	p1 := param.GlobalAllocationParam{
		GlobalGrowthRate:         sdk.NewRat(98, 1000),
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(55, 100),
		DeveloperAllocation:      sdk.NewRat(20, 100),
		ValidatorAllocation:      sdk.NewRat(5, 100),
	}
	p2 := p1
	p2.DeveloperAllocation = sdk.NewRat(25, 100)

	p3 := p1
	p3.GlobalGrowthRate = sdk.NewRat(2, 100)

	p4 := p1
	p4.GlobalGrowthRate = sdk.NewRat(1, 10)

	testCases := []struct {
		testName                       string
		ChangeGlobalAllocationParamMsg ChangeGlobalAllocationParamMsg
		expectedError                  sdk.Error
	}{
		{
			testName:                       "normal case",
			ChangeGlobalAllocationParamMsg: NewChangeGlobalAllocationParamMsg("user1", p1, ""),
			expectedError:                  nil,
		},
		{
			testName:                       "illegal allocation",
			ChangeGlobalAllocationParamMsg: NewChangeGlobalAllocationParamMsg("user1", p2, ""),
			expectedError:                  ErrIllegalParameter(),
		},
		{
			testName:                       "global growth rate lower than lower bound",
			ChangeGlobalAllocationParamMsg: NewChangeGlobalAllocationParamMsg("user1", p3, ""),
			expectedError:                  ErrIllegalParameter(),
		},
		{
			testName:                       "global growth rate exceed than higher bound",
			ChangeGlobalAllocationParamMsg: NewChangeGlobalAllocationParamMsg("user1", p4, ""),
			expectedError:                  ErrIllegalParameter(),
		},
		{
			testName:                       "empty username is illegal",
			ChangeGlobalAllocationParamMsg: NewChangeGlobalAllocationParamMsg("", p1, ""),
			expectedError:                  ErrInvalidUsername(),
		},
		{
			testName: "reason is too long",
			ChangeGlobalAllocationParamMsg: NewChangeGlobalAllocationParamMsg(
				"user1", p1, string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			ChangeGlobalAllocationParamMsg: NewChangeGlobalAllocationParamMsg(
				"user1", p1, tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.ChangeGlobalAllocationParamMsg.ValidateBasic()
		if !assert.Equal(t, tc.expectedError, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
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

	p3 := p1
	p3.StorageAllocation = sdk.NewRat(-1, 100)
	p3.CDNAllocation = sdk.NewRat(101, 100)

	testCases := []struct {
		testName                              string
		ChangeInfraInternalAllocationParamMsg ChangeInfraInternalAllocationParamMsg
		expectedError                         sdk.Error
	}{
		{
			testName: "normal case",
			ChangeInfraInternalAllocationParamMsg: NewChangeInfraInternalAllocationParamMsg("user1", p1, ""),
			expectedError:                         nil,
		},
		{
			testName: "illegal parameter (sum of allocation doesn't equal to 1)",
			ChangeInfraInternalAllocationParamMsg: NewChangeInfraInternalAllocationParamMsg("user1", p2, ""),
			expectedError:                         ErrIllegalParameter(),
		},
		{
			testName: "illegal parameter (negative number)",
			ChangeInfraInternalAllocationParamMsg: NewChangeInfraInternalAllocationParamMsg("user1", p3, ""),
			expectedError:                         ErrIllegalParameter(),
		},
		{
			testName: "empty username is illegal",
			ChangeInfraInternalAllocationParamMsg: NewChangeInfraInternalAllocationParamMsg("", p1, ""),
			expectedError:                         ErrInvalidUsername(),
		},
		{
			testName: "reason is too long",
			ChangeInfraInternalAllocationParamMsg: NewChangeInfraInternalAllocationParamMsg(
				"user1", p1, string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			ChangeInfraInternalAllocationParamMsg: NewChangeInfraInternalAllocationParamMsg(
				"user1", p1, tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.ChangeInfraInternalAllocationParamMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestChangeVoteParamMsg(t *testing.T) {
	p1 := param.VoteParam{
		VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
		VoterCoinReturnTimes:           int64(7),
		DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DelegatorCoinReturnTimes:       int64(7),
	}

	p2 := p1
	p2.VoterCoinReturnIntervalSec = int64(0)

	p3 := p1
	p3.VoterCoinReturnTimes = int64(0)

	p4 := p1
	p4.DelegatorCoinReturnIntervalSec = int64(-1)

	p5 := p1
	p5.DelegatorCoinReturnTimes = int64(0)

	testCases := []struct {
		testName           string
		ChangeVoteParamMsg ChangeVoteParamMsg
		expectedError      sdk.Error
	}{
		{
			testName:           "normal case",
			ChangeVoteParamMsg: NewChangeVoteParamMsg("user1", p1, ""),
			expectedError:      nil,
		},
		{
			testName:           "zero VoterCoinReturnIntervalHr is illegal",
			ChangeVoteParamMsg: NewChangeVoteParamMsg("user1", p2, ""),
			expectedError:      ErrIllegalParameter(),
		},
		{
			testName:           "zero VoterCoinReturnTimes is illegal",
			ChangeVoteParamMsg: NewChangeVoteParamMsg("user1", p3, ""),
			expectedError:      ErrIllegalParameter(),
		},
		{
			testName:           "negative DelegatorCoinReturnIntervalHr is illegal",
			ChangeVoteParamMsg: NewChangeVoteParamMsg("user1", p4, ""),
			expectedError:      ErrIllegalParameter(),
		},
		{
			testName:           "zero DelegatorCoinReturnTimes is illegal",
			ChangeVoteParamMsg: NewChangeVoteParamMsg("user1", p5, ""),
			expectedError:      ErrIllegalParameter(),
		},
		{
			testName:           "empty username is illegal",
			ChangeVoteParamMsg: NewChangeVoteParamMsg("", p1, ""),
			expectedError:      ErrInvalidUsername(),
		},
		{
			testName: "reason is too long",
			ChangeVoteParamMsg: NewChangeVoteParamMsg(
				"user1", p1, string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			ChangeVoteParamMsg: NewChangeVoteParamMsg(
				"user1", p1, tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.ChangeVoteParamMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestChangeDeveloperParamMsg(t *testing.T) {
	p1 := param.DeveloperParam{
		DeveloperCoinReturnIntervalSec: int64(7 * 24 * 3600),
		DeveloperCoinReturnTimes:       int64(7),
		DeveloperMinDeposit:            types.NewCoinFromInt64(1 * types.Decimals),
	}

	p2 := p1
	p2.DeveloperCoinReturnTimes = int64(-7)

	p3 := p1
	p3.DeveloperCoinReturnIntervalSec = int64(0)

	p4 := p1
	p4.DeveloperMinDeposit = types.NewCoinFromInt64(-1 * types.Decimals)

	testCases := []struct {
		testName                string
		ChangeDeveloperParamMsg ChangeDeveloperParamMsg
		expectedError           sdk.Error
	}{
		{
			testName:                "normal case",
			ChangeDeveloperParamMsg: NewChangeDeveloperParamMsg("user1", p1, ""),
			expectedError:           nil,
		},
		{
			testName:                "negative DeveloperCoinReturnTimes is illegal",
			ChangeDeveloperParamMsg: NewChangeDeveloperParamMsg("user1", p2, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "zero DeveloperCoinReturnIntervalHr is illegal",
			ChangeDeveloperParamMsg: NewChangeDeveloperParamMsg("user1", p3, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "negative DeveloperMinDeposit is iilegal",
			ChangeDeveloperParamMsg: NewChangeDeveloperParamMsg("user1", p4, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "empty username is illegal",
			ChangeDeveloperParamMsg: NewChangeDeveloperParamMsg("", p1, ""),
			expectedError:           ErrInvalidUsername(),
		},
		{
			testName: "reason is too long",
			ChangeDeveloperParamMsg: NewChangeDeveloperParamMsg(
				"user1", p1, string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			ChangeDeveloperParamMsg: NewChangeDeveloperParamMsg(
				"user1", p1, tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.ChangeDeveloperParamMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestChangeValidatorParamMsg(t *testing.T) {
	p1 := param.ValidatorParam{
		ValidatorMinWithdraw:           types.NewCoinFromInt64(1 * types.Decimals),
		ValidatorMinVotingDeposit:      types.NewCoinFromInt64(3000 * types.Decimals),
		ValidatorMinCommittingDeposit:  types.NewCoinFromInt64(1000 * types.Decimals),
		ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
		ValidatorCoinReturnTimes:       int64(7),
		PenaltyMissVote:                types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
		PenaltyByzantine:               types.NewCoinFromInt64(1000 * types.Decimals),
		ValidatorListSize:              int64(21),
		AbsentCommitLimitation:         int64(100),
	}

	p2 := p1
	p2.ValidatorMinWithdraw = types.NewCoinFromInt64(-1 * types.Decimals)

	p3 := p1
	p3.ValidatorMinVotingDeposit = types.NewCoinFromInt64(0 * types.Decimals)

	p4 := p1
	p4.ValidatorMinCommittingDeposit = types.NewCoinFromInt64(-1000 * types.Decimals)

	p5 := p1
	p5.ValidatorCoinReturnIntervalSec = int64(-7 * 24 * 3600)

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

	testCases := []struct {
		testName                string
		ChangeValidatorParamMsg ChangeValidatorParamMsg
		expectedError           sdk.Error
	}{
		{
			testName:                "normal case",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p1, ""),
			expectedError:           nil,
		},
		{
			testName:                "negative ValidatorMinWithdraw is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p2, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "zero ValidatorMinVotingDeposit is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p3, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "negative ValidatorMinCommittingDeposit is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p4, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "negative ValidatorCoinReturnIntervalHr is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p5, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "zero ValidatorCoinReturnTimes is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p6, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "negative PenaltyMissVote is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p7, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "negative PenaltyByzantine is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p8, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "zero PenaltyMissCommit is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p9, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "zero AbsentCommitLimitation is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p10, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "negative ValidatorListSize is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("user1", p11, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "empty username is illegal",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg("", p1, ""),
			expectedError:           ErrInvalidUsername(),
		},
		{
			testName: "reason is too long",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg(
				"user1", p1, string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			ChangeValidatorParamMsg: NewChangeValidatorParamMsg(
				"user1", p1, tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.ChangeValidatorParamMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestChangeProposalParamMsg(t *testing.T) {
	p1 := param.ProposalParam{
		ContentCensorshipDecideSec:  int64(24 * 7 * 3600),
		ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
		ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
		ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

		ChangeParamDecideSec:    int64(24 * 7 * 3600),
		ChangeParamExecutionSec: int64(24 * 3600),
		ChangeParamPassRatio:    sdk.NewRat(70, 100),
		ChangeParamPassVotes:    types.NewCoinFromInt64(1000000 * types.Decimals),
		ChangeParamMinDeposit:   types.NewCoinFromInt64(100000 * types.Decimals),

		ProtocolUpgradeDecideSec:  int64(24 * 7 * 3600),
		ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
		ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
		ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
	}

	p2 := p1
	p2.ContentCensorshipDecideSec = int64(-24 * 7 * 3600)

	p3 := p1
	p3.ContentCensorshipPassRatio = sdk.NewRat(150, 100)

	p4 := p1
	p4.ContentCensorshipPassVotes = types.NewCoinFromInt64(-10000 * types.Decimals)

	p5 := p1
	p5.ContentCensorshipMinDeposit = types.NewCoinFromInt64(-100 * types.Decimals)

	p6 := p1
	p6.ChangeParamDecideSec = int64(-24 * 7 * 3600)

	p7 := p1
	p7.ChangeParamPassRatio = sdk.NewRat(0, 8)

	p8 := p1
	p8.ChangeParamPassVotes = types.NewCoinFromInt64(0 * types.Decimals)

	p9 := p1
	p9.ChangeParamMinDeposit = types.NewCoinFromInt64(-100000 * types.Decimals)

	p10 := p1
	p10.ProtocolUpgradeDecideSec = int64(0)

	p11 := p1
	p11.ProtocolUpgradePassRatio = sdk.NewRat(0, 100)

	p12 := p1
	p12.ProtocolUpgradePassVotes = types.NewCoinFromInt64(-10000000 * types.Decimals)

	p13 := p1
	p13.ProtocolUpgradeMinDeposit = types.NewCoinFromInt64(-1000000 * types.Decimals)

	testCases := []struct {
		testName               string
		ChangeProposalParamMsg ChangeProposalParamMsg
		expectedError          sdk.Error
	}{
		{
			testName:               "normal case",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p1, ""),
			expectedError:          nil,
		},
		{
			testName:               "invalid username",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("", p1, ""),
			expectedError:          ErrInvalidUsername(),
		},
		{
			testName:               "negative ContentCensorshipDecideHr is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p2, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "ContentCensorshipPassRatio that is larger than one is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p3, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "negative ContentCensorshipPassVotes is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p4, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "negative ContentCensorshipMinDeposit is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p5, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "negative ChangeParamDecideHr is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p6, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "zero ChangeParamPassRatio is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p7, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "zero ChangeParamPassVotes is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p8, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "negative ChangeParamMinDeposit is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p9, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "zero ProtocolUpgradeDecideHr is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p10, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "zero ProtocolUpgradePassRatio is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p11, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "negative ProtocolUpgradePassVotes is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p12, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName:               "negative ProtocolUpgradeMinDeposit is illegal",
			ChangeProposalParamMsg: NewChangeProposalParamMsg("user1", p13, ""),
			expectedError:          ErrIllegalParameter(),
		},
		{
			testName: "reason is too long",
			ChangeProposalParamMsg: NewChangeProposalParamMsg(
				"user1", p1, string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			ChangeProposalParamMsg: NewChangeProposalParamMsg(
				"user1", p1, tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.ChangeProposalParamMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestChangeAccountParamMsg(t *testing.T) {
	p1 := param.AccountParam{
		MinimumBalance:               types.NewCoinFromInt64(1 * types.Decimals),
		RegisterFee:                  types.NewCoinFromInt64(1 * types.Decimals),
		FirstDepositFullCoinDayLimit: types.NewCoinFromInt64(1 * types.Decimals),
		MaxNumFrozenMoney:            10,
	}

	p2 := p1
	p2.MinimumBalance = types.NewCoinFromInt64(0)

	p3 := p1
	p3.RegisterFee = types.NewCoinFromInt64(0)

	p4 := p1
	p4.RegisterFee = types.NewCoinFromInt64(-1)

	p5 := p1
	p5.FirstDepositFullCoinDayLimit = types.NewCoinFromInt64(-1)

	p6 := p1
	p6.MaxNumFrozenMoney = -1

	testCases := []struct {
		testName              string
		changeAccountParamMsg ChangeAccountParamMsg
		expectedError         sdk.Error
	}{
		{
			testName:              "normal case",
			changeAccountParamMsg: NewChangeAccountParamMsg("user1", p1, ""),
			expectedError:         nil,
		},
		{
			testName:              "too short username is invalid",
			changeAccountParamMsg: NewChangeAccountParamMsg("us", p1, ""),
			expectedError:         ErrInvalidUsername(),
		},
		{
			testName:              "too long username is invalid",
			changeAccountParamMsg: NewChangeAccountParamMsg("user1user1user1user1user1user1", p1, ""),
			expectedError:         ErrInvalidUsername(),
		},
		{
			testName:              "zero MinimumBalance is valid",
			changeAccountParamMsg: NewChangeAccountParamMsg("user1", p2, ""),
			expectedError:         nil,
		},
		{
			testName:              "zero RegisterFee is valid",
			changeAccountParamMsg: NewChangeAccountParamMsg("user1", p3, ""),
			expectedError:         nil,
		},
		{
			testName:              "negative RegisterFee is invalid",
			changeAccountParamMsg: NewChangeAccountParamMsg("user1", p4, ""),
			expectedError:         ErrIllegalParameter(),
		},
		{
			testName:              "negative FirstDepositFullCoinDayLimit is invalid",
			changeAccountParamMsg: NewChangeAccountParamMsg("user1", p5, ""),
			expectedError:         ErrIllegalParameter(),
		},
		{
			testName:              "negative MaxNumFrozenMoney is invalid",
			changeAccountParamMsg: NewChangeAccountParamMsg("user1", p6, ""),
			expectedError:         ErrIllegalParameter(),
		},
		{
			testName: "reason is too long",
			changeAccountParamMsg: NewChangeAccountParamMsg(
				"user1", p1, string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			changeAccountParamMsg: NewChangeAccountParamMsg(
				"user1", p1, tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.changeAccountParamMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestChangeBandwidthParamMsg(t *testing.T) {
	p1 := param.BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
		VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
	}

	p2 := p1
	p2.SecondsToRecoverBandwidth = int64(-1)

	p3 := p1
	p3.CapacityUsagePerTransaction = types.NewCoinFromInt64(-1)

	p4 := p1
	p4.VirtualCoin = types.NewCoinFromInt64(-1)

	testCases := []struct {
		testName                string
		changeBandwidthParamMsg ChangeBandwidthParamMsg
		expectedError           sdk.Error
	}{
		{
			testName:                "normal case",
			changeBandwidthParamMsg: NewChangeBandwidthParamMsg("user1", p1, ""),
			expectedError:           nil,
		},
		{
			testName:                "too short username is illegal",
			changeBandwidthParamMsg: NewChangeBandwidthParamMsg("us", p1, ""),
			expectedError:           ErrInvalidUsername(),
		},
		{
			testName:                "too long username is illegal",
			changeBandwidthParamMsg: NewChangeBandwidthParamMsg("user1user1user1user1user1user1", p1, ""),
			expectedError:           ErrInvalidUsername(),
		},
		{
			testName:                "negative SecondsToRecoverBandwidth is illegal",
			changeBandwidthParamMsg: NewChangeBandwidthParamMsg("user1", p2, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "negative CapacityUsagePerTransaction is illegal",
			changeBandwidthParamMsg: NewChangeBandwidthParamMsg("user1", p3, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName:                "negative VirtualCoin is illegal",
			changeBandwidthParamMsg: NewChangeBandwidthParamMsg("user1", p4, ""),
			expectedError:           ErrIllegalParameter(),
		},
		{
			testName: "reason is too long",
			changeBandwidthParamMsg: NewChangeBandwidthParamMsg(
				"user1", p1, string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			changeBandwidthParamMsg: NewChangeBandwidthParamMsg(
				"user1", p1, tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.changeBandwidthParamMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestChangePostParamMsg(t *testing.T) {
	p1 := param.PostParam{
		ReportOrUpvoteIntervalSec: 1,
		PostIntervalSec:           1,
	}

	p2 := p1
	p2.ReportOrUpvoteIntervalSec = int64(-1)

	p3 := p1
	p3.PostIntervalSec = int64(-1)

	testCases := []struct {
		testName           string
		changePostParamMsg ChangePostParamMsg
		expectedError      sdk.Error
	}{
		{
			testName:           "normal case",
			changePostParamMsg: NewChangePostParamMsg("user1", p1, ""),
			expectedError:      nil,
		},
		{
			testName:           "illegal report or upvote interval",
			changePostParamMsg: NewChangePostParamMsg("user1", p2, ""),
			expectedError:      ErrIllegalParameter(),
		},
		{
			testName:           "illegal post interval",
			changePostParamMsg: NewChangePostParamMsg("user1", p3, ""),
			expectedError:      ErrIllegalParameter(),
		},
		{
			testName:           "username too short",
			changePostParamMsg: NewChangePostParamMsg("us", p1, ""),
			expectedError:      ErrInvalidUsername(),
		},
		{
			testName:           "username too long",
			changePostParamMsg: NewChangePostParamMsg("user1user1user1user1user1", p1, ""),
			expectedError:      ErrInvalidUsername(),
		},
	}

	for _, tc := range testCases {
		result := tc.changePostParamMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
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

	testCases := []struct {
		testName              string
		changeAccountParamMsg ChangeEvaluateOfContentValueParamMsg
		expectedError         sdk.Error
	}{
		{
			testName:              "normal case",
			changeAccountParamMsg: NewChangeEvaluateOfContentValueParamMsg("user1", p1, ""),
			expectedError:         nil,
		},
		{
			testName:              "zero ConsumptionTimeAdjustBase is illegal",
			changeAccountParamMsg: NewChangeEvaluateOfContentValueParamMsg("user1", p2, ""),
			expectedError:         ErrIllegalParameter(),
		},
		{
			testName:              "zero TotalAmountOfConsumptionBase is illegal",
			changeAccountParamMsg: NewChangeEvaluateOfContentValueParamMsg("user1", p3, ""),
			expectedError:         ErrIllegalParameter(),
		},
		{
			testName:              "too short username is illegal",
			changeAccountParamMsg: NewChangeEvaluateOfContentValueParamMsg("us", p1, ""),
			expectedError:         ErrInvalidUsername(),
		},
		{
			testName:              "too long username is illegal",
			changeAccountParamMsg: NewChangeEvaluateOfContentValueParamMsg("user1user1user1user1user1", p1, ""),
			expectedError:         ErrInvalidUsername(),
		},
		{
			testName: "reason is too long",
			changeAccountParamMsg: NewChangeEvaluateOfContentValueParamMsg(
				"user1", p1, string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			changeAccountParamMsg: NewChangeEvaluateOfContentValueParamMsg(
				"user1", p1, tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.changeAccountParamMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestDeletePostContentMsg(t *testing.T) {
	testCases := []struct {
		testName             string
		deletePostContentMsg DeletePostContentMsg
		expectedError        sdk.Error
	}{
		{
			testName:             "normal case",
			deletePostContentMsg: NewDeletePostContentMsg("user1", "permlink", "reason"),
			expectedError:        nil,
		},
		{
			testName:             "too short username is illegal",
			deletePostContentMsg: NewDeletePostContentMsg("us", "permlink", "reason"),
			expectedError:        ErrInvalidUsername(),
		},
		{
			testName:             "too long username is illegal",
			deletePostContentMsg: NewDeletePostContentMsg("user1user1user1user1user1user1", "permlink", "reason"),
			expectedError:        ErrInvalidUsername(),
		},
		{
			testName:             "empty permlink is illegal",
			deletePostContentMsg: NewDeletePostContentMsg("user1", "", "reason"),
			expectedError:        ErrInvalidPermlink(),
		},
		{
			testName: "reason is too long",
			deletePostContentMsg: NewDeletePostContentMsg(
				"user1", "permlink", string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError: ErrReasonTooLong(),
		},
		{
			testName: "utf8 reason is too long",
			deletePostContentMsg: NewDeletePostContentMsg(
				"user1", "permlink", tooLongOfUTF8Reason),
			expectedError: ErrReasonTooLong(),
		},
	}

	for _, tc := range testCases {
		result := tc.deletePostContentMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestUpgradeProtocolMsg(t *testing.T) {
	testCases := []struct {
		testName           string
		upgradeProtocolMsg UpgradeProtocolMsg
		expectedError      sdk.Error
	}{
		{
			testName:           "normal case",
			upgradeProtocolMsg: NewUpgradeProtocolMsg("user1", "link", ""),
			expectedError:      nil,
		},
		{
			testName:           "too short username is illegal",
			upgradeProtocolMsg: NewUpgradeProtocolMsg("us", "link", ""),
			expectedError:      ErrInvalidUsername(),
		},
		{
			testName:           "too long username is illegal",
			upgradeProtocolMsg: NewUpgradeProtocolMsg("user1user1user1user1user1user1", "link", ""),
			expectedError:      ErrInvalidUsername(),
		},
		{
			testName:           "empty link is illegal",
			upgradeProtocolMsg: NewUpgradeProtocolMsg("user1", "", ""),
			expectedError:      ErrInvalidLink(),
		},
		{
			testName:           "reason is too long",
			upgradeProtocolMsg: NewUpgradeProtocolMsg("user1", "", string(make([]byte, types.MaximumLengthOfProposalReason+1))),
			expectedError:      ErrInvalidLink(),
		},
		{
			testName:           "utf8 reason is too long",
			upgradeProtocolMsg: NewUpgradeProtocolMsg("user1", "", tooLongOfUTF8Reason),
			expectedError:      ErrInvalidLink(),
		},
	}

	for _, tc := range testCases {
		result := tc.upgradeProtocolMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestMsgPermission(t *testing.T) {
	testCases := []struct {
		testName         string
		msg              types.Msg
		expectPermission types.Permission
	}{
		{
			testName: "delete post content msg",
			msg: NewDeletePostContentMsg(
				"creator", "perm_link", "reason"),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "upgrade protocol msg",
			msg:              NewUpgradeProtocolMsg("creator", "link", ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change global allocaiton param msg",
			msg: NewChangeGlobalAllocationParamMsg(
				"creator", param.GlobalAllocationParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change evaluate of content value param msg",
			msg: NewChangeEvaluateOfContentValueParamMsg(
				"creator", param.EvaluateOfContentValueParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change infra internal allocation param msg",
			msg: NewChangeInfraInternalAllocationParamMsg(
				"creator", param.InfraInternalAllocationParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change vote param msg",
			msg: NewChangeInfraInternalAllocationParamMsg(
				"creator", param.InfraInternalAllocationParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change proposal param msg",
			msg: NewChangeProposalParamMsg(
				"creator", param.ProposalParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change developer param msg",
			msg: NewChangeDeveloperParamMsg(
				"creator", param.DeveloperParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change validator param msg",
			msg: NewChangeValidatorParamMsg(
				"creator", param.ValidatorParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change bandwidth param msg",
			msg: NewChangeBandwidthParamMsg(
				"creator", param.BandwidthParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change account param msg",
			msg: NewChangeAccountParamMsg(
				"creator", param.AccountParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName: "change post param msg",
			msg: NewChangePostParamMsg(
				"creator", param.PostParam{}, ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "vote proposal msg",
			msg:              NewVoteProposalMsg("voter", 1, true),
			expectPermission: types.TransactionPermission,
		},
	}

	for _, tc := range testCases {
		permission := tc.msg.GetPermission()
		if tc.expectPermission != permission {
			t.Errorf("%s: diff permission, got %v, want %v", tc.testName, permission, tc.expectPermission)
			return
		}
	}
}

func TestGetSignBytes(t *testing.T) {
	testCases := []struct {
		testName string
		msg      types.Msg
	}{
		{
			testName: "delete post content msg",
			msg: NewDeletePostContentMsg(
				"creator", "perm_link", "reason"),
		},
		{
			testName: "upgrade protocol msg",
			msg:      NewUpgradeProtocolMsg("creator", "link", ""),
		},
		{
			testName: "change global allocaiton param msg",
			msg: NewChangeGlobalAllocationParamMsg(
				"creator", param.GlobalAllocationParam{}, ""),
		},
		{
			testName: "change evaluate of content value param msg",
			msg: NewChangeEvaluateOfContentValueParamMsg(
				"creator", param.EvaluateOfContentValueParam{}, ""),
		},
		{
			testName: "change infra internal allocation param msg",
			msg: NewChangeInfraInternalAllocationParamMsg(
				"creator", param.InfraInternalAllocationParam{}, ""),
		},
		{
			testName: "change vote param msg",
			msg: NewChangeInfraInternalAllocationParamMsg(
				"creator", param.InfraInternalAllocationParam{}, ""),
		},
		{
			testName: "change proposal param msg",
			msg: NewChangeProposalParamMsg(
				"creator", param.ProposalParam{}, ""),
		},
		{
			testName: "change developer param msg",
			msg: NewChangeDeveloperParamMsg(
				"creator", param.DeveloperParam{}, ""),
		},
		{
			testName: "change validator param msg",
			msg: NewChangeValidatorParamMsg(
				"creator", param.ValidatorParam{}, ""),
		},
		{
			testName: "change bandwidth param msg",
			msg: NewChangeBandwidthParamMsg(
				"creator", param.BandwidthParam{}, ""),
		},
		{
			testName: "change account param msg",
			msg: NewChangeAccountParamMsg(
				"creator", param.AccountParam{}, ""),
		},
		{
			testName: "change post param msg",
			msg: NewChangePostParamMsg(
				"creator", param.PostParam{}, ""),
		},
		{
			testName: "vote proposal msg",
			msg:      NewVoteProposalMsg("voter", 1, true),
		},
	}

	for _, tc := range testCases {
		require.NotPanics(t, func() { tc.msg.GetSignBytes() }, tc.testName)
	}
}

func TestGetSigners(t *testing.T) {
	testCases := []struct {
		testName      string
		msg           types.Msg
		expectSigners []types.AccountKey
	}{
		{
			testName: "delete post content msg",
			msg: NewDeletePostContentMsg(
				"creator", "perm_link", "reason"),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName:      "upgrade protocol msg",
			msg:           NewUpgradeProtocolMsg("creator", "link", ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change global allocaiton param msg",
			msg: NewChangeGlobalAllocationParamMsg(
				"creator", param.GlobalAllocationParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change evaluate of content value param msg",
			msg: NewChangeEvaluateOfContentValueParamMsg(
				"creator", param.EvaluateOfContentValueParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change infra internal allocation param msg",
			msg: NewChangeInfraInternalAllocationParamMsg(
				"creator", param.InfraInternalAllocationParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change vote param msg",
			msg: NewChangeInfraInternalAllocationParamMsg(
				"creator", param.InfraInternalAllocationParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change proposal param msg",
			msg: NewChangeProposalParamMsg(
				"creator", param.ProposalParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change developer param msg",
			msg: NewChangeDeveloperParamMsg(
				"creator", param.DeveloperParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change validator param msg",
			msg: NewChangeValidatorParamMsg(
				"creator", param.ValidatorParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change bandwidth param msg",
			msg: NewChangeBandwidthParamMsg(
				"creator", param.BandwidthParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change account param msg",
			msg: NewChangeAccountParamMsg(
				"creator", param.AccountParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName: "change post param msg",
			msg: NewChangePostParamMsg(
				"creator", param.PostParam{}, ""),
			expectSigners: []types.AccountKey{"creator"},
		},
		{
			testName:      "vote proposal msg",
			msg:           NewVoteProposalMsg("voter", 1, true),
			expectSigners: []types.AccountKey{"voter"},
		},
	}

	for _, tc := range testCases {
		if len(tc.msg.GetSigners()) != len(tc.expectSigners) {
			t.Errorf("%s: expect number of signers wrong, got %v, want %v", tc.testName, len(tc.msg.GetSigners()), len(tc.expectSigners))
			return
		}
		for i, signer := range tc.msg.GetSigners() {
			if types.AccountKey(signer) != tc.expectSigners[i] {
				t.Errorf("%s: expect signer wrong, got %v, want %v", tc.testName, types.AccountKey(signer), tc.expectSigners[i])
				return
			}
		}
	}
}
