package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker - execute before every block, update signing info and record validator set
func BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock, vm ValidatorKeeper) {
	// update preblock validators
	validatorList := vm.GetValidatorList(ctx)
	vals := vm.GetCommittingValidators(ctx)
	validatorList.PreBlockValidators = vals
	vm.SetValidatorList(ctx, validatorList)

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
