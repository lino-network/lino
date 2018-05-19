package developer

import (
	"strconv"
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestRegistertBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	minBalance := types.NewCoin(1 * types.Decimals)
	developer1 := createTestAccount(ctx, am, "developer1", devParam.DeveloperMinDeposit.Plus(minBalance))
	deposit := strconv.FormatInt(devParam.DeveloperMinDeposit.ToInt64()/types.Decimals, 10)
	msg := NewDeveloperRegisterMsg("developer1", deposit)
	res := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, res)

	// check acc1's money has been withdrawn
	acc1Saving, _ := am.GetSavingFromBank(ctx, developer1)
	assert.Equal(t, true, acc1Saving.IsEqual(minBalance))
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

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	minBalance := types.NewCoin(1 * types.Decimals)
	developer1 := createTestAccount(ctx, am, "developer1", devParam.DeveloperMinDeposit.Plus(minBalance))
	deposit := strconv.FormatInt(devParam.DeveloperMinDeposit.ToInt64()/types.Decimals, 10)
	msg := NewDeveloperRegisterMsg("developer1", deposit)
	handler(ctx, msg)

	msg2 := NewDeveloperRevokeMsg("developer1")
	res2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, res2)
	// check acc1's depoist has not been added back
	acc1Saving, _ := am.GetSavingFromBank(ctx, developer1)
	assert.Equal(t, true, acc1Saving.IsEqual(minBalance))
	assert.Equal(t, false, dm.IsDeveloperExist(ctx, developer1))

	// check acc1 is not in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 0, len(lst.AllDevelopers))
	assert.Equal(t, false, dm.IsDeveloperExist(ctx, "developer1"))
}
