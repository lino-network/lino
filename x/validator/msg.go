package validator

// nolint
import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/tendermint/crypto"
)

var _ types.Msg = ValidatorDepositMsg{}
var _ types.Msg = ValidatorWithdrawMsg{}
var _ types.Msg = ValidatorRevokeMsg{}

// ValidatorDepositMsg - deposit to become validator or add deposit
type ValidatorDepositMsg struct {
	Username  types.AccountKey `json:"username"`
	Deposit   types.LNO        `json:"deposit"`
	ValPubKey crypto.PubKey    `json:"validator_public_key"`
	Link      string           `json:"link"`
}

// ValidatorWithdrawMsg - withdraw validator deposit
type ValidatorWithdrawMsg struct {
	Username types.AccountKey `json:"username"`
	Amount   types.LNO        `json:"amount"`
}

// ValidatorRevokeMsg - revoke validator
type ValidatorRevokeMsg struct {
	Username types.AccountKey `json:"username"`
}

// ValidatorDepositMsg Msg Implementations
func NewValidatorDepositMsg(validator string, deposit types.LNO, pubKey crypto.PubKey, link string) ValidatorDepositMsg {
	return ValidatorDepositMsg{
		Username:  types.AccountKey(validator),
		Deposit:   deposit,
		ValPubKey: pubKey,
		Link:      link,
	}
}

// Type - implement sdk.Msg
func (msg ValidatorDepositMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

// ValidateBasic - implement sdk.Msg
func (msg ValidatorDepositMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if len(msg.Link) > types.MaximumLinkURL {
		return ErrInvalidWebsite()
	}

	_, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err
	}

	return nil
}

func (msg ValidatorDepositMsg) String() string {
	return fmt.Sprintf("ValidatorDepositMsg{Username:%v, Deposit:%v, PubKey:%v}", msg.Username, msg.Deposit, msg.ValPubKey)
}

// GetPermission - implement types.Msg
func (msg ValidatorDepositMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ValidatorDepositMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ValidatorDepositMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implement types.Msg
func (msg ValidatorDepositMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// ValidatorWithdrawMsg Msg Implementations
func NewValidatorWithdrawMsg(validator string, amount types.LNO) ValidatorWithdrawMsg {
	return ValidatorWithdrawMsg{
		Username: types.AccountKey(validator),
		Amount:   amount,
	}
}

// Type - implement sdk.Msg
func (msg ValidatorWithdrawMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

// ValidateBasic - implement sdk.Msg
func (msg ValidatorWithdrawMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (msg ValidatorWithdrawMsg) String() string {
	return fmt.Sprintf("ValidatorWithdrawMsg{Username:%v}", msg.Username)
}

// GetPermission - implement types.Msg
func (msg ValidatorWithdrawMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ValidatorWithdrawMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ValidatorWithdrawMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// Implements Msg.
func (msg ValidatorWithdrawMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// ValidatorRevokeMsg Msg Implementations
func NewValidatorRevokeMsg(validator string) ValidatorRevokeMsg {
	return ValidatorRevokeMsg{
		Username: types.AccountKey(validator),
	}
}

// Type - implement sdk.Msg
func (msg ValidatorRevokeMsg) Type() string { return types.ValidatorRouterName } // TODO: "account/register"

// ValidateBasic - implement sdk.Msg
func (msg ValidatorRevokeMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg ValidatorRevokeMsg) String() string {
	return fmt.Sprintf("ValidatorRevokeMsg{Username:%v}", msg.Username)
}

// GetPermission - implement types.Msg
func (msg ValidatorRevokeMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ValidatorRevokeMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ValidatorRevokeMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implement types.Msg
func (msg ValidatorRevokeMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}
