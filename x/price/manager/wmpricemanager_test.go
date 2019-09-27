package manager

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	mparam "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
	// maccount "github.com/lino-network/lino/x/account/mocks"
	"github.com/lino-network/lino/x/price/model"
	"github.com/lino-network/lino/x/price/types"
	// mglobal "github.com/lino-network/lino/x/global/mocks"
	// mprice "github.com/lino-network/lino/x/price/mocks"
	// mvote "github.com/lino-network/lino/x/vote/mocks"
	// votetypes "github.com/lino-network/lino/x/vote/types"
)

var (
	storeKeyStr = "priceStoreTestKey"
	storeKey    = sdk.NewKVStoreKey(storeKeyStr)
)

type PriceStoreDumper struct{}

func (dumper PriceStoreDumper) NewDumper() *testutils.Dumper {
	return model.NewPriceDumper(model.NewPriceStorage(storeKey))
}

type WMPriceManagerSuite struct {
	testsuites.GoldenTestSuite
	manager      WeightedMedianPriceManager
	mParamKeeper *mparam.ParamKeeper
}

func NewPriceManagerSuite() *WMPriceManagerSuite {
	return &WMPriceManagerSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(PriceStoreDumper{}, storeKey),
	}
}

func (suite *WMPriceManagerSuite) SetupTest() {
	suite.mParamKeeper = new(mparam.ParamKeeper)
	suite.manager = NewPriceManager(storeKey, suite.mParamKeeper, suite.mVoteKeeper, suite.mAccountKeeper, suite.mPriceKeeper, suite.mGlobalKeeper)
	suite.SetupCtx(0, time.Unix(0, 0), storeKey)
}

func TestPriceManagerSuite(t *testing.T) {
	suite.Run(t, NewPriceManagerSuite())
}
