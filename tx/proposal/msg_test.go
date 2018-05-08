package proposal

// import (
// 	"testing"
//
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/lino-network/lino/tx/proposal/model"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestCreateProposalMsg(t *testing.T) {
// 	des1 := model.ChangeParameterDescription{
// 		InfraAllocation:          sdk.Rat{20, 100},
// 		ContentCreatorAllocation: sdk.Rat{55, 100},
// 		DeveloperAllocation:      sdk.Rat{20, 100},
// 		ValidatorAllocation:      sdk.Rat{5, 100},
// 		CDNAllocation:            sdk.Rat{5, 100},
// 		StorageAllocation:        sdk.Rat{95, 100},
// 	}
//
// 	des2 := model.ChangeParameterDescription{
// 		InfraAllocation:          sdk.Rat{20, 100},
// 		ContentCreatorAllocation: sdk.Rat{55, 100},
// 		DeveloperAllocation:      sdk.Rat{25, 100},
// 		ValidatorAllocation:      sdk.Rat{5, 100},
// 		CDNAllocation:            sdk.Rat{5, 100},
// 		StorageAllocation:        sdk.Rat{95, 100},
// 	}
//
// 	des3 := model.ChangeParameterDescription{
// 		InfraAllocation:          sdk.Rat{20, 100},
// 		ContentCreatorAllocation: sdk.Rat{55, 100},
// 		DeveloperAllocation:      sdk.Rat{20, 100},
// 		ValidatorAllocation:      sdk.Rat{5, 100},
// 		CDNAllocation:            sdk.Rat{15, 100},
// 		StorageAllocation:        sdk.Rat{95, 100},
// 	}
// 	cases := []struct {
// 		createProposalMsg CreateProposalMsg
// 		expectError       sdk.Error
// 	}{
// 		{NewCreateProposalMsg("user1", des1), nil},
// 		{NewCreateProposalMsg("user1", des2), ErrIllegalParameter()},
// 		{NewCreateProposalMsg("user1", des3), ErrIllegalParameter()},
// 		{NewCreateProposalMsg("", des1), ErrInvalidUsername()},
// 	}
//
// 	for _, cs := range cases {
// 		result := cs.createProposalMsg.ValidateBasic()
// 		assert.Equal(t, result, cs.expectError)
// 	}
// }
