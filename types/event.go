package types

// Event - event executed in app.go
type Event interface{}

// Minute -> TimeEventList
type TimeEventList struct {
	Events []Event `json:"events"`
}
