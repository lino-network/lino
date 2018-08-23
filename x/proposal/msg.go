package proposal

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
)

var _ types.Msg = DeletePostContentMsg{}
var _ types.Msg = UpgradeProtocolMsg{}
var _ types.Msg = ChangeGlobalAllocationParamMsg{}
var _ types.Msg = ChangeEvaluateOfContentValueParamMsg{}
var _ types.Msg = ChangeInfraInternalAllocationParamMsg{}
var _ types.Msg = ChangeVoteParamMsg{}
var _ types.Msg = ChangeProposalParamMsg{}
var _ types.Msg = ChangeDeveloperParamMsg{}
var _ types.Msg = ChangeValidatorParamMsg{}
var _ types.Msg = ChangeBandwidthParamMsg{}
var _ types.Msg = ChangeAccountParamMsg{}
var _ types.Msg = ChangePostParamMsg{}
var _ types.Msg = VoteProposalMsg{}

var _ ChangeParamMsg = ChangeGlobalAllocationParamMsg{}
var _ ChangeParamMsg = ChangeEvaluateOfContentValueParamMsg{}
var _ ChangeParamMsg = ChangeInfraInternalAllocationParamMsg{}
var _ ChangeParamMsg = ChangeVoteParamMsg{}
var _ ChangeParamMsg = ChangeProposalParamMsg{}
var _ ChangeParamMsg = ChangeDeveloperParamMsg{}
var _ ChangeParamMsg = ChangeValidatorParamMsg{}
var _ ChangeParamMsg = ChangeBandwidthParamMsg{}
var _ ChangeParamMsg = ChangeAccountParamMsg{}
var _ ChangeParamMsg = ChangePostParamMsg{}

var _ ContentCensorshipMsg = DeletePostContentMsg{}

var _ ProtocolUpgradeMsg = UpgradeProtocolMsg{}

type ChangeParamMsg interface {
	GetParameter() param.Parameter
	GetCreator() types.AccountKey
	GetReason() string
}

type ContentCensorshipMsg interface {
	GetCreator() types.AccountKey
	GetPermlink() types.Permlink
	GetReason() string
}

type ProtocolUpgradeMsg interface {
	GetCreator() types.AccountKey
	GetLink() string
	GetReason() string
}

type DeletePostContentMsg struct {
	Creator  types.AccountKey `json:"creator"`
	Permlink types.Permlink   `json:"permlink"`
	Reason   string           `json:"reason"`
}

type UpgradeProtocolMsg struct {
	Creator types.AccountKey `json:"creator"`
	Link    string           `json:"link"`
	Reason  string           `json:"reason"`
}

type ChangeGlobalAllocationParamMsg struct {
	Creator   types.AccountKey            `json:"creator"`
	Parameter param.GlobalAllocationParam `json:"parameter"`
	Reason    string                      `json:"reason"`
}

type ChangeEvaluateOfContentValueParamMsg struct {
	Creator   types.AccountKey                  `json:"creator"`
	Parameter param.EvaluateOfContentValueParam `json:"parameter"`
	Reason    string                            `json:"reason"`
}

type ChangeInfraInternalAllocationParamMsg struct {
	Creator   types.AccountKey                   `json:"creator"`
	Parameter param.InfraInternalAllocationParam `json:"parameter"`
	Reason    string                             `json:"reason"`
}

type ChangeVoteParamMsg struct {
	Creator   types.AccountKey `json:"creator"`
	Parameter param.VoteParam  `json:"parameter"`
	Reason    string           `json:"reason"`
}

type ChangeProposalParamMsg struct {
	Creator   types.AccountKey    `json:"creator"`
	Parameter param.ProposalParam `json:"parameter"`
	Reason    string              `json:"reason"`
}

type ChangeDeveloperParamMsg struct {
	Creator   types.AccountKey     `json:"creator"`
	Parameter param.DeveloperParam `json:"parameter"`
	Reason    string               `json:"reason"`
}

type ChangeValidatorParamMsg struct {
	Creator   types.AccountKey     `json:"creator"`
	Parameter param.ValidatorParam `json:"parameter"`
	Reason    string               `json:"reason"`
}

type ChangeBandwidthParamMsg struct {
	Creator   types.AccountKey     `json:"creator"`
	Parameter param.BandwidthParam `json:"parameter"`
	Reason    string               `json:"reason"`
}

type ChangeAccountParamMsg struct {
	Creator   types.AccountKey   `json:"creator"`
	Parameter param.AccountParam `json:"parameter"`
	Reason    string             `json:"reason"`
}

type ChangePostParamMsg struct {
	Creator   types.AccountKey `json:"creator"`
	Parameter param.PostParam  `json:"parameter"`
	Reason    string           `json:"reason"`
}

type VoteProposalMsg struct {
	Voter      types.AccountKey  `json:"voter"`
	ProposalID types.ProposalKey `json:"proposal_id"`
	Result     bool              `json:"result"`
}

//----------------------------------------
// ChangeGlobalAllocationParamMsg Msg Implementations

func NewDeletePostContentMsg(
	creator string, permlink types.Permlink, reason string) DeletePostContentMsg {
	return DeletePostContentMsg{
		Creator:  types.AccountKey(creator),
		Permlink: permlink,
		Reason:   reason,
	}
}

func (msg DeletePostContentMsg) GetPermlink() types.Permlink  { return msg.Permlink }
func (msg DeletePostContentMsg) GetCreator() types.AccountKey { return msg.Creator }
func (msg DeletePostContentMsg) GetReason() string            { return msg.Reason }
func (msg DeletePostContentMsg) Type() string                 { return types.ProposalRouterName }

func (msg DeletePostContentMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	if len(msg.GetPermlink()) == 0 {
		return ErrInvalidPermlink()
	}
	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg DeletePostContentMsg) String() string {
	return fmt.Sprintf("DeletePostContentMsg{Creator:%v, post:%v}", msg.Creator, msg.GetPermlink())
}

func (msg DeletePostContentMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}
func (msg DeletePostContentMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg DeletePostContentMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg DeletePostContentMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// UpgradeProtocolMsg Msg Implementations

func NewUpgradeProtocolMsg(
	creator, link, reason string) UpgradeProtocolMsg {
	return UpgradeProtocolMsg{
		Creator: types.AccountKey(creator),
		Link:    link,
		Reason:  reason,
	}
}

func (msg UpgradeProtocolMsg) GetCreator() types.AccountKey { return msg.Creator }
func (msg UpgradeProtocolMsg) GetLink() string              { return msg.Link }
func (msg UpgradeProtocolMsg) GetReason() string            { return msg.Reason }
func (msg UpgradeProtocolMsg) Type() string                 { return types.ProposalRouterName }

func (msg UpgradeProtocolMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	if len(msg.GetLink()) == 0 {
		return ErrInvalidLink()
	}
	if len(msg.GetLink()) > types.MaximumLinkURL {
		return ErrInvalidLink()
	}
	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg UpgradeProtocolMsg) String() string {
	return fmt.Sprintf("UpgradeProtocolMsg{Creator:%v, Link:%v}", msg.Creator, msg.GetLink())
}

func (msg UpgradeProtocolMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg UpgradeProtocolMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg UpgradeProtocolMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg UpgradeProtocolMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangeGlobalAllocationParamMsg Msg Implementations

func NewChangeGlobalAllocationParamMsg(
	creator string, parameter param.GlobalAllocationParam, reason string) ChangeGlobalAllocationParamMsg {
	return ChangeGlobalAllocationParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangeGlobalAllocationParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeGlobalAllocationParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeGlobalAllocationParamMsg) GetReason() string             { return msg.Reason }
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
	if msg.Parameter.InfraAllocation.LT(sdk.ZeroRat()) ||
		msg.Parameter.ContentCreatorAllocation.LT(sdk.ZeroRat()) ||
		msg.Parameter.DeveloperAllocation.LT(sdk.ZeroRat()) ||
		msg.Parameter.ValidatorAllocation.LT(sdk.ZeroRat()) {
		return ErrIllegalParameter()
	}

	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg ChangeGlobalAllocationParamMsg) String() string {
	return fmt.Sprintf("ChangeGlobalAllocationParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeGlobalAllocationParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangeGlobalAllocationParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeGlobalAllocationParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangeGlobalAllocationParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangeEvaluateOfContentValueParamMsg Msg Implementations

func NewChangeEvaluateOfContentValueParamMsg(
	creator string, parameter param.EvaluateOfContentValueParam, reason string) ChangeEvaluateOfContentValueParamMsg {
	return ChangeEvaluateOfContentValueParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangeEvaluateOfContentValueParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeEvaluateOfContentValueParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeEvaluateOfContentValueParamMsg) GetReason() string             { return msg.Reason }
func (msg ChangeEvaluateOfContentValueParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeEvaluateOfContentValueParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	if msg.Parameter.ConsumptionTimeAdjustBase <= 0 ||
		msg.Parameter.TotalAmountOfConsumptionBase <= 0 {
		return ErrIllegalParameter()
	}

	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg ChangeEvaluateOfContentValueParamMsg) String() string {
	return fmt.Sprintf("ChangeEvaluateOfContentValueParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeEvaluateOfContentValueParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangeEvaluateOfContentValueParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeEvaluateOfContentValueParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangeEvaluateOfContentValueParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangeInfraInternalAllocationParamMsg Msg Implementations

func NewChangeInfraInternalAllocationParamMsg(
	creator string, parameter param.InfraInternalAllocationParam, reason string) ChangeInfraInternalAllocationParamMsg {
	return ChangeInfraInternalAllocationParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangeInfraInternalAllocationParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeInfraInternalAllocationParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeInfraInternalAllocationParamMsg) GetReason() string             { return msg.Reason }
func (msg ChangeInfraInternalAllocationParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeInfraInternalAllocationParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if !msg.Parameter.CDNAllocation.
		Add(msg.Parameter.StorageAllocation).Equal(sdk.NewRat(1)) {
		return ErrIllegalParameter()
	}

	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg ChangeInfraInternalAllocationParamMsg) String() string {
	return fmt.Sprintf("ChangeInfraInternalAllocationParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeInfraInternalAllocationParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangeInfraInternalAllocationParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeInfraInternalAllocationParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangeInfraInternalAllocationParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangeVoteParamMsg Msg Implementations

func NewChangeVoteParamMsg(
	creator string, parameter param.VoteParam, reason string) ChangeVoteParamMsg {
	return ChangeVoteParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangeVoteParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeVoteParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeVoteParamMsg) GetReason() string             { return msg.Reason }
func (msg ChangeVoteParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeVoteParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if msg.Parameter.DelegatorCoinReturnIntervalSec <= 0 ||
		msg.Parameter.VoterCoinReturnIntervalSec <= 0 ||
		msg.Parameter.DelegatorCoinReturnTimes <= 0 ||
		msg.Parameter.VoterCoinReturnTimes <= 0 {
		return ErrIllegalParameter()
	}

	if !msg.Parameter.DelegatorMinWithdraw.IsPositive() ||
		!msg.Parameter.VoterMinDeposit.IsPositive() ||
		!msg.Parameter.VoterMinWithdraw.IsPositive() {
		return ErrIllegalParameter()
	}
	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg ChangeVoteParamMsg) String() string {
	return fmt.Sprintf("ChangeVoteParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeVoteParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangeVoteParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeVoteParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangeVoteParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangeProposalParamMsg Msg Implementations

func NewChangeProposalParamMsg(
	creator string, parameter param.ProposalParam, reason string) ChangeProposalParamMsg {
	return ChangeProposalParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangeProposalParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeProposalParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeProposalParamMsg) GetReason() string             { return msg.Reason }
func (msg ChangeProposalParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeProposalParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if msg.Parameter.ContentCensorshipDecideSec <= 0 ||
		msg.Parameter.ChangeParamExecutionSec <= 0 ||
		msg.Parameter.ChangeParamDecideSec <= 0 ||
		msg.Parameter.ProtocolUpgradeDecideSec <= 0 {
		return ErrIllegalParameter()
	}

	if !msg.Parameter.ContentCensorshipMinDeposit.IsPositive() ||
		!msg.Parameter.ContentCensorshipPassVotes.IsPositive() ||
		!msg.Parameter.ChangeParamMinDeposit.IsPositive() ||
		!msg.Parameter.ChangeParamPassVotes.IsPositive() ||
		!msg.Parameter.ProtocolUpgradePassVotes.IsPositive() ||
		!msg.Parameter.ProtocolUpgradeMinDeposit.IsPositive() {
		return ErrIllegalParameter()
	}

	if !msg.Parameter.ContentCensorshipPassRatio.GT(sdk.ZeroRat()) ||
		!msg.Parameter.ChangeParamPassRatio.GT(sdk.ZeroRat()) ||
		!msg.Parameter.ProtocolUpgradePassRatio.GT(sdk.ZeroRat()) ||
		msg.Parameter.ProtocolUpgradePassRatio.GT(sdk.NewRat(1, 1)) ||
		msg.Parameter.ChangeParamPassRatio.GT(sdk.NewRat(1, 1)) ||
		msg.Parameter.ContentCensorshipPassRatio.GT(sdk.NewRat(1, 1)) {
		return ErrIllegalParameter()
	}

	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg ChangeProposalParamMsg) String() string {
	return fmt.Sprintf("ChangeProposalParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeProposalParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangeProposalParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeProposalParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangeProposalParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangeDeveloperParamMsg Msg Implementations

func NewChangeDeveloperParamMsg(
	creator string, parameter param.DeveloperParam, reason string) ChangeDeveloperParamMsg {
	return ChangeDeveloperParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangeDeveloperParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeDeveloperParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeDeveloperParamMsg) GetReason() string             { return msg.Reason }
func (msg ChangeDeveloperParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeDeveloperParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if msg.Parameter.DeveloperCoinReturnIntervalSec <= 0 ||
		msg.Parameter.DeveloperCoinReturnTimes <= 0 {
		return ErrIllegalParameter()
	}

	if !msg.Parameter.DeveloperMinDeposit.IsPositive() {
		return ErrIllegalParameter()
	}

	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg ChangeDeveloperParamMsg) String() string {
	return fmt.Sprintf("ChangeDeveloperParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeDeveloperParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangeDeveloperParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeDeveloperParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangeDeveloperParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangeValidatorParamMsg Msg Implementations

func NewChangeValidatorParamMsg(creator string, parameter param.ValidatorParam, reason string) ChangeValidatorParamMsg {
	return ChangeValidatorParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangeValidatorParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeValidatorParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeValidatorParamMsg) GetReason() string             { return msg.Reason }
func (msg ChangeValidatorParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeValidatorParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if msg.Parameter.ValidatorCoinReturnIntervalSec <= 0 ||
		msg.Parameter.ValidatorCoinReturnTimes <= 0 ||
		msg.Parameter.AbsentCommitLimitation <= 0 ||
		msg.Parameter.ValidatorListSize <= 0 {
		return ErrIllegalParameter()
	}

	if !msg.Parameter.ValidatorMinWithdraw.IsPositive() ||
		!msg.Parameter.ValidatorMinVotingDeposit.IsPositive() ||
		!msg.Parameter.ValidatorMinCommittingDeposit.IsPositive() ||
		!msg.Parameter.PenaltyMissVote.IsPositive() ||
		!msg.Parameter.PenaltyMissCommit.IsPositive() ||
		!msg.Parameter.PenaltyByzantine.IsPositive() {
		return ErrIllegalParameter()
	}

	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg ChangeValidatorParamMsg) String() string {
	return fmt.Sprintf("ChangeValidatorParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeValidatorParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangeValidatorParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeValidatorParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangeValidatorParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangeAccountParamMsg Msg Implementations

func NewChangeAccountParamMsg(
	creator string, parameter param.AccountParam, reason string) ChangeAccountParamMsg {
	return ChangeAccountParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangeAccountParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeAccountParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeAccountParamMsg) GetReason() string             { return msg.Reason }
func (msg ChangeAccountParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeAccountParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if types.NewCoinFromInt64(0).IsGT(msg.Parameter.MinimumBalance) ||
		types.NewCoinFromInt64(0).IsGT(msg.Parameter.RegisterFee) {
		return ErrIllegalParameter()
	}
	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg ChangeAccountParamMsg) String() string {
	return fmt.Sprintf("ChangeAccountParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeAccountParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangeAccountParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeAccountParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangeAccountParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangePostParam Msg Implementations

func NewChangePostParamMsg(
	creator string, parameter param.PostParam, reason string) ChangePostParamMsg {
	return ChangePostParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangePostParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangePostParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangePostParamMsg) GetReason() string             { return msg.Reason }
func (msg ChangePostParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangePostParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	if msg.Parameter.PostIntervalSec < 0 || msg.Parameter.ReportOrUpvoteIntervalSec < 0 {
		return ErrIllegalParameter()
	}
	return nil
}

func (msg ChangePostParamMsg) String() string {
	return fmt.Sprintf("ChangePostParamMsg{Creator:%v, param:%v}", msg.Creator, msg.Parameter)
}

func (msg ChangePostParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangePostParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangePostParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangePostParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// ChangeBandwidthParamMsg Msg Implementations

func NewChangeBandwidthParamMsg(
	creator string, parameter param.BandwidthParam, reason string) ChangeBandwidthParamMsg {
	return ChangeBandwidthParamMsg{
		Creator:   types.AccountKey(creator),
		Parameter: parameter,
		Reason:    reason,
	}
}

func (msg ChangeBandwidthParamMsg) GetParameter() param.Parameter { return msg.Parameter }
func (msg ChangeBandwidthParamMsg) GetCreator() types.AccountKey  { return msg.Creator }
func (msg ChangeBandwidthParamMsg) GetReason() string             { return msg.Reason }
func (msg ChangeBandwidthParamMsg) Type() string                  { return types.ProposalRouterName }

func (msg ChangeBandwidthParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if types.NewCoinFromInt64(0).IsGT(msg.Parameter.CapacityUsagePerTransaction) {
		return ErrIllegalParameter()
	}

	if msg.Parameter.SecondsToRecoverBandwidth <= 0 {
		return ErrIllegalParameter()
	}
	if utf8.RuneCountInString(msg.Reason) > types.MaximumLengthOfProposalReason {
		return ErrReasonTooLong()
	}
	return nil
}

func (msg ChangeBandwidthParamMsg) String() string {
	return fmt.Sprintf("ChangeBandwidthParamMsg{Creator:%v}", msg.Creator)
}

func (msg ChangeBandwidthParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg ChangeBandwidthParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChangeBandwidthParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// Implements Msg.
func (msg ChangeBandwidthParamMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

//----------------------------------------
// VoteProposalMsg Msg Implementations

func NewVoteProposalMsg(voter string, proposalID int64, result bool) VoteProposalMsg {
	return VoteProposalMsg{
		Voter:      types.AccountKey(voter),
		ProposalID: types.ProposalKey(strconv.FormatInt(proposalID, 10)),
		Result:     result,
	}
}

func (msg VoteProposalMsg) Type() string { return types.ProposalRouterName }

func (msg VoteProposalMsg) ValidateBasic() sdk.Error {
	if len(msg.Voter) < types.MinimumUsernameLength ||
		len(msg.Voter) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}
	return nil
}

func (msg VoteProposalMsg) String() string {
	return fmt.Sprintf("VoteProposalMsg{Voter:%v, ProposalID:%v, Result:%v}", msg.Voter, msg.ProposalID, msg.Result)
}

func (msg VoteProposalMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

func (msg VoteProposalMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

func (msg VoteProposalMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Voter)}
}

// Implements Msg.
func (msg VoteProposalMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}
