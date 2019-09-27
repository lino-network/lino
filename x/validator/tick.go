package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker - execute before every block, update signing info and record validator set
func BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock, vm ValidatorKeeper) {
	// update preblock validators
	validatorList, err := vm.GetValidatorList(ctx)
	if err != nil {
		panic(err)
	}
	vals, err := vm.GetCommittingValidators(ctx)
	if err != nil {
		panic(err)
	}
	validatorList.PreBlockValidators = vals
	if err := vm.SetValidatorList(ctx, validatorList); err != nil {
		panic(err)
	}

	// update signing stats.
	updateErr := vm.UpdateSigningStats(ctx, req.LastCommitInfo.Votes)
	if updateErr != nil {
		panic(updateErr)
	}

	if err := vm.FireIncompetentValidator(ctx, req.ByzantineValidators); err != nil {
		panic(err)
	}
	return
}
