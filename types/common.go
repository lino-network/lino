package types

import (
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

// AccOrAddr is either an address or an account key.
type AccOrAddr struct {
	AccountKey AccountKey     `json:"account_key,omitempty"`
	Addr       sdk.AccAddress `json:"addr,omitempty"`
	IsAddr     bool           `json:"is_addr,omitempty"`
}

func NewAccOrAddrFromAcc(acc AccountKey) AccOrAddr {
	return AccOrAddr{
		AccountKey: acc,
	}
}

func NewAccOrAddrFromAddr(addr sdk.AccAddress) AccOrAddr {
	return AccOrAddr{
		Addr:   addr,
		IsAddr: true,
	}
}

func (a AccOrAddr) IsValid() bool {
	if a.IsAddr {
		return len(a.Addr) > 0
	}
	return a.AccountKey.IsValid()
}

func (a AccOrAddr) String() string {
	if a.IsAddr {
		return a.Addr.String()
	}
	return string(a.AccountKey)
}

// GetPermlink try to generate Permlink from AccountKey and PostID
func GetPermlink(author AccountKey, postID string) Permlink {
	return Permlink(string(author) + PermlinkSeparator + postID)
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

func (ak AccountKey) IsValid() bool {
	if len(ak) < MinimumUsernameLength ||
		len(ak) > MaximumUsernameLength {
		return false
	}

	match, err := regexp.MatchString(UsernameReCheck, string(ak))
	if !match || err != nil {
		return false
	}

	match, err = regexp.MatchString(IllegalUsernameReCheck, string(ak))
	if match || err != nil {
		return false
	}
	return true
}
