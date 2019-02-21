package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// PrecHoudreds prec of rat(xxx, 10)
	PrecTenths = 1
	// PrecHoudreds prec of rat(xxx, 100)
	PrecHoudreds  = 2
	// PrecThousands prec of rat(xxx, 1000)
	PrecThousands = 3
)

// NewDecFromRat converting a / b to Dec, no float involved.
func NewDecFromRat(a, b int64) sdk.Dec {
	decA := sdk.NewDec(a)
	decB := sdk.NewDec(b)
	return decA.Quo(decB)
}