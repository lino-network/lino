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

// ErrBlockInfoNotFound - error if last block info is not found
func ErrBlockInfoNotFound() sdk.Error {
	return types.NewError(types.CodeBlockInfoNotFound, fmt.Sprintf("last block info is not found"))
}

// ErrFailedToMarshalBlockInfo - error if marshal Last block info failed
func ErrFailedToMarshalBlockInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalBlockInfo, fmt.Sprintf("failed to marshal last block info: %s", err.Error()))
}

// ErrFailedToUnmarshalBlockInfo - error if unmarshal Last block info failed
func ErrFailedToUnmarshalBlockInfo(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalBlockInfo, fmt.Sprintf("failed to unmarshal last block info: %s", err.Error()))
}
