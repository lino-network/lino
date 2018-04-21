package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrGetVoter() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Get voter failed"))
}

func ErrGetVote() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Get vote failed"))
}

func ErrGetProposal() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Get proposal failed"))
}

func ErrGetPenaltyList() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Get penalty failed"))
}

func ErrGetDelegation() sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Get delegation failed"))
}

func ErrVoterMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Voter marshal error: %s", err.Error()))
}

func ErrVoterUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Voter unmarshal error: %s", err.Error()))
}

func ErrVoteMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Vote marshal error: %s", err.Error()))
}

func ErrVoteUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Vote unmarshal error: %s", err.Error()))
}

func ErrProposalMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Proposal marshal error: %s", err.Error()))
}

func ErrProposalUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Proposal unmarshal error: %s", err.Error()))
}

func ErrPenaltyListMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Penalty list marshal error: %s", err.Error()))
}

func ErrPenaltyListUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Penalty list unmarshal error: %s", err.Error()))
}

func ErrDelegationMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Delegation marshal error: %s", err.Error()))
}

func ErrDelegationUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteManagerFailed, fmt.Sprintf("Delegation unmarshal error: %s", err.Error()))
}
