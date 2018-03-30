package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GlobalMeta struct {
	TotalLino        sdk.Coins `json:"total_lino"`
	TotalConsumption sdk.Rat   `json:"total_consumption"`
	GrowthRate       sdk.Rat   `json:"growth_rate"`
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
	TotalLinoInflationPool   sdk.Rat `json:"total_lino_inflation_pool"`
	InfraAllocation          sdk.Rat `json:"infra_allocation"`
	ContentCreatorAllocation sdk.Rat `json:"content_creator_allocation"`
	DeveloperAllocation      sdk.Rat `json:"developer_allocation"`
	ValidatorAllocation      sdk.Rat `json:"validator_allocation"`
}

type InfraInternalAllocation struct {
	StorageAllocation sdk.Rat `json:"storage_allocation"`
	CDNAllocation     sdk.Rat `json:"CDN_allocation"`
}

// ConsumptionMeta
// ConsumptionFrictionRate: percentage the user consumption deducted and added to the TotalLinoInflationPool
// ReportStakeWindow used to evaluate the panelty of the post within the freezing period
// DislikeStakeWindow used to evaluate the panelty of the post within the freezing period
// ConsumptionWindow records all content related consumption within the freezing period
// FreezingPeriodHr is the time content createor can get remain consumption after friction
type ConsumptionMeta struct {
	ConsumptionFrictionRate sdk.Rat `json:"consumption_friction_rate"`
	ReportStakeWindow       sdk.Rat `json:"report_stake_window"`
	DislikeStakeWindow      sdk.Rat `json:"dislike_stake_window"`
	ConsumptionWindow       sdk.Rat `json:"consumption_window"`
	FreezingPeriodHr        int64   `json:"freezing_period"`
}
