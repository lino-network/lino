package app

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
)

type GenesisPool struct {
	Name   types.PoolName `json:"name"`
	Amount types.Coin     `json:"amount"`
}

type GenesisPools struct {
	Pools []GenesisPool `json:"pools"`
	Total types.Coin    `json:"total"`
}

func (g GenesisPools) IsValid() error {
	poolMap := make(map[types.PoolName]bool)
	poolMap[types.InflationDeveloperPool] = true
	poolMap[types.InflationValidatorPool] = true
	poolMap[types.InflationConsumptionPool] = true
	poolMap[types.AccountVestingPool] = true
	poolMap[types.VoteStakeInPool] = true
	poolMap[types.VoteStakeReturnPool] = true
	poolMap[types.VoteFrictionPool] = true
	poolMap[types.DevIDAReservePool] = true

	// checks
	if len(g.Pools) != len(poolMap) {
		return fmt.Errorf("expected number of pools: %d, actual: %d", len(poolMap), len(g.Pools))
	}

	total := types.NewCoinFromInt64(0)
	for _, pool := range g.Pools {
		if !poolMap[pool.Name] {
			return fmt.Errorf("unknown pool: %s", pool.Name)
		}
		total = total.Plus(pool.Amount)
	}
	if !total.IsEqual(g.Total) {
		return fmt.Errorf("expected total: %s, actual: %s", g.Total, total)
	}
	return nil
}

func (g GenesisPools) ReservePool() types.Coin {
	for _, pool := range g.Pools {
		if pool.Name == types.DevIDAReservePool {
			return pool.Amount
		}
	}
	panic("reserve pool not found in genesis file")
}

// genesis state for blockchain
type GenesisState struct {
	LoadPrevStates bool                  `json:"load_prev_states"`
	GenesisPools   GenesisPools          `json:"genesis_pools"`
	InitCoinPrice  types.MiniDollar      `json:"init_coin_price"`
	Accounts       []GenesisAccount      `json:"accounts"`
	Developers     []GenesisAppDeveloper `json:"developers"`
	GenesisParam   GenesisParam          `json:"genesis_param"`
}

// genesis account will get coin to the address and register user
// if genesis account is validator, it will be added to validator list automatically
type GenesisAccount struct {
	Name           string        `json:"name"`
	Coin           types.Coin    `json:"coin"`
	ResetKey       crypto.PubKey `json:"reset_key"`
	TransactionKey crypto.PubKey `json:"transaction_key"`
	IsValidator    bool          `json:"is_validator"`
	ValPubKey      crypto.PubKey `json:"validator_pub_key"`
}

// GenesisAppDeveloper - register developer in genesis phase
type GenesisAppDeveloper struct {
	Name        string `json:"name"`
	Website     string `json:"web_site"`
	Description string `json:"description"`
	AppMetaData string `json:"app_meta_data"`
}

// GenesisParam - genesis parameters
type GenesisParam struct {
	InitFromConfig bool `json:"init_from_config"`
	param.GlobalAllocationParam
	param.VoteParam
	param.ProposalParam
	param.DeveloperParam
	param.ValidatorParam
	param.CoinDayParam
	param.BandwidthParam
	param.AccountParam
	param.PostParam
	param.ReputationParam
	param.PriceParam
}

// LinoBlockchainGenTx - init genesis account
func LinoBlockchainGenTx(cdc *wire.Codec, pk crypto.PubKey) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	resetPriv := secp256k1.GenPrivKey()
	transactionPriv := secp256k1.GenPrivKey()

	fmt.Println("reset private key is:", strings.ToUpper(hex.EncodeToString(resetPriv.Bytes())))
	fmt.Println("transaction private key is:", strings.ToUpper(hex.EncodeToString(transactionPriv.Bytes())))

	totalCoin := types.NewCoinFromInt64(10000000000 * types.Decimals)
	genesisAcc := GenesisAccount{
		Name:           "lino",
		Coin:           totalCoin,
		ResetKey:       resetPriv.PubKey(),
		TransactionKey: transactionPriv.PubKey(),
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

// LinoBlockchainGenState - default genesis file
func LinoBlockchainGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// totalLino := "10000000000"
	genesisState := GenesisState{
		LoadPrevStates: false,
		GenesisPools: GenesisPools{
			Pools: []GenesisPool{
				{Name: types.InflationDeveloperPool},
				{Name: types.InflationValidatorPool},
				{Name: types.InflationConsumptionPool},
				{Name: types.VoteStakeInPool},
				{Name: types.VoteStakeReturnPool},
				{Name: types.VoteFrictionPool},
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
		Accounts:      []GenesisAccount{},
		Developers:    []GenesisAppDeveloper{},
		GenesisParam: GenesisParam{
			true,
			param.GlobalAllocationParam{
				GlobalGrowthRate:         types.NewDecFromRat(98, 1000),
				ContentCreatorAllocation: types.NewDecFromRat(65, 100),
				DeveloperAllocation:      types.NewDecFromRat(10, 100),
				ValidatorAllocation:      types.NewDecFromRat(5, 100),
			},
			param.VoteParam{
				MinStakeIn:                     types.NewCoinFromInt64(1000 * types.Decimals),
				VoterCoinReturnIntervalSec:     int64(7 * 24 * 3600),
				VoterCoinReturnTimes:           int64(7),
				DelegatorCoinReturnIntervalSec: int64(7 * 24 * 3600),
				DelegatorCoinReturnTimes:       int64(7),
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
			param.CoinDayParam{
				SecondsToRecoverCoinDay: int64(7 * 24 * 3600),
			},
			param.BandwidthParam{
				SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
				CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
				VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
				GeneralMsgQuotaRatio:        types.NewDecFromRat(20, 100),
				GeneralMsgEMAFactor:         types.NewDecFromRat(1, 10),
				AppMsgQuotaRatio:            types.NewDecFromRat(80, 100),
				AppMsgEMAFactor:             types.NewDecFromRat(1, 10),
				ExpectedMaxMPS:              types.NewDecFromRat(300, 1),
				MsgFeeFactorA:               types.NewDecFromRat(6, 1),
				MsgFeeFactorB:               types.NewDecFromRat(10, 1),
				MaxMPSDecayRate:             types.NewDecFromRat(99, 100),
				AppBandwidthPoolSize:        types.NewDecFromRat(10, 1),
				AppVacancyFactor:            types.NewDecFromRat(69, 100),
				AppPunishmentFactor:         types.NewDecFromRat(14, 5),
			},
			param.AccountParam{
				MinimumBalance:               types.NewCoinFromInt64(0),
				RegisterFee:                  types.NewCoinFromInt64(1 * types.Decimals),
				FirstDepositFullCoinDayLimit: types.NewCoinFromInt64(1 * types.Decimals),
				MaxNumFrozenMoney:            10,
			},
			param.PostParam{
				ReportOrUpvoteIntervalSec: 24 * 3600,
				PostIntervalSec:           600,
				MaxReportReputation:       types.NewCoinFromInt64(100 * types.Decimals),
			},
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
		Website:     "https://lino.network/",
		Description: "",
		AppMetaData: "",
	}
	genesisState.Developers = append(genesisState.Developers, genesisAppDeveloper)

	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	return
}
