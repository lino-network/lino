package types

const (
	// ModuleKey is the name of the module
	ModuleName = "account"

	// RouterKey is the message route for account
	RouterKey = ModuleName

	// QuerierRoute is the querier route for account
	QuerierRoute = ModuleName

	// query pathes
	QueryAccountInfo            = "info"
	QueryAccountBank            = "bank"
	QueryAccountMeta            = "meta"
	QueryAccountPendingCoinDay  = "pendingCoinDay"
	QueryAccountGrantPubKeys    = "grantPubKey"
	QueryAccountAllGrantPubKeys = "allGrantPubKey"
	QueryTxAndAccountSequence   = "txAndSeq"
)
