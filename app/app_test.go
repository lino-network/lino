package app

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

var (
	user1 = "validator0"
	priv1 = crypto.GenPrivKeyEd25519()
	addr1 = priv1.PubKey().Address()
	priv2 = crypto.GenPrivKeyEd25519()
	addr2 = priv2.PubKey().Address()

	genesisTotalLino    types.LNO  = "10000000000"
	genesisTotalCoin    types.Coin = types.NewCoin(10000000000 * types.Decimals)
	LNOPerValidator     types.LNO  = "100000000"
	growthRate          sdk.Rat    = sdk.NewRat(98, 1000)
	validatorAllocation sdk.Rat    = sdk.NewRat(10, 100)
)

func loggerAndDB() (log.Logger, dbm.DB) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return logger, db
}

func newLinoBlockchain(t *testing.T, numOfValidators int) *LinoBlockchain {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db)

	genesisState := genesis.GenesisState{
		Accounts:  []genesis.GenesisAccount{},
		TotalLino: genesisTotalLino,
	}

	// Generate 21 validators
	genesisAcc := genesis.GenesisAccount{
		Name:           user1,
		Lino:           LNOPerValidator,
		MasterKey:      priv1.PubKey(),
		TransactionKey: crypto.GenPrivKeyEd25519().PubKey(),
		PostKey:        crypto.GenPrivKeyEd25519().PubKey(),
		IsValidator:    true,
		ValPubKey:      priv2.PubKey(),
	}
	genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	for i := 1; i < numOfValidators; i++ {
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

func TestGenesisAcc(t *testing.T) {
	logger, db := loggerAndDB()
	lb := NewLinoBlockchain(logger, db)

	accs := []struct {
		genesisAccountName string
		numOfLino          types.LNO
		masterKey          crypto.PubKey
		transactionKey     crypto.PubKey
		postKey            crypto.PubKey
		isValidator        bool
		valPubKey          crypto.PubKey
	}{
		{"Lino", "9000000000", crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			true, crypto.GenPrivKeyEd25519().PubKey()},
		{"Genesis", "500000000", crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			true, crypto.GenPrivKeyEd25519().PubKey()},
		{"NonValidator", "500000000", crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
			false, crypto.GenPrivKeyEd25519().PubKey()},
	}
	genesisState := genesis.GenesisState{
		Accounts:  []genesis.GenesisAccount{},
		TotalLino: genesisTotalLino,
	}
	for _, acc := range accs {
		genesisAcc := genesis.GenesisAccount{
			Name:           acc.genesisAccountName,
			Lino:           acc.numOfLino,
			MasterKey:      acc.masterKey,
			TransactionKey: acc.transactionKey,
			PostKey:        acc.postKey,
			IsValidator:    acc.isValidator,
			ValPubKey:      acc.valPubKey,
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}

	result, err := genesis.GetGenesisJson(genesisState)
	assert.Nil(t, err)

	vals := []abci.Validator{}
	lb.InitChain(abci.RequestInitChain{vals, json.RawMessage(result)})
	lb.Commit()

	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	for _, acc := range accs {
		expectBalance, err := types.LinoToCoin(acc.numOfLino)
		assert.Nil(t, err)
		if acc.isValidator {
			expectBalance = expectBalance.Minus(
				types.ValidatorMinCommitingDeposit.Plus(types.ValidatorMinVotingDeposit))
		}
		balance, err :=
			lb.accountManager.GetBankBalance(ctx, types.AccountKey(acc.genesisAccountName))
		assert.Nil(t, err)
		assert.Equal(t, expectBalance, balance)
	}

	// reload app and ensure the account is still there
	lb = NewLinoBlockchain(logger, db)
	ctx = lb.BaseApp.NewContext(true, abci.Header{})
	for _, acc := range accs {
		expectBalance, err := types.LinoToCoin(acc.numOfLino)
		assert.Nil(t, err)
		if acc.isValidator {
			expectBalance = expectBalance.Minus(
				types.ValidatorMinCommitingDeposit.Plus(types.ValidatorMinVotingDeposit))
		}
		balance, err :=
			lb.accountManager.GetBankBalance(ctx, types.AccountKey(acc.genesisAccountName))
		assert.Nil(t, err)
		assert.Equal(t, expectBalance, balance)
	}
}

func TestDistributeInflationToValidators(t *testing.T) {
	lb := newLinoBlockchain(t, 21)

	baseTime := time.Now().Unix()
	remainValidatorPool := types.RatToCoin(
		genesisTotalCoin.ToRat().Mul(growthRate).Mul(validatorAllocation))
	coinPerValidator, _ := types.LinoToCoin(LNOPerValidator)
	expectBalance := coinPerValidator.Minus(
		types.ValidatorMinCommitingDeposit.Plus(types.ValidatorMinVotingDeposit))

	testPastMinutes := int64(0)
	for i := baseTime; i < baseTime+3600*20; i += 50 {
		lb.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{ChainID: "Lino", Time: baseTime + i}})
		lb.EndBlock(abci.RequestEndBlock{})
		lb.Commit()
		// simulate app
		if (baseTime+int64(i)-lb.chainStartTime)/int64(60) > testPastMinutes {
			testPastMinutes += 1
			if testPastMinutes%60 == 0 {
				// hourly inflation
				inflationForValidator :=
					types.RatToCoin(remainValidatorPool.ToRat().Mul(
						sdk.NewRat(1, types.HoursPerYear-lb.pastMinutes/60+1)))
				remainValidatorPool = remainValidatorPool.Minus(inflationForValidator)
				// expectBalance for all validators
				expectBalance = expectBalance.Plus(
					types.RatToCoin(inflationForValidator.ToRat().Quo(sdk.NewRat(21))))
				ctx := lb.BaseApp.NewContext(true, abci.Header{})
				for i := 0; i < 21; i++ {
					balance, err :=
						lb.accountManager.GetBankBalance(
							ctx, types.AccountKey("validator"+strconv.Itoa(i)))
					assert.Nil(t, err)
					assert.Equal(t, expectBalance, balance)
				}
			}
		}
	}
}

func TestFireByzantineValidators(t *testing.T) {
	lb := newLinoBlockchain(t, 21)

	lb.BeginBlock(abci.RequestBeginBlock{
		Header: abci.Header{
			ChainID: "Lino", Time: time.Now().Unix()},
		ByzantineValidators: []abci.Evidence{abci.Evidence{PubKey: priv2.PubKey().Bytes()}}})
	lb.EndBlock(abci.RequestEndBlock{})
	lb.Commit()
	ctx := lb.BaseApp.NewContext(true, abci.Header{})
	lst, err := lb.valManager.GetValidatorList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 20, len(lst.OncallValidators))
}
