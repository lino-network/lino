package vote

import (
	"strconv"
	"testing"

	"github.com/lino-network/lino/tx/vote/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestVoterDepositBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(vm, am, gm)

	voteParam, _ := vm.paramHolder.GetVoteParam(ctx)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(voteParam.VoterMinDeposit))

	deposit := coinToString(voteParam.VoterMinDeposit)
	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", deposit)
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)
	handler(ctx, msg)

	// check acc1's money has been withdrawn
	acc1saving, _ := am.GetSavingFromBank(ctx, user1)
	assert.Equal(t, minBalance, acc1saving)
	assert.Equal(t, true, vm.IsVoterExist(ctx, user1))

	// make sure the voter's account info is correct
	voter, _ := vm.storage.GetVoter(ctx, user1)
	assert.Equal(t, voteParam.VoterMinDeposit, voter.Deposit)
}

func TestDelegateBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(vm, am, gm)

	voteParam, _ := vm.paramHolder.GetVoteParam(ctx)
	minBalance := types.NewCoinFromInt64(2000 * types.Decimals)
	// create test users
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(voteParam.VoterMinDeposit))
	user2 := createTestAccount(ctx, am, "user2", minBalance)
	user3 := createTestAccount(ctx, am, "user3", minBalance)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", coinToString(voteParam.VoterMinDeposit))
	handler(ctx, msg)

	delegatedCoin := types.NewCoinFromInt64(100 * types.Decimals)
	// let user2 delegate power to user1 twice
	msg2 := NewDelegateMsg("user2", "user1", coinToString(delegatedCoin))
	handler(ctx, msg2)
	result2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, result2)

	// make sure the voter's voting power is correct
	voter, _ := vm.storage.GetVoter(ctx, user1)
	assert.Equal(t, voteParam.VoterMinDeposit, voter.Deposit)
	assert.Equal(t, delegatedCoin.Plus(delegatedCoin), voter.DelegatedPower)

	votingPower, _ := vm.GetVotingPower(ctx, "user1")
	assert.Equal(t, true, votingPower.IsEqual(voteParam.VoterMinDeposit.Plus(delegatedCoin).Plus(delegatedCoin)))
	acc2Balance, _ := am.GetSavingFromBank(ctx, user2)
	assert.Equal(t, minBalance.Minus(delegatedCoin).Minus(delegatedCoin), acc2Balance)

	// let user3 delegate power to user1
	msg3 := NewDelegateMsg("user3", "user1", coinToString(delegatedCoin))
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
	assert.Equal(t, delegatedCoin.Plus(delegatedCoin), delegation1.Amount)
	assert.Equal(t, delegatedCoin, delegation2.Amount)
}

func TestRevokeBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(vm, am, gm)
	voteParam, _ := vm.paramHolder.GetVoteParam(ctx)
	minBalance := types.NewCoinFromInt64(2000 * types.Decimals)

	// create test users
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(voteParam.VoterMinDeposit))
	user2 := createTestAccount(ctx, am, "user2", minBalance)
	user3 := createTestAccount(ctx, am, "user3", minBalance)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", coinToString(voteParam.VoterMinDeposit))
	handler(ctx, msg)

	delegatedCoin := types.NewCoinFromInt64(100 * types.Decimals)
	// let user2 delegate power to user1
	msg2 := NewDelegateMsg("user2", "user1", coinToString(delegatedCoin))
	handler(ctx, msg2)

	// let user3 delegate power to user1
	msg3 := NewDelegateMsg("user3", "user1", coinToString(delegatedCoin))
	handler(ctx, msg3)

	_, res := vm.storage.GetDelegation(ctx, "user1", "user3")
	assert.Nil(t, res)

	// let user3 reovke delegation
	msg4 := NewRevokeDelegationMsg("user3", "user1")
	result := handler(ctx, msg4)
	assert.Equal(t, sdk.Result{}, result)

	// make sure user3 won't get coins immediately, but user1 power down immediately
	voter, _ := vm.storage.GetVoter(ctx, "user1")
	acc3Balance, _ := am.GetSavingFromBank(ctx, user3)
	_, err := vm.storage.GetDelegation(ctx, "user1", "user3")
	assert.Equal(t, ErrGetDelegation(), err)
	assert.Equal(t, delegatedCoin, voter.DelegatedPower)
	assert.Equal(t, minBalance.Minus(delegatedCoin), acc3Balance)

	// set user1 as validator (cannot revoke)
	referenceList := &model.ReferenceList{
		AllValidators: []types.AccountKey{user1},
	}
	vm.storage.SetReferenceList(ctx, referenceList)
	msg5 := NewVoterRevokeMsg("user1")
	result2 := handler(ctx, msg5)
	assert.Equal(t, ErrValidatorCannotRevoke().Result(), result2)

	// invalid user cannot revoke
	invalidMsg := NewVoterRevokeMsg("wqwdqwdasdsa")
	resultInvalid := handler(ctx, invalidMsg)
	assert.Equal(t, ErrGetVoter().Result(), resultInvalid)

	//  user1  can revoke voter candidancy now
	referenceList = &model.ReferenceList{
		AllValidators: []types.AccountKey{},
	}
	vm.storage.SetReferenceList(ctx, referenceList)
	result3 := handler(ctx, msg5)
	assert.Equal(t, sdk.Result{}, result3)

	// make sure user2 wont get coins immediately, and delegatin was deleted
	_, err2 := vm.storage.GetVoter(ctx, "user1")
	acc1Balance, _ := am.GetSavingFromBank(ctx, user1)
	acc2Balance, _ := am.GetSavingFromBank(ctx, user2)
	assert.Equal(t, ErrGetDelegation(), err)
	assert.Equal(t, ErrGetVoter(), err2)
	assert.Equal(t, minBalance, acc1Balance)
	assert.Equal(t, minBalance.Minus(delegatedCoin), acc2Balance)
}

func TestVoterWithdraw(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(vm, am, gm)
	voteParam, _ := vm.paramHolder.GetVoteParam(ctx)
	minBalance := types.NewCoinFromInt64(30 * types.Decimals)

	// create test users
	createTestAccount(ctx, am, "user1", minBalance.Plus(voteParam.VoterMinDeposit))

	// withdraw will fail if hasn't registed as voter
	illegalWithdrawMsg := NewVoterWithdrawMsg("user1", coinToString(voteParam.VoterMinWithdraw))
	res := handler(ctx, illegalWithdrawMsg)
	assert.Equal(t, ErrIllegalWithdraw().Result(), res)

	// let user1 register as voter
	msg := NewVoterDepositMsg("user1", coinToString(voteParam.VoterMinDeposit.Plus(voteParam.VoterMinWithdraw)))
	handler(ctx, msg)

	msg2 := NewVoterWithdrawMsg("user1", coinToString(minBalance.Plus(voteParam.VoterMinWithdraw)))
	result2 := handler(ctx, msg2)
	assert.Equal(t, ErrIllegalWithdraw().Result(), result2)

	msg3 := NewVoterWithdrawMsg("user1", coinToString(voteParam.VoterMinWithdraw))
	result3 := handler(ctx, msg3)
	assert.Equal(t, sdk.Result{}, result3)

	voter, _ := vm.storage.GetVoter(ctx, "user1")
	assert.Equal(t, voteParam.VoterMinDeposit, voter.Deposit)
}

func TestVoteBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(vm, am, gm)

	proposalID := int64(1)
	voteParam, _ := vm.paramHolder.GetVoteParam(ctx)
	minBalance := types.NewCoinFromInt64(2000 * types.Decimals)

	// create test users
	createTestAccount(ctx, am, "user1", minBalance.Plus(voteParam.VoterMinDeposit))
	user2 := createTestAccount(ctx, am, "user2", minBalance.Plus(voteParam.VoterMinDeposit))
	user3 := createTestAccount(ctx, am, "user3", minBalance.Plus(voteParam.VoterMinDeposit))

	// let user1 create a proposal
	referenceList := &model.ReferenceList{
		OngoingProposal: []types.ProposalKey{types.ProposalKey("1")},
	}
	vm.storage.SetReferenceList(ctx, referenceList)

	// must become a voter before voting
	voteMsg := NewVoteMsg("user2", proposalID, true)
	result2 := handler(ctx, voteMsg)
	assert.Equal(t, ErrGetVoter().Result(), result2)

	depositMsg := NewVoterDepositMsg("user2", coinToString(voteParam.VoterMinDeposit))
	depositMsg2 := NewVoterDepositMsg("user3", coinToString(voteParam.VoterMinDeposit))
	handler(ctx, depositMsg)
	handler(ctx, depositMsg2)

	// invalid deposit
	invalidDepositMsg := NewVoterDepositMsg("1du1i2bdi12bud", coinToString(voteParam.VoterMinDeposit))
	res := handler(ctx, invalidDepositMsg)
	assert.Equal(t, ErrUsernameNotFound().Result(), res)

	// Now user2 can vote, vote on a non exist proposal
	invalidaVoteMsg := NewVoteMsg("user3", 10, true)
	voteRes := handler(ctx, invalidaVoteMsg)
	assert.Equal(t, ErrNotOngoingProposal().Result(), voteRes)

	// successfully vote
	voteMsg2 := NewVoteMsg("user2", proposalID, true)
	voteMsg3 := NewVoteMsg("user3", proposalID, true)
	handler(ctx, voteMsg2)
	handler(ctx, voteMsg3)

	// user cannot vote again
	voteAgainMsg := NewVoteMsg("user3", proposalID, false)
	res = handler(ctx, voteAgainMsg)
	assert.Equal(t, ErrVoteExist().Result(), res)

	// Check vote is correct
	vote, _ := vm.storage.GetVote(ctx, types.ProposalKey(strconv.FormatInt(proposalID, 10)), "user2")
	assert.Equal(t, true, vote.Result)
	assert.Equal(t, user2, vote.Voter)

	voteList, _ := vm.storage.GetAllVotes(ctx, types.ProposalKey(strconv.FormatInt(proposalID, 10)))
	assert.Equal(t, user3, voteList[1].Voter)

	// test delete vote
	vm.storage.DeleteVote(ctx, types.ProposalKey(strconv.FormatInt(proposalID, 10)), "user2")
	vote, err := vm.storage.GetVote(ctx, types.ProposalKey(strconv.FormatInt(proposalID, 10)), "user2")
	assert.Equal(t, ErrGetVote(), err)

}

func TestDelegatorWithdraw(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(2000 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	user2 := createTestAccount(ctx, am, "user2", minBalance)
	handler := NewHandler(vm, am, gm)

	param, _ := vm.paramHolder.GetVoteParam(ctx)
	delegatedCoin := types.NewCoinFromInt64(100 * types.Decimals)
	delta := types.NewCoinFromInt64(1 * types.Decimals)
	vm.AddVoter(ctx, user1, param.VoterMinDeposit)

	cases := []struct {
		addDelegation bool
		delegatedCoin types.Coin
		delegator     types.AccountKey
		voter         types.AccountKey
		withdraw      types.Coin
		expectResult  sdk.Result
	}{
		{false, types.NewCoinFromInt64(0), user2, user1, param.DelegatorMinWithdraw, ErrIllegalWithdraw().Result()},
		{true, delegatedCoin, user2, user1, param.DelegatorMinWithdraw.Minus(delta), ErrIllegalWithdraw().Result()},
		{false, types.NewCoinFromInt64(0), user2, user1, delegatedCoin.Plus(delta), ErrIllegalWithdraw().Result()},
		{false, types.NewCoinFromInt64(0), user2, user1, delegatedCoin.Minus(delta), sdk.Result{}},
	}

	for _, cs := range cases {
		if cs.addDelegation {
			vm.AddDelegation(ctx, cs.voter, cs.delegator, cs.delegatedCoin)
		}
		msg := NewDelegatorWithdrawMsg(string(cs.delegator), string(cs.voter), coinToString(cs.withdraw))
		res := handler(ctx, msg)
		assert.Equal(t, cs.expectResult, res)
	}
}

func TestAddFrozenMoney(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	vm.InitGenesis(ctx)

	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user := createTestAccount(ctx, am, "user", minBalance)

	testCases := []struct {
		testName               string
		times                  int64
		interval               int64
		returnedCoin           types.Coin
		expectedFrozenListLen  int
		expectedFrozenMoney    types.Coin
		expectedFrozenTimes    int64
		expectedFrozenInterval int64
	}{
		{"return coin to user", 10, 2, types.NewCoinFromInt64(100), 1, types.NewCoinFromInt64(100), 10, 2},
		{"return coin to user", 100000, 20000, types.NewCoinFromInt64(100000), 2, types.NewCoinFromInt64(100000), 100000, 20000},
	}

	for _, tc := range testCases {
		err := returnCoinTo(
			ctx, "user", gm, am, tc.times, tc.interval, tc.returnedCoin)
		assert.Equal(t, nil, err)
		lst, err := am.GetFrozenMoneyList(ctx, user)
		assert.Equal(t, tc.expectedFrozenListLen, len(lst))
		assert.Equal(t, tc.expectedFrozenMoney, lst[len(lst)-1].Amount)
		assert.Equal(t, tc.expectedFrozenTimes, lst[len(lst)-1].Times)
		assert.Equal(t, tc.expectedFrozenInterval, lst[len(lst)-1].Interval)

	}
}
