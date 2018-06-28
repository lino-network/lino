package developer

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"
)

func TestDeveloperRegisterMsg(t *testing.T) {
	cases := []struct {
		developerRegisterMsg DeveloperRegisterMsg
		expectError          sdk.Error
	}{
		{NewDeveloperRegisterMsg("user1", "10"), nil},
		{NewDeveloperRegisterMsg("", "10"), ErrInvalidUsername()},
		{NewDeveloperRegisterMsg("user1",
			"-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.developerRegisterMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestDeveloperRevokeMsg(t *testing.T) {
	cases := []struct {
		developerRevokeMsg DeveloperRevokeMsg
		expectError        sdk.Error
	}{
		{NewDeveloperRevokeMsg("user1"), nil},
		{NewDeveloperRevokeMsg(""), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.developerRevokeMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestGrantPermissionMsgMsg(t *testing.T) {
	cases := []struct {
		grantPermissionMsg GrantPermissionMsg
		expectError        sdk.Error
	}{
		{NewGrantPermissionMsg("user1", "app", 10, 1, types.MicropaymentPermission), nil},
		{NewGrantPermissionMsg("user1", "app", 10, 1, types.PostPermission), nil},
		{NewGrantPermissionMsg("user1", "app", 10, 1, types.MasterPermission), ErrGrantPermissionTooHigh()},
		{NewGrantPermissionMsg("user1", "app", 10, 1, types.TransactionPermission), ErrGrantPermissionTooHigh()},
		{NewGrantPermissionMsg("user1", "app", 10, 1, types.GrantMicropaymentPermission), ErrGrantPermissionTooHigh()},
		{NewGrantPermissionMsg("user1", "app", 10, 1, types.GrantPostPermission), ErrGrantPermissionTooHigh()},
		{NewGrantPermissionMsg("user1", "app", -1, 1, types.PostPermission), ErrInvalidValidityPeriod()},
		{NewGrantPermissionMsg("us", "app", 1, 1, types.PostPermission), ErrInvalidUsername()},
		{NewGrantPermissionMsg("user1", "ap", 1, 1, types.PostPermission), ErrInvalidUsername()},
		{NewGrantPermissionMsg("user1user1user1user1user1", "app", 1, 1, types.PostPermission), ErrInvalidUsername()},
		{NewGrantPermissionMsg("user1", "appappappappappappapp", 1, 1, types.PostPermission), ErrInvalidUsername()},
		{NewGrantPermissionMsg("user1", "app", 1, -1, types.PostPermission), ErrInvalidGrantTimes()},
	}

	for _, cs := range cases {
		result := cs.grantPermissionMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestRevokePermissionMsgMsg(t *testing.T) {
	cases := []struct {
		revokePermissionMsg RevokePermissionMsg
		expectError         sdk.Error
	}{
		{NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.MicropaymentPermission), nil},
		{NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission), nil},
		{NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.MasterPermission), ErrGrantPermissionTooHigh()},
		{NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.GrantMicropaymentPermission), ErrGrantPermissionTooHigh()},
		{NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.TransactionPermission), ErrGrantPermissionTooHigh()},
		{NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.GrantPostPermission), ErrGrantPermissionTooHigh()},
		{NewRevokePermissionMsg("us", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission), ErrInvalidUsername()},
		{NewRevokePermissionMsg("user1user1user1user1user1", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.revokePermissionMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              types.Msg
		expectPermission types.Permission
	}{
		"developer register msg": {
			NewDeveloperRegisterMsg("test", types.LNO("1")),
			types.TransactionPermission},
		"developer revoke msg": {
			NewDeveloperRevokeMsg("test"),
			types.TransactionPermission},
		"grant developer post permission msg": {
			NewGrantPermissionMsg("test", "app", 24*3600, 1, types.PostPermission),
			types.GrantPostPermission},
		"grant developer micropayment permission msg": {
			NewGrantPermissionMsg("test", "app", 24*3600, 1, types.MicropaymentPermission),
			types.GrantMicropaymentPermission},
		"revoke developer micropayment permission msg": {
			NewRevokePermissionMsg("test", crypto.GenPrivKeyEd25519().PubKey(), types.MicropaymentPermission),
			types.GrantMicropaymentPermission},
		"revoke developer post permission msg": {
			NewRevokePermissionMsg("test", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission),
			types.GrantPostPermission},
	}

	for testName, cs := range cases {
		permission := cs.msg.GetPermission()
		if cs.expectPermission != permission {
			t.Errorf(
				"%s: expect permission incorrect, expect %v, got %v",
				testName, cs.expectPermission, permission)
			return
		}
	}
}
