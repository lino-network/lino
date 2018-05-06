package account

import (
	"github.com/cosmos/cosmos-sdk/wire"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/lino-network/lino/types"
)

func init() {
	cdc := wire.NewCodec()

	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(ReturnCoinEvent{}, "event/return", nil)
}

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

	if err := am.AddCoin(ctx, event.Username, event.Amount); err != nil {
		return err
	}
	return nil
}
