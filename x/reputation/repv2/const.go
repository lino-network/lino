package repv2

// This package does not support parameter hot change by now.
const (
	DefaultBestContentIndexN = 5
	// unit: coin
	RoundDuration           = 25 // how many hours does game last.
	SampleWindowSize        = 10 // how many rounds is used to sample out user's customer score.
	DecayFactor             = 10 // reputation decay factor %.

	InitialReputation = 1 // initial and minimum score is 10^(-5) lino, one coin.
)
