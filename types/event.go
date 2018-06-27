package types

type Event interface{}

// Minute -> TimeEventList
type TimeEventList struct {
	Events []Event `json:"events"`
}
