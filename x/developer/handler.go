package developer

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/developer/types"
)

// NewHandler - Handle all "developer" type messages.
func NewHandler(dm DeveloperKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.DeveloperRegisterMsg:
			return handleDeveloperRegisterMsg(ctx, dm, msg)
		case types.DeveloperUpdateMsg:
			return handleDeveloperUpdateMsg(ctx, dm, msg)
		// case types.DeveloperRevokeMsg:
		// 	return handleDeveloperRevokeMsg(ctx, dm, am, gm, msg)
		case types.IDAIssueMsg:
			return handleIDAIssueMsg(ctx, dm, msg)
		case types.IDAMintMsg:
			return handleIDAMintMsg(ctx, dm, msg)
		case types.IDATransferMsg:
			return handleIDATransferMsg(ctx, dm, msg)
		case types.IDAAuthorizeMsg:
			return handleIDAAuthorizeMsg(ctx, dm, msg)
		case types.UpdateAffiliatedMsg:
			return handleUpdateAffiliatedMsg(ctx, dm, msg)
		case types.IDAConvertFromLinoMsg:
			if ctx.BlockHeight() >= linotypes.Upgrade5Update2 {
				return handleIDAConvertFromLinoMsg(ctx, dm, msg)
			} else {
				errMsg := fmt.Sprintf(
					"Unrecognized developer msg type: %v", reflect.TypeOf(msg).Name())
				return sdk.ErrUnknownRequest(errMsg).Result()
			}
		default:
			errMsg := fmt.Sprintf("Unrecognized developer msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleDeveloperRegisterMsg(
	ctx sdk.Context, dm DeveloperKeeper, msg types.DeveloperRegisterMsg) sdk.Result {
	if err := dm.RegisterDeveloper(
		ctx, msg.Username, msg.Website, msg.Description, msg.AppMetaData); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleDeveloperUpdateMsg(
	ctx sdk.Context, dm DeveloperKeeper, msg types.DeveloperUpdateMsg) sdk.Result {
	if err := dm.UpdateDeveloper(
		ctx, msg.Username, msg.Website, msg.Description, msg.AppMetaData); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleUpdateAffiliatedMsg(
	ctx sdk.Context, dm DeveloperKeeper, msg types.UpdateAffiliatedMsg) sdk.Result {
	if err := dm.UpdateAffiliated(ctx, msg.App, msg.Username, msg.Activate); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleIDAIssueMsg(
	ctx sdk.Context, dm DeveloperKeeper, msg types.IDAIssueMsg) sdk.Result {
	if err := dm.IssueIDA(ctx, msg.Username, string(msg.Username), msg.IDAPrice); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleIDAMintMsg(
	ctx sdk.Context, dm DeveloperKeeper, msg types.IDAMintMsg) sdk.Result {
	amount, err := linotypes.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if err := dm.MintIDA(ctx, msg.Username, amount); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleIDAConvertFromLinoMsg(
	ctx sdk.Context, dm DeveloperKeeper, msg types.IDAConvertFromLinoMsg) sdk.Result {
	amount, err := linotypes.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if err := dm.IDAConvertFromLino(ctx, msg.Username, msg.App, amount); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleIDATransferMsg(
	ctx sdk.Context, dm DeveloperKeeper, msg types.IDATransferMsg) sdk.Result {
	amount, err := msg.Amount.ToMiniIDA()
	if err != nil {
		return err.Result()
	}
	if err := dm.AppTransferIDA(ctx, msg.App, msg.Signer, amount, msg.From, msg.To); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleIDAAuthorizeMsg(ctx sdk.Context, dm DeveloperKeeper, msg types.IDAAuthorizeMsg) sdk.Result {
	if err := dm.UpdateIDAAuth(ctx, msg.App, msg.Username, msg.Activate); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// func handleDeveloperRevokeMsg(
// 	ctx sdk.Context, dm DeveloperManager, am acc.AccountManager,
// 	gm *global.GlobalManager, msg types.DeveloperRevokeMsg) sdk.Result {
// 	if !dm.DoesDeveloperExist(ctx, msg.Username) {
// 		return types.ErrDeveloperNotFound().Result()
// 	}

// 	if err := dm.RemoveFromDeveloperList(ctx, msg.Username); err != nil {
// 		return err.Result()
// 	}

// 	coin, withdrawErr := dm.WithdrawAll(ctx, msg.Username)
// 	if withdrawErr != nil {
// 		return withdrawErr.Result()
// 	}

// 	param, err := dm.paramHolder.GetDeveloperParam(ctx)
// 	if err != nil {
// 		return err.Result()
// 	}

// 	if err := returnCoinTo(
// 		ctx, msg.Username, gm, am, param.DeveloperCoinReturnTimes, param.DeveloperCoinReturnIntervalSec, coin); err != nil {
// 		return err.Result()
// 	}
// 	return sdk.Result{}
// }
