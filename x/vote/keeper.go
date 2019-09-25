package vote

//go:generate mockery -name VoteKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"

	votemn "github.com/lino-network/lino/x/vote/manager"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/lino-network/lino/x/vote/types"
)

type VoteKeeper interface {
	DoesVoterExist(ctx sdk.Context, username linotypes.AccountKey) bool
	StakeIn(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error
	AddStake(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error
	StakeOut(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error
	MinusStake(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error
	ClaimInterest(ctx sdk.Context, username linotypes.AccountKey) sdk.Error
	GetVoterDuty(ctx sdk.Context, username linotypes.AccountKey) (types.VoterDuty, sdk.Error)
	AssignDuty(
		ctx sdk.Context, username linotypes.AccountKey, duty types.VoterDuty, frozenAmount linotypes.Coin) sdk.Error
	SlashStake(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error
	UnassignDuty(ctx sdk.Context, username linotypes.AccountKey, waitingPeriodSec int64) sdk.Error
	ExecUnassignDutyEvent(ctx sdk.Context, event types.UnassignDutyEvent) sdk.Error
	GetLinoStake(ctx sdk.Context, username linotypes.AccountKey) (linotypes.Coin, sdk.Error)

	// Getter
	GetVoter(ctx sdk.Context, username linotypes.AccountKey) (*model.Voter, sdk.Error)
}

var _ VoteKeeper = votemn.VoteManager{}
