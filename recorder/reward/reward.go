package reward

import "time"

type Reward struct {
	ID              int64     `json:"id"`
	Username        string    `json:"username"`
	TotalIncome     int64     `json:"total_income"`
	OriginalIncome  int64     `json:"original_income"`
	FrictionIncome  int64     `json:"friction_income"`
	InflationIncome int64     `json:"inflation_income"`
	UnclaimReward   int64     `json:"unclaim_reward"`
	CreatedAt       time.Time `json:"created_at"`
}
