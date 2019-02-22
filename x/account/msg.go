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

// RegisterMsg - bind username with public key, need to be referred by others (pay for it)
type RegisterMsg struct {
	Referrer             types.AccountKey `json:"referrer"`
	RegisterFee          types.LNO        `json:"register_fee"`
	NewUser              types.AccountKey `json:"new_username"`
	NewResetPubKey       crypto.PubKey    `json:"new_reset_public_key"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
	NewAppPubKey         crypto.PubKey    `json:"new_app_public_key"`
}

// FollowMsg - follower follow followee
type FollowMsg struct {
	Follower types.AccountKey `json:"follower"`
	Followee types.AccountKey `json:"followee"`
}

// UnfollowMsg - follower unfollow followee
type UnfollowMsg struct {
	Follower types.AccountKey `json:"follower"`
	Followee types.AccountKey `json:"followee"`
}

// ClaimMsg - claim content reward
type ClaimMsg struct {
	Username types.AccountKey `json:"username"`
}

// RecoverMsg - replace three public keys
type RecoverMsg struct {
	Username             types.AccountKey `json:"username"`
	NewResetPubKey       crypto.PubKey    `json:"new_reset_public_key"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
	NewAppPubKey         crypto.PubKey    `json:"new_app_public_key"`
}

// TransferMsg - sender transfer money to receiver
type TransferMsg struct {
	Sender   types.AccountKey `json:"sender"`
	Receiver types.AccountKey `json:"receiver"`
	Amount   types.LNO        `json:"amount"`
	Memo     string           `json:"memo"`
}

// UpdateAccountMsg - update account JSON meta info
type UpdateAccountMsg struct {
	Username types.AccountKey `json:"username"`
	JSONMeta string           `json:"json_meta"`
}

// NewFollowMsg - return a FollowMsg
func NewFollowMsg(follower string, followee string) FollowMsg {
	return FollowMsg{
		Follower: types.AccountKey(follower),
		Followee: types.AccountKey(followee),
	}
}

// Route - implements sdk.Msg
func (msg FollowMsg) Route() string { return types.AccountRouterName }

// Type - implements sdk.Msg
func (msg FollowMsg) Type() string { return "FollowMsg" }

// ValidateBasic - implements sdk.Msg
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

// GetPermission - implements types.Msg
func (msg FollowMsg) GetPermission() types.Permission {
	return types.AppPermission
}

// GetSignBytes - implements sdk.Msg
func (msg FollowMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg FollowMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Follower)}
}

// GetConsumeAmount - implements types.Msg
func (msg FollowMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewUnfollowMsg - return a UnfollowMsg
func NewUnfollowMsg(follower string, followee string) UnfollowMsg {
	return UnfollowMsg{
		Follower: types.AccountKey(follower),
		Followee: types.AccountKey(followee),
	}
}

// Route - implements sdk.Msg
func (msg UnfollowMsg) Route() string { return types.AccountRouterName }

// Type - implements sdk.Msg
func (msg UnfollowMsg) Type() string { return "UnfollowMsg" }

// ValidateBasic - implements sdk.Msg
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

// GetPermission - implements types.Msg
func (msg UnfollowMsg) GetPermission() types.Permission {
	return types.AppPermission
}

// GetSignBytes - implements sdk.Msg
func (msg UnfollowMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg UnfollowMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Follower)}
}

// GetConsumeAmount - implements types.Msg
func (msg UnfollowMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewClaimMsg - return a ClaimMsg
func NewClaimMsg(username string) ClaimMsg {
	return ClaimMsg{
		Username: types.AccountKey(username),
	}
}

// Route - implements sdk.Msg
func (msg ClaimMsg) Route() string { return types.AccountRouterName }

// Type - implements sdk.Msg
func (msg ClaimMsg) Type() string { return "ClaimMsg" }

// ValidateBasic - implements sdk.Msg
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

// GetPermission - implements types.Msg
func (msg ClaimMsg) GetPermission() types.Permission {
	return types.AppPermission
}

// GetSignBytes - implements sdk.Msg
func (msg ClaimMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg ClaimMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg ClaimMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewTransferMsg - return a TransferMsg
func NewTransferMsg(sender, receiver string, amount types.LNO, memo string) TransferMsg {
	return TransferMsg{
		Sender:   types.AccountKey(sender),
		Amount:   amount,
		Memo:     memo,
		Receiver: types.AccountKey(receiver),
	}
}

// Route - implements sdk.Msg
func (msg TransferMsg) Route() string { return types.AccountRouterName }

// Type - implements sdk.Msg
func (msg TransferMsg) Type() string { return "TransferMsg" }

// ValidateBasic - implements sdk.Msg
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

// GetPermission - implements types.Msg
func (msg TransferMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg TransferMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg TransferMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Sender)}
}

// GetConsumeAmount - implements types.Msg
func (msg TransferMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewRecoverMsg - return a recover msg
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

// Route - implements sdk.Msg
func (msg RecoverMsg) Route() string { return types.AccountRouterName }

// Type - implements sdk.Msg
func (msg RecoverMsg) Type() string { return "RecoverMsg" }

// ValidateBasic - implements sdk.Msg
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

// GetPermission - implements types.Msg
func (msg RecoverMsg) GetPermission() types.Permission {
	return types.ResetPermission
}

// GetSignBytes - implements sdk.Msg
func (msg RecoverMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg RecoverMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
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

// Route - implements sdk.Msg
func (msg RegisterMsg) Route() string { return types.AccountRouterName }

// Type - implements sdk.Msg
func (msg RegisterMsg) Type() string { return "RegisterMsg" }

// ValidateBasic - implements sdk.Msg
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

	match, err = regexp.MatchString(types.IllegalUsernameReCheck, string(msg.NewUser))
	if err != nil {
		return ErrInvalidUsername("match error")
	}
	if match {
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

// GetSignBytes - implements sdk.Msg
func (msg RegisterMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetPermission - implements types.Msg
func (msg RegisterMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSigners - implements sdk.Msg
func (msg RegisterMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Referrer)}
}

// GetConsumeAmount - implements types.Msg
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

// Type - implements sdk.Msg
func (msg UpdateAccountMsg) Route() string { return types.AccountRouterName }

// Type - implements sdk.Msg
func (msg UpdateAccountMsg) Type() string { return "UpdateAccountMsg" }

// ValidateBasic - implements sdk.Msg
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

// GetPermission - implements types.Msg
func (msg UpdateAccountMsg) GetPermission() types.Permission {
	return types.AppPermission
}

// GetSignBytes - implements sdk.Msg
func (msg UpdateAccountMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg UpdateAccountMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg UpdateAccountMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}
