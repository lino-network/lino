package global

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrGlobalManagerRegisterEventAtTime(unixTime int64) sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, fmt.Sprintf("tegister event at time %v failed", unixTime))
}

func ErrGlobalManagerRegisterExpiredEvent(unixTime int64) sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, fmt.Sprintf("tegister event at expired time %v", unixTime))
}

func ErrAddConsumptionFrictionToRewardPool() sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, "add consumption friction to reward pool failed")
}

func ErrGetRewardAndPopFromWindow() sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, "get reward from consumption pool and pop from window failed")
}
