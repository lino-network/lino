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
	des model.ChangeParameterDescription, gm global.GlobalManager) (types.ProposalKey, sdk.Error) {
	newID, err := pm.paramHolder.GetNextProposalID(ctx)
	if err != nil {
		return newID, err
	}

	proposal := model.Proposal{
		Creator:      creator,
		ProposalID:   newID,
		AgreeVote:    types.Coin{Amount: 0},
		DisagreeVote: types.Coin{Amount: 0},
	}

	changeParameterProposal := &model.ChangeParameterProposal{
		Proposal:                   proposal,
		ChangeParameterDescription: des,
	}
	if err := pm.storage.SetProposal(ctx, newID, changeParameterProposal); err != nil {
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

func (pm ProposalManager) GetProposalList(ctx sdk.Context) (*model.ProposalList, sdk.Error) {
	return pm.storage.GetProposalList(ctx)
}

// func (vm VoteManager) CreateDecideProposalEvent(ctx sdk.Context, gm global.GlobalManager) sdk.Error {
// 	event := DecideProposalEvent{}
// 	gm.RegisterProposalDecideEvent(ctx, event)
// 	return nil
// }
