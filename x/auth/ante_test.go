package auth

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	"github.com/lino-network/lino/x/global"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var (
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
	TestGlobalKVStoreKey  = sdk.NewKVStoreKey("global")
	TestParamKVStoreKey   = sdk.NewKVStoreKey("param")
)

func createTestAccount(
	ctx sdk.Context, am acc.AccountManager, ph param.ParamHolder, username string) (secp256k1.PrivKeySecp256k1,
	secp256k1.PrivKeySecp256k1, secp256k1.PrivKeySecp256k1, types.AccountKey) {
	resetKey := secp256k1.GenPrivKey()
	transactionKey := secp256k1.GenPrivKey()
	appKey := secp256k1.GenPrivKey()
	accParams, _ := ph.GetAccountParam(ctx)
	am.CreateAccount(ctx, "referrer", types.AccountKey(username),
		resetKey.PubKey(), transactionKey.PubKey(), appKey.PubKey(), accParams.RegisterFee)
	return resetKey, transactionKey, appKey, types.AccountKey(username)
}

func InitGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest() (
	acc.AccountManager, global.GlobalManager, param.ParamHolder, sdk.Context, sdk.AnteHandler) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: 1, Time: time.Now()}, false, log.NewNopLogger())

	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)
	am := acc.NewAccountManager(TestAccountKVStoreKey, ph)
	gm := global.NewGlobalManager(TestGlobalKVStoreKey, ph)
	InitGlobalManager(ctx, gm)
	anteHandler := NewAnteHandler(am, gm)

	return am, gm, ph, ctx, anteHandler
}

type TestMsg struct {
	Signers    []types.AccountKey
	Permission types.Permission
	Amount     types.Coin
}

var _ types.Msg = TestMsg{}

func (msg TestMsg) Route() string                   { return "normal msg" }
func (msg TestMsg) Type() string                    { return "normal msg" }
func (msg TestMsg) GetPermission() types.Permission { return msg.Permission }
func (msg TestMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.Signers)
	if err != nil {
		panic(err)
	}
	return bz
}
func (msg TestMsg) ValidateBasic() sdk.Error { return nil }
func (msg TestMsg) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Signers))
	for i, signer := range msg.Signers {
		addrs[i] = sdk.AccAddress(signer)
	}
	return addrs
}
func (msg TestMsg) GetConsumeAmount() types.Coin {
	return msg.Amount
}

func newTestMsg(accKeys ...types.AccountKey) TestMsg {
	return TestMsg{
		Signers:    accKeys,
		Permission: types.AppPermission,
		Amount:     types.NewCoinFromInt64(10),
	}
}

// run the tx through the anteHandler and ensure its valid
func checkValidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx) {
	_, result, abort := anteHandler(ctx, tx, false)
	assert.False(t, abort)
	assert.True(t, result.Code.IsOK()) // redundent
	assert.True(t, result.IsOK())
}

// run the tx through the anteHandler and ensure it fails with the given code
func checkInvalidTx(
	t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, result sdk.Result) {
	_, r, abort := anteHandler(ctx, tx, false)
	assert.True(t, abort)
	assert.Equal(t, result, r)
}

func newTestTx(
	ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, seqs []uint64) sdk.Tx {
	sigs := make([]auth.StdSignature, len(privs))

	for i, priv := range privs {
		signBytes := auth.StdSignBytes(ctx.ChainID(), 0, seqs[i], auth.StdFee{}, msgs, "")
		bz, _ := priv.Sign(signBytes)
		sigs[i] = auth.StdSignature{
			PubKey: priv.PubKey(), Signature: bz}
	}
	tx := auth.NewStdTx(msgs, auth.StdFee{}, sigs, "")
	return tx
}

// Test various error cases in the AnteHandler control flow.
func TestAnteHandlerSigErrors(t *testing.T) {
	// setup
	am, _, ph, ctx, anteHandler := setupTest()
	// get private key and username
	_, transaction1, _, user1 := createTestAccount(ctx, am, ph, "user1")
	_, transaction2, _, user2 := createTestAccount(ctx, am, ph, "user2")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1, user2)

	// test no signatures
	privs, seqs := []crypto.PrivKey{}, []uint64{}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, ErrNoSignatures().Result())

	// test num sigs less than GetSigners
	privs, seqs = []crypto.PrivKey{transaction1}, []uint64{0}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(
		t, anteHandler, ctx, tx, ErrWrongNumberOfSigners().Result())

	// test sig user mismatch
	privs, seqs = []crypto.PrivKey{transaction2, transaction1}, []uint64{0, 0}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())
}

// Test various error cases in the AnteHandler control flow.
func TestAnteHandlerNormalTx(t *testing.T) {
	am, _, ph, ctx, anteHandler := setupTest()
	// keys and username
	_, transaction1, _, user1 := createTestAccount(ctx, am, ph, "user1")
	_, transaction2, _, _ := createTestAccount(ctx, am, ph, "user2")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1)

	// test valid transaction
	privs, seqs := []crypto.PrivKey{transaction1}, []uint64{0}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)
	seq, err := am.GetSequence(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, seq, uint64(1))

	// test no signatures
	privs, seqs = []crypto.PrivKey{}, []uint64{}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, ErrNoSignatures().Result())

	// test wrong sequence number, now we return signature failed even it's seq number error.
	privs, seqs = []crypto.PrivKey{transaction1}, []uint64{0}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, ErrUnverifiedBytes(
		"signature verification failed, chain-id:Lino, seq:1").Result())

	// test wrong priv key
	privs, seqs = []crypto.PrivKey{transaction2}, []uint64{1}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	// test wrong sig number
	privs, seqs = []crypto.PrivKey{transaction2, transaction1}, []uint64{2, 0}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, ErrWrongNumberOfSigners().Result())
}

// Test grant authentication.
func TestGrantAuthenticationTx(t *testing.T) {
	am, _, ph, ctx, anteHandler := setupTest()
	// keys and username
	_, transaction1, _, user1 := createTestAccount(ctx, am, ph, "user1")
	_, transaction2, post2, user2 := createTestAccount(ctx, am, ph, "user2")
	_, transaction3, post3, user3 := createTestAccount(ctx, am, ph, "user3")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1)

	// test valid transaction
	privs, seqs := []crypto.PrivKey{transaction1}, []uint64{0}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)
	seq, err := am.GetSequence(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, seq, uint64(1))

	// test wrong priv key
	privs, seqs = []crypto.PrivKey{transaction2}, []uint64{1}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	privs, seqs = []crypto.PrivKey{post2}, []uint64{1}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	err = am.AuthorizePermission(ctx, user1, user2, 3600, types.AppPermission, types.NewCoinFromInt64(0))
	assert.Nil(t, err)

	// should still fail by using transaction key
	privs, seqs = []crypto.PrivKey{transaction2}, []uint64{1}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	// should pass authentication check after grant the app permission
	privs, seqs = []crypto.PrivKey{post2}, []uint64{1}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)
	seq, err = am.GetSequence(ctx, user2)
	assert.Nil(t, err)
	assert.Equal(t, seq, uint64(0))
	seq, err = am.GetSequence(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, seq, uint64(2))

	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: ctx.BlockHeader().Time.Add(time.Duration(3601) * time.Second)})
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	// test pre authorization permission
	err = am.AuthorizePermission(ctx, user1, user3, 3600, types.PreAuthorizationPermission, types.NewCoinFromInt64(100))
	assert.Nil(t, err)
	msg.Permission = types.PreAuthorizationPermission
	privs, seqs = []crypto.PrivKey{post3}, []uint64{2}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	privs, seqs = []crypto.PrivKey{transaction3}, []uint64{2}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)
	seq, err = am.GetSequence(ctx, user3)
	assert.Nil(t, err)
	assert.Equal(t, seq, uint64(0))
	seq, err = am.GetSequence(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, seq, uint64(3))

	// test pre authorization exceeds limitation
	msg.Amount = types.NewCoinFromInt64(100)
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(
		t, anteHandler, ctx, tx,
		acc.ErrPreAuthAmountInsufficient(
			user3, msg.Amount.Minus(types.NewCoinFromInt64(10)), msg.Amount).Result())

}

// Test various error cases in the AnteHandler control flow.
func TestTPSCapacity(t *testing.T) {
	am, gm, ph, ctx, anteHandler := setupTest()
	// keys and username
	_, transaction1, _, user1 := createTestAccount(ctx, am, ph, "user1")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1)

	// test valid transaction
	privs, seqs := []crypto.PrivKey{transaction1}, []uint64{0}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)

	seq, err := am.GetSequence(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, seq, uint64(1))

	ctx = ctx.WithBlockHeader(
		abci.Header{ChainID: "Lino", Height: 2, Time: time.Now(), NumTxs: 1000})
	gm.SetLastBlockTime(ctx, time.Now().Unix()-1)
	gm.UpdateTPS(ctx)

	seqs = []uint64{1}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)
	seqs = []uint64{2}
	tx = newTestTx(ctx, []sdk.Msg{msg}, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrAccountTPSCapacityNotEnough(user1).Result())
}
