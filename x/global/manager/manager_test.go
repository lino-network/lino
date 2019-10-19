package manager

import (
	"testing"
	"time"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	mparam "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
	mapp "github.com/lino-network/lino/x/global/manager/mocks"
	"github.com/lino-network/lino/x/global/model"
)

type testEvent struct {
	Id int64 `json:"id"`
}

func regTestCodec(c *wire.Codec) {
	c.RegisterInterface((*linotypes.Event)(nil), nil)
	c.RegisterConcrete(testEvent{}, "lino/testevent", nil)
}

func testCodec() *wire.Codec {
	c := wire.New()
	regTestCodec(c)
	return c
}

var (
	storeKeyStr = "testGlobalStore"
	kvStoreKey  = sdk.NewKVStoreKey(storeKeyStr)
)

type storeDumper struct{}

func (dumper storeDumper) NewDumper() *testutils.Dumper {
	return model.NewDumper(model.NewGlobalStorage(kvStoreKey, testCodec()), regTestCodec)
}

type globalManagerTestSuite struct {
	testsuites.GoldenTestSuite
	global       GlobalManager
	mParamKeeper *mparam.ParamKeeper
	mApp         *mapp.FakeApp
}

func NewGlobalManagerTestSuite() *globalManagerTestSuite {
	return &globalManagerTestSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(storeDumper{}, kvStoreKey),
	}
}

func (suite *globalManagerTestSuite) SetupTest() {
	suite.mParamKeeper = new(mparam.ParamKeeper)
	suite.mApp = new(mapp.FakeApp)
	suite.SetupCtx(0, time.Unix(0, 0), kvStoreKey)
	suite.global = NewGlobalManager(kvStoreKey, suite.mParamKeeper, testCodec(),
		suite.mApp.Hourly,
		suite.mApp.Daily,
		suite.mApp.Monthly,
		suite.mApp.Yearly)
}

func TestGlobalManagerSuite(t *testing.T) {
	suite.Run(t, NewGlobalManagerTestSuite())
}

func (suite *globalManagerTestSuite) TestInitGenesis() {
	init := int64(123456)
	suite.NextBlock(time.Unix(init, 0))
	suite.global.InitGenesis(suite.Ctx)
	suite.Equal(init, suite.global.GetLastBlockTime(suite.Ctx))
	suite.Golden()
}

func (suite *globalManagerTestSuite) TestGetPastDay() {
	init := int64(123456)
	suite.NextBlock(time.Unix(init, 0))
	suite.global.InitGenesis(suite.Ctx)

	suite.Equal(int64(0), suite.global.GetPastDay(suite.Ctx, init+3600*24-1))
	suite.Equal(int64(1), suite.global.GetPastDay(suite.Ctx, init+3600*24))
	suite.Equal(int64(1), suite.global.GetPastDay(suite.Ctx, init+3600*24*2-1))
	suite.Equal(int64(2), suite.global.GetPastDay(suite.Ctx, init+3600*24*2))
}

func (suite *globalManagerTestSuite) TestOnBeginBlockExecuted() {
	init := int64(123456)
	suite.NextBlock(time.Unix(init, 0))
	suite.global.InitGenesis(suite.Ctx)

	suite.mApp.On("Hourly", mock.Anything).Return(nil).Times(linotypes.HoursPerYear)
	suite.mApp.On("Daily", mock.Anything).Return(nil).Times(linotypes.HoursPerYear / 24)
	suite.mApp.On("Monthly", mock.Anything).Return(nil).Times(12)
	suite.mApp.On("Yearly", mock.Anything).Return(nil).Times(1)
	for i := init + 3600; i <= init+linotypes.MinutesPerYear*60; i += 3600 {
		suite.NextBlock(time.Unix(i, 0))
		suite.global.OnBeginBlock(suite.Ctx)
		suite.global.OnEndBlock(suite.Ctx)
	}
	suite.mApp.AssertExpectations(suite.T())
}

func (suite *globalManagerTestSuite) TestOnBeginBlockErrLogged() {
	init := int64(123456)
	suite.NextBlock(time.Unix(init, 0))
	suite.global.InitGenesis(suite.Ctx)

	suite.mApp.On("Hourly", mock.Anything).Return(nil).Times(linotypes.HoursPerYear)
	suite.mApp.On("Daily", mock.Anything).Return(nil).Times(linotypes.HoursPerYear / 24)
	suite.mApp.On("Yearly", mock.Anything).Return(nil).Times(1)

	last := int64(0)
	for i := init + 3600; i <= init+linotypes.MinutesPerYear*60; i += 3600 {
		suite.NextBlock(time.Unix(i, 0))
		if (i-init)/60/linotypes.MinutesPerMonth > last {
			last = (i - init) / 60 / linotypes.MinutesPerMonth
			suite.mApp.On("Monthly", mock.Anything).Return([]linotypes.BCEventErr{
				linotypes.NewBCEventErr(suite.Ctx, linotypes.ErrTestDummyError(), "test"),
			}).Once()
		}
		suite.global.OnBeginBlock(suite.Ctx)
		suite.global.OnEndBlock(suite.Ctx)
	}
	suite.mApp.AssertExpectations(suite.T())
	suite.Equal(12, len(suite.global.GetBCEventErrors(suite.Ctx)))
	suite.Golden()
}

func (suite *globalManagerTestSuite) TestRegisterEventsAndExec() {
	init := int64(123456)
	suite.NextBlock(time.Unix(init, 0))
	suite.global.InitGenesis(suite.Ctx)
	suite.mApp.On("Hourly", mock.Anything).Return(nil)
	suite.mApp.On("Daily", mock.Anything).Return(nil)
	suite.mApp.On("Monthly", mock.Anything).Return(nil)
	suite.mApp.On("Yearly", mock.Anything).Return(nil)


	// a simple counter event
	nExecuted := int64(0)
	exec := func(ctx sdk.Context, event linotypes.Event) sdk.Error {
		e := event.(testEvent)
		nExecuted++
		if e.Id%7 == 0 {
			return linotypes.ErrTestDummyError()
		}
		return nil
	}

	// register event of past time.
	suite.NotNil(suite.global.RegisterEventAtTime(suite.Ctx, init-1, testEvent{}))

	// events will be executed starting from an hour later
	nScheduled := int64(0)
	for i := init + 3; i <= init+7*3600; i += 600 {
		suite.global.RegisterEventAtTime(suite.Ctx, i, testEvent{Id: nScheduled})
		nScheduled++
	}
	// events can be scheduled at same time
	for i := init + 3 + 600*4; i <= init+7*3600; i += 600 {
		suite.global.RegisterEventAtTime(suite.Ctx, i, testEvent{Id: nScheduled})
		nScheduled++
	}

	// while execution, new events are scheduled, but not executed
	for i := init + 1; i <= init+10*3600; i += 577 {
		suite.NextBlock(time.Unix(i, 0))
		suite.global.OnBeginBlock(suite.Ctx)
		suite.global.ExecuteEvents(suite.Ctx, exec)
		suite.global.RegisterEventAtTime(suite.Ctx, i+10000000000, testEvent{Id: -1})
		suite.global.OnEndBlock(suite.Ctx)
	}

	suite.Equal(nScheduled, nExecuted)
	suite.Golden()
}

func (suite *globalManagerTestSuite) TestEventErrIsolation() {
	init := int64(123456)
	suite.NextBlock(time.Unix(init, 0))
	suite.global.InitGenesis(suite.Ctx)
	suite.mApp.On("Hourly", mock.Anything).Return(nil)
	suite.mApp.On("Daily", mock.Anything).Return(nil)
	suite.mApp.On("Monthly", mock.Anything).Return(nil)
	suite.mApp.On("Yearly", mock.Anything).Return(nil)

	bad := func(ctx sdk.Context, event linotypes.Event) sdk.Error {
		// this bad exec will reset global time to zero while return an error.
		suite.global.storage.SetGlobalTime(ctx, &model.GlobalTime{})
		return linotypes.ErrTestDummyError()
	}

	err := suite.global.RegisterEventAtTime(suite.Ctx, init+30, testEvent{})
	suite.Nil(err)

	suite.NextBlock(time.Unix(init+30+10, 0))
	suite.global.ExecuteEvents(suite.Ctx, bad)
	suite.Equal(1, len(suite.global.GetEventErrors(suite.Ctx)))             // err logged
	suite.Equal(init, suite.global.GetGlobalTime(suite.Ctx).ChainStartTime) // state not changed
	suite.Golden()
}
