package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrProposalNotFound - error if proposal is not found in KVStore
func ErrProposalNotFound() sdk.Error {
	return types.NewError(types.CodeProposalNotFound, fmt.Sprintf("proposal is not found"))
}

// ErrNextProposalIDNotFound - error if next proposal ID is not found in KVStore
func ErrNextProposalIDNotFound() sdk.Error {
	return types.NewError(types.CodeNextProposalIDNotFound, fmt.Sprintf("next proposal id is not found"))
}

// ErrFailedToMarshalProposal - error if marshal proposal failed
func ErrFailedToMarshalProposal(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalProposal, fmt.Sprintf("failed to marshal proposal: %s", err.Error()))
}

// ErrFailedToMarshalNextProposalID - error if marshal next proposal id failed
func ErrFailedToMarshalNextProposalID(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalNextProposalID, fmt.Sprintf("failed to marshal next proposal id: %s", err.Error()))
}

// ErrFailedToUnmarshalProposal - error if unmarshal proposal failed
func ErrFailedToUnmarshalProposal(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalProposal, fmt.Sprintf("failed to unmarshal proposal: %s", err.Error()))
}

// ErrFailedToUnmarshalProposal - error if unmarshal next proposal id failed
func ErrFailedToUnmarshalNextProposalID(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalNextProposalID, fmt.Sprintf("failed to unmarshal next proposal id: %s", err.Error()))
}
