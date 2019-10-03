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
