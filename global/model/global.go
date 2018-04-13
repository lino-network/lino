package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type GlobalMeta struct {
	TotalLino             types.LNO  `json:"total_lino"`
	CumulativeConsumption types.Coin `json:"cumulative_consumption"`
	GrowthRate            sdk.Rat    `json:"growth_rate"`
}

type TPS struct {
	CurrentTPS sdk.Rat `json:"current_tps"`
	MaxTPS     sdk.Rat `json:"max_tps"`
}

// GlobalStatistics
type GlobalStatistics struct {
	numOfAccount  int64 `json:"number_of_account"`
	numOfPost     int64 `json:"number_of_post"`
	numOfComment  int64 `json:"number_of_comment"`
	numOfTransfer int64 `json:"number_of_transfer"`
	numOfLike     int64 `json:"number_of_like"`
	numOfDonation int64 `json:"number_of_donation"`
}

// GlobalAllocation
// TotalLinoInflationPool: total Lino inflation for all community roles
// InfraAllocation percentage for all infra related allocation
// ContentCreatorAllocation percentage for all content creator related allocation
// DeveloperAllocation percentage of inflation for developers
// ValidatorAllocation percentage of inflation for validators
type GlobalAllocation struct {
	InfraAllocation          sdk.Rat `json:"infra_allocation"`
	ContentCreatorAllocation sdk.Rat `json:"content_creator_allocation"`
	DeveloperAllocation      sdk.Rat `json:"developer_allocation"`
	ValidatorAllocation      sdk.Rat `json:"validator_allocation"`
}

type InfraInternalAllocation struct {
	StorageAllocation sdk.Rat `json:"storage_allocation"`
	CDNAllocation     sdk.Rat `json:"CDN_allocation"`
}

// InflationPool, determined by GlobalAllocation
// InfraInflationPool inflation pool for infra
// ContentCreatorInflationPool inflation pool for content creator
// DeveloperInflationPool inflation pool for developer
// ValidatorInflationPool inflation pool for validator
type InflationPool struct {
	InfraInflationPool          types.Coin `json:"infra_inflation_pool"`
	ContentCreatorInflationPool types.Coin `json:"content_creator_inflation_pool"`
	DeveloperInflationPool      types.Coin `json:"developer_inflation_pool"`
	ValidatorInflationPool      types.Coin `json:"validator_inflation_pool"`
}

// ConsumptionMeta
// ConsumptionFrictionRate: percentage the user consumption deducted and added to the TotalLinoInflationPool
// ReportStakeWindow used to evaluate the panelty of the post within the freezing period
// DislikeStakeWindow used to evaluate the panelty of the post within the freezing period
// ConsumptionWindow records all content related consumption within the freezing period
// FreezingPeriodHr is the time content createor can get remain consumption after friction
type ConsumptionMeta struct {
	ConsumptionFrictionRate sdk.Rat    `json:"consumption_friction_rate"`
	ReportStakeWindow       sdk.Rat    `json:"report_stake_window"`
	DislikeStakeWindow      sdk.Rat    `json:"dislike_stake_window"`
	ConsumptionWindow       types.Coin `json:"consumption_window"`
	ConsumptionRewardPool   types.Coin `json:"consumption_window"`
	FreezingPeriodHr        int64      `json:"freezing_period"`
}
