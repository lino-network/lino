package developer

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
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

func TestGrantDeveloperMsgMsg(t *testing.T) {
	cases := []struct {
		grantDeveloperMsg GrantDeveloperMsg
		expectError       sdk.Error
	}{
		{NewGrantDeveloperMsg("user1", "app", 10, types.MicropaymentPermission), nil},
		{NewGrantDeveloperMsg("user1", "app", 10, types.PostPermission), nil},
		{NewGrantDeveloperMsg("user1", "app", 10, types.MasterPermission), ErrGrantPermissionTooHigh()},
		{NewGrantDeveloperMsg("user1", "app", 10, types.TransactionPermission), ErrGrantPermissionTooHigh()},
		{NewGrantDeveloperMsg("user1", "app", -1, types.PostPermission), ErrInvalidValidityPeriod()},
		{NewGrantDeveloperMsg("us", "app", 1, types.PostPermission), ErrInvalidUsername()},
		{NewGrantDeveloperMsg("user1", "ap", 1, types.PostPermission), ErrInvalidUsername()},
		{NewGrantDeveloperMsg("user1user1user1user1user1", "app", 1, types.PostPermission), ErrInvalidUsername()},
		{NewGrantDeveloperMsg("user1", "appappappappappappapp", 1, types.PostPermission), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.grantDeveloperMsg.ValidateBasic()
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
		"grant developer post msg": {
			NewGrantDeveloperMsg("test", "app", 24*3600, types.PostPermission),
			types.PostPermission},
		"grant developer micropayment msg": {
			NewGrantDeveloperMsg("test", "app", 24*3600, types.MicropaymentPermission),
			types.MicropaymentPermission},
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
