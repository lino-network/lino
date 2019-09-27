package model

import (
	linotypes "github.com/lino-network/lino/types"
)

// FedPrice is record of price fed by validators.
type FedPrice struct {
	Validator linotypes.AccountKey `json:"validator"`
	Price     linotypes.MiniDollar `json:"price"`
	FedTime   int64                `json:"fed_time"`
}

// TimePrice is time + price
type TimePrice struct {
	Price    linotypes.MiniDollar `json:"price"`
	UpdateAt int64                `json:"update_at"`
}
