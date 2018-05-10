package param

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ChangeParamEvent struct {
	Param Parameter `json:"param"`
}

func (cpe ChangeParamEvent) Execute(ctx sdk.Context, ph ParamHolder) sdk.Error {
	parameter := cpe.Param
	switch parameter := parameter.(type) {
	case GlobalAllocationParam:
		ph.setGlobalAllocationParam(ctx, &parameter)
	case EvaluateOfContentValueParam:
		ph.setEvaluateOfContentValueParam(ctx, &parameter)
	case InfraInternalAllocationParam:
		ph.setInfraInternalAllocationParam(ctx, &parameter)
	case VoteParam:
		ph.setVoteParam(ctx, &parameter)
	case ProposalParam:
		ph.setProposalParam(ctx, &parameter)
	case DeveloperParam:
		ph.setDeveloperParam(ctx, &parameter)
	case ValidatorParam:
		ph.setValidatorParam(ctx, &parameter)
	case CoinDayParam:
		ph.setCoinDayParam(ctx, &parameter)
	case BandwidthParam:
		ph.setBandwidthParam(ctx, &parameter)
	default:
		return ErrInvalidaParameter()
	}
	return nil
}
