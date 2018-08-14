package developer

import (
	"fmt"
	"reflect"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	global "github.com/lino-network/lino/x/global"
)

func NewHandler(dm DeveloperManager, am acc.AccountManager, gm global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case DeveloperRegisterMsg:
			return handleDeveloperRegisterMsg(ctx, dm, am, msg)
		case DeveloperUpdateMsg:
			return handleDeveloperUpdateMsg(ctx, dm, am, msg)
		case GrantPermissionMsg:
			return handleGrantPermissionMsg(ctx, dm, am, msg)
		case PreAuthorizationMsg:
			return handlePreAuthorizationMsg(ctx, dm, am, msg)
		case DeveloperRevokeMsg:
			return handleDeveloperRevokeMsg(ctx, dm, am, gm, msg)
		case RevokePermissionMsg:
			return handleRevokePermissionMsg(ctx, dm, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized developer msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleDeveloperRegisterMsg(
	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager, msg DeveloperRegisterMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.Username) {
		return ErrAccountNotFound().Result()
	}

	if dm.DoesDeveloperExist(ctx, msg.Username) {
		return ErrDeveloperAlreadyExist(msg.Username).Result()
	}

	deposit, err := types.LinoToCoin(msg.Deposit)
	if err != nil {
		return err.Result()
	}

	// withdraw money from developer's bank
	if err = am.MinusSavingCoin(
		ctx, msg.Username, deposit, "", "", types.DeveloperDeposit); err != nil {
		return err.Result()
	}
	if err := dm.RegisterDeveloper(
		ctx, msg.Username, deposit, msg.Website, msg.Description, msg.AppMetaData); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleDeveloperUpdateMsg(
	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager, msg DeveloperUpdateMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.Username) {
		return ErrAccountNotFound().Result()
	}

	if !dm.DoesDeveloperExist(ctx, msg.Username) {
		return ErrDeveloperNotFound().Result()
	}

	if err := dm.UpdateDeveloper(
		ctx, msg.Username, msg.Website, msg.Description, msg.AppMetaData); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleDeveloperRevokeMsg(
	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager,
	gm global.GlobalManager, msg DeveloperRevokeMsg) sdk.Result {
	if !dm.DoesDeveloperExist(ctx, msg.Username) {
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

func handleGrantPermissionMsg(
	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager, msg GrantPermissionMsg) sdk.Result {
	if !dm.DoesDeveloperExist(ctx, msg.AuthorizedApp) {
		return ErrDeveloperNotFound().Result()
	}
	if !am.DoesAccountExist(ctx, msg.Username) {
		return ErrAccountNotFound().Result()
	}

	if err := am.AuthorizePermission(
		ctx, msg.Username, msg.AuthorizedApp, msg.ValidityPeriodSec, msg.GrantLevel, types.NewCoinFromInt64(0)); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleRevokePermissionMsg(
	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager, msg RevokePermissionMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.Username) {
		return ErrAccountNotFound().Result()
	}

	if err := am.RevokePermission(ctx, msg.Username, msg.PubKey); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handlePreAuthorizationMsg(
	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager, msg PreAuthorizationMsg) sdk.Result {
	if !dm.DoesDeveloperExist(ctx, msg.AuthorizedApp) {
		return ErrDeveloperNotFound().Result()
	}
	if !am.DoesAccountExist(ctx, msg.Username) {
		return ErrAccountNotFound().Result()
	}

	amount, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}

	if err := am.AuthorizePermission(
		ctx, msg.Username, msg.AuthorizedApp, msg.ValidityPeriodSec, types.PreAuthorizationPermission, amount); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func returnCoinTo(
	ctx sdk.Context, name types.AccountKey, gm global.GlobalManager,
	am acc.AccountManager, times int64, interval int64, coin types.Coin) sdk.Error {
	if err := am.AddFrozenMoney(
		ctx, name, coin, ctx.BlockHeader().Time.Unix(), interval, times); err != nil {
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
