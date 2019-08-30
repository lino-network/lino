package repv2

const (
	// Inherited from testnet, the unit of 1 reputation is one coin
	// of testnet, which is 10^(-5) * 0.012 USD.
	// Caller need to convert the amount of donation to the number of test coins.
	DefaultRoundDurationSeconds = 25 * 3600 // how many seconds does a round last, default: 25 hours
	DefaultSampleWindowSize     = 10        // how many rounds are used to sample out user's customer score.
	DefaultDecayFactor          = 10        // reputation decay factor %.

	// Initial and minimum score is 10^(-5), one coin.
	DefaultInitialReputation = 1
)
