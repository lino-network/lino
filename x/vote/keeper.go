package vote

//go:generate mockery -name VoteKeeper

import (
	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	votemn "github.com/lino-network/lino/x/vote/manager"
	"github.com/lino-network/lino/x/vote/model"
	"github.com/lino-network/lino/x/vote/types"
)

type VoteKeeper interface {
	InitGenesis(ctx sdk.Context)
	DoesVoterExist(ctx sdk.Context, username linotypes.AccountKey) bool
	StakeIn(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error
	StakeOut(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin) sdk.Error
	ClaimInterest(ctx sdk.Context, username linotypes.AccountKey) sdk.Error
	GetVoterDuty(ctx sdk.Context, username linotypes.AccountKey) (types.VoterDuty, sdk.Error)
	AssignDuty(
		ctx sdk.Context, username linotypes.AccountKey, duty types.VoterDuty, frozenAmount linotypes.Coin) sdk.Error
	// It's caller's duty to move coins from stake-in pool to the destination pool.
	SlashStake(ctx sdk.Context, username linotypes.AccountKey, amount linotypes.Coin, destPool linotypes.PoolName) (linotypes.Coin, sdk.Error)
	UnassignDuty(ctx sdk.Context, username linotypes.AccountKey, waitingPeriodSec int64) sdk.Error
	ExecUnassignDutyEvent(ctx sdk.Context, event types.UnassignDutyEvent) sdk.Error
	GetLinoStake(ctx sdk.Context, username linotypes.AccountKey) (linotypes.Coin, sdk.Error)
	StakeInFor(ctx sdk.Context, sender linotypes.AccountKey, receiver linotypes.AccountKey, amount linotypes.Coin) sdk.Error
	RecordFriction(ctx sdk.Context, friction linotypes.Coin) sdk.Error
	DailyAdvanceLinoStakeStats(ctx sdk.Context) sdk.Error

	// Getter
	GetVoter(ctx sdk.Context, username linotypes.AccountKey) (*model.Voter, sdk.Error)
	GetStakeStatsOfDay(ctx sdk.Context, day int64) (model.LinoStakeStat, sdk.Error)

	// import export
	ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
	ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
}

var _ VoteKeeper = votemn.VoteManager{}
