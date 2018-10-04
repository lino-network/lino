package inflation

type Inflation struct {
	InfraPool     int64 `json:"infraPool"`
	DevPool       int64 `json:"devPool"`
	CreatorPool   int64 `json:"creatorPool"`
	ValidatorPool int64 `json:"validatorPool"`
	Timestamp     int64 `json:"timestamp"`
}
