package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// not found err
func ErrProposalNotFound() sdk.Error {
	return types.NewError(types.CodeProposalNotFound, fmt.Sprintf("proposal is not found"))
}

func ErrProposalListNotFound() sdk.Error {
	return types.NewError(types.CodeProposalListNotFound, fmt.Sprintf("proposal list is not found"))
}

func ErrNextProposalIDNotFound() sdk.Error {
	return types.NewError(types.CodeNextProposalIDNotFound, fmt.Sprintf("next proposal id is not found"))
}

// marshal error
func ErrFailedToMarshalProposal(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalProposal, fmt.Sprintf("failed to marshal proposal: %s", err.Error()))
}

func ErrFailedToMarshalProposalList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalProposalList, fmt.Sprintf("failed to marshal proposal list: %s", err.Error()))
}

func ErrFailedToMarshalNextProposalID(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalNextProposalID, fmt.Sprintf("failed to marshal next proposal id: %s", err.Error()))
}

// unmarshal error
func ErrFailedToUnmarshalProposal(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalProposal, fmt.Sprintf("failed to unmarshal proposal: %s", err.Error()))
}

func ErrFailedToUnmarshalProposalList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalProposalList, fmt.Sprintf("failed to unmarshal proposal list: %s", err.Error()))
}

func ErrFailedToUnmarshalNextProposalID(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalNextProposalID, fmt.Sprintf("failed to unmarshal next proposal id: %s", err.Error()))
}
