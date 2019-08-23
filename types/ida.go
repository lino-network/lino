package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MaxIDAStrLength = 20
)

// IDAStr - string representation of the number of IDA.
type IDAStr string

// IDA is an unsigned integer, >= 0 <= max.
type IDA = sdk.Int

func (i IDAStr) ToIDA() (IDA, sdk.Error) {
	if len(i) > MaxIDAStrLength {
		return IDA(sdk.NewInt(0)), ErrInvalidIDAAmount("IDA string > MaxLength")
	}
	amount, ok := sdk.NewIntFromString(string(i))
	if !ok {
		return IDA(sdk.NewInt(0)), ErrInvalidIDAAmount("not a valid sdk.Int")
	}
	return amount, nil
}

func IDAToMiniDollar(amount IDA, idaPrice MiniDollar) MiniDollar {
	return MiniDollar{idaPrice.Mul(amount)}
}
