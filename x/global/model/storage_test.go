package model

import (
	"testing"
	"time"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
)

type testEvent struct {
	Name string `json:"name"`
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
	return NewDumper(NewGlobalStorage(kvStoreKey, testCodec()), regTestCodec)
}

type globalStoreTestSuite struct {
	testsuites.GoldenTestSuite
	store GlobalStorage
}

func NewGlobalStoreTestSuite() *globalStoreTestSuite {
	return &globalStoreTestSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(storeDumper{}, kvStoreKey),
	}
}

func (suite *globalStoreTestSuite) SetupTest() {
	suite.SetupCtx(0, time.Unix(0, 0), kvStoreKey)
	suite.store = NewGlobalStorage(kvStoreKey, testCodec())
}

func TestGlobalStoreSuite(t *testing.T) {
	suite.Run(t, NewGlobalStoreTestSuite())
}

func (suite *globalStoreTestSuite) TestGetSetDelTimeEventList() {
	// empty not nil
	lst := suite.store.GetTimeEventList(suite.Ctx, 100)
	suite.NotNil(lst)
	suite.Equal(0, len(lst.Events))

	events1 := &linotypes.TimeEventList{
		Events: []linotypes.Event{
			testEvent{
				Name: "event1",
			},
			testEvent{
				Name: "event2",
			},
			testEvent{
				Name: "event3",
			},
		},
	}
	events2 := &linotypes.TimeEventList{
		Events: []linotypes.Event{
			testEvent{
				Name: "event4",
			},
		},
	}
	events3 := &linotypes.TimeEventList{
		Events: []linotypes.Event{
			testEvent{
				Name: "event6",
			},
		},
	}

	suite.store.SetTimeEventList(suite.Ctx, 123, events1)
	suite.store.SetTimeEventList(suite.Ctx, 456, events2)
	suite.store.SetTimeEventList(suite.Ctx, 789, events3)

	suite.Equal(events1, suite.store.GetTimeEventList(suite.Ctx, 123))
	suite.Equal(events2, suite.store.GetTimeEventList(suite.Ctx, 456))
	suite.Equal(events3, suite.store.GetTimeEventList(suite.Ctx, 789))

	suite.store.RemoveTimeEventList(suite.Ctx, 789)
	suite.Equal(0, len(suite.store.GetTimeEventList(suite.Ctx, 789).Events))

	suite.Golden()
}

func (suite *globalStoreTestSuite) TestGetSetGlobalTime() {
	suite.Panics(func() {
		suite.store.GetGlobalTime(suite.Ctx)
	})

	globalTime := &GlobalTime{
		ChainStartTime: 123,
		LastBlockTime:  234,
		PastMinutes:    345,
	}

	suite.store.SetGlobalTime(suite.Ctx, globalTime)
	suite.Equal(globalTime, suite.store.GetGlobalTime(suite.Ctx))

	globalTime.PastMinutes = 999
	suite.store.SetGlobalTime(suite.Ctx, globalTime)
	suite.Equal(globalTime, suite.store.GetGlobalTime(suite.Ctx))

	suite.Golden()
}

func (suite *globalStoreTestSuite) TestGetSetEventErrors() {
	suite.Nil(suite.store.GetEventErrors(suite.Ctx))

	errs := []EventError{
		{
			Time:    123,
			Event:   testEvent{Name: "lol"},
			ErrCode: sdk.CodeType(123),
		},
	}

	suite.store.SetEventErrors(suite.Ctx, errs)
	suite.Equal(errs, suite.store.GetEventErrors(suite.Ctx))

	errs = append(errs, EventError{
		Time:    456,
		Event:   testEvent{Name: "dota"},
		ErrCode: sdk.CodeType(999),
	})
	suite.store.SetEventErrors(suite.Ctx, errs)
	suite.Equal(errs, suite.store.GetEventErrors(suite.Ctx))

	suite.Golden()
}

func (suite *globalStoreTestSuite) TestGetSetBCErrors() {
	suite.Nil(suite.store.GetBCErrors(suite.Ctx))

	errs := []linotypes.BCEventErr{
		{
			Time:         123,
			ErrCode:      sdk.CodeType(123),
			ErrCodeSpace: sdk.CodespaceType(789),
			Reason:       "inflation failed",
		},
	}

	suite.store.SetBCErrors(suite.Ctx, errs)
	suite.Equal(errs, suite.store.GetBCErrors(suite.Ctx))

	errs = append(errs, linotypes.BCEventErr{
		Time:         456,
		ErrCode:      sdk.CodeType(34),
		ErrCodeSpace: sdk.CodespaceType(1585),
		Reason:       "validator failed",
	})
	suite.store.SetBCErrors(suite.Ctx, errs)
	suite.Equal(errs, suite.store.GetBCErrors(suite.Ctx))

	suite.Golden()
}
