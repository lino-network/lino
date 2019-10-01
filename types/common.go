package types

// AccountKey key format in KVStore
type AccountKey string

// Permlink key format in KVStore
type Permlink string

// ProposalKey key format in KVStore
type ProposalKey string

// user permission type to present different permission for different msg
type Permission int

// msg CapacityLevel, different level cost different user capacity
type CapacityLevel int

// indicates the current proposal status
type ProposalResult int

// indicates proposal type
type ProposalType int

// indicates donation type
type DonationType int

// indicates all possible balance behavior types
type TransferDetailType int

// indicates the type of punishment for oncall validators
type PunishType int

// GetPostKey try to generate PostKey from types.AccountKey and PostID
func GetPermlink(author AccountKey, postID string) Permlink {
	return Permlink(string(author) + PermlinkSeparator + postID)
}

// Donation struct, only used in Donation
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}

// PenaltyList - get validator who doesn't vote for proposal
type PenaltyList struct {
	PenaltyList []AccountKey `json:"penalty_list"`
}

// FindAccountInList - find AccountKey in given AccountKey list
func FindAccountInList(me AccountKey, lst []AccountKey) int {
	for index, user := range lst {
		if user == me {
			return index
		}
	}
	return -1
}

func AccountListToSet(lst []AccountKey) map[AccountKey]bool {
	rst := make(map[AccountKey]bool)
	for _, acc := range lst {
		rst[acc] = true
	}
	return rst
}
