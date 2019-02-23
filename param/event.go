package param

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ChangeParamEvent - change parameter event
type ChangeParamEvent struct {
	Param Parameter `json:"param"`
}

// Execute - execute change parameter event
func (cpe ChangeParamEvent) Execute(ctx sdk.Context, ph ParamHolder) sdk.Error {
	parameter := cpe.Param
	switch parameter := parameter.(type) {
	case GlobalAllocationParam:
		return ph.setGlobalAllocationParam(ctx, &parameter)
	case InfraInternalAllocationParam:
		return ph.setInfraInternalAllocationParam(ctx, &parameter)
	case VoteParam:
		return ph.setVoteParam(ctx, &parameter)
	case ProposalParam:
		return ph.setProposalParam(ctx, &parameter)
	case DeveloperParam:
		return ph.setDeveloperParam(ctx, &parameter)
	case ValidatorParam:
		return ph.setValidatorParam(ctx, &parameter)
	case BandwidthParam:
		return ph.setBandwidthParam(ctx, &parameter)
	case AccountParam:
		return ph.setAccountParam(ctx, &parameter)
	case PostParam:
		return ph.setPostParam(ctx, &parameter)
	default:
		return ErrInvalidaParameter()
	}
}
