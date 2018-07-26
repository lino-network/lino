package proposal

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrAccountNotFound() sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("username is not found"))
}

func ErrPostNotFound() sdk.Error {
	return types.NewError(types.CodePostNotFound, fmt.Sprintf("post is not found"))
}

func ErrCensorshipPostIsDeleted(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeCensorshipPostIsDeleted, fmt.Sprintf("censorship post %v is deleted", permlink))
}

func ErrVoterNotFound() sdk.Error {
	return types.NewError(types.CodeVoterNotFound, fmt.Sprintf("voter is not found"))
}

func ErrNotOngoingProposal() sdk.Error {
	return types.NewError(types.CodeNotOngoingProposal, fmt.Sprintf("not ongoing proposal"))
}

func ErrIncorrectProposalType() sdk.Error {
	return types.NewError(types.CodeIncorrectProposalType, fmt.Sprintf("proposal type is wrong"))
}

func ErrOngoingProposalNotFound() sdk.Error {
	return types.NewError(types.CodeOngoingProposalNotFound, fmt.Sprintf("ongoing proposal not found"))
}

func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("invalid username"))
}

func ErrInvalidPermlink() sdk.Error {
	return types.NewError(types.CodeInvalidPermlink, fmt.Sprintf("invalid permlink"))
}

func ErrReasonTooLong() sdk.Error {
	return types.NewError(types.CodeReasonTooLong, fmt.Sprintf("reason length is too long"))
}

func ErrInvalidLink() sdk.Error {
	return types.NewError(types.CodeInvalidLink, fmt.Sprintf("invalid Link"))
}

func ErrCensorshipPostNotFound() sdk.Error {
	return types.NewError(types.CodeCensorshipPostNotFound, fmt.Sprintf("Censorship post not found"))
}

func ErrIllegalParameter() sdk.Error {
	return types.NewError(types.CodeIllegalParameter, fmt.Sprintf("invalid parameter"))
}
