package types

// Validator List Size
var ValidatorListSize = 21

var AbsentCommitLimitation = 100

var AbsentVoteLimitation = 100

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

// TODO need to store into KV
var VoterRegisterFee = Coin{Amount: 1000 * Decimals}
var VoterMinimumWithdraw = Coin{Amount: 50 * Decimals}
var ProposalRegisterFee = Coin{Amount: 2000 * Decimals}
var ValidatorRegisterFee = Coin{Amount: 1000 * Decimals}
var ValidatorMinimumWithdraw = Coin{Amount: 50 * Decimals}
var NextProposalID = int64(0)
var ProposalDecideHr = int64(7 * 24)
var CoinReturnIntervalHr = int64(7 * 24)
var CoinReturnTimes = int64(7)
