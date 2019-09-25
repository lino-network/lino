package vote

import (
	"fmt"
	"reflect"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler - Handle all "vote" type messages.
func NewHandler(vk VoteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.StakeInMsg:
			return handleStakeInMsg(ctx, vk, msg)
		case types.StakeOutMsg:
			return handleStakeOutMsg(ctx, vk, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized vote msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleStakeInMsg(ctx sdk.Context, vk VoteKeeper, msg types.StakeInMsg) sdk.Result {
	coin, err := linotypes.LinoToCoin(msg.Deposit)
	if err != nil {
		return err.Result()
	}

	if err := vk.StakeIn(ctx, msg.Username, coin); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleStakeOutMsg(ctx sdk.Context, vk VoteKeeper, msg types.StakeOutMsg) sdk.Result {
	coin, err := linotypes.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if err := vk.StakeOut(ctx, msg.Username, coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
