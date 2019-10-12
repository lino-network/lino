package types

const (
	// ModuleName is module name
	ModuleName = "proposal"

	// RouterKey is the message route for post
	RouterKey = ModuleName

	// QuerierRoute is the querier route for post
	QuerierRoute = ModuleName

	QueryNextProposal    = "next"
	QueryOngoingProposal = "ongoing"
	QueryExpiredProposal = "expired"
)
