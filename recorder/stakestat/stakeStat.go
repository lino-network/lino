package stakestat

type StakeStat struct {
	TotalConsumptionFriction int64  `json:"totalConsumptionFriction"`
	UnclaimedFriction        int64  `json:"unclaimedFriction"`
	TotalLinoStake           string `json:"totalLinoStake"`
	UnclaimedLinoStake       string `json:"unclaimedLinoStake"`
	Timestamp                int64  `json:"timestamp"`
}
