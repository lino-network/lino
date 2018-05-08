package proposal

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/stretchr/testify/assert"
)

func TestChangeGlobalAllocationMsg(t *testing.T) {
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
		changeGlobalAllocationMsg ChangeGlobalAllocationMsg
		expectError               sdk.Error
	}{
		{NewChangeGlobalAllocationMsg("user1", des1), nil},
		{NewChangeGlobalAllocationMsg("user1", des2), ErrIllegalParameter()},
		{NewChangeGlobalAllocationMsg("", des1), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.changeGlobalAllocationMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}
