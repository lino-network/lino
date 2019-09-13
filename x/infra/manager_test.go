package infra

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctx, im := setupTest(t, 0)
	err := im.InitGenesis(ctx)
	if err != nil {
		panic(err)
	}

	user1 := types.AccountKey("user1")
	err = im.RegisterInfraProvider(ctx, user1)
	if err != nil {
		panic(err)
	}

	_, err = im.storage.GetInfraProvider(ctx, user1)
	assert.Nil(t, err)

}

func TestInfraProviderList(t *testing.T) {
	ctx, im := setupTest(t, 0)
	err := im.InitGenesis(ctx)
	if err != nil {
		panic(err)
	}

	user1 := types.AccountKey("user1")
	err = im.RegisterInfraProvider(ctx, user1)
	if err != nil {
		panic(err)
	}

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
	err := im.InitGenesis(ctx)
	if err != nil {
		panic(err)
	}

	user1 := types.AccountKey("user1")
	err = im.RegisterInfraProvider(ctx, user1)
	if err != nil {
		panic(err)
	}

	user2 := types.AccountKey("user2")
	err = im.RegisterInfraProvider(ctx, user2)
	if err != nil {
		panic(err)
	}

	err = im.AddToInfraProviderList(ctx, "user1")
	if err != nil {
		panic(err)
	}
	err = im.AddToInfraProviderList(ctx, "user2")
	if err != nil {
		panic(err)
	}

	testCases := map[string]struct {
		user1Usage             int64
		user2Usage             int64
		expectUser1UsageWeight sdk.Dec
		expectUser2UsageWeight sdk.Dec
	}{
		"test normal report": {
			user1Usage:             25,
			user2Usage:             75,
			expectUser1UsageWeight: types.NewDecFromRat(1, 4),
			expectUser2UsageWeight: types.NewDecFromRat(3, 4),
		},
		"test empty report": {
			user1Usage:             0,
			user2Usage:             0,
			expectUser1UsageWeight: types.NewDecFromRat(1, 2),
			expectUser2UsageWeight: types.NewDecFromRat(1, 2),
		},
		"issue https://github.com/lino-network/lino/issues/150": {
			user1Usage:             3333333,
			user2Usage:             4444444,
			expectUser1UsageWeight: types.NewDecFromRat(3333333, 7777777),
			expectUser2UsageWeight: types.NewDecFromRat(4444444, 7777777),
		},
	}
	for testName, tc := range testCases {
		err := im.ReportUsage(ctx, "user1", tc.user1Usage)
		if err != nil {
			panic(err)
		}
		err = im.ReportUsage(ctx, "user2", tc.user2Usage)
		if err != nil {
			panic(err)
		}

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
		err = im.ClearUsage(ctx)
		if err != nil {
			panic(err)
		}
	}
}
