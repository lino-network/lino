package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
)

func BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock, vm ValidatorManager) (panelty types.Coin) {
	validatorList, err := vm.GetValidatorList(ctx)
	if err != nil {
		panic(err)
	}
	validatorList.PreBlockValidators = validatorList.OncallValidators
	if err := vm.SetValidatorList(ctx, validatorList); err != nil {
		panic(err)
	}

	vm.UpdateSigningValidator(ctx, req.Validators)

	panelty, _ = vm.FireIncompetentValidator(ctx, req.ByzantineValidators)
	return
}
