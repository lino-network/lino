package model

import (
	crypto "github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
)

// AccountIR account related information when migrate, pk: Username
// CAN NOT use sdk.AccAddress, because there are illegal addresses in
// the state since upgrade-2, and the AccAddress.unmarshalJSON will return err.
type AccountIR struct {
	Username       types.AccountKey `json:"username"`
	CreatedAt      int64            `json:"created_at"`
	SigningKey     crypto.PubKey    `json:"signing_key"`
	TransactionKey crypto.PubKey    `json:"transaction_key"`
	Address        []byte           `json:"address"`
}

// AccountBankIR - user balance
type AccountBankIR struct {
	Address         []byte           `json:"address"` // pk
	Saving          types.Coin       `json:"saving"`
	FrozenMoneyList []FrozenMoneyIR  `json:"frozen_money_list"`
	PubKey          crypto.PubKey    `json:"public_key"`
	Sequence        uint64           `json:"sequence"`
	Username        types.AccountKey `json:"username"`
}

// FrozenMoneyIR - frozen money
type FrozenMoneyIR struct {
	Amount   types.Coin `json:"amount"`
	StartAt  int64      `json:"start_at"`
	Times    int64      `json:"times"`
	Interval int64      `json:"interval"`
}

type PermissionIR struct {
	Permission types.Permission `json:"permission"`
	CreatedAt  int64            `json:"created_at"`
	ExpiresAt  int64            `json:"expires_at"`
	Amount     types.Coin       `json:"amount"`
}

// GrantPubKeyIR also in account, no pk.
// must use AuthorizePermission for now. As this grant permission will be removed
// very soon, it is okay to let this happen.
type GrantPermissionIR struct {
	Username    types.AccountKey `json:"username"`
	GrantTo     types.AccountKey `json:"grant_to"`
	Permissions []PermissionIR   `json:"permissions"`
}

// AccountMetaIR - stores optional fields.
type AccountMetaIR struct {
	Username types.AccountKey `json:"username"`
	JSONMeta string           `json:"json_meta"`
}

// AccountTablesIR -
type AccountTablesIR struct {
	Version  int                 `json:"version"`
	Accounts []AccountIR         `json:"accounts"`
	Banks    []AccountBankIR     `json:"banks"`
	Metas    []AccountMetaIR     `json:"metas"`
	Grants   []GrantPermissionIR `json:"grants"`
}
