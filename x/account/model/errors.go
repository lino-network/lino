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

// ErrRewardNotFound - error if reward is not found
func ErrRewardNotFound() sdk.Error {
	return types.NewError(types.CodeRewardNotFound, fmt.Sprintf("reward is not found"))
}

// ErrPendingStakeQueueNotFound - error if pending stake queue is not found
func ErrPendingStakeQueueNotFound() sdk.Error {
	return types.NewError(types.CodePendingStakeQueueNotFound, fmt.Sprintf("pending stake queue is not found"))
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

// ErrFailedToMarshalFollowerMeta - error if marshal follower meta failed
func ErrFailedToMarshalFollowerMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalFollowerMeta, fmt.Sprintf("failed to marshal follower meta: %s", err.Error()))
}

// ErrFailedToMarshalFollowingMeta - error if marshal following meta failed
func ErrFailedToMarshalFollowingMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalFollowingMeta, fmt.Sprintf("failed to marshal following meta: %s", err.Error()))
}

// ErrFailedToMarshalReward - error if marshal reward failed
func ErrFailedToMarshalReward(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalReward, fmt.Sprintf("failed to marshal reward: %s", err.Error()))
}

// ErrFailedToMarshalPendingStakeQueue - error if marshal pending stake queue failed
func ErrFailedToMarshalPendingStakeQueue(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPendingStakeQueue, fmt.Sprintf("failed to marshal pending stake queue: %s", err.Error()))
}

// ErrFailedToMarshalGrantPubKey - error if marshal grant public key failed
func ErrFailedToMarshalGrantPubKey(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalGrantPubKey, fmt.Sprintf("failed to marshal grant pub key: %s", err.Error()))
}

// ErrFailedToMarshalRelationship - error if marshal relationship failed
func ErrFailedToMarshalRelationship(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalRelationship, fmt.Sprintf("failed to marshal relationship: %s", err.Error()))
}

// ErrFailedToMarshalBalanceHistory - error if marshal balance history failed
func ErrFailedToMarshalBalanceHistory(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalBalanceHistory, fmt.Sprintf("failed to marshal balance history: %s", err.Error()))
}

// ErrFailedToMarshalRewardHistory - error if marshal reward history failed
func ErrFailedToMarshalRewardHistory(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalRewardHistory, fmt.Sprintf("failed to marshal reward history: %s", err.Error()))
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

// ErrFailedToUnmarshalPendingStakeQueue - error if unmarshal pending stake queue failed
func ErrFailedToUnmarshalPendingStakeQueue(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPendingStakeQueue, fmt.Sprintf("failed to unmarshal pending stake queue: %s", err.Error()))
}

// ErrFailedToUnmarshalGrantPubKey - error if unmarshal grant public key failed
func ErrFailedToUnmarshalGrantPubKey(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalGrantPubKey, fmt.Sprintf("failed to unmarshal grant pub key: %s", err.Error()))
}

// ErrFailedToUnmarshalRelationship - error if unmarshal relationship failed
func ErrFailedToUnmarshalRelationship(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalRelationship, fmt.Sprintf("failed to unmarshal relationship: %s", err.Error()))
}

// ErrFailedToUnmarshalBalanceHistory - error if unmarshal balance history failed
func ErrFailedToUnmarshalBalanceHistory(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalBalanceHistory, fmt.Sprintf("failed to unmarshal balance history: %s", err.Error()))
}

// ErrFailedToUnmarshalRewardHistory - error if unmarshal reward history failed
func ErrFailedToUnmarshalRewardHistory(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalRewardHistory, fmt.Sprintf("failed to unmarshal reward history: %s", err.Error()))
}
