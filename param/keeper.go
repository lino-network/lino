package param

//go:generate mockery -name ParamKeeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ParamKeeper interface {
	GetInfraInternalAllocationParam(
		ctx sdk.Context) (*InfraInternalAllocationParam, sdk.Error)
	GetPostParam(ctx sdk.Context) (*PostParam, sdk.Error)
	GetDeveloperParam(ctx sdk.Context) (*DeveloperParam, sdk.Error)
	GetVoteParam(ctx sdk.Context) (*VoteParam, sdk.Error)
	GetProposalParam(ctx sdk.Context) (*ProposalParam, sdk.Error)
	GetValidatorParam(ctx sdk.Context) (*ValidatorParam, sdk.Error)
	GetCoinDayParam(ctx sdk.Context) (*CoinDayParam, sdk.Error)
	GetBandwidthParam(ctx sdk.Context) (*BandwidthParam, sdk.Error)
	GetAccountParam(ctx sdk.Context) (*AccountParam, sdk.Error)
	GetReputationParam(ctx sdk.Context) (*ReputationParam, sdk.Error)
	UpdateGlobalGrowthRate(ctx sdk.Context, growthRate sdk.Dec) sdk.Error
}

var _ ParamKeeper = ParamHolder{}
