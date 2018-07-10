package param

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ErrInvalidaParameter() sdk.Error {
	return types.NewError(types.CodeInvalidaParameter, fmt.Sprintf("invalida parameter"))
}

func ErrParamHolderGenesisFailed() sdk.Error {
	return types.NewError(types.CodeParamHolderGenesisError, fmt.Sprintf("param holder genesis failed"))
}

func ErrDeveloperParamNotFound() sdk.Error {
	return types.NewError(types.CodeDeveloperParamNotFound, fmt.Sprintf("developer param not found"))
}

func ErrValidatorParamNotFound() sdk.Error {
	return types.NewError(types.CodeValidatorParamNotFound, fmt.Sprintf("validator param not found"))
}

func ErrCoinDayParamNotFound() sdk.Error {
	return types.NewError(types.CodeCoinDayParamNotFound, fmt.Sprintf("coin day param not found"))
}

func ErrBandwidthParamNotFound() sdk.Error {
	return types.NewError(types.CodeBandwidthParamNotFound, fmt.Sprintf("bandwidth param not found"))
}

func ErrAccountParamNotFound() sdk.Error {
	return types.NewError(types.CodeAccountParamNotFound, fmt.Sprintf("account param not found"))
}

func ErrVoteParamNotFound() sdk.Error {
	return types.NewError(types.CodeVoteParamNotFound, fmt.Sprintf("vote param not found"))
}

func ErrProposalParamNotFound() sdk.Error {
	return types.NewError(types.CodeProposalParamNotFound, fmt.Sprintf("proposal param not found"))
}

func ErrGlobalAllocationParamNotFound() sdk.Error {
	return types.NewError(types.CodeGlobalAllocationParamNotFound, fmt.Sprintf("global allocation param not found"))
}

func ErrInfraAllocationParamNotFound() sdk.Error {
	return types.NewError(types.CodeInfraAllocationParamNotFound, fmt.Sprintf("infra internal allocation param not found"))
}

func ErrPostParamNotFound() sdk.Error {
	return types.NewError(types.CodePostParamNotFound, fmt.Sprintf("post param not found"))
}

func ErrEvaluateOfContentValueParamNotFound() sdk.Error {
	return types.NewError(types.CodeEvaluateOfContentValueParamNotFound, fmt.Sprintf("evaluate of content value param not found"))
}

// unmarshal error
func ErrFailedToUnmarshalGlobalAllocationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalGlobalAllocationParam, fmt.Sprintf("failed to unmarshal global allocation param: %s", err.Error()))
}

func ErrFailedToUnmarshalPostParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalPostParam, fmt.Sprintf("failed to unmarshal post param: %s", err.Error()))
}

func ErrFailedToUnmarshalValidatorParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalValidatorParam, fmt.Sprintf("failed to unmarshal validator param: %s", err.Error()))
}

func ErrFailedToUnmarshalEvaluateOfContentValueParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalEvaluateOfContentValueParam, fmt.Sprintf("failed to unmarshal evaluate of content value param: %s", err.Error()))
}

func ErrFailedToUnmarshalInfraInternalAllocationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalInfraInternalAllocationParam, fmt.Sprintf("failed to unmarshal infra internal allocation param: %s", err.Error()))
}

func ErrFailedToUnmarshalDeveloperParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalDeveloperParam, fmt.Sprintf("failed to unmarshal developer param: %s", err.Error()))
}

func ErrFailedToUnmarshalVoteParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalVoteParam, fmt.Sprintf("failed to unmarshal vote param: %s", err.Error()))
}

func ErrFailedToUnmarshalProposalParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalProposalParam, fmt.Sprintf("failed to unmarshal proposal param: %s", err.Error()))
}

func ErrFailedToUnmarshalCoinDayParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalCoinDayParam, fmt.Sprintf("failed to unmarshal coin day param: %s", err.Error()))
}

func ErrFailedToUnmarshalBandwidthParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalBandwidthParam, fmt.Sprintf("failed to unmarshal bandwidth param: %s", err.Error()))
}

func ErrFailedToUnmarshalAccountParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalAccountParam, fmt.Sprintf("failed to unmarshal account param: %s", err.Error()))
}

// marshal error
func ErrFailedToMarshalGlobalAllocationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalGlobalAllocationParam, fmt.Sprintf("failed to marshal global allocation param: %s", err.Error()))
}

func ErrFailedToMarshalPostParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalPostParam, fmt.Sprintf("failed to marshal post param: %s", err.Error()))
}

func ErrFailedToMarshalValidatorParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalValidatorParam, fmt.Sprintf("failed to marshal validator param: %s", err.Error()))
}

func ErrFailedToMarshalEvaluateOfContentValueParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalEvaluateOfContentValueParam, fmt.Sprintf("failed to marshal evaluate of content value param: %s", err.Error()))
}

func ErrFailedToMarshalInfraInternalAllocationParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalInfraInternalAllocationParam, fmt.Sprintf("failed to marshal infra internal allocation param: %s", err.Error()))
}

func ErrFailedToMarshalDeveloperParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalDeveloperParam, fmt.Sprintf("failed to marshal developer param: %s", err.Error()))
}

func ErrFailedToMarshalVoteParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalVoteParam, fmt.Sprintf("failed to marshal vote param: %s", err.Error()))
}

func ErrFailedToMarshalProposalParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalProposalParam, fmt.Sprintf("failed to marshal proposal param: %s", err.Error()))
}

func ErrFailedToMarshalCoinDayParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalCoinDayParam, fmt.Sprintf("failed to marshal coin day param: %s", err.Error()))
}

func ErrFailedToMarshalBandwidthParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalBandwidthParam, fmt.Sprintf("failed to marshal bandwidth param: %s", err.Error()))
}

func ErrFailedToMarshalAccountParam(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalAccountParam, fmt.Sprintf("failed to marshal account param: %s", err.Error()))
}
