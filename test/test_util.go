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
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	globalModel "github.com/lino-network/lino/x/global/model"
	post "github.com/lino-network/lino/x/post"
	val "github.com/lino-network/lino/x/validator"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

// construct some global keys and addrs.
var (
	GenesisUser            = "genesis"
	GenesisPriv            = secp256k1.GenPrivKey()
	GenesisTransactionPriv = secp256k1.GenPrivKey()
	GenesisAppPriv         = secp256k1.GenPrivKey()
	GenesisAddr            = GenesisPriv.PubKey().Address()

	DefaultNumOfVal  = 21
	GenesisTotalCoin = types.NewCoinFromInt64(10000000000 * types.Decimals)
	CoinPerValidator = types.NewCoinFromInt64(100000000 * types.Decimals)

	PenaltyMissVote       = types.NewCoinFromInt64(20000 * types.Decimals)
	ChangeParamMinDeposit = types.NewCoinFromInt64(100000 * types.Decimals)

	ProposalDecideSec            int64 = 24 * 7 * 3600
	ParamChangeExecutionSec      int64 = 24 * 3600
	CoinReturnIntervalSec        int64 = 24 * 7 * 3600
	CoinReturnTimes              int64 = 7
	ConsumptionFrictionRate            = sdk.NewRat(5, 100)
	ConsumptionFreezingPeriodSec int64 = 24 * 7 * 3600
	PostIntervalSec              int64 = 600
)

func loggerAndDB() (log.Logger, dbm.DB) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return logger, db
}

func NewTestLinoBlockchain(t *testing.T, numOfValidators int) *app.LinoBlockchain {
	logger, db := loggerAndDB()
	lb := app.NewLinoBlockchain(logger, db, nil)

	genesisState := app.GenesisState{
		Accounts: []app.GenesisAccount{},
	}

	// Generate 21 validators
	for i := 0; i < numOfValidators; i++ {
		genesisAcc := app.GenesisAccount{
			Name:           "validator" + strconv.Itoa(i),
			Coin:           CoinPerValidator,
			ResetKey:       secp256k1.GenPrivKey().PubKey(),
			TransactionKey: secp256k1.GenPrivKey().PubKey(),
			AppKey:         secp256k1.GenPrivKey().PubKey(),
			IsValidator:    true,
			ValPubKey:      secp256k1.GenPrivKey().PubKey(),
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}

	initLNO := GetGenesisAccountCoin(numOfValidators)
	genesisAcc := app.GenesisAccount{
		Name:           GenesisUser,
		Coin:           initLNO,
		ResetKey:       GenesisPriv.PubKey(),
		TransactionKey: GenesisTransactionPriv.PubKey(),
		AppKey:         GenesisAppPriv.PubKey(),
		IsValidator:    false,
		ValPubKey:      GenesisPriv.PubKey(),
	}
	cdc := app.MakeCodec()
	genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	genesisState.InitGlobalMeta = globalModel.InitParamList{
		MaxTPS:                       sdk.NewRat(1000),
		ConsumptionFreezingPeriodSec: 7 * 24 * 3600,
		ConsumptionFrictionRate:      sdk.NewRat(5, 100),
	}
	result, err := wire.MarshalJSONIndent(cdc, genesisState)
	assert.Nil(t, err)

	lb.InitChain(abci.RequestInitChain{ChainId: "Lino", AppStateBytes: json.RawMessage(result)})
	lb.Commit()

	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{ChainID: "Lino", Time: time.Now()}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	return lb
}

// CheckGlobalAllocation - check global allocation parameter
func CheckGlobalAllocation(t *testing.T, lb *app.LinoBlockchain, expectAllocation param.GlobalAllocationParam) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	allocation, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectAllocation, *allocation)
}

// CheckBalance - check account balance
func CheckBalance(t *testing.T, accountName string, lb *app.LinoBlockchain, expectBalance types.Coin) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	accManager := acc.NewAccountManager(lb.CapKeyAccountStore, ph)
	saving, err :=
		accManager.GetSavingFromBank(ctx, types.AccountKey(accountName))
	assert.Nil(t, err)
	assert.Equal(t, expectBalance, saving)
}

// CheckValidatorDeposit - check validator deposit
func CheckValidatorDeposit(t *testing.T, accountName string, lb *app.LinoBlockchain, expectDeposit types.Coin) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	deposit, err := valManager.GetValidatorDeposit(ctx, types.AccountKey(accountName))
	assert.Nil(t, err)
	assert.Equal(t, expectDeposit, deposit)
}

// CheckOncallValidatorList - check if account is in oncall validator set or not
func CheckOncallValidatorList(
	t *testing.T, accountName string, isInOnCallValidatorList bool, lb *app.LinoBlockchain) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	lst, err := valManager.GetValidatorList(ctx)
	assert.Nil(t, err)
	index := types.FindAccountInList(types.AccountKey(accountName), lst.OncallValidators)
	if isInOnCallValidatorList {
		assert.True(t, index > -1)
	} else {
		assert.True(t, index == -1)
	}

}

// CheckAllValidatorList - check if account is in all validator set or not
func CheckAllValidatorList(
	t *testing.T, accountName string, isInAllValidatorList bool, lb *app.LinoBlockchain) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	lst, err := valManager.GetValidatorList(ctx)

	assert.Nil(t, err)
	index := types.FindAccountInList(types.AccountKey(accountName), lst.AllValidators)
	if isInAllValidatorList {
		assert.True(t, index > -1)
	} else {
		assert.True(t, index == -1)
	}
}

// CreateAccount - register account on test blockchain
func CreateAccount(
	t *testing.T, accountName string, lb *app.LinoBlockchain, seq int64,
	resetPriv, transactionPriv, appPriv secp256k1.PrivKeySecp256k1,
	numOfLino string) {

	registerMsg := acc.NewRegisterMsg(
		GenesisUser, accountName, types.LNO(numOfLino),
		resetPriv.PubKey(), transactionPriv.PubKey(), appPriv.PubKey())
	SignCheckDeliver(t, lb, registerMsg, seq, true, GenesisTransactionPriv, time.Now().Unix())
}

// GetGenesisAccountCoin - get genesis account coin
func GetGenesisAccountCoin(numOfValidator int) types.Coin {
	coinPerValidator, _ := CoinPerValidator.ToInt64()
	genesisToken, _ := GenesisTotalCoin.ToInt64()
	initLNO := genesisToken - int64(numOfValidator)*coinPerValidator
	initCoin := types.NewCoinFromInt64(initLNO)
	return initCoin
}

// SignCheckDeliver - sign transaction, simulate and commit a block
func SignCheckDeliver(t *testing.T, lb *app.LinoBlockchain, msg sdk.Msg, seq int64,
	expPass bool, priv secp256k1.PrivKeySecp256k1, headTime int64) {
	// Sign the tx
	tx := genTx(msg, seq, priv)
	res := lb.Simulate(tx)
	if expPass {
		require.Equal(t, sdk.ABCICodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.ABCICodeOK, res.Code, res.Log)
	}

	// Simulate a Block
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			ChainID: "Lino", Time: time.Unix(headTime, 0)}})
	res = lb.Deliver(tx)
	if expPass {
		require.Equal(t, sdk.ABCICodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.ABCICodeOK, res.Code, res.Log)
	}
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
}

// SimulateOneBlock - simulate a empty block and commit
func SimulateOneBlock(lb *app.LinoBlockchain, headTime int64) {
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			ChainID: "Lino", Time: time.Unix(headTime, 0)}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
}

func genTx(msg sdk.Msg, seq int64, priv secp256k1.PrivKeySecp256k1) auth.StdTx {
	bz, _ := priv.Sign(auth.StdSignBytes("Lino", 0, seq, auth.StdFee{}, []sdk.Msg{msg}, ""))
	sigs := []auth.StdSignature{{
		PubKey:    priv.PubKey(),
		Signature: bz,
		Sequence:  seq}}
	return auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, sigs, "")
}

// CreateTestPost - create a test post
func CreateTestPost(
	t *testing.T, lb *app.LinoBlockchain,
	username, postID string, seq int64, priv secp256k1.PrivKeySecp256k1,
	sourceAuthor, sourcePostID string,
	parentAuthor, parentPostID string,
	redistributionSplitRate string, publishTime int64) {

	msg := post.CreatePostMsg{
		PostID:                  postID,
		Title:                   string(make([]byte, 50)),
		Content:                 string(make([]byte, 1000)),
		Author:                  types.AccountKey(username),
		ParentAuthor:            types.AccountKey(parentAuthor),
		ParentPostID:            parentPostID,
		SourceAuthor:            types.AccountKey(sourceAuthor),
		SourcePostID:            sourcePostID,
		Links:                   []types.IDToURLMapping{},
		RedistributionSplitRate: redistributionSplitRate,
	}
	SignCheckDeliver(t, lb, msg, seq, true, priv, publishTime)
}
