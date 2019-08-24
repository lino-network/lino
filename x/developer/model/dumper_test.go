package model

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/testsuites"
	linotypes "github.com/lino-network/lino/types"
)

type DumperTestSuite struct {
	testsuites.CtxTestSuite
	store DeveloperStorage
	key   *sdk.KVStoreKey
}

func TestDumperTestSuite(t *testing.T) {
	suite.Run(t, new(DumperTestSuite))
}

func (suite *DumperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey("dumper")
	suite.key = key
	suite.SetupCtx(0, time.Unix(0, 0), suite.key)
	suite.store = NewDeveloperStorage(suite.key)
}

func (suite *DumperTestSuite) TestDump() {
	file, err := ioutil.TempFile(os.TempDir(), "linodev")
	suite.Require().Nil(err)
	defer os.Remove(file.Name())
	suite.Require().Nil(err)
	suite.store.SetDeveloper(suite.Ctx, Developer{
		Username:    "test",
		Website:     "web1",
		Description: "desc1",
	})
	suite.store.SetDeveloper(suite.Ctx, Developer{
		Username:    "test2",
		Website:     "web2",
		Description: "desc2",
	})
	suite.store.SetIDA(suite.Ctx, AppIDA{
		App:             "app1",
		Name:            "appname1",
		MiniIDAPrice:    linotypes.NewMiniDollar(0),
		IsRevoked:       true,
		RevokeCoinPrice: linotypes.NewMiniDollar(333),
	})
	suite.store.SetIDABank(suite.Ctx, "app1", "user1", &IDABank{
		Balance: linotypes.NewMiniDollar(33),
	})
	suite.store.SetAffiliatedAcc(suite.Ctx, "app1", "user1")
	suite.store.SetReservePool(suite.Ctx, &ReservePool{
		Total:           linotypes.NewCoinFromInt64(444),
		TotalMiniDollar: linotypes.NewMiniDollar(333),
	})
	suite.store.SetUserRole(suite.Ctx, "user2", &Role{
		AffiliatedApp: "app3",
	})
	DumpToFile(suite.Ctx, suite.store, file.Name())
}

func (suite *DumperTestSuite) TestLoad() {
	LoadFromFile(suite.Ctx, suite.store, "./load.json")
	dev, err := suite.store.GetDeveloper(suite.Ctx, "test")
	suite.Require().Nil(err)
	suite.Equal("test", string(dev.Username))
	dev, err = suite.store.GetDeveloper(suite.Ctx, "test2")
	suite.Require().Nil(err)
	suite.Equal("test2", string(dev.Username))
	ida, err := suite.store.GetIDA(suite.Ctx, "app1")
	suite.Require().Nil(err)
	suite.Equal(&AppIDA{
		App:             "app1",
		Name:            "appname1",
		MiniIDAPrice:    linotypes.NewMiniDollar(0),
		IsRevoked:       true,
		RevokeCoinPrice: linotypes.NewMiniDollar(333),
	}, ida)
	suite.True(suite.store.HasAffiliatedAcc(suite.Ctx, "app1", "user1"))
	suite.Equal(&ReservePool{
		Total:           linotypes.NewCoinFromInt64(444),
		TotalMiniDollar: linotypes.NewMiniDollar(333),
	}, suite.store.GetReservePool(suite.Ctx))
	role, err := suite.store.GetUserRole(suite.Ctx, "user2")
	suite.Require().Nil(err)
	suite.Equal(&Role{
		AffiliatedApp: "app3",
	}, role)
}
