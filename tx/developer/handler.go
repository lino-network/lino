package developer

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

func NewHandler(dm DeveloperManager, am acc.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case DeveloperRegisterMsg:
			return handleDeveloperRegisterMsg(ctx, dm, am, msg)
		case DeveloperRevokeMsg:
			return handleDeveloperRevokeMsg(ctx, dm, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized developer Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleDeveloperRegisterMsg(ctx sdk.Context, dm DeveloperManager, am acc.AccountManager, msg DeveloperRegisterMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrUsernameNotFound().Result()
	}

	deposit, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err.Result()
	}

	if err := dm.RegisterDeveloper(ctx, msg.Username, deposit); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleDeveloperRevokeMsg(ctx sdk.Context, dm DeveloperManager, msg DeveloperRevokeMsg) sdk.Result {
	if !dm.IsDeveloperExist(ctx, msg.Username) {
		return ErrDeveloperNotFound().Result()
	}

	if err := dm.RemoveFromDeveloperList(ctx, msg.Username); err != nil {
		return err.Result()
	}

	if err := dm.WithdrawAll(ctx, msg.Username); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
