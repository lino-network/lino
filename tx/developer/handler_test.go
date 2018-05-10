package developer

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

var (
	l0      = types.LNO("0")
	l800000 = types.LNO("800000")

	c0      = types.Coin{0 * types.Decimals}
	c800000 = types.Coin{800000 * types.Decimals}
)

func TestRegistertBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	developer1 := createTestAccount(ctx, am, "developer1")
	am.AddSavingCoin(ctx, developer1, c800000)
	msg := NewDeveloperRegisterMsg("developer1", l800000)
	res := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, res)

	// check acc1's money has been withdrawn
	acc1Saving, _ := am.GetBankSaving(ctx, developer1)
	assert.Equal(t, acc1Saving, c0.Plus(initCoin))
	assert.Equal(t, true, dm.IsDeveloperExist(ctx, developer1))

	// check acc1 is in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 1, len(lst.AllDevelopers))
	assert.Equal(t, developer1, lst.AllDevelopers[0])

}

func TestRevokeBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	developer1 := createTestAccount(ctx, am, "developer1")
	am.AddSavingCoin(ctx, developer1, c800000)
	msg := NewDeveloperRegisterMsg("developer1", l800000)
	handler(ctx, msg)

	msg2 := NewDeveloperRevokeMsg("developer1")
	res2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, res2)
	// check acc1's depoist has not been added back
	acc1Saving, _ := am.GetBankSaving(ctx, developer1)
	assert.Equal(t, acc1Saving, c0.Plus(initCoin))
	assert.Equal(t, false, dm.IsDeveloperExist(ctx, developer1))

	// check acc1 is not in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 0, len(lst.AllDevelopers))
	assert.Equal(t, false, dm.IsDeveloperExist(ctx, "developer1"))
}
