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

// ErrFailedToMarshalVoter - error if marshal voter failed
func ErrFailedToMarshalVoter(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalVoter, fmt.Sprintf("failed to marshal voter: %s", err.Error()))
}

// ErrFailedToUnmarshalVoter - error if unmarshal voter failed
func ErrFailedToUnmarshalVoter(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalVoter, fmt.Sprintf("failed to unmarshal voter: %s", err.Error()))
}
