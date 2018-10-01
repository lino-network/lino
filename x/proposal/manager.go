package proposal

import (
	"strconv"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/proposal/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ProposalManager - proposal manager
type ProposalManager struct {
	storage     model.ProposalStorage
	paramHolder param.ParamHolder
}

// NewProposalManager - return a proposal manager
func NewProposalManager(key sdk.StoreKey, holder param.ParamHolder) ProposalManager {
	return ProposalManager{
		storage:     model.NewProposalStorage(key),
		paramHolder: holder,
	}
}

// InitGenesis - initialize proposal manager
func (pm ProposalManager) InitGenesis(ctx sdk.Context) error {
	if err := pm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

// DoesProposalExist - check given proposal ID exists
func (pm ProposalManager) DoesProposalExist(ctx sdk.Context, proposalID types.ProposalKey) bool {
	return pm.storage.DoesProposalExist(ctx, proposalID)
}

// IsOngoingProposal - check given proposal ID is in ongoing proposal list
func (pm ProposalManager) IsOngoingProposal(ctx sdk.Context, proposalID types.ProposalKey) bool {
	_, err := pm.storage.GetOngoingProposal(ctx, proposalID)
	return err == nil
}

// CreateContentCensorshipProposal - create a content censorship proposal
func (pm ProposalManager) CreateContentCensorshipProposal(
	ctx sdk.Context, permlink types.Permlink, reason string) model.Proposal {
	return &model.ContentCensorshipProposal{
		Permlink: permlink,
		Reason:   reason,
	}
}

// CreateProtocolUpgradeProposal - create a protocol upgrade proposal
func (pm ProposalManager) CreateProtocolUpgradeProposal(ctx sdk.Context, link string, reason string) model.Proposal {
	return &model.ProtocolUpgradeProposal{
		Link:   link,
		Reason: reason,
	}
}

// CreateChangeParamProposal - create a change parameters proposal
func (pm ProposalManager) CreateChangeParamProposal(
	ctx sdk.Context, parameter param.Parameter, reason string) model.Proposal {
	return &model.ChangeParamProposal{
		Param:  parameter,
		Reason: reason,
	}
}

// GetNextProposalID - get next proposal ID from KV store
func (pm ProposalManager) GetNextProposalID(ctx sdk.Context) (types.ProposalKey, sdk.Error) {
	nextProposalID, err := pm.storage.GetNextProposalID(ctx)
	if err != nil {
		return types.ProposalKey(""), err
	}

	return types.ProposalKey(strconv.FormatInt(nextProposalID.NextProposalID, 10)), nil
}

// IncreaseNextProposalID - increase next propsoal ID by 1 in KV store
func (pm ProposalManager) IncreaseNextProposalID(ctx sdk.Context) sdk.Error {
	nextProposalID, err := pm.storage.GetNextProposalID(ctx)
	if err != nil {
		return err
	}

	nextProposalID.NextProposalID++
	if err := pm.storage.SetNextProposalID(ctx, nextProposalID); err != nil {
		return err
	}

	return nil
}

// AddProposal - add a new proposal to ongoing proposal list
func (pm ProposalManager) AddProposal(
	ctx sdk.Context, creator types.AccountKey, proposal model.Proposal, decideSec int64) (types.ProposalKey, sdk.Error) {
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
		CreatedAt:     ctx.BlockHeader().Time.Unix(),
		ExpiredAt:     ctx.BlockHeader().Time.Unix() + decideSec,
	}
	proposal.SetProposalInfo(info)

	if err := pm.storage.SetOngoingProposal(ctx, newID, proposal); err != nil {
		return newID, err
	}

	if err := pm.IncreaseNextProposalID(ctx); err != nil {
		return newID, err
	}

	return newID, nil
}

// GetProposalPassParam - based on proposal type, get pass ratio and pass vote requirement
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

// UpdateProposalVotingStatus - update proposal status after voting
func (pm ProposalManager) UpdateProposalVotingStatus(ctx sdk.Context, proposalID types.ProposalKey,
	voter types.AccountKey, voteResult bool, votingPower types.Coin) sdk.Error {
	proposal, err := pm.storage.GetOngoingProposal(ctx, proposalID)
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
	if err := pm.storage.SetOngoingProposal(ctx, proposalID, proposal); err != nil {
		return err
	}

	return nil
}

// UpdateProposalPassStatus - update proposal pass status when proposal change from ongoing to expired
func (pm ProposalManager) UpdateProposalPassStatus(
	ctx sdk.Context, proposalType types.ProposalType,
	proposalID types.ProposalKey) (types.ProposalResult, sdk.Error) {
	if ctx.BlockHeader().Height > types.LinoBlockchainFirstUpdateHeight {
		return pm.UpdateProposalStatus(ctx, proposalType, proposalID)
	}
	proposal, err := pm.storage.GetOngoingProposal(ctx, proposalID)
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
	actualRatio := proposalInfo.AgreeVotes.ToRat().Quo(totalVotes.ToRat()).Round(types.PrecisionFactor)
	if ratio.LT(actualRatio) {
		proposalInfo.Result = types.ProposalPass
	} else {
		proposalInfo.Result = types.ProposalNotPass
	}

	proposal.SetProposalInfo(proposalInfo)
	if err := pm.storage.SetExpiredProposal(ctx, proposalID, proposal); err != nil {
		return types.ProposalNotPass, err
	}

	if err := pm.storage.DeleteOngoingProposal(ctx, proposalID); err != nil {
		return types.ProposalNotPass, err
	}
	return proposalInfo.Result, nil
}

// UpdateProposalStatus - update proposal pass status when proposal change from ongoing to expired
func (pm ProposalManager) UpdateProposalStatus(
	ctx sdk.Context, proposalType types.ProposalType,
	proposalID types.ProposalKey) (types.ProposalResult, sdk.Error) {
	proposal, err := pm.storage.GetOngoingProposal(ctx, proposalID)
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
	actualRatio := proposalInfo.AgreeVotes.ToRat().Quo(totalVotes.ToRat()).Round(types.PrecisionFactor)

	if !totalVotes.IsGT(minVotes) || !ratio.LT(actualRatio) {
		proposalInfo.Result = types.ProposalNotPass
	} else {
		proposalInfo.Result = types.ProposalPass
	}

	proposal.SetProposalInfo(proposalInfo)
	if err := pm.storage.SetExpiredProposal(ctx, proposalID, proposal); err != nil {
		return types.ProposalNotPass, err
	}

	if err := pm.storage.DeleteOngoingProposal(ctx, proposalID); err != nil {
		return types.ProposalNotPass, err
	}
	return proposalInfo.Result, nil
}

// CreateDecideProposalEvent - create a decide proposal event
func (pm ProposalManager) CreateDecideProposalEvent(
	ctx sdk.Context, proposalType types.ProposalType, proposalID types.ProposalKey) types.Event {
	event := DecideProposalEvent{
		ProposalType: proposalType,
		ProposalID:   proposalID,
	}
	return event
}

// CreateParamChangeEvent - create a parameter change event
func (pm ProposalManager) CreateParamChangeEvent(
	ctx sdk.Context, proposalID types.ProposalKey) (types.Event, sdk.Error) {
	proposal, err := pm.storage.GetExpiredProposal(ctx, proposalID)
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

// GetPermlink - get permlink from expired proposal list
func (pm ProposalManager) GetPermlink(ctx sdk.Context, proposalID types.ProposalKey) (types.Permlink, sdk.Error) {
	proposal, err := pm.storage.GetExpiredProposal(ctx, proposalID)
	if err != nil {
		return types.Permlink(""), err
	}

	p, ok := proposal.(*model.ContentCensorshipProposal)
	if !ok {
		return types.Permlink(""), ErrIncorrectProposalType()
	}
	return p.Permlink, nil
}

// GetOngoingProposalList - get ongoing proposal list
func (pm ProposalManager) GetOngoingProposalList(ctx sdk.Context) ([]model.Proposal, sdk.Error) {
	return pm.storage.GetOngoingProposalList(ctx)
}
