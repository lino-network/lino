package genesis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// State to Unmarshal
type GenesisState struct {
	Accounts    []*GenesisAccount `json:"accounts"`
	GlobalState GlobalState       `json:"global_state"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Name      string        `json:"name"`
	Lino      int64         `json:"lino"`
	PubKey    crypto.PubKey `json:"pub_key"`
	ValPubKey crypto.PubKey `json:"validator_pub_key"`
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
