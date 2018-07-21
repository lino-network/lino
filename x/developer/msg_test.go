package developer

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/tendermint/crypto"
)

func TestDeveloperRegisterMsg(t *testing.T) {
	testCases := []struct {
		testName             string
		developerRegisterMsg DeveloperRegisterMsg
		expectError          sdk.Error
	}{
		{
			testName:             "normal case",
			developerRegisterMsg: NewDeveloperRegisterMsg("user1", "10", "", "", ""),
			expectError:          nil,
		},
		{
			testName:             "invalid username",
			developerRegisterMsg: NewDeveloperRegisterMsg("", "10", "", "", ""),
			expectError:          ErrInvalidUsername(),
		},
		{
			testName:             "invalid coins",
			developerRegisterMsg: NewDeveloperRegisterMsg("user1", "-1", "", "", ""),
			expectError:          types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
	}

	for _, tc := range testCases {
		result := tc.developerRegisterMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectError)
		}
	}
}

func TestDeveloperRevokeMsg(t *testing.T) {
	testCases := []struct {
		testName           string
		developerRevokeMsg DeveloperRevokeMsg
		expectError        sdk.Error
	}{
		{
			testName:           "normal case",
			developerRevokeMsg: NewDeveloperRevokeMsg("user1"),
			expectError:        nil,
		},
		{
			testName:           "invalid username",
			developerRevokeMsg: NewDeveloperRevokeMsg(""),
			expectError:        ErrInvalidUsername(),
		},
	}

	for _, tc := range testCases {
		result := tc.developerRevokeMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectError)
		}
	}
}

func TestGrantPermissionMsgMsg(t *testing.T) {
	testCases := []struct {
		testName           string
		grantPermissionMsg GrantPermissionMsg
		expectError        sdk.Error
	}{
		{
			testName:           "post permission",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.PostPermission),
			expectError:        nil,
		},
		{
			testName:           "reset permission is too high",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.ResetPermission),
			expectError:        ErrGrantPermissionTooHigh(),
		},
		{
			testName:           "transaction permission is too high",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.TransactionPermission),
			expectError:        ErrGrantPermissionTooHigh(),
		},
		{
			testName:           "grant post permission is too high",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.GrantPostPermission),
			expectError:        ErrGrantPermissionTooHigh(),
		},
		{
			testName:           "invalid validity period",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", -1, types.PostPermission),
			expectError:        ErrInvalidValidityPeriod(),
		},
		{
			testName:           "invalid username",
			grantPermissionMsg: NewGrantPermissionMsg("us", "app", 1, types.PostPermission),
			expectError:        ErrInvalidUsername(),
		},
		{
			testName:           "invalid authenticate app, app name is too short",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "ap", 1, types.PostPermission),
			expectError:        ErrInvalidAuthenticateApp(),
		},
		{
			testName:           "invalid username",
			grantPermissionMsg: NewGrantPermissionMsg("user1user1user1user1user1", "app", 1, types.PostPermission),
			expectError:        ErrInvalidUsername(),
		},
		{
			testName:           "invalid authenticate app, app name is too long",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "appappappappappappapp", 1, types.PostPermission),
			expectError:        ErrInvalidAuthenticateApp(),
		},
	}

	for _, tc := range testCases {
		result := tc.grantPermissionMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectError)
		}
	}
}

func TestRevokePermissionMsgMsg(t *testing.T) {
	testCases := []struct {
		testName            string
		revokePermissionMsg RevokePermissionMsg
		expectError         sdk.Error
	}{
		{
			testName:            "revoke post permission",
			revokePermissionMsg: NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission),
			expectError:         nil,
		},
		{
			testName:            "reset permission is too high",
			revokePermissionMsg: NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.ResetPermission),
			expectError:         ErrGrantPermissionTooHigh(),
		},
		{
			testName:            "post permission is too high",
			revokePermissionMsg: NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.GrantPostPermission),
			expectError:         ErrGrantPermissionTooHigh(),
		},
		{
			testName:            "transaction permission is too high",
			revokePermissionMsg: NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.TransactionPermission),
			expectError:         ErrGrantPermissionTooHigh(),
		},
		{
			testName:            "grant post permission is too high",
			revokePermissionMsg: NewRevokePermissionMsg("user1", crypto.GenPrivKeyEd25519().PubKey(), types.GrantPostPermission),
			expectError:         ErrGrantPermissionTooHigh(),
		},
		{
			testName:            "username is too short",
			revokePermissionMsg: NewRevokePermissionMsg("us", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission),
			expectError:         ErrInvalidUsername(),
		},
		{
			testName:            "username is too long",
			revokePermissionMsg: NewRevokePermissionMsg("user1user1user1user1user1", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission),
			expectError:         ErrInvalidUsername(),
		},
	}

	for _, tc := range testCases {
		result := tc.revokePermissionMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectError)
		}
	}
}

func TestMsgPermission(t *testing.T) {
	testCases := []struct {
		testName         string
		msg              types.Msg
		expectPermission types.Permission
	}{
		{
			testName:         "developer register msg",
			msg:              NewDeveloperRegisterMsg("test", types.LNO("1"), "", "", ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "developer revoke msg",
			msg:              NewDeveloperRevokeMsg("test"),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "grant developer post permission msg",
			msg:              NewGrantPermissionMsg("test", "app", 24*3600, types.PostPermission),
			expectPermission: types.GrantPostPermission,
		},
		{
			testName:         "revoke developer post permission msg",
			msg:              NewRevokePermissionMsg("test", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission),
			expectPermission: types.GrantPostPermission,
		},
	}

	for _, tc := range testCases {
		permission := tc.msg.GetPermission()
		if tc.expectPermission != permission {
			t.Errorf(
				"%s: diff permission, got %v, want %v", tc.testName, permission, tc.expectPermission)
			return
		}
	}
}

func TestGetSigners(t *testing.T) {
	testCases := []struct {
		testName      string
		msg           types.Msg
		expectSigners []types.AccountKey
	}{
		{
			testName:      "developer register msg",
			msg:           NewDeveloperRegisterMsg("test", types.LNO("1"), "", "", ""),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "developer revoke msg",
			msg:           NewDeveloperRevokeMsg("test"),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "grant developer post permission msg",
			msg:           NewGrantPermissionMsg("test", "app", 24*3600, types.PostPermission),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "revoke developer post permission msg",
			msg:           NewRevokePermissionMsg("test", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission),
			expectSigners: []types.AccountKey{"test"},
		},
	}

	for _, tc := range testCases {
		if len(tc.msg.GetSigners()) != len(tc.expectSigners) {
			t.Errorf("%s: expect number of signers wrong, got %v, want %v", tc.testName, len(tc.msg.GetSigners()), len(tc.expectSigners))
			return
		}
		for i, signer := range tc.msg.GetSigners() {
			if types.AccountKey(signer) != tc.expectSigners[i] {
				t.Errorf("%s: expect signer wrong, got %v, want %v", tc.testName, types.AccountKey(signer), tc.expectSigners[i])
				return
			}
		}
	}
}

func TestGetSignBytes(t *testing.T) {
	testCases := []struct {
		testName string
		msg      types.Msg
	}{
		{
			testName: "developer register msg",
			msg:      NewDeveloperRegisterMsg("test", types.LNO("1"), "", "", ""),
		},
		{
			testName: "developer revoke msg",
			msg:      NewDeveloperRevokeMsg("test"),
		},
		{
			testName: "grant developer post permission msg",
			msg:      NewGrantPermissionMsg("test", "app", 24*3600, types.PostPermission),
		},
		{
			testName: "revoke developer post permission msg",
			msg:      NewRevokePermissionMsg("test", crypto.GenPrivKeyEd25519().PubKey(), types.PostPermission),
		},
	}

	for _, tc := range testCases {
		require.NotPanics(t, func() { tc.msg.GetSignBytes() }, tc.testName)
	}
}
