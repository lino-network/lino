package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("vote")
)

func setup(t *testing.T) (sdk.Context, VoteStorage) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)
	vs := NewVoteStorage(TestKVStoreKey)
	return ctx, vs
}

func TestVoter(t *testing.T) {
	ctx, vs := setup(t)

	user := types.AccountKey("user")
	voter := Voter{
		Username:       user,
		Deposit:        types.NewCoin(1000),
		DelegatedPower: types.NewCoin(10000),
	}
	err := vs.SetVoter(ctx, user, &voter)
	assert.Nil(t, err)

	voterPtr, err := vs.GetVoter(ctx, user)
	assert.Nil(t, err)
	assert.Equal(t, voter, *voterPtr, "voter should be equal")

	vs.DeleteVoter(ctx, user)
	voterPtr, err = vs.GetVoter(ctx, user)
	assert.Nil(t, voterPtr)
	assert.Equal(t, ErrGetVoter(), err)
}

func TestVote(t *testing.T) {
	ctx, vs := setup(t)

	user1, user2, user3 :=
		types.AccountKey("user1"), types.AccountKey("user2"), types.AccountKey("user3")
	proposalID1, proposalID2 := types.ProposalKey("1"), types.ProposalKey("2")
	votingPower := types.NewCoin(1000)
	cases := []struct {
		isDelete    bool
		voter       types.AccountKey
		result      bool
		votingPower types.Coin
		proposalID  types.ProposalKey
		expectVotes []Vote
	}{
		{false, user1, true, votingPower, proposalID1, []Vote{Vote{user1, votingPower, true}}},
		{false, user2, true, votingPower, proposalID2, []Vote{Vote{user2, votingPower, true}}},
		{false, user2, false, votingPower, proposalID2, []Vote{Vote{user2, votingPower, false}}},
		{false, user3, true, votingPower, proposalID2, []Vote{Vote{user2, votingPower, false}, Vote{user3, votingPower, true}}},
		{true, user1, true, votingPower, proposalID1, nil},
		{true, user2, true, votingPower, proposalID2, []Vote{Vote{user3, votingPower, true}}},
		{false, user3, false, votingPower, proposalID2, []Vote{Vote{user3, votingPower, false}}},
		{false, user2, true, votingPower, proposalID2, []Vote{Vote{user2, votingPower, true}, Vote{user3, votingPower, false}}},
	}

	for _, cs := range cases {
		if cs.isDelete {
			vs.DeleteVote(ctx, cs.proposalID, cs.voter)
			votePtr, err := vs.GetVote(ctx, cs.proposalID, cs.voter)
			assert.Nil(t, votePtr)
			assert.Equal(t, ErrGetVote(), err)
		} else {
			vote := Vote{
				Voter:       cs.voter,
				Result:      cs.result,
				VotingPower: cs.votingPower,
			}
			err := vs.SetVote(ctx, cs.proposalID, cs.voter, &vote)
			assert.Nil(t, err)

			votePtr, err := vs.GetVote(ctx, cs.proposalID, cs.voter)
			assert.Nil(t, err)
			assert.Equal(t, vote, *votePtr, "vote should be equal")
		}
		allVotes, err := vs.GetAllVotes(ctx, cs.proposalID)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectVotes, allVotes)
	}
}

func TestDelegation(t *testing.T) {
	ctx, vs := setup(t)
	user1, user2, user3 :=
		types.AccountKey("user1"), types.AccountKey("user2"), types.AccountKey("user3")

	cases := []struct {
		delegator  types.AccountKey
		amount     types.Coin
		delegateTo types.AccountKey
	}{
		{user1, types.NewCoin(1), user2},
		{user1, types.NewCoin(100), user2},
		{user1, types.NewCoin(100), user1},
		{user3, types.NewCoin(100), user1},
		{user3, types.NewCoin(100), user2},
	}

	for _, cs := range cases {
		err := vs.SetDelegation(ctx, cs.delegateTo, cs.delegator,
			&Delegation{cs.delegator, cs.amount})
		assert.Nil(t, err)
		delegationPtr, err := vs.GetDelegation(ctx, cs.delegateTo, cs.delegator)
		assert.Nil(t, err)
		assert.Equal(t, Delegation{cs.delegator, cs.amount}, *delegationPtr)
		err = vs.DeleteDelegation(ctx, cs.delegateTo, cs.delegator)
		assert.Nil(t, err)
		delegationPtr, err = vs.GetDelegation(ctx, cs.delegateTo, cs.delegator)
		assert.Nil(t, delegationPtr)
		assert.Equal(t, ErrGetDelegation(), err)
	}

	getAllDelegatorsCases := []struct {
		delegator        types.AccountKey
		amount           types.Coin
		delegateTo       types.AccountKey
		expectDelegators []types.AccountKey
	}{
		{user1, types.NewCoin(1), user2, []types.AccountKey{user1}},
		{user1, types.NewCoin(100), user2, []types.AccountKey{user1}},
		{user1, types.NewCoin(1), user3, []types.AccountKey{user1}},
		{user2, types.NewCoin(1), user1, []types.AccountKey{user2}},
		{user3, types.NewCoin(1), user1, []types.AccountKey{user2, user3}},
	}

	for _, cs := range getAllDelegatorsCases {
		err := vs.SetDelegation(ctx, cs.delegateTo, cs.delegator,
			&Delegation{cs.delegator, cs.amount})
		assert.Nil(t, err)
		delegators, err := vs.GetAllDelegators(ctx, cs.delegateTo)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectDelegators, delegators)
	}
}
