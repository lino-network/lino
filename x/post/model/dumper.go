package model

import (
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/testutils"
)

func NewPostDumper(store PostStorage) *testutils.Dumper {
	dumper := testutils.NewDumper(store.key, store.cdc)
	dumper.RegisterType(&Post{}, "lino/post", PostSubStore)
	return dumper
}
