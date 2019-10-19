package model

import (
	"github.com/lino-network/lino/types"
)

// GlobalTime - global time
type GlobalTime struct {
	ChainStartTime int64 `json:"chain_start_time"`
	LastBlockTime  int64 `json:"last_block_time"`
	PastMinutes    int64 `json:"past_minutes"`
}

type EventError = types.EventError
