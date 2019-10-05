package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/testutils"
)

func NewDeveloperDumper(store DeveloperStorage) *testutils.Dumper {
	dumper := testutils.NewDumper(store.key, store.cdc)
	dumper.RegisterType(&Developer{}, "lino/developer", DeveloperSubstore)
	dumper.RegisterType(&AppIDA{}, "lino/appida", IdaSubstore)
	dumper.RegisterType(&Role{}, "lino/role", UserRoleSubstore)
	dumper.RegisterType(&IDABank{}, "lino/bank", IdaBalanceSubstore)
	dumper.RegisterType(&ReservePool{}, "lino/reservepool", ReservePoolSubstore)
	dumper.RegisterType(&AppIDAStats{}, "lino/appidastats", IdaStatsSubstore)
	dumper.RegisterRawString(AffiliatedAccSubstore)
	return dumper
}

func DumpToFile(ctx sdk.Context, store DeveloperStorage, filepath string) {
	dumper := NewDeveloperDumper(store)
	dumper.DumpToFile(ctx, filepath)
}

func LoadFromFile(ctx sdk.Context, store DeveloperStorage, filepath string) {
	dumper := NewDeveloperDumper(store)
	dumper.LoadFromFile(ctx, filepath)
}
