package types

const (
	// ModuleName is module name
	ModuleName = "bandwidth"

	// RouterKey is the message route for bandwidth
	RouterKey = ModuleName

	// QuerierRoute is the querier route for bandwidth
	QuerierRoute = ModuleName

	// querier paths.
	QueryBandwidthInfo    = "bandwidthinfo"
	QueryBlockInfo        = "blockinfo"
	QueryAppBandwidthInfo = "appinfo"
)
