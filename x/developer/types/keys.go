package types

const (
	// ModuleKey is the name of the module
	ModuleName = "developer"

	// RouterKey is the message route for developer
	RouterKey = ModuleName

	// QuerierRoute is the querier route for developer
	QuerierRoute = ModuleName

	// Query sub spaces.
	QueryDeveloper     = "dev"
	QueryDeveloperList = "devList"
	QueryIDA           = "devIDA"
	QueryIDABalance    = "devIDABalance"
	QueryAffiliated    = "devAffiliated"
	QueryReservePool   = "devReservePool"
	QueryIDAStats      = "devIDAStats"
)
