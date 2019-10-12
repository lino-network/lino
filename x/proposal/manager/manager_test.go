package manager

// import (
// 	"testing"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/lino-network/lino/types"
// 	"github.com/lino-network/lino/x/proposal/model"
// 	"github.com/stretchr/testify/assert"
// )

// func TestUpdateProposalVotingStatus(t *testing.T) {
// 	ctx, _, pm, _, _, _, _ := setupTest(t, 0)
// 	permlink := types.Permlink("permlink")
// 	user1 := types.AccountKey("user1")
// 	censorshipReason := "reason"
// 	proposal1 := &model.ContentCensorshipProposal{
// 		Permlink: permlink,
// 		Reason:   censorshipReason,
// 	}

// 	err := pm.InitGenesis(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	curTime := ctx.BlockHeader().Time.Unix()
// 	decideSec := int64(100)
// 	proposalID1, _ := pm.AddProposal(ctx, user1, proposal1, decideSec)

// 	testCases := []struct {
// 		testName     string
// 		proposalID   types.ProposalKey
// 		voter        types.AccountKey
// 		voteResult   bool
// 		votingPower  types.Coin
// 		wantProposal model.Proposal
// 	}{
// 		{
// 			testName:    "agree vote",
// 			proposalID:  proposalID1,
// 			voter:       user1,
// 			voteResult:  true,
// 			votingPower: types.NewCoinFromInt64(1),
// 			wantProposal: &model.ContentCensorshipProposal{
// 				ProposalInfo: model.ProposalInfo{
// 					Creator:       user1,
// 					ProposalID:    proposalID1,
// 					AgreeVotes:    types.NewCoinFromInt64(1),
// 					DisagreeVotes: types.NewCoinFromInt64(0),
// 					CreatedAt:     curTime,
// 					ExpiredAt:     curTime + decideSec,
// 				},
// 				Permlink: permlink,
// 				Reason:   censorshipReason},
// 		},
// 		{
// 			testName:    "one more agree vote",
// 			proposalID:  proposalID1,
// 			voter:       user1,
// 			voteResult:  true,
// 			votingPower: types.NewCoinFromInt64(2),
// 			wantProposal: &model.ContentCensorshipProposal{
// 				ProposalInfo: model.ProposalInfo{
// 					Creator:       user1,
// 					ProposalID:    proposalID1,
// 					AgreeVotes:    types.NewCoinFromInt64(3),
// 					DisagreeVotes: types.NewCoinFromInt64(0),
// 					CreatedAt:     curTime,
// 					ExpiredAt:     curTime + decideSec,
// 				},
// 				Permlink: permlink,
// 				Reason:   censorshipReason},
// 		},
// 		{
// 			testName:    "one disagree vote",
// 			proposalID:  proposalID1,
// 			voter:       user1,
// 			voteResult:  false,
// 			votingPower: types.NewCoinFromInt64(5),
// 			wantProposal: &model.ContentCensorshipProposal{
// 				ProposalInfo: model.ProposalInfo{
// 					Creator:       user1,
// 					ProposalID:    proposalID1,
// 					AgreeVotes:    types.NewCoinFromInt64(3),
// 					DisagreeVotes: types.NewCoinFromInt64(5),
// 					CreatedAt:     curTime,
// 					ExpiredAt:     curTime + decideSec,
// 				},
// 				Permlink: permlink,
// 				Reason:   censorshipReason},
// 		},
// 	}
// 	for _, tc := range testCases {
// 		err := pm.UpdateProposalVotingStatus(ctx, tc.proposalID, tc.voter, tc.voteResult, tc.votingPower)
// 		if err != nil {
// 			t.Errorf("%s: failed to update proposal voting status, got err %v", tc.testName, err)
// 		}

// 		proposal, err := pm.storage.GetOngoingProposal(ctx, tc.proposalID)
// 		if err != nil {
// 			t.Errorf("%s: failed to get proposal, got err %v", tc.testName, err)
// 		}
// 		if !assert.Equal(t, tc.wantProposal, proposal) {
// 			t.Errorf("%s: diff result, got %v, want %v", tc.testName, proposal, tc.wantProposal)
// 		}
// 	}
// }

// func TestUpdateProposalPassStatus(t *testing.T) {
// 	ctx, _, pm, _, _, _, _ := setupTest(t, 100000000)
// 	permlink := types.Permlink("permlink")
// 	user1 := types.AccountKey("user1")
// 	censorshipReason := "reason"
// 	proposal1 := &model.ContentCensorshipProposal{
// 		Permlink: permlink,
// 		Reason:   censorshipReason,
// 	}

// 	proposal2 := &model.ContentCensorshipProposal{
// 		Permlink: permlink,
// 		Reason:   censorshipReason,
// 	}

// 	proposal3 := &model.ContentCensorshipProposal{
// 		Permlink: permlink,
// 		Reason:   censorshipReason,
// 	}
// 	err := pm.InitGenesis(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	curTime := ctx.BlockHeader().Time.Unix()
// 	proposalParam, _ := pm.paramHolder.GetProposalParam(ctx)
// 	decideSec := proposalParam.ContentCensorshipDecideSec

// 	proposalID1, _ := pm.AddProposal(ctx, user1, proposal1, decideSec)
// 	proposalID2, _ := pm.AddProposal(ctx, user1, proposal2, decideSec)
// 	proposalID3, _ := pm.AddProposal(ctx, user1, proposal3, decideSec)

// 	testCases := []struct {
// 		testName        string
// 		agreeVotes      types.Coin
// 		disagreeVotes   types.Coin
// 		proposalType    types.ProposalType
// 		proposalID      types.ProposalKey
// 		wantProposalRes types.ProposalResult
// 		wantProposal    model.Proposal
// 	}{
// 		{
// 			testName:        "test passed proposal has historical data",
// 			agreeVotes:      proposalParam.ContentCensorshipPassVotes,
// 			disagreeVotes:   proposalParam.ContentCensorshipPassVotes,
// 			proposalType:    types.ContentCensorship,
// 			proposalID:      proposalID1,
// 			wantProposalRes: types.ProposalNotPass,
// 			wantProposal: &model.ContentCensorshipProposal{
// 				ProposalInfo: model.ProposalInfo{
// 					Creator:       user1,
// 					ProposalID:    proposalID1,
// 					AgreeVotes:    proposalParam.ContentCensorshipPassVotes,
// 					DisagreeVotes: proposalParam.ContentCensorshipPassVotes,
// 					Result:        types.ProposalNotPass,
// 					CreatedAt:     curTime,
// 					ExpiredAt:     curTime + decideSec,
// 				},
// 				Permlink: permlink,
// 				Reason:   censorshipReason},
// 		},

// 		{
// 			testName:        "test votes don't meet min requirement ",
// 			agreeVotes:      proposalParam.ContentCensorshipPassVotes.Minus(types.NewCoinFromInt64(10)),
// 			disagreeVotes:   types.NewCoinFromInt64(0),
// 			proposalType:    types.ContentCensorship,
// 			proposalID:      proposalID2,
// 			wantProposalRes: types.ProposalNotPass,
// 			wantProposal: &model.ContentCensorshipProposal{
// 				ProposalInfo: model.ProposalInfo{
// 					Creator:       user1,
// 					ProposalID:    proposalID2,
// 					AgreeVotes:    proposalParam.ContentCensorshipPassVotes.Minus(types.NewCoinFromInt64(10)),
// 					DisagreeVotes: types.NewCoinFromInt64(0),
// 					Result:        types.ProposalNotPass,
// 					CreatedAt:     curTime,
// 					ExpiredAt:     curTime + decideSec,
// 				},
// 				Permlink: permlink,
// 				Reason:   censorshipReason},
// 		},

// 		{
// 			testName:        "test votes ratio doesn't meet requirement ",
// 			agreeVotes:      proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoinFromInt64(10)),
// 			disagreeVotes:   proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoinFromInt64(11)),
// 			proposalType:    types.ContentCensorship,
// 			proposalID:      proposalID3,
// 			wantProposalRes: types.ProposalNotPass,
// 			wantProposal: &model.ContentCensorshipProposal{
// 				ProposalInfo: model.ProposalInfo{
// 					Creator:       user1,
// 					ProposalID:    proposalID3,
// 					AgreeVotes:    proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoinFromInt64(10)),
// 					DisagreeVotes: proposalParam.ContentCensorshipPassVotes.Plus(types.NewCoinFromInt64(11)),
// 					Result:        types.ProposalNotPass,
// 					CreatedAt:     curTime,
// 					ExpiredAt:     curTime + decideSec,
// 				},
// 				Permlink: permlink,
// 				Reason:   censorshipReason},
// 		},
// 	}
// 	for _, tc := range testCases {
// 		err := addProposalInfo(ctx, pm, tc.proposalID, tc.agreeVotes, tc.disagreeVotes)
// 		if err != nil {
// 			t.Errorf("%s: failed to add proposal info, got err %v", tc.testName, err)
// 		}

// 		res, err := pm.UpdateProposalPassStatus(ctx, tc.proposalType, tc.proposalID)
// 		if err != nil {
// 			t.Errorf("%s: failed to update proposal pass status, got err %v", tc.testName, err)
// 		}
// 		if res != tc.wantProposalRes {
// 			t.Errorf("%s: test failed, got %v, want %v", tc.testName, res, tc.wantProposalRes)
// 			return
// 		}

// 		ongoingProposal, err := pm.storage.GetOngoingProposal(ctx, tc.proposalID)
// 		assert.NotNil(t, err)
// 		if !assert.Equal(t, nil, ongoingProposal) {
// 			t.Errorf("%s: didn't remove ongoing proposal", tc.testName)
// 		}
// 		proposal, err := pm.storage.GetExpiredProposal(ctx, tc.proposalID)
// 		if err != nil {
// 			t.Errorf("%s: failed to get proposal, got err %v", tc.testName, err)
// 		}
// 		if !assert.Equal(t, tc.wantProposal, proposal) {
// 			t.Errorf("%s: diff result, got %v, want %v", tc.testName, proposal, tc.wantProposal)
// 		}
// 	}
// }

// func TestGetProposalPassParam(t *testing.T) {
// 	ctx, _, pm, _, _, _, _ := setupTest(t, 0)

// 	proposalParam, _ := pm.paramHolder.GetProposalParam(ctx)
// 	testCases := []struct {
// 		testName      string
// 		proposalType  types.ProposalType
// 		wantError     sdk.Error
// 		wantPassRatio sdk.Dec
// 		wantPassVotes types.Coin
// 	}{
// 		{
// 			testName:      "test pass param for changeParamProposal",
// 			proposalType:  types.ChangeParam,
// 			wantError:     nil,
// 			wantPassRatio: proposalParam.ChangeParamPassRatio,
// 			wantPassVotes: proposalParam.ChangeParamPassVotes,
// 		},

// 		{
// 			testName:      "test pass param for contenCensorshipProposal",
// 			proposalType:  types.ContentCensorship,
// 			wantError:     nil,
// 			wantPassRatio: proposalParam.ContentCensorshipPassRatio,
// 			wantPassVotes: proposalParam.ContentCensorshipPassVotes,
// 		},

// 		{
// 			testName:      "test pass param for protocolUpgradeProposal",
// 			proposalType:  types.ProtocolUpgrade,
// 			wantError:     nil,
// 			wantPassRatio: proposalParam.ProtocolUpgradePassRatio,
// 			wantPassVotes: proposalParam.ProtocolUpgradePassVotes,
// 		},

// 		{
// 			testName:      "test wrong proposal type",
// 			proposalType:  23,
// 			wantError:     ErrIncorrectProposalType(),
// 			wantPassRatio: proposalParam.ProtocolUpgradePassRatio,
// 			wantPassVotes: proposalParam.ProtocolUpgradePassVotes,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		ratio, votes, err := pm.GetProposalPassParam(ctx, tc.proposalType)
// 		if !assert.Equal(t, tc.wantError, err) {
// 			t.Errorf("%s: diff err, got %v, want %v", tc.testName, err, tc.wantError)
// 		}

// 		if tc.wantError != nil {
// 			continue
// 		}
// 		if !assert.Equal(t, tc.wantPassRatio, ratio) {
// 			t.Errorf("%s: diff ratio, got %v, want %v", tc.testName, ratio, tc.wantPassRatio)
// 		}
// 		if !assert.Equal(t, tc.wantPassVotes, votes) {
// 			t.Errorf("%s: diff pass votes, got %v, want %v", tc.testName, votes, tc.wantPassVotes)
// 		}
// 	}

// }
