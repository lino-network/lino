package exporter

// AppState contains all informations that needs to migrate when blockchain upgrade.
type AppState struct {
	// Accounts            []AccountRow     `json:"accounts"`
	// AccountGrantPubKeys []GrantPubKeyRow `json:"account_grant_pub_keys"`

	// Developers           []Developer
	// DeveloperList        DeveloperList

	// GlobalTimeEventLists []GlobalTimeEventList
	// GlobalStakeStats     []GlobalStakeStat
	// GlobalMisc           GlobalMisc

	// InfraProviders       []InfraProvider
	// InfraProviderList    InfraProviderList

	// Posts                []Post
	// PostUsers            []PostUser
	// PostComments         []PostComment

	// proposal is skipped
	// vote is skipped

	// Validators    []Validator
	// ValidatorList ValidatorList

	// reputaion has pure math model implementation and can be simply marshal to bytes out.
	// reputation []byte
}
