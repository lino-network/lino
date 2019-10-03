package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

// FeedPriceMsg - Validastors need to send this message to feed price.
type FeedPriceMsg struct {
	Username types.AccountKey `json:"username"`
	Price    types.MiniDollar `json:"price"`
}

var _ types.Msg = FeedPriceMsg{}

// Route - implements sdk.Msg
func (msg FeedPriceMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg FeedPriceMsg) Type() string { return "FeedPriceMsg" }

// ValidateBasic - implements sdk.Msg
func (msg FeedPriceMsg) ValidateBasic() sdk.Error {
	if !msg.Username.IsValid() {
		return types.ErrInvalidUsername(msg.Username)
	}
	if !msg.Price.IsPositive() {
		return ErrInvalidPriceFeed(msg.Price)
	}
	return nil
}

func (msg FeedPriceMsg) String() string {
	return fmt.Sprintf("FeedPriceMsg{%s, %s}", msg.Username, msg.Price)
}

func (msg FeedPriceMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg FeedPriceMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg FeedPriceMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetConsumeAmount - implements types.Msg
func (msg FeedPriceMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// utils
func getSignBytes(msg sdk.Msg) []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
