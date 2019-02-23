package proposal

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(VoteProposalMsg{}, "lino/voteProposal", nil)
	cdc.RegisterConcrete(DeletePostContentMsg{}, "lino/deletePostContent", nil)
	cdc.RegisterConcrete(UpgradeProtocolMsg{}, "lino/upgradeProtocol", nil)
	cdc.RegisterConcrete(ChangeGlobalAllocationParamMsg{}, "lino/changeGlobalAllocation", nil)
	cdc.RegisterConcrete(ChangeInfraInternalAllocationParamMsg{}, "lino/changeInfraAllocation", nil)
	cdc.RegisterConcrete(ChangeVoteParamMsg{}, "lino/changeVoteParam", nil)
	cdc.RegisterConcrete(ChangeProposalParamMsg{}, "lino/changeProposalParam", nil)
	cdc.RegisterConcrete(ChangeDeveloperParamMsg{}, "lino/changeDeveloperParam", nil)
	cdc.RegisterConcrete(ChangeValidatorParamMsg{}, "lino/changeValidatorParam", nil)
	cdc.RegisterConcrete(ChangeBandwidthParamMsg{}, "lino/changeBandwidthParam", nil)
	cdc.RegisterConcrete(ChangeAccountParamMsg{}, "lino/changeAccountParam", nil)
	cdc.RegisterConcrete(ChangePostParamMsg{}, "lino/changePostParam", nil)
}

var msgCdc = wire.New()

func init() {
	RegisterWire(msgCdc)
}
