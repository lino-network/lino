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
	accmn "github.com/lino-network/lino/x/account/manager"
	acctypes "github.com/lino-network/lino/x/account/types"
	bandwidthmn "github.com/lino-network/lino/x/bandwidth/manager"
	bandwidthmodel "github.com/lino-network/lino/x/bandwidth/model"
	"github.com/lino-network/lino/x/global"
	globalModel "github.com/lino-network/lino/x/global/model"
	post "github.com/lino-network/lino/x/post"
	val "github.com/lino-network/lino/x/validator"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
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
	ConsumptionFrictionRate            = types.NewDecFromRat(5, 100)
	ConsumptionFreezingPeriodSec int64 = 24 * 7 * 3600
	PostIntervalSec              int64 = 600
)

func loggerAndDB() (log.Logger, dbm.DB) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return logger, db
}

func NewTestLinoBlockchain(t *testing.T, numOfValidators int, beginBlockTime time.Time) *app.LinoBlockchain {
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
		MaxTPS:                       sdk.NewDec(1000),
		ConsumptionFreezingPeriodSec: 7 * 24 * 3600,
		ConsumptionFrictionRate:      types.NewDecFromRat(5, 100),
	}
	result, err := wire.MarshalJSONIndent(cdc, genesisState)
	assert.Nil(t, err)

	lb.InitChain(abci.RequestInitChain{ChainId: "Lino", AppStateBytes: json.RawMessage(result)})
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{Height: 1, ChainID: "Lino", Time: beginBlockTime}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	bandwidthmn.BandwidthManagerTestMode = true
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
	gm := global.NewGlobalManager(lb.CapKeyGlobalStore, ph)
	accManager := accmn.NewAccountManager(lb.CapKeyAccountStore, ph, &gm)
	saving, err := accManager.GetSavingFromUsername(ctx, types.AccountKey(accountName))
	assert.Nil(t, err)
	assert.Equal(t, expectBalance.Amount.Int64(), saving.Amount.Int64())
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

// CheckAppBandwidthInfo
func CheckAppBandwidthInfo(
	t *testing.T, info bandwidthmodel.AppBandwidthInfo, username types.AccountKey, lb *app.LinoBlockchain) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	bs := bandwidthmodel.NewBandwidthStorage(lb.CapKeyBandwidthStore)
	res, err := bs.GetAppBandwidthInfo(ctx, username)
	assert.Nil(t, err)
	assert.Equal(t, info, *res)
}

// CheckCurBlockInfo
func CheckCurBlockInfo(
	t *testing.T, info bandwidthmodel.BlockInfo, lb *app.LinoBlockchain) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	bs := bandwidthmodel.NewBandwidthStorage(lb.CapKeyBandwidthStore)
	res, err := bs.GetBlockInfo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, info, *res)
}

// CreateAccount - register account on test blockchain
func CreateAccount(
	t *testing.T, accountName string, lb *app.LinoBlockchain, seq uint64,
	resetPriv, transactionPriv, appPriv secp256k1.PrivKeySecp256k1,
	numOfLino string) {

	registerMsg := acctypes.NewRegisterMsg(
		GenesisUser, accountName, numOfLino,
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
func SignCheckDeliver(t *testing.T, lb *app.LinoBlockchain, msg sdk.Msg, seq uint64,
	expPass bool, priv secp256k1.PrivKeySecp256k1, headTime int64) {
	// Sign the tx
	tx := genTx(msg, seq, priv)
	// XXX(yumin): API changed after upgrad-1, new field tx, passing nil, not sure
	// about what is the right way..
	res := lb.Simulate(nil, tx)
	if expPass {
		require.True(t, res.IsOK(), res.Log)
	} else {
		require.False(t, res.IsOK(), res.Log)
	}

	// Simulate a Block
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			Height: lb.LastBlockHeight() + 1, ChainID: "Lino", Time: time.Unix(headTime, 0)}})
	res = lb.Deliver(tx)
	if expPass {
		require.True(t, res.IsOK(), res.Log)
	} else {
		require.False(t, res.IsOK(), res.Log)
	}
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
}

// SignCheckDeliverWithFee - sign transaction with fee, simulate and commit a block
func SignCheckDeliverWithFee(t *testing.T, lb *app.LinoBlockchain, msg sdk.Msg, seq uint64,
	expPass bool, priv secp256k1.PrivKeySecp256k1, headTime int64, fee auth.StdFee) {
	// Sign the tx
	tx := genTxWithFee(msg, seq, priv, fee)
	// XXX(yumin): API changed after upgrad-1, new field tx, passing nil, not sure
	// about what is the right way..
	res := lb.Simulate(nil, tx)
	if expPass {
		require.True(t, res.IsOK(), res.Log)
	} else {
		require.False(t, res.IsOK(), res.Log)
	}

	// Simulate a Block
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			Height: lb.LastBlockHeight() + 1, ChainID: "Lino", Time: time.Unix(headTime, 0)}})
	res = lb.Deliver(tx)
	if expPass {
		require.True(t, res.IsOK(), res.Log)
	} else {
		require.False(t, res.IsOK(), res.Log)
	}
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
}

// RepeatSignCheckDeliver - sign same transaction repeatly, simulate and commit a block
func RepeatSignCheckDeliver(t *testing.T, lb *app.LinoBlockchain, msg sdk.Msg, seq uint64,
	expPass bool, priv secp256k1.PrivKeySecp256k1, headTime int64, times int) {

	// Simulate a Block
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			Height: lb.LastBlockHeight() + 1, ChainID: "Lino", Time: time.Unix(headTime, 0)}})

	for i := 0; i < times; i++ {
		tx := genTx(msg, seq+uint64(i), priv)
		res := lb.Deliver(tx)
		if expPass {
			require.True(t, res.IsOK(), res.Log)
		} else {
			require.False(t, res.IsOK(), res.Log)
		}
	}
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
}

// SimulateOneBlock - simulate a empty block and commit
func SimulateOneBlock(lb *app.LinoBlockchain, headTime int64) {
	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			Height: lb.LastBlockHeight() + 1, ChainID: "Lino", Time: time.Unix(headTime, 0)}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
}

func genTx(msg sdk.Msg, seq uint64, priv secp256k1.PrivKeySecp256k1) auth.StdTx {
	bz, _ := priv.Sign(auth.StdSignBytes("Lino", 0, seq, auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(types.LinoCoinDenom, sdk.NewInt(10000000)))}, []sdk.Msg{msg}, ""))
	sigs := []auth.StdSignature{{
		PubKey:    priv.PubKey(),
		Signature: bz,
	}}
	return auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(types.LinoCoinDenom, sdk.NewInt(10000000)))}, sigs, "")
}

func genTxWithFee(msg sdk.Msg, seq uint64, priv secp256k1.PrivKeySecp256k1, fee auth.StdFee) auth.StdTx {
	bz, _ := priv.Sign(auth.StdSignBytes("Lino", 0, seq, fee, []sdk.Msg{msg}, ""))
	sigs := []auth.StdSignature{{
		PubKey:    priv.PubKey(),
		Signature: bz,
	}}
	return auth.NewStdTx([]sdk.Msg{msg}, fee, sigs, "")
}

// CreateTestPost - create a test post
func CreateTestPost(
	t *testing.T, lb *app.LinoBlockchain,
	username, postID string, seq uint64, priv secp256k1.PrivKeySecp256k1, publishTime int64) {

	msg := post.CreatePostMsg{
		PostID:    postID,
		Title:     string(make([]byte, 50)),
		Content:   string(make([]byte, 1000)),
		Author:    types.AccountKey(username),
		CreatedBy: types.AccountKey(username),
	}
	SignCheckDeliver(t, lb, msg, seq, true, priv, publishTime)
}
