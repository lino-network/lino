package proposal

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrUsernameNotFound() sdk.Error {
	return types.NewError(types.CodeProposalHandlerError, fmt.Sprintf("Username not found"))
}

func ErrPostNotFound() sdk.Error {
	return types.NewError(types.CodeProposalHandlerError, fmt.Sprintf("Username not found"))
}

func ErrOngoingProposalNotFound() sdk.Error {
	return types.NewError(types.CodeProposalManagerError, fmt.Sprintf("Ongoing proposal not found"))
}

func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeProposalManagerError, fmt.Sprintf("Invalid Username"))
}

func ErrInvalidPermlink() sdk.Error {
	return types.NewError(types.CodeProposalMsgError, fmt.Sprintf("Invalid Permlink"))
}

func ErrInvalidLink() sdk.Error {
	return types.NewError(types.CodeProposalMsgError, fmt.Sprintf("Invalid Link"))
}

func ErrCensorshipPostNotFound() sdk.Error {
	return types.NewError(types.CodeProposalEventError, fmt.Sprintf("Censorship post not found"))
}

func ErrCensorshipPostIsDeleted(permLink types.Permlink) sdk.Error {
	return types.NewError(types.CodeProposalEventError, fmt.Sprintf("Censorship post %v is deleted", permLink))
}

func ErrIllegalParameter() sdk.Error {
	return types.NewError(types.CodeProposalManagerError, fmt.Sprintf("Invalid parameter"))
}

func ErrProposalInfoNotFound() sdk.Error {
	return types.NewError(types.CodeProposalManagerError, fmt.Sprintf("Proposal info not found"))
}

func ErrWrongProposalType() sdk.Error {
	return types.NewError(types.CodeProposalManagerError, fmt.Sprintf("Wrong proposal type"))
}

func ErrGetVoter() sdk.Error {
	return types.NewError(types.CodeProposalHandlerError, fmt.Sprintf("Get voter failed"))
}

func ErrNotOngoingProposal() sdk.Error {
	return types.NewError(types.CodeProposalHandlerError, fmt.Sprintf("Not ongoing proposal"))
}
