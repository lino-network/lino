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

// GetPostKey try to generate PostKey from types.AccountKey and PostID
func GetPermLink(author AccountKey, postID string) PermLink {
	return PermLink(string(author) + "#" + postID)
}

// Donation struct, only used in Donation
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}

// TODO need to store into KVStore
var VoterMinDeposit = NewCoin(1000 * Decimals)
var VoterMinWithdraw = NewCoin(1 * Decimals)
var DelegatorMinWithdraw = NewCoin(1 * Decimals)

var NextProposalID = int64(0)
var ProposalDecideHr = int64(7 * 24)
var ProposalRegisterFee = NewCoin(2000 * Decimals)

var ValidatorMinWithdraw = NewCoin(1 * Decimals)
var ValidatorMinVotingDeposit = NewCoin(3000 * Decimals)
var ValidatorMinCommitingDeposit = NewCoin(1000 * Decimals)

var CoinReturnIntervalHr = int64(7 * 24)
var CoinReturnTimes = int64(7)

var PenaltyMissVote = NewCoin(200 * Decimals)
var PenaltyMissCommit = NewCoin(200 * Decimals)
var PenaltyByzantine = NewCoin(1000 * Decimals)

var DeveloperMinDeposit = NewCoin(100000 * Decimals)
