package vote

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestVoterDepositBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(vm, am, gm)
	minBalance := types.NewCoinFromInt64(100 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)

	deposit := types.NewCoinFromInt64(10 * types.Decimals)
	// let user1 register as voter
	msg := NewStakeInMsg("user1", coinToString(deposit))
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// check acc1's money has been withdrawn
	acc1saving, _ := am.GetSavingFromBank(ctx, user1)
	assert.Equal(t, minBalance.Minus(deposit), acc1saving)
	assert.Equal(t, true, vm.DoesVoterExist(ctx, user1))

	// make sure the voter's account info is correct
	voter, _ := vm.storage.GetVoter(ctx, user1)
	assert.Equal(t, deposit, voter.LinoStake)
}

func TestDelegateBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(vm, am, gm)

	minBalance := types.NewCoinFromInt64(2000 * types.Decimals)
	deposit := types.NewCoinFromInt64(10 * types.Decimals)

	// create test users
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(deposit))
	user2 := createTestAccount(ctx, am, "user2", minBalance)
	user3 := createTestAccount(ctx, am, "user3", minBalance)

	// let user1 register as vote
	msg := NewStakeInMsg("user1", coinToString(deposit))
	handler(ctx, msg)

	delegatedCoin := types.NewCoinFromInt64(100 * types.Decimals)
	// let user2 delegate power to user1 twice
	msg2 := NewDelegateMsg("user2", "user1", coinToString(delegatedCoin))
	handler(ctx, msg2)
	result2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, result2)

	// make sure the voter's voting power is correct
	voter, _ := vm.storage.GetVoter(ctx, user1)
	assert.Equal(t, deposit, voter.LinoStake)
	assert.Equal(t, delegatedCoin.Plus(delegatedCoin), voter.DelegatedPower)

	votingPower, _ := vm.GetVotingPower(ctx, "user1")
	assert.Equal(t, true, votingPower.IsEqual(deposit.Plus(delegatedCoin).Plus(delegatedCoin)))
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

	minBalance := types.NewCoinFromInt64(2000 * types.Decimals)
	deposit := types.NewCoinFromInt64(10 * types.Decimals)

	// create test users
	user1 := createTestAccount(ctx, am, "user1", minBalance.Plus(deposit))
	user2 := createTestAccount(ctx, am, "user2", minBalance)
	user3 := createTestAccount(ctx, am, "user3", minBalance)

	// let user1 register as voter
	msg := NewStakeInMsg("user1", coinToString(deposit))
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
	assert.Equal(t, model.ErrDelegationNotFound(), err)
	assert.Equal(t, delegatedCoin, voter.DelegatedPower)
	assert.Equal(t, minBalance.Minus(delegatedCoin), acc3Balance)

	// user1 can revoke voter candidancy now
	referenceList := &model.ReferenceList{
		AllValidators: []types.AccountKey{},
	}
	vm.storage.SetReferenceList(ctx, referenceList)
	msg5 := NewStakeOutMsg("user1", coinToString(deposit))

	vm.storage.SetReferenceList(ctx, referenceList)
	result2 := handler(ctx, msg5)
	assert.Equal(t, sdk.Result{}, result2)

	// make sure user2 wont get coins immediately, and delegatin was deleted
	_, err2 := vm.storage.GetVoter(ctx, "user1")
	acc1Balance, _ := am.GetSavingFromBank(ctx, user1)
	acc2Balance, _ := am.GetSavingFromBank(ctx, user2)
	assert.Equal(t, model.ErrDelegationNotFound(), err)
	assert.Equal(t, model.ErrVoterNotFound(), err2)
	assert.Equal(t, minBalance, acc1Balance)
	assert.Equal(t, minBalance.Minus(delegatedCoin), acc2Balance)
}

func TestVoterWithdraw(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	handler := NewHandler(vm, am, gm)
	voteParam, _ := vm.paramHolder.GetVoteParam(ctx)
	minBalance := types.NewCoinFromInt64(30 * types.Decimals)
	deposit := types.NewCoinFromInt64(10 * types.Decimals)

	// create test users
	createTestAccount(ctx, am, "user1", minBalance.Plus(deposit))

	// withdraw will fail if hasn't registered as voter
	illegalWithdrawMsg := NewStakeOutMsg("user1", coinToString(voteParam.VoterMinWithdraw))
	res := handler(ctx, illegalWithdrawMsg)
	assert.Equal(t, ErrIllegalWithdraw().Result(), res)

	// let user1 register as voter
	msg := NewStakeInMsg("user1", coinToString(deposit.Plus(voteParam.VoterMinWithdraw)))
	handler(ctx, msg)

	// invalid deposit
	invalidDepositMsg := NewStakeInMsg("1du1i2bdi12bud", coinToString(deposit))
	res = handler(ctx, invalidDepositMsg)
	assert.Equal(t, ErrAccountNotFound().Result(), res)

	msg2 := NewStakeOutMsg("user1", coinToString(minBalance.Plus(voteParam.VoterMinWithdraw)))
	result2 := handler(ctx, msg2)
	assert.Equal(t, ErrIllegalWithdraw().Result(), result2)

	msg3 := NewStakeOutMsg("user1", coinToString(voteParam.VoterMinWithdraw))
	result3 := handler(ctx, msg3)
	assert.Equal(t, sdk.Result{}, result3)

	voter, _ := vm.storage.GetVoter(ctx, "user1")
	assert.Equal(t, deposit, voter.LinoStake)
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
	deposit := types.NewCoinFromInt64(10 * types.Decimals)
	vm.AddVoter(ctx, user1, deposit)

	testCases := []struct {
		testName       string
		addDelegation  bool
		delegatedCoin  types.Coin
		delegator      types.AccountKey
		voter          types.AccountKey
		withdraw       types.Coin
		expectedResult sdk.Result
	}{
		{
			testName:       "no delegation exist, can't withdraw",
			addDelegation:  false,
			delegatedCoin:  types.NewCoinFromInt64(0),
			delegator:      user2,
			voter:          user1,
			withdraw:       param.DelegatorMinWithdraw,
			expectedResult: ErrIllegalWithdraw().Result(),
		},
		{
			testName:       "withdraw amount is less than minimum requirement",
			addDelegation:  true,
			delegatedCoin:  delegatedCoin,
			delegator:      user2,
			voter:          user1,
			withdraw:       param.DelegatorMinWithdraw.Minus(delta),
			expectedResult: ErrIllegalWithdraw().Result(),
		},
		{
			testName:       "no delegation exist, can't withdraw delegatedCoin+delta",
			addDelegation:  false,
			delegatedCoin:  types.NewCoinFromInt64(0),
			delegator:      user2,
			voter:          user1,
			withdraw:       delegatedCoin.Plus(delta),
			expectedResult: ErrIllegalWithdraw().Result(),
		},
		{
			testName:       "no delegation exist, can't withdraw delegatedCoin-delta",
			addDelegation:  false,
			delegatedCoin:  types.NewCoinFromInt64(0),
			delegator:      user2,
			voter:          user1,
			withdraw:       delegatedCoin.Minus(delta),
			expectedResult: sdk.Result{},
		},
	}

	for _, tc := range testCases {
		if tc.addDelegation {
			err := vm.AddDelegation(ctx, tc.voter, tc.delegator, tc.delegatedCoin)
			if err != nil {
				t.Errorf("%s: failed to add delegation, got err %v", tc.testName, err)
			}
		}
		msg := NewDelegatorWithdrawMsg(string(tc.delegator), string(tc.voter), coinToString(tc.withdraw))
		res := handler(ctx, msg)
		if !assert.Equal(t, tc.expectedResult, res) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, res, tc.expectedResult)
		}
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
		{
			testName:               "return coin to user",
			times:                  10,
			interval:               2,
			returnedCoin:           types.NewCoinFromInt64(100),
			expectedFrozenListLen:  1,
			expectedFrozenMoney:    types.NewCoinFromInt64(100),
			expectedFrozenTimes:    10,
			expectedFrozenInterval: 2,
		},
		{
			testName:               "return coin to user again",
			times:                  100000,
			interval:               20000,
			returnedCoin:           types.NewCoinFromInt64(100000),
			expectedFrozenListLen:  2,
			expectedFrozenMoney:    types.NewCoinFromInt64(100000),
			expectedFrozenTimes:    100000,
			expectedFrozenInterval: 20000,
		},
	}

	for _, tc := range testCases {
		err := returnCoinTo(
			ctx, "user", gm, am, tc.times, tc.interval, tc.returnedCoin, types.VoteReturnCoin)
		if err != nil {
			t.Errorf("%s: failed to return coin, got err %v", tc.testName, err)
		}
		lst, err := am.GetFrozenMoneyList(ctx, user)
		if err != nil {
			t.Errorf("%s: failed to get frozen money list, got err %v", tc.testName, err)
		}
		if len(lst) != tc.expectedFrozenListLen {
			t.Errorf("%s: diff list len, got %v, want %v", tc.testName, len(lst), tc.expectedFrozenListLen)
		}
		if !lst[len(lst)-1].Amount.IsEqual(tc.expectedFrozenMoney) {
			t.Errorf("%s: diff amount, got %v, want %v", tc.testName, lst[len(lst)-1].Amount, tc.expectedFrozenMoney)
		}
		if lst[len(lst)-1].Times != tc.expectedFrozenTimes {
			t.Errorf("%s: diff times, got %v, want %v", tc.testName, lst[len(lst)-1].Times, tc.expectedFrozenTimes)
		}
		if lst[len(lst)-1].Interval != tc.expectedFrozenInterval {
			t.Errorf("%s: diff interval, got %v, want %v", tc.testName, lst[len(lst)-1].Interval, tc.expectedFrozenInterval)
		}
	}
}

func TestDeleteVoteBasic(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	vm.InitGenesis(ctx)
	handler := NewHandler(vm, am, gm)

	proposalID1 := types.ProposalKey("1")
	minBalance := types.NewCoinFromInt64(2000 * types.Decimals)
	deposit := types.NewCoinFromInt64(10 * types.Decimals)

	// create test users
	user2 := createTestAccount(ctx, am, "user2", minBalance.Plus(deposit))
	depositMsg := NewStakeInMsg("user2", coinToString(deposit))
	handler(ctx, depositMsg)

	// add vote
	_ = vm.AddVote(ctx, proposalID1, user2, true)

	voteList, _ := vm.storage.GetAllVotes(ctx, proposalID1)
	assert.Equal(t, user2, voteList[0].Voter)

	// test delete vote
	vm.storage.DeleteVote(ctx, proposalID1, "user2")
	_, err := vm.storage.GetVote(ctx, proposalID1, "user2")
	assert.Equal(t, model.ErrVoteNotFound(), err)
}
