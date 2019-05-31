package repv2

// UserReputation - pk: Username
type UserReputation struct {
	Username      Uid `json:"username"`
	CustomerScore Rep `json:"customer_score"`
	FreeScore     Rep `json:"free_score"`
}

// UserReputationTable - pk by Username
type UserReputationTable struct {
	Reputations []UserReputation `json:"reputations"`
}
