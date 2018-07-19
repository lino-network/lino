package vote

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestVoterDepositMsg(t *testing.T) {
	testCases := []struct {
		testName        string
		voterDepositMsg VoterDepositMsg
		expectedError   sdk.Error
	}{
		{
			testName:        "normal case",
			voterDepositMsg: NewVoterDepositMsg("user1", "1"),
			expectedError:   nil,
		},
		{
			testName:        "invalid username",
			voterDepositMsg: NewVoterDepositMsg("", "1"),
			expectedError:   ErrInvalidUsername(),
		},
		{
			testName:        "invalid deposit amount",
			voterDepositMsg: NewVoterDepositMsg("user1", "-1"),
			expectedError:   types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
	}

	for _, tc := range testCases {
		result := tc.voterDepositMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, expect %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestVoterWithdrawMsg(t *testing.T) {
	testCases := []struct {
		testName         string
		voterWithdrawMsg VoterWithdrawMsg
		expectedError    sdk.Error
	}{
		{
			testName:         "normal case",
			voterWithdrawMsg: NewVoterWithdrawMsg("user1", "1"),
			expectedError:    nil,
		},
		{
			testName:         "invalid username",
			voterWithdrawMsg: NewVoterWithdrawMsg("", "1"),
			expectedError:    ErrInvalidUsername(),
		},
		{
			testName:         "invalid withdraw amount",
			voterWithdrawMsg: NewVoterWithdrawMsg("user1", "-1"),
			expectedError:    types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
	}

	for _, tc := range testCases {
		result := tc.voterWithdrawMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, expect %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestVoterRevokeMsg(t *testing.T) {
	testCases := []struct {
		testName       string
		voterRevokeMsg VoterRevokeMsg
		expectedError  sdk.Error
	}{
		{
			testName:       "normal case",
			voterRevokeMsg: NewVoterRevokeMsg("user1"),
			expectedError:  nil,
		},
		{
			testName:       "invalid username",
			voterRevokeMsg: NewVoterRevokeMsg(""),
			expectedError:  ErrInvalidUsername(),
		},
	}

	for _, tc := range testCases {
		result := tc.voterRevokeMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, expect %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestDelegateMsg(t *testing.T) {
	testCases := []struct {
		testName      string
		delegateMsg   DelegateMsg
		expectedError sdk.Error
	}{
		{
			testName:      "normal case",
			delegateMsg:   NewDelegateMsg("user1", "user2", "1"),
			expectedError: nil,
		},
		{
			testName:      "invalid delegator",
			delegateMsg:   NewDelegateMsg("", "user2", "1"),
			expectedError: ErrInvalidUsername(),
		},
		{
			testName:      "invalid voter",
			delegateMsg:   NewDelegateMsg("user1", "", "1"),
			expectedError: ErrInvalidUsername(),
		},
		{
			testName:      "both delegator and voter are invalid",
			delegateMsg:   NewDelegateMsg("", "", "1"),
			expectedError: ErrInvalidUsername(),
		},
		{
			testName:      "invalid delegated coin",
			delegateMsg:   NewDelegateMsg("user1", "user2", "-1"),
			expectedError: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
	}

	for _, tc := range testCases {
		result := tc.delegateMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, expect %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestRevokeDelegationMsg(t *testing.T) {
	testCases := []struct {
		testName            string
		revokeDelegationMsg RevokeDelegationMsg
		expectedError       sdk.Error
	}{
		{
			testName:            "normal case",
			revokeDelegationMsg: NewRevokeDelegationMsg("user1", "user2"),
			expectedError:       nil,
		},
		{
			testName:            "invalid delegator",
			revokeDelegationMsg: NewRevokeDelegationMsg("", "user2"),
			expectedError:       ErrInvalidUsername(),
		},
		{
			testName:            "invalid voter",
			revokeDelegationMsg: NewRevokeDelegationMsg("user1", ""),
			expectedError:       ErrInvalidUsername(),
		},
		{
			testName:            "both delegator and voter are invalid",
			revokeDelegationMsg: NewRevokeDelegationMsg("", ""),
			expectedError:       ErrInvalidUsername(),
		},
	}

	for _, tc := range testCases {
		result := tc.revokeDelegationMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, expect %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestDelegatorWithdrawMsg(t *testing.T) {
	testCases := []struct {
		testName             string
		delegatorWithdrawMsg DelegatorWithdrawMsg
		expectedError        sdk.Error
	}{
		{
			testName:             "normal case",
			delegatorWithdrawMsg: NewDelegatorWithdrawMsg("user1", "user2", "1"),
			expectedError:        nil,
		},
		{
			testName:             "invalid username",
			delegatorWithdrawMsg: NewDelegatorWithdrawMsg("", "", "1"),
			expectedError:        ErrInvalidUsername(),
		},
		{
			testName:             "invalid withdraw amount",
			delegatorWithdrawMsg: NewDelegatorWithdrawMsg("user1", "user2", "-1"),
			expectedError:        types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
	}

	for _, tc := range testCases {
		result := tc.delegatorWithdrawMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, expect %v", tc.testName, result, tc.expectedError)
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
			testName:           "vote deposit",
			msg:                NewVoterDepositMsg("test", types.LNO("1")),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "vote withdraw",
			msg:                NewVoterWithdrawMsg("test", types.LNO("1")),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "vote revoke",
			msg:                NewVoterRevokeMsg("test"),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "delegate to voter",
			msg:                NewDelegateMsg("delegator", "voter", types.LNO("1")),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "delegate withdraw",
			msg:                NewDelegatorWithdrawMsg("delegator", "voter", types.LNO("1")),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "revoke delegation",
			msg:                NewRevokeDelegationMsg("delegator", "voter"),
			expectedPermission: types.TransactionPermission,
		},
	}

	for testName, tc := range testCases {
		permission := tc.msg.GetPermission()
		if tc.expectedPermission != permission {
			t.Errorf("%s: diff permission, got %v, want %v", testName, permission, tc.expectedPermission)
		}
	}
}
