package types

// nolint
import (
	"fmt"

	"github.com/lino-network/lino/types"
	crypto "github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

// RecoverMsg - replace three public keys
type RecoverMsg struct {
	Username         types.AccountKey `json:"username"`
	NewTxPubKey      crypto.PubKey    `json:"new_tx_public_key"`
	NewSigningPubKey crypto.PubKey    `json:"new_signing_public_key"`
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
func (msg TransferMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg TransferMsg) Type() string { return TransferMsgType }

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
func (msg RecoverMsg) Type() string { return RecoverMsgType }

// ValidateBasic - implements sdk.Msg
func (msg RecoverMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsUsername() {
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
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg RecoverMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username), sdk.AccAddress(msg.NewTxPubKey.Address())}
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
func (msg RegisterMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg RegisterMsg) Type() string { return RegisterMsgType }

// ValidateBasic - implements sdk.Msg
func (msg RegisterMsg) ValidateBasic() sdk.Error {
	if !msg.NewUser.IsUsername() {
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
func (msg UpdateAccountMsg) Type() string { return UpdateAccountMsgType }

// ValidateBasic - implements sdk.Msg
func (msg UpdateAccountMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsUsername() {
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
