package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewDecFromRat converting a / b to Dec, no float involved.
func NewDecFromRat(a, b int64) sdk.Dec {
	decA := sdk.NewDec(a)
	decB := sdk.NewDec(b)
	return decA.Quo(decB)
}
