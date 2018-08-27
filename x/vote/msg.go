package vote

// nolint
import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = VoterDepositMsg{}
var _ types.Msg = VoterWithdrawMsg{}
var _ types.Msg = VoterRevokeMsg{}
var _ types.Msg = DelegateMsg{}
var _ types.Msg = DelegatorWithdrawMsg{}
var _ types.Msg = RevokeDelegationMsg{}

// VoterDepositMsg - voter deposit
type VoterDepositMsg struct {
	Username types.AccountKey `json:"username"`
	Deposit  types.LNO        `json:"deposit"`
}

// VoterWithdrawMsg - voter withdraw
type VoterWithdrawMsg struct {
	Username types.AccountKey `json:"username"`
	Amount   types.LNO        `json:"amount"`
}

// VoterRevokeMsg - voter revoke
type VoterRevokeMsg struct {
	Username types.AccountKey `json:"username"`
}

type DelegateMsg struct {
	Delegator types.AccountKey `json:"delegator"`
	Voter     types.AccountKey `json:"voter"`
	Amount    types.LNO        `json:"amount"`
}

type DelegatorWithdrawMsg struct {
	Delegator types.AccountKey `json:"delegator"`
	Voter     types.AccountKey `json:"voter"`
	Amount    types.LNO        `json:"amount"`
}

type RevokeDelegationMsg struct {
	Delegator types.AccountKey `json:"delegator"`
	Voter     types.AccountKey `json:"voter"`
}

//----------------------------------------
// VoterDepositMsg Msg Implementations

func NewVoterDepositMsg(username string, deposit types.LNO) VoterDepositMsg {
	return VoterDepositMsg{
		Username: types.AccountKey(username),
		Deposit:  deposit,
	}
}

// Type - implements sdk.Msg
func (msg VoterDepositMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

// ValidateBasic - implements sdk.Msg
func (msg VoterDepositMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	_, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err
	}
	return nil
}

func (msg VoterDepositMsg) String() string {
	return fmt.Sprintf("VoterDepositMsg{Username:%v, Deposit:%v}", msg.Username, msg.Deposit)
}

// GetPermission - implements types.Msg
func (msg VoterDepositMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg VoterDepositMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg.
func (msg VoterDepositMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// Implements Msg.
func (msg VoterDepositMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// VoterWithdrawMsg Msg Implementations
func NewVoterWithdrawMsg(username string, amount types.LNO) VoterWithdrawMsg {
	return VoterWithdrawMsg{
		Username: types.AccountKey(username),
		Amount:   amount,
	}
}

// Type - implements sdk.Msg
func (msg VoterWithdrawMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

// ValidateBasic - implements sdk.Msg
func (msg VoterWithdrawMsg) ValidateBasic() sdk.Error {
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

func (msg VoterWithdrawMsg) String() string {
	return fmt.Sprintf("VoterWithdrawMsg{Username:%v, Amount:%v}", msg.Username, msg.Amount)
}

// GetPermission - implements types.Msg
func (msg VoterWithdrawMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg VoterWithdrawMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg.
func (msg VoterWithdrawMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// Implements Msg.
func (msg VoterWithdrawMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// VoterRevokeMsg Msg Implementations

func NewVoterRevokeMsg(username string) VoterRevokeMsg {
	return VoterRevokeMsg{
		Username: types.AccountKey(username),
	}
}

// Type - implements sdk.Msg
func (msg VoterRevokeMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

// ValidateBasic - implements sdk.Msg
func (msg VoterRevokeMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) < types.MinimumUsernameLength ||
		len(msg.Username) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg VoterRevokeMsg) String() string {
	return fmt.Sprintf("VoterRevokeMsg{Username:%v}", msg.Username)
}

// GetPermission - implements types.Msg
func (msg VoterRevokeMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg VoterRevokeMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg.
func (msg VoterRevokeMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// Implements Msg.
func (msg VoterRevokeMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// DelegateMsg Msg Implementations

func NewDelegateMsg(delegator string, voter string, amount types.LNO) DelegateMsg {
	return DelegateMsg{
		Delegator: types.AccountKey(delegator),
		Voter:     types.AccountKey(voter),
		Amount:    amount,
	}
}

// Type - implements sdk.Msg
func (msg DelegateMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

// ValidateBasic - implements sdk.Msg
func (msg DelegateMsg) ValidateBasic() sdk.Error {
	if len(msg.Delegator) < types.MinimumUsernameLength ||
		len(msg.Delegator) > types.MaximumUsernameLength ||
		len(msg.Voter) < types.MinimumUsernameLength ||
		len(msg.Voter) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (msg DelegateMsg) String() string {
	return fmt.Sprintf("DelegateMsg{Delegator:%v, Voter:%v, Amount:%v}", msg.Delegator, msg.Voter, msg.Amount)
}

// GetPermission - implements types.Msg
func (msg DelegateMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg DelegateMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg.
func (msg DelegateMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Delegator)}
}

// Implements Msg.
func (msg DelegateMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// RevokeDelegation Msg Implementations

func NewRevokeDelegationMsg(delegator string, voter string) RevokeDelegationMsg {
	return RevokeDelegationMsg{
		Delegator: types.AccountKey(delegator),
		Voter:     types.AccountKey(voter),
	}
}

// Type - implements sdk.Msg
func (msg RevokeDelegationMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

// ValidateBasic - implements sdk.Msg
func (msg RevokeDelegationMsg) ValidateBasic() sdk.Error {
	if len(msg.Delegator) < types.MinimumUsernameLength ||
		len(msg.Delegator) > types.MaximumUsernameLength ||
		len(msg.Voter) < types.MinimumUsernameLength ||
		len(msg.Voter) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	return nil
}

func (msg RevokeDelegationMsg) String() string {
	return fmt.Sprintf("RevokeDelegationMsg{Delegator:%v, Voter:%v}", msg.Delegator, msg.Voter)
}

// GetPermission - implements types.Msg
func (msg RevokeDelegationMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg RevokeDelegationMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg.
func (msg RevokeDelegationMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Delegator)}
}

// Implements Msg.
func (msg RevokeDelegationMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// DelegatoWithdrawMsg Msg Implementations
func NewDelegatorWithdrawMsg(delegator string, voter string, amount types.LNO) DelegatorWithdrawMsg {
	return DelegatorWithdrawMsg{
		Delegator: types.AccountKey(delegator),
		Voter:     types.AccountKey(voter),
		Amount:    amount,
	}
}

// Type - implements sdk.Msg
func (msg DelegatorWithdrawMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

// ValidateBasic - implements sdk.Msg
func (msg DelegatorWithdrawMsg) ValidateBasic() sdk.Error {
	if len(msg.Delegator) < types.MinimumUsernameLength ||
		len(msg.Delegator) > types.MaximumUsernameLength ||
		len(msg.Voter) < types.MinimumUsernameLength ||
		len(msg.Voter) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (msg DelegatorWithdrawMsg) String() string {
	return fmt.Sprintf("DelegatorWithdrawMsg{Delegator:%v, Voter:%v, Amount:%v}", msg.Delegator, msg.Voter, msg.Amount)
}

// GetPermission - implements types.Msg
func (msg DelegatorWithdrawMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg DelegatorWithdrawMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg.
func (msg DelegatorWithdrawMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Delegator)}
}

// Implements Msg.
func (msg DelegatorWithdrawMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}
