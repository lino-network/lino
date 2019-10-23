//nolint:unused,deadcode
package auth

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	accmn "github.com/lino-network/lino/x/account/manager"
	accmodel "github.com/lino-network/lino/x/account/model"
	acctypes "github.com/lino-network/lino/x/account/types"
	bandwidthmock "github.com/lino-network/lino/x/bandwidth/mocks"
)

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
		Permission: types.TransactionPermission,
		Amount:     types.NewCoinFromInt64(10),
	}
}

func newTestTx(
	ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, seqs []uint64) sdk.Tx {
	sigs := make([]auth.StdSignature, len(privs))

	for i, priv := range privs {
		signBytes := auth.StdSignBytes(
			ctx.ChainID(), 0, seqs[i],
			auth.StdFee{
				Amount: sdk.NewCoins(sdk.NewCoin(types.LinoCoinDenom, sdk.NewInt(10000000))),
			},
			msgs, "")
		bz, _ := priv.Sign(signBytes)
		sigs[i] = auth.StdSignature{
			PubKey: priv.PubKey(), Signature: bz}
	}
	tx := auth.NewStdTx(
		msgs,
		auth.StdFee{
			Amount: sdk.NewCoins(sdk.NewCoin(types.LinoCoinDenom, sdk.NewInt(10000000))),
		}, sigs, "")
	return tx
}

type AnteTestSuite struct {
	suite.Suite
	am   acc.AccountKeeper
	ph   param.ParamHolder
	ctx  sdk.Context
	ante sdk.AnteHandler
}

func (suite *AnteTestSuite) SetupTest() {
	TestAccountKVStoreKey := sdk.NewKVStoreKey("account")
	TestPostKVStoreKey := sdk.NewKVStoreKey("post")
	TestGlobalKVStoreKey := sdk.NewKVStoreKey("global")
	TestParamKVStoreKey := sdk.NewKVStoreKey("param")
	TestDeveloperKVStoreKey := sdk.NewKVStoreKey("dev")
	TestBandwidthKVStoreKey := sdk.NewKVStoreKey("bandwidth")
	TestVoteKVStoreKey := sdk.NewKVStoreKey("vote")
	TestPriceKVStoreKey := sdk.NewKVStoreKey("price")
	TestValidatorKVStoreKey := sdk.NewKVStoreKey("validator")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestPostKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestDeveloperKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestBandwidthKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestVoteKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestPriceKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestValidatorKVStoreKey, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	ctx := sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: 1, Time: time.Unix(0, 0)}, false, log.NewNopLogger())

	ph := param.NewParamHolder(TestParamKVStoreKey)
	err = ph.InitParam(ctx)
	if err != nil {
		panic(err)
	}
	am := accmn.NewAccountManager(TestAccountKVStoreKey, ph)
	am.InitGenesis(ctx, types.MustLinoToCoin("10000000000"), []accmodel.Pool{
		{
			Name:    types.AccountVestingPool,
			Balance: types.MustLinoToCoin("10000000000"),
		},
	})

	bm := &bandwidthmock.BandwidthKeeper{}
	bm.On("CheckBandwidth", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	anteHandler := NewAnteHandler(am, bm)

	suite.am = am
	suite.ph = ph
	suite.ctx = ctx
	suite.ante = anteHandler
}

func (suite *AnteTestSuite) createTestAccount(username string) (secp256k1.PrivKeySecp256k1, secp256k1.PrivKeySecp256k1, types.AccountKey) {
	signingKey := secp256k1.GenPrivKey()
	transactionKey := secp256k1.GenPrivKey()
	accParams := suite.ph.GetAccountParam(suite.ctx)
	err := suite.am.GenesisAccount(suite.ctx, types.AccountKey(username),
		signingKey.PubKey(), transactionKey.PubKey())
	if err != nil {
		panic(err)
	}
	err = suite.am.MoveFromPool(suite.ctx,
		types.AccountVestingPool,
		types.NewAccOrAddrFromAcc(types.AccountKey(username)), accParams.RegisterFee)
	if err != nil {
		panic(err)
	}
	return signingKey, transactionKey, types.AccountKey(username)
}

// run the tx through the anteHandler and ensure its valid
func (suite *AnteTestSuite) checkValidTx(tx sdk.Tx) {
	_, result, abort := suite.ante(suite.ctx, tx, false)
	suite.Assert().False(abort)
	suite.Assert().True(result.Code.IsOK()) // redundant
	suite.Assert().True(result.IsOK())
}

// run the tx through the anteHandler and ensure it fails with the given code
func (suite *AnteTestSuite) checkInvalidTx(tx sdk.Tx, result sdk.Result) {
	_, r, abort := suite.ante(suite.ctx, tx, false)

	suite.Assert().True(abort)
	suite.Assert().Equal(result, r)
}

// Test various error cases in the AnteHandler control flow.
func (suite *AnteTestSuite) TestAnteHandlerSigErrors() {
	// get private key and username
	_, transaction1, user1 := suite.createTestAccount("user1")
	_, transaction2, user2 := suite.createTestAccount("user2")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1, user2)

	// test no signatures
	privs, seqs := []crypto.PrivKey{}, []uint64{}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(tx, ErrNoSignatures().Result())

	// test num sigs less than GetSigners
	privs, seqs = []crypto.PrivKey{transaction1}, []uint64{0}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(tx, ErrWrongNumberOfSigners().Result())

	// test sig user mismatch
	privs, seqs = []crypto.PrivKey{transaction2, transaction1}, []uint64{0, 0}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(tx, acctypes.ErrCheckAuthenticatePubKeyOwner(user1).Result())
}

// Test various error cases in the AnteHandler control flow.
func (suite *AnteTestSuite) TestAnteHandlerNormalTx() {
	// keys and username
	_, transaction1, user1 := suite.createTestAccount("user1")
	_, transaction2, _ := suite.createTestAccount("user2")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1)

	// test valid transaction
	privs, seqs := []crypto.PrivKey{transaction1}, []uint64{0}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkValidTx(tx)
	addr, err := suite.am.GetAddress(suite.ctx, user1)
	suite.Nil(err)

	seq, err := suite.am.GetSequence(suite.ctx, addr)
	suite.Nil(err)
	suite.Equal(seq, uint64(1))

	// test no signatures
	privs, seqs = []crypto.PrivKey{}, []uint64{}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(tx, ErrNoSignatures().Result())

	// test wrong sequence number, now we return signature failed even it's seq number error.
	privs, seqs = []crypto.PrivKey{transaction1}, []uint64{0}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(tx, ErrUnverifiedBytes(
		"signature verification failed, chain-id:Lino, seq:1").Result())

	// test wrong priv key
	privs, seqs = []crypto.PrivKey{transaction2}, []uint64{1}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(tx, acctypes.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	// test wrong sig number
	privs, seqs = []crypto.PrivKey{transaction2, transaction1}, []uint64{2, 0}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(tx, ErrWrongNumberOfSigners().Result())
}

// TestCheckAccountSigner
func (suite *AnteTestSuite) TestCheckAccountSigner() {
	// keys and username
	_, transaction1, user1 := suite.createTestAccount("user1")
	_, transaction2, _ := suite.createTestAccount("user2")
	privKey := secp256k1.GenPrivKey()
	err := suite.am.MoveFromPool(
		suite.ctx,
		types.AccountVestingPool,
		types.NewAccOrAddrFromAddr(sdk.AccAddress(privKey.PubKey().Address())),
		types.NewCoinFromInt64(1000))
	suite.Nil(err)

	testCases := []struct {
		testName         string
		signer           types.AccountKey
		signKey          crypto.PubKey
		expectSignerAddr sdk.AccAddress
		expectErr        sdk.Error
	}{
		{
			testName:         "get signer from username",
			signer:           user1,
			signKey:          transaction1.PubKey(),
			expectSignerAddr: sdk.AccAddress(transaction1.PubKey().Address()),
			expectErr:        nil,
		},
		{
			testName:         "no permission",
			signer:           user1,
			signKey:          transaction2.PubKey(),
			expectSignerAddr: nil,
			expectErr:        acctypes.ErrCheckAuthenticatePubKeyOwner(user1),
		},
	}

	for _, tc := range testCases {
		signerAddr, err := checkAccountSigner(
			suite.ctx, suite.am, tc.signer, tc.signKey)
		suite.Equal(tc.expectSignerAddr, signerAddr, "%s", tc.testName)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
	}
}

// Test address signer.
func (suite *AnteTestSuite) TestCheckAddrSigner() {
	// keys and username
	privKey := secp256k1.GenPrivKey()
	err := suite.am.MoveFromPool(suite.ctx,
		types.AccountVestingPool,
		types.NewAccOrAddrFromAddr(sdk.AccAddress(privKey.PubKey().Address())),
		types.NewCoinFromInt64(1000))
	suite.Nil(err)

	newPrivKey := secp256k1.GenPrivKey()

	testCases := []struct {
		testName  string
		signer    sdk.AccAddress
		signKey   crypto.PubKey
		expectErr sdk.Error
		isPaid    bool
	}{
		{
			testName:  "get signer from address",
			signer:    sdk.AccAddress(privKey.PubKey().Address()),
			signKey:   privKey.PubKey(),
			expectErr: nil,
		},
		{
			testName: "sign key without bank struct",
			signer:   sdk.AccAddress(newPrivKey.PubKey().Address()),
			signKey:  newPrivKey.PubKey(),
			expectErr: acctypes.ErrAccountBankNotFound(
				sdk.AccAddress(newPrivKey.PubKey().Address())),
		},
		{
			testName:  "sign key without bank struct but paid",
			signer:    sdk.AccAddress(newPrivKey.PubKey().Address()),
			signKey:   newPrivKey.PubKey(),
			expectErr: nil,
			isPaid:    true,
		},
	}

	for _, tc := range testCases {
		err := checkAddrSigner(
			suite.ctx, suite.am, tc.signer, tc.signKey, tc.isPaid)
		suite.Equal(tc.expectErr, err, "%s", tc.testName)
	}
}

// Test multi sig.
func (suite *AnteTestSuite) TestMultiSig() {
	// keys and username
	_, transaction1, user1 := suite.createTestAccount("user1")
	_, transaction2, _ := suite.createTestAccount("user2")
	// _, transaction3, user3 := suite.createTestAccount("user3")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1, user1)

	// test first private key is wrong
	privs, seqs := []crypto.PrivKey{transaction2, transaction1}, []uint64{0, 1}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(
		tx,
		acctypes.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	// test second private key is wrong
	privs, seqs = []crypto.PrivKey{transaction1, transaction2}, []uint64{0, 1}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(
		tx,
		acctypes.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	// test too many sigs
	privs, seqs = []crypto.PrivKey{transaction1, transaction2, transaction1}, []uint64{0, 1, 1}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(
		tx,
		sdk.ErrTooManySignatures("signatures: 3, limit: 2").Result())

	// test valid transaction
	privs, seqs = []crypto.PrivKey{transaction1, transaction1}, []uint64{1, 2}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkValidTx(tx)
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, &AnteTestSuite{})
}
