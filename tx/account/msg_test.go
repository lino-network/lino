package account

import (
	"testing"

	"github.com/lino-network/lino/types"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFollowMsg(t *testing.T) {
	testCases := map[string]struct {
		msg      FollowMsg
		wantCode sdk.CodeType
	}{
		"normal case": {
			msg: FollowMsg{
				Follower: "userA",
				Followee: "userB",
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
				Follower: "userA",
				Followee: "re",
			},
			wantCode: types.CodeInvalidUsername,
		},
		"invalid followee - Username is too long": {
			msg: FollowMsg{
				Follower: "userA",
				Followee: "registerregisterregis",
			},
			wantCode: types.CodeInvalidUsername,
		},
	}

	for testName, tc := range testCases {
		result := tc.msg.ValidateBasic()

		if result == nil {
			if tc.wantCode != sdk.CodeOK {
				t.Errorf("%s: ValidateBasic(%v) error: got %v, want %v", testName, tc.msg, nil, sdk.CodeOK)
			}
			return
		}
		if result.ABCICode() != tc.wantCode {
			t.Errorf("%s: ValidateBasic(%v) errorCode: got %v, want %v", testName, tc.msg, result.ABCICode(), tc.wantCode)
		}
	}
}

func TestUnfollowMsg(t *testing.T) {
	follower := "userA"
	followee := "userB"
	msg := NewUnfollowMsg(follower, followee)
	result := msg.ValidateBasic()
	assert.Nil(t, result)

	// Follower Username length invalid
	follower = "re"
	msg = NewUnfollowMsg(follower, followee)
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())

	follower = "registerregisterregis"
	msg = NewUnfollowMsg(follower, followee)
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())

	// Followee Username length invalid
	follower = "userA"
	followee = "re"
	msg = NewUnfollowMsg(follower, followee)
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())

	followee = "registerregisterregis"
	msg = NewUnfollowMsg(follower, followee)
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())
}

func TestTransferMsg(t *testing.T) {
	// normal transfer to a username
	sender := "userA"
	receiverName := "userB"
	amount := types.LNO("1900")
	memo := "This is a memo!"

	msg := NewTransferMsg(sender, amount, memo, TransferToUser(receiverName))
	result := msg.ValidateBasic()
	assert.Nil(t, result)

	// normal transfer to an address
	receiverAddr := sdk.Address("2137192887931")
	msg = NewTransferMsg(sender, amount, memo, TransferToAddr(receiverAddr))
	result = msg.ValidateBasic()
	assert.Nil(t, result)

	// invalid transfer: no receiver provided
	msg = NewTransferMsg(sender, amount, memo)
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())

	// invalid transfer: both username and address are invalid
	receiverName = ""
	receiverAddr = sdk.Address("")
	msg = NewTransferMsg(sender, amount, memo, TransferToUser(receiverName), TransferToAddr(receiverAddr))
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())

	// invalid transfer: amount is invalid
	receiverName = "userB"
	amount = types.LNO("-1900")
	msg = NewTransferMsg(sender, amount, memo, TransferToUser(receiverName))
	result = msg.ValidateBasic()
	assert.Equal(t, result, sdk.ErrInvalidCoins("LNO can't be less than lower bound"))
}

func TestTransferMsgPermission(t *testing.T) {
	msg := NewTransferMsg("userA", types.LNO("1900"), "This is a memo!", TransferToUser("userB"))
	permissionLevel := msg.Get(types.PermissionLevel)
	permission, ok := permissionLevel.(types.Permission)
	assert.Equal(t, permission, types.TransactionPermission)
	assert.Equal(t, ok, true)
}
