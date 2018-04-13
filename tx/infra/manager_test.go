package infra

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctx, im := setupTest(t, 0)
	im.InitGenesis(ctx)

	user1 := types.AccountKey("user1")
	im.RegisterInfraProvider(ctx, user1)

	_, getErr := im.storage.GetInfraProvider(ctx, user1)
	assert.Nil(t, getErr)

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
