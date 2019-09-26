package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linotypes "github.com/lino-network/lino/types"
)

// Implements StakingHooks interface
var _ StakingHooks = VoteManager{}
var _ StakingHooks = MultiStakingHooks{}

// combine multiple staking hooks, all hook functions are run in array sequence
type MultiStakingHooks []StakingHooks

//go:generate mockery -name StakingHooks
// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	AfterAddingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error
	AfterSubtractingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error
}

// AfterAddingStake - call hook if registered
func (vm VoteManager) AfterAddingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	if vm.hooks != nil {
		return vm.hooks.AfterAddingStake(ctx, username)
	}
	return nil
}

// AfterSubtractingStake - call hook if registered
func (vm VoteManager) AfterSubtractingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	if vm.hooks != nil {
		return vm.hooks.AfterSubtractingStake(ctx, username)
	}
	return nil
}

func NewMultiStakingHooks(hooks ...StakingHooks) MultiStakingHooks {
	return hooks
}

// nolint
func (h MultiStakingHooks) AfterAddingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	for i := range h {
		if err := h[i].AfterAddingStake(ctx, username); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) AfterSubtractingStake(ctx sdk.Context, username linotypes.AccountKey) sdk.Error {
	for i := range h {
		if err := h[i].AfterSubtractingStake(ctx, username); err != nil {
			return err
		}
	}
	return nil
}
