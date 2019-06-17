package repv2

// This package does not support parameter hot change by now.
const (
	// unit: coin
	DefaultRoundDurationSeconds = 25 * 3600 // how many hours does game last.
	DefaultSampleWindowSize     = 10        // how many rounds is used to sample out user's customer score.
	DefaultDecayFactor          = 10        // reputation decay factor %.

	DefaultInitialReputation = 1 // initial and minimum score is 10^(-5) lino, one coin.
)
