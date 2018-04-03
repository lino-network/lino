package account

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
	case types.CodeAccountManagerFail:
		return "Account manager internal error"
	case types.CodeUsernameNotFound:
		return "Username not found"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// Error constructors
func ErrInvalidLinoAmount() sdk.Error {
	return newError(types.CodeInvalidMsg, fmt.Sprintf("Invalid Lino amount"))
}

func ErrUsernameNotFound() sdk.Error {
	return newError(types.CodeUsernameNotFound, fmt.Sprintf("Username not found"))
}

func ErrInvalidUsername() sdk.Error {
	return newError(types.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
}

func ErrAccountCoinNotEnough() sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("Account bank's coins are not enough"))
}

func ErrAccountCreateFail(accKey AccountKey) sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("Account exist: %v", accKey))
}

func ErrUsernameAddressMismatch() sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("Username and address mismatch"))
}

func ErrGetInfoFailed() sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("AccountManager get info failed"))
}

func ErrSetInfoFailed() sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("AccountManager set info failed"))
}

func ErrGetBankFailed() sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("AccountManager get bank failed"))
}

func ErrSetBankFailed() sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("AccountManager set bank failed"))
}

func ErrGetMetaFailed() sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("AccountManager get meta failed"))
}

func ErrSetMetaFailed() sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("AccountManager set meta failed"))
}

func ErrAddMoneyFailed() sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("Add money to bank failed"))
}

func ErrAccountMarshalError(err error) sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("Account marshal error: %s", err.Error()))
}

func ErrAccountUnmarshalError(err error) sdk.Error {
	return newError(types.CodeAccountManagerFail, fmt.Sprintf("Account unmarshal error: %s", err.Error()))
}

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}
