package account

// nolint
import (
	"fmt"
	"regexp"

	"github.com/lino-network/lino/types"
	crypto "github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = FollowMsg{}
var _ types.Msg = UnfollowMsg{}
var _ types.Msg = ClaimMsg{}
var _ types.Msg = TransferMsg{}
var _ types.Msg = RecoverMsg{}
var _ types.Msg = RegisterMsg{}
var _ types.Msg = UpdateAccountMsg{}

// RegisterMsg - bind username with public key, need to be referred by others (pay for it).
type RegisterMsg struct {
	Referrer             types.AccountKey `json:"referrer"`
	RegisterFee          types.LNO        `json:"register_fee"`
	NewUser              types.AccountKey `json:"new_username"`
	NewResetPubKey       crypto.PubKey    `json:"new_reset_public_key"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
	NewAppPubKey         crypto.PubKey    `json:"new_app_public_key"`
}

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

type RecoverMsg struct {
	Username             types.AccountKey `json:"username"`
	NewResetPubKey       crypto.PubKey    `json:"new_reset_public_key"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
	NewAppPubKey         crypto.PubKey    `json:"new_app_public_key"`
}

// we can support to transfer to an user or an address
type TransferMsg struct {
	Sender   types.AccountKey `json:"sender"`
	Receiver types.AccountKey `json:"receiver"`
	Amount   types.LNO        `json:"amount"`
	Memo     string           `json:"memo"`
}

// UpdateAccountMsg - update account JSON meta info.
type UpdateAccountMsg struct {
	Username types.AccountKey `json:"username"`
	JSONMeta string           `json:"json_meta"`
}

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
		return ErrInvalidUsername("illegal length")
	}
	return nil
}

func (msg FollowMsg) String() string {
	return fmt.Sprintf("FollowMsg{Follower:%v, Followee:%v}", msg.Follower, msg.Followee)
}

// Implements Msg.
func (msg FollowMsg) GetPermission() types.Permission {
	return types.AppPermission
}

func (msg FollowMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg FollowMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Follower)}
}

// Implements Msg.
func (msg FollowMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
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
		return ErrInvalidUsername("illegal length")
	}
	return nil
}

func (msg UnfollowMsg) String() string {
	return fmt.Sprintf("UnfollowMsg{Follower:%v, Followee:%v}", msg.Follower, msg.Followee)
}

// Implements Msg.
func (msg UnfollowMsg) GetPermission() types.Permission {
	return types.AppPermission
}

func (msg UnfollowMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg UnfollowMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Follower)}
}

// Implements Msg.
func (msg UnfollowMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
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
		return ErrInvalidUsername("illegal length")
	}
	return nil
}

func (msg ClaimMsg) String() string {
	return fmt.Sprintf("ClaimMsg{Username:%v}", msg.Username)
}

func (msg ClaimMsg) GetPermission() types.Permission {
	return types.AppPermission
}

func (msg ClaimMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ClaimMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// Implements Msg.
func (msg ClaimMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// Transfer Msg Implementations
func NewTransferMsg(sender, receiver string, amount types.LNO, memo string) TransferMsg {
	return TransferMsg{
		Sender:   types.AccountKey(sender),
		Amount:   amount,
		Memo:     memo,
		Receiver: types.AccountKey(receiver),
	}
}

func (msg TransferMsg) Type() string { return types.AccountRouterName }

func (msg TransferMsg) ValidateBasic() sdk.Error {
	if len(msg.Sender) < types.MinimumUsernameLength ||
		len(msg.Sender) > types.MaximumUsernameLength ||
		len(msg.Receiver) < types.MinimumUsernameLength ||
		len(msg.Receiver) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}
	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}

	if len(msg.Memo) > types.MaximumMemoLength {
		return ErrInvalidMemo()
	}
	return nil
}

func (msg TransferMsg) String() string {
	return fmt.Sprintf("TransferMsg{Sender:%v, Receiver:%v, Amount:%v, Memo:%v}",
		msg.Sender, msg.Receiver, msg.Amount, msg.Memo)
}

func (msg TransferMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg TransferMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg TransferMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Sender)}
}

// Implements Msg.
func (msg TransferMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// Recover Msg Implementations
func NewRecoverMsg(
	username string, resetPubkey, transactionPubkey,
	appPubkey crypto.PubKey) RecoverMsg {
	return RecoverMsg{
		Username:             types.AccountKey(username),
		NewResetPubKey:       resetPubkey,
		NewTransactionPubKey: transactionPubkey,
		NewAppPubKey:         appPubkey,
	}
}

func (msg RecoverMsg) Type() string { return types.AccountRouterName }

func (msg RecoverMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}

	return nil
}

func (msg RecoverMsg) String() string {
	return fmt.Sprintf("RecoverMsg{user:%v, new reset key:%v, new app Key:%v, new transaction key:%v}",
		msg.Username, msg.NewResetPubKey, msg.NewAppPubKey, msg.NewTransactionPubKey)
}

func (msg RecoverMsg) GetPermission() types.Permission {
	return types.ResetPermission
}

func (msg RecoverMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg RecoverMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// Implements Msg.
func (msg RecoverMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewRegisterMsg - construct register msg.
func NewRegisterMsg(
	referrer string, newUser string, registerFee types.LNO,
	resetPubkey, transactionPubkey, appPubkey crypto.PubKey) RegisterMsg {
	return RegisterMsg{
		Referrer:             types.AccountKey(referrer),
		NewUser:              types.AccountKey(newUser),
		RegisterFee:          registerFee,
		NewResetPubKey:       resetPubkey,
		NewTransactionPubKey: transactionPubkey,
		NewAppPubKey:         appPubkey,
	}
}

// Implements Msg.
func (msg RegisterMsg) Type() string { return types.AccountRouterName } // TODO: "account/register"

// Implements Msg.
func (msg RegisterMsg) ValidateBasic() sdk.Error {
	if len(msg.NewUser) < types.MinimumUsernameLength ||
		len(msg.NewUser) > types.MaximumUsernameLength ||
		len(msg.Referrer) < types.MinimumUsernameLength ||
		len(msg.Referrer) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}

	match, err := regexp.MatchString(types.UsernameReCheck, string(msg.NewUser))
	if err != nil {
		return ErrInvalidUsername("match error")
	}
	if !match {
		return ErrInvalidUsername("illegal input")
	}

	_, coinErr := types.LinoToCoin(msg.RegisterFee)
	if coinErr != nil {
		return coinErr
	}
	return nil
}

func (msg RegisterMsg) String() string {
	return fmt.Sprintf("RegisterMsg{Newuser:%v, Reset Key:%v, App Key:%v, Transaction Key:%v}",
		msg.NewUser, msg.NewResetPubKey, msg.NewAppPubKey, msg.NewTransactionPubKey)
}

// Implements Msg.
func (msg RegisterMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg RegisterMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// Implements Msg.
func (msg RegisterMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Referrer)}
}

// Implements Msg.
func (msg RegisterMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewUpdateAccountMsg - construct user update msg to update user JSON meta info.
func NewUpdateAccountMsg(username string, JSONMeta string) UpdateAccountMsg {
	return UpdateAccountMsg{
		Username: types.AccountKey(username),
		JSONMeta: JSONMeta,
	}
}

// Implements Msg.
func (msg UpdateAccountMsg) Type() string { return types.AccountRouterName } // TODO: "account/register"

// Implements Msg.
func (msg UpdateAccountMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illegal length")
	}

	if len(msg.JSONMeta) > types.MaximumJSONMetaLength {
		return ErrInvalidJSONMeta()
	}

	return nil
}

func (msg UpdateAccountMsg) String() string {
	return fmt.Sprintf("UpdateAccountMsg{User:%v, JSON meta:%v}", msg.Username, msg.JSONMeta)
}

// Implements Msg.
func (msg UpdateAccountMsg) GetPermission() types.Permission {
	return types.AppPermission
}

// Implements Msg.
func (msg UpdateAccountMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg UpdateAccountMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// Implements Msg.
func (msg UpdateAccountMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}
