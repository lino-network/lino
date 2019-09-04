package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IDAStr - string representation of the number of IDA.
// same as coin, support at most 5 digits precision(log10(Decimals)).
type IDAStr string

// MiniIDA is an integer.
// 100000 MiniIDA =  one IDA(Decimals).
type MiniIDA = sdk.Int

func (i IDAStr) ToMiniIDA() (MiniIDA, sdk.Error) {
	dec, err := sdk.NewDecFromStr(string(i))
	if err != nil {
		return MiniIDA(sdk.NewInt(0)), ErrInvalidIDAAmount()
	}
	if dec.GT(UpperBoundRat) {
		return MiniIDA(sdk.NewInt(0)), ErrInvalidIDAAmount()
	}
	if dec.LT(LowerBoundRat) {
		return MiniIDA(sdk.NewInt(0)), ErrInvalidIDAAmount()
	}
	return MiniIDA(dec.MulInt64(Decimals).RoundInt()), nil
}

func MiniIDAToMiniDollar(amount MiniIDA, miniIDAPrice MiniDollar) MiniDollar {
	return MiniDollar{miniIDAPrice.Mul(amount)}
}
