package developer

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestDeveloperRegisterMsg(t *testing.T) {
	cases := []struct {
		developerRegisterMsg DeveloperRegisterMsg
		expectError          sdk.Error
	}{
		{NewDeveloperRegisterMsg("user1", "10"), nil},
		{NewDeveloperRegisterMsg("", "10"), ErrInvalidUsername()},
		{NewDeveloperRegisterMsg("user1",
			"-1"), sdk.ErrInvalidCoins("LNO can't be less than lower bound")},
	}

	for _, cs := range cases {
		result := cs.developerRegisterMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestDeveloperRevokeMsg(t *testing.T) {
	cases := []struct {
		developerRevokeMsg DeveloperRevokeMsg
		expectError        sdk.Error
	}{
		{NewDeveloperRevokeMsg("user1"), nil},
		{NewDeveloperRevokeMsg(""), ErrInvalidUsername()},
	}

	for _, cs := range cases {
		result := cs.developerRevokeMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}
