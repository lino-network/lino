package auth

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return ms, capKey
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

func (msg *RegisterTestMsg) Type() string                            { return types.RegisterRouterName }
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

// generate a priv key and return it with its address
func privAndBank() (crypto.PrivKey, *types.AccountBank) {
	priv := crypto.GenPrivKeyEd25519()
	accBank := &types.AccountBank{
		Address: priv.PubKey().Address(),
		Coins:   sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	return priv.Wrap(), accBank
}

// run the tx through the anteHandler and ensure its valid
func checkValidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx) {
	_, result, abort := anteHandler(ctx, tx)
	assert.False(t, abort)
	assert.Equal(t, sdk.CodeOK, result.Code)
	assert.True(t, result.IsOK())
}

// run the tx through the anteHandler and ensure it fails with the given code
func checkInvalidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, result sdk.Result) {
	_, r, abort := anteHandler(ctx, tx)
	assert.True(t, abort)
	assert.Equal(t, result, r)
	fmt.Println(r)
}

func newTestTx(ctx sdk.Context, msg sdk.Msg, privs []crypto.PrivKey, seqs []int64) sdk.Tx {
	signBytes := sdk.StdSignBytes(ctx.ChainID(), seqs, msg)
	return newTestTxWithSignBytes(msg, privs, seqs, signBytes)
}

func newTestTxWithSignBytes(msg sdk.Msg, privs []crypto.PrivKey, seqs []int64, signBytes []byte) sdk.Tx {
	sigs := make([]sdk.StdSignature, len(privs))
	for i, priv := range privs {
		sigs[i] = sdk.StdSignature{PubKey: priv.PubKey(), Signature: priv.Sign(signBytes), Sequence: seqs[i]}
	}
	tx := sdk.NewStdTx(msg, sigs)
	return tx
}

// Test various error cases in the AnteHandler control flow.
func TestAnteHandlerSigErrors(t *testing.T) {
	// setup
	ms, capKey := setupMultiStore()
	lam := acc.NewLinoAccountManager(capKey)
	anteHandler := NewAnteHandler(lam)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "mychainid"}, false, nil)

	// keys and addresses
	priv1, bank1 := privAndBank()
	priv2, bank2 := privAndBank()
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")

	_, err := lam.CreateAccount(ctx, user1, priv1.PubKey(), bank1)
	assert.Nil(t, err)
	_, err = lam.CreateAccount(ctx, user2, priv2.PubKey(), bank2)
	assert.Nil(t, err)

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
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("wrong number of signers").Result())

	// test sig user mismatch
	privs, seqs = []crypto.PrivKey{priv2, priv1}, []int64{0, 0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("signer mismatch").Result())
}

// Test various error cases in the AnteHandler control flow.
func TestAnteHandlerRegisterTx(t *testing.T) {
	// setup
	ms, capKey := setupMultiStore()
	lam := acc.NewLinoAccountManager(capKey)
	anteHandler := NewAnteHandler(lam)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "mychainid"}, false, nil)

	// keys and addresses
	priv1, bank1 := privAndBank()
	priv2, _ := privAndBank()
	// user1 := types.AccountKey("user1")
	// user2 := types.AccountKey("user2")

	err := lam.SetBank(ctx, priv1.PubKey().Address(), bank1)
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
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("wrong public key for signer").Result())

	// test wrong sig number
	privs, seqs = []crypto.PrivKey{priv2, priv1}, []int64{0, 0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("wrong number of signers").Result())
}

// Test various error cases in the AnteHandler control flow.
func TestAnteHandlerNormalTx(t *testing.T) {
	// setup
	ms, capKey := setupMultiStore()
	lam := acc.NewLinoAccountManager(capKey)
	anteHandler := NewAnteHandler(lam)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "mychainid"}, false, nil)

	// keys and addresses
	priv1, bank1 := privAndBank()
	priv2, bank2 := privAndBank()
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")

	_, err := lam.CreateAccount(ctx, user1, priv1.PubKey(), bank1)
	assert.Nil(t, err)
	_, err = lam.CreateAccount(ctx, user2, priv2.PubKey(), bank2)
	assert.Nil(t, err)

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1)

	// test valid transaction
	privs, seqs := []crypto.PrivKey{priv1}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkValidTx(t, anteHandler, ctx, tx)

	// test no signatures
	privs, seqs = []crypto.PrivKey{}, []int64{}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("no signers").Result())

	// test wrong sequence number
	privs, seqs = []crypto.PrivKey{priv1}, []int64{0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrInvalidSequence(
		fmt.Sprintf("Invalid sequence. Got %d, expected %d", 0, 1)).Result())

	// test wrong priv key
	privs, seqs = []crypto.PrivKey{priv2}, []int64{1}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("signer mismatch").Result())

	// test wrong sig number
	privs, seqs = []crypto.PrivKey{priv2, priv1}, []int64{1, 0}
	tx = newTestTx(ctx, msg, privs, seqs)
	checkInvalidTx(t, anteHandler, ctx, tx, sdk.ErrUnauthorized("signer mismatch").Result())
}


