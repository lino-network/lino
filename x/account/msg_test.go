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
				t.Errorf("%s: ValidateBasic(%v) error: got %v, want %v", testName, tc.msg, nil, sdk.CodeOK)
			}
			continue
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
			continue
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
			wantCode: sdk.CodeInvalidCoins,
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
				t.Errorf("%s: ValidateBasic(%v) error: got %v, want %v", testName, tc.msg, nil, sdk.CodeOK)
			}
			continue
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
			msg: NewRecoverMsg("test", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"invalid recover - Username is too short": {
			msg: NewRecoverMsg("te", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"invalid recover - Username is too long": {
			msg: NewRecoverMsg("testtesttesttesttesttest", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()

		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: ValidateBasic(%v) error: got %v, want %v", testName, tc.msg, nil, tc.wantCode)
			}
			continue
		}
		if got.ABCICode() != tc.wantCode {
			t.Errorf("%s: ValidateBasic(%v) errorCode: got %v, want %v", testName, tc.msg, got.ABCICode(), tc.wantCode)
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
			wantCode: types.CodeInvalidMsg,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()
		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: ValidateBasic(%v) error: got %v, want %v", testName, tc.msg, nil, tc.wantCode)
				return
			}
			continue
		}
		if got.ABCICode() != tc.wantCode {
			t.Errorf("%s: ValidateBasic(%v) errorCode: got %v, want %v", testName, tc.msg, got.ABCICode(), tc.wantCode)
		}
	}
}

func TestRegisterUsername(t *testing.T) {
	testCases := map[string]struct {
		msg      RegisterMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: NewRegisterMsg("referrer", "newUser", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"register username minimum length": {
			msg: NewRegisterMsg("referrer", "new", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"register username maximum length": {
			msg: NewRegisterMsg("referrer", "newnewnewnewnewnewne", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"register username length exceeds requirement": {
			msg: NewRegisterMsg("referrer", "newnewnewnewnewnewnew", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"register username length doesn't meet requirement": {
			msg: NewRegisterMsg("referrer", "ne", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"referrer invalid": {
			msg: NewRegisterMsg("", "newUser", "1", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"register fee invalid": {
			msg: NewRegisterMsg("", "newUser", "1.", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
	}

	for testName, tc := range testCases {
		got := tc.msg.ValidateBasic()

		if got == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: ValidateBasic(%v) error: got %v, want %v", testName, tc.msg, nil, sdk.CodeOK)
			}
			continue
		}
		if got.ABCICode() != tc.wantCode {
			t.Errorf("%s: ValidateBasic(%v) errorCode: got %v, want %v", testName, tc.msg, got.ABCICode(), tc.wantCode)
		}
	}

	// Illegel character
	registerList := [...]string{"register#", "_register", "-register", "reg@ister",
		"reg*ister", "register!", "register()", "reg$ister", "reg ister", " register",
		"reg=ister", "register^", "register.", "reg$ister,"}
	for _, register := range registerList {
		msg := NewRegisterMsg(
			"referer", register, "0", crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey())
		result := msg.ValidateBasic()
		assert.Equal(t, result, ErrInvalidUsername("illeagle input"))
	}
}

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              sdk.Msg
		expectPermission types.Permission
	}{
		"transfer to user": {
			NewTransferMsg("test", "test_user", types.LNO("1"), "memo"),
			types.TransactionPermission},
		"follow": {
			NewFollowMsg("userA", "userB"),
			types.PostPermission},
		"unfollow": {
			NewUnfollowMsg("userA", "userB"),
			types.PostPermission},
		"recover": {
			NewRecoverMsg(
				"userA", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey()),
			types.MasterPermission},
		"claim": {
			NewClaimMsg("test"), types.PostPermission},
		"register msg": {
			NewRegisterMsg("referrer", "test", "0", crypto.GenPrivKeyEd25519().PubKey(),
				crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey()),
			types.TransactionPermission},
		"update msg": {
			NewUpdateAccountMsg("user", "{'test':'test'}"), types.PostPermission},
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
