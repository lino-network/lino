package vote

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestVoterDepositMsg(t *testing.T) {
	cases := []struct {
		voterDepositMsg VoterDepositMsg
		expectError     sdk.Error
	}{
		{NewVoterDepositMsg("user1", "1"), nil},
		{NewVoterDepositMsg("", "1"), ErrInvalidUsername()},
		{NewVoterDepositMsg("user1", "-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.voterDepositMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestVoterDepositMsgPermission(t *testing.T) {
	msg := NewVoterDepositMsg("user1", "1")
	permissionLevel := msg.Get(types.PermissionLevel)
	permission, ok := permissionLevel.(types.Permission)
	assert.Equal(t, ok, true)
	assert.Equal(t, permission, types.TransactionPermission)
}

func TestVoterWithdrawMsg(t *testing.T) {
	cases := []struct {
		voterWithdrawMsg VoterWithdrawMsg
		expectError      sdk.Error
	}{
		{NewVoterWithdrawMsg("user1", "1"), nil},
		{NewVoterWithdrawMsg("", "1"), ErrInvalidUsername()},
		{NewVoterWithdrawMsg("user1", "-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.voterWithdrawMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestVoterRevokeMsg(t *testing.T) {
	cases := []struct {
		voterRevokeMsg VoterRevokeMsg
		expectError    sdk.Error
	}{
		{NewVoterRevokeMsg("user1"), nil},
		{NewVoterRevokeMsg(""), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.voterRevokeMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestDelegateMsg(t *testing.T) {
	cases := []struct {
		delegateMsg DelegateMsg
		expectError sdk.Error
	}{
		{NewDelegateMsg("user1", "user2", "1"), nil},
		{NewDelegateMsg("", "user2", "1"), ErrInvalidUsername()},
		{NewDelegateMsg("user1", "", "1"), ErrInvalidUsername()},
		{NewDelegateMsg("", "", "1"), ErrInvalidUsername()},
		{NewDelegateMsg("user1", "user2", "-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.delegateMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestDelegateMsgPermission(t *testing.T) {
	msg := NewDelegateMsg("user1", "user2", "1")
	permissionLevel := msg.Get(types.PermissionLevel)
	permission, ok := permissionLevel.(types.Permission)
	assert.Equal(t, ok, true)
	assert.Equal(t, permission, types.TransactionPermission)
}

func TestRevokeDelegationMsg(t *testing.T) {
	cases := []struct {
		revokeDelegationMsg RevokeDelegationMsg
		expectError         sdk.Error
	}{
		{NewRevokeDelegationMsg("user1", "user2"), nil},
		{NewRevokeDelegationMsg("", "user2"), ErrInvalidUsername()},
		{NewRevokeDelegationMsg("user1", ""), ErrInvalidUsername()},
		{NewRevokeDelegationMsg("", ""), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.revokeDelegationMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestDelegatorWithdrawMsg(t *testing.T) {
	cases := []struct {
		delegatorWithdrawMsg DelegatorWithdrawMsg
		expectError          sdk.Error
	}{
		{NewDelegatorWithdrawMsg("user1", "user2", "1"), nil},
		{NewDelegatorWithdrawMsg("", "", "1"), ErrInvalidUsername()},
		{NewDelegatorWithdrawMsg("user1", "user2", "-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.delegatorWithdrawMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              sdk.Msg
		expectPermission types.Permission
	}{
		"vote deposit": {
			NewVoterDepositMsg("test", types.LNO("1")),
			types.TransactionPermission},
		"vote withdraw": {
			NewVoterWithdrawMsg("test", types.LNO("1")),
			types.TransactionPermission},
		"vote revoke": {
			NewVoterRevokeMsg("test"),
			types.TransactionPermission},
		"delegate to voter": {
			NewDelegateMsg("delegator", "voter", types.LNO("1")),
			types.TransactionPermission},
		"delegate withdraw": {
			NewDelegatorWithdrawMsg("delegator", "voter", types.LNO("1")),
			types.TransactionPermission},
		"revoke delegation": {
			NewRevokeDelegationMsg("delegator", "voter"),
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
