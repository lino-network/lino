package types

const (
	// ModuleKey is the name of the module
	ModuleName = "vote"

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	QueryVoter      = "voter"
	QueryStakeStats = "stake-stats"
)
