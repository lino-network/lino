package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wire "github.com/cosmos/cosmos-sdk/wire"
)

func TestNewCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputInt64   int64
		expectedCoin Coin
	}{
		{1, Coin{1}},
		{0, Coin{0}},
		{-1, Coin{-1}},
		{9223372036854775807, Coin{9223372036854775807}},
		{-9223372036854775808, Coin{-9223372036854775808}},
	}

	for _, tc := range cases {
		assert.Equal(NewCoin(tc.inputInt64), tc.expectedCoin)
	}
}

func TestLNOToCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputString  string
		expectedCoin Coin
	}{
		{"1", NewCoin(1 * Decimals)},
		{"92233720368547", NewCoin(92233720368547 * Decimals)},
		{"0.001", NewCoin(0.001 * Decimals)},
		{"0.0001", NewCoin(0.0001 * Decimals)},
		{"0.00001", NewCoin(0.00001 * Decimals)},
	}

	for _, tc := range cases {
		coin, err := LinoToCoin(LNO(tc.inputString))
		assert.Nil(err)
		assert.Equal(coin, tc.expectedCoin)
	}

	invalidCases := []struct {
		inputString string
	}{
		{"92233720368548"},
		{"-1"},
		{"-0.1"},
		{"922337203685470"},
		{"0.000001"},
		{"0"},
	}

	for _, tc := range invalidCases {
		_, err := LinoToCoin(LNO(tc.inputString))
		assert.NotNil(err)
	}
}

func TestIsPositiveCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		expected bool
	}{
		{NewCoin(1), true},
		{NewCoin(0), false},
		{NewCoin(-1), false},
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
		{NewCoin(1), true},
		{NewCoin(0), true},
		{NewCoin(-1), false},
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
		{NewCoin(1), NewCoin(1), true},
		{NewCoin(2), NewCoin(1), true},
		{NewCoin(-1), NewCoin(5), false},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsGTE(tc.inputTwo)
		assert.Equal(tc.expected, res)
	}
}

func TestIsEqualCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewCoin(1), NewCoin(1), true},
		{NewCoin(1), NewCoin(10), false},
		{NewCoin(-11), NewCoin(10), false},
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
		{NewCoin(1), NewCoin(1), NewCoin(2)},
		{NewCoin(-4), NewCoin(5), NewCoin(1)},
		{NewCoin(-1), NewCoin(1), NewCoin(0)},
	}

	for _, tc := range cases {
		res := tc.inputOne.Plus(tc.inputTwo)
		assert.True(tc.expected.IsEqual(res))
	}
}

func TestMinusCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected Coin
	}{
		{NewCoin(1), NewCoin(1), NewCoin(0)},
		{NewCoin(-4), NewCoin(5), NewCoin(-9)},
		{NewCoin(10), NewCoin(1), NewCoin(9)},
	}

	for _, tc := range cases {
		res := tc.inputOne.Minus(tc.inputTwo)
		assert.True(tc.expected.IsEqual(res))
	}
}

var cdc = wire.NewCodec() //var jsonCdc JSONCodec // TODO wire.Codec

func TestSerializationGoWire(t *testing.T) {
	r := NewCoin(100000)

	//bz, err := json.Marshal(r)
	bz, err := cdc.MarshalJSON(r)
	assert.Nil(t, err)

	//str, err := r.MarshalJSON()
	//require.Nil(t, err)

	r2 := NewCoin(0)
	//err = json.Unmarshal([]byte(bz), &r2)
	err = cdc.UnmarshalJSON([]byte(bz), &r2)
	//panic(fmt.Sprintf("debug bz: %v\n", string(bz)))
	assert.Nil(t, err)

	assert.True(t, r.IsEqual(r2), "original: %v, unmarshalled: %v", r, r2)
}
