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

	cases := []struct {
		ChangeParamProposal
	}{
		{ChangeParamProposal{
			ProposalInfo{user, proposalID, types.NewCoinFromInt64(0), types.NewCoinFromInt64(0), res, curTime, curTime + 100},
			param.GlobalAllocationParam{sdk.NewRat(0), sdk.NewRat(0), sdk.NewRat(0),
				sdk.NewRat(0)}}},
	}

	for _, cs := range cases {
		err := ps.SetProposal(ctx, proposalID, &cs.ChangeParamProposal)
		assert.Nil(t, err)
		proposal, err := ps.GetProposal(ctx, proposalID)
		assert.Nil(t, err)
		assert.Equal(t, &cs.ChangeParamProposal, proposal.(*ChangeParamProposal))
		err = ps.DeleteProposal(ctx, proposalID)
		assert.Nil(t, err)
		proposal, err = ps.GetProposal(ctx, proposalID)
		assert.Nil(t, proposal)
		assert.Equal(t, ErrProposalNotFound(), err)
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
