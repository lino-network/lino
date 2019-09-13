package proposal

import (
	"testing"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/proposal/model"
	"github.com/stretchr/testify/assert"
)

func TestDecideProposal(t *testing.T) {
	ctx, am, pm, postManager, voteManager, valManager, gm := setupTest(t, 0)
	err := voteManager.InitGenesis(ctx)
	if err != nil {
		panic(err)
	}
	err = valManager.InitGenesis(ctx)
	if err != nil {
		panic(err)
	}
	err = pm.InitGenesis(ctx)
	if err != nil {
		panic(err)
	}
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

	err = voteManager.AddVoter(ctx, user1, c1)
	if err != nil {
		panic(err)
	}
	err = voteManager.AddVoter(ctx, user2, c2)
	if err != nil {
		panic(err)
	}
	err = voteManager.AddVoter(ctx, user3, c3)
	if err != nil {
		panic(err)
	}
	err = voteManager.AddVoter(ctx, user4, c4)
	if err != nil {
		panic(err)
	}

	param1 := param.GlobalAllocationParam{
		InfraAllocation: types.NewDecFromRat(50, 100),
	}
	param2 := param.GlobalAllocationParam{
		InfraAllocation: types.NewDecFromRat(80, 100),
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
			err := cs.event.Execute(ctx, voteManager, valManager, am, pm, postManager, &gm)
			assert.Nil(t, err)
			proposal, _ := pm.storage.GetExpiredProposal(ctx, cs.proposalID)
			proposalInfo := proposal.GetProposalInfo()

			assert.Equal(t, cs.expectProposalRes, proposalInfo.Result)
			assert.Equal(t, cs.expectAgreeVotes, proposalInfo.AgreeVotes)
			assert.Equal(t, cs.expectDisagreeVotes, proposalInfo.DisagreeVotes)

		} else {
			err := voteManager.AddVote(ctx, cs.proposalID, cs.voter, cs.voterRes)
			if err != nil {
				panic(err)
			}

			err = pm.UpdateProposalVotingStatus(ctx, cs.proposalID, cs.voter, cs.voterRes, cs.votingPower)
			assert.Nil(t, err)
		}

		ongoingList, _ := pm.storage.GetOngoingProposalList(ctx)
		expiredList, _ := pm.storage.GetExpiredProposalList(ctx)

		var expectOngoingProposalList []model.Proposal
		var expectExpiredProposalList []model.Proposal

		for i := 0; i < len(cs.expectOngoingProposal); i++ {
			p, _ := pm.storage.GetOngoingProposal(ctx, cs.expectOngoingProposal[i])
			expectOngoingProposalList = append(expectOngoingProposalList, p)
		}
		for i := 0; i < len(cs.expectDecidedProposal); i++ {
			p, _ := pm.storage.GetExpiredProposal(ctx, cs.expectDecidedProposal[i])
			expectExpiredProposalList = append(expectExpiredProposalList, p)
		}

		assert.Equal(t, expectOngoingProposalList, ongoingList)
		assert.Equal(t, expectExpiredProposalList, expiredList)
	}
}
