package developer

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

var (
	l0      = types.LNO(sdk.NewRat(0))
	l10     = types.LNO(sdk.NewRat(10))
	l11     = types.LNO(sdk.NewRat(11))
	l20     = types.LNO(sdk.NewRat(20))
	l21     = types.LNO(sdk.NewRat(21))
	l40     = types.LNO(sdk.NewRat(40))
	l100    = types.LNO(sdk.NewRat(100))
	l200    = types.LNO(sdk.NewRat(200))
	l400    = types.LNO(sdk.NewRat(400))
	l1000   = types.LNO(sdk.NewRat(1000))
	l1011   = types.LNO(sdk.NewRat(1011))
	l1021   = types.LNO(sdk.NewRat(1021))
	l1022   = types.LNO(sdk.NewRat(1022))
	l1100   = types.LNO(sdk.NewRat(1100))
	l1500   = types.LNO(sdk.NewRat(1500))
	l1600   = types.LNO(sdk.NewRat(1600))
	l1800   = types.LNO(sdk.NewRat(1800))
	l1900   = types.LNO(sdk.NewRat(1900))
	l2000   = types.LNO(sdk.NewRat(2000))
	l800000 = types.LNO(sdk.NewRat(800000))

	c0      = types.Coin{0 * types.Decimals}
	c10     = types.Coin{10 * types.Decimals}
	c11     = types.Coin{11 * types.Decimals}
	c20     = types.Coin{20 * types.Decimals}
	c21     = types.Coin{21 * types.Decimals}
	c100    = types.Coin{100 * types.Decimals}
	c200    = types.Coin{200 * types.Decimals}
	c400    = types.Coin{400 * types.Decimals}
	c600    = types.Coin{600 * types.Decimals}
	c900    = types.Coin{900 * types.Decimals}
	c1000   = types.Coin{1000 * types.Decimals}
	c1011   = types.Coin{1011 * types.Decimals}
	c1021   = types.Coin{1021 * types.Decimals}
	c1022   = types.Coin{1022 * types.Decimals}
	c1500   = types.Coin{1500 * types.Decimals}
	c1600   = types.Coin{1600 * types.Decimals}
	c1800   = types.Coin{1800 * types.Decimals}
	c1900   = types.Coin{1900 * types.Decimals}
	c2000   = types.Coin{2000 * types.Decimals}
	c800000 = types.Coin{800000 * types.Decimals}
)

func TestRegistertBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(*dm, *am, *gm)
	dm.InitGenesis(ctx)

	developer1 := createTestAccount(ctx, am, "developer1")
	am.AddCoin(ctx, developer1, c800000)
	msg := NewDeveloperRegisterMsg("developer1", l800000)
	res := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, res)

	// check acc1's money has been withdrawn
	acc1Balance, _ := am.GetBankBalance(ctx, developer1)
	assert.Equal(t, acc1Balance, c0.Plus(initCoin))
	assert.Equal(t, true, dm.IsDeveloperExist(ctx, developer1))

	// check acc1 is in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 1, len(lst.AllDevelopers))
	assert.Equal(t, developer1, lst.AllDevelopers[0])

}

func TestRevokeBasic(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(*dm, *am, *gm)
	dm.InitGenesis(ctx)

	developer1 := createTestAccount(ctx, am, "developer1")
	am.AddCoin(ctx, developer1, c800000)
	msg := NewDeveloperRegisterMsg("developer1", l800000)
	handler(ctx, msg)

	msg2 := NewDeveloperRevokeMsg("developer1")
	res2 := handler(ctx, msg2)
	assert.Equal(t, sdk.Result{}, res2)
	// check acc1's depoist has not been added back
	acc1Balance, _ := am.GetBankBalance(ctx, developer1)
	assert.Equal(t, acc1Balance, c0.Plus(initCoin))
	assert.Equal(t, false, dm.IsDeveloperExist(ctx, developer1))

	// check acc1 is not in the developer list
	lst, _ := dm.GetDeveloperList(ctx)
	assert.Equal(t, 0, len(lst.AllDevelopers))
}
