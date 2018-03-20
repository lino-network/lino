package account

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type FollowMsg struct {
	Follower types.AccountKey `json:"follower"`
	Followee types.AccountKey `json:"followee"`
}

type UnfollowMsg struct {
	Follower types.AccountKey `json:"follower"`
	Followee types.AccountKey `json:"followee"`
}

// we can support to transfer to an user or an address
type TransferMsg struct {
	Sender       types.AccountKey `json:"sender"`
	ReceiverName types.AccountKey `json:"receiver_name"`
	ReceiverAddr sdk.Address      `json:"receiver_addr"`
	Amount       sdk.Coins        `json:"amount"`
	Memo         []byte           `json:"memo"`
}

type TransferOption func(*TransferMsg)

func TransferToUser(userName string) TransferOption {
	return func(args *TransferMsg) {
		args.ReceiverName = types.AccountKey(userName)
	}
}

func TransferToAddr(addr string) TransferOption {
	return func(args *TransferMsg) {
		args.ReceiverAddr = sdk.Address(addr)
	}
}

type GrantMsg struct {
}

type RecoverMsg struct {
}

type UpdateMsg struct {
}

var _ sdk.Msg = FollowMsg{}
var _ sdk.Msg = UnfollowMsg{}
var _ sdk.Msg = TransferMsg{}
var _ sdk.Msg = GrantMsg{}
var _ sdk.Msg = RecoverMsg{}
var _ sdk.Msg = UpdateMsg{}

//----------------------------------------
// Follow Msg Implementations

func NewFollowMsg(follower string, followee string) FollowMsg {
	return FollowMsg{
		Follower: types.AccountKey(follower),
		Followee: types.AccountKey(followee),
	}
}

func (msg FollowMsg) Type() string { return types.AccountRouterName } // TODO: "account/register"

func (msg FollowMsg) ValidateBasic() sdk.Error {
	if len(msg.Follower) < types.MinimumUsernameLength ||
		len(msg.Followee) < types.MinimumUsernameLength ||
		len(msg.Follower) > types.MaximumUsernameLength ||
		len(msg.Followee) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}
	return nil
}

func (msg FollowMsg) String() string {
	return fmt.Sprintf("FollowMsg{Follower:%v, Followee:%v}", msg.Follower, msg.Followee)
}

func (msg FollowMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg FollowMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg FollowMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Follower)}
}

//----------------------------------------
// Unfollow Msg Implementations

func NewUnfollowMsg(follower string, followee string) UnfollowMsg {
	return UnfollowMsg{
		Follower: types.AccountKey(follower),
		Followee: types.AccountKey(followee),
	}
}

func (msg UnfollowMsg) Type() string { return types.AccountRouterName } // TODO: "account/register"

func (msg UnfollowMsg) ValidateBasic() sdk.Error {
	if len(msg.Follower) < types.MinimumUsernameLength ||
		len(msg.Followee) < types.MinimumUsernameLength ||
		len(msg.Follower) > types.MaximumUsernameLength ||
		len(msg.Followee) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}
	return nil
}

func (msg UnfollowMsg) String() string {
	return fmt.Sprintf("UnfollowMsg{Follower:%v, Followee:%v}", msg.Follower, msg.Followee)
}

func (msg UnfollowMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg UnfollowMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg UnfollowMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Follower)}
}

//----------------------------------------
// Transfer Msg Implementations

func NewTransferMsg(sender string, amount sdk.Coins, memo []byte, setters ...TransferOption) TransferMsg {
	msg := &TransferMsg{
		Sender:       types.AccountKey(sender),
		Amount:       amount,
		Memo:         memo,
		ReceiverName: nil,
		ReceiverAddr: nil,
	}
	for _, setter := range setters {
		setter(msg)
	}
	return msg
}

func (msg TransferMsg) Type() string { return types.AccountRouterName } // TODO: "account/register"

func (msg TransferMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) < types.MinimumUsernameLength ||
		len(msg.Sender) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}

	// should have either receiver's addr or username
	if msg.ReceiverAddr == nil && msg.ReceiverName == nil {
		return ErrInvalidUsername("invalid receiver")
	}

	// cannot transfer a negative amount of money
	if msg.Amount.IsPositive() == false {
		return ErrInvalidCoins("invalid coin amount")
	}

	return nil
}

func (msg TransferMsg) String() string {
	return fmt.Sprintf("TransferMsg{Sender:%v, Receiver:%v}", msg.Sender, msg.Receiver)
}

func (msg TransferMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg TransferMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg TransferMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Sender)}
}
