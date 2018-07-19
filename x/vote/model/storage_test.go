package model

import (
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("vote")
)

func setup(t *testing.T) (sdk.Context, VoteStorage) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	vs := NewVoteStorage(TestKVStoreKey)
	return ctx, vs
}

func TestVoter(t *testing.T) {
	ctx, vs := setup(t)

	user := types.AccountKey("user")
	voter := Voter{
		Username:       user,
		Deposit:        types.NewCoinFromInt64(1000),
		DelegatedPower: types.NewCoinFromInt64(10000),
	}
	err := vs.SetVoter(ctx, user, &voter)
	assert.Nil(t, err)

	voterPtr, err := vs.GetVoter(ctx, user)
	assert.Nil(t, err)
	assert.Equal(t, voter, *voterPtr, "voter should be equal")

	vs.DeleteVoter(ctx, user)
	voterPtr, err = vs.GetVoter(ctx, user)
	assert.Nil(t, voterPtr)
	assert.Equal(t, ErrVoterNotFound(), err)
}

func TestVote(t *testing.T) {
	ctx, vs := setup(t)

	user1, user2, user3 :=
		types.AccountKey("user1"), types.AccountKey("user2"), types.AccountKey("user3")
	proposalID1, proposalID2 := types.ProposalKey("1"), types.ProposalKey("2")
	votingPower := types.NewCoinFromInt64(1000)

	testCases := []struct {
		testName      string
		isDelete      bool
		voter         types.AccountKey
		result        bool
		votingPower   types.Coin
		proposalID    types.ProposalKey
		expectedVotes []Vote
	}{
		{
			testName:    "user1 votes to proposal1 with agree",
			isDelete:    false,
			voter:       user1,
			result:      true,
			votingPower: votingPower,
			proposalID:  proposalID1,
			expectedVotes: []Vote{
				{
					Voter:       user1,
					VotingPower: votingPower,
					Result:      true,
				},
			},
		},
		{
			testName:    "user2 votes to proposal2 with agree",
			isDelete:    false,
			voter:       user2,
			result:      true,
			votingPower: votingPower,
			proposalID:  proposalID2,
			expectedVotes: []Vote{
				{
					Voter:       user2,
					VotingPower: votingPower,
					Result:      true,
				},
			},
		},
		{
			testName:    "user2 votes to proposal2 with disagree",
			isDelete:    false,
			voter:       user2,
			result:      false,
			votingPower: votingPower,
			proposalID:  proposalID2,
			expectedVotes: []Vote{
				{
					Voter:       user2,
					VotingPower: votingPower,
					Result:      false,
				},
			},
		},
		{
			testName:    "user3 votes to proposal2 with agree",
			isDelete:    false,
			voter:       user3,
			result:      true,
			votingPower: votingPower,
			proposalID:  proposalID2,
			expectedVotes: []Vote{
				{
					Voter:       user2,
					VotingPower: votingPower,
					Result:      false,
				},
				{
					Voter:       user3,
					VotingPower: votingPower,
					Result:      true,
				},
			},
		},
		{
			testName:    "user1 removes previous vote to proposal1",
			isDelete:    true,
			voter:       user1,
			result:      true,
			votingPower: votingPower,
			proposalID:  proposalID1,
		},
		{
			testName:    "user2 removes previous vote to proposal2",
			isDelete:    true,
			voter:       user2,
			result:      true,
			votingPower: votingPower,
			proposalID:  proposalID2,
			expectedVotes: []Vote{
				{
					Voter:       user3,
					VotingPower: votingPower,
					Result:      true,
				},
			},
		},
		{
			testName:    "user3 votes to proposal2 with disagree",
			isDelete:    false,
			voter:       user3,
			result:      false,
			votingPower: votingPower,
			proposalID:  proposalID2,
			expectedVotes: []Vote{
				{
					Voter:       user3,
					VotingPower: votingPower,
					Result:      false,
				},
			},
		},
		{
			testName:    "user2 votes to porposal2 with agree again",
			isDelete:    false,
			voter:       user2,
			result:      true,
			votingPower: votingPower,
			proposalID:  proposalID2,
			expectedVotes: []Vote{
				{
					Voter:       user2,
					VotingPower: votingPower,
					Result:      true,
				},
				{
					Voter:       user3,
					VotingPower: votingPower,
					Result:      false,
				},
			},
		},
	}

	for _, tc := range testCases {
		if tc.isDelete {
			vs.DeleteVote(ctx, tc.proposalID, tc.voter)
			votePtr, err := vs.GetVote(ctx, tc.proposalID, tc.voter)
			if err.Code() != types.CodeVoteNotFound {
				t.Errorf("%s: diff err code, got %v, want %v", tc.testName, err.Code(), types.CodeVoteNotFound)
			}
			if votePtr != nil {
				t.Errorf("%s: got non-empty vote, got %v, want nil", tc.testName, votePtr)
			}
		} else {
			vote := Vote{
				Voter:       tc.voter,
				Result:      tc.result,
				VotingPower: tc.votingPower,
			}
			err := vs.SetVote(ctx, tc.proposalID, tc.voter, &vote)
			if err != nil {
				t.Errorf("%s: failed to set vote, got non-empty err: %v", tc.testName, err)
			}

			votePtr, err := vs.GetVote(ctx, tc.proposalID, tc.voter)
			if err != nil {
				t.Errorf("%s: failed to get vote, got non-empty err: %v", tc.testName, err)
			}
			if !reflect.DeepEqual(*votePtr, vote) {
				t.Errorf("%s: diff vote, got %v, want %v", tc.testName, *votePtr, vote)
			}
		}
		allVotes, err := vs.GetAllVotes(ctx, tc.proposalID)
		if err != nil {
			t.Errorf("%s: failed to get all votes, got non-empty err: %v", tc.testName, err)
		}
		if !assert.Equal(t, allVotes, tc.expectedVotes) {
			t.Errorf("%s: diff votes, got %v, want %v", tc.testName, allVotes, tc.expectedVotes)
		}
	}
}

func TestDelegation(t *testing.T) {
	ctx, vs := setup(t)
	user1, user2, user3 :=
		types.AccountKey("user1"), types.AccountKey("user2"), types.AccountKey("user3")

	testCases := []struct {
		testName   string
		delegator  types.AccountKey
		amount     types.Coin
		delegateTo types.AccountKey
	}{
		{
			testName:   "user1 delegates to user2",
			delegator:  user1,
			amount:     types.NewCoinFromInt64(1),
			delegateTo: user2,
		},
		{
			testName:   "user1 delegates to user2 with more coins",
			delegator:  user1,
			amount:     types.NewCoinFromInt64(100),
			delegateTo: user2,
		},
		{
			testName:   "user1 delegates to user1",
			delegator:  user1,
			amount:     types.NewCoinFromInt64(100),
			delegateTo: user1,
		},
		{
			testName:   "user3 delegates to user1",
			delegator:  user3,
			amount:     types.NewCoinFromInt64(100),
			delegateTo: user1,
		},
		{
			testName:   "user3 delegates to user2",
			delegator:  user3,
			amount:     types.NewCoinFromInt64(100),
			delegateTo: user2,
		},
	}

	for _, tc := range testCases {
		err := vs.SetDelegation(ctx, tc.delegateTo, tc.delegator,
			&Delegation{tc.delegator, tc.amount})
		if err != nil {
			t.Errorf("%s: failed to set delegation, got non-empty err: %v", tc.testName, err)
		}

		delegationPtr, err := vs.GetDelegation(ctx, tc.delegateTo, tc.delegator)
		if err != nil {
			t.Errorf("%s: failed to get delegation, got non-empty err: %v", tc.testName, err)
		}
		if !assert.Equal(t, Delegation{tc.delegator, tc.amount}, *delegationPtr) {
			t.Errorf("%s: diff delegation, got %v, want %v", tc.testName, *delegationPtr, Delegation{tc.delegator, tc.amount})
		}

		err = vs.DeleteDelegation(ctx, tc.delegateTo, tc.delegator)
		if err != nil {
			t.Errorf("%s: failed to delete delegation, got non-empty err: %v", tc.testName, err)
		}

		delegationPtr, err = vs.GetDelegation(ctx, tc.delegateTo, tc.delegator)
		if err.Code() != types.CodeDelegationNotFound {
			t.Errorf("%s: diff err code, got %v, want %v", tc.testName, err.Code(), types.CodeDelegationNotFound)
		}
		if delegationPtr != nil {
			t.Errorf("%s: got non-empty delegation %v, want nil", tc.testName, delegationPtr)
		}
	}
}

func TestAllDelegation(t *testing.T) {
	ctx, vs := setup(t)
	user1, user2, user3 :=
		types.AccountKey("user1"), types.AccountKey("user2"), types.AccountKey("user3")

	testCases := []struct {
		testName           string
		delegator          types.AccountKey
		amount             types.Coin
		delegateTo         types.AccountKey
		expectedDelegators []types.AccountKey
	}{
		{
			testName:           "user1 delegates to user2",
			delegator:          user1,
			amount:             types.NewCoinFromInt64(1),
			delegateTo:         user2,
			expectedDelegators: []types.AccountKey{user1},
		},
		{
			testName:           "user1 delegates to user2 with more coins",
			delegator:          user1,
			amount:             types.NewCoinFromInt64(100),
			delegateTo:         user2,
			expectedDelegators: []types.AccountKey{user1},
		},
		{
			testName:           "user1 delegates to user3",
			delegator:          user1,
			amount:             types.NewCoinFromInt64(1),
			delegateTo:         user3,
			expectedDelegators: []types.AccountKey{user1},
		},
		{
			testName:           "user2 delegates to user1",
			delegator:          user2,
			amount:             types.NewCoinFromInt64(1),
			delegateTo:         user1,
			expectedDelegators: []types.AccountKey{user2},
		},
		{
			testName:           "user3 delegates to user1",
			delegator:          user3,
			amount:             types.NewCoinFromInt64(1),
			delegateTo:         user1,
			expectedDelegators: []types.AccountKey{user2, user3},
		},
	}

	for _, tc := range testCases {
		err := vs.SetDelegation(ctx, tc.delegateTo, tc.delegator,
			&Delegation{tc.delegator, tc.amount})
		if err != nil {
			t.Errorf("%s: failed to set delegation, got non-empty err: %v", tc.testName, err)
		}

		delegators, err := vs.GetAllDelegators(ctx, tc.delegateTo)
		if err != nil {
			t.Errorf("%s: failed to get all delegations, got non-empty err: %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.expectedDelegators, delegators) {
			t.Errorf("%s: diff delegators, got %v, want %v", tc.testName, delegators, tc.expectedDelegators)
		}
	}
}
