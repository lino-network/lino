package proposal

import (
	"math/big"

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

func (pm ProposalManager) IsOngoingProposal(ctx sdk.Context, proposalID types.ProposalKey) bool {
	lst, err := pm.storage.GetProposalList(ctx)
	if err != nil {
		return false
	}

	for _, id := range lst.OngoingProposal {
		if id == proposalID {
			return true
		}
	}
	return false
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
		AgreeVotes:    types.NewCoinFromInt64(0),
		DisagreeVotes: types.NewCoinFromInt64(0),
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

func (pm ProposalManager) GetProposalPassParam(
	ctx sdk.Context, proposalType types.ProposalType) (sdk.Rat, types.Coin, sdk.Error) {
	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return sdk.NewRat(1, 1), types.NewCoinFromInt64(0), err
	}
	switch proposalType {
	case types.ChangeParam:
		return param.ChangeParamPassRatio, param.ChangeParamPassVotes, nil
	case types.ContentCensorship:
		return param.ContentCensorshipPassRatio, param.ContentCensorshipPassVotes, nil
	case types.ProtocolUpgrade:
		return param.ProtocolUpgradePassRatio, param.ProtocolUpgradePassVotes, nil
	default:
		return sdk.NewRat(1, 1), types.NewCoinFromInt64(0), ErrWrongProposalType()
	}
}

func (pm ProposalManager) UpdateProposalStatus(
	ctx sdk.Context, res types.VotingResult, proposalType types.ProposalType,
	proposalID types.ProposalKey) (types.ProposalResult, sdk.Error) {
	lst, err := pm.storage.GetProposalList(ctx)
	if err != nil {
		return types.ProposalNotPass, err
	}

	proposal, err := pm.storage.GetProposal(ctx, proposalID)
	if err != nil {
		return types.ProposalNotPass, err
	}

	proposalInfo := proposal.GetProposalInfo()

	proposalInfo.AgreeVotes = res.AgreeVotes
	proposalInfo.DisagreeVotes = res.DisagreeVotes

	// calculate if agree votes meet minimum pass requirement
	ratio, minVotes, err := pm.GetProposalPassParam(ctx, proposalType)
	if err != nil {
		return types.ProposalNotPass, err
	}
	totalVotes := res.AgreeVotes.Plus(res.DisagreeVotes)
	if !totalVotes.IsGT(minVotes) {
		return types.ProposalNotPass, nil
	}
	actualRatio := new(big.Rat).Quo(res.AgreeVotes.ToRat(), totalVotes.ToRat())
	if actualRatio.Cmp(ratio.GetRat()) >= 0 {
		proposalInfo.Result = types.ProposalPass
	} else {
		proposalInfo.Result = types.ProposalNotPass
	}

	proposal.SetProposalInfo(proposalInfo)
	if err := pm.storage.SetProposal(ctx, proposalID, proposal); err != nil {
		return types.ProposalNotPass, err
	}

	// update ongoing and past proposal list
	for index, id := range lst.OngoingProposal {
		if id == proposalID {
			lst.OngoingProposal = append(lst.OngoingProposal[:index], lst.OngoingProposal[index+1:]...)
			break
		}
	}
	lst.PastProposal = append(lst.PastProposal, proposalID)

	if err := pm.storage.SetProposalList(ctx, lst); err != nil {
		return types.ProposalNotPass, err
	}
	return proposalInfo.Result, nil
}

func (pm ProposalManager) CreateDecideProposalEvent(
	ctx sdk.Context, proposalType types.ProposalType, proposalID types.ProposalKey) (types.Event, sdk.Error) {
	event := DecideProposalEvent{
		ProposalType: proposalType,
		ProposalID:   proposalID,
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
