package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/types"
)

// ReturnCoinEvent - return a certain amount of coin to an account
type ReturnCoinEvent struct {
	Username   linotypes.AccountKey         `json:"username"`
	Amount     linotypes.Coin               `json:"amount"`
	ReturnType linotypes.TransferDetailType `json:"return_type"`
}

// Execute - execute coin return events
func (event ReturnCoinEvent) Execute(ctx sdk.Context, am AccountManager) sdk.Error {
	addr, err := am.GetAddress(ctx, event.Username)
	if err != nil {
		return types.ErrAccountNotFound(event.Username)
	}

	if err := am.AddCoinToAddress(ctx, addr, event.Amount); err != nil {
		return err
	}
	return nil
}

// CreateCoinReturnEvents - create coin return events
func CreateCoinReturnEvents(
	ctx sdk.Context, username linotypes.AccountKey, times int64, interval int64, coin linotypes.Coin,
	returnType linotypes.TransferDetailType) ([]linotypes.Event, sdk.Error) {
	events := []linotypes.Event{}
	for i := int64(0); i < times; i++ {
		pieceRat := coin.ToDec().Quo(sdk.NewDec(times - i))
		piece := linotypes.DecToCoin(pieceRat)
		coin = coin.Minus(piece)

		event := ReturnCoinEvent{
			Username:   username,
			Amount:     piece,
			ReturnType: returnType,
		}
		events = append(events, event)
	}
	return events, nil
}
