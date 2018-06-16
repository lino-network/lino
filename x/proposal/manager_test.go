package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/proposal/model"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProposalVotingStatus(t *testing.T) {
	ctx, _, pm, _, _, _, _ := setupTest(t, 0)
	permLink := types.PermLink("postlink")
	user1 := types.AccountKey("user1")
	censorshipReason := "reason"
	proposal1 := &model.ContentCensorshipProposal{
		PermLink: permLink,
		Reason:   censorshipReason,
	}

	pm.InitGenesis(ctx)
	curTime := ctx.BlockHeader().Time
	decideHr := int64(100)
	proposalID1, _ := pm.AddProposal(ctx, user1, proposal1, decideHr)

	testCases := []struct {
		testName     string
		proposalID   types.ProposalKey
		voter        types.AccountKey
		voteResult   bool
		votingPower  types.Coin
		wantProposal model.Proposal
	}{
		{
			testName:    "agree vote",
			proposalID:  proposalID1,
			voter:       user1,
			voteResult:  true,
			votingPower: types.NewCoinFromInt64(1),
			wantProposal: &model.ContentCensorshipProposal{model.ProposalInfo{
				Creator:       user1,
				ProposalID:    proposalID1,
				AgreeVotes:    types.NewCoinFromInt64(1),
				DisagreeVotes: types.NewCoinFromInt64(0),
				CreatedAt:     curTime,
				ExpiredAt:     curTime + decideHr*3600,
			}, permLink, censorshipReason},
		},
		{
			testName:    "one more agree vote",
			proposalID:  proposalID1,
			voter:       user1,
			voteResult:  true,
			votingPower: types.NewCoinFromInt64(2),
			wantProposal: &model.ContentCensorshipProposal{model.ProposalInfo{
				Creator:       user1,
				ProposalID:    proposalID1,
				AgreeVotes:    types.NewCoinFromInt64(3),
				DisagreeVotes: types.NewCoinFromInt64(0),
				CreatedAt:     curTime,
				ExpiredAt:     curTime + decideHr*3600,
			}, permLink, censorshipReason},
		},
		{
			testName:    "one disagree vote",
			proposalID:  proposalID1,
			voter:       user1,
			voteResult:  false,
			votingPower: types.NewCoinFromInt64(5),
			wantProposal: &model.ContentCensorshipProposal{model.ProposalInfo{
				Creator:       user1,
				ProposalID:    proposalID1,
				AgreeVotes:    types.NewCoinFromInt64(3),
				DisagreeVotes: types.NewCoinFromInt64(5),
				CreatedAt:     curTime,
				ExpiredAt:     curTime + decideHr*3600,
			}, permLink, censorshipReason},
		},
	}
	for _, tc := range testCases {
		err := pm.UpdateProposalVotingStatus(ctx, tc.proposalID, tc.voter, tc.voteResult, tc.votingPower)
		assert.Nil(t, err)

		proposal, _ := pm.storage.GetProposal(ctx, tc.proposalID)
		assert.Equal(t, tc.wantProposal, proposal)
	}

}

func TestUpdateProposalPassStatus(t *testing.T) {
	ctx, _, pm, _, _, _, _ := setupTest(t, 0)
	permLink := types.PermLink("postlink")
	user1 := types.AccountKey("user1")
	censorshipReason := "reason"
	proposal1 := &model.ContentCensorshipProposal{
		PermLink: permLink,
		Reason:   censorshipReason,
	}

	proposal2 := &model.ContentCensorshipProposal{
		PermLink: permLink,
		Reason:   censorshipReason,
	}

	proposal3 := &model.ContentCensorshipProposal{
		PermLink: permLink,
		Reason:   censorshipReason,
	}
	pm.InitGenesis(ctx)
	curTime := ctx.BlockHeader().Time
	decideHr := int64(100)
	proposalParam, _ := pm.paramHolder.GetProposalParam(ctx)
	proposalID1, _ := pm.AddProposal(ctx, user1, proposal1, decideHr)
	proposalID2, _ := pm.AddProposal(ctx, user1, proposal2, decideHr)
	proposalID3, _ := pm.AddProposal(ctx, user1, proposal3, decideHr)

	testCases := []struct {
		testName        string
		agreeVotes      types.Coin
		disagreeVotes   types.Coin
		proposalType    types.ProposalType
		proposalID      types.ProposalKey
		wantProposalRes types.ProposalResult
		wantProposal    model.Proposal
	}{
		{testName: "test passed proposal has historical data",
			agreeVotes:      proposalParam.ContentCensorshipPassVotes,
			disagreeVotes:   proposalParam.ContentCensorshipPassVotes,
			proposalType:    types.ContentCensorship,
			proposalID:      proposalID1,
			wantProposalRes: types.ProposalPass,
			wantProposal: &model.ContentCensorshipProposal{model.ProposalInfo{
				Creator:       user1,
				ProposalID:    proposalID1,
				AgreeVotes:    proposalParam.ContentCensorshipPassVotes,
				DisagreeVotes: proposalParam.ContentCensorshipPassVotes,
				Result:        types.ProposalPass,
				CreatedAt:     curTime,
				ExpiredAt:     curTime + decideHr*3600,
			}, permLink, censorshipReason},
		},

		{testName: "test votes don't meet min requirement ",
			agreeVotes:      proposalParam.ContentCensorshipPassVotes.Minus(types.NewCoinFromInt64(10)),
			disagreeVotes:   types.NewCoinFromInt64(0),
			proposalType:    types.ContentCensorship,
			proposalID:      proposalID2,
			wantProposalRes: types.ProposalNotPass,
			wantProposal: &model.ContentCensorshipProposal{model.ProposalInfo{
				Creator:       user1,
				ProposalID:    proposalID2,
				AgreeVotes:    proposalParam.ContentCensorshipPassVotes.Minus(types.NewCoinFromInt64(10)),
				DisagreeVotes: types.NewCoinFromInt64(0),
				Result:        types.ProposalNotPass,
				CreatedAt:     curTime,
				ExpiredAt:     curTime + decideHr*3600,
			}, permLink, censorshipReason},
		},

		{testName: "test votes ratio doesn't meet requirement ",
			agreeVotes:      proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoinFromInt64(10)),
			disagreeVotes:   proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoinFromInt64(11)),
			proposalType:    types.ContentCensorship,
			proposalID:      proposalID3,
			wantProposalRes: types.ProposalNotPass,
			wantProposal: &model.ContentCensorshipProposal{model.ProposalInfo{
				Creator:       user1,
				ProposalID:    proposalID3,
				AgreeVotes:    proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoinFromInt64(10)),
				DisagreeVotes: proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoinFromInt64(11)),
				Result:        types.ProposalPass,
				CreatedAt:     curTime,
				ExpiredAt:     curTime + decideHr*3600,
			}, permLink, censorshipReason},
		},
	}
	for _, tc := range testCases {
		err := addProposalInfo(ctx, pm, tc.proposalID, tc.agreeVotes, tc.disagreeVotes)
		assert.Nil(t, err)

		res, err := pm.UpdateProposalPassStatus(ctx, tc.proposalType, tc.proposalID)
		assert.Nil(t, err)
		assert.Equal(t, tc.wantProposalRes, res)
		if tc.wantProposalRes == types.ProposalNotPass {
			continue
		}
		proposal, _ := pm.storage.GetProposal(ctx, tc.proposalID)
		assert.Equal(t, tc.wantProposal, proposal)
	}
}

func TestGetProposalPassParam(t *testing.T) {
	ctx, _, pm, _, _, _, _ := setupTest(t, 0)

	proposalParam, _ := pm.paramHolder.GetProposalParam(ctx)
	testCases := []struct {
		testName      string
		proposalType  types.ProposalType
		wantError     sdk.Error
		wantPassRatio sdk.Rat
		wantPassVotes types.Coin
	}{
		{testName: "test pass param for changeParamProposal",
			proposalType:  types.ChangeParam,
			wantError:     nil,
			wantPassRatio: proposalParam.ChangeParamPassRatio,
			wantPassVotes: proposalParam.ChangeParamPassVotes,
		},

		{testName: "test pass param for contenCensorshipProposal",
			proposalType:  types.ContentCensorship,
			wantError:     nil,
			wantPassRatio: proposalParam.ContentCensorshipPassRatio,
			wantPassVotes: proposalParam.ContentCensorshipPassVotes,
		},

		{testName: "test pass param for protocolUpgradeProposal",
			proposalType:  types.ProtocolUpgrade,
			wantError:     nil,
			wantPassRatio: proposalParam.ProtocolUpgradePassRatio,
			wantPassVotes: proposalParam.ProtocolUpgradePassVotes,
		},

		{testName: "test wrong proposal type",
			proposalType:  23,
			wantError:     ErrWrongProposalType(),
			wantPassRatio: proposalParam.ProtocolUpgradePassRatio,
			wantPassVotes: proposalParam.ProtocolUpgradePassVotes,
		},
	}
	for _, tc := range testCases {
		ratio, votes, err := pm.GetProposalPassParam(ctx, tc.proposalType)
		assert.Equal(t, tc.wantError, err)
		if tc.wantError != nil {
			continue
		}
		assert.Equal(t, tc.wantPassRatio, ratio)
		assert.Equal(t, tc.wantPassVotes, votes)
	}

}
