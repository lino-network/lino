package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/testutils"
)

func NewDeveloperDumper(store DeveloperStorage) *testutils.Dumper {
	dumper := testutils.NewDumper(store.key, store.cdc)
	dumper.RegisterType(&Developer{}, "lino/developer", developerSubstore)
	dumper.RegisterType(&AppIDA{}, "lino/appida", idaSubstore)
	dumper.RegisterType(&Role{}, "lino/role", userRoleSubstore)
	dumper.RegisterType(&IDABank{}, "lino/bank", idaBalanceSubstore)
	dumper.RegisterType(&ReservePool{}, "lino/reservepool", reservePoolSubstore)
	dumper.RegisterType(&AppIDAStats{}, "lino/appidastats", idaStatsSubstore)
	dumper.RegisterRawString(affiliatedAccSubstore)
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
