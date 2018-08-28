package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrVoterNotFound - error if voter is not found in KVStore
func ErrVoterNotFound() sdk.Error {
	return types.NewError(types.CodeVoterNotFound, fmt.Sprintf("voter is not found"))
}

// ErrVoteNotFound - error if vote is not found in KVStore
func ErrVoteNotFound() sdk.Error {
	return types.NewError(types.CodeVoteNotFound, fmt.Sprintf("vote is not found"))
}

// ErrReferenceListNotFound - error if reference list is not found in KVStore
func ErrReferenceListNotFound() sdk.Error {
	return types.NewError(types.CodeReferenceListNotFound, fmt.Sprintf("reference list is not found"))
}

// ErrDelegationNotFound - error if delegation is not found in KVStore
func ErrDelegationNotFound() sdk.Error {
	return types.NewError(types.CodeDelegationNotFound, fmt.Sprintf("delegation is not found"))
}

// ErrFailedToMarshalVoter - error if marshal voter failed
func ErrFailedToMarshalVoter(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalVoter, fmt.Sprintf("failed to marshal voter: %s", err.Error()))
}

// ErrFailedToMarshalVote - error if marshal vote failed
func ErrFailedToMarshalVote(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalVote, fmt.Sprintf("failed to marshal vote: %s", err.Error()))
}

// ErrFailedToMarshalDelegation - error if marshal delegation failed
func ErrFailedToMarshalDelegation(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalDelegation, fmt.Sprintf("failed to marshal delegation: %s", err.Error()))
}

// ErrFailedToMarshalReferenceList - error if marshal reference list failed
func ErrFailedToMarshalReferenceList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalReferenceList, fmt.Sprintf("failed to marshal reference list: %s", err.Error()))
}

// ErrFailedToUnmarshalVoter - error if unmarshal voter failed
func ErrFailedToUnmarshalVoter(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalVoter, fmt.Sprintf("failed to unmarshal voter: %s", err.Error()))
}

// ErrFailedToUnmarshalVote - error if unmarshal vote failed
func ErrFailedToUnmarshalVote(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalVote, fmt.Sprintf("failed to unmarshal vote: %s", err.Error()))
}

// ErrFailedToUnmarshalDelegation - error if unmarshal delegation failed
func ErrFailedToUnmarshalDelegation(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalDelegation, fmt.Sprintf("failed to unmarshal delegation: %s", err.Error()))
}

// ErrFailedToUnmarshalReferenceList - error if unmarshal reference list failed
func ErrFailedToUnmarshalReferenceList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalReferenceList, fmt.Sprintf("failed to unmarshal reference list: %s", err.Error()))
}
