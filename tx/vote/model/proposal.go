package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/lino-network/lino/types"
)

type Proposal struct {
	Creator      types.AccountKey  `json:"creator"`
	ProposalID   types.ProposalKey `json:"proposal_id"`
	AgreeVote    types.Coin        `json:"agree_vote"`
	DisagreeVote types.Coin        `json:"disagree_vote"`
}

type ChangeParameterDescription struct {
	InfraAllocation          sdk.Rat `json:"infra_allocation"`
	ContentCreatorAllocation sdk.Rat `json:"content_creator_allocation"`
	DeveloperAllocation      sdk.Rat `json:"developer_allocation"`
	ValidatorAllocation      sdk.Rat `json:"validator_allocation"`
	StorageAllocation        sdk.Rat `json:"storage_allocation"`
	CDNAllocation            sdk.Rat `json:"CDN_allocation"`
}

type ChangeParameterProposal struct {
	Proposal
	ChangeParameterDescription
}

type ProposalList struct {
	OngoingProposal []types.ProposalKey `json:"ongoing_proposal"`
	PastProposal    []types.ProposalKey `json:"past_proposal"`
}
