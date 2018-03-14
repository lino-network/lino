package account

import (
	"errors"

	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

//-----------------------------------------------------------
// BaseAccount

// BaseAccount - base account structure.
// Extend this by embedding this in your AppAccount.
// See the examples/basecoin/types/account.go for an example.
type LinoAccount struct {
	Info       types.AccountInfo `json:"info"`
	Bank       types.AccountBank `json:"bank"`
	Meta       types.AccountMeta `json:"meta"`
	Followers  types.Followers   `json:"followers"`
	Followings types.Followings  `json:"followings"`
}

// Implements types.Account.
func (acc LinoAccount) GetUsername() types.AccountKey {
	return acc.Info.Username
}

// Implements types.Account.
func (acc LinoAccount) SetUsername(username types.AccountKey) error {
	if len(acc.Info.Username) != 0 {
		return errors.New("cannot override username")
	}
	acc.Info.Username = username
	return nil
}

// Implements types.Account.
func (acc LinoAccount) GetBankAddress() sdk.Address {
	return acc.Bank.Address
}

// Implements types.Account.
func (acc *LinoAccount) SetBankAddress(addr sdk.Address) error {
	acc.Bank.Address = addr
	return nil
}

// Implements types.Account.
func (acc LinoAccount) GetOwnerKey() crypto.PubKey {
	return acc.Info.OwnerKey
}

// Implements types.Account.
func (acc *LinoAccount) SetOwnerKey(pubKey crypto.PubKey) error {
	acc.Info.OwnerKey = pubKey
	return nil
}

// Implements types.Account.
func (acc LinoAccount) GetPostKey() crypto.PubKey {
	return acc.Info.PostKey
}

// Implements types.Account.
func (acc LinoAccount) SetPostKey(pubKey crypto.PubKey) error {
	acc.Info.PostKey = pubKey
	return nil
}

// Implements types.Account.
func (acc LinoAccount) GetBankBalance() sdk.Coins {
	return acc.Bank.Coins
}

// Implements types.Account.
func (acc LinoAccount) SetBankBalance(coins sdk.Coins) error {
	acc.Bank.Coins = coins
	return nil
}

// Implements types.Account.
func (acc LinoAccount) GetBankSequence() int64 {
	return acc.Bank.Sequence
}

// Implements types.Account.
func (acc LinoAccount) SetBankSequence(seq int64) error {
	acc.Bank.Sequence = seq
	return nil
}

func (acc LinoAccount) GetCreated() types.Height {
	return acc.Info.Created
}

func (acc LinoAccount) SetCreated(created types.Height) error { // errors if already set.
	if acc.Info.Created > 0 {
		return errors.New("cannot override created block height")
	}
	acc.Info.Created = created
	return nil
}

func (acc LinoAccount) GetLastActivity() types.Height {
	return acc.Meta.LastActivity
}

func (acc LinoAccount) SetLastActivity(lastActivity types.Height) error {
	acc.Meta.LastActivity = lastActivity
	return nil
}

func (acc LinoAccount) GetActivityBurden() uint64 {
	return acc.Meta.ActivityBurden
}
func (acc LinoAccount) SetActivityBurden(burden uint64, height types.Height) error { // set AB Block too.
	acc.Meta.ActivityBurden = burden
	acc.Meta.LastABBlock = height
	return nil
}

func (acc LinoAccount) GetLastABBlock() types.Height {
	return acc.Meta.LastABBlock
}

func (acc LinoAccount) GetFollowers() types.Followers {
	return acc.Followers
}
func (acc LinoAccount) SetFollowers(followers types.Followers) error {
	acc.Followers = followers
	return nil
}

func (acc *LinoAccount) GetFollowings() types.Followings {
	return acc.Followings
}
func (acc *LinoAccount) SetFollowings(followings types.Followings) error {
	acc.Followings = followings
	return nil
}

//----------------------------------------
// Wire
