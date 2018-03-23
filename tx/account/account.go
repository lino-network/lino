package account

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

type Memo uint64

// AccountKey key format in KVStore
type AccountKey string

// AccountInfo stores general Lino Account information
type AccountInfo struct {
	Username AccountKey    `json:"key"`
	Created  types.Height  `json:"created"`
	PostKey  crypto.PubKey `json:"post_key"`
	OwnerKey crypto.PubKey `json:"owner_key"`
	Address  sdk.Address   `json:"address"`
}

// AccountBank uses Address as the key instead of Username
type AccountBank struct {
	Address  sdk.Address `json:"address"`
	Balance  sdk.Coins   `json:"coins"`
	Username AccountKey  `json:"Username"`
}

// AccountMeta stores tiny and frequently updated fields.
type AccountMeta struct {
	Sequence       int64        `json:"sequence"`
	LastActivity   types.Height `json:"last_activity"`
	ActivityBurden int64        `json:"activity_burden"`
}

// Follower records all follower belong to one user
type Follower struct {
	Follower []AccountKey `json:"follower"`
}

// Following records all follower belong to one user
type Following struct {
	Following []AccountKey `json:"following"`
}

// linoaccount encapsulates all basic struct
type Account struct {
	username           AccountKey      `json:"username"`
	writeInfoFlag      bool            `json:"write_info_flag"`
	writeBankFlag      bool            `json:"write_bank_flag"`
	writeMetaFlag      bool            `json:"write_meta_flag"`
	writeFollowerFlag  bool            `json:"write_follower_flag"`
	writeFollowingFlag bool            `json:"write_following_flag"`
	accountManager     *AccountManager `json:"account_manager"`
	accountInfo        *AccountInfo    `json:"account_info"`
	accountBank        *AccountBank    `json:"account_bank"`
	accountMeta        *AccountMeta    `json:"account_meta"`
	follower           *Follower       `json:"follower"`
	following          *Following      `json:"following"`
}

func RegisterWireLinoAccount(cdc *wire.Codec) {
	// Register crypto.[PubKey] types.
	wire.RegisterCrypto(cdc)
}

// NewLinoAccount return the account pointer
func NewProxyAccount(username AccountKey, accManager *AccountManager) *Account {
	return &Account{
		username:       username,
		accountManager: accManager,
	}
}

// check if post exist
func (acc *Account) IsAccountExist(ctx sdk.Context) bool {
	if err := acc.checkAccountInfo(ctx); err != nil {
		return false
	}
	return true
}

// Implements types.AccountManager.
func (acc *Account) CreateAccount(ctx sdk.Context, accKey AccountKey, pubkey crypto.PubKey, accBank *AccountBank) sdk.Error {
	if acc.IsAccountExist(ctx) {
		return ErrAccountCreateFail(fmt.Sprintf("account exist: %v", accKey))
	}
	acc.writeInfoFlag = true
	acc.accountInfo = &AccountInfo{
		Username: accKey,
		Created:  types.Height(ctx.BlockHeight()),
		PostKey:  pubkey,
		OwnerKey: pubkey,
		Address:  pubkey.Address(),
	}

	acc.writeBankFlag = true
	accBank.Username = accKey
	acc.accountBank = accBank

	acc.writeMetaFlag = true
	acc.accountMeta = &AccountMeta{
		LastActivity:   types.Height(ctx.BlockHeight()),
		ActivityBurden: types.DefaultActivityBurden,
	}

	acc.writeFollowerFlag = true
	acc.follower = &Follower{Follower: []AccountKey{}}

	acc.writeFollowingFlag = true
	acc.following = &Following{Following: []AccountKey{}}

	return nil
}

func (acc *Account) AddCoins(ctx sdk.Context, coins sdk.Coins) (err sdk.Error) {
	if err := acc.checkAccountBank(ctx); err != nil {
		return err
	}
	acc.accountBank.Balance = acc.accountBank.Balance.Plus(coins)
	acc.writeBankFlag = true
	return nil
}

func (acc *Account) MinusCoins(ctx sdk.Context, coins sdk.Coins) (err sdk.Error) {
	if err := acc.checkAccountBank(ctx); err != nil {
		return err
	}

	if !acc.accountBank.Balance.IsGTE(coins) {
		return ErrAccountManagerFail("Account bank's coins are not enough")
	}

	c0 := sdk.Coins{sdk.Coin{Denom: "lino", Amount: int64(0)}}
	acc.accountBank.Balance = acc.accountBank.Balance.Minus(coins)

	// API return empty when the result is 0 coin
	if len(acc.accountBank.Balance) == 0 {
		acc.accountBank.Balance = c0
	}

	acc.writeBankFlag = true
	return nil
}

func (acc *Account) GetUsername(ctx sdk.Context) AccountKey {
	return acc.username
}

func (acc *Account) GetBankAddress(ctx sdk.Context) (sdk.Address, sdk.Error) {
	if err := acc.checkAccountInfo(ctx); err != nil {
		return nil, err
	}
	return acc.accountInfo.Address, nil
}

func (acc *Account) GetOwnerKey(ctx sdk.Context) (*crypto.PubKey, sdk.Error) {
	if err := acc.checkAccountInfo(ctx); err != nil {
		return nil, err
	}
	return &acc.accountInfo.OwnerKey, nil
}

func (acc *Account) GetPostKey(ctx sdk.Context) (*crypto.PubKey, sdk.Error) {
	if err := acc.checkAccountInfo(ctx); err != nil {
		return nil, err
	}
	return &acc.accountInfo.PostKey, nil
}

func (acc *Account) GetBankBalance(ctx sdk.Context) (sdk.Coins, sdk.Error) {
	if err := acc.checkAccountBank(ctx); err != nil {
		return nil, err
	}
	return acc.accountBank.Balance, nil
}

func (acc *Account) GetSequence(ctx sdk.Context) (int64, sdk.Error) {
	if err := acc.checkAccountMeta(ctx); err != nil {
		return 0, err
	}
	return acc.accountMeta.Sequence, nil
}

func (acc *Account) GetCreated(ctx sdk.Context) (types.Height, sdk.Error) {
	if err := acc.checkAccountInfo(ctx); err != nil {
		return types.Height(0), err
	}
	return acc.accountInfo.Created, nil
}

func (acc *Account) GetLastActivity(ctx sdk.Context) (types.Height, sdk.Error) {
	if err := acc.checkAccountMeta(ctx); err != nil {
		return types.Height(0), err
	}
	return acc.accountMeta.LastActivity, nil
}

func (acc *Account) IncreaseSequenceByOne(ctx sdk.Context) sdk.Error {
	if err := acc.checkAccountMeta(ctx); err != nil {
		return err
	}
	acc.accountMeta.Sequence += 1
	acc.writeMetaFlag = true
	return nil
}

func (acc *Account) GetActivityBurden(ctx sdk.Context) (int64, sdk.Error) {
	if err := acc.checkAccountMeta(ctx); err != nil {
		return 0, err
	}
	return acc.accountMeta.ActivityBurden, nil
}

func (acc *Account) GetFollower(ctx sdk.Context) (*Follower, sdk.Error) {
	if err := acc.checkAccountFollower(ctx); err != nil {
		return nil, err
	}
	return acc.follower, nil
}

func (acc *Account) GetFollowing(ctx sdk.Context) (*Following, sdk.Error) {
	if err := acc.checkAccountFollowing(ctx); err != nil {
		return nil, err
	}
	return acc.following, nil
}

func (acc *Account) UpdateLastActivity(ctx sdk.Context) sdk.Error {
	if err := acc.checkAccountMeta(ctx); err != nil {
		return err
	}
	acc.writeMetaFlag = true
	acc.accountMeta.LastActivity = types.Height(ctx.BlockHeight())
	return nil
}

func (acc *Account) Apply(ctx sdk.Context) sdk.Error {
	if acc.writeInfoFlag {
		if err := acc.accountManager.SetInfo(ctx, acc.username, acc.accountInfo); err != nil {
			return err
		}
	}
	if acc.writeBankFlag {
		if err := acc.checkAccountInfo(ctx); err != nil {
			return err
		}
		if err := acc.accountManager.SetBankFromAddress(ctx, acc.accountInfo.Address, acc.accountBank); err != nil {
			return err
		}
	}
	if acc.writeMetaFlag {
		if err := acc.accountManager.SetMeta(ctx, acc.username, acc.accountMeta); err != nil {
			return err
		}
	}
	if acc.writeFollowerFlag {
		if err := acc.accountManager.SetFollower(ctx, acc.username, acc.follower); err != nil {
			return err
		}
	}
	if acc.writeFollowingFlag {
		if err := acc.accountManager.SetFollowing(ctx, acc.username, acc.following); err != nil {
			return err
		}
	}

	acc.clear()

	return nil
}

func (acc *Account) clear() {
	acc.writeInfoFlag = false
	acc.writeBankFlag = false
	acc.writeMetaFlag = false
	acc.writeFollowerFlag = false
	acc.writeFollowingFlag = false
	acc.accountInfo = nil
	acc.accountBank = nil
	acc.accountMeta = nil
	acc.follower = nil
	acc.following = nil
}

func (acc *Account) checkAccountInfo(ctx sdk.Context) (err sdk.Error) {
	if acc.accountInfo == nil {
		acc.accountInfo, err = acc.accountManager.GetInfo(ctx, acc.username)
	}
	return err
}

func (acc *Account) checkAccountBank(ctx sdk.Context) (err sdk.Error) {
	if err = acc.checkAccountInfo(ctx); err != nil {
		return err
	}
	if acc.accountBank == nil {
		acc.accountBank, err = acc.accountManager.GetBankFromAddress(ctx, acc.accountInfo.Address)
	}
	return err
}

func (acc *Account) checkAccountMeta(ctx sdk.Context) (err sdk.Error) {
	if acc.accountMeta == nil {
		acc.accountMeta, err = acc.accountManager.GetMeta(ctx, acc.username)
	}
	return err
}

func (acc *Account) checkAccountFollower(ctx sdk.Context) (err sdk.Error) {
	if acc.follower == nil {
		acc.follower, err = acc.accountManager.GetFollower(ctx, acc.username)
	}
	return err
}

func (acc *Account) checkAccountFollowing(ctx sdk.Context) (err sdk.Error) {
	if acc.following == nil {
		acc.following, err = acc.accountManager.GetFollowing(ctx, acc.username)
	}
	return err
}
