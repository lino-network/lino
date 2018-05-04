package vote

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/vote/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestCanBecomeValidator(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	user1 := createTestAccount(ctx, am, "user1")
	voterMinDeposit, _ := gm.GetVoterMinDeposit(ctx)
	validatorMinVotingDeposit, _ := gm.GetValidatorMinVotingDeposit(ctx)
	cases := []struct {
		addVoter     bool
		username     types.AccountKey
		coin         types.Coin
		expectResult bool
	}{
		{false, user1, types.NewCoin(0), false},
		{true, user1, voterMinDeposit, false},
		{true, user1, validatorMinVotingDeposit, true},
	}

	for _, cs := range cases {
		if cs.addVoter {
			vm.AddVoter(ctx, cs.username, cs.coin, gm)
		}
		actualRes := vm.CanBecomeValidator(ctx, cs.username, gm)
		assert.Equal(t, cs.expectResult, actualRes)
	}
}

func TestAddVoter(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	user1 := createTestAccount(ctx, am, "user1")
	voterMinDeposit, _ := gm.GetVoterMinDeposit(ctx)

	cases := []struct {
		username     types.AccountKey
		coin         types.Coin
		expectResult sdk.Error
	}{
		{user1, types.NewCoin(100 * types.Decimals), ErrRegisterFeeNotEnough()},
		{user1, voterMinDeposit, nil},
	}

	for _, cs := range cases {
		res := vm.AddVoter(ctx, cs.username, cs.coin, gm)
		assert.Equal(t, cs.expectResult, res)
	}
}

func TestIsInValidatorList(t *testing.T) {
	ctx, am, vm, _ := setupTest(t, 0)
	user1 := createTestAccount(ctx, am, "user1")
	user2 := createTestAccount(ctx, am, "user2")
	user3 := createTestAccount(ctx, am, "user3")

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
		referenceList := &model.ValidatorReferenceList{
			AllValidators: cs.allValidators,
		}
		vm.storage.SetValidatorReferenceList(ctx, referenceList)
		res := vm.IsInValidatorList(ctx, cs.username)
		assert.Equal(t, cs.expectResult, res)
	}
}

func TestIsLegalVoterWithdraw(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	user1 := createTestAccount(ctx, am, "user1")
	voterMinDeposit, _ := gm.GetVoterMinDeposit(ctx)
	voterMinWithdraw, _ := gm.GetVoterMinWithdraw(ctx)
	vm.AddVoter(ctx, user1, voterMinDeposit.Plus(types.NewCoin(100*types.Decimals)), gm)

	cases := []struct {
		allValidators []types.AccountKey
		username      types.AccountKey
		withdraw      types.Coin
		expectResult  bool
	}{
		{[]types.AccountKey{}, user1, voterMinWithdraw.Minus(types.NewCoin(1 * types.Decimals)), false},
		{[]types.AccountKey{}, user1, voterMinWithdraw, true},
		{[]types.AccountKey{user1}, user1, voterMinWithdraw, false},
		{[]types.AccountKey{}, user1, types.NewCoin(100), false},
	}

	for _, cs := range cases {
		referenceList := &model.ValidatorReferenceList{
			AllValidators: cs.allValidators,
		}
		vm.storage.SetValidatorReferenceList(ctx, referenceList)
		res := vm.IsLegalVoterWithdraw(ctx, cs.username, cs.withdraw, gm)
		assert.Equal(t, cs.expectResult, res)
	}
}

func TestIsLegalDelegatorWithdraw(t *testing.T) {
	ctx, am, vm, gm := setupTest(t, 0)
	user1 := createTestAccount(ctx, am, "user1")
	user2 := createTestAccount(ctx, am, "user2")
	delegatorMinWithdraw, _ := gm.GetDelegatorMinWithdraw(ctx)
	voterMinDeposit, _ := gm.GetVoterMinDeposit(ctx)
	vm.AddVoter(ctx, user1, voterMinDeposit, gm)

	cases := []struct {
		addDelegation bool
		delegatedCoin types.Coin
		delegator     types.AccountKey
		voter         types.AccountKey
		withdraw      types.Coin
		expectResult  bool
	}{
		{false, types.NewCoin(0), user2, user1, delegatorMinWithdraw, false},
		{true, types.NewCoin(100 * types.Decimals), user2, user1, delegatorMinWithdraw, true},
		{false, types.NewCoin(0), user2, user1, types.NewCoin(0), false},
		{false, types.NewCoin(0), user2, user1, types.NewCoin(101 * types.Decimals), false},
	}

	for _, cs := range cases {
		if cs.addDelegation {
			vm.AddDelegation(ctx, cs.voter, cs.delegator, cs.delegatedCoin)
		}
		res := vm.IsLegalDelegatorWithdraw(ctx, cs.voter, cs.delegator, cs.withdraw, gm)
		assert.Equal(t, cs.expectResult, res)
	}
}
