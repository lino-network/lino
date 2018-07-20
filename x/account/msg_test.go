package account

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	userA = types.AccountKey("userA")
	userB = types.AccountKey("userB")

	memo1       = "This is a memo!"
	invalidMemo = "Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! "
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
				t.Errorf("%s: diff error: got %v, want %v", testName, sdk.CodeOK, tc.wantCode)
			}
			continue
		}
		if got.Code() != tc.wantCode {
			t.Errorf("%s: diff error code: got %v, want %v", testName, got.Code(), tc.wantCode)
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
				t.Errorf("%s: diff error: got %v, want %v", testName, sdk.CodeOK, tc.wantCode)
			}
			continue
		}
		if got.Code() != tc.wantCode {
			t.Errorf("%s: diff error code: got %v, want %v", testName, got.Code(), tc.wantCode)
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
				Sender:   userA,
				Receiver: userB,
				Amount:   types.LNO("1900"),
				Memo:     memo1,
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
		"invalid transfer -  amount is invalid": {
			msg: TransferMsg{
				Sender:   userA,
				Receiver: userB,
				Amount:   types.LNO("-1900"),
				Memo:     memo1,
			},
			wantCode: types.CodeInvalidCoins,
		},
		"invalid transfer -  memo is invalid": {
			msg: TransferMsg{
				Sender:   userA,
				Receiver: userB,
				Amount:   types.LNO("1900"),
				Memo:     invalidMemo,
			},
			wantCode: types.CodeInvalidMemo,
		},
		"valid lino": {
			msg: TransferMsg{
				Sender:   userA,
				Receiver: userB,
				Amount:   types.LNO("100"),
				Memo:     memo1,
			},
			wantCode: sdk.CodeOK,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()

		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: diff error: got %v, want %v", testName, sdk.CodeOK, tc.wantCode)
			}
			continue
		}
		if got.Code() != tc.wantCode {
			t.Errorf("%s: diff error code: got %v, want %v", testName, got.Code(), tc.wantCode)
		}
	}
}

func TestRecoverMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      RecoverMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: NewRecoverMsg("test", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"invalid recover - Username is too short": {
			msg: NewRecoverMsg("te", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"invalid recover - Username is too long": {
			msg: NewRecoverMsg("testtesttesttesttesttest", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
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

func TestClaimMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      ClaimMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg:      NewClaimMsg("test"),
			wantCode: sdk.CodeOK,
		},
		"invalid claim - Username is too short": {
			msg:      NewClaimMsg("te"),
			wantCode: types.CodeInvalidUsername,
		},
		"invalid claim - Username is too long": {
			msg:      NewClaimMsg("testtesttesttesttesttest"),
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

func TestUpdateAccountMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      UpdateAccountMsg
		wantCode sdk.CodeType
	}{
		"normal case - update JSON Meta": {
			msg: UpdateAccountMsg{
				Username: userA,
				JSONMeta: "{'test':'test'}",
			},
			wantCode: sdk.CodeOK,
		},
		"normal case - update JSON Meta too long": {
			msg: UpdateAccountMsg{
				Username: userA,
				JSONMeta: string(make([]byte, 501)),
			},
			wantCode: types.CodeInvalidJSONMeta,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()
		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: diff error: got %v, want %v", testName, tc.wantCode, tc.wantCode)
				return
			}
			continue
		}
		if got.Code() != tc.wantCode {
			t.Errorf("%s: diff error code: got %v, want %v", testName, got.Code(), tc.wantCode)
		}
	}
}

func TestRegisterUsername(t *testing.T) {
	testCases := map[string]struct {
		msg      RegisterMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: NewRegisterMsg("referrer", "newuser", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"register username minimum length": {
			msg: NewRegisterMsg("referrer", "new", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"register username maximum length": {
			msg: NewRegisterMsg("referrer", "newnewnewnewnewnewne", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"register username length exceeds requirement": {
			msg: NewRegisterMsg("referrer", "newnewnewnewnewnewnew", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"register username length doesn't meet requirement": {
			msg: NewRegisterMsg("referrer", "ne", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"referrer invalid": {
			msg: NewRegisterMsg("", "newuser", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"register fee invalid": {
			msg: NewRegisterMsg("", "newuser", "1.", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()

		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: diff error: got %v, want %v", testName, sdk.CodeOK, tc.wantCode)
			}
			continue
		}
		if got.Code() != tc.wantCode {
			t.Errorf("%s: diff errorCode: got %v, want %v", testName, got.Code(), tc.wantCode)
		}
	}

	// Illegel character
	registerList := [...]string{"register#", "_register", "-register", "reg@ister",
		"reg*ister", "register!", "register()", "reg$ister", "reg ister", " register",
		"reg=ister", "register^", "register.", "reg$ister,", "Register"}
	for _, register := range registerList {
		msg := NewRegisterMsg(
			"referer", register, "0", crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey())
		result := msg.ValidateBasic()
		assert.Equal(t, result, ErrInvalidUsername("illeagle input"))
	}
}

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              types.Msg
		expectPermission types.Permission
	}{
		"transfer to user": {
			msg:              NewTransferMsg("test", "test_user", types.LNO("1"), "memo"),
			expectPermission: types.TransactionPermission,
		},
		"follow": {
			msg:              NewFollowMsg("userA", "userB"),
			expectPermission: types.PostPermission,
		},
		"unfollow": {
			msg:              NewUnfollowMsg("userA", "userB"),
			expectPermission: types.PostPermission,
		},
		"recover": {
			msg: NewRecoverMsg(
				"userA", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey()),
			expectPermission: types.ResetPermission,
		},
		"claim": {
			msg:              NewClaimMsg("test"),
			expectPermission: types.PostPermission,
		},
		"register msg": {
			msg: NewRegisterMsg("referrer", "test", "0", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey()),
			expectPermission: types.TransactionPermission,
		},
		"update msg": {
			msg:              NewUpdateAccountMsg("user", "{'test':'test'}"),
			expectPermission: types.PostPermission,
		},
	}

	for testName, cs := range cases {
		permission := cs.msg.GetPermission()
		if cs.expectPermission != permission {
			t.Errorf("%s: diff permission, got %v, want %v", testName, permission, cs.expectPermission)
			return
		}
	}
}
