package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestDecideProposal(t *testing.T) {
	ctx, am, pm, postManager, voteManager, valManager, gm := setupTest(t, 0)
	voteManager.InitGenesis(ctx)
	valManager.InitGenesis(ctx)
	pm.InitGenesis(ctx)
	proposalParam, _ := pm.paramHolder.GetProposalParam(ctx)

	c1, c2, c3, c4 :=
		proposalParam.ChangeParamPassVotes.Plus(types.NewCoinFromInt64(20)),
		proposalParam.ChangeParamPassVotes.Plus(types.NewCoinFromInt64(30)),
		proposalParam.ChangeParamPassVotes.Plus(types.NewCoinFromInt64(50)),
		proposalParam.ChangeParamPassVotes.Plus(types.NewCoinFromInt64(10))

	user1 := createTestAccount(ctx, am, "user1", c1)
	user2 := createTestAccount(ctx, am, "user2", c2)
	user3 := createTestAccount(ctx, am, "user3", c3)
	user4 := createTestAccount(ctx, am, "user4", c4)

	voteManager.AddVoter(ctx, user1, c1)
	voteManager.AddVoter(ctx, user2, c2)
	voteManager.AddVoter(ctx, user3, c3)
	voteManager.AddVoter(ctx, user4, c4)

	param1 := param.GlobalAllocationParam{
		InfraAllocation: sdk.NewRat(50, 100),
	}
	param2 := param.GlobalAllocationParam{
		InfraAllocation: sdk.NewRat(80, 100),
	}
	p1 := pm.CreateChangeParamProposal(ctx, param1)
	p2 := pm.CreateChangeParamProposal(ctx, param2)
	id1, _ := pm.AddProposal(ctx, types.AccountKey("c1"), p1)
	id2, _ := pm.AddProposal(ctx, types.AccountKey("c2"), p2)

	e1 := DecideProposalEvent{
		ProposalType: types.ChangeParam,
		ProposalID:   id1,
	}

	e2 := DecideProposalEvent{
		ProposalType: types.ChangeParam,
		ProposalID:   id2,
	}

	cases := []struct {
		event                 DecideProposalEvent
		decideProposal        bool
		voter                 types.AccountKey
		proposalID            types.ProposalKey
		voterRes              bool
		expectOngoingProposal []types.ProposalKey
		expectDecidedProposal []types.ProposalKey
		expectProposalRes     types.ProposalResult
		expectAgreeVotes      types.Coin
		expectDisagreeVotes   types.Coin
	}{
		{e1, false, user1, id1, true, []types.ProposalKey{id1, id2}, nil, types.ProposalNotPass,
			c1, c2},
		{e1, false, user2, id1, false, []types.ProposalKey{id1, id2}, nil, types.ProposalNotPass,
			c1, c2},
		{e1, true, types.AccountKey(""), id1, false, []types.ProposalKey{id2}, []types.ProposalKey{id1},
			types.ProposalNotPass, c1, c2},
		{e2, false, user1, id2, true, []types.ProposalKey{id2}, []types.ProposalKey{id1},
			types.ProposalNotPass, c1.Plus(c2).Plus(c4), c3},
		{e2, false, user2, id2, true, []types.ProposalKey{id2}, []types.ProposalKey{id1},
			types.ProposalNotPass, c1.Plus(c2).Plus(c4), c3},
		{e2, false, user4, id2, true, []types.ProposalKey{id2}, []types.ProposalKey{id1},
			types.ProposalNotPass, c1.Plus(c2).Plus(c4), c3},
		{e2, false, user3, id2, false, []types.ProposalKey{id2}, []types.ProposalKey{id1},
			types.ProposalNotPass, c1.Plus(c2).Plus(c4), c3},
		{e2, true, types.AccountKey(""), id2, false, nil, []types.ProposalKey{id1, id2},
			types.ProposalPass, c1.Plus(c2).Plus(c4), c3},
	}

	for _, cs := range cases {
		if cs.decideProposal {
			cs.event.Execute(ctx, voteManager, valManager, am, pm, postManager, gm)
			proposal, _ := pm.storage.GetProposal(ctx, cs.proposalID)
			proposalInfo := proposal.GetProposalInfo()

			assert.Equal(t, cs.expectProposalRes, proposalInfo.Result)
			assert.Equal(t, cs.expectAgreeVotes, proposalInfo.AgreeVotes)
			assert.Equal(t, cs.expectDisagreeVotes, proposalInfo.DisagreeVotes)

		} else {
			voteManager.AddVote(ctx, cs.proposalID, cs.voter, cs.voterRes)
		}

		lst, err := pm.storage.GetProposalList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectOngoingProposal, lst.OngoingProposal)
		assert.Equal(t, cs.expectDecidedProposal, lst.PastProposal)
	}
}
