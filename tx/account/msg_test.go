package account

import (
	"testing"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	userA = types.AccountKey("userA")
	userB = types.AccountKey("userB")

	memo1 = "This is a memo!"
)

func TestFollowMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      FollowMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: FollowMsg{
				Follower: userA,
				Followee: userB,
			},
			wantCode: sdk.CodeOK,
		},
		"invalid follower - Username is too short": {
			msg: FollowMsg{
				Follower: "re",
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid follower - Username is too long": {
			msg: FollowMsg{
				Follower: "registerregisterregis",
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid followee - Username is too short": {
			msg: FollowMsg{
				Follower: userA,
				Followee: "re",
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid followee - Username is too long": {
			msg: FollowMsg{
				Follower: userA,
				Followee: "registerregisterregis",
			},
			wantCode: types.CodeInvalidUsername,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()

		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: ValidateBasic(%v) error: got %v, want %v", testName, tc.msg, nil, sdk.CodeOK)
			}
			return
		}
		if got.ABCICode() != tc.wantCode {
			t.Errorf("%s: ValidateBasic(%v) errorCode: got %v, want %v", testName, tc.msg, got.ABCICode(), tc.wantCode)
		}
	}
}

func TestUnfollowMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      UnfollowMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: UnfollowMsg{
				Follower: userA,
				Followee: userB,
			},
			wantCode: sdk.CodeOK,
		},
		"invalid follower - Username is too short": {
			msg: UnfollowMsg{
				Follower: "re",
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid follower - Username is too long": {
			msg: UnfollowMsg{
				Follower: "registerregisterregis",
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid followee - Username is too short": {
			msg: UnfollowMsg{
				Follower: userA,
				Followee: "re",
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid followee - Username is too long": {
			msg: UnfollowMsg{
				Follower: userA,
				Followee: "registerregisterregis",
			},
			wantCode: types.CodeInvalidUsername,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()

		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: ValidateBasic(%v) error: got %v, want %v", testName, tc.msg, nil, sdk.CodeOK)
			}
			return
		}
		if got.ABCICode() != tc.wantCode {
			t.Errorf("%s: ValidateBasic(%v) errorCode: got %v, want %v", testName, tc.msg, got.ABCICode(), tc.wantCode)
		}
	}
}

func TestTransferMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      TransferMsg
		wantCode sdk.CodeType
	}{
		"normal case - transfer to an username": {
			msg: TransferMsg{
				Sender:       userA,
				ReceiverName: userB,
				Amount:       types.LNO("1900"),
				Memo:         memo1,
			},
			wantCode: sdk.CodeOK,
		},
		"normal case - transfer to an address": {
			msg: TransferMsg{
				Sender:       userA,
				ReceiverAddr: sdk.Address("2137192887931"),
				Amount:       types.LNO("1900"),
				Memo:         memo1,
			},
			wantCode: sdk.CodeOK,
		},
		"invalid transfer - no receiver provided": {
			msg: TransferMsg{
				Sender: userA,
				Amount: types.LNO("1900"),
				Memo:   memo1,
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid transfer - both username and address are invalid": {
			msg: TransferMsg{
				Sender:       userA,
				ReceiverName: types.AccountKey(""),
				ReceiverAddr: sdk.Address(""),
				Amount:       types.LNO("1900"),
				Memo:         memo1,
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid transfer -  amount is invalid": {
			msg: TransferMsg{
				Sender:       userA,
				ReceiverName: userB,
				ReceiverAddr: sdk.Address(""),
				Amount:       types.LNO("-1900"),
				Memo:         memo1,
			},
			wantCode: sdk.CodeInvalidCoins,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()

		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: ValidateBasic(%v) error: got %v, want %v", testName, tc.msg, nil, sdk.CodeOK)
			}
			return
		}
		if got.ABCICode() != tc.wantCode {
			t.Errorf("%s: ValidateBasic(%v) errorCode: got %v, want %v", testName, tc.msg, got.ABCICode(), tc.wantCode)
		}
	}
}

func TestTransferMsgPermission(t *testing.T) {
	testCases := map[string]struct {
		msg            TransferMsg
		wantOK         bool
		wantPermission types.Permission
	}{
		"normal case": {
			msg: TransferMsg{
				Sender:       userA,
				ReceiverName: userB,
				Amount:       types.LNO("1900"),
				Memo:         memo1,
			},
			wantOK: true,
		},
	}

	for testName, tc := range testCases {
		gotPermissionLevel := tc.msg.Get(types.PermissionLevel)
		gotPermission, ok := gotPermissionLevel.(types.Permission)

		if ok != tc.wantOK {
			t.Errorf("%s: Get(%v): got %v, want %v", testName, tc.msg, ok, tc.wantOK)
			return
		}
		if gotPermission != types.TransactionPermission {
			t.Errorf("%s: Get(%v): got %v, want %v", testName, tc.msg, gotPermission, types.TransactionPermission)
		}
	}
}
