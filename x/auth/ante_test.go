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
	acctypes "github.com/lino-network/lino/x/account/types"
	bandwidthmock "github.com/lino-network/lino/x/bandwidth/mocks"

	// bandwidthmn "github.com/lino-network/lino/x/bandwidth/manager"
	// dev "github.com/lino-network/lino/x/developer"
	devmn "github.com/lino-network/lino/x/developer/manager"
	"github.com/lino-network/lino/x/global"
	post "github.com/lino-network/lino/x/post"
	postmn "github.com/lino-network/lino/x/post/manager"
	pricemn "github.com/lino-network/lino/x/price/manager"
	vote "github.com/lino-network/lino/x/vote"
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
		Permission: types.AppPermission,
		Amount:     types.NewCoinFromInt64(10),
	}
}

func newTestTx(
	ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, seqs []uint64) sdk.Tx {
	sigs := make([]auth.StdSignature, len(privs))

	for i, priv := range privs {
		signBytes := auth.StdSignBytes(ctx.ChainID(), 0, seqs[i], auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(types.LinoCoinDenom, sdk.NewInt(10000000)))}, msgs, "")
		bz, _ := priv.Sign(signBytes)
		sigs[i] = auth.StdSignature{
			PubKey: priv.PubKey(), Signature: bz}
	}
	tx := auth.NewStdTx(msgs, auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(types.LinoCoinDenom, sdk.NewInt(10000000)))}, sigs, "")
	return tx
}

func initGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

type AnteTestSuite struct {
	suite.Suite
	am   acc.AccountKeeper
	pm   post.PostKeeper
	gm   global.GlobalManager
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

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestPostKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestDeveloperKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestBandwidthKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestVoteKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: 1, Time: time.Now()}, false, log.NewNopLogger())

	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)
	gm := global.NewGlobalManager(TestGlobalKVStoreKey, ph)

	am := accmn.NewAccountManager(TestAccountKVStoreKey, ph, &gm)
	vm := vote.NewVoteManager(TestVoteKVStoreKey, ph)
	price := pricemn.TestnetPriceManager{}
	dm := devmn.NewDeveloperManager(TestDeveloperKVStoreKey, ph, vm, am, price, &gm)
	pm := postmn.NewPostManager(TestPostKVStoreKey, am, &gm, dm, nil, price)

	bm := &bandwidthmock.BandwidthKeeper{}

	initGlobalManager(ctx, gm)
	anteHandler := NewAnteHandler(am, bm)

	suite.am = am
	suite.pm = pm
	suite.gm = gm
	suite.ph = ph
	suite.ctx = ctx
	suite.ante = anteHandler
	bm.On("CheckBandwidth", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
}

func (suite *AnteTestSuite) createTestAccount(username string) (secp256k1.PrivKeySecp256k1, secp256k1.PrivKeySecp256k1, types.AccountKey) {
	signingKey := secp256k1.GenPrivKey()
	transactionKey := secp256k1.GenPrivKey()
	accParams, _ := suite.ph.GetAccountParam(suite.ctx)
	suite.am.CreateAccount(suite.ctx, types.AccountKey(username),
		signingKey.PubKey(), transactionKey.PubKey())
	suite.am.AddCoinToUsername(suite.ctx, types.AccountKey(username), accParams.RegisterFee)
	return signingKey, transactionKey, types.AccountKey(username)
}

func (suite *AnteTestSuite) createTestPost(postid string, author types.AccountKey) {
	msg := post.CreatePostMsg{
		PostID:    postid,
		Title:     "testTitle",
		Content:   "qqqqqqq",
		Author:    author,
		CreatedBy: author,
	}
	err := suite.pm.CreatePost(suite.ctx, msg.Author, msg.PostID, msg.CreatedBy, msg.Content, msg.Title)
	suite.Require().Nil(err)
}

// run the tx through the anteHandler and ensure its valid
func (suite *AnteTestSuite) checkValidTx(tx sdk.Tx) {
	_, result, abort := suite.ante(suite.ctx, tx, false)
	suite.Assert().False(abort)
	suite.Assert().True(result.Code.IsOK()) // redundent
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

// Test grant authentication.
func (suite *AnteTestSuite) TestGrantAuthenticationTx() {
	// keys and username
	_, transaction1, user1 := suite.createTestAccount("user1")
	_, transaction2, user2 := suite.createTestAccount("user2")
	_, transaction3, user3 := suite.createTestAccount("user3")

	// msg and signatures
	var tx sdk.Tx
	msg := newTestMsg(user1)

	// test valid transaction
	privs, seqs := []crypto.PrivKey{transaction1}, []uint64{0}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkValidTx(tx)
	addr1, err := suite.am.GetAddress(suite.ctx, user1)
	suite.Nil(err)
	addr2, err := suite.am.GetAddress(suite.ctx, user2)
	suite.Nil(err)
	seq, err := suite.am.GetSequence(suite.ctx, addr1)
	suite.Nil(err)
	suite.Equal(seq, uint64(1))

	// test wrong priv key
	privs, seqs = []crypto.PrivKey{transaction2}, []uint64{1}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(tx, acctypes.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	privs, seqs = []crypto.PrivKey{transaction3}, []uint64{1}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(tx, acctypes.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	err = suite.am.AuthorizePermission(suite.ctx, user1, user2, 3600, types.AppPermission, types.NewCoinFromInt64(0))
	suite.Nil(err)

	privs, seqs = []crypto.PrivKey{transaction2}, []uint64{1}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkValidTx(tx)

	// should pass authentication check after grant the app permission
	// privs, seqs = []crypto.PrivKey{post2}, []uint64{1}
	// tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	// suite.checkValidTx(tx)
	seq, err = suite.am.GetSequence(suite.ctx, addr2)
	suite.Nil(err)
	suite.Equal(seq, uint64(0))
	seq, err = suite.am.GetSequence(suite.ctx, addr1)
	suite.Nil(err)
	suite.Equal(seq, uint64(2))

	suite.ctx = suite.ctx.WithBlockHeader(abci.Header{
		ChainID: "Lino", Height: 2,
		Time: suite.ctx.BlockHeader().Time.Add(time.Duration(3601) * time.Second)})
	suite.checkInvalidTx(tx, acctypes.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	// test pre authorization permission
	err = suite.am.AuthorizePermission(suite.ctx, user1, user3, 3600, types.PreAuthorizationPermission, types.NewCoinFromInt64(100))
	suite.Nil(err)
	msg.Permission = types.PreAuthorizationPermission
	msg.Amount = types.NewCoinFromInt64(100)
	// privs, seqs = []crypto.PrivKey{post3}, []uint64{2}
	// tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	// suite.checkInvalidTx(tx, acc.ErrCheckAuthenticatePubKeyOwner(user1).Result())

	privs, seqs = []crypto.PrivKey{transaction3}, []uint64{2}
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkValidTx(tx)
	// seq, err = suite.am.GetSequence(suite.ctx, user3)
	// suite.Nil(err)
	// suite.Equal(seq, uint64(0))
	// seq, err = suite.am.GetSequence(suite.ctx, user1)
	// suite.Nil(err)
	// suite.Equal(seq, uint64(3))

	// test pre authorization exceeds limitation
	tx = newTestTx(suite.ctx, []sdk.Msg{msg}, privs, seqs)
	suite.checkInvalidTx(
		tx,
		acctypes.ErrPreAuthAmountInsufficient(
			user3, types.NewCoinFromInt64(0), msg.Amount).Result())
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, &AnteTestSuite{})
}
