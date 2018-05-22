package developer

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestReportConsumption(t *testing.T) {
	ctx, _, dm, _ := setupTest(t, 0)
	dm.InitGenesis(ctx)

	devParam, _ := dm.paramHolder.GetDeveloperParam(ctx)
	dm.RegisterDeveloper(ctx, "developer1", devParam.DeveloperMinDeposit)
	dm.RegisterDeveloper(ctx, "developer2", devParam.DeveloperMinDeposit)

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
	assert.Equal(t, true, p3.Equal(sdk.NewRat(1, 2)))

}
