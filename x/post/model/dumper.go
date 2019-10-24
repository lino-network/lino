package model

import (
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/testutils"
	"github.com/lino-network/lino/types"
)

func NewPostDumper(store PostStorage) *testutils.Dumper {
	dumper := testutils.NewDumper(store.key, store.cdc)
	dumper.RegisterType(&Post{}, "lino/post", PostSubStore)
	dumper.RegisterType(&types.MiniDollar{}, "lino/minidollar", ConsumptionWindowSubStore)
	return dumper
}
