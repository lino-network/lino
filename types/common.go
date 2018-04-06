package types

// Validator List Size
var ValidatorListSize = 21

var AbsentLimitation = 100

// AccountKey key format in KVStore
type AccountKey string

// PostKey key format in KVStore
type PostKey string

// ProposalKey key format in KVStore
type ProposalKey string

// GetPostKey try to generate PostKey from types.AccountKey and PostID
func GetPostKey(author AccountKey, postID string) PostKey {
	return PostKey(string(author) + "#" + postID)
}

// Donation struct, only used in Donation
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}
