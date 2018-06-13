package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrGetVoter() sdk.Error {
	return sdk.NewError(types.CodeVoteStorageFailed, fmt.Sprintf("Get voter failed"))
}

func ErrGetVote() sdk.Error {
	return sdk.NewError(types.CodeVoteStorageFailed, fmt.Sprintf("Get vote failed"))
}

func ErrGetReferenceList() sdk.Error {
	return sdk.NewError(types.CodeVoteStorageFailed, fmt.Sprintf("Get reference list failed"))
}

func ErrGetDelegateeList() sdk.Error {
	return sdk.NewError(types.CodeVoteStorageFailed, fmt.Sprintf("Get delegatee list failed"))
}

func ErrGetDelegation() sdk.Error {
	return sdk.NewError(types.CodeVoteStorageFailed, fmt.Sprintf("Get delegation failed"))
}

func ErrMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteStorageFailed, fmt.Sprintf("Vote storage marshal error: %s", err.Error()))
}

func ErrUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeVoteStorageFailed, fmt.Sprintf("Vote storage unmarshal error: %s", err.Error()))
}
