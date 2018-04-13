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
