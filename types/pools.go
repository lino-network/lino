package types

const (
	// inflation
	InflationDeveloperPool   PoolName = "inflation/developer"
	InflationValidatorPool   PoolName = "inflation/validator"
	InflationConsumptionPool PoolName = "inflation/consumption"

	// account
	AccountVestingPool PoolName = "account/vesting"

	// vote
	VoteStakeInPool     PoolName = "vote/stake-in"
	VoteStakeReturnPool PoolName = "vote/stake-return"
	VoteFrictionPool    PoolName = "vote/friction"

	// developer
	DevIDAReservePool PoolName = "dev/ida-reserve-pool"
)

func ListPools() []PoolName {
	return []PoolName{
		InflationDeveloperPool,
		InflationValidatorPool,
		InflationConsumptionPool,
		AccountVestingPool,
		VoteStakeInPool,
		VoteStakeReturnPool,
		VoteFrictionPool,
		DevIDAReservePool,
	}
}
