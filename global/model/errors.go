package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrGlobalStorageGenesisInflationFailed() sdk.Error {
	return sdk.NewError(types.CodeGlobalStorageGenesisError, fmt.Sprintf("inflation allocation more than 100 percent"))
}

func ErrGlobalStorageGenesisFailed() sdk.Error {
	return sdk.NewError(types.CodeGlobalStorageGenesisError, fmt.Sprintf("GlobalStorage genesis failed"))
}

func ErrEventUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeGlobalStorageError, fmt.Sprintf("Event unmarshal error: %s", err.Error()))
}

func ErrEventMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeGlobalStorageError, fmt.Sprintf("Event marshal error: %s", err.Error()))
}

func ErrGlobalStatisticsNotFound() sdk.Error {
	return sdk.NewError(types.CodeGlobalStorageError, fmt.Sprintf("Global statistic not found"))
}

func ErrGlobalMetaNotFound() sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, fmt.Sprintf("Global meta not found"))
}

func ErrGlobalAllocationNotFound() sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, fmt.Sprintf("Global allocation not found"))
}

func ErrInfraAllocationNotFound() sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, fmt.Sprintf("Infra internal allocation not found"))
}

func ErrGlobalConsumptionMetaNotFound() sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, fmt.Sprintf("Global consumption meta not found"))
}

func ErrGlobalTPSNotFound() sdk.Error {
	return sdk.NewError(types.CodeGlobalManagerError, fmt.Sprintf("Global tps not found"))
}
