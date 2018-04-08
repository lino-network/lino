package vote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case types.CodeInvalidUsername:
		return "Invalid username format"
	case types.CodeAccRegisterFailed:
		return "Voter register failed"
	case types.CodeVoteHandlerFailed:
		return "Vote handler failed"
	case types.CodeVoteManagerFailed:
		return "Vote manager failed"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// Error constructors
func ErrGetVoter() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get voter failed"))
}

func ErrGetVote() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get vote failed"))
}

func ErrGetProposal() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get proposal failed"))
}

func ErrGetDelegation() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get delegation failed"))
}

func ErrUsernameNotFound() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Username not found"))
}

func ErrIllegalWithdraw() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Illegal withdraw"))
}

func ErrRegisterFeeNotEnough() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Register fee not enough"))
}

func ErrInvalidUsername() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Invalida Username"))
}

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	} else {
		return codeToDefaultMsg(code)
	}
}

func newError(code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}
