package types

const (
	// ModuleName is module name
	ModuleName = "price"

	// RouterKey is the message route for post
	RouterKey = ModuleName

	// QuerierRoute is the querier route for post
	QuerierRoute = ModuleName

	// query stores
	QueryPriceCurrent = "current"
	QueryPriceHistory = "history"
)
