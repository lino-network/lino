package types

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LNO = string

var (
	LowerBoundRat = sdk.NewRat(1, Decimals)
	UpperBoundRat = sdk.NewRat(math.MaxInt64 / Decimals)
)

// Coin holds some amount of one currency
type Coin struct {
	// Amount *big.Int `json:"amount"`
	Amount int64 `json:"amount"`
}

func NewCoin(amount int64) Coin {
	// return Coin{big.NewInt(amount)}
	return Coin{amount}
}

func LinoToCoin(lino LNO) (Coin, sdk.Error) {
	num, err := sdk.NewRatFromDecimal(lino)
	if err != nil {
		return NewCoin(0), sdk.ErrInvalidCoins("Illegal LNO")
	}
	if num.GT(UpperBoundRat) {
		return NewCoin(0), sdk.ErrInvalidCoins("LNO overflow")
	}
	if num.LT(LowerBoundRat) {
		return NewCoin(0), sdk.ErrInvalidCoins("LNO can't be less than lower bound")
	}
	return RatToCoin(num.Mul(sdk.NewRat(Decimals))), nil
}

func RatToCoin(rat sdk.Rat) Coin {
	//return Coin{rat.EvaluateBig()}
	return Coin{rat.Evaluate()}
}

func (coin Coin) ToRat() sdk.Rat {
	return sdk.NewRat(coin.Amount)
}

func (coin Coin) ToInt64() int64 {
	return coin.Amount
}

// String provides a human-readable representation of a coin
func (coin Coin) String() string {
	return fmt.Sprintf("coin:%v", coin.Amount)
}

// IsZero returns if this represents no money
func (coin Coin) IsZero() bool {
	return coin.Amount == 0
}

// IsGT returns true if the receiver is greater value
func (coin Coin) IsGT(other Coin) bool {
	return coin.Amount > other.Amount
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
	return (coin.Amount > 0)
}

// IsNotNegative returns true if coin amount is not negative
func (coin Coin) IsNotNegative() bool {
	return (coin.Amount >= 0)
}

// Adds amounts of two coins with same denom
func (coin Coin) Plus(coinB Coin) Coin {
	return Coin{coin.Amount + coinB.Amount}
}

// Subtracts amounts of two coins with same denom
func (coin Coin) Minus(coinB Coin) Coin {
	return Coin{coin.Amount - coinB.Amount}
}

// TODO(Lino) wait until https://github.com/cosmos/cosmos-sdk/issues/785 pass

// // IsZero returns if this contains 0 amount of coin
// func (coin Coin) IsZero() bool {
// 	return coin.Amount.Cmp(big.NewInt(0)) == 0
// }

// // IsGTE returns true if the receiver is an equal or greater value
// func (coin Coin) IsGTE(other Coin) bool {
// 	return coin.Amount.Cmp(other.Amount) >= 0
// }

// // IsEqual returns true if the two coin have the same value
// func (coin Coin) IsEqual(other Coin) bool {
// 	return coin.Amount.Cmp(other.Amount) == 0
// }

// // IsPositive returns true if coin amount is positive
// func (coin Coin) IsPositive() bool {
// 	return coin.Amount.Sign() > 0
// }

// // IsNotNegative returns true if coin amount is not negative
// func (coin Coin) IsNotNegative() bool {
// 	return coin.Amount.Sign() >= 0
// }

// // Adds amounts of two coins with same denom
// func (coin Coin) Plus(coinB Coin) Coin {
// 	return Coin{new(big.Int).Add(coin.Amount, coinB.Amount)}
// }

// // Subtracts amounts of two coins with same denom
// func (coin Coin) Minus(coinB Coin) Coin {
// 	return Coin{new(big.Int).Sub(coin.Amount, coinB.Amount)}
// }

// func (coin Coin) UnmarshalJSON(coinBytes []byte) error {
// 	fmt.Println(string(coinBytes))
// 	bigint, ok := new(big.Int).SetString(string(coinBytes), 10)
// 	if !ok {
// 		return sdk.ErrInvalidCoins("parse coin failed")
// 	}
// 	coin.Amount = bigint
// 	return nil
// }
