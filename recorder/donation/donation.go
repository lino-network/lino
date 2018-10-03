package donation

type Donation struct {
	Username       string `json:"username"`
	Seq            int64  `json:"seq"`
	Dp             int64  `json:"dp"`
	Permlink       string `json:"permlink"`
	Amount         int64  `json:"amount"`
	FromApp        string `json:"fromApp"`
	CoinDayDonated int64  `json:"coinDayDonated"`
	Reputation     int64  `json:"reputation"`
	Timestamp      int64  `json:"timestamp"`
	EvaluateResult int64  `json:"evaluateResult"`
}
