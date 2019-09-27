package price

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// linotypes "github.com/lino-network/lino/types"
	types "github.com/lino-network/lino/x/price/types"
)

type FeedPriceMsg = types.FeedPriceMsg

// NewHandler - Handle all "price" type messages.
func NewHandler(pm PriceKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case FeedPriceMsg:
			return handleFeedPriceMsg(ctx, msg, pm)
		default:
			errMsg := fmt.Sprintf("unknown price msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleFeedPriceMsg feed price message
func handleFeedPriceMsg(ctx sdk.Context, msg FeedPriceMsg, pm PriceKeeper) sdk.Result {
	err := pm.FeedPrice(ctx, msg.Username, msg.Price)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
