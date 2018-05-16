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

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              sdk.Msg
		expectPermission types.Permission
	}{
		"validator deposit msg": {
			NewValidatorDepositMsg(
				"test", types.LNO("1"), crypto.GenPrivKeyEd25519().PubKey(), "https://lino.network"),
			types.TransactionPermission},
		"validator withdraw msg": {
			NewValidatorWithdrawMsg("test", types.LNO("1")),
			types.TransactionPermission},
		"validator revoke msg": {
			NewValidatorRevokeMsg("test"),
			types.TransactionPermission},
	}

	for testName, cs := range cases {
		permissionLevel := cs.msg.Get(types.PermissionLevel)
		if permissionLevel == nil {
			if cs.expectPermission != types.PostPermission {
				t.Errorf(
					"%s: expect permission incorrect, expect %v, got %v",
					testName, cs.expectPermission, types.PostPermission)
				return
			} else {
				continue
			}
		}
		permission, ok := permissionLevel.(types.Permission)
		assert.Equal(t, ok, true)
		if cs.expectPermission != permission {
			t.Errorf(
				"%s: expect permission incorrect, expect %v, got %v",
				testName, cs.expectPermission, permission)
			return
		}
	}
}
