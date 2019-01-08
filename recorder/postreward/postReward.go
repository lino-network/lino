package postreward

type PostReward struct {
	ID           int64  `json:"id"`
	Permlink     string `json:"permlink"`
	Reward       int64  `json:"reward"`
	PenaltyScore string `json:"penaltyScore"`
	Timestamp    int64  `json:"timestamp"`
	Evaluate     int64  `json:"evaluate"`
	Original     int64  `json:"original"`
	Consumer     string `json:"consumer"`
}
