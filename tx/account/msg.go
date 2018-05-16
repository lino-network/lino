package account

// nolint
import (
	"encoding/json"
	"fmt"

	"github.com/lino-network/lino/types"
	crypto "github.com/tendermint/go-crypto"

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

type RecoverMsg struct {
	Username             types.AccountKey `json:"username"`
	NewMasterPubKey      crypto.PubKey    `json:"new_master_public_key"`
	NewPostPubKey        crypto.PubKey    `json:"new_post_public_key"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
}

type SavingToCheckingMsg struct {
	Username types.AccountKey `json:"username"`
	Amount   types.LNO        `json:"amount"`
}

type CheckingToSavingMsg struct {
	Username types.AccountKey `json:"username"`
	Amount   types.LNO        `json:"amount"`
}

// we can support to transfer to an user or an address
type TransferMsg struct {
	Sender       types.AccountKey `json:"sender"`
	ReceiverName types.AccountKey `json:"receiver_name"`
	ReceiverAddr sdk.Address      `json:"receiver_addr"`
	Amount       types.LNO        `json:"amount"`
	Memo         string           `json:"memo"`
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
var _ sdk.Msg = RecoverMsg{}
var _ sdk.Msg = SavingToCheckingMsg{}
var _ sdk.Msg = CheckingToSavingMsg{}

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
func NewTransferMsg(sender string, amount types.LNO, memo string, setters ...TransferOption) TransferMsg {
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
		return ErrInvalidUsername()
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

// SavingToChecking Msg Implementations
func NewSavingToCheckingMsg(username string, amount types.LNO) SavingToCheckingMsg {
	return SavingToCheckingMsg{
		Username: types.AccountKey(username),
		Amount:   amount,
	}
}

func (msg SavingToCheckingMsg) Type() string { return types.AccountRouterName }

func (msg SavingToCheckingMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.Username) == 0 {
		return ErrInvalidUsername()
	}

	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (msg SavingToCheckingMsg) String() string {
	return fmt.Sprintf("SavingToCheckingMsg{user:%v, amount:%v}", msg.Username, msg.Amount)
}

func (msg SavingToCheckingMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
	return nil
}

func (msg SavingToCheckingMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg SavingToCheckingMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

// CheckingToSaving Msg Implementations
func NewCheckingToSavingMsg(username string, amount types.LNO) CheckingToSavingMsg {
	return CheckingToSavingMsg{
		Username: types.AccountKey(username),
		Amount:   amount,
	}
}

func (msg CheckingToSavingMsg) Type() string { return types.AccountRouterName }

func (msg CheckingToSavingMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.Username) == 0 {
		return ErrInvalidUsername()
	}

	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (msg CheckingToSavingMsg) String() string {
	return fmt.Sprintf("CheckingToSavingMsg{user:%v, amount:%v}", msg.Username, msg.Amount)
}

func (msg CheckingToSavingMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		return types.TransactionPermission
	}
	return nil
}

func (msg CheckingToSavingMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg CheckingToSavingMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}
