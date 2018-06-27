package types

// AccountKey key format in KVStore
type AccountKey string

// Permlink key format in KVStore
type Permlink string

// ProposalKey key format in KVStore
type ProposalKey string

// user permission type to present different permission for different msg
type Permission int

// indicates the current proposal status
type ProposalResult int

// indicates proposal type
type ProposalType int

// indicates donation type
type DonationType int

// indicates all possible balance behavior types
type TransferDetailType int

// GetPostKey try to generate PostKey from types.AccountKey and PostID
func GetPermlink(author AccountKey, postID string) Permlink {
	return Permlink(string(author) + "#" + postID)
}

// Donation struct, only used in Donation
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}

type PenaltyList struct {
	PenaltyList []AccountKey `json:"penalty_list"`
}
