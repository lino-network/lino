package infra

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctx, im := setupTest(t, 0)
	im.InitGenesis(ctx)

	user1 := types.AccountKey("user1")
	im.RegisterInfraProvider(ctx, user1)

	_, err := im.storage.GetInfraProvider(ctx, user1)
	assert.Nil(t, err)

}

func TestInfraProviderList(t *testing.T) {
	ctx, im := setupTest(t, 0)
	im.InitGenesis(ctx)

	user1 := types.AccountKey("user1")
	im.RegisterInfraProvider(ctx, user1)

	addErr := im.AddToInfraProviderList(ctx, "user1")
	assert.Nil(t, addErr)

	lst, _ := im.storage.GetInfraProviderList(ctx)
	assert.Equal(t, 1, len(lst.AllInfraProviders))
	assert.Equal(t, user1, lst.AllInfraProviders[0])

	removeErr := im.RemoveFromProviderList(ctx, "user1")
	assert.Nil(t, removeErr)

	lst2, _ := im.storage.GetInfraProviderList(ctx)
	assert.Equal(t, 0, len(lst2.AllInfraProviders))

}

func TestReportUsage(t *testing.T) {
	ctx, im := setupTest(t, 0)
	im.InitGenesis(ctx)

	user1 := types.AccountKey("user1")
	im.RegisterInfraProvider(ctx, user1)

	user2 := types.AccountKey("user2")
	im.RegisterInfraProvider(ctx, user2)

	im.AddToInfraProviderList(ctx, "user1")
	im.AddToInfraProviderList(ctx, "user2")

	testCases := map[string]struct {
		user1Usage             int64
		user2Usage             int64
		expectUser1UsageWeight sdk.Rat
		expectUser2UsageWeight sdk.Rat
	}{
		"test normal report": {
			user1Usage:             25,
			user2Usage:             75,
			expectUser1UsageWeight: sdk.NewRat(1, 4),
			expectUser2UsageWeight: sdk.NewRat(3, 4),
		},
		"test empty report": {
			user1Usage:             0,
			user2Usage:             0,
			expectUser1UsageWeight: sdk.NewRat(1, 2),
			expectUser2UsageWeight: sdk.NewRat(1, 2),
		},
		"issue https://github.com/lino-network/lino/issues/150": {
			user1Usage:             3333333,
			user2Usage:             4444444,
			expectUser1UsageWeight: sdk.NewRat(429, 1000),
			expectUser2UsageWeight: sdk.NewRat(571, 1000),
		},
	}
	for testName, tc := range testCases {
		im.ReportUsage(ctx, "user1", tc.user1Usage)
		im.ReportUsage(ctx, "user2", tc.user2Usage)

		w1, _ := im.GetUsageWeight(ctx, "user1")
		if !tc.expectUser1UsageWeight.Equal(w1) {
			t.Errorf("%s: diff user1 usage weight, got %v, want %v", testName, w1, tc.expectUser1UsageWeight)
			return
		}

		w2, _ := im.GetUsageWeight(ctx, "user2")
		if !tc.expectUser2UsageWeight.Equal(w2) {
			t.Errorf("%s: diff user2 usage weight, got %v, want %v", testName, w2, tc.expectUser2UsageWeight)
			return
		}
		im.ClearUsage(ctx)
	}
}
