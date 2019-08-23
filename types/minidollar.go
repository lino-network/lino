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

func NewMiniDollarFromInt(i sdk.Int) MiniDollar {
	return MiniDollar{i}
}

func NewMiniDollarFromBig(v *big.Int) MiniDollar {
	return MiniDollar{sdk.NewIntFromBigInt(v)}
}

// TODO(yumin): MUST DELETE on upgrade-3
func NewMiniDollarFromTestnetCoin(c Coin) MiniDollar {
	rst := NewMiniDollarFromBig(c.Amount.BigInt())
	rst.Mul(sdk.NewInt(12))
	return rst
}
