package param

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ErrInvalidaParameter - error when propose a proposal with invalid parameter.
func ErrInvalidaParameter() sdk.Error {
	return types.NewError(types.CodeInvalidaParameter, fmt.Sprintf("invalida parameter"))
}

// ErrParamHolderGenesisFailed - error when genesis failed.
func ErrParamHolderGenesisFailed() sdk.Error {
	return types.NewError(types.CodeParamHolderGenesisError, fmt.Sprintf("param holder genesis failed"))
}

// ErrDeveloperParamNotFound - error when developer param is empty.
func ErrDeveloperParamNotFound() sdk.Error {
	return types.NewError(types.CodeDeveloperParamNotFound, fmt.Sprintf("developer param not found"))
}

// ErrValidatorParamNotFound - error when validator param is empty.
func ErrValidatorParamNotFound() sdk.Error {
	return types.NewError(types.CodeValidatorParamNotFound, fmt.Sprintf("validator param not found"))
}

// ErrCoinDayParamNotFound - error when coin day param is empty.
func ErrCoinDayParamNotFound() sdk.Error {
	return types.NewError(types.CodeCoinDayParamNotFound, fmt.Sprintf("coin day param not found"))
}

// ErrBandwidthParamNotFound - error when bandwidth param is empty.
func ErrBandwidthParamNotFound() sdk.Error {
	return types.NewError(types.CodeBandwidthParamNotFound, fmt.Sprintf("bandwidth param not found"))
}

// ErrAccountParamNotFound - error when account param is empty.
func ErrAccountParamNotFound() sdk.Error {
	return types.NewError(types.CodeAccountParamNotFound, fmt.Sprintf("account param not found"))
}

// ErrVoteParamNotFound - error when vote param is empty.
func ErrVoteParamNotFound() sdk.Error {
	return types.NewError(types.CodeVoteParamNotFound, fmt.Sprintf("vote param not found"))
}

// ErrProposalParamNotFound - error when proposal param is empty.
func ErrProposalParamNotFound() sdk.Error {
	return types.NewError(types.CodeProposalParamNotFound, fmt.Sprintf("proposal param not found"))
}

// ErrGlobalAllocationParamNotFound - error when global allocation param is empty.
func ErrGlobalAllocationParamNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalAllocationParamNotFound, fmt.Sprintf("global allocation param not found"))
}

// ErrInfraAllocationParamNotFound - error when infra allocation param is empty.
func ErrInfraAllocationParamNotFound() sdk.Error {
	return types.NewError(types.CodeInfraAllocationParamNotFound, fmt.Sprintf("infra internal allocation param not found"))
}

// ErrPostParamNotFound - error when post param is empty.
func ErrPostParamNotFound() sdk.Error {
	return types.NewError(types.CodePostParamNotFound, fmt.Sprintf("post param not found"))
}

// ErrEvaluateOfContentValueParamNotFound - error when evaluate of content value param is empty.
func ErrEvaluateOfContentValueParamNotFound() sdk.Error {
	return types.NewError(types.CodeEvaluateOfContentValueParamNotFound, fmt.Sprintf("evaluate of content value param not found"))
}

// ErrReputationParamNotFound - error when reputation param is empty.
func ErrReputationParamNotFound() sdk.Error {
	return types.NewError(types.CodeReputationParamNotFound, fmt.Sprintf("reputation param not found"))
}

// ErrFailedToUnmarshalGlobalAllocationParam - error when unmarshal global allocation param failed.
func ErrFailedToUnmarshalGlobalAllocationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalGlobalAllocationParam, fmt.Sprintf("failed to unmarshal global allocation param: %s", err.Error()))
}

// ErrFailedToUnmarshalPostParam - error when unmarshal post param failed.
func ErrFailedToUnmarshalPostParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostParam, fmt.Sprintf("failed to unmarshal post param: %s", err.Error()))
}

// ErrFailedToUnmarshalValidatorParam - error when unmarshal validator param failed.
func ErrFailedToUnmarshalValidatorParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalValidatorParam, fmt.Sprintf("failed to unmarshal validator param: %s", err.Error()))
}

// ErrFailedToUnmarshalEvaluateOfContentValueParam - error when unmarshal evaluate of content value param failed.
func ErrFailedToUnmarshalEvaluateOfContentValueParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalEvaluateOfContentValueParam, fmt.Sprintf("failed to unmarshal evaluate of content value param: %s", err.Error()))
}

// ErrFailedToUnmarshalInfraInternalAllocationParam - error when unmarshal infra internal allocation param failed.
func ErrFailedToUnmarshalInfraInternalAllocationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalInfraInternalAllocationParam, fmt.Sprintf("failed to unmarshal infra internal allocation param: %s", err.Error()))
}

// ErrFailedToUnmarshalDeveloperParam - error when unmarshal developer param failed.
func ErrFailedToUnmarshalDeveloperParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalDeveloperParam, fmt.Sprintf("failed to unmarshal developer param: %s", err.Error()))
}

// ErrFailedToUnmarshalVoteParam - error when unmarshal vote param failed.
func ErrFailedToUnmarshalVoteParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalVoteParam, fmt.Sprintf("failed to unmarshal vote param: %s", err.Error()))
}

// ErrFailedToUnmarshalProposalParam - error when unmarshal proposal param failed.
func ErrFailedToUnmarshalProposalParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalProposalParam, fmt.Sprintf("failed to unmarshal proposal param: %s", err.Error()))
}

// ErrFailedToUnmarshalCoinDayParam - error when unmarshal coin day param failed.
func ErrFailedToUnmarshalCoinDayParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalCoinDayParam, fmt.Sprintf("failed to unmarshal coin day param: %s", err.Error()))
}

// ErrFailedToUnmarshalBandwidthParam - error when unmarshal bandwidth param failed.
func ErrFailedToUnmarshalBandwidthParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalBandwidthParam, fmt.Sprintf("failed to unmarshal bandwidth param: %s", err.Error()))
}

// ErrFailedToUnmarshalAccountParam - error when unmarshal account param failed.
func ErrFailedToUnmarshalAccountParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalAccountParam, fmt.Sprintf("failed to unmarshal account param: %s", err.Error()))
}

// ErrFailedToUnmarshalReputationParam - error when unmarshal reputation param failed.
func ErrFailedToUnmarshalReputationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalReputationParam, fmt.Sprintf("failed to unmarshal account param: %s", err.Error()))
}

// ErrFailedToUnmarshalAccountParam - error when marshal global allocation param failed.
func ErrFailedToMarshalGlobalAllocationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalGlobalAllocationParam, fmt.Sprintf("failed to marshal global allocation param: %s", err.Error()))
}

// ErrFailedToMarshalPostParam - error when marshal post param failed.
func ErrFailedToMarshalPostParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostParam, fmt.Sprintf("failed to marshal post param: %s", err.Error()))
}

// ErrFailedToMarshalValidatorParam - error when marshal validator param failed.
func ErrFailedToMarshalValidatorParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalValidatorParam, fmt.Sprintf("failed to marshal validator param: %s", err.Error()))
}

// ErrFailedToMarshalEvaluateOfContentValueParam - error when marshal evaluate of content value param failed.
func ErrFailedToMarshalEvaluateOfContentValueParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalEvaluateOfContentValueParam, fmt.Sprintf("failed to marshal evaluate of content value param: %s", err.Error()))
}

// ErrFailedToMarshalInfraInternalAllocationParam - error when marshal infra internal allocation param failed.
func ErrFailedToMarshalInfraInternalAllocationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalInfraInternalAllocationParam, fmt.Sprintf("failed to marshal infra internal allocation param: %s", err.Error()))
}

// ErrFailedToMarshalDeveloperParam - error when marshal developer param failed.
func ErrFailedToMarshalDeveloperParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalDeveloperParam, fmt.Sprintf("failed to marshal developer param: %s", err.Error()))
}

// ErrFailedToMarshalVoteParam - error when marshal vote param failed.
func ErrFailedToMarshalVoteParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalVoteParam, fmt.Sprintf("failed to marshal vote param: %s", err.Error()))
}

// ErrFailedToMarshalProposalParam - error when marshal proposal param failed.
func ErrFailedToMarshalProposalParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalProposalParam, fmt.Sprintf("failed to marshal proposal param: %s", err.Error()))
}

// ErrFailedToMarshalCoinDayParam - error when marshal coin day param failed.
func ErrFailedToMarshalCoinDayParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalCoinDayParam, fmt.Sprintf("failed to marshal coin day param: %s", err.Error()))
}

// ErrFailedToMarshalBandwidthParam - error when marshal bandwidth day param failed.
func ErrFailedToMarshalBandwidthParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalBandwidthParam, fmt.Sprintf("failed to marshal bandwidth param: %s", err.Error()))
}

// ErrFailedToMarshalAccountParam - error when marshal account param failed.
func ErrFailedToMarshalAccountParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalAccountParam, fmt.Sprintf("failed to marshal account param: %s", err.Error()))
}

// ErrFailedToMarshalReputationParam - error when marshal reputation failed.
func ErrFailedToMarshalReputationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalReputationParam, fmt.Sprintf("failed to marshal reputation param: %s", err.Error()))
}

// ErrQueryFailed - error when query paramter store failed
func ErrQueryFailed() sdk.Error {
	return types.NewError(types.CodeParamQueryFailed, fmt.Sprintf("query paramter store failed"))
}
