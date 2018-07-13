package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// not found error
func ErrAccountInfoNotFound() sdk.Error {
	return types.NewError(types.CodeAccountInfoNotFound, fmt.Sprintf("account info is not found"))
}

func ErrAccountBankNotFound() sdk.Error {
	return types.NewError(types.CodeAccountBankNotFound, fmt.Sprintf("account bank is not found"))
}

func ErrAccountMetaNotFound() sdk.Error {
	return types.NewError(types.CodeAccountMetaNotFound, fmt.Sprintf("account meta is not found"))
}

func ErrRewardNotFound() sdk.Error {
	return types.NewError(types.CodeRewardNotFound, fmt.Sprintf("reward is not found"))
}

func ErrPendingStakeQueueNotFound() sdk.Error {
	return types.NewError(types.CodePendingStakeQueueNotFound, fmt.Sprintf("pending stake queue is not found"))
}

func ErrGrantPubKeyNotFound() sdk.Error {
	return types.NewError(types.CodeGrantPubKeyNotFound, fmt.Sprintf("grant public key is not found"))
}

// marshal error
func ErrFailedToMarshalAccountInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalAccountInfo, fmt.Sprintf("failed to marshal account info: %s", err.Error()))
}

func ErrFailedToMarshalAccountBank(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalAccountBank, fmt.Sprintf("failed to marshal account bank: %s", err.Error()))
}

func ErrFailedToMarshalAccountMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalAccountMeta, fmt.Sprintf("failed to marshal account meta: %s", err.Error()))
}

func ErrFailedToMarshalFollowerMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalFollowerMeta, fmt.Sprintf("failed to marshal follower meta: %s", err.Error()))
}

func ErrFailedToMarshalFollowingMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalFollowingMeta, fmt.Sprintf("failed to marshal following meta: %s", err.Error()))
}

func ErrFailedToMarshalReward(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalReward, fmt.Sprintf("failed to marshal reward: %s", err.Error()))
}

func ErrFailedToMarshalPendingStakeQueue(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPendingStakeQueue, fmt.Sprintf("failed to marshal pending stake queue: %s", err.Error()))
}

func ErrFailedToMarshalGrantPubKey(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalGrantPubKey, fmt.Sprintf("failed to marshal grant pub key: %s", err.Error()))
}

func ErrFailedToMarshalRelationship(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalRelationship, fmt.Sprintf("failed to marshal relationship: %s", err.Error()))
}

func ErrFailedToMarshalBalanceHistory(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalBalanceHistory, fmt.Sprintf("failed to marshal balance history: %s", err.Error()))
}

func ErrFailedToMarshalRewardHistory(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalRewardHistory, fmt.Sprintf("failed to marshal reward history: %s", err.Error()))
}

// unmarshal error
func ErrFailedToUnmarshalAccountInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalAccountInfo, fmt.Sprintf("failed to unmarshal account info: %s", err.Error()))
}

func ErrFailedToUnmarshalAccountBank(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalAccountBank, fmt.Sprintf("failed to unmarshal account bank: %s", err.Error()))
}

func ErrFailedToUnmarshalAccountMeta(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalAccountMeta, fmt.Sprintf("failed to unmarshal account meta: %s", err.Error()))
}

func ErrFailedToUnmarshalReward(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalReward, fmt.Sprintf("failed to unmarshal reward: %s", err.Error()))
}

func ErrFailedToUnmarshalPendingStakeQueue(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPendingStakeQueue, fmt.Sprintf("failed to unmarshal pending stake queue: %s", err.Error()))
}

func ErrFailedToUnmarshalGrantPubKey(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalGrantPubKey, fmt.Sprintf("failed to unmarshal grant pub key: %s", err.Error()))
}

func ErrFailedToUnmarshalRelationship(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalRelationship, fmt.Sprintf("failed to unmarshal relationship: %s", err.Error()))
}

func ErrFailedToUnmarshalBalanceHistory(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalBalanceHistory, fmt.Sprintf("failed to unmarshal balance history: %s", err.Error()))
}

func ErrFailedToUnmarshalRewardHistory(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalRewardHistory, fmt.Sprintf("failed to unmarshal reward history: %s", err.Error()))
}
