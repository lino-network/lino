package developer

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
	return types.NewError(types.CodeInsufficientDeveloperDeposit, fmt.Sprintf("developer deposit not enough"))
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
