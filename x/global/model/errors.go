package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrGlobalMetaNotFound - error if global meta not found in KVStore
func ErrGlobalMetaNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalMetaNotFound, fmt.Sprintf("global meta not found"))
}

// ErrInflationPoolNotFound - error if inflation pool is not found in KVStore
func ErrInflationPoolNotFound() sdk.Error {
	return types.NewError(types.CodeInflationPoolNotFound, fmt.Sprintf("inflation pool not found"))
}

// ErrGlobalConsumptionMetaNotFound - error if global consumption meta is not found in KVStore
func ErrGlobalConsumptionMetaNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalConsumptionMetaNotFound, fmt.Sprintf("global consumption meta not found"))
}

// ErrGlobalTPSNotFound - error if global tps is not found in KVStore
func ErrGlobalTPSNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalTPSNotFound, fmt.Sprintf("global tps not found"))
}

// ErrGlobalTimeNotFound - error if global time is not found in KVStore
func ErrGlobalTimeNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalTimeNotFound, fmt.Sprintf("global time not found"))
}

// ErrFailedToMarshalTimeEventList - error if marshal time event list failed
func ErrFailedToMarshalTimeEventList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalTimeEventList, fmt.Sprintf("failed to marshal time event list: %s", err.Error()))
}

// ErrFailedToMarshalGlobalMeta - error if marshal global meta failed
func ErrFailedToMarshalGlobalMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalGlobalMeta, fmt.Sprintf("failed to marshal global meta: %s", err.Error()))
}

// ErrFailedToMarshalInflationPool - error if marshal inflation pool failed
func ErrFailedToMarshalInflationPool(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalInflationPoll, fmt.Sprintf("failed to marshal inflation pool: %s", err.Error()))
}

// ErrFailedToMarshalConsumptionMeta - error if marshal consumption meta failed
func ErrFailedToMarshalConsumptionMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalConsumptionMeta, fmt.Sprintf("failed to marshal consumption meta: %s", err.Error()))
}

// ErrFailedToMarshalTPS - error if marshal tps failed
func ErrFailedToMarshalTPS(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalTPS, fmt.Sprintf("failed to marshal tps: %s", err.Error()))
}

// ErrFailedToMarshalTime - error if marshal time failed
func ErrFailedToMarshalTime(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalTime, fmt.Sprintf("failed to marshal time: %s", err.Error()))
}

// ErrFailedToUnmarshalTimeEventList - error if unmarshal time event list failed
func ErrFailedToUnmarshalTimeEventList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalTimeEventList, fmt.Sprintf("failed to unmarshal time event list: %s", err.Error()))
}

// ErrFailedToUnmarshalGlobalMeta - error if unmarshal global meta failed
func ErrFailedToUnmarshalGlobalMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalGlobalMeta, fmt.Sprintf("failed to unmarshal global meta: %s", err.Error()))
}

// ErrFailedToUnmarshalInflationPool - error if unmarshal inflation pool failed
func ErrFailedToUnmarshalInflationPool(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalInflationPool, fmt.Sprintf("failed to unmarshal inflation pool: %s", err.Error()))
}

// ErrFailedToUnmarshalConsumptionMeta - error if unmarshal consumption meta failed
func ErrFailedToUnmarshalConsumptionMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalConsumptionMeta, fmt.Sprintf("failed to unmarshal consumption meta: %s", err.Error()))
}

// ErrFailedToUnmarshalTPS - error if unmarshal tps failed
func ErrFailedToUnmarshalTPS(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalTPS, fmt.Sprintf("failed to unmarshal tps: %s", err.Error()))
}

// ErrFailedToUnmarshalTime - error if unmarshal time failed
func ErrFailedToUnmarshalTime(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalTime, fmt.Sprintf("failed to unmarshal time: %s", err.Error()))
}
