package proposal

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
)

type ChangeParamMsg interface {
	GetParameter() param.Parameter
	GetCreator() types.AccountKey
}

type ContentCensorshipMsg interface {
	GetCreator() types.AccountKey
	GetPermLink() types.PermLink
}

type ProtocolUpgradeMsg interface {
	GetCreator() types.AccountKey
	GetLink() string
}

type ChangeGlobalAllocationParamMsg struct {
	Creator   types.AccountKey            `json:"creator"`
	Parameter param.GlobalAllocationParam `json:"parameter"`
}

type ChangeEvaluateOfContentValueParamMsg struct {
	Creator   types.AccountKey                  `json:"creator"`
	Parameter param.EvaluateOfContentValueParam `json:"parameter"`
}

type ChangeInfraInternalAllocationParamMsg struct {
	Creator   types.AccountKey                   `json:"creator"`
	Parameter param.InfraInternalAllocationParam `json:"parameter"`
}

type ChangeVoteParamMsg struct {
	Creator   types.AccountKey `json:"creator"`
	Parameter param.VoteParam  `json:"parameter"`
}

type ChangeProposalParamMsg struct {
	Creator   types.AccountKey    `json:"creator"`
	Parameter param.ProposalParam `json:"parameter"`
}

type ChangeDeveloperParamMsg struct {
	Creator   types.AccountKey     `json:"creator"`
	Parameter param.DeveloperParam `json:"parameter"`
}

type ChangeValidatorParamMsg struct {
	Creator   types.AccountKey     `json:"creator"`
	Parameter param.ValidatorParam `json:"parameter"`
}

type ChangeCoinDayParamMsg struct {
	Creator   types.AccountKey   `json:"creator"`
	Parameter param.CoinDayParam `json:"parameter"`
}

type ChangeBandwidthParamMsg struct {
	Creator   types.AccountKey     `json:"creator"`
	Parameter param.BandwidthParam `json:"parameter"`
}

//----------------------------------------
// ChangeGlobalAllocationParamMsg Msg Implementations

func NewChangeGlobalAllocationParamMsg(creator string, parameter param.GlobalAllocationParam) ChangeGlobalAllocationParamMsg {
	return ChangeGlobalAllocationParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
	}
}

func (msg ChangeGlobalAllocationParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeGlobalAllocationParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeGlobalAllocationParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeGlobalAllocationParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if !msg.Parameter.InfraAllocation.
		Add(msg.Parameter.ContentCreatorAllocation).
		Add(msg.Parameter.DeveloperAllocation).
		Add(msg.Parameter.ValidatorAllocation).Equal(sdk.NewRat(1)) {
		return ErrIllegalParameter()
	}

	return nil
}

func (msg ChangeGlobalAllocationParamMsg) String() string {
	return fmt.Sprintf("ChangeGlobalAllocationParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeGlobalAllocationParamMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeGlobalAllocationParamMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeGlobalAllocationParamMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}

//----------------------------------------
// ChangeEvaluateOfContentValueParamMsg Msg Implementations

func NewChangeEvaluateOfContentValueParamMsg(creator string, parameter param.EvaluateOfContentValueParam) ChangeEvaluateOfContentValueParamMsg {
	return ChangeEvaluateOfContentValueParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
	}
}

func (msg ChangeEvaluateOfContentValueParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeEvaluateOfContentValueParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeEvaluateOfContentValueParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeEvaluateOfContentValueParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	//TODO

	return nil
}

func (msg ChangeEvaluateOfContentValueParamMsg) String() string {
	return fmt.Sprintf("ChangeEvaluateOfContentValueParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeEvaluateOfContentValueParamMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeEvaluateOfContentValueParamMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeEvaluateOfContentValueParamMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}

//----------------------------------------
// ChangeInfraInternalAllocationParamMsg Msg Implementations

func NewChangeInfraInternalAllocationParamMsg(creator string, parameter param.InfraInternalAllocationParam) ChangeInfraInternalAllocationParamMsg {
	return ChangeInfraInternalAllocationParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
	}
}

func (msg ChangeInfraInternalAllocationParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeInfraInternalAllocationParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeInfraInternalAllocationParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeInfraInternalAllocationParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	//TODO

	return nil
}

func (msg ChangeInfraInternalAllocationParamMsg) String() string {
	return fmt.Sprintf("ChangeInfraInternalAllocationParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeInfraInternalAllocationParamMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeInfraInternalAllocationParamMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeInfraInternalAllocationParamMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}

//----------------------------------------
// ChangeVoteParamMsg Msg Implementations

func NewChangeVoteParamMsg(creator string, parameter param.VoteParam) ChangeVoteParamMsg {
	return ChangeVoteParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
	}
}

func (msg ChangeVoteParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeVoteParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeVoteParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeVoteParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	//TODO
	return nil
}

func (msg ChangeVoteParamMsg) String() string {
	return fmt.Sprintf("ChangeVoteParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeVoteParamMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeVoteParamMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeVoteParamMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}

//----------------------------------------
// ChangeProposalParamMsg Msg Implementations

func NewChangeProposalParamMsg(creator string, parameter param.ProposalParam) ChangeProposalParamMsg {
	return ChangeProposalParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
	}
}

func (msg ChangeProposalParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeProposalParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeProposalParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeProposalParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	//TODO

	return nil
}

func (msg ChangeProposalParamMsg) String() string {
	return fmt.Sprintf("ChangeProposalParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeProposalParamMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeProposalParamMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeProposalParamMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}

//----------------------------------------
// ChangeDeveloperParamMsg Msg Implementations

func NewChangeDeveloperParamMsg(creator string, parameter param.DeveloperParam) ChangeDeveloperParamMsg {
	return ChangeDeveloperParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
	}
}

func (msg ChangeDeveloperParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeDeveloperParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeDeveloperParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeDeveloperParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	//TODO
	return nil
}

func (msg ChangeDeveloperParamMsg) String() string {
	return fmt.Sprintf("ChangeDeveloperParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeDeveloperParamMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeDeveloperParamMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeDeveloperParamMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}

//----------------------------------------
// ChangeValidatorParamMsg Msg Implementations

func NewChangeValidatorParamMsg(creator string, parameter param.ValidatorParam) ChangeValidatorParamMsg {
	return ChangeValidatorParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
	}
}

func (msg ChangeValidatorParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeValidatorParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeValidatorParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeValidatorParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	//TODO
	return nil
}

func (msg ChangeValidatorParamMsg) String() string {
	return fmt.Sprintf("ChangeValidatorParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeValidatorParamMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeValidatorParamMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeValidatorParamMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}

//----------------------------------------
// ChangeCoinDayParamMsg Msg Implementations

func NewChangeCoinDayParamMsg(creator string, parameter param.CoinDayParam) ChangeCoinDayParamMsg {
	return ChangeCoinDayParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
	}
}

func (msg ChangeCoinDayParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeCoinDayParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeCoinDayParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeCoinDayParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	//TODO

	return nil
}

func (msg ChangeCoinDayParamMsg) String() string {
	return fmt.Sprintf("ChangeCoinDayParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeCoinDayParamMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeCoinDayParamMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeCoinDayParamMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}

//----------------------------------------
// ChangeBandwidthParamMsg Msg Implementations

func NewChangeBandwidthParamMsg(creator string, parameter param.BandwidthParam) ChangeBandwidthParamMsg {
	return ChangeBandwidthParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
	}
}

func (msg ChangeBandwidthParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeBandwidthParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeBandwidthParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeBandwidthParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	//TODO
	return nil
}

func (msg ChangeBandwidthParamMsg) String() string {
	return fmt.Sprintf("ChangeBandwidthParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeBandwidthParamMsg) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg ChangeBandwidthParamMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeBandwidthParamMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Creator)}
}
