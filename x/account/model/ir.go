package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
)

// AccountIR account related information when migrate, pk: Username
type AccountIR struct {
	Username       types.AccountKey `json:"username"`
	CreatedAt      int64            `json:"created_at"`
	SigningKey     crypto.PubKey    `json:"signing_key"`
	TransactionKey crypto.PubKey    `json:"transaction_key"`
	Address        sdk.AccAddress   `json:"address"`
}

// AccountBankIR - user balance
type AccountBankIR struct {
	Address  []byte           `json:"address"` // pk
	Saving   types.Coin       `json:"saving"`
	Pending  types.Coin       `json:"pending"`
	PubKey   crypto.PubKey    `json:"public_key"`
	Sequence uint64           `json:"sequence"`
	Username types.AccountKey `json:"username"`
}

// AccountMetaIR - stores optional fields.
type AccountMetaIR struct {
	Username types.AccountKey `json:"username"`
	JSONMeta string           `json:"json_meta"`
}

// PoolIR - the module account.
type PoolIR struct {
	Name    types.PoolName `json:"name"`
	Balance types.Coin     `json:"balance"`
}

// SupplyIR - stats of lino supply.
type SupplyIR struct {
	LastYearTotal     types.Coin `json:"last_year_total"`
	Total             types.Coin `json:"total"`
	ChainStartTime    int64      `json:"chain_start_time"`
	LastInflationTime int64      `json:"last_inflation_time"`
}

// AccountTablesIR -
type AccountTablesIR struct {
	Version  int             `json:"version"`
	Accounts []AccountIR     `json:"accounts"`
	Banks    []AccountBankIR `json:"banks"`
	Metas    []AccountMetaIR `json:"metas"`
	Pools    []PoolIR        `json:"pools"`
	Supply   SupplyIR        `json:"supply"`
}
