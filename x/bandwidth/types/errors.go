package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

// ErrInvalidMsgQuota - error when message fee is not valid
func ErrInvalidMsgQuota() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidMsgQuota, fmt.Sprintf("invalid message quota"))
}

// ErrInvalidExpectedMPS - error when message fee is not valid
func ErrInvalidExpectedMPS() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidExpectedMPS, fmt.Sprintf("invalid expected mps"))
}

// ErrAppBandwidthNotEnough - error when app bandwidth not enough
func ErrAppBandwidthNotEnough() sdk.Error {
	return linotypes.NewError(linotypes.CodeAppBandwidthNotEnough, fmt.Sprintf("app bandwidth not enough"))
}

// ErrUserMsgFeeNotEnough - error when app bandwidth not enough
func ErrUserMsgFeeNotEnough() sdk.Error {
	return linotypes.NewError(linotypes.CodeUserMsgFeeNotEnough, fmt.Sprintf("user message fee not enough"))
}
