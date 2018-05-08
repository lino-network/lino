package param

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ChangeGlobalAllocationParamEvent struct {
	Param GlobalAllocationParam `json:"param"`
}

func (event ChangeGlobalAllocationParamEvent) Execute(ctx sdk.Context, ph ParamHolder) sdk.Error {
	if err := ph.setGlobalAllocationParam(ctx, &event.Param); err != nil {
		return err
	}
	return nil
}
