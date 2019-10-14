package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestValidatorRevokeMsg(t *testing.T) {
	testCases := []struct {
		testName           string
		validatorRevokeMsg ValidatorRevokeMsg
		expectedError      sdk.Error
	}{
		{
			testName:           "normal case",
			validatorRevokeMsg: NewValidatorRevokeMsg("user1"),
			expectedError:      nil,
		},
		{
			testName:           "invalid username",
			validatorRevokeMsg: NewValidatorRevokeMsg(""),
			expectedError:      ErrInvalidUsername(),
		},
	}

	for _, tc := range testCases {
		result := tc.validatorRevokeMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestValidatorRegisterMsg(t *testing.T) {
	testCases := []struct {
		testName             string
		validatorRegisterMsg ValidatorRegisterMsg
		expectedError        sdk.Error
	}{
		{
			testName:             "normal case",
			validatorRegisterMsg: NewValidatorRegisterMsg("user1", secp256k1.GenPrivKey().PubKey(), ""),
			expectedError:        nil,
		},
		{
			testName:             "invalid username",
			validatorRegisterMsg: NewValidatorRegisterMsg("", secp256k1.GenPrivKey().PubKey(), ""),
			expectedError:        ErrInvalidUsername(),
		},
		{
			testName: "invalid Website",
			validatorRegisterMsg: NewValidatorRegisterMsg(
				"user", secp256k1.GenPrivKey().PubKey(), string(make([]byte, types.MaximumLinkURL+1))),
			expectedError: ErrInvalidWebsite(),
		},
	}

	for _, tc := range testCases {
		result := tc.validatorRegisterMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestValidatorUpdateMsg(t *testing.T) {
	testCases := []struct {
		testName           string
		validatorUpdateMsg ValidatorUpdateMsg
		expectedError      sdk.Error
	}{
		{
			testName:           "normal case",
			validatorUpdateMsg: NewValidatorUpdateMsg("user1", "123123"),
			expectedError:      nil,
		},
	}

	for _, tc := range testCases {
		result := tc.validatorUpdateMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestVoteValidatorMsg(t *testing.T) {
	testCases := []struct {
		testName      string
		msg           VoteValidatorMsg
		expectedError sdk.Error
	}{
		{
			testName:      "normal case",
			msg:           NewVoteValidatorMsg("user1", []string{"val1"}),
			expectedError: nil,
		},
		{
			testName:      "invalid username",
			msg:           NewVoteValidatorMsg("", []string{"val1"}),
			expectedError: ErrInvalidUsername(),
		},
		{
			testName:      "invalid voted validators",
			msg:           NewVoteValidatorMsg("user1", []string{}),
			expectedError: ErrInvalidVotedValidators(),
		},
	}

	for _, tc := range testCases {
		result := tc.msg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestMsgPermission(t *testing.T) {
	testCases := []struct {
		testName           string
		msg                types.Msg
		expectedPermission types.Permission
	}{
		{
			testName: "validator deposit msg",
			msg: NewValidatorRegisterMsg(
				"test", secp256k1.GenPrivKey().PubKey(), "https://lino.network"),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "validator revoke msg",
			msg:                NewValidatorRevokeMsg("test"),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "vote validator msg",
			msg:                NewVoteValidatorMsg("test", []string{"val1"}),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "uddate validator msg",
			msg:                NewValidatorUpdateMsg("test", "asd"),
			expectedPermission: types.TransactionPermission,
		},
	}

	for _, tc := range testCases {
		permission := tc.msg.GetPermission()
		if tc.expectedPermission != permission {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, permission, tc.expectedPermission)
			return
		}
	}
}

func TestGetSignBytes(t *testing.T) {
	testCases := []struct {
		testName string
		msg      types.Msg
	}{
		{
			testName: "validator register msg",
			msg: NewValidatorRegisterMsg(
				"test", secp256k1.GenPrivKey().PubKey(), "https://lino.network"),
		},
		{
			testName: "validator revoke msg",
			msg:      NewValidatorRevokeMsg("test"),
		},
		{
			testName: "vote validator msg",
			msg:      NewVoteValidatorMsg("test", []string{"val1", "val2"}),
		},
		{
			testName: "validator update msg",
			msg: NewValidatorUpdateMsg(
				"test", "https://lino.network"),
		},
	}

	for testName, tc := range testCases {
		require.NotPanics(t, func() { tc.msg.GetSignBytes() }, testName)
	}
}

func TestGetSigners(t *testing.T) {
	testCases := []struct {
		testName      string
		msg           types.Msg
		expectSigners []types.AccountKey
	}{
		{
			testName: "validator register msg",
			msg: NewValidatorRegisterMsg(
				"test", secp256k1.GenPrivKey().PubKey(), "https://lino.network"),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "validator revoke msg",
			msg:           NewValidatorRevokeMsg("test"),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "vote validator msg",
			msg:           NewVoteValidatorMsg("test", []string{"val1", "val2"}),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "validator update msg",
			msg:           NewValidatorUpdateMsg("test", "https://lino.network"),
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
