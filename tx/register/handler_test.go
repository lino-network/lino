package register

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestKVStoreKey = sdk.NewKVStoreKey("account")
)

func setupTest(t *testing.T) (*acc.AccountManager, sdk.Context, sdk.Handler) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	am := acc.NewAccountManager(capKey)
	handler := NewHandler(*am)
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)

	return am, ctx, handler
}

func createBank(t *testing.T, ctx sdk.Context, am *acc.AccountManager, coin types.Coin) crypto.PrivKey {
	priv := crypto.GenPrivKeyEd25519()
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), coin)
	assert.Nil(t, err)

	return priv.Wrap()
}

func TestRegisterBankDoesntExist(t *testing.T) {
	_, ctx, handler := setupTest(t)
	priv := crypto.GenPrivKeyEd25519()

	msg := NewRegisterMsg("register", priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, acc.ErrAccountCreateFailed(types.AccountKey("register")).Result().Code, result.Code)
}

func TestRegister(t *testing.T) {
	register := "register"
	am, ctx, handler := setupTest(t)
	priv := createBank(t, ctx, am, types.NewCoin(123*types.Decimals))

	assert.False(t, am.IsAccountExist(ctx, types.AccountKey(register)))

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	assert.True(t, am.IsAccountExist(ctx, types.AccountKey(register)))
}

func TestRegisterFeeInsufficient(t *testing.T) {
	register := "register"
	am, ctx, handler := setupTest(t)
	priv := createBank(t, ctx, am, types.NewCoin(23*types.Decimals))

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, acc.ErrRegisterFeeInsufficient().Result().Code, result.Code)
}

func TestRegisterDuplicate(t *testing.T) {
	am, ctx, handler := setupTest(t)
	priv := createBank(t, ctx, am, types.NewCoin(123*types.Decimals))
	register := "register"

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	result = handler(ctx, msg)
	assert.Equal(t, result.Code, acc.ErrAccountAlreadyExists(types.AccountKey(register)).Result().Code)
}

func TestReRegister(t *testing.T) {
	am, ctx, handler := setupTest(t)
	priv := createBank(t, ctx, am, types.NewCoin(123*types.Decimals))
	register := "register"
	newRegister := "newRegister"

	msg := NewRegisterMsg(register, priv.PubKey())
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	msg = NewRegisterMsg(newRegister, priv.PubKey())
	result = handler(ctx, msg)
	assert.Equal(t, result.Code, acc.ErrBankAlreadyRegistered().Result().Code)
}
