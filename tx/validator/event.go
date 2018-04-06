package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	types "github.com/lino-network/lino/types"
)

type ReturnCoinEvent struct {
	Username types.AccountKey `json:"username"`
	Amount   types.Coin       `json:"amount"`
}

func (event ReturnCoinEvent) Execute(ctx sdk.Context, vm ValidatorManager, am acc.AccountManager, gm global.GlobalManager) sdk.Error {
	if !am.IsAccountExist(ctx, event.Username) {
		return acc.ErrUsernameNotFound()
	}

	if err := am.AddCoin(ctx, event.Username, event.Amount); err != nil {
		return err
	}
	return nil
}
