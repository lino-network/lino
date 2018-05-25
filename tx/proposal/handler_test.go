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
	c460000 = types.NewCoinFromInt64(460000 * types.Decimals)
	c4600   = types.NewCoinFromInt64(4600 * types.Decimals)
	c46     = types.NewCoinFromInt64(46 * types.Decimals)
)

func TestChangeParamProposal(t *testing.T) {
	ctx, am, proposalManager, postManager, _, _, gm := setupTest(t, 0)
	handler := NewHandler(am, proposalManager, postManager, gm)
	proposalManager.InitGenesis(ctx)

	allocation := param.GlobalAllocationParam{
		ContentCreatorAllocation: sdk.Rat{Denom: 10, Num: 5},
	}
	proposalID1 := types.ProposalKey(strconv.FormatInt(int64(1), 10))
	proposalID2 := types.ProposalKey(strconv.FormatInt(int64(2), 10))

	user1 := createTestAccount(ctx, am, "user1", c460000)
	user2 := createTestAccount(ctx, am, "user2", c4600)
	proposalParam, _ := proposalManager.paramHolder.GetProposalParam(ctx)
	proposal1 := &model.ChangeParamProposal{model.ProposalInfo{
		Creator:       user1,
		ProposalID:    proposalID1,
		AgreeVotes:    types.NewCoinFromInt64(0),
		DisagreeVotes: types.NewCoinFromInt64(0),
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

		proposalList, _ := proposalManager.GetProposalList(ctx)
		assert.Equal(t, tc.wantOngoingProposal, proposalList.OngoingProposal)
		proposal, _ := proposalManager.storage.GetProposal(ctx, tc.proposalID)
		assert.Equal(t, tc.wantProposal, proposal)
	}
}

func TestContentCensorshipProposal(t *testing.T) {
	ctx, am, proposalManager, postManager, _, _, gm := setupTest(t, 0)
	handler := NewHandler(am, proposalManager, postManager, gm)
	proposalParam, _ := proposalManager.paramHolder.GetProposalParam(ctx)

	proposalManager.InitGenesis(ctx)

	proposalID1 := types.ProposalKey(strconv.FormatInt(int64(1), 10))
	//proposalID2 := types.ProposalKey(strconv.FormatInt(int64(2), 10))

	user1, postID1 := createTestPost(t, ctx, "user1", "postID", c460000, am, postManager, "0")
	user2, postID2 := createTestPost(t, ctx, "user2", "postID", c4600, am, postManager, "0")
	user3 := createTestAccount(
		ctx, am, "user3", proposalParam.ContentCensorshipMinDeposit.Minus(types.NewCoinFromInt64((1))))
	postManager.DeletePost(ctx, types.GetPermLink(user2, postID2))
	proposal1 := &model.ContentCensorshipProposal{model.ProposalInfo{
		Creator:       user2,
		ProposalID:    proposalID1,
		AgreeVotes:    types.NewCoinFromInt64(0),
		DisagreeVotes: types.NewCoinFromInt64(0),
		Result:        types.ProposalNotPass,
	}, types.GetPermLink(user1, postID1)}

	testCases := []struct {
		testName            string
		creator             types.AccountKey
		permLink            types.PermLink
		proposalID          types.ProposalKey
		wantOK              bool
		wantRes             sdk.Result
		wantCreatorBalance  types.Coin
		wantOngoingProposal []types.ProposalKey
		wantProposal        model.Proposal
	}{
		{testName: "user2 censorship user1's post successfully",
			creator:             user2,
			permLink:            types.GetPermLink(user1, postID1),
			proposalID:          proposalID1,
			wantOK:              true,
			wantRes:             sdk.Result{},
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{testName: "target post is not exist",
			creator:             user2,
			permLink:            types.GetPermLink(user1, "invalid"),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             ErrPostNotFound().Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{testName: "target post is deleted",
			creator:             user1,
			permLink:            types.GetPermLink(user2, postID2),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             ErrCensorshipPostIsDeleted(types.GetPermLink(user2, postID2)).Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{testName: "proposal is invalid",
			creator:             "invalid",
			permLink:            types.GetPermLink(user1, postID1),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             ErrUsernameNotFound().Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{testName: "user3 doesn't have enough money to create proposal",
			creator:             user3,
			permLink:            types.GetPermLink(user1, postID1),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             acc.ErrAccountSavingCoinNotEnough().Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
	}
	for _, tc := range testCases {
		msg := NewDeletePostContentMsg(string(tc.creator), tc.permLink)
		result := handler(ctx, msg)
		assert.Equal(t, tc.wantRes, result)

		if !tc.wantOK {
			continue
		}

		creatorBalance, _ := am.GetSavingFromBank(ctx, tc.creator)
		if !creatorBalance.IsEqual(tc.wantCreatorBalance) {
			t.Errorf("%s get creator bank balance(%v): got %v, want %v",
				tc.testName, msg.Creator, creatorBalance, tc.wantCreatorBalance)
		}

		proposalList, _ := proposalManager.GetProposalList(ctx)
		assert.Equal(t, tc.wantOngoingProposal, proposalList.OngoingProposal)
		proposal, _ := proposalManager.storage.GetProposal(ctx, tc.proposalID)
		assert.Equal(t, tc.wantProposal, proposal)
		permLink, err := proposalManager.GetPermLink(ctx, tc.proposalID)
		assert.Nil(t, err)
		assert.Equal(t, tc.permLink, permLink)
	}
}
