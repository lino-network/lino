package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

// linoaccount encapsulates all basic struct
type linoaccount struct {
	username           types.AccountKey     `json:"username"`
	writeInfoFlag      bool                 `json:"write_info_flag"`
	writeBankFlag      bool                 `json:"write_bank_flag"`
	writeMetaFlag      bool                 `json:"write_meta_flag"`
	writeFollowerFlag  bool                 `json:"write_follower_flag"`
	writeFollowingFlag bool                 `json:"write_following_flag"`
	accountManager     types.AccountManager `json:"account_manager"`
	accountInfo        *types.AccountInfo   `json:"account_info"`
	accountBank        *types.AccountBank   `json:"account_bank"`
	accountMeta        *types.AccountMeta   `json:"account_meta"`
	follower           *types.Follower      `json:"follower"`
	following          *types.Following     `json:"following"`
}

// NewLinoAccount return the account pointer
func NewLinoAccount(username types.AccountKey, accManager types.AccountManager) *linoaccount {
	return &linoaccount{
		username:       username,
		accountManager: accManager,
	}
}

func (acc *linoaccount) GetUsername(ctx sdk.Context) types.AccountKey {
	return acc.username
}

func (acc *linoaccount) GetBankAddress(ctx sdk.Context) (sdk.Address, sdk.Error) {
	if err := acc.checkAccountInfo(ctx); err != nil {
		return nil, err
	}
	return acc.accountInfo.Address, nil
}

func (acc *linoaccount) GetOwnerKey(ctx sdk.Context) (*crypto.PubKey, sdk.Error) {
	if err := acc.checkAccountInfo(ctx); err != nil {
		return nil, err
	}
	return &acc.accountInfo.OwnerKey, nil
}

func (acc *linoaccount) GetPostKey(ctx sdk.Context) (*crypto.PubKey, sdk.Error) {
	if err := acc.checkAccountInfo(ctx); err != nil {
		return nil, err
	}
	return &acc.accountInfo.PostKey, nil
}

func (acc *linoaccount) GetBankBalance(ctx sdk.Context) (sdk.Coins, sdk.Error) {
	if err := acc.checkAccountBank(ctx); err != nil {
		return nil, err
	}
	return acc.accountBank.Balance, nil
}

func (acc *linoaccount) GetSequence(ctx sdk.Context) (int64, sdk.Error) {
	if err := acc.checkAccountMeta(ctx); err != nil {
		return 0, err
	}
	return acc.accountMeta.Sequence, nil
}

func (acc *linoaccount) GetCreated(ctx sdk.Context) (types.Height, sdk.Error) {
	if err := acc.checkAccountInfo(ctx); err != nil {
		return types.Height(0), err
	}
	return acc.accountInfo.Created, nil
}

func (acc *linoaccount) GetLastActivity(ctx sdk.Context) (types.Height, sdk.Error) {
	if err := acc.checkAccountMeta(ctx); err != nil {
		return types.Height(0), err
	}
	return acc.accountMeta.LastActivity, nil
}

func (acc *linoaccount) GetActivityBurden(ctx sdk.Context) (int64, sdk.Error) {
	if err := acc.checkAccountMeta(ctx); err != nil {
		return 0, err
	}
	return acc.accountMeta.ActivityBurden, nil
}

func (acc *linoaccount) GetFollower(ctx sdk.Context) (*types.Follower, sdk.Error) {
	if err := acc.checkAccountFollower(ctx); err != nil {
		return nil, err
	}
	return acc.follower, nil
}

func (acc *linoaccount) GetFollowing(ctx sdk.Context) (*types.Following, sdk.Error) {
	if err := acc.checkAccountFollowing(ctx); err != nil {
		return nil, err
	}
	return acc.following, nil
}

func (acc *linoaccount) UpdateLastActivity(ctx sdk.Context) sdk.Error {
	if err := acc.checkAccountMeta(ctx); err != nil {
		return err
	}
	acc.writeMetaFlag = true
	acc.accountMeta.LastActivity = types.Height(ctx.BlockHeight())
	return nil
}

func (acc *linoaccount) Apply(ctx sdk.Context) sdk.Error {
	if acc.writeInfoFlag {
		if err := acc.accountManager.SetInfo(ctx, acc.username, acc.accountInfo); err != nil {
			return err
		}
	}
	if acc.writeBankFlag {
		if err := acc.checkAccountInfo(ctx); err != nil {
			return err
		}
		if err := acc.accountManager.SetBank(ctx, acc.accountInfo.Address, acc.accountBank); err != nil {
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
	return nil
}

func (acc *linoaccount) checkAccountInfo(ctx sdk.Context) (err sdk.Error) {
	if acc.accountInfo == nil {
		acc.accountInfo, err = acc.accountManager.GetInfo(ctx, acc.username)
	}
	return err
}

func (acc *linoaccount) checkAccountBank(ctx sdk.Context) (err sdk.Error) {
	if err = acc.checkAccountInfo(ctx); err != nil {
		return err
	}
	if acc.accountBank == nil {
		acc.accountBank, err = acc.accountManager.GetBankFromAddress(ctx, acc.accountInfo.Address)
	}
	return err
}

func (acc *linoaccount) checkAccountMeta(ctx sdk.Context) (err sdk.Error) {
	if acc.accountMeta == nil {
		acc.accountMeta, err = acc.accountManager.GetMeta(ctx, acc.username)
	}
	return err
}

func (acc *linoaccount) checkAccountFollower(ctx sdk.Context) (err sdk.Error) {
	if acc.follower == nil {
		acc.follower, err = acc.accountManager.GetFollower(ctx, acc.username)
	}
	return err
}

func (acc *linoaccount) checkAccountFollowing(ctx sdk.Context) (err sdk.Error) {
	if acc.following == nil {
		acc.following, err = acc.accountManager.GetFollowing(ctx, acc.username)
	}
	return err
}
