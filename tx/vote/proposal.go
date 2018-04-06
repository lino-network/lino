package vote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	types "github.com/lino-network/lino/types"
)

type ProposalKey string

type Proposal struct {
	Creator      acc.AccountKey `json:"creator"`
	ProposalID   ProposalKey    `json:"proposal_id"`
	AgreeVote    types.Coin     `json:"agree_vote"`
	DisagreeVote types.Coin     `json:"disagree_vote"`
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
	OngoingProposal []ProposalKey `json:"ongoing_proposal"`
	PastProposal    []ProposalKey `json:"past_proposal"`
}

var nextProposalID = int64(0)
var ProposalDecideHr = int64(7 * 24)
var CoinReturnIntervalHr = int64(7 * 24)
var CoinReturnTimes = int64(7)
