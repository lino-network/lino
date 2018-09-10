package model

import (
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
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

func TestDoesVoterExist(t *testing.T) {
	ctx, vs := setup(t)
	user := types.AccountKey("user")
	voter := &Voter{
		Username:       user,
		LinoPower:      types.NewCoinFromInt64(1000),
		DelegatedPower: types.NewCoinFromInt64(10000),
	}
	err := vs.SetVoter(ctx, user, voter)
	if err != nil {
		t.Errorf("%s: failed to set voter, got err %v", "TestDoesVoterExist", err)
	}

	testCases := []struct {
		testName  string
		accKey    types.AccountKey
		wantExist bool
	}{
		{
			testName:  "voter exist",
			accKey:    user,
			wantExist: true,
		},
		{
			testName:  "voter doesn't exist",
			accKey:    types.AccountKey("acc"),
			wantExist: false,
		},
	}
	for _, tc := range testCases {
		gotExist := vs.DoesVoterExist(ctx, tc.accKey)
		if gotExist != tc.wantExist {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, gotExist, tc.wantExist)
		}
	}
}

func TestDoesVoteExist(t *testing.T) {
	ctx, vs := setup(t)

	user1 := types.AccountKey("user1")
	proposalID1 := types.ProposalKey("1")
	votingPower := types.NewCoinFromInt64(1000)

	vote := &Vote{
		Voter:       user1,
		Result:      true,
		VotingPower: votingPower,
	}
	err := vs.SetVote(ctx, proposalID1, user1, vote)
	if err != nil {
		t.Errorf("%s: failed to set vote, got err: %v", "TestDoesVoteExist", err)
	}

	testCases := []struct {
		testName   string
		proposalID types.ProposalKey
		voter      types.AccountKey
		wantExist  bool
	}{
		{
			testName:   "vote exist",
			proposalID: proposalID1,
			voter:      user1,
			wantExist:  true,
		},
		{
			testName:   "vote doesn't exist because voter doesn't exist",
			proposalID: proposalID1,
			voter:      types.AccountKey("acc"),
			wantExist:  false,
		},
		{
			testName:   "vote doesn't exist because proposal doesn't exist",
			proposalID: types.ProposalKey("100"),
			voter:      user1,
			wantExist:  false,
		},
	}
	for _, tc := range testCases {
		gotExist := vs.DoesVoteExist(ctx, tc.proposalID, tc.voter)
		if gotExist != tc.wantExist {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, gotExist, tc.wantExist)
		}
	}
}

func TestDoesDelegationExist(t *testing.T) {
	ctx, vs := setup(t)

	user1, user2 := types.AccountKey("user1"), types.AccountKey("user2")
	delegation := &Delegation{
		Delegator: user1,
		Amount:    types.NewCoinFromInt64(1),
	}
	err := vs.SetDelegation(ctx, user1, user2, delegation)
	if err != nil {
		t.Errorf("%s: failed to set delegation, got non-empty err: %v", "TestDoesDelegationExist", err)
	}

	testCases := []struct {
		testName  string
		voter     types.AccountKey
		delegator types.AccountKey
		wantExist bool
	}{
		{
			testName:  "delegation exist",
			voter:     user1,
			delegator: user2,
			wantExist: true,
		},
		{
			testName:  "delegation doesn't exist",
			voter:     user2,
			delegator: user1,
			wantExist: false,
		},
		{
			testName:  "delegation doesn't exist because voter doesn't exist",
			voter:     types.AccountKey("voter"),
			delegator: user2,
			wantExist: false,
		},
		{
			testName:  "delegation doesn't exist because delegator doesn't exist",
			voter:     user1,
			delegator: types.AccountKey("delegator"),
			wantExist: false,
		},
	}
	for _, tc := range testCases {
		gotExist := vs.DoesDelegationExist(ctx, tc.voter, tc.delegator)
		if gotExist != tc.wantExist {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, gotExist, tc.wantExist)
		}
	}
}

func TestVoter(t *testing.T) {
	ctx, vs := setup(t)

	user := types.AccountKey("user")
	voter := Voter{
		Username:       user,
		LinoPower:      types.NewCoinFromInt64(1000),
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
