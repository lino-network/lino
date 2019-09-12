package repv2

import (
	"encoding/json"
	"math/big"
)

type Int struct {
	*big.Int
}

func NewInt(v int64) Int {
	return Int{
		big.NewInt(v),
	}
}

func NewIntFromBig(v *big.Int) Int {
	return Int{big.NewInt(0).Set(v)}
}

func (i Int) Clone() Int {
	return Int{big.NewInt(0).Set(i.Int)}
}

// a bunch of helper functions that takes two bitInt and returns
// a newly allocated bigInt.
func IntMax(a, b Int) Int {
	if a.Int.Cmp(b.Int) > 0 {
		return a.Clone()
	} else {
		return b.Clone()
	}
}

func IntMin(a, b Int) Int {
	if a.Int.Cmp(b.Int) < 0 {
		return a.Clone()
	} else {
		return b.Clone()
	}
}

func IntAdd(a, b Int) Int {
	rst := big.NewInt(0)
	return Int{rst.Add(a.Int, b.Int)}
}

func IntSub(a, b Int) Int {
	rst := big.NewInt(0)
	return Int{rst.Sub(a.Int, b.Int)}
}

func IntMul(a, b Int) Int {
	rst := big.NewInt(0)
	return Int{rst.Mul(a.Int, b.Int)}
}

func IntDiv(a, b Int) Int {
	rst := big.NewInt(0)
	return Int{rst.Div(a.Int, b.Int)}
}

func IntGreater(a, b Int) bool {
	return a.Cmp(b) > 0
}

func IntGTE(a, b Int) bool {
	return a.Cmp(b) >= 0
}

func IntLess(a, b Int) bool {
	return a.Cmp(b) < 0
}

// return v / (num / denom)
func IntDivFrac(v Int, num, denom int64) Int {
	if num == 0 || denom == 0 {
		panic("bigIntDivFrac zero num or denom")
	}
	return IntMulFrac(v, denom, num)
}

// return v * (num / denom)
func IntMulFrac(v Int, num, denom int64) Int {
	if denom == 0 {
		panic("bigIntMulFrac zero denom")
	}
	return IntDiv(IntMul(v, NewInt(num)), NewInt(denom))
}

// func bigIntLTE(a, b *big.Int) bool {
// 	return a.Cmp(b) <= 0
// }

func (i Int) Add(b Int) {
	i.Int.Add(i.Int, b.Int)
}

func (i Int) Sub(b Int) {
	i.Int.Sub(i.Int, b.Int)
}

func (i Int) Div(b Int) {
	i.Int.Div(i.Int, b.Int)
}

func (i Int) Mul(b Int) {
	i.Int.Mul(i.Int, b.Int)
}

func (i Int) Cmp(b Int) int {
	return i.Int.Cmp(b.Int)
}

// MarshalAmino defines custom encoding scheme
func (i Int) MarshalAmino() (string, error) {
	if i.Int == nil { // Necessary since default Uint initialization has i.i as nil
		i.Int = new(big.Int)
	}
	return marshalAmino(i.Int)
}

// UnmarshalAmino defines custom decoding scheme
func (i *Int) UnmarshalAmino(text string) error {
	if i.Int == nil { // Necessary since default Int initialization has i.i as nil
		i.Int = new(big.Int)
	}
	return unmarshalAmino(i.Int, text)
}

// MarshalJSON defines custom encoding scheme
func (i Int) MarshalJSON() ([]byte, error) {
	if i.Int == nil { // Necessary since default Uint initialization has i.i as nil
		i.Int = new(big.Int)
	}
	return marshalJSON(i.Int)
}

// UnmarshalJSON defines custom decoding scheme
func (i *Int) UnmarshalJSON(bz []byte) error {
	if i.Int == nil { // Necessary since default Int initialization has i.i as nil
		i.Int = new(big.Int)
	}
	return unmarshalJSON(i.Int, bz)
}

// MarshalAmino for custom encoding scheme
func marshalAmino(i *big.Int) (string, error) {
	bz, err := i.MarshalText()
	return string(bz), err
}

func unmarshalText(i *big.Int, text string) error {
	if err := i.UnmarshalText([]byte(text)); err != nil {
		return err
	}
	return nil
}

// UnmarshalAmino for custom decoding scheme
func unmarshalAmino(i *big.Int, text string) (err error) {
	return unmarshalText(i, text)
}

// MarshalJSON for custom encoding scheme
// Must be encoded as a string for JSON precision
func marshalJSON(i *big.Int) ([]byte, error) {
	text, err := i.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(text))
}

// UnmarshalJSON for custom decoding scheme
// Must be encoded as a string for JSON precision
func unmarshalJSON(i *big.Int, bz []byte) error {
	var text string
	err := json.Unmarshal(bz, &text)
	if err == nil {
		return unmarshalText(i, text)
	}

	// backward compatibility for old big int.
	num := big.NewInt(0)
	err = json.Unmarshal(bz, &num)
	if err != nil {
		return err
	}
	return unmarshalText(i, num.String())
}
