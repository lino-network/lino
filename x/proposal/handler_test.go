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

	curTime := ctx.BlockHeader().Time.Unix()
	proposalParam, _ := proposalManager.paramHolder.GetProposalParam(ctx)

	proposal1 := &model.ChangeParamProposal{model.ProposalInfo{
		Creator:       user1,
		ProposalID:    proposalID1,
		AgreeVotes:    types.NewCoinFromInt64(0),
		DisagreeVotes: types.NewCoinFromInt64(0),
		Result:        types.ProposalNotPass,
		CreatedAt:     curTime,
		ExpiredAt:     curTime + proposalParam.ChangeParamDecideSec,
	}, allocation, ""}

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
		{
			testName: "user1 creates change param msg successfully",
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

		{
			testName: "user2 doesn't have enough money to create proposal",
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
		if !assert.Equal(t, tc.wantRes, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.wantRes)
		}

		if !tc.wantOK {
			continue
		}

		creatorBalance, _ := am.GetSavingFromBank(ctx, tc.msg.GetCreator())
		if !creatorBalance.IsEqual(tc.wantCreatorBalance) {
			t.Errorf("%s: diff bank balance: got %v, want %v", tc.testName, creatorBalance, tc.wantCreatorBalance)
		}

		proposalList, err := proposalManager.GetProposalList(ctx)
		if err != nil {
			t.Errorf("%s: failed to get proposal list, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.wantOngoingProposal, proposalList.OngoingProposal) {
			t.Errorf("%s: diff ongoing proposal, got %v, want %v", tc.testName, proposalList.OngoingProposal, tc.wantOngoingProposal)
		}

		proposal, err := proposalManager.storage.GetProposal(ctx, tc.proposalID)
		if err != nil {
			t.Errorf("%s: failed to get proposal, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.wantProposal, proposal) {
			t.Errorf("%s: diff proposal, got %v, want %v", tc.testName, proposal, tc.wantProposal)
		}
	}
}

func TestContentCensorshipProposal(t *testing.T) {
	ctx, am, proposalManager, postManager, vm, _, gm := setupTest(t, 0)
	handler := NewHandler(am, proposalManager, postManager, gm, vm)
	curTime := ctx.BlockHeader().Time.Unix()
	proposalParam, _ := proposalManager.paramHolder.GetProposalParam(ctx)

	proposalManager.InitGenesis(ctx)

	proposalID1 := types.ProposalKey(strconv.FormatInt(int64(1), 10))

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
		ExpiredAt:     curTime + proposalParam.ChangeParamDecideSec,
	}, types.GetPermlink(user1, postID1), censorshipReason}

	testCases := []struct {
		testName            string
		creator             types.AccountKey
		permlink            types.Permlink
		proposalID          types.ProposalKey
		wantOK              bool
		wantRes             sdk.Result
		wantCreatorBalance  types.Coin
		wantOngoingProposal []types.ProposalKey
		wantProposal        model.Proposal
	}{
		{
			testName:            "user2 censorship user1's post successfully",
			creator:             user2,
			permlink:            types.GetPermlink(user1, postID1),
			proposalID:          proposalID1,
			wantOK:              true,
			wantRes:             sdk.Result{},
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{
			testName:            "target post is not exist",
			creator:             user2,
			permlink:            types.GetPermlink(user1, "invalid"),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             ErrPostNotFound().Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{
			testName:            "target post is deleted",
			creator:             user1,
			permlink:            types.GetPermlink(user2, postID2),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             ErrCensorshipPostIsDeleted(types.GetPermlink(user2, postID2)).Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{
			testName:            "proposal is invalid",
			creator:             "invalid",
			permlink:            types.GetPermlink(user1, postID1),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             ErrAccountNotFound().Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
		{
			testName:            "user3 doesn't have enough money to create proposal",
			creator:             user3,
			permlink:            types.GetPermlink(user1, postID1),
			proposalID:          proposalID1,
			wantOK:              false,
			wantRes:             acc.ErrAccountSavingCoinNotEnough().Result(),
			wantCreatorBalance:  c4600.Minus(proposalParam.ContentCensorshipMinDeposit),
			wantOngoingProposal: []types.ProposalKey{proposalID1},
			wantProposal:        proposal1,
		},
	}
	for _, tc := range testCases {
		msg := NewDeletePostContentMsg(string(tc.creator), tc.permlink, censorshipReason)
		result := handler(ctx, msg)
		if !assert.Equal(t, tc.wantRes, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.wantRes)
		}

		if !tc.wantOK {
			continue
		}

		creatorBalance, _ := am.GetSavingFromBank(ctx, tc.creator)
		if !creatorBalance.IsEqual(tc.wantCreatorBalance) {
			t.Errorf("%s: diff bank balance: got %v, want %v",
				tc.testName, creatorBalance, tc.wantCreatorBalance)
		}

		proposalList, err := proposalManager.GetProposalList(ctx)
		if err != nil {
			t.Errorf("%s: failed to get proposal list, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.wantOngoingProposal, proposalList.OngoingProposal) {
			t.Errorf("%s: diff ongoing proposal, got %v, want %v", tc.testName, proposalList.OngoingProposal, tc.wantOngoingProposal)
		}

		proposal, err := proposalManager.storage.GetProposal(ctx, tc.proposalID)
		if err != nil {
			t.Errorf("%s: failed to get proposal, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.wantProposal, proposal) {
			t.Errorf("%s: diff proposal, got %v, want %v", tc.testName, proposal, tc.wantProposal)
		}

		permlink, err := proposalManager.GetPermlink(ctx, tc.proposalID)
		if err != nil {
			t.Errorf("%s: failed to get permlink, get err %v", tc.testName, err)
		}
		if permlink != tc.permlink {
			t.Errorf("%s: diff permlink, got %v, want %v", tc.testName, permlink, tc.permlink)
		}
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
		{
			testName:               "return coin to user",
			times:                  10,
			interval:               2,
			returnedCoin:           types.NewCoinFromInt64(100),
			expectedFrozenListLen:  1,
			expectedFrozenMoney:    types.NewCoinFromInt64(100),
			expectedFrozenTimes:    10,
			expectedFrozenInterval: 2,
		},
		{
			testName:               "return coin to user multiple times",
			times:                  100000,
			interval:               20000,
			returnedCoin:           types.NewCoinFromInt64(100000),
			expectedFrozenListLen:  2,
			expectedFrozenMoney:    types.NewCoinFromInt64(100000),
			expectedFrozenTimes:    100000,
			expectedFrozenInterval: 20000,
		},
	}

	for _, tc := range testCases {
		err := returnCoinTo(
			ctx, "user", gm, am, tc.times, tc.interval, tc.returnedCoin)
		if err != nil {
			t.Errorf("%s: failed to return coin, got err %v", tc.testName, err)
		}

		lst, err := am.GetFrozenMoneyList(ctx, user)
		if err != nil {
			t.Errorf("%s: failed to get frozen money list, got err %v", tc.testName, err)
		}
		if len(lst) != tc.expectedFrozenListLen {
			t.Errorf("%s: diff list len, got %v, want %v", tc.testName, len(lst), tc.expectedFrozenListLen)
		}
		if !lst[len(lst)-1].Amount.IsEqual(tc.expectedFrozenMoney) {
			t.Errorf("%s: diff amount, got %v, want %v", tc.testName, lst[len(lst)-1].Amount, tc.expectedFrozenMoney)
		}
		if lst[len(lst)-1].Times != tc.expectedFrozenTimes {
			t.Errorf("%s: diff times, got %v, want %v", tc.testName, lst[len(lst)-1].Times, tc.expectedFrozenTimes)
		}
		if lst[len(lst)-1].Interval != tc.expectedFrozenInterval {
			t.Errorf("%s: diff interval, got %v, want %v", tc.testName, lst[len(lst)-1].Interval, tc.expectedFrozenInterval)
		}
	}
}

func TestVoteProposalBasic(t *testing.T) {
	ctx, am, proposalManager, postManager, vm, _, gm := setupTest(t, 0)
	handler := NewHandler(am, proposalManager, postManager, gm, vm)
	curTime := ctx.BlockHeader().Time.Unix()
	proposalManager.InitGenesis(ctx)

	user2 := types.AccountKey("user2")

	// create voter
	user1 := types.AccountKey("user1")
	createTestAccount(ctx, am, "user1", c4600)
	_ = vm.AddVoter(ctx, user1, c4600)

	// create proposal
	permlink := types.Permlink("postlink")
	censorshipReason := "reason"
	proposal1 := &model.ContentCensorshipProposal{
		Permlink: permlink,
		Reason:   censorshipReason,
	}
	decideSec := int64(100)
	proposalID1, _ := proposalManager.AddProposal(ctx, user1, proposal1, decideSec)

	testCases := []struct {
		testName            string
		msg                 VoteProposalMsg
		wantRes             sdk.Result
		wantOK              bool
		wantOngoingProposal []types.ProposalKey
		wantProposal        model.Proposal
	}{
		{
			testName: "Must become a voter before voting",
			msg: VoteProposalMsg{
				Voter:      user2,
				ProposalID: proposalID1,
				Result:     true,
			},
			wantRes:             ErrVoterNotFound().Result(),
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
					ExpiredAt:     curTime + decideSec,
				}, permlink, censorshipReason},
		},
		{
			testName: "Vote on a non-exist proposal should fail",
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
					ExpiredAt:     curTime + decideSec,
				}, permlink, censorshipReason},
		},
		{
			testName: "vote successfully",
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
					ExpiredAt:     curTime + decideSec,
				}, permlink, censorshipReason},
		},
		{
			testName: "user can't double-vote",
			msg: VoteProposalMsg{
				Voter:      user1,
				ProposalID: proposalID1,
				Result:     false,
			},
			wantRes:             vote.ErrVoteAlreadyExist().Result(),
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
					ExpiredAt:     curTime + decideSec,
				}, permlink, censorshipReason},
		},
	}
	for _, tc := range testCases {
		result := handler(ctx, tc.msg)
		if !assert.Equal(t, tc.wantRes, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.wantRes)
		}

		if !tc.wantOK {
			continue
		}

		proposalList, err := proposalManager.GetProposalList(ctx)
		if err != nil {
			t.Errorf("%s: failed to get proposal list, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.wantOngoingProposal, proposalList.OngoingProposal) {
			t.Errorf("%s: diff ongoing proposal, got %v, want %v", tc.testName, proposalList.OngoingProposal, tc.wantOngoingProposal)
		}

		proposal, err := proposalManager.storage.GetProposal(ctx, tc.msg.ProposalID)
		if err != nil {
			t.Errorf("%s: failed to get proposal, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.wantProposal, proposal) {
			t.Errorf("%s: diff proposal, got %v, want %v", tc.testName, proposal, tc.wantProposal)
		}
	}
}
