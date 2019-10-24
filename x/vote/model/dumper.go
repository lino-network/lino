package model

import (
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/testutils"
)

func NewVoteDumper(store VoteStorage) *testutils.Dumper {
	dumper := testutils.NewDumper(store.key, store.cdc)
	dumper.RegisterType(&Voter{}, "lino/voter", VoterSubstore)
	dumper.RegisterType(&LinoStakeStat{}, "lino/stakestats", LinoStakeStatSubStore)
	return dumper
}
