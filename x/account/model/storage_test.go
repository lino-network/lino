package model

import (
	"testing"
	"time"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/types"
)

var (
	storeKeyStr = "testAccountStore"
	kvStoreKey  = sdk.NewKVStoreKey(storeKeyStr)
)

type AccountStoreDumper struct{}

func (dumper AccountStoreDumper) NewDumper() *testutils.Dumper {
	return NewAccountDumper(NewAccountStorage(kvStoreKey))
}

type accountStoreTestSuite struct {
	testsuites.GoldenTestSuite
	store AccountStorage
}

func NewAccountStoreTestSuite() *accountStoreTestSuite {
	return &accountStoreTestSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(AccountStoreDumper{}, kvStoreKey),
	}
}

func (suite *accountStoreTestSuite) SetupTest() {
	suite.SetupCtx(0, time.Unix(0, 0), kvStoreKey)
	suite.store = NewAccountStorage(kvStoreKey)
}

func TestAccountStoreSuite(t *testing.T) {
	suite.Run(t, NewAccountStoreTestSuite())
}

func (suite *accountStoreTestSuite) TestInfo() {
	user1 := linotypes.AccountKey("user1")
	user2 := linotypes.AccountKey("user2")

	store := suite.store
	ctx := suite.Ctx
	keys := sampleKeys()

	info1 := &AccountInfo{
		Username:       user1,
		CreatedAt:      123,
		SigningKey:     keys[9],
		TransactionKey: keys[0],
		Address:        sdk.AccAddress(keys[0].Address()),
	}
	info2 := &AccountInfo{
		Username:       user2,
		CreatedAt:      456,
		SigningKey:     keys[8],
		TransactionKey: keys[1],
		Address:        sdk.AccAddress(keys[1].Address()),
	}

	suite.False(store.DoesAccountExist(ctx, user1))
	suite.False(store.DoesAccountExist(ctx, user2))

	_, err := store.GetInfo(ctx, user1)
	suite.Equal(types.ErrAccountNotFound(user1), err)

	suite.store.SetInfo(ctx, info1)
	suite.store.SetInfo(ctx, info2)

	r1, err := store.GetInfo(ctx, user1)
	suite.Nil(err)
	suite.Equal(info1, r1)

	r2, err := store.GetInfo(ctx, user2)
	suite.Nil(err)
	suite.Equal(info2, r2)

	suite.True(store.DoesAccountExist(ctx, user1))
	suite.True(store.DoesAccountExist(ctx, user2))

	suite.Golden()
}

func (suite *accountStoreTestSuite) TestBank() {
	user1 := linotypes.AccountKey("user1")
	user2 := linotypes.AccountKey("user2")

	store := suite.store
	ctx := suite.Ctx
	keys := sampleKeys()

	key1 := keys[0]
	addr1 := sdk.AccAddress(key1.Address())
	bank1 := &AccountBank{
		Saving:   linotypes.NewCoinFromInt64(1234),
		Pending:  linotypes.NewCoinFromInt64(0),
		PubKey:   key1,
		Sequence: 123,
		Username: user1,
	}

	key2 := keys[1]
	addr2 := sdk.AccAddress(key2.Address())
	bank2 := &AccountBank{
		Saving:   linotypes.NewCoinFromInt64(2345),
		Pending:  linotypes.NewCoinFromInt64(789),
		PubKey:   key2,
		Sequence: 456,
		Username: user2,
	}

	_, err := store.GetBank(ctx, addr1)
	suite.Equal(types.ErrAccountBankNotFound(addr1), err)
	_, err = store.GetBank(ctx, addr2)
	suite.Equal(types.ErrAccountBankNotFound(addr2), err)

	suite.store.SetBank(ctx, addr1, bank1)
	suite.store.SetBank(ctx, addr2, bank2)

	r1, err := store.GetBank(ctx, addr1)
	suite.Nil(err)
	suite.Equal(bank1, r1)

	r2, err := store.GetBank(ctx, addr2)
	suite.Nil(err)
	suite.Equal(bank2, r2)

	suite.Golden()
}

func (suite *accountStoreTestSuite) TestMeta() {
	user1 := linotypes.AccountKey("user1")
	user2 := linotypes.AccountKey("user2")

	store := suite.store
	ctx := suite.Ctx

	meta1 := &AccountMeta{
		JSONMeta: "lol",
	}
	meta2 := &AccountMeta{
		JSONMeta: "dota",
	}
	store.SetMeta(ctx, user1, meta1)
	store.SetMeta(ctx, user2, meta2)

	r1 := store.GetMeta(ctx, user1)
	suite.Equal(meta1, r1)

	r2 := store.GetMeta(ctx, user2)
	suite.Equal(meta2, r2)

	suite.Equal(&AccountMeta{
		JSONMeta: "",
	}, store.GetMeta(ctx, "others"))

	suite.Golden()
}

func (suite *accountStoreTestSuite) TestPool() {
	store := suite.store
	ctx := suite.Ctx

	_, err := store.GetPool(ctx, linotypes.InflationConsumptionPool)
	suite.NotNil(err)

	pool1 := &Pool{
		Name:    linotypes.InflationValidatorPool,
		Balance: linotypes.NewCoinFromInt64(1234),
	}
	pool2 := &Pool{
		Name:    linotypes.InflationDeveloperPool,
		Balance: linotypes.NewCoinFromInt64(45567),
	}

	store.SetPool(ctx, pool1)
	store.SetPool(ctx, pool2)

	r1, err := store.GetPool(ctx, linotypes.InflationValidatorPool)
	suite.Nil(err)
	suite.Equal(pool1, r1)

	r2, err := store.GetPool(ctx, linotypes.InflationDeveloperPool)
	suite.Nil(err)
	suite.Equal(pool2, r2)

	suite.Golden()
}

func (suite *accountStoreTestSuite) TestSupply() {
	store := suite.store
	ctx := suite.Ctx

	suite.Panics(func() { store.GetSupply(ctx) })

	supply := &Supply{
		LastYearTotal:     linotypes.NewCoinFromInt64(10000),
		Total:             linotypes.NewCoinFromInt64(512354),
		ChainStartTime:    123,
		LastInflationTime: 3245,
	}

	store.SetSupply(ctx, supply)

	suite.Equal(store.GetSupply(ctx), supply)

	suite.Golden()
}

// cdc := wire.New()
// wire.RegisterCrypto(cdc)
// keys := make([]crypto.PubKey, 0)
// for i := 0 ; i < 10; i++ {
// keys = append(keys, secp256k1.GenPrivKey().PubKey())
// }
// fmt.Print(string(cdc.MustMarshalJSON(keys)))
func sampleKeys() []crypto.PubKey {
	json := `
[{"type":"tendermint/PubKeySecp256k1","value":"Aot3u5m7vuxUOszkS6IZW5XYVu6ATvZsfSQIjtQo9tML"},{"type":"tendermint/PubKeySecp256k1","value":"AoFqbXKmblwKVggqb8Cqo30gRKs9EfqwhOhuyOKlGCuD"},{"type":"tendermint/PubKeySecp256k1","value":"Aj/1EOLKUKUPhp+mx3fLNoZOEEsY+tjPeTW4nOPbqwwq"},{"type":"tendermint/PubKeySecp256k1","value":"A1SxTVyDiXljmHeimniCQiNZQ3dcDsgppP0gDCMgJtdp"},{"type":"tendermint/PubKeySecp256k1","value":"Ax8b6HzTh9el9/NfE8fI4awCvMZWGQkjl+rYOGWeGJc9"},{"type":"tendermint/PubKeySecp256k1","value":"A4r+RjYEc2V9p43J4CovoktRTXNY9vvcQbx0aOW9bhoq"},{"type":"tendermint/PubKeySecp256k1","value":"AwFSpofxlQGAQv167WveHyeUvTh/3fukkJU7gkEW+iMm"},{"type":"tendermint/PubKeySecp256k1","value":"AjglddkWGGlMZck7uvWMDCtyqpNWSBy9HmnJV9vPnu2k"},{"type":"tendermint/PubKeySecp256k1","value":"A+KW7obJ0BpKqUWmY33svTBxGdTfRhmOym7A5imWWwGm"},{"type":"tendermint/PubKeySecp256k1","value":"A6P8IUdt9DKrYCe3/Tflt7DBdgFokRcCKkixt+UbhjZ8"}]
`

	keys := make([]crypto.PubKey, 0)
	cdc := wire.New()
	wire.RegisterCrypto(cdc)
	cdc.MustUnmarshalJSON([]byte(json), &keys)
	return keys
}
