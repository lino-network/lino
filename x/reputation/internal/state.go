package internal

// UserReputation - pk: Username
type UserReputation struct {
	Username          Uid
	CustomerScore     Rep
	FreeScore         Rep
}

// UserReputationTable -
type UserReputationTable struct {
	reputations []UserReputation
}
