package reward

import "time"

type Reward struct {
	ID              int64     `json:"id"`
	Username        string    `json:"username"`
	TotalIncome     string    `json:"total_income"`
	OriginalIncome  string    `json:"original_income"`
	FrictionIncome  string    `json:"friction_income"`
	InflationIncome string    `json:"inflation_income"`
	UnclaimReward   string    `json:"unclaim_reward"`
	CreatedAt       time.Time `json:"created_at"`
}
