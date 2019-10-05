package model

import (
	types "github.com/lino-network/lino/types"
)

// DeveloperIR - apps
type DeveloperIR struct {
	Username       types.AccountKey `json:"username"` // pk
	AppConsumption types.MiniDollar `json:"app_consumption"`
	Website        string           `json:"web_site"`
	Description    string           `json:"description"`
	AppMetaData    string           `json:"app_meta_data"`
	IsDeleted      bool             `json:"is_deleted"`
	NAffiliated    int64            `json:"n_affiliated"`
}

// AppIDAIR - app issued IDA.
type AppIDAIR struct {
	App             types.AccountKey `json:"app"` // pk
	Name            string           `json:"name"`
	MiniIDAPrice    types.MiniDollar `json:"price"`
	IsRevoked       bool             `json:"is_revoked"`
	RevokeCoinPrice types.MiniDollar `json:"revoke_coin_price"` // the price of one coin upon revoke.
}

type IDABankIR struct {
	App      types.AccountKey `json:"app"`  // pk
	User     types.AccountKey `json:"user"` // pk
	Balance  types.MiniDollar `json:"b"`
	Unauthed bool             `json:"unauthed,omitempty"`
}

type ReservePoolIR struct {
	Total           types.Coin       `json:"total"`
	TotalMiniDollar types.MiniDollar `json:"total_minidollar"`
}

type AffiliatedAccIR struct {
	App  types.AccountKey `json:"app"`  // pk
	User types.AccountKey `json:"user"` // pk
}

type UserRoleIR struct {
	User          types.AccountKey `json:"user"` // pk
	AffiliatedApp types.AccountKey `json:"aa"`
}

type IDAStatsIR struct {
	App   types.AccountKey `json:"app"` // pk
	Total types.MiniDollar `json:"total"`
}

// DeveloperTablesIR
type DeveloperTablesIR struct {
	Version        int               `json:"version"`
	Developers     []DeveloperIR     `json:"developers"`
	IDAs           []AppIDAIR        `json:"id_as"`
	IDABanks       []IDABankIR       `json:"ida_banks"`
	ReservePool    ReservePoolIR     `json:"reserve_pool"`
	AffiliatedAccs []AffiliatedAccIR `json:"affiliated_accs"`
	UserRoles      []UserRoleIR      `json:"user_roles"`
	IDAStats       []IDAStatsIR      `json:"ida_stats"`
}
