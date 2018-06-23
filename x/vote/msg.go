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

type VoterDepositMsg struct {
	Username types.AccountKey `json:"username"`
	Deposit  types.LNO        `json:"deposit"`
}

type VoterWithdrawMsg struct {
	Username types.AccountKey `json:"username"`
	Amount   types.LNO        `json:"amount"`
}

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

func (msg VoterDepositMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

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

func (msg VoterDepositMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg VoterDepositMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg VoterDepositMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

//----------------------------------------
// VoterWithdrawMsg Msg Implementations
func NewVoterWithdrawMsg(username string, amount types.LNO) VoterWithdrawMsg {
	return VoterWithdrawMsg{
		Username: types.AccountKey(username),
		Amount:   amount,
	}
}

func (msg VoterWithdrawMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

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

func (msg VoterWithdrawMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg VoterWithdrawMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg VoterWithdrawMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

//----------------------------------------
// VoterRevokeMsg Msg Implementations

func NewVoterRevokeMsg(username string) VoterRevokeMsg {
	return VoterRevokeMsg{
		Username: types.AccountKey(username),
	}
}

func (msg VoterRevokeMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

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

func (msg VoterRevokeMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg VoterRevokeMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg VoterRevokeMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
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

func (msg DelegateMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

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

func (msg DelegateMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg DelegateMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg DelegateMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Delegator)}
}

//----------------------------------------
// RevokeDelegation Msg Implementations

func NewRevokeDelegationMsg(delegator string, voter string) RevokeDelegationMsg {
	return RevokeDelegationMsg{
		Delegator: types.AccountKey(delegator),
		Voter:     types.AccountKey(voter),
	}
}

func (msg RevokeDelegationMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

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

func (msg RevokeDelegationMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg RevokeDelegationMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg RevokeDelegationMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Delegator)}
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

func (msg DelegatorWithdrawMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

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

func (msg DelegatorWithdrawMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg DelegatorWithdrawMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg DelegatorWithdrawMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Delegator)}
}
