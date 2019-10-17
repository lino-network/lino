package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/testutils"
	types "github.com/lino-network/lino/types"
)

func NewDumper(store GlobalStorage, options ...testutils.OptionCodec) *testutils.Dumper {
	dumper := testutils.NewDumper(store.key, store.cdc, options...)
	dumper.RegisterType(&types.TimeEventList{}, "lino/global/eventlist", TimeEventListSubStore)
	dumper.RegisterType(&GlobalTime{}, "lino/global/time", TimeSubStore)
	dumper.RegisterType(&[]EventError{}, "lino/global/eventerr", EventErrorSubStore)
	dumper.RegisterType(&[]types.BCEventErr{}, "lino/global/bceventerr", BCErrorSubStore)
	return dumper
}

func DumpToFile(ctx sdk.Context, store GlobalStorage, filepath string) {
	dumper := NewDumper(store)
	dumper.DumpToFile(ctx, filepath)
}

func LoadFromFile(ctx sdk.Context, store GlobalStorage, filepath string) {
	dumper := NewDumper(store)
	dumper.LoadFromFile(ctx, filepath)
}
