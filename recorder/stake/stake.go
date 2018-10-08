package stake

type Stake struct {
	Username  string `json:"username"`
	Amount    int64  `json:"amount"`
	Timestamp int64  `json:"timestamp"`
	Op        string `json:"op"`
}
