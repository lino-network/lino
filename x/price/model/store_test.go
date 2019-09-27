package model

import (
	// "fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
)

var (
	storeKeyStr = "testPriceStore"
	kvStoreKey  = sdk.NewKVStoreKey(storeKeyStr)
)

type PriceStoreDumper struct{}

func (dumper PriceStoreDumper) NewDumper() *testutils.Dumper {
	return NewPriceDumper(NewPriceStorage(kvStoreKey))
}

type priceStoreTestSuite struct {
	testsuites.GoldenTestSuite
	store PriceStorage
}

func NewPriceStoreTestSuite() *priceStoreTestSuite {
	return &priceStoreTestSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(PriceStoreDumper{}, kvStoreKey),
	}
}

func (suite *priceStoreTestSuite) SetupTest() {
	suite.SetupCtx(0, time.Unix(0, 0), kvStoreKey)
	suite.store = NewPriceStorage(kvStoreKey)
}

func TestPriceStoreSuite(t *testing.T) {
	suite.Run(t, NewPriceStoreTestSuite())
}

func (suite *priceStoreTestSuite) TestGetSetFedPrice() {
	store := suite.store
	ctx := suite.Ctx
	user1 := linotypes.AccountKey("user1")
	user2 := linotypes.AccountKey("user2")
	time1 := time.Unix(1000, 0)
	time2 := time.Unix(2000, 0)
	fed1 := &FedPrice{
		Validator: user1,
		Price:     linotypes.NewMiniDollar(123),
		UpdateAt:   time1.Unix(),
	}
	fed2 := &FedPrice{
		Validator: user2,
		Price:     linotypes.NewMiniDollar(456),
		UpdateAt:   time2.Unix(),
	}
	_, err := store.GetFedPrice(ctx, user1)
	suite.Require().NotNil(err)

	store.SetFedPrice(ctx, fed1)
	store.SetFedPrice(ctx, fed2)

	price, err := store.GetFedPrice(ctx, user1)
	suite.Require().Nil(err)
	suite.Equal(fed1, price)
	price, err = store.GetFedPrice(ctx, user2)
	suite.Require().Nil(err)
	suite.Equal(fed2, price)

	suite.Golden()
}

func (suite *priceStoreTestSuite) TestGetSetPriceHistory() {
	store := suite.store
	ctx := suite.Ctx
	time1 := time.Unix(1000, 0)
	time2 := time.Unix(2000, 0)
	time3 := time.Unix(3000, 0)
	tp1 := TimePrice{
		Price:    linotypes.NewMiniDollar(123),
		UpdateAt: time1.Unix(),
	}
	tp2 := TimePrice{
		Price:    linotypes.NewMiniDollar(456),
		UpdateAt: time2.Unix(),
	}
	tp3 := TimePrice{
		Price:    linotypes.NewMiniDollar(789),
		UpdateAt: time3.Unix(),
	}

	history1 := []TimePrice{tp1, tp2}
	history2 := []TimePrice{tp1, tp2, tp3}
	history3 := []TimePrice{tp2, tp3}

	rst := store.GetPriceHistory(ctx)
	suite.Require().Equal(0, len(rst))

	store.SetPriceHistory(ctx, history1)
	rst = store.GetPriceHistory(ctx)
	suite.Equal(history1, rst)

	store.SetPriceHistory(ctx, history2)
	rst = store.GetPriceHistory(ctx)
	suite.Equal(history2, rst)

	store.SetPriceHistory(ctx, history3)
	rst = store.GetPriceHistory(ctx)
	suite.Equal(history3, rst)

	suite.Golden()
}

func (suite *priceStoreTestSuite) TestGetSetCurrentPrice() {
	store := suite.store
	ctx := suite.Ctx
	time1 := time.Unix(1000, 0)
	time2 := time.Unix(2000, 0)
	tp1 := &TimePrice{
		Price:    linotypes.NewMiniDollar(123),
		UpdateAt: time1.Unix(),
	}
	tp2 := &TimePrice{
		Price:    linotypes.NewMiniDollar(456),
		UpdateAt: time2.Unix(),
	}

	_, err := store.GetCurrentPrice(ctx)
	suite.Require().NotNil(err)

	store.SetCurrentPrice(ctx, tp1)
	rst, err := store.GetCurrentPrice(ctx)
	suite.Require().Nil(err)
	suite.Require().Equal(tp1, rst)

	store.SetCurrentPrice(ctx, tp2)
	rst, err = store.GetCurrentPrice(ctx)
	suite.Require().Nil(err)
	suite.Require().Equal(tp2, rst)

	suite.Golden()
}

func (suite *priceStoreTestSuite) TestGetSetLastValidators() {
	store := suite.store
	ctx := suite.Ctx
	eg1 := []linotypes.AccountKey{"user1", "user2"}
	eg2 := []linotypes.AccountKey{"user3", "user4"}

	last := store.GetLastValidators(ctx)
	suite.Equal(0, len(last))

	store.SetLastValidators(ctx, eg1)
	suite.Equal(eg1, store.GetLastValidators(ctx))
	store.SetLastValidators(ctx, eg2)
	suite.Equal(eg2, store.GetLastValidators(ctx))

	suite.Golden()
}
