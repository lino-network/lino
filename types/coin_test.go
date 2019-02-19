package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/codec"
)

var (
	bigString = "1000000000000000000000000"
	bigInt, _ = new(big.Int).SetString(bigString, 10)

	doubleBigString  = "2000000000000000000000000"
	doubleBigInt, _  = new(big.Int).SetString(doubleBigString, 10)
	bigCoin, _       = NewCoinFromString(bigString)
	doubleBigCoin, _ = NewCoinFromString(doubleBigString)
)

func TestNewCoinFromInt64(t *testing.T) {
	testCases := []struct {
		testName     string
		inputInt64   int64
		expectedCoin Coin
	}{
		{
			testName:     "amount 1",
			inputInt64:   1,
			expectedCoin: NewCoinFromBigInt(new(big.Int).SetInt64(1)),
		},
		{
			testName:     "amount 0",
			inputInt64:   0,
			expectedCoin: NewCoinFromBigInt(new(big.Int).SetInt64(0)),
		},
		{
			testName:     "amount -1",
			inputInt64:   -1,
			expectedCoin: NewCoinFromBigInt(new(big.Int).SetInt64(-1)),
		},
		{
			testName:     "amount 9223372036854775807",
			inputInt64:   9223372036854775807,
			expectedCoin: NewCoinFromBigInt(new(big.Int).SetInt64(9223372036854775807)),
		},
		{
			testName:     "amount -9223372036854775808",
			inputInt64:   -9223372036854775808,
			expectedCoin: NewCoinFromBigInt(new(big.Int).SetInt64(-9223372036854775808)),
		},
	}

	for _, tc := range testCases {
		coin := NewCoinFromInt64(tc.inputInt64)
		if !coin.IsEqual(tc.expectedCoin) {
			t.Errorf("%s: diff coin, got %v, want %v", tc.testName, coin.Amount, tc.expectedCoin)
		}
	}
}

func TestNewCoinFromBigInt(t *testing.T) {
	testCases := []struct {
		testName    string
		inputBigInt *big.Int
		expectCoin  Coin
	}{
		{
			testName:    "amount 1",
			inputBigInt: new(big.Int).SetInt64(1),
			expectCoin:  NewCoinFromInt64(1),
		},
		{
			testName:    "amount 0",
			inputBigInt: new(big.Int).SetInt64(0),
			expectCoin:  NewCoinFromInt64(0),
		},
		{
			testName:    "amount -1",
			inputBigInt: new(big.Int).SetInt64(-1),
			expectCoin:  NewCoinFromInt64(-1),
		},
		{
			testName:    "amount 100",
			inputBigInt: new(big.Int).SetInt64(100),
			expectCoin:  NewCoinFromInt64(100),
		},
		{
			testName:    "amount 9223372036854775807",
			inputBigInt: new(big.Int).SetInt64(9223372036854775807),
			expectCoin:  NewCoinFromInt64(9223372036854775807),
		},
		{
			testName:    "amount -9223372036854775808",
			inputBigInt: new(big.Int).SetInt64(-9223372036854775808),
			expectCoin:  NewCoinFromInt64(-9223372036854775808),
		},
		{
			testName:    "amount bigInt",
			inputBigInt: bigInt,
			expectCoin:  bigCoin,
		},
	}

	for _, tc := range testCases {
		coin := NewCoinFromBigInt(tc.inputBigInt)
		if !coin.IsEqual(tc.expectCoin) {
			t.Errorf("%s: diff amount, got %v, want %v", tc.testName, coin.Amount, tc.expectCoin.Amount)
		}
	}
}

func TestLNOToCoin(t *testing.T) {
	testCases := []struct {
		testName     string
		inputString  string
		expectCoin   Coin
		expectResult sdk.Error
	}{
		{
			testName:     "LNO 1",
			inputString:  "1",
			expectCoin:   NewCoinFromInt64(1 * Decimals),
			expectResult: nil,
		},
		{
			testName:     "LNO 92233720368547",
			inputString:  "92233720368547",
			expectCoin:   NewCoinFromInt64(92233720368547 * Decimals),
			expectResult: nil,
		},
		{
			testName:     "LNO 0.001",
			inputString:  "0.001",
			expectCoin:   NewCoinFromInt64(0.001 * Decimals),
			expectResult: nil,
		},
		{
			testName:     "LNO 0.0001",
			inputString:  "0.0001",
			expectCoin:   NewCoinFromInt64(0.0001 * Decimals),
			expectResult: nil,
		},
		{
			testName:     "LNO 0.00001",
			inputString:  "0.00001",
			expectCoin:   NewCoinFromInt64(0.00001 * Decimals),
			expectResult: nil,
		},
		{
			testName:     "less than lower bound LNO is invalid",
			inputString:  "0.000001",
			expectCoin:   NewCoinFromInt64(0),
			expectResult: ErrInvalidCoins("Illegal LNO"),
		},
		{
			testName:     "0 LNO is invalid",
			inputString:  "0",
			expectCoin:   NewCoinFromInt64(0),
			expectResult: ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName:     "negative LNO is invalid",
			inputString:  "-1",
			expectCoin:   NewCoinFromInt64(0),
			expectResult: ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName:     "negative -0.1 LNO is invalid",
			inputString:  "-0.1",
			expectCoin:   NewCoinFromInt64(0),
			expectResult: ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName:     "overflow LNO is invalid",
			inputString:  "92233720368548",
			expectCoin:   NewCoinFromInt64(0),
			expectResult: ErrInvalidCoins("LNO overflow"),
		},
		{
			testName:     "illegal coin",
			inputString:  "1$",
			expectCoin:   NewCoinFromInt64(0),
			expectResult: ErrInvalidCoins("Illegal LNO"),
		},
		{
			testName:     "large amount of coin",
			inputString:  "1e9999999999999999",
			expectCoin:   NewCoinFromInt64(0),
			expectResult: ErrInvalidCoins("Illegal LNO"),
		},
	}

	for _, tc := range testCases {
		coin, err := LinoToCoin(LNO(tc.inputString))
		if !assert.Equal(t, tc.expectResult, err) {
			t.Errorf("%s: diff err, got %v, want %v", tc.testName, err, tc.expectResult)
		}
		if !coin.IsEqual(tc.expectCoin) {
			t.Errorf("%s: diff coin, got %v, want %v", tc.testName, coin, tc.expectCoin)
		}
	}
}

func TestRatToCoin(t *testing.T) {
	testCases := []struct {
		testName    string
		inputString string
		expectCoin  Coin
	}{
		{
			testName:    "Coin 1",
			inputString: "1",
			expectCoin:  NewCoinFromInt64(1),
		},
		{
			testName:    "Coin 0",
			inputString: "0",
			expectCoin:  NewCoinFromInt64(0),
		},
		{
			testName:    "Coin -1",
			inputString: "-1",
			expectCoin:  NewCoinFromInt64(-1),
		},
		{
			testName:    "Coin 100000",
			inputString: "100000",
			expectCoin:  NewCoinFromInt64(100000),
		},
		{
			testName:    "Coin 0.5",
			inputString: "0.5",
			expectCoin:  NewCoinFromInt64(0),
		},
		{
			testName:    "Coin 0.6 will be rounded to 1",
			inputString: "0.6",
			expectCoin:  NewCoinFromInt64(1),
		},
		{
			testName:    "Coin 1.4 will be rounded to 1",
			inputString: "1.4",
			expectCoin:  NewCoinFromInt64(1),
		},
		{
			testName:    "Coin 1 will be rounded to 2",
			inputString: "1.5",
			expectCoin:  NewCoinFromInt64(2),
		},
		{
			testName:    "Coin 9223372036854775807",
			inputString: "9223372036854775807",
			expectCoin:  NewCoinFromInt64(9223372036854775807),
		},
		{
			testName:    "Coin -9223372036854775807",
			inputString: "-9223372036854775807",
			expectCoin:  NewCoinFromInt64(-9223372036854775807),
		},
		{
			testName:    "Coin bigString",
			inputString: bigString,
			expectCoin:  bigCoin,
		},
	}

	for _, tc := range testCases {
		bigRat, success := new(big.Rat).SetString(tc.inputString)
		if !success {
			t.Errorf("%s: failed to convert input to big rat", tc.testName)
		}

		rat := sdk.Rat{Rat: bigRat}
		coin := RatToCoin(rat)
		if !coin.IsEqual(tc.expectCoin) {
			t.Errorf("%s: diff coin, got %v, want %v", tc.testName, coin, tc.expectCoin)
		}
	}
}

func TestIsPositiveCoin(t *testing.T) {
	testCases := []struct {
		testName     string
		inputOne     Coin
		expectResult bool
	}{
		{
			testName:     "1 is positive",
			inputOne:     NewCoinFromInt64(1),
			expectResult: true,
		},
		{
			testName:     "0 is not positive",
			inputOne:     NewCoinFromInt64(0),
			expectResult: false,
		},
		{
			testName:     "-1 is not positive",
			inputOne:     NewCoinFromInt64(-1),
			expectResult: false,
		},
		{
			testName:     "bigInt128 is positive",
			inputOne:     bigCoin,
			expectResult: true,
		},
		{
			testName:     "doubleBigInt128 is not positive",
			inputOne:     doubleBigCoin,
			expectResult: true,
		},
	}

	for _, tc := range testCases {
		res := tc.inputOne.IsPositive()
		if res != tc.expectResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectResult)
		}
	}
}

func TestIsNotNegativeCoin(t *testing.T) {
	testCases := []struct {
		testName     string
		inputOne     Coin
		expectResult bool
	}{
		{
			testName:     "1 is not negative",
			inputOne:     NewCoinFromInt64(1),
			expectResult: true,
		},
		{
			testName:     "0 is not negative",
			inputOne:     NewCoinFromInt64(0),
			expectResult: true,
		},
		{
			testName:     "-1 is negative",
			inputOne:     NewCoinFromInt64(-1),
			expectResult: false,
		},
		{
			testName:     "bigInt128 is not negative",
			inputOne:     bigCoin,
			expectResult: true,
		},
		{
			testName:     "doubleBigInt128 is not negative",
			inputOne:     doubleBigCoin,
			expectResult: true,
		},
	}

	for _, tc := range testCases {
		res := tc.inputOne.IsNotNegative()
		if res != tc.expectResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectResult)
		}
	}
}

func TestIsGTECoin(t *testing.T) {
	testCases := []struct {
		testName     string
		inputOne     Coin
		inputTwo     Coin
		expectResult bool
	}{
		{
			testName:     "inputs 1 are equal",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: true,
		},
		{
			testName:     "inputOne is bigger than inputTwo",
			inputOne:     NewCoinFromInt64(2),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: true,
		},
		{
			testName:     "inputOne is less than inputTwo",
			inputOne:     NewCoinFromInt64(-1),
			inputTwo:     NewCoinFromInt64(5),
			expectResult: false,
		},
		{
			testName:     "bigInt128 is bigger than 5",
			inputOne:     bigCoin,
			inputTwo:     NewCoinFromInt64(5),
			expectResult: true,
		},
		{
			testName:     "inputs bigInt128 are equal",
			inputOne:     bigCoin,
			inputTwo:     bigCoin,
			expectResult: true,
		},
		{
			testName:     "doubleBigInt128 is bigger than bigInt128",
			inputOne:     doubleBigCoin,
			inputTwo:     bigCoin,
			expectResult: true,
		},
	}

	for _, tc := range testCases {
		res := tc.inputOne.IsGTE(tc.inputTwo)
		if res != tc.expectResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectResult)
		}
	}
}

func TestIsGTCoin(t *testing.T) {
	testCases := []struct {
		testName     string
		inputOne     Coin
		inputTwo     Coin
		expectResult bool
	}{
		{
			testName:     "1 equals to 1",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: false,
		},
		{
			testName:     "2 is bigger than 1",
			inputOne:     NewCoinFromInt64(2),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: true,
		},
		{
			testName:     "-1 is less than 5",
			inputOne:     NewCoinFromInt64(-1),
			inputTwo:     NewCoinFromInt64(5),
			expectResult: false,
		},
		{
			testName:     "bigInt128 is bigger than 5",
			inputOne:     bigCoin,
			inputTwo:     NewCoinFromInt64(5),
			expectResult: true,
		},
		{
			testName:     "bigInt128 is not bigger than bigInt128",
			inputOne:     bigCoin,
			inputTwo:     bigCoin,
			expectResult: false,
		},
		{
			testName:     "doubleBigInt128 is bigger than bigInt128",
			inputOne:     doubleBigCoin,
			inputTwo:     bigCoin,
			expectResult: true,
		},
	}

	for _, tc := range testCases {
		res := tc.inputOne.IsGT(tc.inputTwo)
		if res != tc.expectResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectResult)
		}
	}
}

func TestIsEqualCoin(t *testing.T) {
	testCases := []struct {
		testName     string
		inputOne     Coin
		inputTwo     Coin
		expectResult bool
	}{
		{
			testName:     "bigInt128 equals to bigInt128",
			inputOne:     bigCoin,
			inputTwo:     bigCoin,
			expectResult: true,
		},
		{
			testName:     "1 equals to 1",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: true,
		},
		{
			testName:     "1 is not equal to 10",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromInt64(10),
			expectResult: false,
		},
		{
			testName:     "-11 is not equal to 10",
			inputOne:     NewCoinFromInt64(-11),
			inputTwo:     NewCoinFromInt64(10),
			expectResult: false,
		},
		{
			testName:     "1 equals to int128 1",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromBigInt(new(big.Int).SetInt64(1)),
			expectResult: true,
		},
		{
			testName:     "1 equals to coin1",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromBigInt(new(big.Int).SetInt64(1)),
			expectResult: true,
		},
	}

	for _, tc := range testCases {
		res := tc.inputOne.IsEqual(tc.inputTwo)
		if res != tc.expectResult {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectResult)
		}
	}
}

func TestPlusCoin(t *testing.T) {
	testCases := []struct {
		testName     string
		inputOne     Coin
		inputTwo     Coin
		expectResult Coin
	}{
		{
			testName:     "0 + 1 = 1",
			inputOne:     NewCoinFromInt64(0),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: NewCoinFromInt64(1),
		},
		{
			testName:     "1 + 0 = 1",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromInt64(0),
			expectResult: NewCoinFromInt64(1),
		},
		{
			testName:     "1 + 1 = 2",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: NewCoinFromInt64(2),
		},
		{
			testName:     "-4 + 5 = 1",
			inputOne:     NewCoinFromInt64(-4),
			inputTwo:     NewCoinFromInt64(5),
			expectResult: NewCoinFromInt64(1),
		},
		{
			testName:     "-1 + 1 = 0",
			inputOne:     NewCoinFromInt64(-1),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: NewCoinFromInt64(0),
		},
		{
			testName:     "bigInt128 + bigInt128 = doubleBigInt128",
			inputOne:     bigCoin,
			inputTwo:     bigCoin,
			expectResult: doubleBigCoin,
		},
	}

	for _, tc := range testCases {
		res := tc.inputOne.Plus(tc.inputTwo)
		if !res.IsEqual(tc.expectResult) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectResult)
		}
	}
}

func TestMinusCoin(t *testing.T) {
	testCases := []struct {
		testName     string
		inputOne     Coin
		inputTwo     Coin
		expectResult Coin
	}{
		{
			testName:     "0 - 0 = 0",
			inputOne:     NewCoinFromInt64(0),
			inputTwo:     NewCoinFromInt64(0),
			expectResult: NewCoinFromInt64(0),
		},
		{
			testName:     "1 - 0 = 1",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromInt64(0),
			expectResult: NewCoinFromInt64(1),
		},
		{
			testName:     "1 - 1 = 0",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: NewCoinFromInt64(0),
		},
		{
			testName:     "1 - (-1) = 2",
			inputOne:     NewCoinFromInt64(1),
			inputTwo:     NewCoinFromInt64(-1),
			expectResult: NewCoinFromInt64(2),
		},
		{
			testName:     "-1 - (-1) = 0",
			inputOne:     NewCoinFromInt64(-1),
			inputTwo:     NewCoinFromInt64(-1),
			expectResult: NewCoinFromInt64(0),
		},
		{
			testName:     "-4 - 5 = -9",
			inputOne:     NewCoinFromInt64(-4),
			inputTwo:     NewCoinFromInt64(5),
			expectResult: NewCoinFromInt64(-9),
		},
		{
			testName:     "10 - 1 = 9",
			inputOne:     NewCoinFromInt64(10),
			inputTwo:     NewCoinFromInt64(1),
			expectResult: NewCoinFromInt64(9),
		},
		{
			testName:     "bigInt128 - bigInt128 = 0",
			inputOne:     bigCoin,
			inputTwo:     bigCoin,
			expectResult: NewCoinFromInt64(0),
		},
	}

	for _, tc := range testCases {
		res := tc.inputOne.Minus(tc.inputTwo)
		if !res.IsEqual(tc.expectResult) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectResult)
		}
	}
}

var cdc = wire.NewCodec() //var jsonCdc JSONCodec // TODO wire.Codec

func TestSerializationGoWire(t *testing.T) {
	r := bigCoin

	//bz, err := json.Marshal(r)
	bz, err := cdc.MarshalJSON(r)
	if err != nil {
		t.Errorf("TestSerializationGoWire: failed to marshal, got err %v", err)
	}

	//str, err := r.MarshalJSON()
	//require.Nil(t, err)

	r2 := NewCoinFromInt64(0)
	//err = json.Unmarshal([]byte(bz), &r2)
	err = cdc.UnmarshalJSON([]byte(bz), &r2)
	//panic(fmt.Sprintf("debug bz: %v\n", string(bz)))
	if err != nil {
		t.Errorf("TestSerializationGoWire: failed to unmarshal, got err %v", err)
	}

	if !r2.IsEqual(r) {
		t.Errorf("TestSerializationGoWire: diff result, got %v, want %v", r2, r)
	}
}
