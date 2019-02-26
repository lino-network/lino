package developer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

var (
	// len of 1000
	maxLengthUTF8Str = `
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌123`

	// len of 1001
	tooLongUTF8Str = `
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌
	12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌12345 67890 你好👌1234`
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
			testName: "utf8 description",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "10", "", maxLengthUTF8Str, ""),
			expectError: nil,
		},
		{
			testName: "utf8 app metadata",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "10", "", "", "a:bb,你:👌 。"),
			expectError: nil,
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
		{
			testName: "website is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "10", string(make([]byte, types.MaximumLengthOfDeveloperWebsite+1)), "", ""),
			expectError: ErrInvalidWebsite(),
		},
		{
			testName: "description is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "10", "", string(make([]byte, types.MaximumLengthOfDeveloperDesctiption+1)), ""),
			expectError: ErrInvalidDescription(),
		},
		{
			testName: "utf8 description is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "10", "", tooLongUTF8Str, ""),
			expectError: ErrInvalidDescription(),
		},
		{
			testName: "app metadata is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "10", "", "", string(make([]byte, types.MaximumLengthOfAppMetadata+1))),
			expectError: ErrInvalidAppMetadata(),
		},
		{
			testName: "utf8 app metadata is too long",
			developerRegisterMsg: NewDeveloperRegisterMsg(
				"user1", "10", "", "", tooLongUTF8Str),
			expectError: ErrInvalidAppMetadata(),
		},
	}

	for _, tc := range testCases {
		result := tc.developerRegisterMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectError)
		}
	}
}
func TestDeveloperUpdateMsg(t *testing.T) {
	testCases := []struct {
		testName           string
		developerUpdateMsg DeveloperUpdateMsg
		expectError        sdk.Error
	}{
		{
			testName:           "normal case",
			developerUpdateMsg: NewDeveloperUpdateMsg("user1", "", "", ""),
			expectError:        nil,
		},
		{
			testName: "utf8 description",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", maxLengthUTF8Str, ""),
			expectError: nil,
		},
		{
			testName: "uft8 app metadata",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", "", maxLengthUTF8Str),
			expectError: nil,
		},
		{
			testName:           "invalid username",
			developerUpdateMsg: NewDeveloperUpdateMsg("", "", "", ""),
			expectError:        ErrInvalidUsername(),
		},
		{
			testName: "website is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", string(make([]byte, types.MaximumLengthOfDeveloperWebsite+1)), "", ""),
			expectError: ErrInvalidWebsite(),
		},
		{
			testName: "description is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", string(make([]byte, types.MaximumLengthOfDeveloperDesctiption+1)), ""),
			expectError: ErrInvalidDescription(),
		},
		{
			testName: "utf8 description is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", tooLongUTF8Str, ""),
			expectError: ErrInvalidDescription(),
		},
		{
			testName: "app metadata is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", "", string(make([]byte, types.MaximumLengthOfAppMetadata+1))),
			expectError: ErrInvalidAppMetadata(),
		},
		{
			testName: "utf8 app metadata is too long",
			developerUpdateMsg: NewDeveloperUpdateMsg(
				"user1", "", "", tooLongUTF8Str),
			expectError: ErrInvalidAppMetadata(),
		},
	}

	for _, tc := range testCases {
		result := tc.developerUpdateMsg.ValidateBasic()
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
			testName:           "app permission",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.AppPermission, "1"),
			expectError:        nil,
		},
		{
			testName:           "reset permission is too high",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.ResetPermission, "1"),
			expectError:        ErrGrantPermissionTooHigh(),
		},
		{
			testName:           "transaction permission is too high",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.TransactionPermission, "1"),
			expectError:        ErrGrantPermissionTooHigh(),
		},
		{
			testName:           "grant app permission is too high",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.GrantAppPermission, "1"),
			expectError:        ErrGrantPermissionTooHigh(),
		},
		{
			testName:           "grant pre-auth permission with invalid coin",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.PreAuthorizationPermission, "0"),
			expectError:        types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName:           "grant both app and pre-auth permission with invalid coin",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", 10, types.AppAndPreAuthorizationPermission, "0"),
			expectError:        types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
		{
			testName:           "invalid validity period",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "app", -1, types.AppPermission, "0"),
			expectError:        ErrInvalidValidityPeriod(),
		},
		{
			testName:           "invalid username",
			grantPermissionMsg: NewGrantPermissionMsg("us", "app", 1, types.AppPermission, "0"),
			expectError:        ErrInvalidUsername(),
		},
		{
			testName:           "invalid authenticate app, app name is too short",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "ap", 1, types.AppPermission, "0"),
			expectError:        ErrInvalidAuthorizedApp(),
		},
		{
			testName:           "invalid username",
			grantPermissionMsg: NewGrantPermissionMsg("user1user1user1user1user1", "app", 1, types.AppPermission, "0"),
			expectError:        ErrInvalidUsername(),
		},
		{
			testName:           "invalid authenticate app, app name is too long",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "appappappappappappapp", 1, types.AppPermission, "0"),
			expectError:        ErrInvalidAuthorizedApp(),
		},
		{
			testName:           "invalid authenticate app, app name is too long",
			grantPermissionMsg: NewGrantPermissionMsg("user1", "appappappappappappapp", 1, types.AppPermission, "0"),
			expectError:        ErrInvalidAuthorizedApp(),
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
			testName:            "revoke permission",
			revokePermissionMsg: NewRevokePermissionMsg("user1", secp256k1.GenPrivKey().PubKey()),
			expectError:         nil,
		},
		{
			testName:            "username is too short",
			revokePermissionMsg: NewRevokePermissionMsg("us", secp256k1.GenPrivKey().PubKey()),
			expectError:         ErrInvalidUsername(),
		},
		{
			testName:            "username is too long",
			revokePermissionMsg: NewRevokePermissionMsg("user1user1user1user1user1", secp256k1.GenPrivKey().PubKey()),
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
func TestPreAuthorizationMsgMsg(t *testing.T) {
	testCases := []struct {
		testName            string
		preAuthorizationMsg PreAuthorizationMsg
		expectError         sdk.Error
	}{
		{
			testName:            "normal preauthorization",
			preAuthorizationMsg: NewPreAuthorizationMsg("user1", "app", 1000, "1"),
			expectError:         nil,
		},
		{
			testName:            "invalid validity second",
			preAuthorizationMsg: NewPreAuthorizationMsg("user1", "app", -1, "1"),
			expectError:         ErrInvalidValidityPeriod(),
		},
		{
			testName:            "illegal LNO",
			preAuthorizationMsg: NewPreAuthorizationMsg("user1", "app", 1000, "*"),
			expectError:         types.ErrInvalidCoins("Illegal LNO"),
		},
		{
			testName:            "username is too short",
			preAuthorizationMsg: NewPreAuthorizationMsg("us", "app", 1000, "1"),
			expectError:         ErrInvalidUsername(),
		},
		{
			testName:            "username is too long",
			preAuthorizationMsg: NewPreAuthorizationMsg("user1user1user1user1user1", "app", 1000, "1"),
			expectError:         ErrInvalidUsername(),
		},
		{
			testName:            "app name is too short",
			preAuthorizationMsg: NewPreAuthorizationMsg("user1", "ap", 1000, "1"),
			expectError:         ErrInvalidAuthorizedApp(),
		},
		{
			testName:            "app name is too long",
			preAuthorizationMsg: NewPreAuthorizationMsg("user1", "appappappappappappappapp", 1000, "1"),
			expectError:         ErrInvalidAuthorizedApp(),
		},
	}

	for _, tc := range testCases {
		result := tc.preAuthorizationMsg.ValidateBasic()
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
			testName:         "developer register msg",
			msg:              NewDeveloperUpdateMsg("test", "", "", ""),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "developer revoke msg",
			msg:              NewDeveloperRevokeMsg("test"),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "grant developer app permission msg",
			msg:              NewGrantPermissionMsg("test", "app", 24*3600, types.AppPermission, "0"),
			expectPermission: types.GrantAppPermission,
		},
		{
			testName:         "grant developer pre-auth permission msg",
			msg:              NewGrantPermissionMsg("test", "app", 24*3600, types.PreAuthorizationPermission, "1"),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "grant developer both app and pre-auth permission msg",
			msg:              NewGrantPermissionMsg("test", "app", 24*3600, types.AppAndPreAuthorizationPermission, "1"),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "grant developer app permission msg",
			msg:              NewGrantPermissionMsg("test", "app", 24*3600, types.AppPermission, "0"),
			expectPermission: types.GrantAppPermission,
		},
		{
			testName:         "revoke developer app permission msg",
			msg:              NewRevokePermissionMsg("test", secp256k1.GenPrivKey().PubKey()),
			expectPermission: types.TransactionPermission,
		},
		{
			testName:         "pre authorization msg",
			msg:              NewPreAuthorizationMsg("test", "app", 1000, "1"),
			expectPermission: types.TransactionPermission,
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
			testName:      "developer update msg",
			msg:           NewDeveloperUpdateMsg("test", "", "", ""),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "developer revoke msg",
			msg:           NewDeveloperRevokeMsg("test"),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "grant developer app permission msg",
			msg:           NewGrantPermissionMsg("test", "app", 24*3600, types.AppPermission, "0"),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "revoke developer post permission msg",
			msg:           NewRevokePermissionMsg("test", secp256k1.GenPrivKey().PubKey()),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "pre authorization msg",
			msg:           NewPreAuthorizationMsg("test", "app", 1000, "1"),
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
			testName: "developer register msg",
			msg:      NewDeveloperUpdateMsg("test", "", "", ""),
		},
		{
			testName: "developer revoke msg",
			msg:      NewDeveloperRevokeMsg("test"),
		},
		{
			testName: "grant developer app permission msg",
			msg:      NewGrantPermissionMsg("test", "app", 24*3600, types.AppPermission, "0"),
		},
		{
			testName: "revoke developer post permission msg",
			msg:      NewRevokePermissionMsg("test", secp256k1.GenPrivKey().PubKey()),
		},
		{
			testName: "preauth msg",
			msg:      NewPreAuthorizationMsg("test", "app", 1000, "1"),
		},
	}

	for _, tc := range testCases {
		require.NotPanics(t, func() { tc.msg.GetSignBytes() }, tc.testName)
	}
}
