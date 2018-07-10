package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/cznic/mathutil"
)

var (
	bigString    = "1000000000000000000000000"
	bigInt, _    = new(big.Int).SetString(bigString, 10)
	bigInt128, _ = new(mathutil.Int128).SetBigInt(bigInt)

	doubleBigString    = "2000000000000000000000000"
	doubleBigInt, _    = new(big.Int).SetString(doubleBigString, 10)
	doubleBigInt128, _ = new(mathutil.Int128).SetBigInt(doubleBigInt)
)

func TestNewCoinFromInt64(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputInt64     int64
		expectedAmount mathutil.Int128
	}{
		{1, new(mathutil.Int128).SetInt64(1)},
		{0, new(mathutil.Int128).SetInt64(0)},
		{-1, new(mathutil.Int128).SetInt64(-1)},
		{9223372036854775807, new(mathutil.Int128).SetInt64(9223372036854775807)},
		{-9223372036854775808, new(mathutil.Int128).SetInt64(-9223372036854775808)},
	}

	for _, tc := range cases {
		assert.Equal(NewCoinFromInt64(tc.inputInt64).Amount, tc.expectedAmount)
	}
}
func TestNewCoinFromBigInt(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputBigInt    *big.Int
		expectedAmount mathutil.Int128
	}{
		{new(big.Int).SetInt64(1), new(mathutil.Int128).SetInt64(1)},
		{new(big.Int).SetInt64(0), new(mathutil.Int128).SetInt64(0)},
		{new(big.Int).SetInt64(-1), new(mathutil.Int128).SetInt64(-1)},
		{new(big.Int).SetInt64(100), new(mathutil.Int128).SetInt64(100)},
		{new(big.Int).SetInt64(9223372036854775807),
			new(mathutil.Int128).SetInt64(9223372036854775807)},
		{new(big.Int).SetInt64(-9223372036854775808),
			new(mathutil.Int128).SetInt64(-9223372036854775808)},
		{bigInt, bigInt128},
	}

	for _, tc := range cases {
		coin, err := NewCoinFromBigInt(tc.inputBigInt)
		assert.Nil(err)
		assert.Equal(coin.Amount, tc.expectedAmount)
	}
}

func TestLNOToCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputString  string
		expectCoin   Coin
		expectResult sdk.Error
	}{
		{"1", NewCoinFromInt64(1 * Decimals), nil},
		{"92233720368547", NewCoinFromInt64(92233720368547 * Decimals), nil},
		{"0.001", NewCoinFromInt64(0.001 * Decimals), nil},
		{"0.0001", NewCoinFromInt64(0.0001 * Decimals), nil},
		{"0.00001", NewCoinFromInt64(0.00001 * Decimals), nil},
		{"0.000001", NewCoinFromInt64(0),
			ErrInvalidCoins("LNO can't be less than lower bound")},
		{"0", NewCoinFromInt64(0),
			ErrInvalidCoins("LNO can't be less than lower bound")},
		{"-1", NewCoinFromInt64(0),
			ErrInvalidCoins("LNO can't be less than lower bound")},
		{"-0.1", NewCoinFromInt64(0),
			ErrInvalidCoins("LNO can't be less than lower bound")},
		{"92233720368548", NewCoinFromInt64(0),
			ErrInvalidCoins("LNO overflow")},
		{"1$", NewCoinFromInt64(0), ErrInvalidCoins("Illegal LNO")},
	}

	for _, tc := range cases {
		coin, err := LinoToCoin(LNO(tc.inputString))
		assert.Equal(tc.expectResult, err)
		assert.Equal(coin, tc.expectCoin)
	}
}

func TestRatToCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputString string
		expectCoin  Coin
	}{
		{"1", NewCoinFromInt64(1)},
		{"0", NewCoinFromInt64(0)},
		{"-1", NewCoinFromInt64(-1)},
		{"100000", NewCoinFromInt64(100000)},
		{"0.5", NewCoinFromInt64(0)},
		{"0.6", NewCoinFromInt64(1)},
		{"1.4", NewCoinFromInt64(1)},
		{"1.5", NewCoinFromInt64(2)},
		{"9223372036854775807", NewCoinFromInt64(9223372036854775807)},
		{"-9223372036854775807", NewCoinFromInt64(-9223372036854775807)},
		{bigString, NewCoin(bigInt128)},
	}

	for _, tc := range cases {
		bigRat, success := new(big.Rat).SetString(tc.inputString)
		assert.True(success)
		rat := sdk.Rat{*bigRat}
		coin, changeErr := RatToCoin(rat)
		assert.Nil(changeErr)
		assert.Equal(tc.expectCoin, coin)
	}
}

func TestIsPositiveCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		expected bool
	}{
		{NewCoinFromInt64(1), true},
		{NewCoinFromInt64(0), false},
		{NewCoinFromInt64(-1), false},
		{NewCoin(bigInt128), true},
		{NewCoin(doubleBigInt128), true},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsPositive()
		assert.Equal(tc.expected, res)
	}
}

func TestIsNotNegativeCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		expected bool
	}{
		{NewCoinFromInt64(1), true},
		{NewCoinFromInt64(0), true},
		{NewCoinFromInt64(-1), false},
		{NewCoin(bigInt128), true},
		{NewCoin(doubleBigInt128), true},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsNotNegative()
		assert.Equal(tc.expected, res)
	}
}

func TestIsGTECoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewCoinFromInt64(1), NewCoinFromInt64(1), true},
		{NewCoinFromInt64(2), NewCoinFromInt64(1), true},
		{NewCoinFromInt64(-1), NewCoinFromInt64(5), false},
		{NewCoin(bigInt128), NewCoinFromInt64(5), true},
		{NewCoin(bigInt128), NewCoin(bigInt128), true},
		{NewCoin(doubleBigInt128), NewCoin(bigInt128), true},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsGTE(tc.inputTwo)
		assert.Equal(tc.expected, res)
	}
}

func TestIsGTCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewCoinFromInt64(1), NewCoinFromInt64(1), false},
		{NewCoinFromInt64(2), NewCoinFromInt64(1), true},
		{NewCoinFromInt64(-1), NewCoinFromInt64(5), false},
		{NewCoin(bigInt128), NewCoinFromInt64(5), true},
		{NewCoin(bigInt128), NewCoin(bigInt128), false},
		{NewCoin(doubleBigInt128), NewCoin(bigInt128), true},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsGT(tc.inputTwo)
		assert.Equal(tc.expected, res)
	}
}

func TestIsEqualCoin(t *testing.T) {
	assert := assert.New(t)

	coin1, err := NewCoinFromBigInt(new(big.Int).SetInt64(1))
	assert.Nil(err)

	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewCoin(bigInt128), NewCoin(bigInt128), true},
		{NewCoinFromInt64(1), NewCoinFromInt64(1), true},
		{NewCoinFromInt64(1), NewCoinFromInt64(10), false},
		{NewCoinFromInt64(-11), NewCoinFromInt64(10), false},
		{NewCoinFromInt64(1), NewCoin(new(mathutil.Int128).SetInt64(1)), true},
		{NewCoinFromInt64(1), coin1, true},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsEqual(tc.inputTwo)
		assert.Equal(tc.expected, res)
	}
}

func TestPlusCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected Coin
	}{
		{NewCoinFromInt64(0), NewCoinFromInt64(1), NewCoinFromInt64(1)},
		{NewCoinFromInt64(1), NewCoinFromInt64(0), NewCoinFromInt64(1)},
		{NewCoinFromInt64(1), NewCoinFromInt64(1), NewCoinFromInt64(2)},
		{NewCoinFromInt64(-4), NewCoinFromInt64(5), NewCoinFromInt64(1)},
		{NewCoinFromInt64(-1), NewCoinFromInt64(1), NewCoinFromInt64(0)},
		{NewCoin(bigInt128), NewCoin(bigInt128), NewCoin(doubleBigInt128)},
	}

	for _, tc := range cases {
		res := tc.inputOne.Plus(tc.inputTwo)
		assert.Equal(tc.expected, res)
	}
}

func TestMinusCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected Coin
	}{
		{NewCoinFromInt64(0), NewCoinFromInt64(0), NewCoinFromInt64(0)},
		{NewCoinFromInt64(1), NewCoinFromInt64(0), NewCoinFromInt64(1)},
		{NewCoinFromInt64(1), NewCoinFromInt64(1), NewCoinFromInt64(0)},
		{NewCoinFromInt64(1), NewCoinFromInt64(-1), NewCoinFromInt64(2)},
		{NewCoinFromInt64(-1), NewCoinFromInt64(-1), NewCoinFromInt64(0)},
		{NewCoinFromInt64(-4), NewCoinFromInt64(5), NewCoinFromInt64(-9)},
		{NewCoinFromInt64(10), NewCoinFromInt64(1), NewCoinFromInt64(9)},
		{NewCoin(bigInt128), NewCoin(bigInt128), NewCoinFromInt64(0)},
	}

	for _, tc := range cases {
		res := tc.inputOne.Minus(tc.inputTwo)
		assert.True(tc.expected.IsEqual(res))
	}
}

var cdc = wire.NewCodec() //var jsonCdc JSONCodec // TODO wire.Codec

func TestSerializationGoWire(t *testing.T) {
	r := NewCoin(bigInt128)

	//bz, err := json.Marshal(r)
	bz, err := cdc.MarshalJSON(r)
	assert.Nil(t, err)

	//str, err := r.MarshalJSON()
	//require.Nil(t, err)

	r2 := NewCoinFromInt64(0)
	//err = json.Unmarshal([]byte(bz), &r2)
	err = cdc.UnmarshalJSON([]byte(bz), &r2)
	//panic(fmt.Sprintf("debug bz: %v\n", string(bz)))
	assert.Nil(t, err)

	assert.True(t, r.IsEqual(r2), "original: %v, unmarshalled: %v", r, r2)
}
