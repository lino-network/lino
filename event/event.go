package event

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/lino-network/lino/types"
)

type EventListKey string

type Event interface {
	execute() sdk.Error
}

type EventList struct {
	Events []Event `json:"events"`
}

type PostRewardEvent struct {
	PostID int64
}

type DonateRewardEvent struct {
	DonateID int64
}

func (e PostRewardEvent) execute() sdk.Error {
	fmt.Println("Execute post reward event")
	return nil
}

func (e DonateRewardEvent) execute() sdk.Error {
	fmt.Println("Execute donate reward event")
	return nil
}

func HeightToEventListKey(height types.Height) EventListKey {
	return EventListKey(strconv.FormatInt(int64(height), 10))
}
