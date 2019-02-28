package global

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrRegisterExpiredEvent - error when register event time is before current timestamp
func ErrRegisterExpiredEvent(unixTime int64) sdk.Error {
	return types.NewError(types.CodeRegisterExpiredEvent, fmt.Sprintf("register event at expired time %v", unixTime))
}

// ErrGetPastDay - error if get past day is negative
func ErrGetPastDay() sdk.Error {
	return types.NewError(types.CodeFailedToGetAmountOfConsumptionExponent, "get past day failed")
}

// ErrParseEventCacheList - error if parse event cache list failed
func ErrParseEventCacheList() sdk.Error {
	return types.NewError(types.CodeFailedToParseEventCacheList, "parse event list failed")
}

// ErrQueryFailed - error when query global store failed
func ErrQueryFailed() sdk.Error {
	return types.NewError(types.CodeGlobalQueryFailed, fmt.Sprintf("query global store failed"))
}
