package developer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrAccountNotFound() sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("account not found"))
}

func ErrDeveloperAlreadyExist(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeDeveloperAlreadyExist, fmt.Sprintf("developer %v already exist", username))
}

func ErrDeveloperNotFound() sdk.Error {
	return types.NewError(types.CodeDeveloperNotFound, fmt.Sprintf("developer not found"))
}

func ErrInsufficientDeveloperDeposit() sdk.Error {
	return types.NewError(types.CodeInsufficientDeveloperDeposit, fmt.Sprintf("developer deposit not enough"))
}

func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("Invalid Username"))
}

func ErrInvalidWebsite() sdk.Error {
	return types.NewError(types.CodeInvalidWebsite, fmt.Sprintf("Invalid website"))
}

func ErrInvalidDescription() sdk.Error {
	return types.NewError(types.CodeInvalidDescription, fmt.Sprintf("Invalid description"))
}

func ErrInvalidAppMetadata() sdk.Error {
	return types.NewError(types.CodeInvalidAppMetadata, fmt.Sprintf("Invalid metadata"))
}

func ErrInvalidAuthorizedApp() sdk.Error {
	return types.NewError(types.CodeInvalidAuthorizedApp, fmt.Sprintf("invalid authorized app"))
}

func ErrInvalidValidityPeriod() sdk.Error {
	return types.NewError(types.CodeInvalidValidityPeriod, fmt.Sprintf("invalid grant validity period"))
}

func ErrGrantPermissionTooHigh() sdk.Error {
	return types.NewError(types.CodeGrantPermissionTooHigh, fmt.Sprintf("grant permission is too high"))
}
