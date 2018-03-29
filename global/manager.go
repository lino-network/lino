package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"
)

var EventListPrefix = []byte("EventList/")

const eventTypePostReward = 0x1
const eventTypeDonateReward = 0x2

var _ = oldwire.RegisterInterface(
	struct{ Event }{},
	oldwire.ConcreteType{PostRewardEvent{}, eventTypePostReward},
	oldwire.ConcreteType{DonateRewardEvent{}, eventTypeDonateReward},
)

type GlobalManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewGlobalManager(key sdk.StoreKey) GlobalManager {
	cdc := wire.NewCodec()
	gm := GlobalManager{
		key: key,
		cdc: cdc,
	}
	return gm
}

func (gm GlobalManager) GetEventList(ctx sdk.Context, key EventListKey) (*EventList, sdk.Error) {
	store := ctx.KVStore(gm.key)
	listByte := store.Get(eventListKey(key))
	if listByte == nil {
		return nil, ErrEventNotFound(eventListKey(key))
	}
	lst := new(EventList)
	if err := gm.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return lst, nil
}

func (gm GlobalManager) SetEventList(ctx sdk.Context, key EventListKey, lst *EventList) sdk.Error {
	store := ctx.KVStore(gm.key)
	listByte, err := gm.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(eventListKey(key), listByte)
	return nil
}

func (gm GlobalManager) removeEventList(ctx sdk.Context, key EventListKey) sdk.Error {
	store := ctx.KVStore(gm.key)
	store.Delete(eventListKey(key))
	return nil
}

func (gm GlobalManager) ExecuteEvents(ctx sdk.Context, key EventListKey) sdk.Error {
	lst, err := gm.GetEventList(ctx, key)
	if err != nil {
		return err
	}

	for _, event := range lst.Events {
		switch event := event.(type) {
		case PostRewardEvent:
			if err := event.execute(); err != nil {
				return err
			}
		case DonateRewardEvent:
			if err := event.execute(); err != nil {
				return err
			}
		default:
			return ErrWrongEventType()
		}
	}

	if err := gm.removeEventList(ctx, key); err != nil {
		return err
	}
	return nil
}

func eventListKey(eventListKey EventListKey) []byte {
	return append(EventListPrefix, eventListKey...)
}
