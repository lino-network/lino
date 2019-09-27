package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	votemn "github.com/lino-network/lino/x/vote/manager"
)

func (vm ValidatorManager) AfterAddingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	return vm.onStakeChange(ctx, username)
}

func (vm ValidatorManager) AfterSubtractingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	return vm.onStakeChange(ctx, username)
}

type Hooks struct {
	vm ValidatorManager
}

var _ votemn.StakingHooks = Hooks{}

// Return the wrapper struct
func (vm ValidatorManager) Hooks() Hooks {
	return Hooks{vm}
}

func (h Hooks) AfterAddingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	return h.vm.AfterAddingStake(ctx, username)
}

func (h Hooks) AfterSubtractingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	return h.vm.AfterSubtractingStake(ctx, username)
}
