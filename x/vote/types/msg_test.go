package types

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

func TestStakeInForMsg(t *testing.T) {
	testCases := []struct {
		testName      string
		StakeInForMsg StakeInForMsg
		expectedError sdk.Error
	}{
		{
			testName:      "normal case",
			StakeInForMsg: NewStakeInForMsg("user1", "user2", "1"),
			expectedError: nil,
		},
		{
			testName:      "invalid username1",
			StakeInForMsg: NewStakeInForMsg("", "user2", "1"),
			expectedError: ErrInvalidUsername(),
		},
		{
			testName:      "invalid username2",
			StakeInForMsg: NewStakeInForMsg("user1", "", "1"),
			expectedError: ErrInvalidUsername(),
		},
		{
			testName:      "invalid deposit amount",
			StakeInForMsg: NewStakeInForMsg("user1", "user2", "-1"),
			expectedError: types.ErrInvalidCoins("LNO can't be less than lower bound"),
		},
	}

	for _, tc := range testCases {
		result := tc.StakeInForMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectedError) {
			t.Errorf("%s: diff result, got %v, expect %v", tc.testName, result, tc.expectedError)
		}
	}
}

func TestClaimInterestMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      ClaimInterestMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg:      NewClaimInterestMsg("test"),
			wantCode: sdk.CodeOK,
		},
		"invalid claim interest - Username is too short": {
			msg:      NewClaimInterestMsg("te"),
			wantCode: types.CodeInvalidUsername,
		},
		"invalid claim interest - Username is too long": {
			msg:      NewClaimInterestMsg("testtesttesttesttesttest"),
			wantCode: types.CodeInvalidUsername,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()

		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: diff error: got %v, want %v", testName, tc.wantCode, tc.wantCode)
			}
			continue
		}
		if got.Code() != tc.wantCode {
			t.Errorf("%s: diff error code: got %v, want %v", testName, got.Code(), tc.wantCode)
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
