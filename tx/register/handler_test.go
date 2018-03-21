package register

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"
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

	accBank := acc.AccountBank{
		Address: priv.PubKey().Address(),
		Balance: sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	err := lam.SetBankFromAddress(ctx, priv.PubKey().Address(), &accBank)
	assert.Nil(t, err)

	handler := NewHandler(lam)

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	accInfo := acc.AccountInfo{
		Username: acc.AccountKey(register),
		Created:  types.Height(0),
		PostKey:  priv.PubKey(),
		OwnerKey: priv.PubKey(),
		Address:  priv.PubKey().Address(),
	}
	infoPtr, err := lam.GetInfo(ctx, acc.AccountKey(register))
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *infoPtr, "Account info should be equal")

	bankPtr, err := lam.GetBankFromAccountKey(ctx, acc.AccountKey(register))
	assert.Nil(t, err)
	accBank.Username = acc.AccountKey(register)
	assert.Equal(t, accBank, *bankPtr, "Account bank should be equal")

	accMeta := acc.AccountMeta{
		LastActivity:   types.Height(ctx.BlockHeight()),
		ActivityBurden: types.DefaultActivityBurden,
	}
	metaPtr, err := lam.GetMeta(ctx, acc.AccountKey(register))
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *metaPtr, "Account meta should be equal")

	follower := acc.Follower{Follower: []acc.AccountKey{}}
	followerPtr, err := lam.GetFollower(ctx, acc.AccountKey(register))
	assert.Nil(t, err)
	assert.Equal(t, follower, *followerPtr, "Account follower should be equal")

	following := acc.Following{Following: []acc.AccountKey{}}
	followingPtr, err := lam.GetFollowing(ctx, acc.AccountKey(register))
	assert.Nil(t, err)
	assert.Equal(t, following, *followingPtr, "Account follower should be equal")
}

func TestRegisterFeeInsufficient(t *testing.T) {
	register := "register"
	lam := newLinoAccountManager()
	ctx := getContext()
	priv := crypto.GenPrivKeyEd25519()

	accBank := acc.AccountBank{
		Address: priv.PubKey().Address(),
		Balance: RegisterFee.Minus(sdk.Coins{sdk.Coin{Denom: "Lino", Amount: 1}}),
	}
	err := lam.SetBankFromAddress(ctx, priv.PubKey().Address(), &accBank)
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

	accBank := acc.AccountBank{
		Address: priv.PubKey().Address(),
		Balance: sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	err := lam.SetBankFromAddress(ctx, priv.PubKey().Address(), &accBank)
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

	accBank := acc.AccountBank{
		Address: priv.PubKey().Address(),
		Balance: sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	err := lam.SetBankFromAddress(ctx, priv.PubKey().Address(), &accBank)
	assert.Nil(t, err)

	handler := NewHandler(lam)

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	msg = NewRegisterMsg(newRegister, priv.PubKey())
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrAccRegisterFail("Already registered").Result())
}
