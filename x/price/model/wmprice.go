package model

import (
	linotypes "github.com/lino-network/lino/types"
)

// FedPrice is record of price fed by validators.
type FedPrice struct {
	Validator linotypes.AccountKey `json:"validator"`
	Price     linotypes.MiniDollar `json:"price"`
	UpdateAt  int64                `json:"update_at"`
}

// TimePrice is time + price
type TimePrice struct {
	Price    linotypes.MiniDollar `json:"price"`
	UpdateAt int64                `json:"update_at"`
}

// FedRecord - power and price.
type FedRecord struct {
	Validator linotypes.AccountKey `json:"validator"`
	Price     linotypes.MiniDollar `json:"price"`
	Power     linotypes.Coin       `json:"power"`
	UpdateAt  int64                `json:"update_at"`
}

// FeedHistory the history of price feed of one price update.
// Used by querier for now and governance in future.
type FeedHistory struct {
	Price    linotypes.MiniDollar `json:"price"`
	Feeded   []FedRecord          `json:"feeded"`
	UpdateAt int64                `json:"update_at"`
}
