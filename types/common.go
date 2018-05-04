package types

// Validator List Size
var ValidatorListSize = 21

var AbsentCommitLimitation = 100

var AbsentVoteLimitation = 100

// AccountKey key format in KVStore
type AccountKey string

// PostKey key format in KVStore
type PermLink string

// ProposalKey key format in KVStore
type ProposalKey string

// user permission type to present different permission for different msg
type Permission int

// coin return type to present different coin return parameters for different user
type CoinReturn int

// GetPostKey try to generate PostKey from types.AccountKey and PostID
func GetPermLink(author AccountKey, postID string) PermLink {
	return PermLink(string(author) + "#" + postID)
}

// Donation struct, only used in Donation
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}
