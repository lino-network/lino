package account

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestFollowMsg(t *testing.T) {
	follower := "userA"
	followee := "userB"
	msg := NewFollowMsg(follower, followee)
	result := msg.ValidateBasic()
	assert.Nil(t, result)

	// Follower Username length invalid
	follower = "re"
	msg = NewFollowMsg(follower, followee)
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())

	follower = "registerregisterregis"
	msg = NewFollowMsg(follower, followee)
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())

	// Followee Username length invalid
	follower = "userA"
	followee = "re"
	msg = NewFollowMsg(follower, followee)
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())

	followee = "registerregisterregis"
	msg = NewFollowMsg(follower, followee)
	result = msg.ValidateBasic()
	assert.Equal(t, result, ErrInvalidUsername())
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
	amount := types.TestLNO(sdk.NewRat(1900))
	memo := []byte("This is a memo!")

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
	amount = types.TestLNO(sdk.NewRat(-1900))
	msg = NewTransferMsg(sender, amount, memo, TransferToUser(receiverName))
	result = msg.ValidateBasic()
	assert.Equal(t, result, sdk.ErrInvalidCoins("TestLNO can't be less than lower bound"))

}
