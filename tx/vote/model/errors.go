package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case types.CodeInvalidUsername:
		return "Invalid username format"
	case types.CodeAccRegisterFailed:
		return "Voter register failed"
	case types.CodeVoteHandlerFailed:
		return "Vote handler failed"
	case types.CodeVoteManagerFailed:
		return "Vote manager failed"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// Error constructors
func ErrSetVoter() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Set voter failed"))
}

func ErrGetVoter() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get voter failed"))
}

func ErrSetVote() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Set vote failed"))
}

func ErrGetVote() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get vote failed"))
}

func ErrSetProposal() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Set proposal failed"))
}

func ErrGetProposal() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get proposal failed"))
}

func ErrGetDelegation() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get delegation failed"))
}

func ErrSetDelegation() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Set delegation failed"))
}

func ErrVoterMarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Voter marshal error: %s", err.Error()))
}

func ErrVoterUnmarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Voter unmarshal error: %s", err.Error()))
}

func ErrVoteMarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Vote marshal error: %s", err.Error()))
}

func ErrVoteUnmarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Vote unmarshal error: %s", err.Error()))
}

func ErrProposalMarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Proposal marshal error: %s", err.Error()))
}

func ErrProposalUnmarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Proposal unmarshal error: %s", err.Error()))
}

func ErrDelegationMarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Delegation marshal error: %s", err.Error()))
}

func ErrDelegationUnmarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Delegation unmarshal error: %s", err.Error()))
}

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	} else {
		return codeToDefaultMsg(code)
	}
}

func newError(code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}
