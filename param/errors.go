package param

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrParamHolderGenesisFailed() sdk.Error {
	return sdk.NewError(types.CodeParamHolderGenesisError, fmt.Sprintf("Param holder genesis failed"))
}

func ErrDeveloperParamNotFound() sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Developer param not found"))
}

func ErrValidatorParamNotFound() sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Validator param not found"))
}

func ErrCoinDayParamNotFound() sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Coin day param not found"))
}

func ErrBandwidthParamNotFound() sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Bandwidth param not found"))
}

func ErrVoteParamNotFound() sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Vote param not found"))
}

func ErrProposalParamNotFound() sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Proposal param not found"))
}

func ErrGlobalAllocationParamNotFound() sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Global allocation param not found"))
}

func ErrInfraAllocationParamNotFound() sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Infra internal allocation param not found"))
}

func ErrEvaluateOfContentValueParamNotFound() sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Evaluate of content value param not found"))
}

func ErrEventUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Event unmarshal error: %s", err.Error()))
}

func ErrEventMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeParamStoreError, fmt.Sprintf("Event marshal error: %s", err.Error()))
}
