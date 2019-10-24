package global

//go:generate mockery -name GlobalKeeper

import (
	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global/model"
)

// GlobalKeeper - aka global event manager.
type GlobalKeeper interface {
	InitGenesis(ctx sdk.Context)

	// blockchain scheduled events.
	OnBeginBlock(ctx sdk.Context)
	OnEndBlock(ctx sdk.Context)

	// module events
	RegisterEventAtTime(ctx sdk.Context, unixTime int64, event linotypes.Event) sdk.Error
	ExecuteEvents(ctx sdk.Context, exec linotypes.EventExec)

	// Getter
	//// global time
	GetLastBlockTime(ctx sdk.Context) int64
	GetPastDay(ctx sdk.Context, unixTime int64) int64
	GetBCEventErrors(ctx sdk.Context) []linotypes.BCEventErr
	GetEventErrors(ctx sdk.Context) []model.EventError
	GetGlobalTime(ctx sdk.Context) model.GlobalTime

	// import export
	ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
	ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
}
