package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GlobalMetaKey string

// Rete uses 0 to 100 to present 0% to 100%
type GlobalMeta struct {
	RedistriSplitRate       int64 `json:"redistribution_split_rate"`
	ConsumptionFrictionRate int64 `json:"consumption_friction_rate"`
}

// Rete uses 0 to 100 to present 0% to 100%
type ConsumpotionMeta struct {
	TotalReportStake  int64     `json:"total_report_stake"`
	TotalDisLikeStake int64     `json:"total_like_stake"`
	ConsumpotionPool  sdk.Coins `json:"consumption_pool"`
}

func GlobalMetaPrefixWithKey() GlobalMetaKey {
	return GlobalMetaKey("GlobalMetaKey")
}
