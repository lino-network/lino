package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrProposalMarshalError(err error) sdk.Error {
	return types.NewError(types.CodeProposalStoreError, fmt.Sprintf("Proposal marshal error: %s", err.Error()))
}

func ErrProposalUnmarshalError(err error) sdk.Error {
	return types.NewError(types.CodeProposalStoreError, fmt.Sprintf("Proposal unmarshal error: %s", err.Error()))
}

func ErrGetProposal() sdk.Error {
	return types.NewError(types.CodeProposalStoreError, fmt.Sprintf("Get proposal failed"))
}
