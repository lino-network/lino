package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

// GlobalTime - global time
type GlobalTime struct {
	ChainStartTime int64 `json:"chain_start_time"`
	LastBlockTime  int64 `json:"last_block_time"`
	PastMinutes    int64 `json:"past_minutes"`
}

// EventError - event and errors
type EventError struct {
	Time    int64        `json:"time"`
	Event   types.Event  `json:"event"`
	ErrCode sdk.CodeType `json:"err_code"`
}
