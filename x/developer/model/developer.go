package model

import (
	types "github.com/lino-network/lino/types"
)

// Developer - developer is account with developer deposit, can get developer inflation
type Developer struct {
	Username types.AccountKey `json:"username"`
	// Deprecated field, use stakein amount as requirement since upgrade2.
	Deposit        types.Coin       `json:"deposit"`
	AppConsumption types.MiniDollar `json:"app_consumption"`
	Website        string           `json:"web_site"`
	Description    string           `json:"description"`
	AppMetaData    string           `json:"app_meta_data"`
	IsDeleted      bool             `json:"is_deleted"`
	NAffiliated    int64            `json:"n_affiliated"`
}

// AppIDA - app issued IDA.
type AppIDA struct {
	App             types.AccountKey `json:"app"`
	Name            string           `json:"name"`
	MiniIDAPrice    types.MiniDollar `json:"mini_ida_price"`
	IsRevoked       bool             `json:"is_revoked"`
	RevokeCoinPrice types.MiniDollar `json:"revoke_coin_price"` // the price of one coin upon revoke.
}

// AppIDAStats - app ida stats
type AppIDAStats struct {
	Total types.MiniDollar `json:"total"`
}

// Role - User Role
type Role struct {
	AffiliatedApp types.AccountKey `json:"aa"`
}

// IDABank - IDA's bank
type IDABank struct {
	Balance  types.MiniDollar `json:"b"`
	Unauthed bool             `json:"unauthed,omitempty"`
}

type ReservePool struct {
	Total           types.Coin       `json:"total"`
	TotalMiniDollar types.MiniDollar `json:"total_minidollar"`
}
