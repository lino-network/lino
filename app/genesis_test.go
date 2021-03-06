package app

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	wire "github.com/cosmos/cosmos-sdk/codec"
	// sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGetGenesisJson(t *testing.T) {
	resetPriv := secp256k1.GenPrivKey()
	transactionPriv := secp256k1.GenPrivKey()
	validatorPriv := secp256k1.GenPrivKey()

	totalLino := types.NewCoinFromInt64(10000000000 * types.Decimals)
	genesisAcc := GenesisAccount{
		Name:        "Lino",
		Coin:        totalLino,
		TxKey:       resetPriv.PubKey(),
		SignKey:     transactionPriv.PubKey(),
		IsValidator: true,
		ValPubKey:   validatorPriv.PubKey(),
	}

	genesisAppDeveloper := GenesisAppDeveloper{
		Name: "Lino",
	}
	genesisState := GenesisState{
		LoadPrevStates: false,
		GenesisPools: GenesisPools{
			Pools: []GenesisPool{
				{
					Name:   types.InflationDeveloperPool,
					Amount: types.NewCoinFromInt64(0),
				},
				{
					Name:   types.InflationValidatorPool,
					Amount: types.NewCoinFromInt64(0),
				},
				{
					Name:   types.InflationConsumptionPool,
					Amount: types.NewCoinFromInt64(0),
				},
				{
					Name:   types.VoteStakeInPool,
					Amount: types.NewCoinFromInt64(0),
				},
				{
					Name:   types.VoteStakeReturnPool,
					Amount: types.NewCoinFromInt64(0),
				},
				{
					Name:   types.VoteFrictionPool,
					Amount: types.NewCoinFromInt64(0),
				},
				{
					Name:   types.DevIDAReservePool,
					Amount: types.MustLinoToCoin("2000000000"),
				},
				{
					Name:   types.AccountVestingPool,
					Amount: types.MustLinoToCoin("8000000000"),
				},
			},
			Total: types.MustLinoToCoin("10000000000"),
		},
		InitCoinPrice: types.NewMiniDollar(1200),
		Accounts:      []GenesisAccount{genesisAcc},
		Developers:    []GenesisAppDeveloper{genesisAppDeveloper},
		GenesisParam: GenesisParam{
			true,
			param.GlobalAllocationParam{
				GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
				ContentCreatorAllocation: types.NewDecFromRat(10, 100),
				DeveloperAllocation:      types.NewDecFromRat(70, 100),
				ValidatorAllocation:      types.NewDecFromRat(20, 100),
			},
			param.VoteParam{
				MinStakeIn:                 types.NewCoinFromInt64(1000 * types.Decimals),
				VoterCoinReturnIntervalSec: int64(7 * 24 * 3600),
				VoterCoinReturnTimes:       int64(7),
			},
			param.ProposalParam{
				ContentCensorshipDecideSec:  int64(24 * 7 * 3600),
				ContentCensorshipPassRatio:  types.NewDecFromRat(50, 100),
				ContentCensorshipPassVotes:  types.NewCoinFromInt64(10000 * types.Decimals),
				ContentCensorshipMinDeposit: types.NewCoinFromInt64(100 * types.Decimals),

				ChangeParamDecideSec:  int64(24 * 7 * 3600),
				ChangeParamPassRatio:  types.NewDecFromRat(70, 100),
				ChangeParamPassVotes:  types.NewCoinFromInt64(1000000 * types.Decimals),
				ChangeParamMinDeposit: types.NewCoinFromInt64(100000 * types.Decimals),

				ProtocolUpgradeDecideSec:  int64(24 * 7 * 3600),
				ProtocolUpgradePassRatio:  types.NewDecFromRat(80, 100),
				ProtocolUpgradePassVotes:  types.NewCoinFromInt64(10000000 * types.Decimals),
				ProtocolUpgradeMinDeposit: types.NewCoinFromInt64(1000000 * types.Decimals),
			},
			param.DeveloperParam{
				DeveloperMinDeposit:            types.NewCoinFromInt64(1000000 * types.Decimals),
				DeveloperCoinReturnIntervalSec: int64(7 * 24 * 3600),
				DeveloperCoinReturnTimes:       int64(7),
			},
			param.ValidatorParam{
				ValidatorMinDeposit:            types.NewCoinFromInt64(200000 * types.Decimals),
				ValidatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
				ValidatorCoinReturnTimes:       int64(7),
				PenaltyMissCommit:              types.NewCoinFromInt64(200 * types.Decimals),
				PenaltyByzantine:               types.NewCoinFromInt64(1000 * types.Decimals),
				AbsentCommitLimitation:         int64(600), // 30min
				OncallSize:                     int64(22),
				StandbySize:                    int64(7),
				ValidatorRevokePendingSec:      int64(7 * 24 * 3600),
				OncallInflationWeight:          int64(2),
				StandbyInflationWeight:         int64(1),
				MaxVotedValidators:             int64(3),
				SlashLimitation:                int64(5),
			},
			param.BandwidthParam{
				SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
				CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
				VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
				GeneralMsgQuotaRatio:        types.NewDecFromRat(20, 100),
				GeneralMsgEMAFactor:         types.NewDecFromRat(1, 10),
				AppMsgQuotaRatio:            types.NewDecFromRat(80, 100),
				AppMsgEMAFactor:             types.NewDecFromRat(1, 10),
				ExpectedMaxMPS:              types.NewDecFromRat(1000, 1),
				MsgFeeFactorA:               types.NewDecFromRat(6, 1),
				MsgFeeFactorB:               types.NewDecFromRat(10, 1),
				MaxMPSDecayRate:             types.NewDecFromRat(99, 100),
				AppBandwidthPoolSize:        types.NewDecFromRat(10, 1),
				AppVacancyFactor:            types.NewDecFromRat(69, 100),
				AppPunishmentFactor:         types.NewDecFromRat(14, 5),
			},
			param.AccountParam{
				MinimumBalance: types.NewCoinFromInt64(0),
				RegisterFee:    types.NewCoinFromInt64(1 * types.Decimals),
			},
			param.PostParam{},
			param.ReputationParam{
				BestContentIndexN: 200,
				UserMaxN:          50,
			},
			param.PriceParam{
				TestnetMode:     true,
				UpdateEverySec:  int64(time.Hour.Seconds()),
				FeedEverySec:    int64((10 * time.Minute).Seconds()),
				HistoryMaxLen:   71,
				PenaltyMissFeed: types.NewCoinFromInt64(10000 * types.Decimals),
			},
		},
	}

	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	appState, err := wire.MarshalJSONIndent(cdc, genesisState)
	assert.Nil(t, err)
	appGenesisState := new(GenesisState)
	err = cdc.UnmarshalJSON(appState, appGenesisState)
	assert.Nil(t, err)

	assert.Equal(t, genesisState, *appGenesisState)
}

func TestLinoBlockchainGenTx(t *testing.T) {
	cdc := MakeCodec()
	pk := secp256k1.GenPrivKey().PubKey()
	appGenTx, _, validator, err := LinoBlockchainGenTx(cdc, pk)
	assert.Nil(t, err)
	var genesisAcc GenesisAccount
	err = cdc.UnmarshalJSON(appGenTx, &genesisAcc)
	assert.Nil(t, err)
	assert.Equal(t, genesisAcc.Name, "lino")
	assert.Equal(t, genesisAcc.Coin, types.NewCoinFromInt64(100000000*types.Decimals))
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
			Name:        "validator" + strconv.Itoa(i),
			Coin:        coinPerValidator,
			TxKey:       secp256k1.GenPrivKey().PubKey(),
			SignKey:     secp256k1.GenPrivKey().PubKey(),
			IsValidator: true,
			ValPubKey:   secp256k1.GenPrivKey().PubKey(),
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
}
