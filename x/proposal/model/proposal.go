package model

import (
	"github.com/lino-network/lino/param"
	types "github.com/lino-network/lino/types"
)

// Proposal - there are three proposal types
// 1) change parameter proposal
// 2) content censorship proposal
// 3) protocol upgrade proposal
type Proposal interface {
	GetProposalInfo() ProposalInfo
	SetProposalInfo(ProposalInfo)
}

// ProposalInfo - basic proposal info
type ProposalInfo struct {
	Creator       types.AccountKey     `json:"creator"`
	ProposalID    types.ProposalKey    `json:"proposal_id"`
	AgreeVotes    types.Coin           `json:"agree_vote"`
	DisagreeVotes types.Coin           `json:"disagree_vote"`
	Result        types.ProposalResult `json:"result"`
	CreatedAt     int64                `json:"created_at"`
	ExpiredAt     int64                `json:"expired_at"`
	Reason        string               `json:"reason"`
}

// ChangeParamProposal - change parameter proposal
type ChangeParamProposal struct {
	ProposalInfo
	Param  param.Parameter `json:"param"`
	Reason string          `json:"reason"` //nolint:govet
}

// GetProposalInfo - implements Proposal
func (p *ChangeParamProposal) GetProposalInfo() ProposalInfo { return p.ProposalInfo }

// SetProposalInfo - implements Proposal
func (p *ChangeParamProposal) SetProposalInfo(info ProposalInfo) { p.ProposalInfo = info }

// ContentCensorshipProposal - content censorship proposal
type ContentCensorshipProposal struct {
	ProposalInfo
	Permlink types.Permlink `json:"permlink"`
	Reason   string         `json:"reason"` //nolint:govet
}

// GetProposalInfo - implements Proposal
func (p *ContentCensorshipProposal) GetProposalInfo() ProposalInfo { return p.ProposalInfo }

// SetProposalInfo - implements Proposal
func (p *ContentCensorshipProposal) SetProposalInfo(info ProposalInfo) { p.ProposalInfo = info }

// ProtocolUpgradeProposal - protocol upgrade proposal
type ProtocolUpgradeProposal struct {
	ProposalInfo
	Link   string `json:"link"`
	Reason string `json:"reason"` //nolint:govet
}

// GetProposalInfo - implements Proposal
func (p *ProtocolUpgradeProposal) GetProposalInfo() ProposalInfo { return p.ProposalInfo }

// SetProposalInfo - implements Proposal
func (p *ProtocolUpgradeProposal) SetProposalInfo(info ProposalInfo) { p.ProposalInfo = info }

// NextProposalID - store next proposal ID to KVStore
type NextProposalID struct {
	NextProposalID int64 `json:"next_proposal_id"`
}
