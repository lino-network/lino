package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

// ReturnCoinEvent - return a certain amount of coin to an account
type ReturnCoinEvent struct {
	Username   linotypes.AccountKey         `json:"username"`
	Amount     linotypes.Coin               `json:"amount"`
	ReturnType linotypes.TransferDetailType `json:"return_type"`
	FromPool   linotypes.PoolName           `json:"from_pool"`
	At         int64                        `json:"at"`
}

// Execute - execute coin return events
func (event ReturnCoinEvent) Execute(ctx sdk.Context, am AccountManager) sdk.Error {
	err := am.AddPending(ctx, username, event.Amount.Neg())
	if err != nil {
		return err
	}
	return am.MoveFromPool(
		ctx, event.FromPool, linotypes.NewAccOrAddrFromAcc(event.Username), event.Amount)
}

// CreateCoinReturnEvents - create coin return events
// The return interval list is expected to be executed at [start + interval, start + 2 * interval...]
// If [start, start + interval...] is expected, pass int (startAt - interval) as start at instead.
func CreateCoinReturnEvents(username linotypes.AccountKey, startAt, interval, times int64, coin linotypes.Coin, returnType linotypes.TransferDetailType, pool linotypes.PoolName) []ReturnCoinEvent {
	events := []ReturnCoinEvent{}
	for i := int64(0); i < times; i++ {
		pieceDec := coin.ToDec().Quo(sdk.NewDec(times - i))
		piece := linotypes.DecToCoin(pieceDec)
		coin = coin.Minus(piece)

		event := ReturnCoinEvent{
			Username:   username,
			Amount:     piece,
			ReturnType: returnType,
			FromPool:   pool,
			At:         startAt + (i+1)*interval,
		}
		events = append(events, event)
	}
	return events
}
