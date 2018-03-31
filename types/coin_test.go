package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
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
		{"1", Coin{1 * Decimals}},
		{"92233720368547", Coin{92233720368547 * Decimals}},
		{"0.001", Coin{0.001 * Decimals}},
		{"0.0001", Coin{0.0001 * Decimals}},
		{"0.00001", Coin{0.00001 * Decimals}},
	}

	for _, tc := range cases {
		rat, err := sdk.NewRatFromDecimal(tc.inputString)
		assert.Nil(err)
		coin, err := LinoToCoin(LNO(rat))
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
		rat, err := sdk.NewRatFromDecimal(tc.inputString)
		assert.Nil(err)
		_, err = LinoToCoin(LNO(rat))
		assert.NotNil(err)
	}
}

func TestIsPositiveCoin(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		inputOne Coin
		expected bool
	}{
		{Coin{1}, true},
		{Coin{0}, false},
		{Coin{-1}, false},
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
		{Coin{1}, true},
		{Coin{0}, true},
		{Coin{-1}, false},
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
		{Coin{1}, Coin{1}, true},
		{Coin{2}, Coin{1}, true},
		{Coin{-1}, Coin{5}, false},
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
		{Coin{1}, Coin{1}, true},
		{Coin{1}, Coin{10}, false},
		{Coin{-11}, Coin{10}, false},
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
		{Coin{1}, Coin{1}, Coin{2}},
		{Coin{-4}, Coin{5}, Coin{1}},
		{Coin{-1}, Coin{1}, Coin{0}},
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
		{Coin{1}, Coin{1}, Coin{0}},
		{Coin{-4}, Coin{5}, Coin{-9}},
		{Coin{10}, Coin{1}, Coin{9}},
	}

	for _, tc := range cases {
		res := tc.inputOne.Minus(tc.inputTwo)
		assert.Equal(tc.expected, res)
	}
}
