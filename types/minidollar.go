package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// 1 MiniDollar = 10^(-8) USD.
type MiniDollar struct {
	sdk.Int
}

func NewMiniDollar(v int64) MiniDollar {
	return MiniDollar{sdk.NewInt(v)}
}

func NewMiniDollarFromBig(v *big.Int) MiniDollar {
	return MiniDollar{sdk.NewIntFromBigInt(v)}
}
