package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
)

// TransferMsg - sender transfer money to receiver
type TransferMsg struct {
	Sender   types.AccountKey `json:"sender"`
	Receiver types.AccountKey `json:"receiver"`
	Amount   types.LNO        `json:"amount"`
	Memo     string           `json:"memo"`
}

var _ types.Msg = TransferMsg{}

// NewTransferMsg - return a TransferMsg
func NewTransferMsg(sender, receiver string, amount types.LNO, memo string) TransferMsg {
	return TransferMsg{
		Sender:   types.AccountKey(sender),
		Receiver: types.AccountKey(receiver),
		Amount:   amount,
		Memo:     memo,
	}
}

// Route - implements sdk.Msg
func (msg TransferMsg) Route() string { return RouterKey }

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

// TransferV2Msg - support account/addr to account/addr
type TransferV2Msg struct {
	Sender   types.AccOrAddr `json:"sender"`
	Receiver types.AccOrAddr `json:"receiver"`
	Amount   types.LNO       `json:"amount"`
	Memo     string          `json:"memo"`
}

var _ types.AddrMsg = TransferV2Msg{}

// NewTransferV2Msg - return a TransferV2Msg
func NewTransferV2Msg(sender, receiver types.AccOrAddr, amount types.LNO, memo string) TransferV2Msg {
	return TransferV2Msg{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
		Memo:     memo,
	}
}

// Route - implements sdk.Msg
func (msg TransferV2Msg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg TransferV2Msg) Type() string { return "TransferV2Msg" }

// ValidateBasic - implements sdk.Msg
func (msg TransferV2Msg) ValidateBasic() sdk.Error {
	if !msg.Sender.IsValid() {
		return ErrInvalidUsername(msg.Sender.String())
	}
	if !msg.Receiver.IsValid() {
		return ErrInvalidUsername(msg.Receiver.String())
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

func (msg TransferV2Msg) String() string {
	return fmt.Sprintf("TransferV2Msg{Sender:%s,Receiver:%s,Amount:%s,Memo:%s}",
		msg.Sender, msg.Receiver, msg.Amount, msg.Memo)
}

// GetSignBytes - implements sdk.Msg
func (msg TransferV2Msg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
// SHOULD NOT BE USED.
func (msg TransferV2Msg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Sender.String())}
}

// GetAccOrAddrSigners - implements types.AddrMsg
func (msg TransferV2Msg) GetAccOrAddrSigners() []types.AccOrAddr {
	return []types.AccOrAddr{msg.Sender}
}

// RecoverMsg - replace two keys
type RecoverMsg struct {
	Username         types.AccountKey `json:"username"`
	NewTxPubKey      crypto.PubKey    `json:"new_tx_public_key"`
	NewSigningPubKey crypto.PubKey    `json:"new_signing_public_key"`
}

var _ types.Msg = RecoverMsg{}

// NewRecoverMsg - return a recover msg
func NewRecoverMsg(
	username string, transactionPubkey, signingPubkey crypto.PubKey) RecoverMsg {
	return RecoverMsg{
		Username:         types.AccountKey(username),
		NewTxPubKey:      transactionPubkey,
		NewSigningPubKey: signingPubkey,
	}
}

// Route - implements sdk.Msg
func (msg RecoverMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg RecoverMsg) Type() string { return "RecoverMsg" }

// ValidateBasic - implements sdk.Msg
func (msg RecoverMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername("illegal username")
	}

	return nil
}

func (msg RecoverMsg) String() string {
	return fmt.Sprintf(
		"RecoverMsg{user:%v, new tx key:%v, new signing Key:%v}",
		msg.Username, msg.NewTxPubKey, msg.NewSigningPubKey)
}

// GetPermission - implements types.Msg
func (msg RecoverMsg) GetPermission() types.Permission {
	return types.ResetPermission
}

// GetSignBytes - implements sdk.Msg
func (msg RecoverMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
// SHOULD NOT BE USED.
func (msg RecoverMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username), sdk.AccAddress(msg.NewTxPubKey.Address())}
}

// GetAccOrAddrSigners - implements types.AddrMsg
func (msg RecoverMsg) GetAccOrAddrSigners() []types.AccOrAddr {
	return []types.AccOrAddr{
		types.NewAccOrAddrFromAcc(msg.Username),
		types.NewAccOrAddrFromAddr(sdk.AccAddress(msg.NewTxPubKey.Address()))}
}

// GetConsumeAmount - implements types.Msg
func (msg RecoverMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// RegisterMsg - bind username with public key, need to be referred by others (pay for it)
type RegisterMsg struct {
	Referrer             types.AccountKey `json:"referrer"`
	RegisterFee          types.LNO        `json:"register_fee"`
	NewUser              types.AccountKey `json:"new_username"`
	NewResetPubKey       crypto.PubKey    `json:"new_reset_public_key"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
	NewAppPubKey         crypto.PubKey    `json:"new_app_public_key"`
}

var _ types.Msg = RegisterMsg{}

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
func (msg RegisterMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg RegisterMsg) Type() string { return "RegisterMsg" }

// ValidateBasic - implements sdk.Msg
func (msg RegisterMsg) ValidateBasic() sdk.Error {
	if !msg.NewUser.IsValid() {
		return ErrInvalidUsername("illegal username")
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

// UpdateAccountMsg - update account JSON meta info
type UpdateAccountMsg struct {
	Username types.AccountKey `json:"username"`
	JSONMeta string           `json:"json_meta"`
}

var _ types.Msg = UpdateAccountMsg{}

// NewUpdateAccountMsg - construct user update msg to update user JSON meta info.
func NewUpdateAccountMsg(username string, jsonMeta string) UpdateAccountMsg {
	return UpdateAccountMsg{
		Username: types.AccountKey(username),
		JSONMeta: jsonMeta,
	}
}

// Type - implements sdk.Msg
func (msg UpdateAccountMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg UpdateAccountMsg) Type() string { return "UpdateAccountMsg" }

// ValidateBasic - implements sdk.Msg
func (msg UpdateAccountMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername("illegal username")
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

// RegisterV2Msg - bind username with public key, need to be referred by others (pay for it)
type RegisterV2Msg struct {
	Referrer             types.AccOrAddr  `json:"referrer"`
	RegisterFee          types.LNO        `json:"register_fee"`
	NewUser              types.AccountKey `json:"new_username"`
	NewTransactionPubKey crypto.PubKey    `json:"new_transaction_public_key"`
	NewSigningPubKey     crypto.PubKey    `json:"new_signing_public_key"`
}

var _ types.AddrMsg = RegisterV2Msg{}

// NewRegisterV2Msg - construct register msg.
func NewRegisterV2Msg(
	referrer types.AccOrAddr, newUser string, registerFee types.LNO,
	transactionPubkey, signingPubKey crypto.PubKey) RegisterV2Msg {
	return RegisterV2Msg{
		Referrer:             referrer,
		NewUser:              types.AccountKey(newUser),
		RegisterFee:          registerFee,
		NewTransactionPubKey: transactionPubkey,
		NewSigningPubKey:     signingPubKey,
	}
}

// Route - implements sdk.Msg
func (msg RegisterV2Msg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg RegisterV2Msg) Type() string { return "RegisterV2Msg" }

// ValidateBasic - implements sdk.Msg
func (msg RegisterV2Msg) ValidateBasic() sdk.Error {
	if !msg.Referrer.IsValid() {
		return ErrInvalidUsername(msg.Referrer.String())
	}
	if !msg.NewUser.IsValid() {
		return ErrInvalidUsername(string(msg.NewUser))
	}

	_, coinErr := types.LinoToCoin(msg.RegisterFee)
	if coinErr != nil {
		return coinErr
	}
	return nil
}

func (msg RegisterV2Msg) String() string {
	return fmt.Sprintf("RegisterV2Msg{Newuser:%v, Referrer: %s, Tx Key:%v, Signing Key:%v}",
		msg.NewUser, msg.Referrer, msg.NewTransactionPubKey, msg.NewSigningPubKey)
}

// GetSignBytes - implements sdk.Msg
func (msg RegisterV2Msg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
// SHOULD NOT BE USED
func (msg RegisterV2Msg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		sdk.AccAddress(msg.Referrer.String()),
		sdk.AccAddress(msg.NewTransactionPubKey.Address())}
}

// GetAccOrAddrSigners - implements types.AddrMsg
func (msg RegisterV2Msg) GetAccOrAddrSigners() []types.AccOrAddr {
	return []types.AccOrAddr{
		msg.Referrer,
		types.NewAccOrAddrFromAddr(sdk.AccAddress(msg.NewTransactionPubKey.Address()))}
}

// utils
func getSignBytes(msg sdk.Msg) []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}
