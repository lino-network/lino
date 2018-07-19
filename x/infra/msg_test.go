package infra

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestProviderReportMsg(t *testing.T) {
	testCases := []struct {
		testName          string
		providerReportMsg ProviderReportMsg
		expectError       sdk.Error
	}{
		{
			testName:          "normal case",
			providerReportMsg: NewProviderReportMsg("user1", 100),
			expectError:       nil,
		},
		{
			testName:          "invalid username",
			providerReportMsg: NewProviderReportMsg("", 100),
			expectError:       ErrInvalidUsername(),
		},
		{
			testName:          "invalid usage",
			providerReportMsg: NewProviderReportMsg("user1", -100),
			expectError:       ErrInvalidUsage(),
		},
	}

	for _, tc := range testCases {
		result := tc.providerReportMsg.ValidateBasic()
		if !assert.Equal(t, result, tc.expectError) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectError)
		}
	}
}

func TestMsgPermission(t *testing.T) {
	testCases := map[string]struct {
		msg              types.Msg
		expectPermission types.Permission
	}{
		"provider report msg": {
			msg:              NewProviderReportMsg("test", 1),
			expectPermission: types.TransactionPermission,
		},
	}

	for testName, tc := range testCases {
		permission := tc.msg.GetPermission()
		if tc.expectPermission != permission {
			t.Errorf("%s: diff permission,  got %v, want %v", testName, permission, tc.expectPermission)
			return
		}
	}
}
