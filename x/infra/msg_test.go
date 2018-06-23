package infra

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestProviderReportMsg(t *testing.T) {
	cases := []struct {
		providerReportMsg ProviderReportMsg
		expectError       sdk.Error
	}{
		{NewProviderReportMsg("user1", 100), nil},
		{NewProviderReportMsg("", 100), ErrInvalidUsername()},
		{NewProviderReportMsg("user1", -100), ErrInvalidUsage()},
	}

	for _, cs := range cases {
		result := cs.providerReportMsg.ValidateBasic()
		assert.Equal(t, result, cs.expectError)
	}
}

func TestMsgPermission(t *testing.T) {
	cases := map[string]struct {
		msg              types.Msg
		expectPermission types.Permission
	}{
		"provider report msg": {
			NewProviderReportMsg("test", 1),
			types.TransactionPermission},
	}

	for testName, cs := range cases {
		permission := cs.msg.GetPermission()
		if cs.expectPermission != permission {
			t.Errorf(
				"%s: expect permission incorrect, expect %v, got %v",
				testName, cs.expectPermission, permission)
			return
		}
	}
}
