package validator

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"
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
		{NewValidatorWithdrawMsg("user1", "1"), nil},
		{NewValidatorWithdrawMsg("", "1"), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.validatorWithdrawMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestValidatorDepositPermission(t *testing.T) {
	priv := crypto.GenPrivKeyEd25519()
	msg := NewValidatorDepositMsg("user1", "1", priv.PubKey(), "")
	permissionLevel := msg.Get(types.PermissionLevel)
	permission, ok := permissionLevel.(types.Permission)
	assert.Equal(t, ok, true)
	assert.Equal(t, permission, types.TransactionPermission)
}
