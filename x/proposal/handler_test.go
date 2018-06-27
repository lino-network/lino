package proposal

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	"github.com/lino-network/lino/x/proposal/model"
	"github.com/lino-network/lino/x/vote"
	"github.com/stretchr/testify/assert"
)

var (
	c460000 = types.NewCoinFromInt64(460000 * types.Decimals)
	c4600   = types.NewCoinFromInt64(4600 * types.Decimals)
	c46     = types.NewCoinFromInt64(46 * types.Decimals)
)

func TestChangeParamProposal(t *testing.T) {
	ctx, am, proposalManager, postManager, vm, _, gm := setupTest(t, 0)
	handler := NewHandler(am, proposalManager, postManager, gm, vm)
	proposalManager.InitGenesis(ctx)

	allocation := param.GlobalAllocationParam{
		DeveloperAllocation:      sdk.ZeroRat(),
		ValidatorAllocation:      sdk.ZeroRat(),
		InfraAllocation:          sdk.ZeroRat(),
		ContentCreatorAllocation: sdk.NewRat(5, 10),
	}
	proposalID1 := types.ProposalKey(strconv.FormatInt(int64(1), 10))
	proposalID2 := types.ProposalKey(strconv.FormatInt(int64(2), 10))

	user1 := createTestAccount(ctx, am, "user1", c460000)
	user2 := createTestAccount(ctx, am, "user2", c4600)

	curTime := ctx.BlockHeader().Time
	proposalParam, _ := proposalManager.paramHolder.GetProposalParam(ctx)

	proposal1 := &model.ChangeParamProposal{model.ProposalInfo{
		Creator:       user1,
		ProposalID:    proposalID1,
		AgreeVotes:    types.NewCoinFromInt64(0),
		DisagreeVotes: types.NewCoinFromInt64(0),
		Result:        types.ProposalNotPass,
		CreatedAt:     curTime,
		ExpiredAt:     curTime + proposalParam.ChangeParamDecideHr*3600,
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
	ctx, am, proposalManager, postManager, vm, _, gm := setupTest(t, 0)
	handler := NewHandler(am, proposalManager, postManager, gm, vm)
	curTime := ctx.BlockHeader().Time
	proposalParam, _ := proposalManager.paramHolder.GetProposalParam(ctx)

	proposalManager.InitGenesis(ctx)

	proposalID1 := types.ProposalKey(strconv.FormatInt(int64(1), 10))
	//proposalID2 := types.ProposalKey(strconv.FormatInt(int64(2), 10))

	user1, postID1 := createTestPost(t, ctx, "user1", "postID", c460000, am, postManager, "0")
	user2, postID2 := createTestPost(t, ctx, "user2", "postID", c4600, am, postManager, "0")
	user3 := createTestAccount(
		ctx, am, "user3", proposalParam.ContentCensorshipMinDeposit.Minus(types.NewCoinFromInt64((1))))
	postManager.DeletePost(ctx, types.GetPermlink(user2, postID2))
	censorshipReason := "reason"
	proposal1 := &model.ContentCensorshipProposal{model.ProposalInfo{
		Creator:       user2,
		ProposalID:    proposalID1,
		AgreeVotes:    types.NewCoinFromInt64(0),
		DisagreeVotes: types.NewCoinFromInt64(0),
		Result:        types.ProposalNotPass,
		CreatedAt:     curTime,
		ExpiredAt:     curTime + proposalParam.ChangeParamDecideHr*3600,
	}, types.GetPermlink(user1, postID1), censorshipReason}

	testCases := []struct {
		testName            string
		creator             types.AccountKey
		permLink            types.Permlink
		proposalID          types.ProposalKey
		wantOK              bool
		wantRes             sdk.Result
		wantCreatorBalance  types.Coin
		wantOngoingProposal []types.ProposalKey
		wantProposal        model.Proposal
	}{
		{testName: "user2 censorship user1's post successfully",
			creator:             user2,
			permLink:            types.GetPermlink(user1, postID1),
			proposalID:          proposalID1,
			wantOK:              true,
			wantRes:             sdk.Result{},
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{testName: "target post is not exist",
			creator:             user2,
			permLink:            types.GetPermlink(user1, "invalid"),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             ErrPostNotFound().Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{testName: "target post is deleted",
			creator:             user1,
			permLink:            types.GetPermlink(user2, postID2),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             ErrCensorshipPostIsDeleted(types.GetPermlink(user2, postID2)).Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{testName: "proposal is invalid",
			creator:             "invalid",
			permLink:            types.GetPermlink(user1, postID1),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             ErrUsernameNotFound().Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{testName: "user3 doesn't have enough money to create proposal",
			creator:             user3,
			permLink:            types.GetPermlink(user1, postID1),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             acc.ErrAccountSavingCoinNotEnough().Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
	}
	for _, tc := range testCases {
		msg := NewDeletePostContentMsg(string(tc.creator), tc.permLink, censorshipReason)
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
		permLink, err := proposalManager.GetPermlink(ctx, tc.proposalID)
		assert.Nil(t, err)
		assert.Equal(t, tc.permLink, permLink)
	}
}

func TestAddFrozenMoney(t *testing.T) {
	ctx, am, proposalManager, _, _, _, gm := setupTest(t, 0)
	proposalManager.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user := createTestAccount(ctx, am, "user", minBalance)

	testCases := []struct {
		testName               string
		times                  int64
		interval               int64
		returnedCoin           types.Coin
		expectedFrozenListLen  int
		expectedFrozenMoney    types.Coin
		expectedFrozenTimes    int64
		expectedFrozenInterval int64
	}{
		{"return coin to user", 10, 2, types.NewCoinFromInt64(100), 1, types.NewCoinFromInt64(100), 10, 2},
		{"return coin to user multiple times", 100000, 20000, types.NewCoinFromInt64(100000), 2, types.NewCoinFromInt64(100000), 100000, 20000},
	}

	for _, tc := range testCases {
		err := returnCoinTo(
			ctx, "user", gm, am, tc.times, tc.interval, tc.returnedCoin)
		assert.Equal(t, nil, err)
		lst, err := am.GetFrozenMoneyList(ctx, user)
		assert.Equal(t, tc.expectedFrozenListLen, len(lst))
		assert.Equal(t, tc.expectedFrozenMoney, lst[len(lst)-1].Amount)
		assert.Equal(t, tc.expectedFrozenTimes, lst[len(lst)-1].Times)
		assert.Equal(t, tc.expectedFrozenInterval, lst[len(lst)-1].Interval)
	}
}

func TestVoteProposalBasic(t *testing.T) {
	ctx, am, proposalManager, postManager, vm, _, gm := setupTest(t, 0)
	handler := NewHandler(am, proposalManager, postManager, gm, vm)
	curTime := ctx.BlockHeader().Time
	proposalManager.InitGenesis(ctx)

	user2 := types.AccountKey("user2")

	// create voter
	user1 := types.AccountKey("user1")
	createTestAccount(ctx, am, "user1", c4600)
	_ = vm.AddVoter(ctx, user1, c4600)

	// create proposal
	permLink := types.Permlink("postlink")
	censorshipReason := "reason"
	proposal1 := &model.ContentCensorshipProposal{
		Permlink: permLink,
		Reason:   censorshipReason,
	}
	decideHr := int64(100)
	proposalID1, _ := proposalManager.AddProposal(ctx, user1, proposal1, decideHr)

	testCases := []struct {
		testName            string
		msg                 VoteProposalMsg
		wantRes             sdk.Result
		wantOK              bool
		wantOngoingProposal []types.ProposalKey
		wantProposal        model.Proposal
	}{
		{testName: "Must become a voter before voting",
			msg: VoteProposalMsg{
				Voter:      user2,
				ProposalID: proposalID1,
				Result:     true,
			},
			wantRes:             ErrGetVoter().Result(),
			wantOK:              true,
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal: &model.ContentCensorshipProposal{
				model.ProposalInfo{
					Creator:       user1,
					ProposalID:    proposalID1,
					AgreeVotes:    types.NewCoinFromInt64(0),
					DisagreeVotes: types.NewCoinFromInt64(0),
					Result:        types.ProposalNotPass,
					CreatedAt:     curTime,
					ExpiredAt:     curTime + decideHr*3600,
				}, permLink, censorshipReason},
		},
		{testName: "Vote on a non-exist proposal should fail",
			msg: VoteProposalMsg{
				Voter:      user1,
				ProposalID: types.ProposalKey(100),
				Result:     true,
			},
			wantRes:             ErrNotOngoingProposal().Result(),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal: &model.ContentCensorshipProposal{
				model.ProposalInfo{
					Creator:       user1,
					ProposalID:    proposalID1,
					AgreeVotes:    types.NewCoinFromInt64(0),
					DisagreeVotes: types.NewCoinFromInt64(0),
					Result:        types.ProposalNotPass,
					CreatedAt:     curTime,
					ExpiredAt:     curTime + decideHr*3600,
				}, permLink, censorshipReason},
		},
		{testName: "vote successfully",
			msg: VoteProposalMsg{
				Voter:      user1,
				ProposalID: proposalID1,
				Result:     true,
			},
			wantRes:             sdk.Result{},
			wantOK:              true,
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal: &model.ContentCensorshipProposal{
				model.ProposalInfo{
					Creator:       user1,
					ProposalID:    proposalID1,
					AgreeVotes:    c4600,
					DisagreeVotes: types.NewCoinFromInt64(0),
					Result:        types.ProposalNotPass,
					CreatedAt:     curTime,
					ExpiredAt:     curTime + decideHr*3600,
				}, permLink, censorshipReason},
		},
		{testName: "user can't double-vote",
			msg: VoteProposalMsg{
				Voter:      user1,
				ProposalID: proposalID1,
				Result:     false,
			},
			wantRes:             vote.ErrVoteExist().Result(),
			wantOK:              true,
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal: &model.ContentCensorshipProposal{
				model.ProposalInfo{
					Creator:       user1,
					ProposalID:    proposalID1,
					AgreeVotes:    c4600,
					DisagreeVotes: types.NewCoinFromInt64(0),
					Result:        types.ProposalNotPass,
					CreatedAt:     curTime,
					ExpiredAt:     curTime + decideHr*3600,
				}, permLink, censorshipReason},
		},
	}
	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		assert.Equal(t, tc.wantRes, result)

		if !tc.wantOK {
			continue
		}

		proposalList, _ := proposalManager.GetProposalList(ctx)
		assert.Equal(t, tc.wantOngoingProposal, proposalList.OngoingProposal)
		proposal, _ := proposalManager.storage.GetProposal(ctx, tc.msg.ProposalID)
		assert.Equal(t, tc.wantProposal, proposal)
	}
}
