package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/lino-network/lino/types"
)

type ReturnCoinEvent struct {
	Username types.AccountKey `json:"username"`
	Amount   types.Coin       `json:"amount"`
}

func (event ReturnCoinEvent) Execute(ctx sdk.Context, am AccountManager) sdk.Error {
	if !am.IsAccountExist(ctx, event.Username) {
		return ErrUsernameNotFound()
	}

	if err := am.AddCoin(ctx, event.Username, event.Amount); err != nil {
		return err
	}
	return nil
}
