package vote

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

var (
	l0    = types.LNO(sdk.NewRat(0))
	l100  = types.LNO(sdk.NewRat(100))
	l200  = types.LNO(sdk.NewRat(200))
	l400  = types.LNO(sdk.NewRat(400))
	l1000 = types.LNO(sdk.NewRat(1000))
	l1600 = types.LNO(sdk.NewRat(1600))
	l2000 = types.LNO(sdk.NewRat(2000))

	c0    = types.Coin{0 * types.Decimals}
	c100  = types.Coin{100 * types.Decimals}
	c200  = types.Coin{200 * types.Decimals}
	c400  = types.Coin{400 * types.Decimals}
	c600  = types.Coin{600 * types.Decimals}
	c1000 = types.Coin{1000 * types.Decimals}
	c1200 = types.Coin{1200 * types.Decimals}
	c1600 = types.Coin{1600 * types.Decimals}
	c1900 = types.Coin{1900 * types.Decimals}
	c2000 = types.Coin{2000 * types.Decimals}
	c2600 = types.Coin{2600 * types.Decimals}
	c3200 = types.Coin{3200 * types.Decimals}
	c3600 = types.Coin{3600 * types.Decimals}
	c4600 = types.Coin{4600 * types.Decimals}
)

func TestVoterDepositBasic(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newVoteManager()
	gm := newGlobalProxy()
	ctx := getContext()
	handler := NewHandler(vm, lam, gm)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc1.AddCoin(ctx, c3600)
	acc1.Apply(ctx)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", l1600)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)
	handler(ctx, msg)

	// check acc1's money has been withdrawn
	acc1Balance, _ := acc1.GetBankBalance(ctx)
	assert.Equal(t, c400, acc1Balance)
	assert.Equal(t, true, vm.IsVoterExist(ctx, acc.AccountKey("user1")))

	// make sure the voter's account info is correct
	voter, _ := vm.GetVoter(ctx, acc.AccountKey("user1"))
	assert.Equal(t, c3200, voter.Deposit)
}

func TestDelegateBasic(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newVoteManager()
	gm := newGlobalProxy()
	ctx := getContext()
	handler := NewHandler(vm, lam, gm)

	// create test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc1.AddCoin(ctx, c2000)
	acc1.Apply(ctx)

	acc2 := createTestAccount(ctx, lam, "user2")
	acc2.AddCoin(ctx, c2000)
	acc2.Apply(ctx)

	acc3 := createTestAccount(ctx, lam, "user3")
	acc3.AddCoin(ctx, c2000)
	acc3.Apply(ctx)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", l1600)
	handler(ctx, msg)

	// let user2 delegate power to user1 twice
	msg2 := NewDelegateMsg("user2", "user1", l1000)
	handler(ctx, msg2)
	result2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, result2)

	// make sure the voter's voting power is correct
	voter, _ := vm.GetVoter(ctx, acc.AccountKey("user1"))
	assert.Equal(t, c1600, voter.Deposit)
	assert.Equal(t, c2000, voter.DelegatedPower)

	votingPower, _ := vm.GetVotingPower(ctx, "user1")
	assert.Equal(t, true, votingPower.IsEqual(c3600))
	acc2Balance, _ := acc2.GetBankBalance(ctx)
	assert.Equal(t, true, acc2Balance.IsEqual(c0))

	// let user3 delegate power to user1
	msg3 := NewDelegateMsg("user3", "user1", l1000)
	result3 := handler(ctx, msg3)
	assert.Equal(t, sdk.Result{}, result3)

	// check delegator list is correct
	delegators, _ := vm.GetAllDelegators(ctx, "user1")
	assert.Equal(t, 2, len(delegators))
	assert.Equal(t, acc.AccountKey("user2"), delegators[0])
	assert.Equal(t, acc.AccountKey("user3"), delegators[1])

	// check delegation are correct
	delegation1, _ := vm.GetDelegation(ctx, "user1", "user2")
	delegation2, _ := vm.GetDelegation(ctx, "user1", "user3")
	assert.Equal(t, c2000, delegation1.Amount)
	assert.Equal(t, c1000, delegation2.Amount)
}

func TestRevokeBasic(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newVoteManager()
	gm := newGlobalProxy()
	ctx := getContext()
	handler := NewHandler(vm, lam, gm)

	// create test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc1.AddCoin(ctx, c2000)
	acc1.Apply(ctx)

	acc2 := createTestAccount(ctx, lam, "user2")
	acc2.AddCoin(ctx, c2000)
	acc2.Apply(ctx)

	acc3 := createTestAccount(ctx, lam, "user3")
	acc3.AddCoin(ctx, c2000)
	acc3.Apply(ctx)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", l1600)
	handler(ctx, msg)

	// let user2 delegate power to user1 twice
	msg2 := NewDelegateMsg("user2", "user1", l1000)
	handler(ctx, msg2)

	// let user3 delegate power to user1
	msg3 := NewDelegateMsg("user3", "user1", l1000)
	handler(ctx, msg3)
	_, res := vm.GetDelegation(ctx, "user1", "user3")
	assert.Nil(t, res)

	// let user3 reovke delegation
	msg4 := NewRevokeDelegationMsg("user3", "user1")
	result := handler(ctx, msg4)
	assert.Equal(t, sdk.Result{}, result)

	// make sure user3 won't get coins immediately, but user1 power down immediately
	voter, _ := vm.GetVoter(ctx, "user1")
	acc3Balance, _ := acc3.GetBankBalance(ctx)
	_, getErr := vm.GetDelegation(ctx, "user1", "user3")
	assert.Equal(t, ErrGetDelegation(), getErr)
	assert.Equal(t, c1000, voter.DelegatedPower)
	assert.Equal(t, true, acc3Balance.IsEqual(c1000))

	// let user1 revoke voter candidancy
	msg5 := NewVoterRevokeMsg("user1")
	result2 := handler(ctx, msg5)
	assert.Equal(t, sdk.Result{}, result2)

	// make sure user2 wont get coins immediately, and delegatin was deleted
	_, err := vm.GetDelegation(ctx, "user1", "user2")
	_, err2 := vm.GetVoter(ctx, "user1")
	acc1Balance, _ := acc1.GetBankBalance(ctx)
	acc2Balance, _ := acc2.GetBankBalance(ctx)
	assert.Equal(t, ErrGetDelegation(), err)
	assert.Equal(t, ErrGetVoter(), err2)
	assert.Equal(t, c400, acc1Balance)
	assert.Equal(t, c1000, acc2Balance)
}

func TestWithdrawBasic(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newVoteManager()
	gm := newGlobalProxy()
	ctx := getContext()
	handler := NewHandler(vm, lam, gm)

	acc1 := createTestAccount(ctx, lam, "user1")
	acc1.AddCoin(ctx, c3600)
	acc1.Apply(ctx)

	// withdraw will fail if hasn't registed as voter
	illegalWithdrawMsg := NewVoterWithdrawMsg("user1", l1600)
	res := handler(ctx, illegalWithdrawMsg)
	assert.Equal(t, ErrIllegalWithdraw().Result(), res)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", l1600)
	handler(ctx, msg)

	msg2 := NewVoterWithdrawMsg("user1", l1000)
	result2 := handler(ctx, msg2)
	assert.Equal(t, ErrIllegalWithdraw().Result(), result2)

	msg3 := NewVoterWithdrawMsg("user1", l400)
	result3 := handler(ctx, msg3)
	assert.Equal(t, sdk.Result{}, result3)

	voter, _ := vm.GetVoter(ctx, "user1")
	assert.Equal(t, c1200, voter.Deposit)
}

func TestProposalBasic(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newVoteManager()
	gm := newGlobalProxy()
	ctx := getContext()
	handler := NewHandler(vm, lam, gm)
	vm.InitGenesis(ctx)

	rat := sdk.Rat{Denom: 10, Num: 5}
	para := ChangeParameterDescription{
		CDNAllocation: rat,
	}
	proposalID1 := ProposalKey(strconv.FormatInt(int64(1), 10))
	proposalID2 := ProposalKey(strconv.FormatInt(int64(2), 10))

	acc1 := createTestAccount(ctx, lam, "user1")
	acc1.AddCoin(ctx, c4600)
	acc1.Apply(ctx)

	// let user1 create a proposal
	msg := NewCreateProposalMsg("user1", para)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// invalid create
	invalidMsg := NewCreateProposalMsg("wqdkqwndkqwd", para)
	resultInvalid := handler(ctx, invalidMsg)
	assert.Equal(t, ErrUsernameNotFound().Result(), resultInvalid)

	result2 := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result2)

	proposal, _ := vm.GetProposal(ctx, proposalID1)
	assert.Equal(t, true, proposal.CDNAllocation.Equal(rat))

	// check use1's money has been reduced
	acc1Balance, _ := acc1.GetBankBalance(ctx)
	assert.Equal(t, true, acc1Balance.IsEqual(c600))

	// check proposal list is correct
	lst, _ := vm.GetProposalList(ctx)
	assert.Equal(t, 2, len(lst.OngoingProposal))
	assert.Equal(t, proposalID1, lst.OngoingProposal[0])
	assert.Equal(t, proposalID2, lst.OngoingProposal[1])

	// test delete proposal
	vm.DeleteProposal(ctx, proposalID2)
	_, getErr := vm.GetProposal(ctx, proposalID2)
	assert.Equal(t, ErrGetProposal(), getErr)

}

func TestVoteBasic(t *testing.T) {
	lam := newLinoAccountManager()
	vm := newVoteManager()
	gm := newGlobalProxy()
	ctx := getContext()
	handler := NewHandler(vm, lam, gm)

	rat := sdk.Rat{Denom: 10, Num: 5}
	para := ChangeParameterDescription{
		CDNAllocation: rat,
	}
	proposalID := int64(3)
	acc1 := createTestAccount(ctx, lam, "user1")
	acc1.AddCoin(ctx, c2000)
	acc1.Apply(ctx)

	acc2 := createTestAccount(ctx, lam, "user2")
	acc2.AddCoin(ctx, c2000)
	acc2.Apply(ctx)

	acc3 := createTestAccount(ctx, lam, "user3")
	acc3.AddCoin(ctx, c2000)
	acc3.Apply(ctx)

	// let user1 create a proposal
	msg := NewCreateProposalMsg("user1", para)
	handler(ctx, msg)

	// must become a voter before voting
	voteMsg := NewVoteMsg("user2", proposalID, true)
	result2 := handler(ctx, voteMsg)
	assert.Equal(t, ErrGetVoter().Result(), result2)

	depositMsg := NewVoterDepositMsg("user2", l1000)
	depositMsg2 := NewVoterDepositMsg("user3", l2000)
	handler(ctx, depositMsg)
	handler(ctx, depositMsg2)

	// invalid deposit
	invalidDepositMsg := NewVoterDepositMsg("1du1i2bdi12bud", l2000)
	res := handler(ctx, invalidDepositMsg)
	assert.Equal(t, ErrUsernameNotFound().Result(), res)

	// Now user2 can vote, vote on a non exist proposal
	invalidaVoteMsg := NewVoteMsg("user3", 10, true)
	voteRes := handler(ctx, invalidaVoteMsg)
	assert.Equal(t, ErrGetProposal().Result(), voteRes)

	// successfully vote
	voteMsg2 := NewVoteMsg("user2", proposalID, true)
	voteMsg3 := NewVoteMsg("user3", proposalID, true)
	handler(ctx, voteMsg2)
	handler(ctx, voteMsg3)

	// Check vote is correct
	vote, _ := vm.GetVote(ctx, ProposalKey(strconv.FormatInt(proposalID, 10)), "user2")
	assert.Equal(t, true, vote.Result)
	assert.Equal(t, acc.AccountKey("user2"), vote.Voter)

	voteList, _ := vm.GetAllVotes(ctx, ProposalKey(strconv.FormatInt(proposalID, 10)))
	assert.Equal(t, acc.AccountKey("user3"), voteList[1].Voter)

	// test delete vote
	vm.DeleteVote(ctx, ProposalKey(strconv.FormatInt(proposalID, 10)), "user2")
	vote, getErr := vm.GetVote(ctx, ProposalKey(strconv.FormatInt(proposalID, 10)), "user2")
	assert.Equal(t, ErrGetVote(), getErr)

}
