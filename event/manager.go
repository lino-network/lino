package event

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

type EventManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	cdc *wire.Codec
}

func NewEventManager(key sdk.StoreKey) EventManager {
	cdc := wire.NewCodec()
	em := EventManager{
		key: key,
		cdc: cdc,
	}
	return em
}

func (em EventManager) GetEventList(ctx sdk.Context, key EventListKey) (*EventList, sdk.Error) {
	store := ctx.KVStore(em.key)
	listByte := store.Get(eventListKey(key))
	if listByte == nil {
		return nil, ErrEventNotFound(eventListKey(key))
	}
	lst := new(EventList)
	if err := em.cdc.UnmarshalJSON(listByte, lst); err != nil {
		return nil, ErrEventUnmarshalError(err)
	}
	return lst, nil
}

func (em EventManager) SetEventList(ctx sdk.Context, key EventListKey, lst *EventList) sdk.Error {
	store := ctx.KVStore(em.key)
	listByte, err := em.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrEventMarshalError(err)
	}
	store.Set(eventListKey(key), listByte)
	return nil
}

func (em EventManager) removeEventList(ctx sdk.Context, key EventListKey) sdk.Error {
	store := ctx.KVStore(em.key)
	store.Delete(eventListKey(key))
	return nil
}

func (em EventManager) ExecuteEvents(ctx sdk.Context, key EventListKey) sdk.Error {
	lst, err := em.GetEventList(ctx, key)
	if err != nil {
		return err
	}

	for _, event := range lst.Events {
		if err := event.execute(); err != nil {
			return err
		}
	}

	if err := em.removeEventList(ctx, key); err != nil {
		return err
	}
	return nil
}

func eventListKey(eventListKey EventListKey) []byte {
	return append(EventListPrefix, eventListKey...)
}
