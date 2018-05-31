package developer

import (
	"fmt"
	"reflect"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	global "github.com/lino-network/lino/tx/global"
)

func NewHandler(dm DeveloperManager, am acc.AccountManager, gm global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case DeveloperRegisterMsg:
			return handleDeveloperRegisterMsg(ctx, dm, am, msg)
		case GrantDeveloperMsg:
			return handleGrantDeveloperMsg(ctx, dm, am, msg)
		case DeveloperRevokeMsg:
			return handleDeveloperRevokeMsg(ctx, dm, am, gm, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized developer msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleDeveloperRegisterMsg(
	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager, msg DeveloperRegisterMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrUsernameNotFound().Result()
	}

	deposit, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err.Result()
	}

	// withdraw money from developer's bank
	if err = am.MinusSavingCoin(ctx, msg.Username, deposit, types.DeveloperDeposit); err != nil {
		return err.Result()
	}
	if err := dm.RegisterDeveloper(ctx, msg.Username, deposit); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleDeveloperRevokeMsg(
	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager, gm global.GlobalManager, msg DeveloperRevokeMsg) sdk.Result {
	if !dm.IsDeveloperExist(ctx, msg.Username) {
		return ErrDeveloperNotFound().Result()
	}

	if err := dm.RemoveFromDeveloperList(ctx, msg.Username); err != nil {
		return err.Result()
	}

	coin, withdrawErr := dm.WithdrawAll(ctx, msg.Username)
	if withdrawErr != nil {
		return withdrawErr.Result()
	}

	param, err := dm.paramHolder.GetDeveloperParam(ctx)
	if err != nil {
		return err.Result()
	}

	if err := returnCoinTo(
		ctx, msg.Username, gm, am, param.DeveloperCoinReturnTimes, param.DeveloperCoinReturnIntervalHr, coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleGrantDeveloperMsg(
	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager, msg GrantDeveloperMsg) sdk.Result {
	if !dm.IsDeveloperExist(ctx, msg.AuthenticateApp) {
		return ErrDeveloperNotFound().Result()
	}
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrUsernameNotFound().Result()
	}

	if err := am.AuthorizePermission(
		ctx, msg.Username, msg.AuthenticateApp, msg.ValidityPeriod, msg.GrantLevel); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func returnCoinTo(
	ctx sdk.Context, name types.AccountKey, gm global.GlobalManager,
	am acc.AccountManager, times int64, interval int64, coin types.Coin) sdk.Error {
	if err := am.AddFrozenMoney(
		ctx, name, coin, ctx.BlockHeader().Time, interval, times); err != nil {
		return err
	}

	events, err := acc.CreateCoinReturnEvents(name, times, interval, coin, types.DeveloperReturnCoin)
	if err != nil {
		return err
	}

	if err := gm.RegisterCoinReturnEvent(ctx, events, times, interval); err != nil {
		return err
	}
	return nil
}
