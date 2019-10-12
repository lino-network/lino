package manager

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/param"
	linotypes "github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	accmn "github.com/lino-network/lino/x/account/manager"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/post"
	"github.com/lino-network/lino/x/proposal/model"
	"github.com/lino-network/lino/x/proposal/types"
	"github.com/lino-network/lino/x/vote"
)

// ProposalManager - proposal manager
type ProposalManager struct {
	storage model.ProposalStorage

	// deps
	paramHolder param.ParamHolder
	acc         acc.AccountKeeper
	post        post.PostKeeper
	global      global.GlobalKeeper
	vote        vote.VoteKeeper
}

// NewProposalManager - return a proposal manager
func NewProposalManager(key sdk.StoreKey, holder param.ParamHolder, vote vote.VoteKeeper,
	global global.GlobalKeeper, acc acc.AccountKeeper, post post.PostKeeper) ProposalManager {
	return ProposalManager{
		storage:     model.NewProposalStorage(key),
		paramHolder: holder,
		vote:        vote,
		global:      global,
		acc:         acc,
		post:        post,
	}
}

// InitGenesis - initialize proposal manager
func (pm ProposalManager) InitGenesis(ctx sdk.Context) error {
	if err := pm.storage.InitGenesis(ctx); err != nil {
		return err
	}
	return nil
}

func (pm ProposalManager) ChangeParam(ctx sdk.Context, creator linotypes.AccountKey,
	reason string, p param.Parameter) sdk.Error {
	if !pm.acc.DoesAccountExist(ctx, creator) {
		return types.ErrAccountNotFound()
	}

	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err
	}

	proposal := pm.CreateChangeParamProposal(ctx, p, reason)
	proposalID, err := pm.AddProposal(ctx, creator, proposal, param.ChangeParamDecideSec)
	if err != nil {
		return err
	}
	//  set a time event to decide the proposal
	event := pm.CreateDecideProposalEvent(ctx, linotypes.ChangeParam, proposalID)

	if err := pm.global.RegisterEventAtTime(
		ctx, ctx.BlockHeader().Time.Unix()+param.ChangeParamDecideSec, event); err != nil {
		return err
	}

	// minus coin from account and return when deciding the proposal
	if err = pm.acc.MinusCoinFromUsername(ctx, creator, param.ChangeParamMinDeposit); err != nil {
		return err
	}

	if err := pm.returnCoinTo(ctx, creator, int64(1),
		param.ChangeParamDecideSec, param.ChangeParamMinDeposit); err != nil {
		return err
	}
	return nil
}

func (pm ProposalManager) ProtocolUpgrade(ctx sdk.Context, creator linotypes.AccountKey,
	reason, link string) sdk.Error {
	if !pm.acc.DoesAccountExist(ctx, creator) {
		return types.ErrAccountNotFound()
	}

	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err
	}

	proposal := pm.CreateProtocolUpgradeProposal(ctx, link, reason)
	proposalID, err := pm.AddProposal(ctx, creator, proposal, param.ProtocolUpgradeDecideSec)
	if err != nil {
		return err
	}
	//  set a time event to decide the proposal
	event := pm.CreateDecideProposalEvent(ctx, linotypes.ProtocolUpgrade, proposalID)

	if err := pm.global.RegisterEventAtTime(
		ctx, ctx.BlockHeader().Time.Unix()+param.ProtocolUpgradeDecideSec, event); err != nil {
		return err
	}

	// minus coin from account and return when deciding the proposal
	if err = pm.acc.MinusCoinFromUsername(ctx, creator, param.ProtocolUpgradeMinDeposit); err != nil {
		return err
	}

	if err := pm.returnCoinTo(ctx, creator, int64(1),
		param.ProtocolUpgradeDecideSec, param.ProtocolUpgradeMinDeposit); err != nil {
		return err
	}
	return nil
}

func (pm ProposalManager) ContentCensorship(ctx sdk.Context, creator linotypes.AccountKey,
	reason string, permlink linotypes.Permlink) sdk.Error {
	if !pm.acc.DoesAccountExist(ctx, creator) {
		return types.ErrAccountNotFound()
	}

	if !pm.post.DoesPostExist(ctx, permlink) {
		return types.ErrPostNotFound()
	}

	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return err
	}

	proposal := pm.CreateContentCensorshipProposal(ctx, permlink, reason)
	proposalID, err := pm.AddProposal(ctx, creator, proposal, param.ContentCensorshipDecideSec)
	if err != nil {
		return err
	}
	//  set a time event to decide the proposal
	event := pm.CreateDecideProposalEvent(ctx, linotypes.ContentCensorship, proposalID)
	// minus coin from account and return when deciding the proposal
	if err = pm.acc.MinusCoinFromUsername(ctx, creator, param.ContentCensorshipMinDeposit); err != nil {
		return err
	}

	if err := pm.global.RegisterEventAtTime(
		ctx, ctx.BlockHeader().Time.Unix()+param.ContentCensorshipDecideSec, event); err != nil {
		return err
	}

	if err := pm.returnCoinTo(ctx, creator, int64(1),
		param.ContentCensorshipDecideSec, param.ContentCensorshipMinDeposit); err != nil {
		return err
	}
	return nil
}

// DoesProposalExist - check given proposal ID exists
func (pm ProposalManager) DoesProposalExist(ctx sdk.Context, proposalID linotypes.ProposalKey) bool {
	return pm.storage.DoesProposalExist(ctx, proposalID)
}

// IsOngoingProposal - check given proposal ID is in ongoing proposal list
func (pm ProposalManager) IsOngoingProposal(ctx sdk.Context, proposalID linotypes.ProposalKey) bool {
	_, err := pm.storage.GetOngoingProposal(ctx, proposalID)
	return err == nil
}

// CreateContentCensorshipProposal - create a content censorship proposal
func (pm ProposalManager) CreateContentCensorshipProposal(
	ctx sdk.Context, permlink linotypes.Permlink, reason string) model.Proposal {
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
func (pm ProposalManager) GetNextProposalID(ctx sdk.Context) (linotypes.ProposalKey, sdk.Error) {
	nextProposalID, err := pm.storage.GetNextProposalID(ctx)
	if err != nil {
		return linotypes.ProposalKey(""), err
	}

	return linotypes.ProposalKey(strconv.FormatInt(nextProposalID.NextProposalID, 10)), nil
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
	ctx sdk.Context, creator linotypes.AccountKey, proposal model.Proposal, decideSec int64) (linotypes.ProposalKey, sdk.Error) {
	newID, err := pm.GetNextProposalID(ctx)
	if err != nil {
		return newID, err
	}

	info := model.ProposalInfo{
		Creator:       creator,
		ProposalID:    newID,
		AgreeVotes:    linotypes.NewCoinFromInt64(0),
		DisagreeVotes: linotypes.NewCoinFromInt64(0),
		Result:        linotypes.ProposalNotPass,
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
	ctx sdk.Context, proposalType linotypes.ProposalType) (sdk.Dec, linotypes.Coin, sdk.Error) {
	param, err := pm.paramHolder.GetProposalParam(ctx)
	if err != nil {
		return sdk.NewDec(1), linotypes.NewCoinFromInt64(0), err
	}
	switch proposalType {
	case linotypes.ChangeParam:
		return param.ChangeParamPassRatio, param.ChangeParamPassVotes, nil
	case linotypes.ContentCensorship:
		return param.ContentCensorshipPassRatio, param.ContentCensorshipPassVotes, nil
	case linotypes.ProtocolUpgrade:
		return param.ProtocolUpgradePassRatio, param.ProtocolUpgradePassVotes, nil
	default:
		return sdk.NewDec(1), linotypes.NewCoinFromInt64(0), types.ErrIncorrectProposalType()
	}
}

// UpdateProposalVotingStatus - update proposal status after voting
func (pm ProposalManager) UpdateProposalVotingStatus(ctx sdk.Context, proposalID linotypes.ProposalKey,
	voter linotypes.AccountKey, voteResult bool, votingPower linotypes.Coin) sdk.Error {
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
	ctx sdk.Context, proposalType linotypes.ProposalType,
	proposalID linotypes.ProposalKey) (linotypes.ProposalResult, sdk.Error) {
	return pm.UpdateProposalStatus(ctx, proposalType, proposalID)
}

// UpdateProposalStatus - update proposal pass status when proposal change from ongoing to expired
func (pm ProposalManager) UpdateProposalStatus(
	ctx sdk.Context, proposalType linotypes.ProposalType,
	proposalID linotypes.ProposalKey) (linotypes.ProposalResult, sdk.Error) {
	proposal, err := pm.storage.GetOngoingProposal(ctx, proposalID)
	if err != nil {
		return linotypes.ProposalNotPass, err
	}

	proposalInfo := proposal.GetProposalInfo()

	// calculate if agree votes meet minimum pass requirement
	ratio, minVotes, err := pm.GetProposalPassParam(ctx, proposalType)
	if err != nil {
		return linotypes.ProposalNotPass, err
	}
	totalVotes := proposalInfo.AgreeVotes.Plus(proposalInfo.DisagreeVotes)
	actualRatio, err := sdk.NewDecFromStr("0")
	if err != nil {
		return linotypes.ProposalNotPass, err
	}
	if !totalVotes.IsZero() {
		actualRatio = proposalInfo.AgreeVotes.ToDec().Quo(totalVotes.ToDec())
	}

	if !totalVotes.IsGT(minVotes) || !ratio.LT(actualRatio) {
		proposalInfo.Result = linotypes.ProposalNotPass
	} else {
		proposalInfo.Result = linotypes.ProposalPass
	}

	proposal.SetProposalInfo(proposalInfo)
	if err := pm.storage.SetExpiredProposal(ctx, proposalID, proposal); err != nil {
		return linotypes.ProposalNotPass, err
	}

	if err := pm.storage.DeleteOngoingProposal(ctx, proposalID); err != nil {
		return linotypes.ProposalNotPass, err
	}
	return proposalInfo.Result, nil
}

// CreateDecideProposalEvent - create a decide proposal event
func (pm ProposalManager) CreateDecideProposalEvent(
	ctx sdk.Context, proposalType linotypes.ProposalType, proposalID linotypes.ProposalKey) linotypes.Event {
	event := DecideProposalEvent{
		ProposalType: proposalType,
		ProposalID:   proposalID,
	}
	return event
}

// CreateParamChangeEvent - create a parameter change event
func (pm ProposalManager) CreateParamChangeEvent(
	ctx sdk.Context, proposalID linotypes.ProposalKey) (linotypes.Event, sdk.Error) {
	proposal, err := pm.storage.GetExpiredProposal(ctx, proposalID)
	if err != nil {
		return nil, err
	}

	p, ok := proposal.(*model.ChangeParamProposal)
	if !ok {
		return nil, types.ErrIncorrectProposalType()
	}

	event := param.ChangeParamEvent{
		Param: p.Param,
	}
	return event, nil
}

// GetPermlink - get permlink from expired proposal list
func (pm ProposalManager) GetPermlink(ctx sdk.Context, proposalID linotypes.ProposalKey) (linotypes.Permlink, sdk.Error) {
	proposal, err := pm.storage.GetExpiredProposal(ctx, proposalID)
	if err != nil {
		return linotypes.Permlink(""), err
	}

	p, ok := proposal.(*model.ContentCensorshipProposal)
	if !ok {
		return linotypes.Permlink(""), types.ErrIncorrectProposalType()
	}
	return p.Permlink, nil
}

// GetOngoingProposalList - get ongoing proposal list
func (pm ProposalManager) GetOngoingProposalList(ctx sdk.Context) ([]model.Proposal, sdk.Error) {
	return pm.storage.GetOngoingProposalList(ctx)
}

func (pm ProposalManager) GetOngoingProposal(ctx sdk.Context, proposalID linotypes.ProposalKey) (model.Proposal, sdk.Error) {
	return pm.storage.GetOngoingProposal(ctx, proposalID)
}

func (pm ProposalManager) GetExpiredProposal(ctx sdk.Context, proposalID linotypes.ProposalKey) (model.Proposal, sdk.Error) {
	return pm.storage.GetExpiredProposal(ctx, proposalID)
}

func (pm ProposalManager) returnCoinTo(
	ctx sdk.Context, name linotypes.AccountKey, times int64, interval int64, coin linotypes.Coin) sdk.Error {
	if err := pm.acc.AddFrozenMoney(
		ctx, name, coin, ctx.BlockHeader().Time.Unix(), interval, times); err != nil {
		return err
	}

	events, err := accmn.CreateCoinReturnEvents(ctx, name, times, interval, coin, linotypes.ProposalReturnCoin)
	if err != nil {
		return err
	}

	if err := pm.global.RegisterCoinReturnEvent(ctx, events, times, interval); err != nil {
		return err
	}
	return nil
}
