package register

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
	TestParamKVStoreKey   = sdk.NewKVStoreKey("account")
)

func setupTest(t *testing.T) (acc.AccountManager, sdk.Context, sdk.Handler) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)
	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)
	am := acc.NewAccountManager(TestAccountKVStoreKey, ph)
	handler := NewHandler(am)

	return am, ctx, handler
}

func createBank(t *testing.T, ctx sdk.Context, am acc.AccountManager, coin types.Coin) crypto.PrivKeyEd25519 {
	priv := crypto.GenPrivKeyEd25519()
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), coin)
	assert.Nil(t, err)

	return priv
}

func TestRegisterBankDoesntExist(t *testing.T) {
	_, ctx, handler := setupTest(t)
	priv := crypto.GenPrivKeyEd25519()

	msg := NewRegisterMsg("register", priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, acc.ErrAccountCreateFailed(types.AccountKey("register")).Result().Code, result.Code)
}

func TestRegister(t *testing.T) {
	register := "register"
	am, ctx, handler := setupTest(t)
	priv := createBank(t, ctx, am, types.NewCoin(123*types.Decimals))

	assert.False(t, am.IsAccountExist(ctx, types.AccountKey(register)))

	msg := NewRegisterMsg(register, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	assert.True(t, am.IsAccountExist(ctx, types.AccountKey(register)))
}

func TestDuplicateRegister(t *testing.T) {
	am, ctx, handler := setupTest(t)
	priv := createBank(t, ctx, am, types.NewCoin(123*types.Decimals))
	register := "register"

	msg := NewRegisterMsg(register, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	result = handler(ctx, msg)
	assert.Equal(t, result.Code, acc.ErrAccountAlreadyExists(types.AccountKey(register)).Result().Code)
}

func TestBankReRegister(t *testing.T) {
	am, ctx, handler := setupTest(t)
	priv := createBank(t, ctx, am, types.NewCoin(123*types.Decimals))
	register := "register"
	newRegister := "newRegister"

	msg := NewRegisterMsg(register, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	msg = NewRegisterMsg(newRegister, priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	result = handler(ctx, msg)
	assert.Equal(t, result.Code, acc.ErrBankAlreadyRegistered().Result().Code)
}
