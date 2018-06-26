package vote

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestCanBecomeValidator(t *testing.T) {
	ctx, am, vm, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	voteParam, _ := vm.paramHolder.GetVoteParam(ctx)
	valParam, _ := vm.paramHolder.GetValidatorParam(ctx)
	cases := []struct {
		addVoter     bool
		username     types.AccountKey
		coin         types.Coin
		expectResult bool
	}{
		{false, user1, types.NewCoinFromInt64(0), false},
		{true, user1, voteParam.VoterMinDeposit, false},
		{true, user1, valParam.ValidatorMinVotingDeposit, true},
	}

	for _, cs := range cases {
		if cs.addVoter {
			vm.AddVoter(ctx, cs.username, cs.coin)
		}
		actualRes := vm.CanBecomeValidator(ctx, cs.username)
		assert.Equal(t, cs.expectResult, actualRes)
	}
}

func TestAddVoter(t *testing.T) {
	ctx, am, vm, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	param, _ := vm.paramHolder.GetVoteParam(ctx)

	cases := []struct {
		username     types.AccountKey
		coin         types.Coin
		expectResult sdk.Error
	}{
		{user1, types.NewCoinFromInt64(100 * types.Decimals), ErrRegisterFeeNotEnough()},
		{user1, param.VoterMinDeposit, nil},
	}

	for _, cs := range cases {
		res := vm.AddVoter(ctx, cs.username, cs.coin)
		assert.Equal(t, cs.expectResult, res)
	}
}

func TestIsInValidatorList(t *testing.T) {
	ctx, am, vm, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	user2 := createTestAccount(ctx, am, "user2", minBalance)
	user3 := createTestAccount(ctx, am, "user3", minBalance)

	cases := []struct {
		username      types.AccountKey
		allValidators []types.AccountKey
		expectResult  bool
	}{
		{user1, []types.AccountKey{}, false},
		{user1, []types.AccountKey{user2, user3}, false},
		{user1, []types.AccountKey{user1}, true},
	}

	for _, cs := range cases {
		referenceList := &model.ReferenceList{
			AllValidators: cs.allValidators,
		}
		vm.storage.SetReferenceList(ctx, referenceList)
		res := vm.IsInValidatorList(ctx, cs.username)
		assert.Equal(t, cs.expectResult, res)
	}
}

func TestIsLegalVoterWithdraw(t *testing.T) {
	ctx, am, vm, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	param, _ := vm.paramHolder.GetVoteParam(ctx)

	vm.AddVoter(ctx, user1, param.VoterMinDeposit.Plus(types.NewCoinFromInt64(100*types.Decimals)))

	cases := []struct {
		allValidators []types.AccountKey
		username      types.AccountKey
		withdraw      types.Coin
		expectResult  bool
	}{
		{[]types.AccountKey{}, user1, param.VoterMinWithdraw.Minus(types.NewCoinFromInt64(1 * types.Decimals)), false},
		{[]types.AccountKey{}, user1, param.VoterMinWithdraw, true},
		{[]types.AccountKey{user1}, user1, param.VoterMinWithdraw, false},
		{[]types.AccountKey{}, user1, types.NewCoinFromInt64(100), false},
	}

	for _, cs := range cases {
		referenceList := &model.ReferenceList{
			AllValidators: cs.allValidators,
		}
		vm.storage.SetReferenceList(ctx, referenceList)
		res := vm.IsLegalVoterWithdraw(ctx, cs.username, cs.withdraw)
		assert.Equal(t, cs.expectResult, res)
	}
}

func TestIsLegalDelegatorWithdraw(t *testing.T) {
	ctx, am, vm, _ := setupTest(t, 0)
	minBalance := types.NewCoinFromInt64(1 * types.Decimals)
	user1 := createTestAccount(ctx, am, "user1", minBalance)
	user2 := createTestAccount(ctx, am, "user2", minBalance)
	param, _ := vm.paramHolder.GetVoteParam(ctx)

	vm.AddVoter(ctx, user1, param.VoterMinDeposit)

	cases := []struct {
		addDelegation bool
		delegatedCoin types.Coin
		delegator     types.AccountKey
		voter         types.AccountKey
		withdraw      types.Coin
		expectResult  bool
	}{
		{false, types.NewCoinFromInt64(0), user2, user1, param.DelegatorMinWithdraw, false},
		{true, types.NewCoinFromInt64(100 * types.Decimals), user2, user1, param.DelegatorMinWithdraw, true},
		{false, types.NewCoinFromInt64(0), user2, user1, types.NewCoinFromInt64(0), false},
		{false, types.NewCoinFromInt64(0), user2, user1, types.NewCoinFromInt64(101 * types.Decimals), false},
	}

	for _, cs := range cases {
		if cs.addDelegation {
			vm.AddDelegation(ctx, cs.voter, cs.delegator, cs.delegatedCoin)
		}
		res := vm.IsLegalDelegatorWithdraw(ctx, cs.voter, cs.delegator, cs.withdraw)
		assert.Equal(t, cs.expectResult, res)
	}
}
