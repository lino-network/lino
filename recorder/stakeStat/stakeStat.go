package stakestat

type StakeStat struct {
	TotalConsumptionFriction int64 `json:"totalConsumptionFriction"`
	UnclaimedFriction        int64 `json:"unclaimedFriction"`
	TotalLinoStake           int64 `json:"totalLinoStake"`
	UnclaimedLinoStake       int64 `json:"unclaimedLinoStake"`
	Timestamp                int64 `json:"timestamp"`
}
