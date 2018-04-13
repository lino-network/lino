package infra

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestReportBasic(t *testing.T) {
	ctx, im := setupTest(t, 0)
	handler := NewHandler(*im)
	im.InitGenesis(ctx)

	user1 := types.AccountKey("user1")
	usage := int64(100)
	im.RegisterInfraProvider(ctx, user1)

	// infra provider not exist
	msg1 := NewProviderReportMsg("qwdqwdqw", usage)
	res := handler(ctx, msg1)
	assert.Equal(t, ErrProviderNotFound().Result(), res)

	msg2 := NewProviderReportMsg("user1", usage)
	res2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, res2)

	provider, _ := im.storage.GetInfraProvider(ctx, user1)
	assert.Equal(t, usage, provider.Usage)

}
