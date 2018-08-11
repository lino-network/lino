package app

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/spf13/pflag"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	flagName       = "name"
	flagClientHome = "home-client"
	flagOWK        = "owk"
)

// get app init parameters for server init command
func LinoBlockchainInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)

	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagName, "", "validator moniker, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(flagOWK, false, "overwrite the accounts created")

	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         LinoBlockchainGenTx,
		AppGenState:      LinoBlockchainGenState,
	}
}

// genesis state for blockchain
type GenesisState struct {
	Accounts     []GenesisAccount       `json:"accounts"`
	Developers   []GenesisAppDeveloper  `json:"developers"`
	Infra        []GenesisInfraProvider `json:"infra"`
	GenesisParam GenesisParam           `json:"genesis_param"`
}

// genesis account will get coin to the address and register user
// if genesis account is validator, it will be added to validator list automatically
type GenesisAccount struct {
	Name           string        `json:"name"`
	Lino           types.LNO     `json:"lino"`
	ResetKey       crypto.PubKey `json:"reset_key"`
	TransactionKey crypto.PubKey `json:"transaction_key"`
	AppKey         crypto.PubKey `json:"app_key"`
	IsValidator    bool          `json:"is_validator"`
	ValPubKey      crypto.PubKey `json:"validator_pub_key"`
}

// register developer in genesis phase
type GenesisAppDeveloper struct {
	Name        string    `json:"name"`
	Deposit     types.LNO `json:"deposit"`
	Website     string    `json:"web_site"`
	Description string    `json:"description"`
	AppMetaData string    `json:"app_meta_data"`
}

// register infra provider in genesis phase
type GenesisInfraProvider struct {
	Name string `json:"name"`
}

type GenesisParam struct {
	InitFromConfig bool `json:"init_from_config"`
	param.EvaluateOfContentValueParam
	param.GlobalAllocationParam
	param.InfraInternalAllocationParam
	param.VoteParam
	param.ProposalParam
	param.DeveloperParam
	param.ValidatorParam
	param.CoinDayParam
	param.BandwidthParam
	param.AccountParam
	param.PostParam
}

func LinoBlockchainGenTx(cdc *wire.Codec, pk crypto.PubKey, genTxConfig config.GenTx) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	resetPriv := secp256k1.GenPrivKey()
	transactionPriv := secp256k1.GenPrivKey()
	appPriv := secp256k1.GenPrivKey()

	fmt.Println("reset private key is:", strings.ToUpper(hex.EncodeToString(resetPriv.Bytes())))
	fmt.Println("transaction private key is:", strings.ToUpper(hex.EncodeToString(transactionPriv.Bytes())))
	fmt.Println("app private key is:", strings.ToUpper(hex.EncodeToString(appPriv.Bytes())))

	totalLino := "10000000000"
	genesisAcc := GenesisAccount{
		Name:           "lino",
		Lino:           totalLino,
		ResetKey:       resetPriv.PubKey(),
		TransactionKey: transactionPriv.PubKey(),
		AppKey:         appPriv.PubKey(),
		IsValidator:    true,
		ValPubKey:      pk,
	}

	var bz []byte
	bz, err = wire.MarshalJSONIndent(cdc, genesisAcc)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  1000,
	}
	return
}

// default genesis file, only have one genesis account
func LinoBlockchainGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// totalLino := "10000000000"
	genesisState := GenesisState{
		Accounts:   []GenesisAccount{},
		Developers: []GenesisAppDeveloper{},
		Infra:      []GenesisInfraProvider{},
		GenesisParam: GenesisParam{
			true,
			param.EvaluateOfContentValueParam{
				ConsumptionTimeAdjustBase:      3153600,
				ConsumptionTimeAdjustOffset:    5,
				NumOfConsumptionOnAuthorOffset: 7,
				TotalAmountOfConsumptionBase:   1000 * types.Decimals,
				TotalAmountOfConsumptionOffset: 5,
				AmountOfConsumptionExponent:    sdk.NewRat(8, 10),
			},
			param.GlobalAllocationParam{
				InfraAllocation:          sdk.NewRat(20, 100),
				ContentCreatorAllocation: sdk.NewRat(65, 100),
				DeveloperAllocation:      sdk.NewRat(10, 100),
				ValidatorAllocation:      sdk.NewRat(5, 100),
			},
			param.InfraInternalAllocationParam{
				StorageAllocation: sdk.NewRat(50, 100),
				CDNAllocation:     sdk.NewRat(50, 100),
			},
			param.VoteParam{
				VoterMinDeposit:               types.NewCoinFromInt64(2000 * types.Decimals),
				VoterMinWithdraw:              types.NewCoinFromInt64(2 * types.Decimals),
				DelegatorMinWithdraw:          types.NewCoinFromInt64(2 * types.Decimals),
				VoterCoinReturnIntervalHr:     int64(7 * 24),
				VoterCoinReturnTimes:          int64(7),
				DelegatorCoinReturnIntervalHr: int64(7 * 24),
				DelegatorCoinReturnTimes:      int64(7),
			},
			param.ProposalParam{
				ContentCensorshipDecideHr:   int64(24 * 7),
				ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
				ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
				ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

				ChangeParamDecideHr:   int64(24 * 7),
				ChangeParamPassRatio:  sdk.NewRat(70, 100),
				ChangeParamPassVotes:  types.NewCoinFromInt64(1000000 * types.Decimals),
				ChangeParamMinDeposit: types.NewCoinFromInt64(100000 * types.Decimals),

				ProtocolUpgradeDecideHr:   int64(24 * 7),
				ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
				ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
				ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
			},
			param.DeveloperParam{
				DeveloperMinDeposit:           types.NewCoinFromInt64(1000000 * types.Decimals),
				DeveloperCoinReturnIntervalHr: int64(7 * 24),
				DeveloperCoinReturnTimes:      int64(7),
			},
			param.ValidatorParam{
				ValidatorMinWithdraw:          types.NewCoinFromInt64(1 * types.Decimals),
				ValidatorMinVotingDeposit:     types.NewCoinFromInt64(300000 * types.Decimals),
				ValidatorMinCommitingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
				ValidatorCoinReturnIntervalHr: int64(7 * 24),
				ValidatorCoinReturnTimes:      int64(7),
				PenaltyMissVote:               types.NewCoinFromInt64(20000 * types.Decimals),
				PenaltyMissCommit:             types.NewCoinFromInt64(200 * types.Decimals),
				PenaltyByzantine:              types.NewCoinFromInt64(1000000 * types.Decimals),
				ValidatorListSize:             int64(21),
				AbsentCommitLimitation:        int64(600), // 30min
			},
			param.CoinDayParam{
				DaysToRecoverCoinDayStake:    int64(7),
				SecondsToRecoverCoinDayStake: int64(7 * 24 * 3600),
			},
			param.BandwidthParam{
				SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
				CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
			},
			param.AccountParam{
				MinimumBalance:             types.NewCoinFromInt64(0),
				RegisterFee:                types.NewCoinFromInt64(1 * types.Decimals),
				FirstDepositFullStakeLimit: types.NewCoinFromInt64(1 * types.Decimals),
			},
			param.PostParam{
				ReportOrUpvoteInterval: 24 * 3600,
			},
		},
	}

	for _, genesisAccRaw := range appGenTxs {
		var genesisAcc GenesisAccount
		err = cdc.UnmarshalJSON(genesisAccRaw, &genesisAcc)
		if err != nil {
			return
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}
	genesisAppDeveloper := GenesisAppDeveloper{
		Name:        "lino",
		Deposit:     "1000000",
		Website:     "https://lino.network/",
		Description: "",
		AppMetaData: "",
	}
	genesisState.Developers = append(genesisState.Developers, genesisAppDeveloper)
	genesisInfraProvider := GenesisInfraProvider{
		Name: "lino",
	}
	genesisState.Infra = append(genesisState.Infra, genesisInfraProvider)

	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	return
}
