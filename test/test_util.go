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

	DefaultNumOfVal  int       = 21
	GenesisTotalLino types.LNO = "10000000000"
	LNOPerValidator  types.LNO = "100000000"

	PenaltyMissVote       types.Coin = types.NewCoinFromInt64(20000 * types.Decimals)
	ChangeParamMinDeposit types.Coin = types.NewCoinFromInt64(100000 * types.Decimals)

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
	lb := app.NewLinoBlockchain(logger, db, nil)

	genesisState := app.GenesisState{
		Accounts: []app.GenesisAccount{},
	}

	// Generate 21 validators
	for i := 0; i < numOfValidators; i++ {
		genesisAcc := app.GenesisAccount{
			Name:           "validator" + strconv.Itoa(i),
			Lino:           LNOPerValidator,
			ResetKey:       secp256k1.GenPrivKey().PubKey(),
			TransactionKey: secp256k1.GenPrivKey().PubKey(),
			AppKey:         secp256k1.GenPrivKey().PubKey(),
			IsValidator:    true,
			ValPubKey:      secp256k1.GenPrivKey().PubKey(),
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}

	totalAmt, _ := strconv.ParseInt(GenesisTotalLino, 10, 64)
	validatorAmt, _ := strconv.ParseInt(LNOPerValidator, 10, 64)
	initLNO := strconv.FormatInt(totalAmt-int64(numOfValidators)*validatorAmt, 10)
	genesisAcc := app.GenesisAccount{
		Name:           GenesisUser,
		Lino:           initLNO,
		ResetKey:       GenesisPriv.PubKey(),
		TransactionKey: GenesisTransactionPriv.PubKey(),
		AppKey:         GenesisAppPriv.PubKey(),
		IsValidator:    false,
		ValPubKey:      GenesisPriv.PubKey(),
	}
	cdc := app.MakeCodec()
	genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	genesisState.InitGlobalMeta = globalModel.InitParamList{
		GrowthRate: sdk.NewRat(98, 1000),
		Ceiling:    sdk.NewRat(98, 1000),
		Floor:      sdk.NewRat(3, 100),
		MaxTPS:     sdk.NewRat(1000),
		ConsumptionFreezingPeriodHr: 7 * 24,
		ConsumptionFrictionRate:     sdk.NewRat(5, 100),
	}
	result, err := wire.MarshalJSONIndent(cdc, genesisState)
	assert.Nil(t, err)

	lb.InitChain(abci.RequestInitChain{ChainId: "Lino", AppStateBytes: json.RawMessage(result)})
	lb.Commit()

	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	return lb
}

func CheckGlobalAllocation(t *testing.T, lb *app.LinoBlockchain, expectAllocation param.GlobalAllocationParam) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	allocation, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectAllocation, *allocation)
}

func CheckBalance(t *testing.T, accountName string, lb *app.LinoBlockchain, expectBalance types.Coin) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	accManager := acc.NewAccountManager(lb.CapKeyAccountStore, ph)
	saving, err :=
		accManager.GetSavingFromBank(ctx, types.AccountKey(accountName))
	assert.Nil(t, err)
	assert.Equal(t, expectBalance, saving)
}

func CheckValidatorDeposit(t *testing.T, accountName string, lb *app.LinoBlockchain, expectDeposit types.Coin) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
	ph := param.NewParamHolder(lb.CapKeyParamStore)
	valManager := val.NewValidatorManager(lb.CapKeyValStore, ph)
	deposit, err := valManager.GetValidatorDeposit(ctx, types.AccountKey(accountName))
	assert.Nil(t, err)
	assert.Equal(t, expectDeposit, deposit)
}

func CheckOncallValidatorList(
	t *testing.T, accountName string, isInOnCallValidatorList bool, lb *app.LinoBlockchain) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
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
	ctx := lb.BaseApp.NewContext(true, abci.Header{ChainID: "Lino", Time: time.Unix(0, 0)})
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
	resetPriv, transactionPriv, appPriv secp256k1.PrivKeySecp256k1,
	numOfLino string) {

	registerMsg := acc.NewRegisterMsg(
		GenesisUser, accountName, types.LNO(numOfLino),
		resetPriv.PubKey(), transactionPriv.PubKey(), appPriv.PubKey())
	SignCheckDeliver(t, lb, registerMsg, seq, true, GenesisTransactionPriv, time.Now().Unix())
}

func GetGenesisAccountCoin(numOfValidator int) types.Coin {
	totalAmt, _ := strconv.ParseInt(GenesisTotalLino, 10, 64)
	validatorAmt, _ := strconv.ParseInt(LNOPerValidator, 10, 64)
	initLNO := strconv.FormatInt(totalAmt-int64(numOfValidator)*validatorAmt, 10)
	initCoin, _ := types.LinoToCoin(initLNO)
	return initCoin
}

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

func CreateTestPost(
	t *testing.T, lb *app.LinoBlockchain,
	username, postID string, seq int64, priv secp256k1.PrivKeySecp256k1,
	sourceAuthor, sourcePostID string,
	parentAuthor, parentPostID string,
	redistributionSplitRate string, publishTime int64) {

	msg := post.CreatePostMsg{
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
	SignCheckDeliver(t, lb, msg, seq, true, priv, publishTime)
}

func CoinToString(coin types.Coin) string {
	return strconv.FormatInt(coin.ToInt64()/types.Decimals, 10)
}
