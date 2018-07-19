package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("proposal")
)

func setup(t *testing.T) (sdk.Context, ProposalStorage) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
	ps := NewProposalStorage(TestKVStoreKey)
	err := ps.InitGenesis(ctx)
	assert.Nil(t, err)
	return ctx, ps
}

func TestProposalList(t *testing.T) {
	ctx, ps := setup(t)

	proposalList, err := ps.GetProposalList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, ProposalList{}, *proposalList)
	proposalList.OngoingProposal =
		append(proposalList.OngoingProposal, types.ProposalKey("test1"))
	proposalList.OngoingProposal =
		append(proposalList.OngoingProposal, types.ProposalKey("test2"))
	proposalList.PastProposal = append(proposalList.PastProposal, types.ProposalKey("test3"))

	err = ps.SetProposalList(ctx, proposalList)
	assert.Nil(t, err)

	proposalListPtr, err := ps.GetProposalList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, proposalList, proposalListPtr)
}

func TestProposal(t *testing.T) {
	ctx, ps := setup(t)
	user := types.AccountKey("user")
	proposalID := types.ProposalKey("123")
	res := types.ProposalPass
	curTime := ctx.BlockHeader().Time

	testCases := []struct {
		testName            string
		changeParamProposal ChangeParamProposal
	}{
		{
			testName: "change param proposal",
			changeParamProposal: ChangeParamProposal{
				ProposalInfo: ProposalInfo{
					Creator:       user,
					ProposalID:    proposalID,
					AgreeVotes:    types.NewCoinFromInt64(0),
					DisagreeVotes: types.NewCoinFromInt64(0),
					Result:        res,
					CreatedAt:     curTime,
					ExpiredAt:     curTime + 100,
				},
				Param: param.GlobalAllocationParam{
					InfraAllocation:          sdk.NewRat(0),
					ContentCreatorAllocation: sdk.NewRat(0),
					DeveloperAllocation:      sdk.NewRat(0),
					ValidatorAllocation:      sdk.NewRat(0),
				},
			},
		},
	}

	for _, tc := range testCases {
		err := ps.SetProposal(ctx, proposalID, &tc.changeParamProposal)
		if err != nil {
			t.Errorf("%s: failed to set proposal, get err %v", tc.testName, err)
		}

		proposal, err := ps.GetProposal(ctx, proposalID)
		if err != nil {
			t.Errorf("%s: failed to get proposal, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, &tc.changeParamProposal, proposal.(*ChangeParamProposal)) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, proposal.(*ChangeParamProposal), &tc.changeParamProposal)
		}

		err = ps.DeleteProposal(ctx, proposalID)
		if err != nil {
			t.Errorf("%s: failed to delete proposal, get err %v", tc.testName, err)
		}

		proposal, err = ps.GetProposal(ctx, proposalID)
		if err == nil {
			t.Errorf("%s: failed to get proposal after deletion, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, ErrProposalNotFound(), err) {
			t.Errorf("%s: diff err, got %v, want %v", tc.testName, err, ErrProposalNotFound())
		}
	}
}

func TestNextProposalID(t *testing.T) {
	ctx, ps := setup(t)

	id, err := ps.GetNextProposalID(ctx)
	assert.Nil(t, err)
	assert.Equal(t, NextProposalID{1}, *id)

	id.NextProposalID = 2
	err = ps.SetNextProposalID(ctx, id)
	assert.Nil(t, err)

	nextProposalID, err := ps.GetNextProposalID(ctx)
	assert.Nil(t, err)
	assert.Equal(t, nextProposalID, id)
}
