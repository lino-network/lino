package register

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"
	"testing"
)

func TestRegisterBankDoesntExist(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	priv := crypto.GenPrivKeyEd25519()
	handler := NewHandler(lam)

	msg := NewRegisterMsg("register", priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, ErrAccRegisterFail("Get bank failed").Result())
}

func TestRegister(t *testing.T) {
	register := "register"
	lam := newLinoAccountManager()
	ctx := getContext()
	priv := crypto.GenPrivKeyEd25519()

	accBank := types.AccountBank{
		Address: priv.PubKey().Address(),
		PubKey:  priv.PubKey(),
		Coins:   sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	err := lam.SetBank(ctx, priv.PubKey().Address(), &accBank)
	assert.Nil(t, err)

	handler := NewHandler(lam)

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	accInfo := types.AccountInfo{
		Username: types.AccountKey(register),
		Created:  types.Height(0),
		PostKey:  priv.PubKey(),
		OwnerKey: priv.PubKey(),
		Address:  priv.PubKey().Address(),
	}
	infoPtr, err := lam.GetInfo(ctx, types.AccountKey(register))
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *infoPtr, "Account info should be equal")

	bankPtr, err := lam.GetBankFromAccountKey(ctx, types.AccountKey(register))
	assert.Nil(t, err)
	accBank.Username = types.AccountKey(register)
	assert.Equal(t, accBank, *bankPtr, "Account bank should be equal")

	accMeta := types.AccountMeta{
		LastActivity:   types.Height(ctx.BlockHeight()),
		ActivityBurden: types.DefaultActivityBurden,
		LastABBlock:    types.Height(ctx.BlockHeight()),
	}
	metaPtr, err := lam.GetMeta(ctx, types.AccountKey(register))
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *metaPtr, "Account meta should be equal")

	follower := types.Follower{Follower: []types.AccountKey{}}
	followerPtr, err := lam.GetFollower(ctx, types.AccountKey(register))
	assert.Nil(t, err)
	assert.Equal(t, follower, *followerPtr, "Account follower should be equal")

	following := types.Following{Following: []types.AccountKey{}}
	followingPtr, err := lam.GetFollowing(ctx, types.AccountKey(register))
	assert.Nil(t, err)
	assert.Equal(t, following, *followingPtr, "Account follower should be equal")
}

func TestRegisterFeeInsufficient(t *testing.T) {
	register := "register"
	lam := newLinoAccountManager()
	ctx := getContext()
	priv := crypto.GenPrivKeyEd25519()

	accBank := types.AccountBank{
		Address: priv.PubKey().Address(),
		PubKey:  priv.PubKey(),
		Coins:   RegisterFee.Minus(sdk.Coins{sdk.Coin{Denom: "Lino", Amount: 1}}),
	}
	err := lam.SetBank(ctx, priv.PubKey().Address(), &accBank)
	assert.Nil(t, err)

	handler := NewHandler(lam)

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, ErrAccRegisterFail("Register Fee Doesn't enough").Result())
}

func TestRegisterDuplicate(t *testing.T) {
	register := "register"
	lam := newLinoAccountManager()
	ctx := getContext()
	priv := crypto.GenPrivKeyEd25519()

	accBank := types.AccountBank{
		Address: priv.PubKey().Address(),
		PubKey:  priv.PubKey(),
		Coins:   sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	err := lam.SetBank(ctx, priv.PubKey().Address(), &accBank)
	assert.Nil(t, err)

	handler := NewHandler(lam)

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrAccRegisterFail("Username exist").Result())
}

func TestReRegister(t *testing.T) {
	register := "register"
	newRegister := "newRegister"
	lam := newLinoAccountManager()
	ctx := getContext()
	priv := crypto.GenPrivKeyEd25519()

	accBank := types.AccountBank{
		Address: priv.PubKey().Address(),
		PubKey:  priv.PubKey(),
		Coins:   sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	err := lam.SetBank(ctx, priv.PubKey().Address(), &accBank)
	assert.Nil(t, err)

	handler := NewHandler(lam)

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	msg = NewRegisterMsg(newRegister, priv.PubKey())
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrAccRegisterFail("Already registered").Result())
}
