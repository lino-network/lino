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
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("proposal")
)

func setup(t *testing.T) (sdk.Context, ProposalStorage) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)
	ps := NewProposalStorage(TestKVStoreKey)
	err := ps.InitGenesis(ctx)
	assert.Nil(t, err)
	return ctx, ps
}

func TestProposalList(t *testing.T) {
	ctx, vs := setup(t)

	proposalList, err := vs.GetProposalList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, ProposalList{}, *proposalList)
	proposalList.OngoingProposal =
		append(proposalList.OngoingProposal, types.ProposalKey("test1"))
	proposalList.OngoingProposal =
		append(proposalList.OngoingProposal, types.ProposalKey("test2"))
	proposalList.PastProposal = append(proposalList.PastProposal, types.ProposalKey("test3"))

	err = vs.SetProposalList(ctx, proposalList)
	assert.Nil(t, err)

	proposalListPtr, err := vs.GetProposalList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, proposalList, proposalListPtr)
}

func TestProposal(t *testing.T) {
	ctx, vs := setup(t)
	user := types.AccountKey("user")
	proposalID := types.ProposalKey("123")
	res := types.ProposalPass

	cases := []struct {
		ChangeParamProposal
	}{
		{ChangeParamProposal{
			ProposalInfo{user, proposalID, types.NewCoin(0), types.NewCoin(0), res},
			param.GlobalAllocationParam{sdk.NewRat(0), sdk.NewRat(0), sdk.NewRat(0),
				sdk.NewRat(0)}}},
	}

	for _, cs := range cases {
		err := vs.SetProposal(ctx, proposalID, &cs.ChangeParamProposal)
		assert.Nil(t, err)
		proposal, err := vs.GetProposal(ctx, proposalID)
		assert.Nil(t, err)
		assert.Equal(t, &cs.ChangeParamProposal, proposal.(*ChangeParamProposal))
		err = vs.DeleteProposal(ctx, proposalID)
		assert.Nil(t, err)
		proposal, err = vs.GetProposal(ctx, proposalID)
		assert.Nil(t, proposal)
		assert.Equal(t, ErrGetProposal(), err)
	}
}
