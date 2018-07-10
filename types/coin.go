package types

import (
	"fmt"
	"math"

	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cznic/mathutil"
)

type LNO = string

var (
	LowerBoundRat = big.NewRat(1, Decimals)
	UpperBoundRat = big.NewRat(math.MaxInt64/Decimals, 1)
)

// Coin holds some amount of one currency
type Coin struct {
	// Amount *big.Int `json:"amount"`
	Amount mathutil.Int128 `json:"amount"`
}

func NewCoinFromInt64(amount int64) Coin {
	// return Coin{big.NewInt(amount)}
	return Coin{new(mathutil.Int128).SetInt64(amount)}
}

func NewCoinFromBigInt(amount *big.Int) (Coin, sdk.Error) {
	// return Coin{big.NewInt(amount)}
	r, err := new(mathutil.Int128).SetBigInt(amount)
	if err != nil {
		return NewCoinFromInt64(0), ErrInvalidCoins("Invalid rat")
	}
	return NewCoin(r), nil
}

func NewCoin(amount mathutil.Int128) Coin {
	// return Coin{big.NewInt(amount)}
	return Coin{amount}
}

func LinoToCoin(lino LNO) (Coin, sdk.Error) {
	num, success := new(big.Rat).SetString(lino)
	if !success {
		return NewCoinFromInt64(0), ErrInvalidCoins("Illegal LNO")
	}
	if num.Cmp(UpperBoundRat) > 0 {
		return NewCoinFromInt64(0), ErrInvalidCoins("LNO overflow")
	}
	if num.Cmp(LowerBoundRat) < 0 {
		return NewCoinFromInt64(0), ErrInvalidCoins("LNO can't be less than lower bound")
	}
	return RatToCoin(sdk.Rat{*new(big.Rat).Mul(num, big.NewRat(Decimals, 1))})
}

var (
	zero  = big.NewInt(0)
	one   = big.NewInt(1)
	two   = big.NewInt(2)
	five  = big.NewInt(5)
	nFive = big.NewInt(-5)
	ten   = big.NewInt(10)
)

func RatToCoin(rat sdk.Rat) (Coin, sdk.Error) {
	//return Coin{rat.EvaluateBig()}

	// num := rat.Num()
	// denom := rat.Denom()

	// d, rem := new(big.Int), new(big.Int)
	// d.QuoRem(num, denom, rem)
	// if rem.Cmp(zero) == 0 { // is the remainder zero
	// 	return NewCoinFromBigInt(d)
	// }

	// // evaluate the remainder using bankers rounding
	// tenNum := new(big.Int).Mul(num, ten)
	// tenD := new(big.Int).Mul(d, ten)
	// remainderDigit := new(big.Int).Sub(new(big.Int).Quo(tenNum, denom), tenD) // get the first remainder digit
	// isFinalDigit := (new(big.Int).Rem(tenNum, denom).Cmp(zero) == 0)          // is this the final digit in the remainder?

	// switch {
	// case isFinalDigit && (remainderDigit.Cmp(five) == 0 || remainderDigit.Cmp(nFive) == 0):
	// 	dRem2 := new(big.Int).Rem(d, two)
	// 	return NewCoinFromBigInt(new(big.Int).Add(d, dRem2)) // always rounds to the even number
	// case remainderDigit.Cmp(five) != -1: //remainderDigit >= 5:
	// 	d.Add(d, one)
	// case remainderDigit.Cmp(nFive) != 1: //remainderDigit <= -5:
	// 	d.Sub(d, one)
	// }
	return NewCoinFromBigInt(rat.EvaluateBig())
}

func (coin Coin) ToRat() sdk.Rat {
	return sdk.Rat{*new(big.Rat).SetInt(coin.Amount.BigInt())}
}

func (coin Coin) ToInt64() int64 {
	return coin.Amount.BigInt().Int64()
}

// String provides a human-readable representation of a coin
func (coin Coin) String() string {
	return fmt.Sprintf("coin:%v", coin.Amount)
}

// IsZero returns if this represents no money
func (coin Coin) IsZero() bool {
	return coin.Amount.Sign() == 0
}

// IsGT returns true if the receiver is greater value
func (coin Coin) IsGT(other Coin) bool {
	return coin.Amount.Cmp(other.Amount) > 0
}

// IsGTE returns true if they are the same type and the receiver is
// an equal or greater value
func (coin Coin) IsGTE(other Coin) bool {
	return coin.Amount.Cmp(other.Amount) >= 0
}

// IsEqual returns true if the two sets of Coins have the same value
func (coin Coin) IsEqual(other Coin) bool {
	return coin.Amount.Cmp(other.Amount) == 0
}

// IsPositive returns true if coin amount is positive
func (coin Coin) IsPositive() bool {
	return coin.Amount.Sign() > 0
}

// IsNotNegative returns true if coin amount is not negative
func (coin Coin) IsNotNegative() bool {
	return coin.Amount.Sign() >= 0
}

// Adds amounts of two coins with same denom
func (coin Coin) Plus(coinB Coin) Coin {
	r, cy := coin.Amount.Add(coinB.Amount)
	if cy {
		panic("overflow")
	}
	return NewCoin(r)
}

// Subtracts amounts of two coins with same denom
func (coin Coin) Minus(coinB Coin) Coin {
	negNum, success := coinB.Amount.Neg()
	if !success {
		panic("overflow")
	}
	r, cy := coin.Amount.Add(negNum)
	if cy {
		panic("overflow")
	}
	return NewCoin(r)
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
