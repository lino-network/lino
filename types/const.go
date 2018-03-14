package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ABCI Response Codes
	// Base SDK reserves 0 ~ 99.
	// Coin errors reserve 100 ~ 199.
	// Lino authentication errors reserve 200 ~ 299.
	// Lino account handler errors reserve 300 ~ 399.
	// CodeInvalidUsername indicates the username format is invalid.
	CodeInvalidUsername sdk.CodeType = 301
)
