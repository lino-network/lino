package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"
)

type TestLNO = sdk.Rat

var ZeroTestLNO = sdk.NewRat(0)
var LowerBoundTestLNO = sdk.NewRat(1, Decimals)
var UpperBoundTestLNO = sdk.NewRat(math.MaxInt64 / Decimals)

// Coin hold some amount of one currency
type Coin struct {
	Amount int64 `json:"amount"`
}

func NewCoin(amount int64) Coin {
	return Coin{Amount: amount}
}

func LinoToCoin(lino TestLNO) (Coin, sdk.Error) {
	if lino.GT(UpperBoundTestLNO) {
		return Coin{}, sdk.ErrInvalidCoins("TestLNO overflow")
	}
	if lino.LT(LowerBoundTestLNO) {
		return Coin{}, sdk.ErrInvalidCoins("TestLNO can't be less than lower bound")
	}
	return Coin{Amount: sdk.Rat(lino).Mul(sdk.NewRat(Decimals)).Evaluate()}, nil
}

func RatToCoin(amount sdk.Rat) Coin {
	return Coin{Amount: amount.Evaluate()}
}

func (coin Coin) ToRat() sdk.Rat {
	return sdk.NewRat(coin.Amount)
}

// String provides a human-readable representation of a coin
func (coin Coin) String() string {
	return fmt.Sprintf("coin:%v", coin.Amount)
}

// IsZero returns if this represents no money
func (coin Coin) IsZero() bool {
	return coin.Amount == 0
}

// IsGTE returns true if they are the same type and the receiver is
// an equal or greater value
func (coin Coin) IsGTE(other Coin) bool {
	return coin.Amount >= other.Amount
}

// IsEqual returns true if the two sets of Coins have the same value
func (coin Coin) IsEqual(other Coin) bool {
	return coin.Amount == other.Amount
}

// IsPositive returns true if coin amount is positive
func (coin Coin) IsPositive() bool {
	return coin.Amount > 0
}

// IsNotNegative returns true if coin amount is not negative
func (coin Coin) IsNotNegative() bool {
	return coin.Amount >= 0
}

// Adds amounts of two coins with same denom
func (coin Coin) Plus(coinB Coin) Coin {
	return Coin{coin.Amount + coinB.Amount}
}

// Subtracts amounts of two coins with same denom
func (coin Coin) Minus(coinB Coin) Coin {
	return Coin{coin.Amount - coinB.Amount}
}
