package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IDAStr - string representation of the number of IDA.
// same as coin, support at most 5 digits precision(log10(Decimals)).
type IDAStr string

// MiniIDA is an unsigned integer, >= 0 <= max.
// One MiniIDA = 100000 IDA(Decimals).
type MiniIDA = sdk.Int

func (i IDAStr) ToIDA() (MiniIDA, sdk.Error) {
	dec, err := sdk.NewDecFromStr(string(i))
	if err != nil {
		return MiniIDA(sdk.NewInt(0)), ErrInvalidIDAAmount("Illegal IDA amount")
	}
	if dec.GT(UpperBoundRat) {
		return MiniIDA(sdk.NewInt(0)), ErrInvalidIDAAmount("IDA overflow")
	}
	if dec.LT(LowerBoundRat) {
		return MiniIDA(sdk.NewInt(0)), ErrInvalidIDAAmount("IDA can't be less than lower bound")
	}
	return MiniIDA(dec.MulInt64(Decimals).RoundInt()), nil
}

func MiniIDAToMiniDollar(amount MiniIDA, miniIDAPrice MiniDollar) MiniDollar {
	return MiniDollar{miniIDAPrice.Mul(amount)}
}
