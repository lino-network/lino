package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ErrBandwidthInfoNotFound - error if bandwidth info is not found
func ErrBandwidthInfoNotFound() sdk.Error {
	return types.NewError(types.CodeBandwidthInfoNotFound, fmt.Sprintf("bandwidth info is not found"))
}

// ErrFailedToMarshalBandwidthInfo - error if marshal bandwidth info failed
func ErrFailedToMarshalBandwidthInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalBandwidthInfo, fmt.Sprintf("failed to marshal bandwidth info: %s", err.Error()))
}

// ErrFailedToUnmarshalBandwidthInfo - error if unmarshal bandwidth info failed
func ErrFailedToUnmarshalBandwidthInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalBandwidthInfo, fmt.Sprintf("failed to unmarshal bandwidth info: %s", err.Error()))
}

// ErrCurBlockInfoNotFound - error if cur block info is not found
func ErrCurBlockInfoNotFound() sdk.Error {
	return types.NewError(types.CodeCurBlockInfoNotFound, fmt.Sprintf("cur block info is not found"))
}

// ErrFailedToMarshalCurBlockInfo - error if marshal cur block info failed
func ErrFailedToMarshalCurBlockInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalCurBlockInfo, fmt.Sprintf("failed to marshal cur block info: %s", err.Error()))
}

// ErrFailedToUnmarshalCurBlockInfo - error if unmarshal cur block info failed
func ErrFailedToUnmarshalCurBlockInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalCurBlockInfo, fmt.Sprintf("failed to unmarshal cur block info: %s", err.Error()))
}
