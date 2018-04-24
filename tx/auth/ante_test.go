package auth

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

func createTestAccount(
	ctx sdk.Context, am acc.AccountManager, username string) (crypto.PrivKey, types.AccountKey) {
	priv := crypto.GenPrivKeyEd25519()
	am.AddCoinToAddress(ctx, priv.PubKey().Address(), types.NewCoin(100))
	am.CreateAccount(ctx, types.AccountKey(username), priv.PubKey(), types.NewCoin(0))
	return priv.Wrap(), types.AccountKey(username)
}

func InitGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoin(10000*types.Decimals))
}

func setupTest() (acc.AccountManager, global.GlobalManager, sdk.Context, sdk.AnteHandler) {
	db := dbm.NewMemDB()
	accountCapKey := sdk.NewKVStoreKey("account")
	globalCapKey := sdk.NewKVStoreKey("global")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accountCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(globalCapKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: 1, Time: time.Now().Unix()}, false, nil)
	am := acc.NewAccountManager(accountCapKey)
	gm := global.NewGlobalManager(globalCapKey)
	InitGlobalManager(ctx, gm)
	anteHandler := NewAnteHandler(am, gm)

	return am, gm, ctx, anteHandler
}

type TestMsg struct {
	signers []types.AccountKey
}

func (msg *TestMsg) Type() string                            { return "normal msg" }
func (msg *TestMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg *TestMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.signers)
	if err != nil {
		panic(err)
	}
	return bz
}
func (msg *TestMsg) ValidateBasic() sdk.Error { return nil }
func (msg *TestMsg) GetSigners() []sdk.Address {
	addrs := make([]sdk.Address, len(msg.signers))
	for i, signer := range msg.signers {
		addrs[i] = sdk.Address(signer)
	}
	return addrs
}

func newTestMsg(accKeys ...types.AccountKey) *TestMsg {
	return &TestMsg{
		signers: accKeys,
	}
}

type RegisterTestMsg struct {
	Register sdk.Address
}

func (msg *RegisterTestMsg) Type() string {
	return types.RegisterRouterName
}

func (msg *RegisterTestMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg *RegisterTestMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.Register)
	if err != nil {
		panic(err)
	}
	return bz
}
func (msg *RegisterTestMsg) ValidateBasic() sdk.Error { return nil }
func (msg *RegisterTestMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Register}
}

func newRegisterTestMsg(addr sdk.Address) *RegisterTestMsg {
	return &RegisterTestMsg{
		Register: addr,
	}
}

// run the tx through the anteHandler and ensure its valid
func checkValidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx) {
	_, result, abort := anteHandler(ctx, tx)
	assert.False(t, abort)
	fmt.Println(result)
	assert.Equal(t, sdk.CodeOK, result.Code)
	assert.True(t, result.IsOK())
}

// run the tx through the anteHandler and ensure it fails with the given code
func checkInvalidTx(
	t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, result sdk.Result) {
	_, r, abort := anteHandler(ctx, tx)
	assert.True(t, abort)
	assert.Equal(t, result, r)
}

func newTestTx(ctx sdk.Context, msg sdk.Msg, privs []crypto.PrivKey, seqs []int64) sdk.Tx {
	signBytes := sdk.StdSignBytes(ctx.ChainID(), seqs, sdk.StdFee{}, msg)
	return newTestTxWithSignBytes(msg, privs, seqs, signBytes)
}

func newTestTxWithSignBytes(
	msg sdk.Msg, privs []crypto.PrivKey, seqs []int64, signBytes []byte) sdk.Tx {
	sigs := make([]sdk.StdSignature, len(privs))
	for i, priv := range privs {
		sigs[i] = sdk.StdSignature{
			PubKey: priv.PubKey(), Signature: priv.Sign(signBytes), Sequence: seqs[i]}
	}
	tx := sdk.NewStdTx(msg, sdk.StdFee{}, sigs)
	return tx
}

// Test various error cases in the AnteHandler control flow.
func TestAnteHandlerSigErrors(t *testing.T) {
	// setup
	am, _, ctx, anteHandler := setupTest()
	// get private key and username
	priv1, user1 := createTestAccount(ctx, am, "user1")
	priv2, user2 := createTestAccount(ctx, am, "user2")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1, user2)

	// test no signatures
	privs, seqs := []crypto.PrivKey{}, []int64{}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("no signers").Result())

	// test num sigs less than GetSigners
	privs, seqs = []crypto.PrivKey{priv1}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(
		t, anteHandler, ctx, tx, sdk.ErrUnauthorized("wrong number of signers").Result())

	// test sig user mismatch
	privs, seqs = []crypto.PrivKey{priv2, priv1}, []int64{0, 0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())
}

// Test various error cases in the AnteHandler control flow.
func TestAnteHandlerRegisterTx(t *testing.T) {
	am, _, ctx, anteHandler := setupTest()
	priv1 := crypto.GenPrivKeyEd25519().Wrap()
	priv2 := crypto.GenPrivKeyEd25519().Wrap()
	err := am.AddCoinToAddress(ctx, priv1.PubKey().Address(), types.NewCoin(0))
	assert.Nil(t, err)

	// msg and signatures
	var tx sdk.Tx
	msg := newRegisterTestMsg(priv1.PubKey().Address())

	// test valid transaction
	privs, seqs := []crypto.PrivKey{priv1}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)

	// test no signatures
	privs, seqs = []crypto.PrivKey{}, []int64{}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("no signers").Result())

	// test wrong priv key
	privs, seqs = []crypto.PrivKey{priv2}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(
		t, anteHandler, ctx, tx, sdk.ErrUnauthorized("wrong public key for signer").Result())

	// test wrong sig number
	privs, seqs = []crypto.PrivKey{priv2, priv1}, []int64{0, 0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(
		t, anteHandler, ctx, tx, sdk.ErrUnauthorized("wrong number of signers").Result())
}

// Test various error cases in the AnteHandler control flow.
func TestAnteHandlerNormalTx(t *testing.T) {
	am, _, ctx, anteHandler := setupTest()
	// keys and username
	priv1, user1 := createTestAccount(ctx, am, "user1")
	priv2, _ := createTestAccount(ctx, am, "user2")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1)

	// test valid transaction
	privs, seqs := []crypto.PrivKey{priv1}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)
	seq, err := am.GetSequence(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, seq, int64(1))

	// test no signatures
	privs, seqs = []crypto.PrivKey{}, []int64{}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("no signers").Result())

	// test wrong sequence number
	privs, seqs = []crypto.PrivKey{priv1}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrInvalidSequence(
		fmt.Sprintf("Invalid sequence for signer %v. Got %d, expected %d", user1, 0, 1)).Result())

	// test wrong priv key
	privs, seqs = []crypto.PrivKey{priv2}, []int64{1}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	// test wrong sig number
	privs, seqs = []crypto.PrivKey{priv2, priv1}, []int64{2, 0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())
}

// Test grant authentication.
func TestGrantAuthenticationTx(t *testing.T) {
	am, _, ctx, anteHandler := setupTest()
	// keys and username
	priv1, user1 := createTestAccount(ctx, am, "user1")
	priv2, user2 := createTestAccount(ctx, am, "user2")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1)

	// test valid transaction
	privs, seqs := []crypto.PrivKey{priv1}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)
	seq, err := am.GetSequence(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, seq, int64(1))

	// test wrong priv key
	privs, seqs = []crypto.PrivKey{priv2}, []int64{1}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	err = am.GrantPubKeyToUser(ctx, user1, user2, 3600, 1)
	assert.Nil(t, err)

	// should pass authentication check after grant
	privs, seqs = []crypto.PrivKey{priv2}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)
	seq, err = am.GetSequence(ctx, user2)
	assert.Nil(t, err)
	assert.Equal(t, seq, int64(1))

	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: ctx.BlockHeader().Time + 3601})
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())
}

// Test various error cases in the AnteHandler control flow.
func TestTPSCapacity(t *testing.T) {
	am, gm, ctx, anteHandler := setupTest()
	// keys and username
	priv1, user1 := createTestAccount(ctx, am, "user1")
	_, _ = createTestAccount(ctx, am, "user2")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1)

	// test valid transaction
	privs, seqs := []crypto.PrivKey{priv1}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)
	seq, err := am.GetSequence(ctx, user1)
	assert.Nil(t, err)
	assert.Equal(t, seq, int64(1))

	ctx = ctx.WithBlockHeader(
		abci.Header{ChainID: "Lino", Height: 2, Time: time.Now().Unix(), NumTxs: 1000})
	gm.UpdateTPS(ctx, time.Now().Unix()-1)
	seqs = []int64{1}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrAccountTPSCapacityNotEnough(user1).Result())
	seqs = []int64{2}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, acc.ErrAccountTPSCapacityNotEnough(user1).Result())
}
