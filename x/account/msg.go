package account

// nolint
import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/lino-network/lino/types"
	crypto "github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = FollowMsg{}
var _ sdk.Msg = UnfollowMsg{}
var _ sdk.Msg = ClaimMsg{}
var _ sdk.Msg = TransferMsg{}
var _ sdk.Msg = RecoverMsg{}
var _ sdk.Msg = RegisterMsg{}
var _ sdk.Msg = UpdateAccountMsg{}

// RegisterMsg - bind username with public key, need to be referred by others (pay for it).
type RegisterMsg struct {
	Referrer             types.AccountKey `json:"referrer"`
	RegisterFee          types.LNO        `json:"register_fee"`
	NewUser              types.AccountKey `json:"new_username"`
	NewMasterPubKey      crypto.PubKey    `json:"new_master_public_key"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
	NewPostPubKey        crypto.PubKey    `json:"new_post_public_key"`
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
	NewMasterPubKey      crypto.PubKey    `json:"new_master_public_key"`
	NewPostPubKey        crypto.PubKey    `json:"new_post_public_key"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
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
		return ErrInvalidUsername("illeagle length")
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
		return ErrInvalidUsername("illeagle length")
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

func (msg TransferMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
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

// Recover Msg Implementations
func NewRecoverMsg(
	username string, masterPubkey, transactionPubkey, postPubkey crypto.PubKey) RecoverMsg {
	return RecoverMsg{
		Username:             types.AccountKey(username),
		NewMasterPubKey:      masterPubkey,
		NewTransactionPubKey: transactionPubkey,
		NewPostPubKey:        postPubkey,
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
	return fmt.Sprintf("RecoverMsg{user:%v, new master key:%v, new post Key:%v, new transaction key:%v}",
		msg.Username, msg.NewMasterPubKey, msg.NewPostPubKey, msg.NewTransactionPubKey)
}

func (msg RecoverMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.MasterPermission
	}
	return nil
}

func (msg RecoverMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg RecoverMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

// NewRegisterMsg - construct register msg.
func NewRegisterMsg(
	referrer string,
	newUser string,
	registerFee types.LNO,
	masterPubkey crypto.PubKey,
	transactionPubkey crypto.PubKey,
	postPubkey crypto.PubKey) RegisterMsg {
	return RegisterMsg{
		Referrer:             types.AccountKey(referrer),
		NewUser:              types.AccountKey(newUser),
		RegisterFee:          registerFee,
		NewMasterPubKey:      masterPubkey,
		NewTransactionPubKey: transactionPubkey,
		NewPostPubKey:        postPubkey,
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
	return fmt.Sprintf("RegisterMsg{Newuser:%v, Master Key:%v, Post Key:%v, Transaction Key:%v}",
		msg.NewUser, msg.NewMasterPubKey, msg.NewPostPubKey, msg.NewTransactionPubKey)
}

// Implements Msg.
func (msg RegisterMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	// the permission will not be checked at auth
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
	return nil
}

// Implements Msg.
func (msg RegisterMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg RegisterMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Referrer)}
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
func (msg UpdateAccountMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	// the permission will not be checked at auth
	if keyStr == types.PermissionLevel {
		return types.PostPermission
	}
	return nil
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
func (msg UpdateAccountMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}
