package genesis

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// State to Unmarshal
type GenesisState struct {
	Accounts    []GenesisAccount `json:"accounts"`
	GlobalState GlobalState      `json:"global_state"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Name        string        `json:"name"`
	Lino        int64         `json:"lino"`
	PubKey      crypto.PubKey `json:"pub_key"`
	IsValidator bool          `json:"is_validator"`
	ValPubKey   crypto.PubKey `json:"validator_pub_key"`
}

type GlobalState struct {
	TotalLino                int64   `json:"total_lino"`
	GrowthRate               sdk.Rat `json:"growth_rate"`
	InfraAllocation          sdk.Rat `json:"infra_allocation"`
	ContentCreatorAllocation sdk.Rat `json:"content_creator_allocation"`
	DeveloperAllocation      sdk.Rat `json:"developer_allocation"`
	ValidatorAllocation      sdk.Rat `json:"validator_allocation"`
	ConsumptionFrictionRate  sdk.Rat `json:"consumption_friction_rate"`
	FreezingPeriodHr         int64   `json:"freezing_period_hr"`
}

func GetGenesisJson(genesisState GenesisState) (string, error) {
	output, err := json.MarshalIndent(genesisState, "", "\t")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func GetDefaultGenesis(pubkey crypto.PubKey, validatorPubKey crypto.PubKey) (string, error) {
	totalLino := int64(10000000000)
	genesisAcc := GenesisAccount{
		Name:        "Lino",
		Lino:        totalLino,
		PubKey:      pubkey,
		IsValidator: true,
		ValPubKey:   validatorPubKey,
	}
	globalState := GlobalState{
		TotalLino:                totalLino,
		GrowthRate:               sdk.NewRat(98, 1000),
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(50, 100),
		DeveloperAllocation:      sdk.NewRat(20, 100),
		ValidatorAllocation:      sdk.NewRat(10, 100),
		ConsumptionFrictionRate:  sdk.NewRat(1, 100),
		FreezingPeriodHr:         24 * 7,
	}
	genesisState := GenesisState{
		Accounts:    []GenesisAccount{genesisAcc},
		GlobalState: globalState,
	}

	result, err := GetGenesisJson(genesisState)
	if err != nil {
		return "", err
	}
	return result, err
}
