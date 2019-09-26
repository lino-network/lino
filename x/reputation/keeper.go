package reputation

//go:generate mockery -name ReputationKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

type ReputationKeeper interface {
	// upon donation, record this donation and return the impact factor.
	DonateAt(
		ctx sdk.Context,
		username types.AccountKey,
		post types.Permlink,
		amount types.MiniDollar) (types.MiniDollar, sdk.Error)

	// get user's latest reputation, which is the largest impact factor a user can
	// make in a window.
	GetReputation(ctx sdk.Context, username types.AccountKey) (types.MiniDollar, sdk.Error)

	// update game status on block end.
	Update(ctx sdk.Context) sdk.Error

	// return the current round start time
	GetCurrentRound(ctx sdk.Context) (int64, sdk.Error)

	GetReputationDetail(ctx sdk.Context, username types.AccountKey) (string, sdk.Error)

	// import/export this module to files
	ExportToFile(ctx sdk.Context, file string) error
	ImportFromFile(ctx sdk.Context, file string) error
}
