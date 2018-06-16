package proposal

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrUsernameNotFound() sdk.Error {
	return sdk.NewError(types.CodeProposalHandlerError, fmt.Sprintf("Username not found"))
}

func ErrPostNotFound() sdk.Error {
	return sdk.NewError(types.CodeProposalHandlerError, fmt.Sprintf("Username not found"))
}

func ErrOngoingProposalNotFound() sdk.Error {
	return sdk.NewError(types.CodeProposalManagerError, fmt.Sprintf("Ongoing proposal not found"))
}

func ErrInvalidUsername() sdk.Error {
	return sdk.NewError(types.CodeProposalManagerError, fmt.Sprintf("Invalid Username"))
}

func ErrInvalidPermLink() sdk.Error {
	return sdk.NewError(types.CodeProposalMsgError, fmt.Sprintf("Invalid PermLink"))
}

func ErrInvalidLink() sdk.Error {
	return sdk.NewError(types.CodeProposalMsgError, fmt.Sprintf("Invalid Link"))
}

func ErrCensorshipPostNotFound() sdk.Error {
	return sdk.NewError(types.CodeProposalEventError, fmt.Sprintf("Censorship post not found"))
}

func ErrCensorshipPostIsDeleted(permLink types.PermLink) sdk.Error {
	return sdk.NewError(types.CodeProposalEventError, fmt.Sprintf("Censorship post %v is deleted", permLink))
}

func ErrIllegalParameter() sdk.Error {
	return sdk.NewError(types.CodeProposalManagerError, fmt.Sprintf("Invalid parameter"))
}

func ErrProposalInfoNotFound() sdk.Error {
	return sdk.NewError(types.CodeProposalManagerError, fmt.Sprintf("Proposal info not found"))
}

func ErrWrongProposalType() sdk.Error {
	return sdk.NewError(types.CodeProposalManagerError, fmt.Sprintf("Wrong proposal type"))
}

func ErrGetVoter() sdk.Error {
	return sdk.NewError(types.CodeProposalHandlerError, fmt.Sprintf("Get voter failed"))
}

func ErrNotOngoingProposal() sdk.Error {
	return sdk.NewError(types.CodeProposalHandlerError, fmt.Sprintf("Not ongoing proposal"))
}
