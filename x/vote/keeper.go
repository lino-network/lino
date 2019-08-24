package vote

//go:generate mockery -name VoteKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"

	"github.com/lino-network/lino/x/vote/types"
)

type VoteKeeper interface {
	DoesVoterExist(ctx sdk.Context, accKey linotypes.AccountKey) bool
	GetVoterDuty(ctx sdk.Context, accKey linotypes.AccountKey) types.VoterDuty
}

var _ VoteKeeper = VoteManager{}
