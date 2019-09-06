package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	linotypes "github.com/lino-network/lino/types"

	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

type DevStorageTestSuite struct {
	suite.Suite
	ctx sdk.Context
	ds  DeveloperStorage
}

func TestDevStorageTestSuite(t *testing.T) {
	suite.Run(t, new(DevStorageTestSuite))
}

func (suite *DevStorageTestSuite) SetupTest() {
	TestKVStoreKey := sdk.NewKVStoreKey("dev")
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()
	suite.ctx = sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	suite.ds = NewDeveloperStorage(TestKVStoreKey)
}

func (suite *DevStorageTestSuite) TestDeveloperGetSetHas() {
	app1 := linotypes.AccountKey("app1")
	app2 := linotypes.AccountKey("app2")
	app3 := linotypes.AccountKey("app3")
	suite.False(suite.ds.HasDeveloper(suite.ctx, app1))
	suite.False(suite.ds.HasDeveloper(suite.ctx, app2))

	dev1 := Developer{
		Username:       app1,
		Deposit:        linotypes.NewCoinFromInt64(123),
		AppConsumption: linotypes.NewMiniDollar(456),
		Website:        "web.com",
		Description:    "test app",
		AppMetaData:    "empty1",
		IsDeleted:      false,
		NAffiliated:    0,
	}
	dev2 := Developer{
		Username:       app2,
		Deposit:        linotypes.NewCoinFromInt64(333),
		AppConsumption: linotypes.NewMiniDollar(444),
		Website:        "xyz",
		Description:    "test app 2",
		AppMetaData:    "nonnon",
		IsDeleted:      true,
		NAffiliated:    100,
	}

	suite.ds.SetDeveloper(suite.ctx, dev1)
	suite.ds.SetDeveloper(suite.ctx, dev2)

	rst1, err := suite.ds.GetDeveloper(suite.ctx, app1)
	suite.Require().Nil(err)
	suite.Equal(&dev1, rst1)
	rst2, err := suite.ds.GetDeveloper(suite.ctx, app2)
	suite.Require().Nil(err)
	suite.Equal(&dev2, rst2)
	_, err = suite.ds.GetDeveloper(suite.ctx, app3)
	suite.Require().NotNil(err)

	suite.True(suite.ds.HasDeveloper(suite.ctx, app1))
	suite.True(suite.ds.HasDeveloper(suite.ctx, app2))
	suite.False(suite.ds.HasDeveloper(suite.ctx, app3))

	alldevs := suite.ds.GetAllDevelopers(suite.ctx)
	suite.Require().Equal(2, len(alldevs))
	suite.Equal(dev1, alldevs[0])
	suite.Equal(dev2, alldevs[1])
}

func (suite *DevStorageTestSuite) TestIDAGetSetHas() {
	app1 := linotypes.AccountKey("app1")
	app2 := linotypes.AccountKey("app2")
	app3 := linotypes.AccountKey("app3")
	suite.False(suite.ds.HasIDA(suite.ctx, app1))
	suite.False(suite.ds.HasIDA(suite.ctx, app2))

	ida1 := AppIDA{
		App:             app1,
		Name:            "COOL",
		MiniIDAPrice:    linotypes.NewMiniDollar(99),
		IsRevoked:       false,
		RevokeCoinPrice: linotypes.NewMiniDollar(0),
	}
	ida2 := AppIDA{
		App:             app2,
		Name:            "DUDE",
		MiniIDAPrice:    linotypes.NewMiniDollar(5),
		IsRevoked:       true,
		RevokeCoinPrice: linotypes.NewMiniDollar(7),
	}

	suite.ds.SetIDA(suite.ctx, ida1)
	suite.ds.SetIDA(suite.ctx, ida2)

	rst1, err := suite.ds.GetIDA(suite.ctx, app1)
	suite.Require().Nil(err)
	suite.Equal(&ida1, rst1)
	rst2, err := suite.ds.GetIDA(suite.ctx, app2)
	suite.Require().Nil(err)
	suite.Equal(&ida2, rst2)
	_, err = suite.ds.GetIDA(suite.ctx, app3)
	suite.Require().NotNil(err)

	_, err = suite.ds.GetDeveloper(suite.ctx, app3)
	suite.Require().NotNil(err)

	suite.True(suite.ds.HasIDA(suite.ctx, app1))
	suite.True(suite.ds.HasIDA(suite.ctx, app2))
	suite.False(suite.ds.HasIDA(suite.ctx, app3))
}

func (suite *DevStorageTestSuite) TestIDAStatsGetSet() {
	app1 := linotypes.AccountKey("app1")
	suite.Equal(&AppIDAStats{
		Total: linotypes.NewMiniDollar(0),
	}, suite.ds.GetIDAStats(suite.ctx, app1))
	suite.ds.SetIDAStats(suite.ctx, app1, AppIDAStats{
		Total: linotypes.NewMiniDollar(3),
	})
	suite.Equal(&AppIDAStats{
		Total: linotypes.NewMiniDollar(3),
	}, suite.ds.GetIDAStats(suite.ctx, app1))
}

func (suite *DevStorageTestSuite) TestUserRoleGetSetHasDel() {
	app1 := linotypes.AccountKey("app1")
	user1 := linotypes.AccountKey("user1")
	app2 := linotypes.AccountKey("app2")
	user2 := linotypes.AccountKey("user2")
	user3 := linotypes.AccountKey("user3")
	role1 := &Role{
		AffiliatedApp: app1,
	}
	role2 := &Role{
		AffiliatedApp: app2,
	}

	_, err := suite.ds.GetUserRole(suite.ctx, user1)
	suite.Require().NotNil(err)
	_, err = suite.ds.GetUserRole(suite.ctx, user2)
	suite.Require().NotNil(err)
	suite.False(suite.ds.HasUserRole(suite.ctx, user1))
	suite.False(suite.ds.HasUserRole(suite.ctx, user2))
	suite.False(suite.ds.HasUserRole(suite.ctx, user3))

	suite.ds.SetUserRole(suite.ctx, user1, role1)
	suite.ds.SetUserRole(suite.ctx, user2, role2)

	rst1, err := suite.ds.GetUserRole(suite.ctx, user1)
	suite.Require().Nil(err)
	suite.Equal(role1, rst1)
	rst2, err := suite.ds.GetUserRole(suite.ctx, user2)
	suite.Require().Nil(err)
	suite.Equal(role2, rst2)
	suite.True(suite.ds.HasUserRole(suite.ctx, user1))
	suite.True(suite.ds.HasUserRole(suite.ctx, user2))
	suite.False(suite.ds.HasUserRole(suite.ctx, user3))

	suite.ds.DelUserRole(suite.ctx, user2)
	suite.ds.DelUserRole(suite.ctx, user3)

	rst1, err = suite.ds.GetUserRole(suite.ctx, user1)
	suite.Require().Nil(err)
	suite.Equal(role1, rst1)
	_, err = suite.ds.GetUserRole(suite.ctx, user2)
	suite.Require().NotNil(err)
	suite.True(suite.ds.HasUserRole(suite.ctx, user1))
	suite.False(suite.ds.HasUserRole(suite.ctx, user2))
	suite.False(suite.ds.HasUserRole(suite.ctx, user3))
}

func (suite *DevStorageTestSuite) TestGetSetIDABank() {
	app1 := linotypes.AccountKey("app1")
	user1 := linotypes.AccountKey("user1")
	app2 := linotypes.AccountKey("app2")
	user2 := linotypes.AccountKey("user2")

	empty := suite.ds.GetIDABank(suite.ctx, app1, user1)
	suite.Equal(&IDABank{
		Balance:  linotypes.NewMiniDollar(0),
		Unauthed: false,
	}, empty)

	suite.ds.SetIDABank(suite.ctx, app1, user1, &IDABank{
		Balance:  linotypes.NewMiniDollar(123),
		Unauthed: false,
	})
	suite.Equal(&IDABank{
		Balance:  linotypes.NewMiniDollar(123),
		Unauthed: false,
	}, suite.ds.GetIDABank(suite.ctx, app1, user1))

	suite.Equal(&IDABank{
		Balance:  linotypes.NewMiniDollar(0),
		Unauthed: false,
	}, suite.ds.GetIDABank(suite.ctx, app2, user2))
}

func (suite *DevStorageTestSuite) TestSetHasDelItrAffiliated() {
	app1 := linotypes.AccountKey("app1")
	user1 := linotypes.AccountKey("user1")
	app2 := linotypes.AccountKey("app2")
	user2 := linotypes.AccountKey("user2")

	suite.False(suite.ds.HasAffiliatedAcc(suite.ctx, app1, user1))
	suite.False(suite.ds.HasAffiliatedAcc(suite.ctx, app2, user2))

	suite.ds.SetAffiliatedAcc(suite.ctx, app1, user1)
	suite.ds.SetAffiliatedAcc(suite.ctx, app1, user2)
	suite.True(suite.ds.HasAffiliatedAcc(suite.ctx, app1, user1))
	suite.True(suite.ds.HasAffiliatedAcc(suite.ctx, app1, user2))
	suite.False(suite.ds.HasAffiliatedAcc(suite.ctx, app2, user2))

	all := suite.ds.GetAllAffiliatedAcc(suite.ctx, app1)
	suite.Require().Equal(2, len(all))
	suite.Equal([]linotypes.AccountKey{user1, user2}, all)

	all2 := suite.ds.GetAllAffiliatedAcc(suite.ctx, app2)
	suite.Equal(0, len(all2))

	suite.ds.DelAffiliatedAcc(suite.ctx, app1, user1)
	suite.ds.DelAffiliatedAcc(suite.ctx, app2, user2)
	suite.False(suite.ds.HasAffiliatedAcc(suite.ctx, app1, user1))
	suite.True(suite.ds.HasAffiliatedAcc(suite.ctx, app1, user2))

	all = suite.ds.GetAllAffiliatedAcc(suite.ctx, app1)
	suite.Require().Equal(1, len(all))
	suite.Equal([]linotypes.AccountKey{user2}, all)
}

func (suite *DevStorageTestSuite) TestGetSetReservePool() {
	suite.Panics(func() {
		suite.ds.GetReservePool(suite.ctx)
	})
	suite.ds.SetReservePool(suite.ctx, &ReservePool{
		Total:           linotypes.NewCoinFromInt64(123),
		TotalMiniDollar: linotypes.NewMiniDollar(456),
	})
	suite.Equal(&ReservePool{
		Total:           linotypes.NewCoinFromInt64(123),
		TotalMiniDollar: linotypes.NewMiniDollar(456),
	}, suite.ds.GetReservePool(suite.ctx))
}

func (suite *DevStorageTestSuite) TestGetSetEmptyReservePool() {
	suite.ds.SetReservePool(suite.ctx, &ReservePool{})
	suite.Equal(&ReservePool{
		Total:           linotypes.NewCoinFromInt64(0),
		TotalMiniDollar: linotypes.NewMiniDollar(0),
	}, suite.ds.GetReservePool(suite.ctx))
}
