package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrGlobalMetaNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalMetaNotFound, fmt.Sprintf("global meta not found"))
}

func ErrInflationPoolNotFound() sdk.Error {
	return types.NewError(types.CodeInflationPoolNotFound, fmt.Sprintf("inflation pool not found"))
}

func ErrGlobalConsumptionMetaNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalConsumptionMetaNotFound, fmt.Sprintf("global consumption meta not found"))
}

func ErrGlobalTPSNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalTPSNotFound, fmt.Sprintf("global tps not found"))
}

func ErrGlobalTimeNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalTimeNotFound, fmt.Sprintf("global time not found"))
}

func ErrInfraInflationCoinConversion() sdk.Error {
	return types.NewError(types.CodeInfraInflationCoinConversion, fmt.Sprintf("failed to convert infra inflation coin"))
}

func ErrContentCreatorCoinConversion() sdk.Error {
	return types.NewError(types.CodeContentCreatorCoinConversion, fmt.Sprintf("failed to convert content creator coin"))
}

func ErrDeveloperCoinConversion() sdk.Error {
	return types.NewError(types.CodeDeveloperCoinConversion, fmt.Sprintf("failed to convert developer coin"))
}

func ErrValidatorCoinConversion() sdk.Error {
	return types.NewError(types.CodeValidatorCoinConversion, fmt.Sprintf("failed to convert validator coin"))
}

// marshal error
func ErrFailedToMarshalTimeEventList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalTimeEventList, fmt.Sprintf("failed to marshal time event list: %s", err.Error()))
}

func ErrFailedToMarshalGlobalMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalGlobalMeta, fmt.Sprintf("failed to marshal global meta: %s", err.Error()))
}

func ErrFailedToMarshalInflationPool(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalInflationPoll, fmt.Sprintf("failed to marshal inflation pool: %s", err.Error()))
}

func ErrFailedToMarshalConsumptionMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalConsumptionMeta, fmt.Sprintf("failed to marshal consumption meta: %s", err.Error()))
}

func ErrFailedToMarshalTPS(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalTPS, fmt.Sprintf("failed to marshal tps: %s", err.Error()))
}

func ErrFailedToMarshalTime(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalTime, fmt.Sprintf("failed to marshal time: %s", err.Error()))
}

// unmarshal error
func ErrFailedToUnmarshalTimeEventList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalTimeEventList, fmt.Sprintf("failed to unmarshal time event list: %s", err.Error()))
}

func ErrFailedToUnmarshalGlobalMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalGlobalMeta, fmt.Sprintf("failed to unmarshal global meta: %s", err.Error()))
}

func ErrFailedToUnmarshalInflationPool(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalInflationPool, fmt.Sprintf("failed to unmarshal inflation pool: %s", err.Error()))
}

func ErrFailedToUnmarshalConsumptionMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalConsumptionMeta, fmt.Sprintf("failed to unmarshal consumption meta: %s", err.Error()))
}

func ErrFailedToUnmarshalTPS(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalTPS, fmt.Sprintf("failed to unmarshal tps: %s", err.Error()))
}

func ErrFailedToUnmarshalTime(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalTime, fmt.Sprintf("failed to unmarshal time: %s", err.Error()))
}
