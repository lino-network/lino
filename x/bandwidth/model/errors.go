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

// ErrLastBlockInfoNotFound - error if last block info is not found
func ErrLastBlockInfoNotFound() sdk.Error {
	return types.NewError(types.CodeLastBlockInfoNotFound, fmt.Sprintf("last block info is not found"))
}

// ErrFailedToMarshalLastBlockInfo - error if marshal Last block info failed
func ErrFailedToMarshalLastBlockInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalLastBlockInfo, fmt.Sprintf("failed to marshal last block info: %s", err.Error()))
}

// ErrFailedToUnmarshalLastBlockInfo - error if unmarshal Last block info failed
func ErrFailedToUnmarshalLastBlockInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalLastBlockInfo, fmt.Sprintf("failed to unmarshal last block info: %s", err.Error()))
}
