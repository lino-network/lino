package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/proposal/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProposalStatus(t *testing.T) {
	ctx, _, pm, _, _, _, _ := setupTest(t, 0)
	permLink := types.PermLink("postlink")
	user1 := types.AccountKey("user1")
	proposal1 := &model.ContentCensorshipProposal{
		PermLink: permLink,
	}

	proposal2 := &model.ContentCensorshipProposal{
		PermLink: permLink,
	}

	proposal3 := &model.ContentCensorshipProposal{
		PermLink: permLink,
	}
	pm.InitGenesis(ctx)
	proposalParam, _ := pm.paramHolder.GetProposalParam(ctx)
	proposalID1, _ := pm.AddProposal(ctx, "user1", proposal1)
	proposalID2, _ := pm.AddProposal(ctx, "user1", proposal2)
	proposalID3, _ := pm.AddProposal(ctx, "user1", proposal3)
	testCases := []struct {
		testName        string
		votingRes       types.VotingResult
		proposalType    types.ProposalType
		proposalID      types.ProposalKey
		wantProposalRes types.ProposalResult
		wantProposal    model.Proposal
	}{
		{testName: "test passed proposal has historical data",
			votingRes: types.VotingResult{
				AgreeVotes:    proposalParam.ContentCensorshipPassVotes,
				DisagreeVotes: proposalParam.ContentCensorshipPassVotes,
			},
			proposalType:    types.ContentCensorship,
			proposalID:      proposalID1,
			wantProposalRes: types.ProposalPass,
			wantProposal: &model.ContentCensorshipProposal{model.ProposalInfo{
				Creator:       user1,
				ProposalID:    proposalID1,
				AgreeVotes:    proposalParam.ContentCensorshipPassVotes,
				DisagreeVotes: proposalParam.ContentCensorshipPassVotes,
				Result:        types.ProposalPass,
			}, permLink},
		},

		{testName: "test votes don't meet min requirement ",
			votingRes: types.VotingResult{
				AgreeVotes:    proposalParam.ContentCensorshipPassVotes.Minus(types.NewCoin(10)),
				DisagreeVotes: types.NewCoin(0),
			},
			proposalType:    types.ContentCensorship,
			proposalID:      proposalID2,
			wantProposalRes: types.ProposalNotPass,
			wantProposal: &model.ContentCensorshipProposal{model.ProposalInfo{
				Creator:       user1,
				ProposalID:    proposalID2,
				AgreeVotes:    proposalParam.ContentCensorshipPassVotes.Minus(types.NewCoin(10)),
				DisagreeVotes: types.NewCoin(0),
				Result:        types.ProposalNotPass,
			}, permLink},
		},

		{testName: "test votes ratio doesn't meet requirement ",
			votingRes: types.VotingResult{
				AgreeVotes:    proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoin(10)),
				DisagreeVotes: proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoin(11)),
			},
			proposalType:    types.ContentCensorship,
			proposalID:      proposalID3,
			wantProposalRes: types.ProposalNotPass,
			wantProposal: &model.ContentCensorshipProposal{model.ProposalInfo{
				Creator:       user1,
				ProposalID:    proposalID3,
				AgreeVotes:    proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoin(10)),
				DisagreeVotes: proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoin(11)),
				Result:        types.ProposalPass,
			}, permLink},
		},
	}
	for _, tc := range testCases {
		res, err := pm.UpdateProposalStatus(ctx, tc.votingRes, tc.proposalType, tc.proposalID)
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
