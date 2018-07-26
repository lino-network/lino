package validator

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

func TestValidatorWithdrawMsg(t *testing.T) {
	testCases := []struct {
		testName             string
		validatorWithdrawMsg ValidatorWithdrawMsg
		expectedError        sdk.Error
	}{
		{
			testName:             "normal case",
			validatorWithdrawMsg: NewValidatorWithdrawMsg("user1", "1"),
			expectedError:        nil,
		},
		{
			testName:             "invalid username",
			validatorWithdrawMsg: NewValidatorWithdrawMsg("", "1"),
			expectedError:        ErrInvalidUsername(),
		},
	}

	for _, tc := range testCases {
		result := tc.validatorWithdrawMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestValidatorDepositMsg(t *testing.T) {
	testCases := []struct {
		testName            string
		validatorDepositMsg ValidatorDepositMsg
		expectedError       sdk.Error
	}{
		{
			testName:            "normal case",
			validatorDepositMsg: NewValidatorDepositMsg("user1", "1", secp256k1.GenPrivKey().PubKey(), ""),
			expectedError:       nil,
		},
		{
			testName:            "invalid username",
			validatorDepositMsg: NewValidatorDepositMsg("", "1", secp256k1.GenPrivKey().PubKey(), ""),
			expectedError:       ErrInvalidUsername(),
		},
		{
			testName:            "invalid LNO",
			validatorDepositMsg: NewValidatorDepositMsg("user", ".", secp256k1.GenPrivKey().PubKey(), ""),
			expectedError:       types.ErrInvalidCoins("Illegal LNO"),
		},
		{
			testName: "invalid Website",
			validatorDepositMsg: NewValidatorDepositMsg(
				"user", "1", secp256k1.GenPrivKey().PubKey(), string(make([]byte, types.MaximumLinkURL+1))),
			expectedError: ErrInvalidWebsite(),
		},
	}

	for _, tc := range testCases {
		result := tc.validatorDepositMsg.ValidateBasic()
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
			msg: NewValidatorDepositMsg(
				"test", types.LNO("1"), secp256k1.GenPrivKey().PubKey(), "https://lino.network"),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "validator withdraw msg",
			msg:                NewValidatorWithdrawMsg("test", types.LNO("1")),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "validator revoke msg",
			msg:                NewValidatorRevokeMsg("test"),
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
			testName: "validator deposit msg",
			msg: NewValidatorDepositMsg(
				"test", types.LNO("1"), secp256k1.GenPrivKey().PubKey(), "https://lino.network"),
		},
		{
			testName: "validator withdraw msg",
			msg:      NewValidatorWithdrawMsg("test", types.LNO("1")),
		},
		{
			testName: "validator revoke msg",
			msg:      NewValidatorRevokeMsg("test"),
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
			testName: "validator deposit msg",
			msg: NewValidatorDepositMsg(
				"test", types.LNO("1"), secp256k1.GenPrivKey().PubKey(), "https://lino.network"),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "validator withdraw msg",
			msg:           NewValidatorWithdrawMsg("test", types.LNO("1")),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "validator revoke msg",
			msg:           NewValidatorRevokeMsg("test"),
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
