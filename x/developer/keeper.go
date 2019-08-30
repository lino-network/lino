package developer

//go:generate mockery -name DeveloperKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type DeveloperKeeper interface {
	MoveIDA(app types.AccountKey, from types.AccountKey, to types.AccountKey, amount types.MiniDollar) sdk.Error
	GetMiniIDAPrice(dev types.AccountKey) (types.MiniDollar, sdk.Error)
	DoesDeveloperExist(ctx sdk.Context, username types.AccountKey) bool
	ReportConsumption(
		ctx sdk.Context, username types.AccountKey, consumption types.Coin) sdk.Error
}

var _ DeveloperKeeper = DeveloperManager{}
