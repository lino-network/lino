package developer

//go:generate mockery -name DeveloperKeeper

import (
	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/developer/model"
)

type DeveloperKeeper interface {
	// developer
	DoesDeveloperExist(ctx sdk.Context, username linotypes.AccountKey) bool

	RegisterDeveloper(ctx sdk.Context, username linotypes.AccountKey, website, description, appMetaData string) sdk.Error
	UpdateDeveloper(
		ctx sdk.Context, username linotypes.AccountKey, website, description, appMetadata string) sdk.Error
	GetDeveloper(ctx sdk.Context, username linotypes.AccountKey) (model.Developer, sdk.Error)
	GetLiveDevelopers(ctx sdk.Context) []model.Developer

	// affiliated account
	UpdateAffiliated(ctx sdk.Context, appname, username linotypes.AccountKey, activate bool) sdk.Error
	GetAffiliatingApp(ctx sdk.Context, username linotypes.AccountKey) (linotypes.AccountKey, sdk.Error)
	GetAffiliated(ctx sdk.Context, app linotypes.AccountKey) []linotypes.AccountKey

	// IDA
	IssueIDA(ctx sdk.Context, appname linotypes.AccountKey, idaName string, idaPrice int64) sdk.Error
	MintIDA(ctx sdk.Context, appname linotypes.AccountKey, amount linotypes.Coin) sdk.Error
	GetMiniIDAPrice(ctx sdk.Context, app linotypes.AccountKey) (linotypes.MiniDollar, sdk.Error)
	AppTransferIDA(ctx sdk.Context, appname, signer linotypes.AccountKey, amount linotypes.MiniIDA, from, to linotypes.AccountKey) sdk.Error
	MoveIDA(ctx sdk.Context, app linotypes.AccountKey, from, to linotypes.AccountKey, amount linotypes.MiniDollar) sdk.Error
	BurnIDA(ctx sdk.Context, app, user linotypes.AccountKey, amount linotypes.MiniDollar) (linotypes.Coin, sdk.Error)
	UpdateIDAAuth(ctx sdk.Context, app, username linotypes.AccountKey, active bool) sdk.Error
	GetIDABank(ctx sdk.Context, app, user linotypes.AccountKey) (model.IDABank, sdk.Error)
	GetIDA(ctx sdk.Context, app linotypes.AccountKey) (model.AppIDA, sdk.Error)
	GetReservePool(ctx sdk.Context) model.ReservePool
	GetIDAStats(ctx sdk.Context, app linotypes.AccountKey) (model.AppIDAStats, sdk.Error)

	// consumption stats
	ReportConsumption(
		ctx sdk.Context, username linotypes.AccountKey, consumption linotypes.MiniDollar) sdk.Error
	MonthlyDistributeDevInflation(ctx sdk.Context) sdk.Error

	// Genesis
	InitGenesis(ctx sdk.Context, reservePoolAmount linotypes.Coin) sdk.Error

	// importer exporter
	ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
	ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
}

// var _ DeveloperKeeper = DeveloperManager{}
