package app

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	globalModel "github.com/lino-network/lino/x/global/model"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

func TestGetGenesisJson(t *testing.T) {
	resetPriv := secp256k1.GenPrivKey()
	transactionPriv := secp256k1.GenPrivKey()
	appPriv := secp256k1.GenPrivKey()
	validatorPriv := secp256k1.GenPrivKey()

	totalLino := types.NewCoinFromInt64(10000000000 * types.Decimals)
	genesisAcc := GenesisAccount{
		Name:           "Lino",
		Coin:           totalLino,
		ResetKey:       resetPriv.PubKey(),
		TransactionKey: transactionPriv.PubKey(),
		AppKey:         appPriv.PubKey(),
		IsValidator:    true,
		ValPubKey:      validatorPriv.PubKey(),
	}

	genesisAppDeveloper := GenesisAppDeveloper{
		Name:    "Lino",
		Deposit: types.NewCoinFromInt64(1000000 * types.Decimals),
	}
	genesisInfraProvider := GenesisInfraProvider{
		Name: "Lino",
	}
	genesisState := GenesisState{
		Accounts:   []GenesisAccount{genesisAcc},
		Developers: []GenesisAppDeveloper{genesisAppDeveloper},
		Infra:      []GenesisInfraProvider{genesisInfraProvider},
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
				Ceiling:                  sdk.NewRat(98, 1000),
				Floor:                    sdk.NewRat(3, 100),
				GlobalGrowthRate:         sdk.NewRat(98, 1000),
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
				VoterMinDeposit:                types.NewCoinFromInt64(2000 * types.Decimals),
				VoterMinWithdraw:               types.NewCoinFromInt64(2 * types.Decimals),
				DelegatorMinWithdraw:           types.NewCoinFromInt64(2 * types.Decimals),
				VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
				VoterCoinReturnTimes:           int64(7),
				DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
				DelegatorCoinReturnTimes:       int64(7),
			},
			param.ProposalParam{
				ContentCensorshipDecideSec:  int64(24 * 7 * 3600),
				ContentCensorshipPassRatio:  sdk.NewRat(50, 100),
				ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
				ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

				ChangeParamDecideSec:  int64(24 * 7 * 3600),
				ChangeParamPassRatio:  sdk.NewRat(70, 100),
				ChangeParamPassVotes:  types.NewCoinFromInt64(1000000 * types.Decimals),
				ChangeParamMinDeposit: types.NewCoinFromInt64(100000 * types.Decimals),

				ProtocolUpgradeDecideSec:  int64(24 * 7 * 3600),
				ProtocolUpgradePassRatio:  sdk.NewRat(80, 100),
				ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
				ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
			},
			param.DeveloperParam{
				DeveloperMinDeposit:            types.NewCoinFromInt64(1000000 * types.Decimals),
				DeveloperCoinReturnIntervalSec: int64(7 * 24 * 3600),
				DeveloperCoinReturnTimes:       int64(7),
			},
			param.ValidatorParam{
				ValidatorMinWithdraw:           types.NewCoinFromInt64(1 * types.Decimals),
				ValidatorMinVotingDeposit:      types.NewCoinFromInt64(300000 * types.Decimals),
				ValidatorMinCommittingDeposit:  types.NewCoinFromInt64(100000 * types.Decimals),
				ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
				ValidatorCoinReturnTimes:       int64(7),
				PenaltyMissVote:                types.NewCoinFromInt64(20000 * types.Decimals),
				PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
				PenaltyByzantine:               types.NewCoinFromInt64(1000000 * types.Decimals),
				ValidatorListSize:              int64(21),
				AbsentCommitLimitation:         int64(600), // 30min
			},
			param.CoinDayParam{
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
				ReportOrUpvoteIntervalSec: 24 * 3600,
				PostIntervalSec:           600,
			},
		},
		InitGlobalMeta: globalModel.InitParamList{
			MaxTPS: sdk.NewRat(1000),
			ConsumptionFreezingPeriodSec: 7 * 24 * 3600,
			ConsumptionFrictionRate:      sdk.NewRat(5, 100),
		},
	}

	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	appState, err := wire.MarshalJSONIndent(cdc, genesisState)
	assert.Nil(t, err)
	appGenesisState := new(GenesisState)
	err = cdc.UnmarshalJSON([]byte(appState), appGenesisState)
	assert.Nil(t, err)

	assert.Equal(t, genesisState, *appGenesisState)
}

func TestLinoBlockchainGenTx(t *testing.T) {
	cdc := MakeCodec()
	pk := secp256k1.GenPrivKey().PubKey()
	var genTxConfig config.GenTx
	appGenTx, _, validator, err := LinoBlockchainGenTx(cdc, pk, genTxConfig)
	assert.Nil(t, err)
	var genesisAcc GenesisAccount
	err = cdc.UnmarshalJSON(appGenTx, &genesisAcc)
	assert.Nil(t, err)
	assert.Equal(t, genesisAcc.Name, "lino")
	assert.Equal(t, genesisAcc.Coin, types.NewCoinFromInt64(10000000000*types.Decimals))
	assert.Equal(t, genesisAcc.IsValidator, true)
	assert.Equal(t, genesisAcc.ValPubKey, pk)
	assert.Equal(t, validator.PubKey, pk)
}

func TestLinoBlockchainGenState(t *testing.T) {
	cdc := MakeCodec()
	appGenTxs := []json.RawMessage{}
	coinPerValidator := types.NewCoinFromInt64(100000000 * types.Decimals)
	for i := 1; i < 21; i++ {
		genesisAcc := GenesisAccount{
			Name:           "validator" + strconv.Itoa(i),
			Coin:           coinPerValidator,
			ResetKey:       secp256k1.GenPrivKey().PubKey(),
			TransactionKey: secp256k1.GenPrivKey().PubKey(),
			AppKey:         secp256k1.GenPrivKey().PubKey(),
			IsValidator:    true,
			ValPubKey:      secp256k1.GenPrivKey().PubKey(),
		}
		marshalResult, err := wire.MarshalJSONIndent(cdc, genesisAcc)
		assert.Nil(t, err)
		appGenTxs = append(appGenTxs, json.RawMessage(marshalResult))
	}
	appState, err := LinoBlockchainGenState(cdc, appGenTxs)
	assert.Nil(t, err)

	genesisState := new(GenesisState)
	if err := cdc.UnmarshalJSON(appState, genesisState); err != nil {
		panic(err)
	}
	for i, gacc := range genesisState.Accounts {
		assert.Equal(t, gacc.Name, "validator"+strconv.Itoa(i+1))
		assert.Equal(t, gacc.Coin, coinPerValidator)
	}
	assert.Equal(t, 1, len(genesisState.Developers))
	assert.Equal(t, 1, len(genesisState.Infra))
}
