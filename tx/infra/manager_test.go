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

	im.ReportUsage(ctx, "user1", int64(25))
	im.ReportUsage(ctx, "user2", int64(75))

	w1, _ := im.GetUsageWeight(ctx, "user1")
	assert.Equal(t, true, sdk.NewRat(1, 4).Equal(w1))

	im.ClearUsage(ctx)
	w2, _ := im.GetUsageWeight(ctx, "user1")
	assert.Equal(t, true, sdk.NewRat(1, 2).Equal(w2))
}
