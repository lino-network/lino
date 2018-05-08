package model

import (
	"github.com/lino-network/lino/param"
	types "github.com/lino-network/lino/types"
)

type Proposal interface {
	GetProposalInfo() ProposalInfo
	SetProposalInfo(ProposalInfo)
}

type Description interface{}

type ProposalInfo struct {
	Creator       types.AccountKey     `json:"creator"`
	ProposalID    types.ProposalKey    `json:"proposal_id"`
	AgreeVotes    types.Coin           `json:"agree_vote"`
	DisagreeVotes types.Coin           `json:"disagree_vote"`
	Result        types.ProposalResult `json:"result"`
}

type ChangeGlobalAllocationParamProposal struct {
	ProposalInfo
	Description param.GlobalAllocationParam `json:"description"`
}

func (p *ChangeGlobalAllocationParamProposal) GetProposalInfo() ProposalInfo {
	return p.ProposalInfo
}

func (p *ChangeGlobalAllocationParamProposal) SetProposalInfo(info ProposalInfo) {
	p.ProposalInfo = info
}

type ProposalList struct {
	OngoingProposal []types.ProposalKey `json:"ongoing_proposal"`
	PastProposal    []types.ProposalKey `json:"past_proposal"`
}
