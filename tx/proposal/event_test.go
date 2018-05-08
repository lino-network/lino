package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestDecideProposal(t *testing.T) {
	ctx, am, pm, voteManager, valManager, gm := setupTest(t, 0)
	voteManager.InitGenesis(ctx)
	valManager.InitGenesis(ctx)
	pm.InitGenesis(ctx)

	user1 := createTestAccount(ctx, am, "user1")
	user2 := createTestAccount(ctx, am, "user2")
	user3 := createTestAccount(ctx, am, "user3")
	user4 := createTestAccount(ctx, am, "user4")

	voteParam, _ := pm.paramHolder.GetVoteParam(ctx)

	voteManager.AddVoter(ctx, user1, voteParam.VoterMinDeposit.Plus(types.NewCoin(20)))
	voteManager.AddVoter(ctx, user2, voteParam.VoterMinDeposit.Plus(types.NewCoin(30)))
	voteManager.AddVoter(ctx, user3, voteParam.VoterMinDeposit.Plus(types.NewCoin(50)))
	voteManager.AddVoter(ctx, user4, voteParam.VoterMinDeposit.Plus(types.NewCoin(10)))
	e := DecideProposalEvent{}
	des1 := param.GlobalAllocationParam{
		InfraAllocation: sdk.NewRat(50, 100),
	}
	des2 := param.GlobalAllocationParam{
		InfraAllocation: sdk.NewRat(80, 100),
	}
	id1, _ := pm.AddProposal(ctx, types.AccountKey("c1"), des1, gm)
	id2, _ := pm.AddProposal(ctx, types.AccountKey("c2"), des2, gm)

	cases := []struct {
		decideProposal        bool
		voter                 types.AccountKey
		proposalID            types.ProposalKey
		voterRes              bool
		expectOngoingProposal []types.ProposalKey
		expectDecidedProposal []types.ProposalKey
		expectProposalRes     types.ProposalResult
	}{
		{false, user1, id1, true, []types.ProposalKey{id1, id2}, nil, types.ProposalNotPass},
		{false, user2, id1, false, []types.ProposalKey{id1, id2}, nil, types.ProposalNotPass},
		{true, types.AccountKey(""), id1, false, []types.ProposalKey{id2}, []types.ProposalKey{id1}, types.ProposalNotPass},
		{false, user1, id2, true, []types.ProposalKey{id2}, []types.ProposalKey{id1}, types.ProposalNotPass},
		{false, user2, id2, true, []types.ProposalKey{id2}, []types.ProposalKey{id1}, types.ProposalNotPass},
		{false, user4, id2, true, []types.ProposalKey{id2}, []types.ProposalKey{id1}, types.ProposalNotPass},
		{false, user3, id2, false, []types.ProposalKey{id2}, []types.ProposalKey{id1}, types.ProposalNotPass},
		{true, types.AccountKey(""), id2, false, nil, []types.ProposalKey{id1, id2}, types.ProposalNotPass},
	}

	for _, cs := range cases {
		if cs.decideProposal {
			e.Execute(ctx, voteManager, valManager, am, pm, gm)
		} else {
			voteManager.AddVote(ctx, cs.proposalID, cs.voter, cs.voterRes)
		}

		proposal, _ := pm.storage.GetProposal(ctx, cs.proposalID)
		proposalInfoPtr := proposal.GetProposalInfo()
		assert.NotNil(t, proposalInfoPtr)
		assert.Equal(t, cs.expectProposalRes, proposalInfoPtr.Result)
		lst, err := pm.storage.GetProposalList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectOngoingProposal, lst.OngoingProposal)
		assert.Equal(t, cs.expectDecidedProposal, lst.PastProposal)
	}
}

// func TestForceValidatorVote(t *testing.T) {
// 	ctx, am, pm, voteManager, valManager, gm := setupTest(t, 0)
// 	vm.InitGenesis(ctx)
// 	user1 := createTestAccount(ctx, am, "user1")
// 	user2 := createTestAccount(ctx, am, "user2")
//
// 	voterMinDeposit, _ := gm.GetVoterMinDeposit(ctx)
// 	vm.AddVoter(ctx, user1, voterMinDeposit.Plus(types.NewCoin(20)), gm)
// 	vm.AddVoter(ctx, user2, voterMinDeposit.Plus(types.NewCoin(30)), gm)
//
// 	referenceList := &model.ValidatorReferenceList{
// 		OncallValidators: []types.AccountKey{user2, user1},
// 	}
// 	vm.storage.SetValidatorReferenceList(ctx, referenceList)
//
// 	e := DecideProposalEvent{}
//
// 	des1 := &model.ChangeParameterDescription{
// 		InfraAllocation: sdk.NewRat(50, 100),
// 	}
// 	id1, _ := vm.AddProposal(ctx, types.AccountKey("c1"), des1, gm)
// 	cases := []struct {
// 		decideProposal    bool
// 		voter             types.AccountKey
// 		proposalID        types.ProposalKey
// 		voterRes          bool
// 		expectPenaltyList []types.AccountKey
// 	}{
// 		{false, user1, id1, true, nil},
// 		{true, types.AccountKey(""), id1, true, []types.AccountKey{user2}},
// 	}
//
// 	for _, cs := range cases {
// 		if cs.decideProposal {
// 			e.Execute(ctx, vm, am, gm)
// 		} else {
// 			vm.AddVote(ctx, cs.proposalID, cs.voter, cs.voterRes)
// 		}
// 		lst, _ := vm.GetValidatorReferenceList(ctx)
// 		assert.Equal(t, cs.expectPenaltyList, lst.PenaltyValidators)
// 	}
// }
