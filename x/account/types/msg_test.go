package types

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	userA = types.AccountKey("usera")
	userB = types.AccountKey("userb")

	memo1       = "This is a memo!"
	invalidMemo = "Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! Memo is too long!!! "
)

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
			msg: NewRecoverMsg("test", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey()),
			wantCode: sdk.CodeOK,
		},
		"invalid recover - Username is too short": {
			msg:      NewRecoverMsg("te", secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey()),
			wantCode: types.CodeInvalidUsername,
		},
		"invalid recover - Username is too long": {
			msg: NewRecoverMsg("testtesttesttesttesttest",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
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
		"invalid username": {
			msg: UpdateAccountMsg{
				Username: "us",
				JSONMeta: "{'test':'test'}",
			},
			wantCode: types.CodeInvalidUsername,
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
		msgv2    RegisterMsgV2
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: NewRegisterMsg("referrer", "newuser", "1", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			msgv2: NewRegisterMsgV2("referrer", "newuser", "1",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"normal case with dot": {
			msg: NewRegisterMsg("zhimao.liu", "newuser", "1", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			msgv2: NewRegisterMsgV2("zhimao.liu", "newuser", "1",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"register username minimum length": {
			msg: NewRegisterMsg("referrer", "new", "1", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			msgv2: NewRegisterMsgV2("referrer", "new", "1",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"register username maximum length": {
			msg: NewRegisterMsg("referrer", "newnewnewnewnewnewne", "1", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			msgv2: NewRegisterMsgV2("referrer", "newnewnewnewnewnewne", "1",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			wantCode: sdk.CodeOK,
		},
		"register username length exceeds requirement": {
			msg: NewRegisterMsg("referrer", "newnewnewnewnewnewnew", "1", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			msgv2: NewRegisterMsgV2("referrer", "newnewnewnewnewnewnew", "1",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"register username length doesn't meet requirement": {
			msg: NewRegisterMsg("referrer", "ne", "1", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			msgv2: NewRegisterMsgV2("referrer", "ne", "1",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			wantCode: types.CodeInvalidUsername,
		},
		"register fee invalid": {
			msg: NewRegisterMsg("newuser", "newuser", "1.", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			msgv2: NewRegisterMsgV2("newuser", "newuser", "1.",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
			),
			wantCode: types.CodeInvalidCoins,
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

	for testName, tc := range testCases {
		got := tc.msgv2.ValidateBasic()

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
	registerList := [...]string{"register#", "_register", "-register", "reg@ister", "re--gister",
		"reg*ister", "register!", "register()", "reg$ister", "reg ister", " register", "re_-gister",
		"reg=ister", "register^", "register.", "reg$ister,", "Register", "r__egister", "reGister",
		"r_--gister", "re.-gister", ".re-gister", "re-gister.", "register_", "register-", "a.2.2.-.-..2"}
	for _, register := range registerList {
		msg := NewRegisterMsg(
			"referer", register, "1", secp256k1.GenPrivKey().PubKey(),
			secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey())
		result := msg.ValidateBasic()
		assert.Equal(t, ErrInvalidUsername("illegal username"), result, "%s", register)
		msgv2 := NewRegisterMsgV2(
			"referer", register, "1", secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey())
		result = msgv2.ValidateBasic()
		assert.Equal(t, ErrInvalidUsername("illegal username"), result, "%s", register)
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
		"recover": {
			msg: NewRecoverMsg(
				"userA", secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey()),
			expectPermission: types.ResetPermission,
		},
		"register msg": {
			msg: NewRegisterMsg("referrer", "test", "0", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey()),
			expectPermission: types.TransactionPermission,
		},
		"register msgv2": {
			msg: NewRegisterMsgV2("referrer", "test", "0",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey()),
			expectPermission: types.TransactionPermission,
		},
		"update msg": {
			msg:              NewUpdateAccountMsg("user", "{'test':'test'}"),
			expectPermission: types.AppPermission,
		},
	}

	for testName, tc := range cases {
		permission := tc.msg.GetPermission()
		if tc.expectPermission != permission {
			t.Errorf("%s: diff permission, got %v, want %v", testName, permission, tc.expectPermission)
			return
		}
	}
}

func TestGetSignBytes(t *testing.T) {
	cases := map[string]struct {
		msg types.Msg
	}{
		"transfer to user": {
			msg: NewTransferMsg("test", "test_user", types.LNO("1"), "memo"),
		},
		"recover msg with public key type Secp256k1": {
			msg: NewRecoverMsg(
				"userA",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey()),
		},
		"register msg with public key type Secp256k1": {
			msg: NewRegisterMsg("referrer", "test", "0", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey()),
		},
		"register msgv2 with public key type Secp256k1": {
			msg: NewRegisterMsgV2("referrer", "test", "0",
				secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey()),
		},
		"update msg": {
			msg: NewUpdateAccountMsg("user", "{'test':'test'}"),
		},
	}

	for testName, tc := range cases {
		require.NotPanics(t, func() { tc.msg.GetSignBytes() }, testName)
	}
}

func TestGetSigners(t *testing.T) {
	privKey := secp256k1.GenPrivKey()
	cases := map[string]struct {
		msg           types.Msg
		expectSigners []types.AccountKey
	}{
		"transfer to user": {
			msg:           NewTransferMsg("test", "test_user", types.LNO("1"), "memo"),
			expectSigners: []types.AccountKey{"test"},
		},
		"transfer to address": {
			msg:           NewTransferMsg("test", string(secp256k1.GenPrivKey().PubKey().Address()), types.LNO("1"), "memo"),
			expectSigners: []types.AccountKey{"test"},
		},
		"transfer from address": {
			msg:           NewTransferMsg(string(privKey.PubKey().Address()), "test", types.LNO("1"), "memo"),
			expectSigners: []types.AccountKey{types.AccountKey(privKey.PubKey().Address())},
		},
		"recover msg with public key type Secp256k1": {
			msg: NewRecoverMsg(
				"userA",
				privKey.PubKey(),
				secp256k1.GenPrivKey().PubKey()),
			expectSigners: []types.AccountKey{"userA", types.AccountKey(privKey.PubKey().Address())},
		},
		"register msg with public key type Secp256k1": {
			msg: NewRegisterMsg("referrer", "test", "0", secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey()),
			expectSigners: []types.AccountKey{"referrer"},
		},
		"register msgv2 with public key type Secp256k1": {
			msg: NewRegisterMsgV2("referrer", "test", "0",
				privKey.PubKey(),
				secp256k1.GenPrivKey().PubKey()),
			expectSigners: []types.AccountKey{"referrer", types.AccountKey(privKey.PubKey().Address())},
		},
		"update msg": {
			msg:           NewUpdateAccountMsg("user", "{'test':'test'}"),
			expectSigners: []types.AccountKey{"user"},
		},
	}

	for testName, tc := range cases {
		if len(tc.msg.GetSigners()) != len(tc.expectSigners) {
			t.Errorf("%s: expect number of signers wrong, got %v, want %v", testName, len(tc.msg.GetSigners()), len(tc.expectSigners))
			return
		}
		for i, signer := range tc.msg.GetSigners() {
			if types.AccountKey(signer) != tc.expectSigners[i] {
				t.Errorf("%s: expect signer wrong, got %v, want %v", testName, types.AccountKey(signer), tc.expectSigners[i])
				return
			}
		}
	}
}
