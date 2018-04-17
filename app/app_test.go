package app

import (
	"encoding/json"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

var (
	priv1 = crypto.GenPrivKeyEd25519()
	addr1 = priv1.PubKey().Address()
	priv2 = crypto.GenPrivKeyEd25519()
	addr2 = priv2.PubKey().Address()
	priv3 = crypto.GenPrivKeyEd25519()
	addr3 = priv3.PubKey().Address()
	priv4 = crypto.GenPrivKeyEd25519()
	addr4 = priv4.PubKey().Address()
	priv5 = crypto.GenPrivKeyEd25519()
	addr5 = priv3.PubKey().Address()
	priv6 = crypto.GenPrivKeyEd25519()
	addr6 = priv4.PubKey().Address()
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

func newLinoBlockchain() *LinoBlockchain {
	logger, dbs := loggerAndDBs()
	return NewLinoBlockchain(logger, dbs)
}

func TestGenesisAcc(t *testing.T) {
	logger, dbs := loggerAndDBs()
	lb := NewLinoBlockchain(logger, dbs)

	accs := []struct {
		genesisAccountName string
		numOfLino          int64
		pubKey             crypto.PubKey
		isValidator        bool
		valPubKey          crypto.PubKey
	}{
		{"Lino", 9000000000, priv1.PubKey(), true, priv2.PubKey()},
		{"Genesis", 500000000, priv3.PubKey(), true, priv4.PubKey()},
		{"NonValidator", 500000000, priv5.PubKey(), false, priv6.PubKey()},
	}

	totalLino := int64(10000000000)

	globalState := genesis.GlobalState{
		TotalLino:                totalLino,
		GrowthRate:               sdk.NewRat(98, 1000),
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(50, 100),
		DeveloperAllocation:      sdk.NewRat(20, 100),
		ValidatorAllocation:      sdk.NewRat(10, 100),
		ConsumptionFrictionRate:  sdk.NewRat(1, 100),
		FreezingPeriodHr:         24 * 7,
	}

	genesisState := genesis.GenesisState{
		Accounts:    []genesis.GenesisAccount{},
		GlobalState: globalState,
	}
	for _, acc := range accs {
		genesisAcc := genesis.GenesisAccount{
			Name:        acc.genesisAccountName,
			Lino:        acc.numOfLino,
			PubKey:      acc.pubKey,
			IsValidator: acc.isValidator,
			ValPubKey:   acc.valPubKey,
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
		expectBalance, err := types.LinoToCoin(types.LNO(sdk.NewRat(acc.numOfLino)))
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
	lb = NewLinoBlockchain(logger, dbs)
	ctx = lb.BaseApp.NewContext(true, abci.Header{})
	for _, acc := range accs {
		expectBalance, err := types.LinoToCoin(types.LNO(sdk.NewRat(acc.numOfLino)))
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
