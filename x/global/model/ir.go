package model

import (
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

// GlobalTimeIR - global time
type GlobalTimeIR struct {
	ChainStartTime int64 `json:"chain_start_time"`
	LastBlockTime  int64 `json:"last_block_time"`
	PastMinutes    int64 `json:"past_minutes"`
}

// GlobalTimeEventsIR - events, pk: UnixTime
type GlobalTimeEventsIR struct {
	UnixTime      int64               `json:"unix_time"`
	TimeEventList types.TimeEventList `json:"time_event_list"`
}

// GlobalTablesIR - state
type GlobalTablesIR struct {
	Version              int                  `json:"version"`
	GlobalTimeEventLists []GlobalTimeEventsIR `json:"global_time_event_lists"`
	Time                 GlobalTimeIR         `json:"time"`
}
