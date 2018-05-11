package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestChangeGlobalAllocationParamMsg(t *testing.T) {
	des1 := param.GlobalAllocationParam{
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
	}

	des2 := param.GlobalAllocationParam{
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{25, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
	}

	cases := []struct {
		ChangeGlobalAllocationParamMsg ChangeGlobalAllocationParamMsg
		expectError                    sdk.Error
	}{
		{NewChangeGlobalAllocationParamMsg("user1", des1), nil},
		{NewChangeGlobalAllocationParamMsg("user1", des2), ErrIllegalParameter()},
		{NewChangeGlobalAllocationParamMsg("", des1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.ChangeGlobalAllocationParamMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              sdk.Msg
		expectPermission types.Permission
	}{
		"change evaluate of content value param": {
			NewChangeEvaluateOfContentValueParamMsg("creator",
				param.EvaluateOfContentValueParam{}),
			types.TransactionPermission},
		"change global allocation param": {
			NewChangeGlobalAllocationParamMsg("creator",
				param.GlobalAllocationParam{}),
			types.TransactionPermission},
		"change infra internal allocation param": {
			NewChangeInfraInternalAllocationParamMsg("creator",
				param.InfraInternalAllocationParam{}),
			types.TransactionPermission},
		"change vote param": {
			NewChangeVoteParamMsg("creator", param.VoteParam{}),
			types.TransactionPermission},
		"change proposal param": {
			NewChangeProposalParamMsg("creator", param.ProposalParam{}),
			types.TransactionPermission},
		"change developer param": {
			NewChangeDeveloperParamMsg("creator", param.DeveloperParam{}),
			types.TransactionPermission},
		"change validator param": {
			NewChangeValidatorParamMsg("creator", param.ValidatorParam{}),
			types.TransactionPermission},
		"change coinday param": {
			NewChangeCoinDayParamMsg("creator", param.CoinDayParam{}),
			types.TransactionPermission},
		"change account param": {
			NewChangeAccountParamMsg("creator", param.AccountParam{}),
			types.TransactionPermission},
	}

	for testName, cs := range cases {
		permissionLevel := cs.msg.Get(types.PermissionLevel)
		if permissionLevel == nil {
			if cs.expectPermission != types.PostPermission {
				t.Errorf(
					"%s: expect permission incorrect, expect %v, got %v",
					testName, cs.expectPermission, types.PostPermission)
				return
			} else {
				continue
			}
		}
		permission, ok := permissionLevel.(types.Permission)
		assert.Equal(t, ok, true)
		if cs.expectPermission != permission {
			t.Errorf(
				"%s: expect permission incorrect, expect %v, got %v",
				testName, cs.expectPermission, permission)
			return
		}
	}
}
