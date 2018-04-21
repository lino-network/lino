package infra

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
