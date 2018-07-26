package proposal

import (
	"strconv"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/proposal/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (pm ProposalManager) DoesProposalExist(ctx sdk.Context, proposalID types.ProposalKey) bool {
	return pm.storage.DoesProposalExist(ctx, proposalID)
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
	ctx sdk.Context, permlink types.Permlink, reason string) model.Proposal {
	return &model.ContentCensorshipProposal{
		Permlink: permlink,
		Reason:   reason,
	}
}

func (pm ProposalManager) CreateProtocolUpgradeProposal(ctx sdk.Context, link string, reason string) model.Proposal {
	return &model.ProtocolUpgradeProposal{
		Link:   link,
		Reason: reason,
	}
}

func (pm ProposalManager) CreateChangeParamProposal(
	ctx sdk.Context, parameter param.Parameter, reason string) model.Proposal {
	return &model.ChangeParamProposal{
		Param:  parameter,
		Reason: reason,
	}
}

func (pm ProposalManager) GetNextProposalID(ctx sdk.Context) (types.ProposalKey, sdk.Error) {
	nextProposalID, err := pm.storage.GetNextProposalID(ctx)
	if err != nil {
		return types.ProposalKey(""), err
	}

	return types.ProposalKey(strconv.FormatInt(nextProposalID.NextProposalID, 10)), nil
}

func (pm ProposalManager) IncreaseNextProposalID(ctx sdk.Context) sdk.Error {
	nextProposalID, err := pm.storage.GetNextProposalID(ctx)
	if err != nil {
		return err
	}

	nextProposalID.NextProposalID += 1
	if err := pm.storage.SetNextProposalID(ctx, nextProposalID); err != nil {
		return err
	}

	return nil
}

func (pm ProposalManager) AddProposal(
	ctx sdk.Context, creator types.AccountKey, proposal model.Proposal, decideHr int64) (types.ProposalKey, sdk.Error) {
	newID, err := pm.GetNextProposalID(ctx)
	if err != nil {
		return newID, err
	}

	info := model.ProposalInfo{
		Creator:       creator,
		ProposalID:    newID,
		AgreeVotes:    types.NewCoinFromInt64(0),
		DisagreeVotes: types.NewCoinFromInt64(0),
		Result:        types.ProposalNotPass,
		CreatedAt:     ctx.BlockHeader().Time,
		ExpiredAt:     ctx.BlockHeader().Time + decideHr*3600,
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

	if err := pm.IncreaseNextProposalID(ctx); err != nil {
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
		return sdk.NewRat(1, 1), types.NewCoinFromInt64(0), ErrIncorrectProposalType()
	}
}

func (pm ProposalManager) UpdateProposalVotingStatus(ctx sdk.Context, proposalID types.ProposalKey,
	voter types.AccountKey, voteResult bool, votingPower types.Coin) sdk.Error {
	proposal, err := pm.storage.GetProposal(ctx, proposalID)
	if err != nil {
		return err
	}
	proposalInfo := proposal.GetProposalInfo()

	if voteResult == true {
		proposalInfo.AgreeVotes = proposalInfo.AgreeVotes.Plus(votingPower)
	} else {
		proposalInfo.DisagreeVotes = proposalInfo.DisagreeVotes.Plus(votingPower)
	}

	proposal.SetProposalInfo(proposalInfo)
	if err := pm.storage.SetProposal(ctx, proposalID, proposal); err != nil {
		return err
	}

	return nil
}

func (pm ProposalManager) UpdateProposalPassStatus(
	ctx sdk.Context, proposalType types.ProposalType,
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

	// calculate if agree votes meet minimum pass requirement
	ratio, minVotes, err := pm.GetProposalPassParam(ctx, proposalType)
	if err != nil {
		return types.ProposalNotPass, err
	}
	totalVotes := proposalInfo.AgreeVotes.Plus(proposalInfo.DisagreeVotes)
	if !totalVotes.IsGT(minVotes) {
		return types.ProposalNotPass, nil
	}
	actualRatio := proposalInfo.AgreeVotes.ToRat().Quo(totalVotes.ToRat())
	if ratio.LT(actualRatio) {
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
		return nil, ErrIncorrectProposalType()
	}

	event := param.ChangeParamEvent{
		Param: p.Param,
	}
	return event, nil
}

func (pm ProposalManager) GetPermlink(ctx sdk.Context, proposalID types.ProposalKey) (types.Permlink, sdk.Error) {
	proposal, err := pm.storage.GetProposal(ctx, proposalID)
	if err != nil {
		return types.Permlink(""), err
	}

	p, ok := proposal.(*model.ContentCensorshipProposal)
	if !ok {
		return types.Permlink(""), ErrIncorrectProposalType()
	}
	return p.Permlink, nil
}

func (pm ProposalManager) GetProposalList(ctx sdk.Context) (*model.ProposalList, sdk.Error) {
	return pm.storage.GetProposalList(ctx)
}
