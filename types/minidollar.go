package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// 1 LINO  = 0.012 USD
// 10^5 Coin = 12 * 10^7 MiniDollar
// 1 coin = 1200 minidollar
var TestnetPrice = NewMiniDollar(1200)

// 1 MiniDollar = 10^(-10) USD.
type MiniDollar struct {
	// embeding sdk.Int, inheriting marshal/unmarshal function from
	// sdk.Int, so DO NOT add any other field in this struct.
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

// TODO(yumin): MUST DELETE on upgrade3
func NewMiniDollarFromTestnetCoin(c Coin) MiniDollar {
	return NewMiniDollarFromInt(c.Amount.Mul(TestnetPrice.Int))
}

func (m MiniDollar) Plus(other MiniDollar) MiniDollar {
	return MiniDollar{m.Add(other.Int)}
}

func (m MiniDollar) Minus(other MiniDollar) MiniDollar {
	return MiniDollar{m.Sub(other.Int)}
}

func (m MiniDollar) Multiply(other MiniDollar) MiniDollar {
	return MiniDollar{m.Mul(other.Int)}
}
