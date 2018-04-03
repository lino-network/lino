package global

import (
	"strconv"

	types "github.com/lino-network/lino/types"
)

type EventListKey string

type Event interface{}

// Height -> HeightEventList
type HeightEventList struct {
	Events []Event `json:"events"`
}

// Minute -> TimeEventList
type TimeEventList struct {
	Events []Event `json:"events"`
}

func HeightToEventListKey(height types.Height) EventListKey {
	return EventListKey(strconv.FormatInt(int64(height), 10))
}

func UnixTimeToEventListKey(unixTime int64) EventListKey {
	return EventListKey(strconv.FormatInt(unixTime, 10))
}
