package proposal

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	proposalmn "github.com/lino-network/lino/x/proposal/manager"
	"github.com/lino-network/lino/x/proposal/model"
)

type ProposalKeeper interface {
	InitGenesis(ctx sdk.Context) error
	ChangeParam(ctx sdk.Context, creator linotypes.AccountKey, reason string, p param.Parameter) sdk.Error
	ProtocolUpgrade(ctx sdk.Context, creator linotypes.AccountKey, reason, link string) sdk.Error
	ContentCensorship(ctx sdk.Context, creator linotypes.AccountKey, reason string, permlink linotypes.Permlink) sdk.Error
	GetOngoingProposal(ctx sdk.Context, proposalID linotypes.ProposalKey) (model.Proposal, sdk.Error)
	GetExpiredProposal(ctx sdk.Context, proposalID linotypes.ProposalKey) (model.Proposal, sdk.Error)
}

var _ ProposalKeeper = proposalmn.ProposalManager{}
