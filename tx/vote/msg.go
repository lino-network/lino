package vote

// nolint
import (
	"encoding/json"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/vote/model"
	"github.com/lino-network/lino/types"
)

type VoteMsg struct {
	Voter      types.AccountKey  `json:"voter"`
	ProposalID types.ProposalKey `json:"proposal_id"`
	Result     bool              `json:"result"`
}

type CreateProposalMsg struct {
	Creator types.AccountKey `json:"creator"`
	model.ChangeParameterDescription
}

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
	return fmt.Sprintf("ValidatorDepositMsg{Username:%v, Deposit:%v}", msg.Username, msg.Deposit)
}

func (msg VoterDepositMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg VoterDepositMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
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

func (msg VoterWithdrawMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg VoterWithdrawMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
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

func (msg VoterRevokeMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg VoterRevokeMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
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

func (msg DelegateMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg DelegateMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
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

func (msg RevokeDelegationMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg RevokeDelegationMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg RevokeDelegationMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Delegator)}
}

//----------------------------------------
// VoteMsg Msg Implementations

func NewVoteMsg(voter string, proposalID int64, result bool) VoteMsg {
	return VoteMsg{
		Voter:      types.AccountKey(voter),
		ProposalID: types.ProposalKey(strconv.FormatInt(proposalID, 10)),
		Result:     result,
	}
}

func (msg VoteMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

func (msg VoteMsg) ValidateBasic() sdk.Error {
	if len(msg.Voter) < types.MinimumUsernameLength ||
		len(msg.Voter) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg VoteMsg) String() string {
	return fmt.Sprintf("VoterMsg{Voter:%v, ProposalID:%v, Result:%v}", msg.Voter, msg.ProposalID, msg.Result)
}

func (msg VoteMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg VoteMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg VoteMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Voter)}
}

//----------------------------------------
// CreateProposalMsg Msg Implementations

func NewCreateProposalMsg(voter string, para model.ChangeParameterDescription) CreateProposalMsg {
	return CreateProposalMsg{
		Creator:                    types.AccountKey(voter),
		ChangeParameterDescription: para,
	}
}

func (msg CreateProposalMsg) Type() string { return types.VoteRouterName } // TODO: "account/register"

func (msg CreateProposalMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg CreateProposalMsg) String() string {
	return fmt.Sprintf("CreateProposalMsg{Creator:%v}", msg.Creator)
}

func (msg CreateProposalMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg CreateProposalMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg CreateProposalMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
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

func (msg DelegatorWithdrawMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg DelegatorWithdrawMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg DelegatorWithdrawMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Delegator)}
}
