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

// BCEventErrReason - blockchain event error. deterministic.
type BCEventErr struct {
	ErrCode      sdk.CodeType      `json:"err_code"`
	ErrCodeSpace sdk.CodespaceType `json:"err_code_space"`
	Reason       string            `json:"module"`
}

func NewBCEventErr(err sdk.Error, reason string) *BCEventErr {
	return &BCEventErr{
		ErrCode:      err.Code(),
		ErrCodeSpace: err.Codespace(),
		Reason:       reason,
	}
}
