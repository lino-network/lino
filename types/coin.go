package types

import (
	"fmt"
	"math"

	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LNO - exposed type
type LNO = string

var (
	// LowerBoundRat - the lower bound of Rat
	LowerBoundRat = NewDecFromRat(1, Decimals)
	// UpperBoundRat - the upper bound of Rat
	UpperBoundRat = sdk.NewDec(math.MaxInt64 / Decimals)
)

// Coin - 10^5 Coin = 1 LNO
type Coin struct {
	// Amount *big.Int `json:"amount"`
	Amount sdk.Int `json:"amount"`
}

// NewCoin - return coin from sdk.Int.
func NewCoin(amount sdk.Int) Coin {
	return Coin{Amount: amount}
}

// NewCoinFromInt64 - return int64 amount of Coin
func NewCoinFromInt64(amount int64) Coin {
	// return Coin{big.NewInt(amount)}
	return Coin{sdk.NewInt(amount)}
}

// NewCoinFromBigInt - return big.Int amount of Coin
func NewCoinFromBigInt(amount *big.Int) Coin {
	sdkInt := sdk.NewIntFromBigInt(amount)
	return Coin{sdkInt}
}

// NewCoinFromString - return string amount of Coin
func NewCoinFromString(amount string) (Coin, bool) {
	res, ok := sdk.NewIntFromString(amount)
	return Coin{res}, ok
}

// LinoToCoin - convert 1 LNO to 10^5 Coin
func LinoToCoin(lino LNO) (Coin, sdk.Error) {
	rat, err := sdk.NewDecFromStr(lino)
	if err != nil {
		return NewCoinFromInt64(0), ErrInvalidCoins("Illegal LNO")
	}
	if rat.GT(UpperBoundRat) {
		return NewCoinFromInt64(0), ErrInvalidCoins("LNO overflow")
	}
	if rat.LT(LowerBoundRat) {
		return NewCoinFromInt64(0), ErrInvalidCoins("LNO can't be less than lower bound")
	}
	return DecToCoin(rat.Mul(sdk.NewDec(Decimals))), nil
}

func MustLinoToCoin(lino LNO) Coin {
	c, err := LinoToCoin(lino)
	if err != nil {
		panic(err)
	}
	return c
}

// DecToCoin - convert sdk.Dec to LNO coin
// XXX(yumin): the unit of @p rat must be coin.
func DecToCoin(rat sdk.Dec) Coin {
	return NewCoinFromBigInt(rat.RoundInt().BigInt())
}

// ToDec - convert Coin to sdk.Dec
func (coin Coin) ToDec() sdk.Dec {
	return sdk.NewDecFromBigInt(coin.Amount.BigInt())
}

// ToInt64 - convert Coin to int64
func (coin Coin) ToInt64() (int64, sdk.Error) {
	if !coin.Amount.BigInt().IsInt64() {
		return 0, ErrAmountOverflow()
	}
	return coin.Amount.BigInt().Int64(), nil
}

// String - provides a human-readable representation of a coin
func (coin Coin) String() string {
	return fmt.Sprintf("coin:%v", coin.Amount)
}

// IsZero - returns if this represents no money
func (coin Coin) IsZero() bool {
	return coin.Amount.Sign() == 0
}

// IsGT - returns true if the receiver is greater value
func (coin Coin) IsGT(other Coin) bool {
	return coin.Amount.GT(other.Amount)
}

// IsGTE - returns true if they are the same type and the receiver is
// an equal or greater value
func (coin Coin) IsGTE(other Coin) bool {
	return coin.Amount.GT(other.Amount) || coin.Amount.Equal(other.Amount)
}

// IsEqual - returns true if the two sets of Coins have the same value
func (coin Coin) IsEqual(other Coin) bool {
	return coin.Amount.Equal(other.Amount)
}

// IsPositive - returns true if coin amount is positive
func (coin Coin) IsPositive() bool {
	return coin.Amount.IsPositive()
}

// IsPositive - returns true if coin amount is positive
func (coin Coin) IsNegative() bool {
	return coin.Amount.IsNegative()
}

// IsNotNegative - returns true if coin amount is not negative
func (coin Coin) IsNotNegative() bool {
	return coin.Amount.Sign() >= 0
}

// Plus - Adds amounts of two coins with same denom
func (coin Coin) Plus(coinB Coin) Coin {
	r := coin.Amount.Add(coinB.Amount)
	return Coin{r}
}

// Minus - Subtracts amounts of two coins with same denom
func (coin Coin) Minus(coinB Coin) Coin {
	sdkInt := coin.Amount.Sub(coinB.Amount)
	return Coin{sdkInt}
}

// Neg - return a new neged coin.
func (coin Coin) Neg() Coin {
	return Coin{coin.Amount.Neg()}
}

// ToLino - return LINO str
func (coin Coin) ToLino() LNO {
	return coin.Amount.QuoRaw(Decimals).String()
}
