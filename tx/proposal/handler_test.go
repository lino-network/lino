package proposal

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/proposal/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

var (
	c460000 = types.Coin{460000 * types.Decimals}
	c4600   = types.Coin{4600 * types.Decimals}
	c46     = types.Coin{46 * types.Decimals}
)

func TestChangeParamProposal(t *testing.T) {
	ctx, am, pm, _, _, _, gm := setupTest(t, 0)
	handler := NewHandler(am, pm, gm)
	pm.InitGenesis(ctx)

	allocation := param.GlobalAllocationParam{
		ContentCreatorAllocation: sdk.Rat{Denom: 10, Num: 5},
	}
	proposalID1 := types.ProposalKey(strconv.FormatInt(int64(1), 10))
	proposalID2 := types.ProposalKey(strconv.FormatInt(int64(2), 10))

	user1 := createTestAccount(ctx, am, "user1", c460000)
	user2 := createTestAccount(ctx, am, "user2", c4600)
	proposalParam, _ := pm.paramHolder.GetProposalParam(ctx)
	proposal1 := &model.ChangeParamProposal{model.ProposalInfo{
		Creator:       user1,
		ProposalID:    proposalID1,
		AgreeVotes:    types.Coin{Amount: 0},
		DisagreeVotes: types.Coin{Amount: 0},
		Result:        types.ProposalNotPass,
	}, allocation}

	testCases := []struct {
		testName            string
		msg                 ChangeGlobalAllocationParamMsg
		proposalID          types.ProposalKey
		wantOK              bool
		wantRes             sdk.Result
		wantCreatorBalance  types.Coin
		wantOngoingProposal []types.ProposalKey
		wantProposal        model.Proposal
	}{
		{testName: "user1 creates change param msg successfully",
			msg: ChangeGlobalAllocationParamMsg{
				Creator:   user1,
				Parameter: allocation,
			},
			proposalID:          proposalID1,
			wantOK:              true,
			wantRes:             sdk.Result{},
			wantCreatorBalance:  c460000.Minus(proposalParam.ChangeParamMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},

		{testName: "user2 doesn't have enough money to create proposal",
			msg: ChangeGlobalAllocationParamMsg{
				Creator:   user2,
				Parameter: allocation,
			},
			proposalID:          proposalID2,
			wantOK:              false,
			wantRes:             acc.ErrAccountSavingCoinNotEnough().Result(),
			wantCreatorBalance:  c4600,
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        nil,
		},
	}
	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		assert.Equal(t, tc.wantRes, result)

		if !tc.wantOK {
			continue
		}

		creatorBalance, _ := am.GetSavingFromBank(ctx, tc.msg.GetCreator())
		if !creatorBalance.IsEqual(tc.wantCreatorBalance) {
			t.Errorf("%s get creator bank balance(%v): got %v, want %v", tc.testName, tc.msg.Creator, creatorBalance, tc.wantCreatorBalance)
		}

		proposalList, _ := pm.GetProposalList(ctx)
		assert.Equal(t, tc.wantOngoingProposal, proposalList.OngoingProposal)
		proposal, _ := pm.storage.GetProposal(ctx, tc.proposalID)
		assert.Equal(t, tc.wantProposal, proposal)
	}

}

func TestContenCencorshipProposal(t *testing.T) {
	ctx, am, pm, _, _, _, gm := setupTest(t, 0)
	handler := NewHandler(am, pm, gm)
	pm.InitGenesis(ctx)

	permLink := types.PermLink("postlink")
	proposalID1 := types.ProposalKey(strconv.FormatInt(int64(1), 10))
	proposalID2 := types.ProposalKey(strconv.FormatInt(int64(2), 10))

	user1 := createTestAccount(ctx, am, "user1", c460000)
	user2 := createTestAccount(ctx, am, "user2", c46)
	proposalParam, _ := pm.paramHolder.GetProposalParam(ctx)
	proposal1 := &model.ContentCensorshipProposal{model.ProposalInfo{
		Creator:       user1,
		ProposalID:    proposalID1,
		AgreeVotes:    types.Coin{Amount: 0},
		DisagreeVotes: types.Coin{Amount: 0},
		Result:        types.ProposalNotPass,
	}, permLink}

	testCases := []struct {
		testName            string
		msg                 DeletePostContentMsg
		proposalID          types.ProposalKey
		wantOK              bool
		wantRes             sdk.Result
		wantCreatorBalance  types.Coin
		wantOngoingProposal []types.ProposalKey
		wantProposal        model.Proposal
	}{
		{testName: "user1 creates conten censorship msg successfully",
			msg: DeletePostContentMsg{
				Creator:  user1,
				PermLink: permLink,
			},
			proposalID:          proposalID1,
			wantOK:              true,
			wantRes:             sdk.Result{},
			wantCreatorBalance:  c460000.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},

		{testName: "user2 doesn't have enough money to create proposal",
			msg: DeletePostContentMsg{
				Creator:  user2,
				PermLink: permLink,
			},
			proposalID:          proposalID2,
			wantOK:              false,
			wantRes:             acc.ErrAccountSavingCoinNotEnough().Result(),
			wantCreatorBalance:  c46,
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        nil,
		},
	}
	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		assert.Equal(t, tc.wantRes, result)

		if !tc.wantOK {
			continue
		}

		creatorBalance, _ := am.GetSavingFromBank(ctx, tc.msg.GetCreator())
		if !creatorBalance.IsEqual(tc.wantCreatorBalance) {
			t.Errorf("%s get creator bank balance(%v): got %v, want %v", tc.testName, tc.msg.Creator, creatorBalance, tc.wantCreatorBalance)
		}

		proposalList, _ := pm.GetProposalList(ctx)
		assert.Equal(t, tc.wantOngoingProposal, proposalList.OngoingProposal)
		proposal, _ := pm.storage.GetProposal(ctx, tc.proposalID)
		assert.Equal(t, tc.wantProposal, proposal)
	}

}
