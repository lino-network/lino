package test

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/genesis"
	acc "github.com/lino-network/lino/tx/account"
	post "github.com/lino-network/lino/tx/post"
	reg "github.com/lino-network/lino/tx/register"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

var (
	GenesisUser = "genesis"
	GenesisPriv = crypto.GenPrivKeyEd25519()
	GenesisAddr = GenesisPriv.PubKey().Address()

	DefaultNumOfVal  int        = 21
	GenesisTotalLino int64      = 10000000000
	LNOPerValidator  int64      = 100000000
	GenesisTotalCoin types.Coin = types.NewCoin(GenesisTotalLino * types.Decimals)

	ConsumptionFrictionRate sdk.Rat = sdk.NewRat(5, 100)
	FreezingPeriodHr        int64   = 24 * 7
)

func loggerAndDBs() (log.Logger, map[string]dbm.DB) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	dbs := map[string]dbm.DB{
		"acc":       dbm.NewMemDB(),
		"post":      dbm.NewMemDB(),
		"val":       dbm.NewMemDB(),
		"vote":      dbm.NewMemDB(),
		"infra":     dbm.NewMemDB(),
		"developer": dbm.NewMemDB(),
		"global":    dbm.NewMemDB(),
	}
	return logger, dbs
}

func NewTestLinoBlockchain(t *testing.T, numOfValidators int) *app.LinoBlockchain {
	logger, dbs := loggerAndDBs()
	lb := app.NewLinoBlockchain(logger, dbs)
	globalState := genesis.GlobalState{
		TotalLino:                GenesisTotalLino,
		GrowthRate:               sdk.NewRat(98, 1000),
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(50, 100),
		DeveloperAllocation:      sdk.NewRat(20, 100),
		ValidatorAllocation:      sdk.NewRat(10, 100),
		ConsumptionFrictionRate:  ConsumptionFrictionRate,
		FreezingPeriodHr:         FreezingPeriodHr,
	}

	genesisState := genesis.GenesisState{
		Accounts:    []genesis.GenesisAccount{},
		GlobalState: globalState,
	}

	// Generate 21 validators
	for i := 0; i < numOfValidators; i++ {
		privKey := crypto.GenPrivKeyEd25519()
		valPrivKey := crypto.GenPrivKeyEd25519()
		genesisAcc := genesis.GenesisAccount{
			Name:        "validator" + strconv.Itoa(i),
			Lino:        LNOPerValidator,
			PubKey:      privKey.PubKey(),
			IsValidator: true,
			ValPubKey:   valPrivKey.PubKey(),
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}

	genesisAcc := genesis.GenesisAccount{
		Name:        GenesisUser,
		Lino:        GenesisTotalLino - int64(numOfValidators)*LNOPerValidator,
		PubKey:      GenesisPriv.PubKey(),
		IsValidator: false,
		ValPubKey:   GenesisPriv.PubKey(),
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

func CheckBalance(t *testing.T, accountName string, lb *app.LinoBlockchain, expectBalance types.Coin) {
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	accManager := acc.NewAccountManager(lb.CapKeyAccountStore)
	balance, err :=
		accManager.GetBankBalance(ctx, types.AccountKey(accountName))
	assert.Nil(t, err)
	assert.Equal(t, expectBalance, balance)
}

func CreateAccount(
	t *testing.T, accountName string, lb *app.LinoBlockchain, seq int64,
	priv crypto.PrivKeyEd25519, numOfLino int64) {

	transferMsg := acc.NewTransferMsg(
		GenesisUser, types.LNO(sdk.NewRat(numOfLino)),
		[]byte{}, acc.TransferToAddr(priv.PubKey().Address()))

	SignCheckDeliver(t, lb, transferMsg, seq, true, GenesisPriv, time.Now().Unix())

	registerMsg := reg.NewRegisterMsg(accountName, priv.PubKey())
	SignCheckDeliver(t, lb, registerMsg, 0, true, priv, time.Now().Unix())
}

func GetGenesisAccountCoin(numOfValidator int) types.Coin {
	return types.NewCoin((GenesisTotalLino - LNOPerValidator*21) * types.Decimals)
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
	redistributionSplitRate sdk.Rat, publishTime int64) {

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
