package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/lino-network/lino/types"
)

// ReturnCoin Event return a certain amount of coin to an account
type ReturnCoinEvent struct {
	Username   types.AccountKey         `json:"username"`
	Amount     types.Coin               `json:"amount"`
	ReturnType types.TransferDetailType `json:"return_type"`
}

// execute return coin event
func (event ReturnCoinEvent) Execute(ctx sdk.Context, am AccountManager) sdk.Error {
	if !am.DoesAccountExist(ctx, event.Username) {
		return ErrAccountNotFound(event.Username)
	}

	if err := am.AddSavingCoin(
		ctx, event.Username, event.Amount, "", "",
		event.ReturnType); err != nil {
		return err
	}
	return nil
}

// create coin return events
func CreateCoinReturnEvents(
	username types.AccountKey, times int64, interval int64, coin types.Coin,
	returnType types.TransferDetailType) ([]types.Event, sdk.Error) {
	events := []types.Event{}
	for i := int64(0); i < times; i++ {
		pieceRat := coin.ToRat().Quo(sdk.NewRat(times-i, 1))
		piece, err := types.RatToCoin(pieceRat)
		if err != nil {
			return nil, err
		}
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
