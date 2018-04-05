package global

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrGlobalManagerRegisterEventAtHeight(height int64) sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, fmt.Sprintf("Register event at height %v", height))
}

func ErrGlobalManagerRegisterEventAtTime(unixTime int64) sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, fmt.Sprintf("Register event at time %v", unixTime))
}

func ErrAddConsumptionFrictionToRewardPool() sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, "add consumption friction to reward pool failed")
}

func ErrGetRewardAndPopFromWindow() sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, "get reward from consumption pool and pop from window failed")
}
