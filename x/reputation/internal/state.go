package internal

// UserReputation - pk: Username
type UserReputation struct {
	Username      Uid `json:"username"`
	CustomerScore Rep `json:"customer_score"`
	FreeScore     Rep `json:"free_score"`
}

// UserReputationTable -
type UserReputationTable struct {
	Reputations []UserReputation `json:"reputations"`
}
