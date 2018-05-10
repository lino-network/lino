package proposal

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/tx/proposal/model"
	"github.com/lino-network/lino/types"
)

type ProposalManager struct {
	storage     model.ProposalStorage `json:"proposal_storage"`
	paramHolder param.ParamHolder     `json:"param_holder"`
}

func NewProposalManager(key sdk.StoreKey, holder param.ParamHolder) ProposalManager {
	return ProposalManager{
		storage:     model.NewProposalStorage(key),
		paramHolder: holder,
	}
}

func (pm ProposalManager) InitGenesis(ctx sdk.Context) error {
	if err := pm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (pm ProposalManager) IsProposalExist(ctx sdk.Context, proposalID types.ProposalKey) bool {
	proposalByte, _ := pm.storage.GetProposal(ctx, proposalID)
	return proposalByte != nil
}

func (pm ProposalManager) CreateContentCensorshipProposal(
	ctx sdk.Context, permLink types.PermLink) model.Proposal {
	return &model.ContentCensorshipProposal{
		PermLink: permLink,
	}
}

func (pm ProposalManager) CreateProtocolUpgradeProposal(ctx sdk.Context, link string) model.Proposal {
	return &model.ProtocolUpgradeProposal{
		Link: link,
	}
}

func (pm ProposalManager) CreateChangeParamProposal(
	ctx sdk.Context, parameter param.Parameter) model.Proposal {
	return &model.ChangeParamProposal{
		Param: parameter,
	}
}

func (pm ProposalManager) AddProposal(
	ctx sdk.Context, creator types.AccountKey, proposal model.Proposal) (types.ProposalKey, sdk.Error) {
	newID, err := pm.paramHolder.GetNextProposalID(ctx)
	if err != nil {
		return newID, err
	}

	info := model.ProposalInfo{
		Creator:       creator,
		ProposalID:    newID,
		AgreeVotes:    types.Coin{Amount: 0},
		DisagreeVotes: types.Coin{Amount: 0},
		Result:        types.ProposalNotPass,
	}
	proposal.SetProposalInfo(info)

	if err := pm.storage.SetProposal(ctx, newID, proposal); err != nil {
		return newID, err
	}

	lst, err := pm.storage.GetProposalList(ctx)
	if err != nil {
		return newID, err
	}
	lst.OngoingProposal = append(lst.OngoingProposal, newID)
	if err := pm.storage.SetProposalList(ctx, lst); err != nil {
		return newID, err
	}

	return newID, nil
}

func (pm ProposalManager) GetCurrentProposal(ctx sdk.Context) (types.ProposalKey, sdk.Error) {
	lst, err := pm.storage.GetProposalList(ctx)
	if err != nil {
		return types.ProposalKey(""), err
	}

	if len(lst.OngoingProposal) == 0 {
		return types.ProposalKey(""), ErrOngoingProposalNotFound()
	}
	return lst.OngoingProposal[0], nil
}

func (pm ProposalManager) UpdateProposalStatus(
	ctx sdk.Context, res types.VotingResult, proposalType types.ProposalType) (types.ProposalResult, sdk.Error) {
	lst, err := pm.storage.GetProposalList(ctx)
	if err != nil {
		return types.ProposalNotPass, err
	}

	curID := lst.OngoingProposal[0]
	proposal, err := pm.storage.GetProposal(ctx, curID)
	if err != nil {
		return types.ProposalNotPass, err
	}

	proposalInfo := proposal.GetProposalInfo()

	proposalInfo.AgreeVotes = res.AgreeVotes
	proposalInfo.DisagreeVotes = res.DisagreeVotes

	// TODO consider different types of propsal
	if proposalInfo.AgreeVotes.IsGT(proposalInfo.DisagreeVotes) {
		proposalInfo.Result = types.ProposalPass
	}

	proposal.SetProposalInfo(proposalInfo)
	if err := pm.storage.SetProposal(ctx, curID, proposal); err != nil {
		return types.ProposalNotPass, err
	}

	lst.OngoingProposal = lst.OngoingProposal[1:]
	lst.PastProposal = append(lst.PastProposal, curID)

	if err := pm.storage.SetProposalList(ctx, lst); err != nil {
		return types.ProposalNotPass, err
	}
	return proposalInfo.Result, nil
}

func (pm ProposalManager) CreateDecideProposalEvent(
	ctx sdk.Context, proposalType types.ProposalType) (types.Event, sdk.Error) {
	event := DecideProposalEvent{
		ProposalType: proposalType,
	}
	return event, nil
}

func (pm ProposalManager) CreateParamChangeEvent(
	ctx sdk.Context, proposalID types.ProposalKey) (types.Event, sdk.Error) {
	proposal, err := pm.storage.GetProposal(ctx, proposalID)
	if err != nil {
		return nil, err
	}

	p, ok := proposal.(*model.ChangeParamProposal)
	if !ok {
		return nil, ErrWrongProposalType()
	}

	event := param.ChangeParamEvent{
		Param: p.Param,
	}
	return event, nil
}

func (pm ProposalManager) GetPermLink(ctx sdk.Context, proposalID types.ProposalKey) (types.PermLink, sdk.Error) {
	proposal, err := pm.storage.GetProposal(ctx, proposalID)
	if err != nil {
		return types.PermLink(""), err
	}

	p, ok := proposal.(*model.ContentCensorshipProposal)
	if !ok {
		return types.PermLink(""), ErrWrongProposalType()
	}
	return p.PermLink, nil
}

func (pm ProposalManager) GetProposalList(ctx sdk.Context) (*model.ProposalList, sdk.Error) {
	return pm.storage.GetProposalList(ctx)
}
