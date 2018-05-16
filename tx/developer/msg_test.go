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

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              sdk.Msg
		expectPermission types.Permission
	}{
		"developer register msg": {
			NewDeveloperRegisterMsg("test", types.LNO("1")),
			types.TransactionPermission},
		"developer revoke msg": {
			NewDeveloperRevokeMsg("test"),
			types.TransactionPermission},
		"grant developer msg": {
			NewGrantDeveloperMsg("test", "app", 24*3600, types.PostPermission),
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
