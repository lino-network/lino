package validator

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidatorRevokeMsg(t *testing.T) {
	cases := []struct {
		validatorRevokeMsg ValidatorRevokeMsg
		expectError        sdk.Error
	}{
		{NewValidatorRevokeMsg("user1"), nil},
		{NewValidatorRevokeMsg(""), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.validatorRevokeMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestValidatorWithdrawMsg(t *testing.T) {
	cases := []struct {
		validatorWithdrawMsg ValidatorWithdrawMsg
		expectError          sdk.Error
	}{
		{NewValidatorWithdrawMsg("user1", sdk.NewRat(1)), nil},
		{NewValidatorWithdrawMsg("", sdk.NewRat(1)), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.validatorWithdrawMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}
