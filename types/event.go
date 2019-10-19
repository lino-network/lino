package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event - event executed in app.go
type Event interface{}

// Minute -> TimeEventList
type TimeEventList struct {
	Events []Event `json:"events"`
}

// EventError - event and errors
// Module scheduled events
type EventError struct {
	Time    int64        `json:"time"`
	Event   Event        `json:"event"`
	ErrCode sdk.CodeType `json:"err_code"`
}

// EventExec is a function that can execute events.
type EventExec = func(ctx sdk.Context, event Event) sdk.Error

// BCEventErrReason - blockchain event error. deterministic.
type BCEventErr struct {
	Time         int64             `json:"time"`
	ErrCode      sdk.CodeType      `json:"err_code"`
	ErrCodeSpace sdk.CodespaceType `json:"err_code_space"`
	Reason       string            `json:"reason"`
}

func NewBCEventErr(ctx sdk.Context, err sdk.Error, reason string) BCEventErr {
	return BCEventErr{
		Time:         ctx.BlockTime().Unix(),
		ErrCode:      err.Code(),
		ErrCodeSpace: err.Codespace(),
		Reason:       reason,
	}
}

// BCEvent execute blockchain scheduled events.
type BCEventExec = func(ctx sdk.Context) []BCEventErr
