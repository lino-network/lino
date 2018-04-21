package types

type Event interface{}

// Height -> HeightEventList
type HeightEventList struct {
	Events []Event `json:"events"`
}

// Minute -> TimeEventList
type TimeEventList struct {
	Events []Event `json:"events"`
}
