package account

import (
	"fmt"
	"reflect"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler - Handle all "account" type messages.
func NewHandler(am AccountKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.TransferMsg:
			return handleTransferMsg(ctx, am, msg)
		case types.TransferV2Msg:
			return handleTransferV2Msg(ctx, am, msg)
		case types.RecoverMsg:
			return handleRecoverMsg(ctx, am, msg)
		case types.RegisterV2Msg:
			return handleRegisterV2Msg(ctx, am, msg)
		case types.UpdateAccountMsg:
			return handleUpdateAccountMsg(ctx, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleTransferMsg(ctx sdk.Context, am AccountKeeper, msg types.TransferMsg) sdk.Result {
	// withdraw money from sender's bank
	coin, err := linotypes.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if err := am.MoveCoin(ctx,
		linotypes.NewAccOrAddrFromAcc(msg.Sender),
		linotypes.NewAccOrAddrFromAcc(msg.Receiver),
		coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleTransferV2Msg(ctx sdk.Context, am AccountKeeper, msg types.TransferV2Msg) sdk.Result {
	// withdraw money from sender's bank
	coin, err := linotypes.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if err := am.MoveCoin(ctx, msg.Sender, msg.Receiver, coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleRecoverMsg(ctx sdk.Context, am AccountKeeper, msg types.RecoverMsg) sdk.Result {
	if err := am.RecoverAccount(ctx, msg.Username, msg.NewTxPubKey, msg.NewSigningPubKey); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle RegisterV2Msg
func handleRegisterV2Msg(ctx sdk.Context, am AccountKeeper, msg types.RegisterV2Msg) sdk.Result {
	coin, err := linotypes.LinoToCoin(msg.RegisterFee)
	if err != nil {
		return err.Result()
	}
	if err := am.RegisterAccount(
		ctx, msg.Referrer, coin, msg.NewUser, msg.NewSigningPubKey, msg.NewTransactionPubKey); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle RegisterMsg
func handleUpdateAccountMsg(ctx sdk.Context, am AccountKeeper, msg types.UpdateAccountMsg) sdk.Result {
	if err := am.UpdateJSONMeta(ctx, msg.Username, msg.JSONMeta); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
