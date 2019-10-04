package types

// nolint
import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = StakeInMsg{}
var _ types.Msg = StakeOutMsg{}
var _ types.Msg = ClaimInterestMsg{}
var _ types.Msg = StakeInForMsg{}

// StakeInMsg - voter deposit
type StakeInMsg struct {
	Username types.AccountKey `json:"username"`
	Deposit  types.LNO        `json:"deposit"`
}

// StakeOutMsg - voter withdraw
type StakeOutMsg struct {
	Username types.AccountKey `json:"username"`
	Amount   types.LNO        `json:"amount"`
}

// ClaimInterestMsg - claim interest generated from lino power
type ClaimInterestMsg struct {
	Username types.AccountKey `json:"username"`
}

// StakeInForMsg - stake in for other people
type StakeInForMsg struct {
	Sender   types.AccountKey `json:"username"`
	Receiver types.AccountKey `json:"receiver"`
	Deposit  types.LNO        `json:"deposit"`
}

// NewStakeInMsg - return a StakeInMsg
func NewStakeInMsg(username string, deposit types.LNO) StakeInMsg {
	return StakeInMsg{
		Username: types.AccountKey(username),
		Deposit:  deposit,
	}
}

// Route - implements sdk.Msg
func (msg StakeInMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg StakeInMsg) Type() string { return "StakeInMsg" }

// ValidateBasic - implements sdk.Msg
func (msg StakeInMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}

	_, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err
	}
	return nil
}

func (msg StakeInMsg) String() string {
	return fmt.Sprintf("StakeInMsg{Username:%v, Deposit:%v}", msg.Username, msg.Deposit)
}

// GetPermission - implements types.Msg
func (msg StakeInMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg StakeInMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg StakeInMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implement types.Msg
func (msg StakeInMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewStakeOutMsg - return StakeOutMsg
func NewStakeOutMsg(username string, amount types.LNO) StakeOutMsg {
	return StakeOutMsg{
		Username: types.AccountKey(username),
		Amount:   amount,
	}
}

// Route - implements sdk.Msg
func (msg StakeOutMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg StakeOutMsg) Type() string { return "StakeOutMsg" }

// ValidateBasic - implements sdk.Msg
func (msg StakeOutMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}
	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (msg StakeOutMsg) String() string {
	return fmt.Sprintf("StakeOutMsg{Username:%v, Amount:%v}", msg.Username, msg.Amount)
}

// GetPermission - implements types.Msg
func (msg StakeOutMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg StakeOutMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg StakeOutMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implement types.Msg
func (msg StakeOutMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewClaimInterestMsg - return a ClaimInterestMsg
func NewClaimInterestMsg(username string) ClaimInterestMsg {
	return ClaimInterestMsg{
		Username: types.AccountKey(username),
	}
}

// Route - implements sdk.Msg
func (msg ClaimInterestMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg ClaimInterestMsg) Type() string { return "ClaimInterestMsg" }

// ValidateBasic - implements sdk.Msg
func (msg ClaimInterestMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg ClaimInterestMsg) String() string {
	return fmt.Sprintf("ClaimInterestMsg{Username:%v}", msg.Username)
}

// GetPermission - implements types.Msg
func (msg ClaimInterestMsg) GetPermission() types.Permission {
	return types.AppPermission
}

// GetSignBytes - implements sdk.Msg
func (msg ClaimInterestMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg ClaimInterestMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg ClaimInterestMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// NewStakeInForMsg - return a StakeInForMsg
func NewStakeInForMsg(sender, receiver string, deposit types.LNO) StakeInForMsg {
	return StakeInForMsg{
		Sender:   types.AccountKey(sender),
		Receiver: types.AccountKey(receiver),
		Deposit:  deposit,
	}
}

// Route - implements sdk.Msg
func (msg StakeInForMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg StakeInForMsg) Type() string { return "StakeInForMsg" }

// ValidateBasic - implements sdk.Msg
func (msg StakeInForMsg) ValidateBasic() sdk.Error {
	if !msg.Sender.IsValid() || !msg.Receiver.IsValid() {
		return ErrInvalidUsername()
	}

	_, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err
	}
	return nil
}

func (msg StakeInForMsg) String() string {
	return fmt.Sprintf("StakeInForMsg{Sender:%v, Receiver:%v, Deposit:%v}", msg.Sender, msg.Receiver, msg.Deposit)
}

// GetPermission - implements types.Msg
func (msg StakeInForMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg StakeInForMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg
func (msg StakeInForMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Sender)}
}

// GetConsumeAmount - implement types.Msg
func (msg StakeInForMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}
