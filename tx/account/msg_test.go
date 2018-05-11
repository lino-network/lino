package account

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"

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

func TestRecoverMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      RecoverMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: RecoverMsg{
				Username:             "test",
				NewPostPubKey:        crypto.GenPrivKeyEd25519().PubKey(),
				NewTransactionPubKey: crypto.GenPrivKeyEd25519().PubKey(),
			},
			wantCode: sdk.CodeOK,
		},
		"invalid recover - Username is too short": {
			msg: RecoverMsg{
				Username:             "te",
				NewPostPubKey:        crypto.GenPrivKeyEd25519().PubKey(),
				NewTransactionPubKey: crypto.GenPrivKeyEd25519().PubKey(),
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid recover - Username is too long": {
			msg: RecoverMsg{
				Username:             "testtesttesttesttesttest",
				NewPostPubKey:        crypto.GenPrivKeyEd25519().PubKey(),
				NewTransactionPubKey: crypto.GenPrivKeyEd25519().PubKey(),
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

func TestCheckingToSavingMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      CheckingToSavingMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: CheckingToSavingMsg{
				Username: "test",
				Amount:   types.LNO("100"),
			},
			wantCode: sdk.CodeOK,
		},
		"username is too short": {
			msg: CheckingToSavingMsg{
				Username: "",
				Amount:   types.LNO("100"),
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid amount - zero": {
			msg: CheckingToSavingMsg{
				Username: "test",
				Amount:   types.LNO("0"),
			},
			wantCode: sdk.CodeInvalidCoins,
		},
		"invalid amount - negative": {
			msg: CheckingToSavingMsg{
				Username: "test",
				Amount:   types.LNO("-1"),
			},
			wantCode: sdk.CodeInvalidCoins,
		},
		"invalid amount - overflow": {
			msg: CheckingToSavingMsg{
				Username: "test",
				Amount:   types.LNO("1000000000000000"),
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

func TestSavingToCheckingMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      SavingToCheckingMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: SavingToCheckingMsg{
				Username: "test",
				Amount:   types.LNO("100"),
			},
			wantCode: sdk.CodeOK,
		},
		"username is too short": {
			msg: SavingToCheckingMsg{
				Username: "",
				Amount:   types.LNO("100"),
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid amount - zero": {
			msg: SavingToCheckingMsg{
				Username: "test",
				Amount:   types.LNO("0"),
			},
			wantCode: sdk.CodeInvalidCoins,
		},
		"invalid amount - overflow": {
			msg: SavingToCheckingMsg{
				Username: "test",
				Amount:   types.LNO("1000000000000000"),
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

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              sdk.Msg
		expectPermission types.Permission
	}{
		"saving to checking": {
			NewSavingToCheckingMsg("test", types.LNO("1")),
			types.TransactionPermission},
		"checking to saving": {
			NewCheckingToSavingMsg("test", types.LNO("1")),
			types.TransactionPermission},
		"transfer to user": {
			NewTransferMsg("test", types.LNO("1"), "memo", TransferToUser("test_user")),
			types.TransactionPermission},
		"transfer to address": {
			NewTransferMsg("test", types.LNO("1"), "memo", TransferToAddr(sdk.Address("test_address"))),
			types.TransactionPermission},
		"follow": {
			NewFollowMsg("userA", "userB"),
			types.PostPermission},
		"unfollow": {
			NewUnfollowMsg("userA", "userB"),
			types.PostPermission},
		"recover": {
			NewRecoverMsg(
				"userA", crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey()),
			types.MasterPermission},
		"claim": {
			NewClaimMsg("test"), types.PostPermission},
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
