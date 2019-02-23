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

// ChangeParamMsg - change parameter msg
type ChangeParamMsg interface {
	GetParameter() param.Parameter
	GetCreator() types.AccountKey
	GetReason() string
}

// ContentCensorshipMsg - content censorship msg
type ContentCensorshipMsg interface {
	GetCreator() types.AccountKey
	GetPermlink() types.Permlink
	GetReason() string
}

// ProtocolUpgradeMsg - protocol upgrade msg
type ProtocolUpgradeMsg interface {
	GetCreator() types.AccountKey
	GetLink() string
	GetReason() string
}

// DeletePostContentMsg - implement of content censorship msg
type DeletePostContentMsg struct {
	Creator  types.AccountKey `json:"creator"`
	Permlink types.Permlink   `json:"permlink"`
	Reason   string           `json:"reason"`
}

// UpgradeProtocolMsg - implement of protocol upgrade msg
type UpgradeProtocolMsg struct {
	Creator types.AccountKey `json:"creator"`
	Link    string           `json:"link"`
	Reason  string           `json:"reason"`
}

// ChangeGlobalAllocationParamMsg - implement of change parameter msg
type ChangeGlobalAllocationParamMsg struct {
	Creator   types.AccountKey            `json:"creator"`
	Parameter param.GlobalAllocationParam `json:"parameter"`
	Reason    string                      `json:"reason"`
}

// ChangeInfraInternalAllocationParamMsg - implement of change parameter msg
type ChangeInfraInternalAllocationParamMsg struct {
	Creator   types.AccountKey                   `json:"creator"`
	Parameter param.InfraInternalAllocationParam `json:"parameter"`
	Reason    string                             `json:"reason"`
}

// ChangeVoteParamMsg - implement of change parameter msg
type ChangeVoteParamMsg struct {
	Creator   types.AccountKey `json:"creator"`
	Parameter param.VoteParam  `json:"parameter"`
	Reason    string           `json:"reason"`
}

// ChangeProposalParamMsg - implement of change parameter msg
type ChangeProposalParamMsg struct {
	Creator   types.AccountKey    `json:"creator"`
	Parameter param.ProposalParam `json:"parameter"`
	Reason    string              `json:"reason"`
}

// ChangeDeveloperParamMsg - implement of change parameter msg
type ChangeDeveloperParamMsg struct {
	Creator   types.AccountKey     `json:"creator"`
	Parameter param.DeveloperParam `json:"parameter"`
	Reason    string               `json:"reason"`
}

// ChangeValidatorParamMsg - implement of change parameter msg
type ChangeValidatorParamMsg struct {
	Creator   types.AccountKey     `json:"creator"`
	Parameter param.ValidatorParam `json:"parameter"`
	Reason    string               `json:"reason"`
}

// ChangeBandwidthParamMsg - implement of change parameter msg
type ChangeBandwidthParamMsg struct {
	Creator   types.AccountKey     `json:"creator"`
	Parameter param.BandwidthParam `json:"parameter"`
	Reason    string               `json:"reason"`
}

// ChangeAccountParamMsg - implement of change parameter msg
type ChangeAccountParamMsg struct {
	Creator   types.AccountKey   `json:"creator"`
	Parameter param.AccountParam `json:"parameter"`
	Reason    string             `json:"reason"`
}

// ChangePostParamMsg - implement of change parameter msg
type ChangePostParamMsg struct {
	Creator   types.AccountKey `json:"creator"`
	Parameter param.PostParam  `json:"parameter"`
	Reason    string           `json:"reason"`
}

// VoteProposalMsg - implement of change parameter msg
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

// GetPermlink - implement DeletePostContentMsg
func (msg DeletePostContentMsg) GetPermlink() types.Permlink { return msg.Permlink }

// GetCreator - implement DeletePostContentMsg
func (msg DeletePostContentMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement DeletePostContentMsg
func (msg DeletePostContentMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg DeletePostContentMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg DeletePostContentMsg) Type() string { return "DeletePostContentMsg" }

// ValidateBasic - implement sdk.Msg
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

// GetPermission - implement types.Msg
func (msg DeletePostContentMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg DeletePostContentMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg DeletePostContentMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// GetCreator - implement UpgradeProtocolMsg
func (msg UpgradeProtocolMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetLink - implement UpgradeProtocolMsg
func (msg UpgradeProtocolMsg) GetLink() string { return msg.Link }

// GetReason - implement UpgradeProtocolMsg
func (msg UpgradeProtocolMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg UpgradeProtocolMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg UpgradeProtocolMsg) Type() string { return "UpgradeProtocolMsg" }

// ValidateBasic - implement sdk.Msg
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

// GetPermission - implement types.Msg
func (msg UpgradeProtocolMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg UpgradeProtocolMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg UpgradeProtocolMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// GetParameter - implement ChangeParamMsg
func (msg ChangeGlobalAllocationParamMsg) GetParameter() param.Parameter { return msg.Parameter }

// GetCreator - implement ChangeParamMsg
func (msg ChangeGlobalAllocationParamMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement ChangeParamMsg
func (msg ChangeGlobalAllocationParamMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg ChangeGlobalAllocationParamMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg ChangeGlobalAllocationParamMsg) Type() string { return "ChangeGlobalAllocationParamMsg" }

// ValidateBasic - implement sdk.Msg
func (msg ChangeGlobalAllocationParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if !msg.Parameter.InfraAllocation.
		Add(msg.Parameter.ContentCreatorAllocation).
		Add(msg.Parameter.DeveloperAllocation).
		Add(msg.Parameter.ValidatorAllocation).Equal(sdk.NewDec(1)) {
		return ErrIllegalParameter()
	}
	if msg.Parameter.InfraAllocation.LT(sdk.ZeroDec()) ||
		msg.Parameter.ContentCreatorAllocation.LT(sdk.ZeroDec()) ||
		msg.Parameter.DeveloperAllocation.LT(sdk.ZeroDec()) ||
		msg.Parameter.ValidatorAllocation.LT(sdk.ZeroDec()) {
		return ErrIllegalParameter()
	}
	if msg.Parameter.GlobalGrowthRate.GT(param.AnnualInflationCeiling) {
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

// GetPermission - implement types.Msg
func (msg ChangeGlobalAllocationParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ChangeGlobalAllocationParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ChangeGlobalAllocationParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
func (msg ChangeGlobalAllocationParamMsg) GetConsumeAmount() types.Coin {
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

// GetParameter - implement ChangeParamMsg
func (msg ChangeInfraInternalAllocationParamMsg) GetParameter() param.Parameter { return msg.Parameter }

// GetCreator - implement ChangeParamMsg
func (msg ChangeInfraInternalAllocationParamMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement ChangeParamMsg
func (msg ChangeInfraInternalAllocationParamMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg ChangeInfraInternalAllocationParamMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg ChangeInfraInternalAllocationParamMsg) Type() string {
	return "ChangeInfraInternalAllocationParamMsg"
}

// ValidateBasic - implement sdk.Msg
func (msg ChangeInfraInternalAllocationParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if !msg.Parameter.CDNAllocation.
		Add(msg.Parameter.StorageAllocation).Equal(sdk.NewDec(1)) ||
		msg.Parameter.CDNAllocation.LT(sdk.ZeroDec()) ||
		msg.Parameter.StorageAllocation.LT(sdk.ZeroDec()) {
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

// GetPermission - implement types.Msg
func (msg ChangeInfraInternalAllocationParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ChangeInfraInternalAllocationParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ChangeInfraInternalAllocationParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// GetParameter - implement ChangeParamMsg
func (msg ChangeVoteParamMsg) GetParameter() param.Parameter { return msg.Parameter }

// GetCreator - implement ChangeParamMsg
func (msg ChangeVoteParamMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement ChangeParamMsg
func (msg ChangeVoteParamMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg ChangeVoteParamMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg ChangeVoteParamMsg) Type() string { return "ChangeVoteParamMsg" }

// ValidateBasic - implement sdk.Msg
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

	if !msg.Parameter.MinStakeIn.IsPositive() {
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

// GetPermission - implement types.Msg
func (msg ChangeVoteParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ChangeVoteParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ChangeVoteParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// GetParameter - implement ChangeParamMsg
func (msg ChangeProposalParamMsg) GetParameter() param.Parameter { return msg.Parameter }

// GetCreator - implement ChangeParamMsg
func (msg ChangeProposalParamMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement ChangeParamMsg
func (msg ChangeProposalParamMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg ChangeProposalParamMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg ChangeProposalParamMsg) Type() string { return "ChangeProposalParamMsg" }

// ValidateBasic - implement sdk.Msg
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

	if !msg.Parameter.ContentCensorshipPassRatio.GT(sdk.ZeroDec()) ||
		!msg.Parameter.ChangeParamPassRatio.GT(sdk.ZeroDec()) ||
		!msg.Parameter.ProtocolUpgradePassRatio.GT(sdk.ZeroDec()) ||
		msg.Parameter.ProtocolUpgradePassRatio.GT(sdk.NewDec(1)) ||
		msg.Parameter.ChangeParamPassRatio.GT(sdk.NewDec(1)) ||
		msg.Parameter.ContentCensorshipPassRatio.GT(sdk.NewDec(1)) {
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

// GetPermission - implement types.Msg
func (msg ChangeProposalParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ChangeProposalParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ChangeProposalParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// GetParameter - implement ChangeParamMsg
func (msg ChangeDeveloperParamMsg) GetParameter() param.Parameter { return msg.Parameter }

// GetCreator - implement ChangeParamMsg
func (msg ChangeDeveloperParamMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement ChangeParamMsg
func (msg ChangeDeveloperParamMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg ChangeDeveloperParamMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg ChangeDeveloperParamMsg) Type() string { return "ChangeDeveloperParamMsg" }

// ValidateBasic - implement sdk.Msg
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

// GetPermission - implement types.Msg
func (msg ChangeDeveloperParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ChangeDeveloperParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ChangeDeveloperParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// GetParameter - implement ChangeParamMsg
func (msg ChangeValidatorParamMsg) GetParameter() param.Parameter { return msg.Parameter }

// GetCreator - implement ChangeParamMsg
func (msg ChangeValidatorParamMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement ChangeParamMsg
func (msg ChangeValidatorParamMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg ChangeValidatorParamMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg ChangeValidatorParamMsg) Type() string { return "ChangeValidatorParamMsg" }

// ValidateBasic - implement sdk.Msg
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

// GetPermission - implement types.Msg
func (msg ChangeValidatorParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ChangeValidatorParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ChangeValidatorParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// GetParameter - implement ChangeParamMsg
func (msg ChangeAccountParamMsg) GetParameter() param.Parameter { return msg.Parameter }

// GetCreator - implement ChangeParamMsg
func (msg ChangeAccountParamMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement ChangeParamMsg
func (msg ChangeAccountParamMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg ChangeAccountParamMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg ChangeAccountParamMsg) Type() string { return "ChangeAccountParamMsg" }

// ValidateBasic - implement sdk.Msg
func (msg ChangeAccountParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if !msg.Parameter.MinimumBalance.IsNotNegative() ||
		!msg.Parameter.RegisterFee.IsNotNegative() ||
		!msg.Parameter.FirstDepositFullCoinDayLimit.IsNotNegative() ||
		msg.Parameter.MaxNumFrozenMoney <= 0 {
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

// GetPermission - implement types.Msg
func (msg ChangeAccountParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ChangeAccountParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ChangeAccountParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// GetParameter - implement ChangeParamMsg
func (msg ChangePostParamMsg) GetParameter() param.Parameter { return msg.Parameter }

// GetCreator - implement ChangeParamMsg
func (msg ChangePostParamMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement ChangeParamMsg
func (msg ChangePostParamMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg ChangePostParamMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg ChangePostParamMsg) Type() string { return "ChangePostParamMsg" }

// ValidateBasic - implement sdk.Msg
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

// GetPermission - implement types.Msg
func (msg ChangePostParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ChangePostParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ChangePostParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// GetParameter - implement ChangeParamMsg
func (msg ChangeBandwidthParamMsg) GetParameter() param.Parameter { return msg.Parameter }

// GetCreator - implement ChangeParamMsg
func (msg ChangeBandwidthParamMsg) GetCreator() types.AccountKey { return msg.Creator }

// GetReason - implement ChangeParamMsg
func (msg ChangeBandwidthParamMsg) GetReason() string { return msg.Reason }

// Route - implement sdk.Msg
func (msg ChangeBandwidthParamMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg ChangeBandwidthParamMsg) Type() string { return "ChangeBandwidthParamMsg" }

// ValidateBasic - implement sdk.Msg
func (msg ChangeBandwidthParamMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) < types.MinimumUsernameLength ||
		len(msg.Creator) > types.MaximumUsernameLength {
		return ErrInvalidUsername()
	}

	if !msg.Parameter.CapacityUsagePerTransaction.IsNotNegative() {
		return ErrIllegalParameter()
	}
	if !msg.Parameter.VirtualCoin.IsNotNegative() {
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

// GetPermission - implement types.Msg
func (msg ChangeBandwidthParamMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg ChangeBandwidthParamMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg ChangeBandwidthParamMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetConsumeAmount - implement types.Msg
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

// Route - implement sdk.Msg
func (msg VoteProposalMsg) Route() string { return types.ProposalRouterName }

// Type - implement sdk.Msg
func (msg VoteProposalMsg) Type() string { return "VoteProposalMsg" }

// ValidateBasic - implement sdk.Msg
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

// GetPermission - implement types.Msg
func (msg VoteProposalMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implement sdk.Msg
func (msg VoteProposalMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implement sdk.Msg
func (msg VoteProposalMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Voter)}
}

// GetConsumeAmount - implement types.Msg
func (msg VoteProposalMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}
