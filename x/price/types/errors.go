package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

// ErrFedPriceNotFound - error when fed price is not found.
func ErrFedPriceNotFound(u linotypes.AccountKey) sdk.Error {
	return linotypes.NewError(
		linotypes.CodeFedPriceNotFound,
		fmt.Sprintf("fed price of %v is not found", u))
}

// ErrCurrentPriceNotFound - error current price is not found.
func ErrCurrentPriceNotFound() sdk.Error {
	return linotypes.NewError(
		linotypes.CodeFedPriceNotFound, "current price not found")
}

// ErrNoValidator - error when no validator is found.
func ErrNoValidator() sdk.Error {
	return linotypes.NewError(
		linotypes.CodeNoValidatorSet, fmt.Sprintf("no validator set found"))
}

// ErrNotAValidator -
func ErrNotAValidator(u linotypes.AccountKey) sdk.Error {
	return linotypes.NewError(
		linotypes.CodeNotAValidator, fmt.Sprintf("%s is not a validator", u))
}

// ErrInvalidPriceFeed -
func ErrInvalidPriceFeed(price linotypes.MiniDollar) sdk.Error {
	return linotypes.NewError(
		linotypes.CodeInvalidPriceFeed, fmt.Sprintf("invalid price: %s", price))
}

// ErrPriceFeedRateLimited -
func ErrPriceFeedRateLimited() sdk.Error {
	return linotypes.NewError(
		linotypes.CodePriceFeedRateLimited, fmt.Sprintf(""))
}
