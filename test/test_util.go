package test

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/param"
	acc "github.com/lino-network/lino/tx/account"
	post "github.com/lino-network/lino/tx/post"
	reg "github.com/lino-network/lino/tx/register"
	val "github.com/lino-network/lino/tx/validator"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

// construct some global keys and addrs.
var (
	GenesisUser            = "genesis"
	GenesisPriv            = crypto.GenPrivKeyEd25519()
	GenesisTransactionPriv = crypto.GenPrivKeyEd25519()
	GenesisPostPriv        = crypto.GenPrivKeyEd25519()
	GenesisAddr            = GenesisPriv.PubKey().Address()

	DefaultNumOfVal  int       = 21
	GenesisTotalLino types.LNO = "10000000000"
	LNOPerValidator  types.LNO = "100000000"

	PenaltyMissVote       types.Coin = types.NewCoin(20000 * types.Decimals)
	ChangeParamMinDeposit types.Coin = types.NewCoin(100000 * types.Decimals)

	ProposalDecideHr            int64   = 24 * 7
	ParamChangeHr               int64   = 24
	CoinReturnIntervalHr        int64   = 24 * 7
	CoinReturnTimes             int64   = 7
	ConsumptionFrictionRate     sdk.Rat = sdk.NewRat(5, 100)
	ConsumptionFreezingPeriodHr int64   = 24 * 7
)

func loggerAndDB() (log.Logger, dbm.DB) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return logger, db
}

func NewTestLinoBlockchain(t *testing.T, numOfValidators int) *app.LinoBlockchain {
	logger, db := loggerAndDB()
	lb := app.NewLinoBlockchain(logger, db)

	genesisState := genesis.GenesisState{
		Accounts:  []genesis.GenesisAccount{},
		TotalLino: GenesisTotalLino,
	}

	// Generate 21 validators
	for i := 0; i < numOfValidators; i++ {
		genesisAcc := genesis.GenesisAccount{
			Name:           "validator" + strconv.Itoa(i),
			Lino:           LNOPerValidator,
			MasterKey:      crypto.GenPrivKeyEd25519().PubKey(),
			TransactionKey: crypto.GenPrivKeyEd25519().PubKey(),
			PostKey:        crypto.GenPrivKeyEd25519().PubKey(),
			IsValidator:    true,
			ValPubKey:      crypto.GenPrivKeyEd25519().PubKey(),
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}

	totalAmt, _ := strconv.ParseInt(GenesisTotalLino, 10, 64)
	validatorAmt, _ := strconv.ParseInt(LNOPerValidator, 10, 64)
	initLNO := strconv.FormatInt(totalAmt-validatorAmt, 10)
	genesisAcc := genesis.GenesisAccount{
		Name:           GenesisUser,
		Lino:           initLNO,
		MasterKey:      GenesisPriv.PubKey(),
		TransactionKey: GenesisTransactionPriv.PubKey(),
		PostKey:        GenesisPostPriv.PubKey(),
		IsValidator:    false,
		ValPubKey:      GenesisPriv.PubKey(),
	}
	genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	result, err := genesis.GetGenesisJson(genesisState)
	assert.Nil(t, err)

	vals := []abci.Validator{}
	lb.InitChain(abci.RequestInitChain{vals, json.RawMessage(result)})
	lb.Commit()

	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{ChainID: "Lino"}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	return lb
}

func CheckGlobalAllocation(t *testing.T, lb *app.LinoBlockchain, expectAllocation param.GlobalAllocationParam) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	allocation, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectAllocation, *allocation)
}

func CheckBalance(t *testing.T, accountName string, lb *app.LinoBlockchain, expectBalance types.Coin) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	accManager := acc.NewAccountManager(lb.CapKeyAccountStore, ph)
	saving, err :=
		accManager.GetSavingFromBank(ctx, types.AccountKey(accountName))
	assert.Nil(t, err)
	assert.Equal(t, expectBalance, saving)
}

func CheckValidatorDeposit(t *testing.T, accountName string, lb *app.LinoBlockchain, expectDeposit types.Coin) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	deposit, err := valManager.GetValidatorDeposit(ctx, types.AccountKey(accountName))
	assert.Nil(t, err)
	assert.Equal(t, expectDeposit, deposit)
}

func CheckOncallValidatorList(
	t *testing.T, accountName string, isInOnCallValidatorList bool, lb *app.LinoBlockchain) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	lst, err := valManager.GetValidatorList(ctx)
	assert.Nil(t, err)
	index := val.FindAccountInList(types.AccountKey(accountName), lst.OncallValidators)
	if isInOnCallValidatorList {
		assert.True(t, index > -1)
	} else {
		assert.True(t, index == -1)
	}

}

func CheckAllValidatorList(
	t *testing.T, accountName string, isInAllValidatorList bool, lb *app.LinoBlockchain) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	lst, err := valManager.GetValidatorList(ctx)

	assert.Nil(t, err)
	index := val.FindAccountInList(types.AccountKey(accountName), lst.AllValidators)
	if isInAllValidatorList {
		assert.True(t, index > -1)
	} else {
		assert.True(t, index == -1)
	}
}

func CreateAccount(
	t *testing.T, accountName string, lb *app.LinoBlockchain, seq int64,
	masterPriv crypto.PrivKeyEd25519, transactionPriv crypto.PrivKeyEd25519, postPriv crypto.PrivKeyEd25519,
	numOfLino string) {

	transferMsg := acc.NewTransferMsg(
		GenesisUser, types.LNO(numOfLino),
		"", acc.TransferToAddr(masterPriv.PubKey().Address()))

	SignCheckDeliver(t, lb, transferMsg, seq, true, GenesisTransactionPriv, time.Now().Unix())

	registerMsg := reg.NewRegisterMsg(accountName, masterPriv.PubKey(), transactionPriv.PubKey(), postPriv.PubKey())
	SignCheckDeliver(t, lb, registerMsg, 0, true, masterPriv, time.Now().Unix())
}

func GetGenesisAccountCoin(numOfValidator int) types.Coin {
	totalAmt, _ := strconv.ParseInt(GenesisTotalLino, 10, 64)
	validatorAmt, _ := strconv.ParseInt(LNOPerValidator, 10, 64)
	initLNO := strconv.FormatInt(totalAmt-validatorAmt, 10)
	initCoin, _ := types.LinoToCoin(initLNO)
	return initCoin
}

func SignCheckDeliver(t *testing.T, lb *app.LinoBlockchain, msg sdk.Msg, seq int64,
	expPass bool, priv crypto.PrivKeyEd25519, headTime int64) {
	// Sign the tx
	tx := genTx(msg, seq, priv)
	// Run a Check
	res := lb.Check(tx)
	if expPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}

	// Simulate a Block
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			ChainID: "Lino", Time: headTime}})
	res = lb.Deliver(tx)
	if expPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
}

func SimulateOneBlock(lb *app.LinoBlockchain, headTime int64) {
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			ChainID: "Lino", Time: headTime}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
}

func genTx(msg sdk.Msg, seq int64, priv crypto.PrivKeyEd25519) sdk.StdTx {
	sigs := []sdk.StdSignature{{
		PubKey:    priv.PubKey(),
		Signature: priv.Sign(sdk.StdSignBytes("Lino", []int64{seq}, sdk.StdFee{}, msg)),
		Sequence:  seq}}

	return sdk.NewStdTx(msg, sdk.StdFee{}, sigs)

}

func CreateTestPost(
	t *testing.T, lb *app.LinoBlockchain,
	username, postID string, seq int64, priv crypto.PrivKeyEd25519,
	sourceAuthor, sourcePostID string,
	parentAuthor, parentPostID string,
	redistributionSplitRate string, publishTime int64) {

	postCreateParams := post.PostCreateParams{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       types.AccountKey(username),
		ParentAuthor: types.AccountKey(parentAuthor),
		ParentPostID: parentPostID,
		SourceAuthor: types.AccountKey(sourceAuthor),
		SourcePostID: sourcePostID,
		Links:        []types.IDToURLMapping{},
		RedistributionSplitRate: redistributionSplitRate,
	}
	msg := post.NewCreatePostMsg(postCreateParams)
	SignCheckDeliver(t, lb, msg, seq, true, priv, publishTime)
}

func CoinToString(coin types.Coin) string {
	return strconv.FormatInt(coin.ToInt64()/types.Decimals, 10)
}
