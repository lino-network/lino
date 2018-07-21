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
	Referrer              types.AccountKey `json:"referrer"`
	RegisterFee           types.LNO        `json:"register_fee"`
	NewUser               types.AccountKey `json:"new_username"`
	NewRecoveryPubKey     crypto.PubKey    `json:"new_recovery_public_key"`
	NewTransactionPubKey  crypto.PubKey    `json:"new_transaction_public_key"`
	NewMicropaymentPubKey crypto.PubKey    `json:"new_micropayment_public_key"`
	NewPostPubKey         crypto.PubKey    `json:"new_post_public_key"`
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
	Username              types.AccountKey `json:"username"`
	NewRecoveryPubKey     crypto.PubKey    `json:"new_recovery_public_key"`
	NewTransactionPubKey  crypto.PubKey    `json:"new_transaction_public_key"`
	NewMicropaymentPubKey crypto.PubKey    `json:"new_micropayment_public_key"`
	NewPostPubKey         crypto.PubKey    `json:"new_post_public_key"`
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
		return ErrInvalidUsername("illeagle length")
	}
	return nil
}

func (msg FollowMsg) String() string {
	return fmt.Sprintf("FollowMsg{Follower:%v, Followee:%v}", msg.Follower, msg.Followee)
}

// Implements Msg.
func (msg FollowMsg) GetPermission() types.Permission {
	return types.PostPermission
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
		return ErrInvalidUsername("illeagle length")
	}
	return nil
}

func (msg UnfollowMsg) String() string {
	return fmt.Sprintf("UnfollowMsg{Follower:%v, Followee:%v}", msg.Follower, msg.Followee)
}

// Implements Msg.
func (msg UnfollowMsg) GetPermission() types.Permission {
	return types.PostPermission
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
		return ErrInvalidUsername("illeagle length")
	}
	return nil
}

func (msg ClaimMsg) String() string {
	return fmt.Sprintf("ClaimMsg{Username:%v}", msg.Username)
}

func (msg ClaimMsg) GetPermission() types.Permission {
	return types.PostPermission
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
		return ErrInvalidUsername("illeagle length")
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

// Recover Msg Implementations
func NewRecoverMsg(
	username string, recoveryPubkey, transactionPubkey,
	micropaymentPubkey, postPubkey crypto.PubKey) RecoverMsg {
	return RecoverMsg{
		Username:              types.AccountKey(username),
		NewRecoveryPubKey:     recoveryPubkey,
		NewTransactionPubKey:  transactionPubkey,
		NewMicropaymentPubKey: micropaymentPubkey,
		NewPostPubKey:         postPubkey,
	}
}

func (msg RecoverMsg) Type() string { return types.AccountRouterName }

func (msg RecoverMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername("illeagle length")
	}

	return nil
}

func (msg RecoverMsg) String() string {
	return fmt.Sprintf("RecoverMsg{user:%v, new recovery key:%v, new post Key:%v, new transaction key:%v}",
		msg.Username, msg.NewRecoveryPubKey, msg.NewPostPubKey, msg.NewTransactionPubKey)
}

func (msg RecoverMsg) GetPermission() types.Permission {
	return types.RecoveryPermission
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

// NewRegisterMsg - construct register msg.
func NewRegisterMsg(
	referrer string, newUser string, registerFee types.LNO,
	recoveryPubkey, transactionPubkey, micropaymentPubkey,
	postPubkey crypto.PubKey) RegisterMsg {
	return RegisterMsg{
		Referrer:              types.AccountKey(referrer),
		NewUser:               types.AccountKey(newUser),
		RegisterFee:           registerFee,
		NewRecoveryPubKey:     recoveryPubkey,
		NewTransactionPubKey:  transactionPubkey,
		NewMicropaymentPubKey: micropaymentPubkey,
		NewPostPubKey:         postPubkey,
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
		return ErrInvalidUsername("illeagle length")
	}

	match, err := regexp.MatchString(types.UsernameReCheck, string(msg.NewUser))
	if err != nil {
		return ErrInvalidUsername("match error")
	}
	if !match {
		return ErrInvalidUsername("illeagle input")
	}

	_, coinErr := types.LinoToCoin(msg.RegisterFee)
	if coinErr != nil {
		return coinErr
	}
	return nil
}

func (msg RegisterMsg) String() string {
	return fmt.Sprintf("RegisterMsg{Newuser:%v, Recovery Key:%v, Post Key:%v, Transaction Key:%v}",
		msg.NewUser, msg.NewRecoveryPubKey, msg.NewPostPubKey, msg.NewTransactionPubKey)
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
		return ErrInvalidUsername("illeagle length")
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
	return types.PostPermission
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
