package account

// nolint
import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type FollowMsg struct {
	Follower AccountKey `json:"follower"`
	Followee AccountKey `json:"followee"`
}

type UnfollowMsg struct {
	Follower AccountKey `json:"follower"`
	Followee AccountKey `json:"followee"`
}

// we can support to transfer to an user or an address
type TransferMsg struct {
	Sender       AccountKey  `json:"sender"`
	ReceiverName AccountKey  `json:"receiver_name"`
	ReceiverAddr sdk.Address `json:"receiver_addr"`
	Amount       sdk.Coins   `json:"amount"`
	Memo         []byte      `json:"memo"`
}

type TransferOption func(*TransferMsg)

func TransferToUser(userName string) TransferOption {
	return func(args *TransferMsg) {
		args.ReceiverName = AccountKey(userName)
	}
}

func TransferToAddr(addr sdk.Address) TransferOption {
	return func(args *TransferMsg) {
		args.ReceiverAddr = addr
	}
}

var _ sdk.Msg = FollowMsg{}
var _ sdk.Msg = UnfollowMsg{}
var _ sdk.Msg = TransferMsg{}

//----------------------------------------
// Follow Msg Implementations

func NewFollowMsg(follower string, followee string) FollowMsg {
	return FollowMsg{
		Follower: AccountKey(follower),
		Followee: AccountKey(followee),
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
		Follower: AccountKey(follower),
		Followee: AccountKey(followee),
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
		Sender: AccountKey(sender),
		Amount: amount,
		Memo:   memo,
	}
	for _, setter := range setters {
		setter(msg)
	}
	return *msg
}

func (msg TransferMsg) Type() string { return types.AccountRouterName } // TODO: "account/register"

func (msg TransferMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) < types.MinimumUsernameLength ||
		len(msg.Sender) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}

	// should have either receiver's addr or username
	if len(msg.ReceiverAddr) == 0 && len(msg.ReceiverName) == 0 {
		return ErrInvalidUsername("invalid receiver")
	}

	// cannot transfer a negative amount of money
	if msg.Amount.IsPositive() == false {
		return sdk.ErrInvalidCoins("invalid coin amount")
	}

	// cannot transfer othe coin types
	if len(msg.Amount) != 1 || msg.Amount[0].Denom != "lino" {
		return sdk.ErrInvalidCoins("invalid coin type")
	}

	return nil
}

func (msg TransferMsg) String() string {
	return fmt.Sprintf("TransferMsg{Sender:%v, ReceiverName:%v, ReceiverAddr:%v}",
		msg.Sender, msg.ReceiverName, msg.ReceiverAddr)
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
