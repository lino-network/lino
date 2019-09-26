package model

import (
	"time"

	linotypes "github.com/lino-network/lino/types"
)

// FedPrice is record of price fed by validators.
type FedPrice struct {
	Validator linotypes.AccountKey `json:"validator"`
	Price     linotypes.MiniDollar `json:"price"`
	FedTime   time.Time            `json:"fed_time"`
}

// TimePrice is time + price
type TimePrice struct {
	Price    linotypes.MiniDollar `json:"price"`
	UpdateAt time.Time            `json:"update_at"`
}
