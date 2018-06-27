package infra

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(im InfraManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case ProviderReportMsg:
			return handleProviderReportMsg(ctx, im, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized infra msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleProviderReportMsg(ctx sdk.Context, im InfraManager, msg ProviderReportMsg) sdk.Result {
	if !im.DoesInfraProviderExist(ctx, msg.Username) {
		return ErrProviderNotFound().Result()
	}

	if err := im.ReportUsage(ctx, msg.Username, msg.Usage); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
