package account

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"
)

func TestAccountInfo(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	priv := crypto.GenPrivKeyEd25519()
	accInfo := AccountInfo{
		Username: AccountKey("test"),
		Created:  0,
		PostKey:  priv.PubKey(),
		OwnerKey: priv.PubKey(),
		Address:  priv.PubKey().Address(),
	}
	err := lam.SetInfo(ctx, AccountKey("test"), &accInfo)
	assert.Nil(t, err)

	resultPtr, err := lam.GetInfo(ctx, AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *resultPtr, "Account info should be equal")
}

func TestInvalidAccountInfo(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	resultPtr, err := lam.GetInfo(ctx, AccountKey("test"))
	assert.Nil(t, resultPtr)
	assert.Equal(t, err, ErrAccountManagerFail("AccountManager get info failed: info doesn't exist"))
}

func TestAccountBank(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	priv := crypto.GenPrivKeyEd25519()
	accInfo := AccountInfo{
		Username: AccountKey("test"),
		Created:  0,
		PostKey:  priv.PubKey(),
		OwnerKey: priv.PubKey(),
		Address:  priv.PubKey().Address(),
	}
	err := lam.SetInfo(ctx, AccountKey("test"), &accInfo)
	assert.Nil(t, err)

	accBank := AccountBank{
		Address: priv.PubKey().Address(),
		Balance: sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	err = lam.SetBankFromAddress(ctx, priv.PubKey().Address(), &accBank)
	assert.Nil(t, err)

	resultPtr, err := lam.GetBankFromAccountKey(ctx, AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accBank, *resultPtr, "Account bank should be equal")

	resultPtr, err = lam.GetBankFromAddress(ctx, priv.PubKey().Address())
	assert.Nil(t, err)
	assert.Equal(t, accBank, *resultPtr, "Account bank should be equal")
}

func TestAccountMeta(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	accMeta := AccountMeta{}
	err := lam.SetMeta(ctx, AccountKey("test"), &accMeta)
	assert.Nil(t, err)

	resultPtr, err := lam.GetMeta(ctx, AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *resultPtr, "Account meta should be equal")
}

func TestAccountFollower(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	follower := Follower{Follower: []AccountKey{}}
	err := lam.SetFollower(ctx, AccountKey("test"), &follower)
	assert.Nil(t, err)

	resultPtr, err := lam.GetFollower(ctx, AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, follower, *resultPtr, "Account follower should be equal")
}

func TestAccountFollowing(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	following := Following{Following: []AccountKey{}}
	err := lam.SetFollowing(ctx, AccountKey("test"), &following)
	assert.Nil(t, err)

	resultPtr, err := lam.GetFollowing(ctx, AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, following, *resultPtr, "Account follower should be equal")
}
