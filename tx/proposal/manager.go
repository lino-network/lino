package proposal

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/tx/global"
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

// only support change parameter proposal now
func (pm ProposalManager) AddProposal(ctx sdk.Context, creator types.AccountKey,
	des model.Description, gm global.GlobalManager) (types.ProposalKey, sdk.Error) {
	newID, err := pm.paramHolder.GetNextProposalID(ctx)
	if err != nil {
		return newID, err
	}

	var proposal model.Proposal
	proposalInfo := model.ProposalInfo{
		Creator:       creator,
		ProposalID:    newID,
		AgreeVotes:    types.Coin{Amount: 0},
		DisagreeVotes: types.Coin{Amount: 0},
		Result:        types.ProposalNotPass,
	}

	switch des := des.(type) {
	case param.GlobalAllocationParam:
		proposal = model.ChangeGlobalAllocationParamProposal{proposalInfo, des}
	default:
		panic(des)
	}

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
	ctx sdk.Context, res types.VotingResult) (types.ProposalResult, sdk.Error) {
	lst, err := pm.storage.GetProposalList(ctx)
	if err != nil {
		return types.ProposalNotPass, err
	}

	curID := lst.OngoingProposal[0]
	proposal, err := pm.storage.GetProposal(ctx, curID)
	if err != nil {
		return types.ProposalNotPass, err
	}

	proposalInfoPtr := proposal.GetProposalInfo()
	if proposalInfoPtr == nil {
		return types.ProposalNotPass, ErrProposalInfoNotFound()
	}

	proposalInfoPtr.AgreeVotes = res.AgreeVotes
	proposalInfoPtr.DisagreeVotes = res.DisagreeVotes

	// TODO consider different types of propsal
	if proposalInfoPtr.AgreeVotes.IsGT(proposalInfoPtr.DisagreeVotes) {
		proposalInfoPtr.Result = types.ProposalPass
	}

	if err := pm.storage.SetProposal(ctx, curID, proposal); err != nil {
		return types.ProposalNotPass, err
	}

	lst.OngoingProposal = lst.OngoingProposal[1:]
	lst.PastProposal = append(lst.PastProposal, curID)

	if err := pm.storage.SetProposalList(ctx, lst); err != nil {
		return types.ProposalNotPass, err
	}
	return proposalInfoPtr.Result, nil
}

func (pm ProposalManager) CreateDecideProposalEvent(ctx sdk.Context, gm global.GlobalManager) sdk.Error {
	event := DecideProposalEvent{}
	if err := gm.RegisterProposalDecideEvent(ctx, event); err != nil {
		return err
	}
	return nil
}

func (pm ProposalManager) CreateParamChangeEvent(
	ctx sdk.Context, proposalID types.ProposalKey, gm global.GlobalManager) sdk.Error {
	proposal, err := pm.storage.GetProposal(ctx, proposalID)
	if err != nil {
		return err
	}

	var event types.Event
	switch proposal := proposal.(type) {
	case model.ChangeGlobalAllocationParamProposal:
		event = param.ChangeGlobalAllocationParamEvent{proposal.Description}
	default:
		panic("err")
	}

	if err := gm.RegisterParamChangeEvent(ctx, event); err != nil {
		return err
	}
	return nil
}

func (pm ProposalManager) GetProposalList(ctx sdk.Context) (*model.ProposalList, sdk.Error) {
	return pm.storage.GetProposalList(ctx)
}
