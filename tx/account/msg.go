package account

// nolint
import (
	"encoding/json"
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FollowMsg struct {
	Follower types.AccountKey `json:"follower"`
	Followee types.AccountKey `json:"followee"`
}

type UnfollowMsg struct {
	Follower types.AccountKey `json:"follower"`
	Followee types.AccountKey `json:"followee"`
}

type ClaimMsg struct {
	Username types.AccountKey `json:"username"`
}

// we can support to transfer to an user or an address
type TransferMsg struct {
	Sender       types.AccountKey `json:"sender"`
	ReceiverName types.AccountKey `json:"receiver_name"`
	ReceiverAddr sdk.Address      `json:"receiver_addr"`
	Amount       types.LNO        `json:"amount"`
	Memo         []byte           `json:"memo"`
}

type TransferOption func(*TransferMsg)

func TransferToUser(userName string) TransferOption {
	return func(args *TransferMsg) {
		args.ReceiverName = types.AccountKey(userName)
	}
}

func TransferToAddr(addr sdk.Address) TransferOption {
	return func(args *TransferMsg) {
		args.ReceiverAddr = addr
	}
}

var _ sdk.Msg = FollowMsg{}
var _ sdk.Msg = UnfollowMsg{}
var _ sdk.Msg = ClaimMsg{}
var _ sdk.Msg = TransferMsg{}

// Follow Msg Implementations
func NewFollowMsg(follower string, followee string) FollowMsg {
	return FollowMsg{
		Follower: types.AccountKey(follower),
		Followee: types.AccountKey(followee),
	}
}

func (msg FollowMsg) Type() string { return types.AccountRouterName }

func (msg FollowMsg) ValidateBasic() sdk.Error {
	if len(msg.Follower) < types.MinimumUsernameLength ||
		len(msg.Followee) < types.MinimumUsernameLength ||
		len(msg.Follower) > types.MaximumUsernameLength ||
		len(msg.Followee) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
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

// Unfollow Msg Implementations
func NewUnfollowMsg(follower string, followee string) UnfollowMsg {
	return UnfollowMsg{
		Follower: types.AccountKey(follower),
		Followee: types.AccountKey(followee),
	}
}

func (msg UnfollowMsg) Type() string { return types.AccountRouterName }

func (msg UnfollowMsg) ValidateBasic() sdk.Error {
	if len(msg.Follower) < types.MinimumUsernameLength ||
		len(msg.Followee) < types.MinimumUsernameLength ||
		len(msg.Follower) > types.MaximumUsernameLength ||
		len(msg.Followee) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
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

// Claim Msg Implementations
func NewClaimMsg(username string) ClaimMsg {
	return ClaimMsg{
		Username: types.AccountKey(username),
	}
}

func (msg ClaimMsg) Type() string { return types.AccountRouterName }

func (msg ClaimMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg ClaimMsg) String() string {
	return fmt.Sprintf("ClaimMsg{Username:%v}", msg.Username)
}

func (msg ClaimMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ClaimMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ClaimMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

// Transfer Msg Implementations
func NewTransferMsg(sender string, amount types.LNO, memo []byte, setters ...TransferOption) TransferMsg {
	msg := &TransferMsg{
		Sender: types.AccountKey(sender),
		Amount: amount,
		Memo:   memo,
	}
	for _, setter := range setters {
		setter(msg)
	}
	return *msg
}

func (msg TransferMsg) Type() string { return types.AccountRouterName }

func (msg TransferMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) < types.MinimumUsernameLength ||
		len(msg.Sender) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	// should have either receiver's addr or username
	if len(msg.ReceiverAddr) == 0 && len(msg.ReceiverName) == 0 {
		return ErrInvalidUsername()
	}
	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
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
