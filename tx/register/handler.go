package register

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

func NewHandler(am acc.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case RegisterMsg:
			return handleRegisterMsg(ctx, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle RegisterMsg
func handleRegisterMsg(ctx sdk.Context, am acc.AccountManager, msg RegisterMsg) sdk.Result {
	if err := am.CreateAccount(ctx, msg.NewUser, msg.NewPubKey, types.NewCoin(100*types.Decimals)); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
