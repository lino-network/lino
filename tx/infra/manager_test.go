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

	cases := map[string]struct {
		User1Usage             int64
		User2Usage             int64
		ExpectUser1UsageWeight sdk.Rat
		ExpectUser2UsageWeight sdk.Rat
	}{
		"test normal report": {
			25, 75, sdk.NewRat(1, 4), sdk.NewRat(3, 4),
		},
		"test empty report": {
			0, 0, sdk.NewRat(1, 2), sdk.NewRat(1, 2),
		},
		"issue https://github.com/lino-network/lino/issues/150": {
			3333333, 4444444, sdk.NewRat(429, 1000), sdk.NewRat(571, 1000),
		},
	}
	for testName, cs := range cases {
		im.ReportUsage(ctx, "user1", cs.User1Usage)
		im.ReportUsage(ctx, "user2", cs.User2Usage)

		w1, _ := im.GetUsageWeight(ctx, "user1")
		if !cs.ExpectUser1UsageWeight.Equal(w1) {
			t.Errorf(
				"%s: expect user1 usage weight %v, got %v",
				testName, cs.ExpectUser1UsageWeight, w1)
			return
		}

		w2, _ := im.GetUsageWeight(ctx, "user2")
		if !cs.ExpectUser2UsageWeight.Equal(w2) {
			t.Errorf(
				"%s: expect user2 usage weight %v, got %v",
				testName, cs.ExpectUser2UsageWeight, w2)
			return
		}
		im.ClearUsage(ctx)
	}
}
