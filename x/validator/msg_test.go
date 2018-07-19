package validator

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/tendermint/crypto"
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

func TestMsgPermission(t *testing.T) {
	testCases := []struct {
		testName           string
		msg                types.Msg
		expectedPermission types.Permission
	}{
		{
			testName: "validator deposit msg",
			msg: NewValidatorDepositMsg(
				"test", types.LNO("1"), crypto.GenPrivKeyEd25519().PubKey(), "https://lino.network"),
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
