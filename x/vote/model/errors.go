package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// not found error
func ErrVoterNotFound() sdk.Error {
	return types.NewError(types.CodeVoterNotFound, fmt.Sprintf("voter is not found"))
}

func ErrVoteNotFound() sdk.Error {
	return types.NewError(types.CodeVoteNotFound, fmt.Sprintf("vote is not found"))
}

func ErrReferenceListNotFound() sdk.Error {
	return types.NewError(types.CodeReferenceListNotFound, fmt.Sprintf("reference list is not found"))
}

func ErrDelegationNotFound() sdk.Error {
	return types.NewError(types.CodeDelegationNotFound, fmt.Sprintf("delegation is not found"))
}

// marshal error
func ErrFailedToMarshalVoter(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalVoter, fmt.Sprintf("failed to marshal voter: %s", err.Error()))
}

func ErrFailedToMarshalVote(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalVote, fmt.Sprintf("failed to marshal vote: %s", err.Error()))
}

func ErrFailedToMarshalDelegation(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalDelegation, fmt.Sprintf("failed to marshal delegation: %s", err.Error()))
}

func ErrFailedToMarshalReferenceList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalReferenceList, fmt.Sprintf("failed to marshal reference list: %s", err.Error()))
}

// unmarshal error
func ErrFailedToUnmarshalVoter(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalVoter, fmt.Sprintf("failed to unmarshal voter: %s", err.Error()))
}

func ErrFailedToUnmarshalVote(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalVote, fmt.Sprintf("failed to unmarshal vote: %s", err.Error()))
}

func ErrFailedToUnmarshalDelegation(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalDelegation, fmt.Sprintf("failed to unmarshal delegation: %s", err.Error()))
}

func ErrFailedToUnmarshalReferenceList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalReferenceList, fmt.Sprintf("failed to unmarshal reference list: %s", err.Error()))
}
