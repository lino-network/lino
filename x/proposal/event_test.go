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

	p1 := pm.CreateChangeParamProposal(ctx, param1, "")
	p2 := pm.CreateChangeParamProposal(ctx, param2, "")
	id1, _ := pm.AddProposal(ctx, types.AccountKey("c1"), p1, 10)
	id2, _ := pm.AddProposal(ctx, types.AccountKey("c2"), p2, 10)

	e1 := DecideProposalEvent{
		ProposalType: types.ChangeParam,
		ProposalID:   id1,
	}

	e2 := DecideProposalEvent{
		ProposalType: types.ChangeParam,
		ProposalID:   id2,
	}

	cases := []struct {
		testName              string
		event                 DecideProposalEvent
		decideProposal        bool
		voter                 types.AccountKey
		proposalID            types.ProposalKey
		voterRes              bool
		votingPower           types.Coin
		expectOngoingProposal []types.ProposalKey
		expectDecidedProposal []types.ProposalKey
		expectProposalRes     types.ProposalResult
		expectAgreeVotes      types.Coin
		expectDisagreeVotes   types.Coin
	}{
		{
			testName:              "test1",
			event:                 e1,
			decideProposal:        false,
			voter:                 user1,
			proposalID:            id1,
			voterRes:              true,
			votingPower:           c1,
			expectOngoingProposal: []types.ProposalKey{id1, id2},
			expectDecidedProposal: nil,
			expectProposalRes:     types.ProposalNotPass,
			expectAgreeVotes:      c1,
			expectDisagreeVotes:   c2,
		},
		{
			testName:              "test2",
			event:                 e1,
			decideProposal:        false,
			voter:                 user2,
			proposalID:            id1,
			voterRes:              false,
			votingPower:           c2,
			expectOngoingProposal: []types.ProposalKey{id1, id2},
			expectDecidedProposal: nil,
			expectProposalRes:     types.ProposalNotPass,
			expectAgreeVotes:      c1,
			expectDisagreeVotes:   c2,
		},
		{
			testName:              "test3",
			event:                 e1,
			decideProposal:        true,
			voter:                 types.AccountKey(""),
			proposalID:            id1,
			voterRes:              false,
			expectOngoingProposal: []types.ProposalKey{id2},
			expectDecidedProposal: []types.ProposalKey{id1},
			expectProposalRes:     types.ProposalNotPass,
			expectAgreeVotes:      c1,
			expectDisagreeVotes:   c2,
		},
		{
			testName:              "test4",
			event:                 e2,
			decideProposal:        false,
			voter:                 user1,
			proposalID:            id2,
			voterRes:              true,
			votingPower:           c1,
			expectOngoingProposal: []types.ProposalKey{id2},
			expectDecidedProposal: []types.ProposalKey{id1},
			expectProposalRes:     types.ProposalNotPass,
			expectAgreeVotes:      c1.Plus(c2).Plus(c4),
			expectDisagreeVotes:   c3,
		},
		{
			testName:              "test5",
			event:                 e2,
			decideProposal:        false,
			voter:                 user2,
			proposalID:            id2,
			voterRes:              true,
			votingPower:           c2,
			expectOngoingProposal: []types.ProposalKey{id2},
			expectDecidedProposal: []types.ProposalKey{id1},
			expectProposalRes:     types.ProposalNotPass,
			expectAgreeVotes:      c1.Plus(c2).Plus(c4),
			expectDisagreeVotes:   c3,
		},
		{
			testName:              "test6",
			event:                 e2,
			decideProposal:        false,
			voter:                 user4,
			proposalID:            id2,
			voterRes:              true,
			votingPower:           c4,
			expectOngoingProposal: []types.ProposalKey{id2},
			expectDecidedProposal: []types.ProposalKey{id1},
			expectProposalRes:     types.ProposalNotPass,
			expectAgreeVotes:      c1.Plus(c2).Plus(c4),
			expectDisagreeVotes:   c3,
		},
		{
			testName:              "test7",
			event:                 e2,
			decideProposal:        false,
			voter:                 user3,
			proposalID:            id2,
			voterRes:              false,
			votingPower:           c3,
			expectOngoingProposal: []types.ProposalKey{id2},
			expectDecidedProposal: []types.ProposalKey{id1},
			expectProposalRes:     types.ProposalNotPass,
			expectAgreeVotes:      c1.Plus(c2).Plus(c4),
			expectDisagreeVotes:   c3,
		},
		{
			testName:              "test8",
			event:                 e2,
			decideProposal:        true,
			voter:                 types.AccountKey(""),
			proposalID:            id2,
			voterRes:              false,
			expectOngoingProposal: nil,
			expectDecidedProposal: []types.ProposalKey{id1, id2},
			expectProposalRes:     types.ProposalPass,
			expectAgreeVotes:      c1.Plus(c2).Plus(c4),
			expectDisagreeVotes:   c3,
		},
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

			err := pm.UpdateProposalVotingStatus(ctx, cs.proposalID, cs.voter, cs.voterRes, cs.votingPower)
			assert.Nil(t, err)
		}

		lst, err := pm.storage.GetProposalList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectOngoingProposal, lst.OngoingProposal)
		assert.Equal(t, cs.expectDecidedProposal, lst.PastProposal)
	}
}
