package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ErrAccountInfoNotFound - error if account info is not found
func ErrAccountInfoNotFound() sdk.Error {
	return types.NewError(types.CodeAccountInfoNotFound, fmt.Sprintf("account info is not found"))
}

// ErrAccountBankNotFound - error if account bank is not found
func ErrAccountBankNotFound() sdk.Error {
	return types.NewError(types.CodeAccountBankNotFound, fmt.Sprintf("account bank is not found"))
}

// ErrAccountMetaNotFound - error if account meta is not found
func ErrAccountMetaNotFound() sdk.Error {
	return types.NewError(types.CodeAccountMetaNotFound, fmt.Sprintf("account meta is not found"))
}

// ErrPoolNotFound - error if pool is not found
func ErrPoolNotFound() sdk.Error {
	return types.NewError(types.CodePoolNotFound, fmt.Sprintf("pool is not found"))
}

// ErrRewardNotFound - error if reward is not found
func ErrRewardNotFound() sdk.Error {
	return types.NewError(types.CodeRewardNotFound, fmt.Sprintf("reward is not found"))
}

// ErrPendingCoinDayQueueNotFound - error if pending coin day queue is not found
func ErrPendingCoinDayQueueNotFound() sdk.Error {
	return types.NewError(types.CodePendingCoinDayQueueNotFound, fmt.Sprintf("pending coin day queue is not found"))
}

// ErrGrantPubKeyNotFound - error if grant public key is not found
func ErrGrantPubKeyNotFound() sdk.Error {
	return types.NewError(types.CodeGrantPubKeyNotFound, fmt.Sprintf("grant public key is not found"))
}

// ErrFailedToMarshalAccountInfo - error if marshal account info failed
func ErrFailedToMarshalAccountInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalAccountInfo, fmt.Sprintf("failed to marshal account info: %s", err.Error()))
}

// ErrFailedToMarshalAccountBank - error if marshal account bank failed
func ErrFailedToMarshalAccountBank(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalAccountBank, fmt.Sprintf("failed to marshal account bank: %s", err.Error()))
}

// ErrFailedToMarshalAccountMeta - error if marshal account meta failed
func ErrFailedToMarshalAccountMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalAccountMeta, fmt.Sprintf("failed to marshal account meta: %s", err.Error()))
}

// ErrFailedToMarshalReward - error if marshal reward failed
func ErrFailedToMarshalReward(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalReward, fmt.Sprintf("failed to marshal reward: %s", err.Error()))
}

// ErrFailedToMarshalPendingCoinDayQueue - error if marshal pending coin day queue failed
func ErrFailedToMarshalPendingCoinDayQueue(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPendingCoinDayQueue, fmt.Sprintf("failed to marshal pending coin day queue: %s", err.Error()))
}

// ErrFailedToMarshalGrantPubKey - error if marshal grant public key failed
func ErrFailedToMarshalGrantPubKey(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalGrantPubKey, fmt.Sprintf("failed to marshal grant pub key: %s", err.Error()))
}

// ErrFailedToUnmarshalAccountInfo - error if unmarshal account info failed
func ErrFailedToUnmarshalAccountInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalAccountInfo, fmt.Sprintf("failed to unmarshal account info: %s", err.Error()))
}

// ErrFailedToUnmarshalAccountBank - error if unmarshal account bank failed
func ErrFailedToUnmarshalAccountBank(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalAccountBank, fmt.Sprintf("failed to unmarshal account bank: %s", err.Error()))
}

// ErrFailedToUnmarshalAccountMeta - error if unmarshal account meta failed
func ErrFailedToUnmarshalAccountMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalAccountMeta, fmt.Sprintf("failed to unmarshal account meta: %s", err.Error()))
}

// ErrFailedToUnmarshalReward - error if unmarshal account reward failed
func ErrFailedToUnmarshalReward(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalReward, fmt.Sprintf("failed to unmarshal reward: %s", err.Error()))
}

// ErrFailedToUnmarshalPendingCoinDayQueue - error if unmarshal pending coin day queue failed
func ErrFailedToUnmarshalPendingCoinDayQueue(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPendingCoinDayQueue, fmt.Sprintf("failed to unmarshal pending coin day queue: %s", err.Error()))
}

// ErrFailedToUnmarshalGrantPubKey - error if unmarshal grant public key failed
func ErrFailedToUnmarshalGrantPubKey(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalGrantPubKey, fmt.Sprintf("failed to unmarshal grant pub key: %s", err.Error()))
}

