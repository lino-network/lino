package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// add test for NewDecFromRat
func TestDecFromRat(t *testing.T) {
	assert := assert.New(t)
	rst := NewDecFromRat(1, 3)
	assert.Equal(sdk.MustNewDecFromStr("0.333333333333333333"), rst)
}

// precision lost case.
// NOTE(yumin): this test case does not guarantee anything, instead, it shows
// a case where precision can be lost if you do:
// rst = (a / b) * c
// instead of,
// rst = (a * c) / b
// we should use the latter form throughout the project, if possible, namely @mul-first-form.
func TestErrors(t *testing.T) {
	assert := assert.New(t)
	expected := NewDecFromRat(316250000, 63)
	a := sdk.NewDec(316250000)
	b := sdk.NewDec(63)
	ratio := sdk.NewDec(1).Quo(b)
	rst := a.Mul(ratio)
	assert.NotEqual(expected, rst)
}
