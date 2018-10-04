package inflation

type Inflation struct {
	InfraPool          int64 `json:"infraPool"`
	DevPool            int64 `json:"devPool"`
	CreatorPool        int64 `json:"creatorPool"`
	ValidatorPool      int64 `json:"validatorPool"`
	InfraInflation     int64 `json:"infraInflation"`
	DevInflation       int64 `json:"devInflation"`
	CreatorInflation   int64 `json:"creatorInflation"`
	ValidatorInflation int64 `json:"validatorInflation"`
	Timestamp          int64 `json:"timestamp"`
}
