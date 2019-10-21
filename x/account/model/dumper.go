package model

import (
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/testutils"
)

func NewAccountDumper(store AccountStorage) *testutils.Dumper {
	dumper := testutils.NewDumper(store.key, store.cdc)
	dumper.RegisterType(&AccountInfo{}, "lino/account/info", AccountInfoSubstore)
	dumper.RegisterType(&AccountBank{}, "lino/account/bank", AccountBankSubstore)
	dumper.RegisterType(&AccountMeta{}, "lino/account/meta", AccountMetaSubstore)
	dumper.RegisterType(&Pool{}, "lino/account/pool", AccountPoolSubstore)
	dumper.RegisterType(&Supply{}, "lino/account/supply", AccountSupplySubstore)
	return dumper
}
