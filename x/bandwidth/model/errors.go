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

// ErrBlockInfoNotFound - error if last block info is not found
func ErrBlockInfoNotFound() sdk.Error {
	return types.NewError(types.CodeBlockInfoNotFound, fmt.Sprintf("last block info is not found"))
}
