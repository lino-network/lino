package developer

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestReportConsumption(t *testing.T) {
	ctx, am, dm, gm := setupTest(t, 0)
	handler := NewHandler(dm, am, gm)
	dm.InitGenesis(ctx)

	developer1 := createTestAccount(ctx, am, "developer1")
	am.AddCoin(ctx, developer1, c800000)
	msg := NewDeveloperRegisterMsg("developer1", l800000)
	handler(ctx, msg)

	developer2 := createTestAccount(ctx, am, "developer2")
	am.AddCoin(ctx, developer2, c800000)
	msg2 := NewDeveloperRegisterMsg("developer2", l800000)
	handler(ctx, msg2)

	con1 := types.NewCoin(100)
	dm.ReportConsumption(ctx, "developer1", con1)
	p1, _ := dm.GetConsumptionWeight(ctx, "developer1")
	assert.Equal(t, int64(1), p1.Evaluate())

	con2 := types.NewCoin(100)
	dm.ReportConsumption(ctx, "developer2", con2)
	p2, _ := dm.GetConsumptionWeight(ctx, "developer1")
	assert.Equal(t, true, p2.Equal(sdk.NewRat(1, 2)))

	dm.ClearConsumption(ctx)
	p3, _ := dm.GetConsumptionWeight(ctx, "developer1")
	assert.Equal(t, true, p3.IsZero())

}
