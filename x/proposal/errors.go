package proposal

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrAccountNotFound - error when account is not found
func ErrAccountNotFound() sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("username is not found"))
}

// ErrPostNotFound - error when post is not found
func ErrPostNotFound() sdk.Error {
	return types.NewError(types.CodePostNotFound, fmt.Sprintf("post is not found"))
}

// ErrCensorshipPostIsDeleted - error when censorship post is already deleted
func ErrCensorshipPostIsDeleted(permlink types.Permlink) sdk.Error {
	return types.NewError(types.CodeCensorshipPostIsDeleted, fmt.Sprintf("censorship post %v is deleted", permlink))
}

// ErrVoterNotFound - error when voter is not found
func ErrVoterNotFound() sdk.Error {
	return types.NewError(types.CodeVoterNotFound, fmt.Sprintf("voter is not found"))
}

// ErrNotOngoingProposal - error if vote to an expired proposal
func ErrNotOngoingProposal() sdk.Error {
	return types.NewError(types.CodeNotOngoingProposal, fmt.Sprintf("not ongoing proposal"))
}

// ErrIncorrectProposalType - error if check proposal type failed
func ErrIncorrectProposalType() sdk.Error {
	return types.NewError(types.CodeIncorrectProposalType, fmt.Sprintf("proposal type is wrong"))
}

// ErrOngoingProposalNotFound - error if ongoing proposal is not found
func ErrOngoingProposalNotFound() sdk.Error {
	return types.NewError(types.CodeOngoingProposalNotFound, fmt.Sprintf("ongoing proposal not found"))
}

// ErrInvalidUsername - error if username is invalid
func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("invalid username"))
}

// ErrInvalidPermlink - error if permlink is invalid
func ErrInvalidPermlink() sdk.Error {
	return types.NewError(types.CodeInvalidPermlink, fmt.Sprintf("invalid permlink"))
}

// ErrReasonTooLong - error if proposal reason is invalid
func ErrReasonTooLong() sdk.Error {
	return types.NewError(types.CodeReasonTooLong, fmt.Sprintf("reason length is too long"))
}

// ErrInvalidLink - error if proposal link is invalid
func ErrInvalidLink() sdk.Error {
	return types.NewError(types.CodeInvalidLink, fmt.Sprintf("invalid Link"))
}

// ErrCensorshipPostNotFound - error if content censhorship post is not found
func ErrCensorshipPostNotFound() sdk.Error {
	return types.NewError(types.CodeCensorshipPostNotFound, fmt.Sprintf("Censorship post not found"))
}

// ErrIllegalParameter - error if parameter is invalid
func ErrIllegalParameter() sdk.Error {
	return types.NewError(types.CodeIllegalParameter, fmt.Sprintf("invalid parameter"))
}

// ErrQueryFailed - error when query proposal store failed
func ErrQueryFailed() sdk.Error {
	return types.NewError(types.CodeProposalQueryFailed, fmt.Sprintf("query proposal store failed"))
}
