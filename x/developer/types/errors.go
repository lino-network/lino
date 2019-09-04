package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrAccountNotFound - error if account doesn't exist
func ErrAccountNotFound() sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("account not found"))
}

// ErrDeveloperAlreadyExist - error if developer is already registered
func ErrDeveloperAlreadyExist(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeDeveloperAlreadyExist, fmt.Sprintf("developer %v already exist", username))
}

// ErrDeveloperNotFound - error if developer not found
func ErrDeveloperNotFound() sdk.Error {
	return types.NewError(types.CodeDeveloperNotFound, fmt.Sprintf("developer not found"))
}

// ErrInsufficientDeveloperDeposit - error if developer deposit is insufficient
func ErrInsufficientDeveloperDeposit() sdk.Error {
	return types.NewError(types.CodeInsufficientDeveloperDeposit, fmt.Sprintf("developer deposit(stake-in) not enough"))
}

// ErrInvalidUsername - error if username invalid
func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("Invalid Username"))
}

// ErrInvalidWebsite - error if website length invalid
func ErrInvalidWebsite() sdk.Error {
	return types.NewError(types.CodeInvalidWebsite, fmt.Sprintf("Invalid website"))
}

// ErrInvalidDescription - error if description length invalid
func ErrInvalidDescription() sdk.Error {
	return types.NewError(types.CodeInvalidDescription, fmt.Sprintf("Invalid description"))
}

// ErrInvalidAppMetadata - error if app metadata length invalid
func ErrInvalidAppMetadata() sdk.Error {
	return types.NewError(types.CodeInvalidAppMetadata, fmt.Sprintf("Invalid metadata"))
}

// ErrInvalidAuthorizedApp - error if auth app target is invalid
func ErrInvalidAuthorizedApp() sdk.Error {
	return types.NewError(types.CodeInvalidAuthorizedApp, fmt.Sprintf("invalid authorized app"))
}

// ErrInvalidValidityPeriod - error if validity is invalid
func ErrInvalidValidityPeriod() sdk.Error {
	return types.NewError(types.CodeInvalidValidityPeriod, fmt.Sprintf("invalid grant validity period"))
}

// ErrGrantPermissionTooHigh - error if grant permission is not supported
func ErrGrantPermissionTooHigh() sdk.Error {
	return types.NewError(types.CodeGrantPermissionTooHigh, fmt.Sprintf("grant permission is too high"))
}

// ErrInvalidGrantPermission - error if grant permission is not supported
func ErrInvalidGrantPermission() sdk.Error {
	return types.NewError(types.CodeInvalidGrantPermission, fmt.Sprintf("grant permission is invalid"))
}

// ErrQueryFailed - error when query developer store failed
func ErrQueryFailed() sdk.Error {
	return types.NewError(types.CodeDeveloperQueryFailed, fmt.Sprintf("query developer store failed"))
}

// ErrInvalidReserveAmount - error when reserve pool amount is invalid.
func ErrInvalidReserveAmount(amount types.Coin) sdk.Error {
	return types.NewError(
		types.CodeInvalidReserveAmount, fmt.Sprintf("reserve amount is invalid: %s", amount.String()))
}

// ErrInvalidVoterDuty - error when developer attempting to be regsitered is not a voter.
func ErrInvalidVoterDuty() sdk.Error {
	return types.NewError(
		types.CodeInvalidVoterDuty, fmt.Sprintf("user's duty is not voter"))
}

// ErrInvalidUserRole - error when user's role is not valid(like it's an affiliaed account)
func ErrInvalidUserRole() sdk.Error {
	return types.NewError(
		types.CodeInvalidUserRole, fmt.Sprintf("user role is not valid to become an app"))
}

// ErrInvalidIDAName - ida name not valid.
func ErrInvalidIDAName() sdk.Error {
	return types.NewError(
		types.CodeInvalidIDAName, fmt.Sprintf("ida name must be all uppercased letter, 3 to 10"))
}

// ErrInvalidIDAPrice - IDA price is not valid.
func ErrInvalidIDAPrice() sdk.Error {
	return types.NewError(
		types.CodeInvalidIDAPrice, fmt.Sprintf("ida price must be [1,1000] int"))
}

// ErrIDATransferSelf -
func ErrIDATransferSelf() sdk.Error {
	return types.NewError(
		types.CodeIDATransferSelf, fmt.Sprintf("ida transfer receiver and sender are the same"))
}

// ErrIDAIssuedBefore - ida has been issued before.
func ErrIDAIssuedBefore() sdk.Error {
	return types.NewError(
		types.CodeIDAIssuedBefore, fmt.Sprintf("IDA has been issued"))
}

// ErrIDARevoked - ida was revoked before
func ErrIDARevoked() sdk.Error {
	return types.NewError(
		types.CodeIDARevoked, fmt.Sprintf("ida revoked"))
}

// ErrIDAUnauthed - app's authorization of ida on user is revoked.
func ErrIDAUnauthed() sdk.Error {
	return types.NewError(
		types.CodeIDAUnauthed, fmt.Sprintf("ida is unauthed by user"))
}

// ErrExchangeMiniDollarZeroAmount -
func ErrExchangeMiniDollarZeroAmount() sdk.Error {
	return types.NewError(
		types.CodeExchangeMiniDollarZeroAmount, fmt.Sprintf("trying to exchange 0 mini dollar"))
}

// ErrNotEnoughIDA -
func ErrNotEnoughIDA() sdk.Error {
	return types.NewError(
		types.CodeNotEnoughIDA, fmt.Sprintf("ida balance not enough"))
}

// ErrBurnZeroIDA -
func ErrBurnZeroIDA() sdk.Error {
	return types.NewError(
		types.CodeBurnZeroIDA, fmt.Sprintf("trying to burn zero amount of IDA"))
}

// ErrInvalidTransferTarget -
func ErrInvalidTransferTarget() sdk.Error {
	return types.NewError(
		types.CodeInvalidTransferTarget, fmt.Sprintf("can only transfer from or to app"))
}

// ErrInvalidAffiliatedAccount -
func ErrInvalidAffiliatedAccount(reason string) sdk.Error {
	return types.NewError(
		types.CodeInvalidAffiliatedAccount, reason)
}

// ErrMaxAffiliatedExceeded -
func ErrMaxAffiliatedExceeded() sdk.Error {
	return types.NewError(
		types.CodeMaxAffiliatedExceeded, fmt.Sprintf("max affiliated account exceeded"))
}

// ErrInvalidIDAAuth -
func ErrInvalidIDAAuth() sdk.Error {
	return types.NewError(
		types.CodeInvalidIDAAuth, fmt.Sprintf("invalid ida authorization update"))
}

// ErrIDANotFound -
func ErrIDANotFound() sdk.Error {
	return types.NewError(
		types.CodeIDANotFound, fmt.Sprintf("ida not found"))
}

// ErrInvalidSigner -
func ErrInvalidSigner() sdk.Error {
	return types.NewError(
		types.CodeInvalidSigner, fmt.Sprintf("invalid signer of developer"))
}
