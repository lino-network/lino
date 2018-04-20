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
	err := vs.InitGenesis(ctx)
	assert.Nil(t, err)
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

	cases := []struct {
		isDelete    bool
		voter       types.AccountKey
		result      bool
		proposalID  types.ProposalKey
		expectVotes []Vote
	}{
		{false, user1, true, proposalID1, []Vote{Vote{user1, true}}},
		{false, user2, true, proposalID2, []Vote{Vote{user2, true}}},
		{false, user2, false, proposalID2, []Vote{Vote{user2, false}}},
		{false, user3, true, proposalID2, []Vote{Vote{user2, false}, Vote{user3, true}}},
		{true, user1, true, proposalID1, nil},
		{true, user2, true, proposalID2, []Vote{Vote{user3, true}}},
		{false, user3, false, proposalID2, []Vote{Vote{user3, false}}},
		{false, user2, true, proposalID2, []Vote{Vote{user2, true}, Vote{user3, false}}},
	}

	for _, cs := range cases {
		if cs.isDelete {
			vs.DeleteVote(ctx, cs.proposalID, cs.voter)
			votePtr, err := vs.GetVote(ctx, cs.proposalID, cs.voter)
			assert.Nil(t, votePtr)
			assert.Equal(t, ErrGetVote(), err)
		} else {
			vote := Vote{
				Voter:  cs.voter,
				Result: cs.result,
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

func TestProposalList(t *testing.T) {
	ctx, vs := setup(t)

	proposalList, err := vs.GetProposalList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, ProposalList{[]types.ProposalKey{}, []types.ProposalKey{}}, *proposalList)
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

func TestPenaltyList(t *testing.T) {
	ctx, vs := setup(t)
	lst, err := vs.GetValidatorReferenceList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, ValidatorReferenceList{[]types.AccountKey{},
		[]types.AccountKey{}, []types.AccountKey{}}, *lst)
	lst.PenaltyValidators =
		append(lst.PenaltyValidators, types.AccountKey("test1"))

	err = vs.SetValidatorReferenceList(ctx, lst)
	assert.Nil(t, err)

	lstPtr, err := vs.GetValidatorReferenceList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, lst, lstPtr)
}

func TestProposal(t *testing.T) {
	ctx, vs := setup(t)
	user := types.AccountKey("user")
	proposalID := types.ProposalKey("123")

	cases := []struct {
		ChangeParameterProposal
	}{
		{ChangeParameterProposal{
			Proposal{user, proposalID, types.NewCoin(0), types.NewCoin(0)},
			ChangeParameterDescription{sdk.NewRat(0), sdk.NewRat(0), sdk.NewRat(0),
				sdk.NewRat(0), sdk.NewRat(0), sdk.NewRat(0)}}},
	}

	for _, proposal := range cases {
		err := vs.SetProposal(ctx, proposalID, &proposal.ChangeParameterProposal)
		assert.Nil(t, err)
		proposlPtr, err := vs.GetProposal(ctx, proposalID)
		assert.Nil(t, err)
		assert.Equal(t, proposal.ChangeParameterProposal, *proposlPtr)
		err = vs.DeleteProposal(ctx, proposalID)
		assert.Nil(t, err)
		proposlPtr, err = vs.GetProposal(ctx, proposalID)
		assert.Nil(t, proposlPtr)
		assert.Equal(t, ErrGetProposal(), err)
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
