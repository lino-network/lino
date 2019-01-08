package stake

type Stake struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Amount    string `json:"amount"`
	Timestamp int64  `json:"timestamp"`
	Op        string `json:"op"`
}
