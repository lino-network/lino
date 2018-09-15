package vote

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStakeInMsg(t *testing.T) {
	testCases := []struct {
		testName      string
		StakeInMsg    StakeInMsg
		expectedError sdk.Error
	}{
		{
			testName:      "normal case",
			StakeInMsg:    NewStakeInMsg("user1", "1"),
			expectedError: nil,
		},
		{
			testName:      "invalid username",
			StakeInMsg:    NewStakeInMsg("", "1"),
			expectedError: ErrInvalidUsername(),
		},
		{
			testName:      "invalid deposit amount",
			StakeInMsg:    NewStakeInMsg("user1", "-1"),
			expectedError: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
	}

	for _, tc := range testCases {
		result := tc.StakeInMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, expect %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestStakeOutMsg(t *testing.T) {
	testCases := []struct {
		testName      string
		StakeOutMsg   StakeOutMsg
		expectedError sdk.Error
	}{
		{
			testName:      "normal case",
			StakeOutMsg:   NewStakeOutMsg("user1", "1"),
			expectedError: nil,
		},
		{
			testName:      "invalid username",
			StakeOutMsg:   NewStakeOutMsg("", "1"),
			expectedError: ErrInvalidUsername(),
		},
		{
			testName:      "invalid withdraw amount",
			StakeOutMsg:   NewStakeOutMsg("user1", "-1"),
			expectedError: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
	}

	for _, tc := range testCases {
		result := tc.StakeOutMsg.ValidateBasic()
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
			msg:                NewStakeInMsg("test", types.LNO("1")),
			expectedPermission: types.TransactionPermission,
		},
		{
			testName:           "vote withdraw",
			msg:                NewStakeOutMsg("test", types.LNO("1")),
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
	}

	for _, tc := range testCases {
		permission := tc.msg.GetPermission()
		if tc.expectedPermission != permission {
			t.Errorf("%s: diff permission, got %v, want %v", tc.testName, permission, tc.expectedPermission)
		}
	}
}

func TestGetSignBytes(t *testing.T) {
	testCases := []struct {
		testName string
		msg      types.Msg
	}{
		{
			testName: "vote deposit",
			msg:      NewStakeInMsg("test", types.LNO("1")),
		},
		{
			testName: "vote withdraw",
			msg:      NewStakeOutMsg("test", types.LNO("1")),
		},
		{
			testName: "delegate to voter",
			msg:      NewDelegateMsg("delegator", "voter", types.LNO("1")),
		},
		{
			testName: "delegate withdraw",
			msg:      NewDelegatorWithdrawMsg("delegator", "voter", types.LNO("1")),
		},
	}

	for _, tc := range testCases {
		require.NotPanics(t, func() { tc.msg.GetSignBytes() }, tc.testName)
	}
}

func TestGetSigners(t *testing.T) {
	testCases := []struct {
		testName      string
		msg           types.Msg
		expectSigners []types.AccountKey
	}{
		{
			testName:      "vote deposit",
			msg:           NewStakeInMsg("test", types.LNO("1")),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "vote withdraw",
			msg:           NewStakeOutMsg("test", types.LNO("1")),
			expectSigners: []types.AccountKey{"test"},
		},
		{
			testName:      "delegate to voter",
			msg:           NewDelegateMsg("delegator", "voter", types.LNO("1")),
			expectSigners: []types.AccountKey{"delegator"},
		},
		{
			testName:      "delegate withdraw",
			msg:           NewDelegatorWithdrawMsg("delegator", "voter", types.LNO("1")),
			expectSigners: []types.AccountKey{"delegator"},
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
