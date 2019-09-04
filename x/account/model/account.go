package model

import (
	"github.com/lino-network/lino/types"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ttypes "github.com/tendermint/tendermint/types"
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
	Saving          types.Coin       `json:"saving"`
	FrozenMoneyList []FrozenMoney    `json:"frozen_money_list"`
	PubKey          crypto.PubKey    `json:"public_key"`
	Sequence        uint64           `json:"sequence"`
	Username        types.AccountKey `json:"username"`
}

// FrozenMoney - frozen money
type FrozenMoney struct {
	Amount   types.Coin `json:"amount"`
	StartAt  int64      `json:"start_at"`
	Times    int64      `json:"times"`
	Interval int64      `json:"interval"`
}

// GrantPermission - user grant permission to a user with a certain permission
type GrantPermission struct {
	GrantTo    types.AccountKey `json:"grant_to"`
	Permission types.Permission `json:"permission"`
	CreatedAt  int64            `json:"created_at"`
	ExpiresAt  int64            `json:"expires_at"`
	Amount     types.Coin       `json:"amount"`
}

// AccountMeta - stores tiny and frequently updated fields.
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
