package vote

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/vote/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

var (
	l400  = types.LNO(sdk.NewRat(400))
	l1000 = types.LNO(sdk.NewRat(1000))
	l1600 = types.LNO(sdk.NewRat(1600))
	l2000 = types.LNO(sdk.NewRat(2000))

	c400  = types.Coin{400 * types.Decimals}
	c600  = types.Coin{600 * types.Decimals}
	c1000 = types.Coin{1000 * types.Decimals}
	c1200 = types.Coin{1200 * types.Decimals}
	c1600 = types.Coin{1600 * types.Decimals}
	c2000 = types.Coin{2000 * types.Decimals}
	c3200 = types.Coin{3200 * types.Decimals}
	c3600 = types.Coin{3600 * types.Decimals}
	c4600 = types.Coin{4600 * types.Decimals}
)

func TestVoterDepositBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(*vm, *am, *gm)

	// create two test users
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c3600)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", l1600)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)
	handler(ctx, msg)

	// check acc1's money has been withdrawn
	acc1Balance, _ := am.GetBankBalance(ctx, user1)
	assert.Equal(t, c400.Plus(initCoin), acc1Balance)
	assert.Equal(t, true, vm.IsVoterExist(ctx, user1))

	// make sure the voter's account info is correct
	voter, _ := vm.storage.GetVoter(ctx, user1)
	assert.Equal(t, c3200, voter.Deposit)
}

func TestDelegateBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(*vm, *am, *gm)

	// create test users
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)

	user2 := createTestAccount(ctx, am, "user2")
	am.AddCoin(ctx, user2, c2000)

	user3 := createTestAccount(ctx, am, "user3")
	am.AddCoin(ctx, user3, c2000)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", l1600)
	handler(ctx, msg)

	// let user2 delegate power to user1 twice
	msg2 := NewDelegateMsg("user2", "user1", l1000)
	handler(ctx, msg2)
	result2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, result2)

	// make sure the voter's voting power is correct
	voter, _ := vm.storage.GetVoter(ctx, user1)
	assert.Equal(t, c1600, voter.Deposit)
	assert.Equal(t, c2000, voter.DelegatedPower)

	votingPower, _ := vm.GetVotingPower(ctx, "user1")
	assert.Equal(t, true, votingPower.IsEqual(c3600))
	acc2Balance, _ := am.GetBankBalance(ctx, user2)
	assert.Equal(t, acc2Balance, initCoin)

	// let user3 delegate power to user1
	msg3 := NewDelegateMsg("user3", "user1", l1000)
	result3 := handler(ctx, msg3)
	assert.Equal(t, sdk.Result{}, result3)

	// check delegator list is correct
	delegators, _ := vm.storage.GetAllDelegators(ctx, "user1")
	assert.Equal(t, 2, len(delegators))
	assert.Equal(t, user2, delegators[0])
	assert.Equal(t, user3, delegators[1])

	// check delegation are correct
	delegation1, _ := vm.storage.GetDelegation(ctx, "user1", "user2")
	delegation2, _ := vm.storage.GetDelegation(ctx, "user1", "user3")
	assert.Equal(t, c2000, delegation1.Amount)
	assert.Equal(t, c1000, delegation2.Amount)
}

func TestRevokeBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(*vm, *am, *gm)

	// create test users
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)

	user2 := createTestAccount(ctx, am, "user2")
	am.AddCoin(ctx, user2, c2000)

	user3 := createTestAccount(ctx, am, "user3")
	am.AddCoin(ctx, user3, c2000)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", l1600)
	handler(ctx, msg)

	// let user2 delegate power to user1
	msg2 := NewDelegateMsg("user2", "user1", l1000)
	handler(ctx, msg2)

	// let user3 delegate power to user1
	msg3 := NewDelegateMsg("user3", "user1", l1000)
	handler(ctx, msg3)
	_, res := vm.storage.GetDelegation(ctx, "user1", "user3")
	assert.Nil(t, res)

	// let user3 reovke delegation
	msg4 := NewRevokeDelegationMsg("user3", "user1")
	result := handler(ctx, msg4)
	assert.Equal(t, sdk.Result{}, result)

	// make sure user3 won't get coins immediately, but user1 power down immediately
	voter, _ := vm.storage.GetVoter(ctx, "user1")
	acc3Balance, _ := am.GetBankBalance(ctx, user3)
	_, getErr := vm.storage.GetDelegation(ctx, "user1", "user3")
	assert.Equal(t, ErrGetDelegation(), getErr)
	assert.Equal(t, c1000, voter.DelegatedPower)
	assert.Equal(t, acc3Balance, c1000.Plus(initCoin))

	// set user1 as validator (cannot revoke)
	ctx = WithAllValidators(ctx, []types.AccountKey{user1})
	msg5 := NewVoterRevokeMsg("user1")
	result2 := handler(ctx, msg5)
	assert.Equal(t, ErrValidatorCannotRevoke().Result(), result2)

	// invalid user cannot revoke
	invalidMsg := NewVoterRevokeMsg("wqwdqwdasdsa")
	resultInvalid := handler(ctx, invalidMsg)
	assert.Equal(t, ErrGetVoter().Result(), resultInvalid)

	//  user1  can revoke voter candidancy now
	ctx = WithAllValidators(ctx, []types.AccountKey{})
	result3 := handler(ctx, msg5)
	assert.Equal(t, sdk.Result{}, result3)

	// make sure user2 wont get coins immediately, and delegatin was deleted
	_, err := vm.storage.GetDelegation(ctx, "user1", "user2")
	_, err2 := vm.storage.GetVoter(ctx, "user1")
	acc1Balance, _ := am.GetBankBalance(ctx, user1)
	acc2Balance, _ := am.GetBankBalance(ctx, user2)
	assert.Equal(t, ErrGetDelegation(), err)
	assert.Equal(t, ErrGetVoter(), err2)
	assert.Equal(t, c400.Plus(initCoin), acc1Balance)
	assert.Equal(t, c1000.Plus(initCoin), acc2Balance)
}

func TestVoterWithdraw(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(*vm, *am, *gm)

	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c3600)

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

	voter, _ := vm.storage.GetVoter(ctx, "user1")
	assert.Equal(t, c1200, voter.Deposit)
}

func TestProposalBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(*vm, *am, *gm)
	vm.InitGenesis(ctx)

	rat := sdk.Rat{Denom: 10, Num: 5}
	para := model.ChangeParameterDescription{
		CDNAllocation: rat,
	}
	proposalID1 := types.ProposalKey(strconv.FormatInt(int64(4), 10))
	proposalID2 := types.ProposalKey(strconv.FormatInt(int64(5), 10))

	user1 := createTestAccount(ctx, am, "user1")

	// let user1 create a proposal (not enough coins)
	msg := NewCreateProposalMsg("user1", para)
	result := handler(ctx, msg)
	assert.Equal(t, acc.ErrAccountCoinNotEnough().Result(), result)

	am.AddCoin(ctx, user1, c4600)
	resultPass := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, resultPass)

	// invalid create
	invalidMsg := NewCreateProposalMsg("wqdkqwndkqwd", para)
	resultInvalid := handler(ctx, invalidMsg)
	assert.Equal(t, ErrUsernameNotFound().Result(), resultInvalid)

	result2 := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result2)

	proposal, _ := vm.storage.GetProposal(ctx, proposalID1)
	assert.Equal(t, true, proposal.CDNAllocation.Equal(rat))

	// check use1's money has been reduced
	acc1Balance, _ := am.GetBankBalance(ctx, user1)
	assert.Equal(t, acc1Balance, c600.Plus(initCoin))

	// check proposal list is correct
	lst, _ := vm.storage.GetProposalList(ctx)
	assert.Equal(t, 2, len(lst.OngoingProposal))
	assert.Equal(t, proposalID1, lst.OngoingProposal[0])
	assert.Equal(t, proposalID2, lst.OngoingProposal[1])

	// test delete proposal
	vm.storage.DeleteProposal(ctx, proposalID2)
	_, getErr := vm.storage.GetProposal(ctx, proposalID2)
	assert.Equal(t, ErrGetProposal(), getErr)

}

func TestVoteBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(*vm, *am, *gm)

	rat := sdk.Rat{Denom: 10, Num: 5}
	para := model.ChangeParameterDescription{
		CDNAllocation: rat,
	}
	proposalID := int64(6)
	user1 := createTestAccount(ctx, am, "user1")
	am.AddCoin(ctx, user1, c2000)

	user2 := createTestAccount(ctx, am, "user2")
	am.AddCoin(ctx, user2, c2000)

	user3 := createTestAccount(ctx, am, "user3")
	am.AddCoin(ctx, user3, c2000)

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
	vote, _ := vm.storage.GetVote(ctx, types.ProposalKey(strconv.FormatInt(proposalID, 10)), "user2")
	assert.Equal(t, true, vote.Result)
	assert.Equal(t, user2, vote.Voter)

	voteList, _ := vm.storage.GetAllVotes(ctx, types.ProposalKey(strconv.FormatInt(proposalID, 10)))
	assert.Equal(t, user3, voteList[1].Voter)

	// test delete vote
	vm.storage.DeleteVote(ctx, types.ProposalKey(strconv.FormatInt(proposalID, 10)), "user2")
	vote, getErr := vm.storage.GetVote(ctx, types.ProposalKey(strconv.FormatInt(proposalID, 10)), "user2")
	assert.Equal(t, ErrGetVote(), getErr)

}

func TestDelegatorWithdraw(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	user1 := createTestAccount(ctx, am, "user1")
	user2 := createTestAccount(ctx, am, "user2")
	handler := NewHandler(*vm, *am, *gm)
	vm.AddVoter(ctx, user1, types.VoterMinDeposit)

	cases := []struct {
		addDelegation bool
		delegatedCoin types.Coin
		delegator     types.AccountKey
		voter         types.AccountKey
		withdraw      types.LNO
		expectResult  sdk.Result
	}{
		{false, types.NewCoin(0), user2, user1, types.DelegatorMinWithdraw.ToRat(), ErrIllegalWithdraw().Result()},
		{true, types.NewCoin(100 * types.Decimals), user2, user1, sdk.NewRat(1, 10), ErrIllegalWithdraw().Result()},
		{false, types.NewCoin(0), user2, user1, sdk.NewRat(101), ErrIllegalWithdraw().Result()},
		{false, types.NewCoin(0), user2, user1, sdk.NewRat(10), sdk.Result{}},
	}

	for _, cs := range cases {
		if cs.addDelegation {
			vm.AddDelegation(ctx, cs.voter, cs.delegator, cs.delegatedCoin)
		}
		msg := NewDelegatorWithdrawMsg(string(cs.delegator), string(cs.voter), cs.withdraw)
		res := handler(ctx, msg)
		assert.Equal(t, cs.expectResult, res)
	}
}
