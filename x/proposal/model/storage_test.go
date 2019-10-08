package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("proposal")
)

func setup(t *testing.T) (sdk.Context, ProposalStorage) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	ps := NewProposalStorage(TestKVStoreKey)
	err = ps.InitGenesis(ctx)
	assert.Nil(t, err)
	return ctx, ps
}

func TestProposalList(t *testing.T) {
	ctx, ps := setup(t)

	proposalList, err := ps.GetOngoingProposalList(ctx)
	assert.Nil(t, err)

	p1 := ChangeParamProposal{
		ProposalInfo: ProposalInfo{
			Creator:       types.AccountKey("user"),
			ProposalID:    types.ProposalKey("123"),
			AgreeVotes:    types.NewCoinFromInt64(0),
			DisagreeVotes: types.NewCoinFromInt64(0),
		},
		Param: param.GlobalAllocationParam{
			GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
			ContentCreatorAllocation: sdk.NewDec(0),
			DeveloperAllocation:      sdk.NewDec(0),
			ValidatorAllocation:      sdk.NewDec(0),
		},
	}

	p2 := p1
	p2.ProposalInfo.ProposalID = types.ProposalKey("321")

	err = ps.SetOngoingProposal(ctx, types.ProposalKey("123"), &p1)
	if err != nil {
		panic(err)
	}
	err = ps.SetOngoingProposal(ctx, types.ProposalKey("321"), &p2)
	if err != nil {
		panic(err)
	}

	proposalList =
		append(proposalList, &p1)
	proposalList =
		append(proposalList, &p2)

	ongoingProposalList, err := ps.GetOngoingProposalList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, proposalList, ongoingProposalList)
}

func TestProposal(t *testing.T) {
	ctx, ps := setup(t)
	user := types.AccountKey("user")
	proposalID := types.ProposalKey("123")
	res := types.ProposalPass
	curTime := ctx.BlockHeader().Time.Unix()

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
					GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
					ContentCreatorAllocation: sdk.NewDec(0),
					DeveloperAllocation:      sdk.NewDec(0),
					ValidatorAllocation:      sdk.NewDec(0),
				},
			},
		},
	}

	for _, tc := range testCases {
		err := ps.SetOngoingProposal(ctx, proposalID, &tc.changeParamProposal)
		if err != nil {
			t.Errorf("%s: failed to set proposal, get err %v", tc.testName, err)
		}

		err = ps.SetExpiredProposal(ctx, proposalID, &tc.changeParamProposal)
		if err != nil {
			t.Errorf("%s: failed to set proposal, get err %v", tc.testName, err)
		}

		proposal, err := ps.GetOngoingProposal(ctx, proposalID)
		if err != nil {
			t.Errorf("%s: failed to get proposal, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, &tc.changeParamProposal, proposal.(*ChangeParamProposal)) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, proposal.(*ChangeParamProposal), &tc.changeParamProposal)
		}

		proposal, err = ps.GetExpiredProposal(ctx, proposalID)
		if err != nil {
			t.Errorf("%s: failed to get proposal, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, &tc.changeParamProposal, proposal.(*ChangeParamProposal)) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, proposal.(*ChangeParamProposal), &tc.changeParamProposal)
		}

		err = ps.DeleteOngoingProposal(ctx, proposalID)
		if err != nil {
			t.Errorf("%s: failed to delete proposal, get err %v", tc.testName, err)
		}

		err = ps.DeleteExpiredProposal(ctx, proposalID)
		if err != nil {
			t.Errorf("%s: failed to delete proposal, get err %v", tc.testName, err)
		}

		_, err = ps.GetOngoingProposal(ctx, proposalID)
		if err == nil {
			t.Errorf("%s: failed to get proposal after deletion, get err %v", tc.testName, err)
		}
		if !assert.Equal(t, ErrProposalNotFound(), err) {
			t.Errorf("%s: diff err, got %v, want %v", tc.testName, err, ErrProposalNotFound())
		}

		_, err = ps.GetExpiredProposal(ctx, proposalID)
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
