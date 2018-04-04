package vote

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/lino-network/lino/types"
)

type ProposalKey string

type Proposal struct {
	ProposalID   ProposalKey `json:"proposal_id"`
	AgreeVote    types.Coin  `json:"agree_vote"`
	DisagreeVote types.Coin  `json:"disagree_vote"`
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

var nextProposalID = big.NewInt(0)
var ProposalDecideHr = int64(7 * 24)
