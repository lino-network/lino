package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
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
