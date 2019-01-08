package inflation

type Inflation struct {
	ID                 int64  `json:"id"`
	InfraPool          string `json:"infraPool"`
	DevPool            string `json:"devPool"`
	CreatorPool        string `json:"creatorPool"`
	ValidatorPool      string `json:"validatorPool"`
	InfraInflation     string `json:"infraInflation"`
	DevInflation       string `json:"devInflation"`
	CreatorInflation   string `json:"creatorInflation"`
	ValidatorInflation string `json:"validatorInflation"`
	Timestamp          int64  `json:"timestamp"`
}
