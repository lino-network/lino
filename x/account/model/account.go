package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	ttypes "github.com/tendermint/tendermint/types"

	"github.com/lino-network/lino/types"
)

// AccountInfo - user information
type AccountInfo struct {
	Username       types.AccountKey `json:"username"`
	CreatedAt      int64            `json:"created_at"`
	SigningKey     crypto.PubKey    `json:"signing_key"`
	TransactionKey crypto.PubKey    `json:"transaction_key"`
	Address        sdk.AccAddress   `json:"address"`
}

// AccountBank - user balance
type AccountBank struct {
	Saving   types.Coin       `json:"saving"`
	Pending  types.Coin       `json:"pending"`
	PubKey   crypto.PubKey    `json:"public_key"`
	Sequence uint64           `json:"sequence"`
	Username types.AccountKey `json:"username"`
}

// Pool - the pool for modules
type Pool struct {
	Name    types.PoolName `json:"name"`
	Balance types.Coin     `json:"balance"`
}

// Supply - stats of lino supply.
type Supply struct {
	LastYearTotal     types.Coin `json:"last_year_total"`
	Total             types.Coin `json:"total"`
	ChainStartTime    int64      `json:"chain_start_time"`
	LastInflationTime int64      `json:"last_inflation_time"`
}

// AccountMeta - stores optional fields.
type AccountMeta struct {
	JSONMeta string `json:"json_meta"`
}

type TxAndSequenceNumber struct {
	Username string       `json:"username"`
	Sequence uint64       `json:"sequence"`
	Tx       *Transaction `json:"tx"`
}

type Transaction struct {
	Hash   string    `json:"hash"`
	Height int64     `json:"height"`
	Tx     ttypes.Tx `json:"tx"`
	Code   uint32    `json:"code"`
	Log    string    `json:"log"`
}
