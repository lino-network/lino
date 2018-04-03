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
func ErrSetVoter() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Set voter failed"))
}

func ErrGetVoter() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get voter failed"))
}

func ErrGetDelegation() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Get delegation failed"))
}

func ErrSetDelegation() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Set delegation failed"))
}

func ErrVoterMarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Voter marshal error: %s", err.Error()))
}

func ErrVoterUnmarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Voter unmarshal error: %s", err.Error()))
}

func ErrDelegationMarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Delegation marshal error: %s", err.Error()))
}

func ErrDelegationUnmarshalError(err error) sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Delegation unmarshal error: %s", err.Error()))
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

func ErrAccountCoinNotEnough() sdk.Error {
	return newError(types.CodeVoteManagerFailed, fmt.Sprintf("Account bank's coins are not enough"))
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
