package developer

//go:generate mockery -name DeveloperKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/developer/model"
)

type DeveloperKeeper interface {
	MoveIDA(app types.AccountKey, from types.AccountKey, to types.AccountKey, amount types.MiniDollar) sdk.Error
	GetMiniIDAPrice(dev types.AccountKey) (types.MiniDollar, sdk.Error)
	DoesDeveloperExist(ctx sdk.Context, username types.AccountKey) bool
	ReportConsumption(
		ctx sdk.Context, username types.AccountKey, consumption types.Coin) sdk.Error
	GetLiveDevelopers(ctx sdk.Context) []model.Developer
	GetAffiliatingApp(ctx sdk.Context, username types.AccountKey) (types.AccountKey, sdk.Error)
}

var _ DeveloperKeeper = DeveloperManager{}
