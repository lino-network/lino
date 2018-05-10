package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/lino-network/lino/types"
)

// ReturnCoin Event return a certain amount of coin to an account
type ReturnCoinEvent struct {
	Username types.AccountKey `json:"username"`
	Amount   types.Coin       `json:"amount"`
}

// execute return coin event
func (event ReturnCoinEvent) Execute(ctx sdk.Context, am AccountManager) sdk.Error {
	if !am.IsAccountExist(ctx, event.Username) {
		return ErrUsernameNotFound()
	}

	if err := am.AddSavingCoin(ctx, event.Username, event.Amount); err != nil {
		return err
	}
	return nil
}
