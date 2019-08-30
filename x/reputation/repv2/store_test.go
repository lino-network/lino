package repv2

import (
	"testing"

	"github.com/lino-network/lino/x/reputation/repv2/internal"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	mockDB *internal.MockStore
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, &StoreTestSuite{})
}

func (suite *StoreTestSuite) SetupTest() {
	suite.mockDB = internal.NewMockStore()
}

func (suite *StoreTestSuite) TestPrefix() {
	suite.Equal(append(repUserMetaPrefix, []byte("qwe")...), getUserMetaKey("qwe"))
	suite.Equal(append(repRoundMetaPrefix, []byte("3")...), getRoundMetaKey(3))
	suite.Equal(append(repRoundPostMetaPrefix,
		[]byte{byte('b'), byte('/'), byte('x'), byte('y')}...), getRoundPostMetaKey(11, "xy"))
	suite.Equal(append(repRoundPostMetaPrefix,
		[]byte{byte('z'), byte('/'), byte('x'), byte('y')}...), getRoundPostMetaKey(35, "xy"))
	suite.Equal(append(repRoundPostMetaPrefix,
		[]byte{byte('2'), byte('f'), byte('/'), byte('a'), byte('b'), byte('c'), byte('d')}...),
		getRoundPostMetaKey(87, "abcd"))
	suite.Equal(repGameMetaPrefix, getGameKey())
}

func (suite *StoreTestSuite) TestInitValues() {
	store := NewReputationStore(suite.mockDB, DefaultInitialReputation)

	// user
	user := store.GetUserMeta("no")
	suite.Equal(NewInt(DefaultInitialReputation), user.Reputation)
	suite.Equal(RoundId(0), user.LastSettledRound)
	suite.Equal(RoundId(0), user.LastDonationRound)
	suite.Empty(user.Unsettled)

	// round meta
	round := store.GetRoundMeta(333)
	suite.Empty(round.Result)
	suite.Equal(NewInt(0), round.SumIF)
	suite.Equal(Time(0), round.StartAt)
	suite.Empty(round.TopN)

	// current round
	suite.Equal(RoundId(1), store.GetCurrentRound())

	// post
	post := store.GetRoundPostMeta(33, "xxx")
	suite.Equal(NewInt(0), post.SumIF)

	// game
	game := store.GetGameMeta()
	suite.Equal(RoundId(1), game.CurrentRound)
}

func (suite *StoreTestSuite) TestStoreGetSet() {
	store := NewReputationStore(suite.mockDB, DefaultInitialReputation)

	var user1 Uid = "test"
	var user2 Uid = "test2"
	var post1 Pid = "post1"
	var post2 Pid = "post2"

	u1 := &userMeta{
		Consumption:       NewInt(0),
		Hold:              NewInt(0),
		Reputation:        NewInt(123),
		LastSettledRound:  3,
		LastDonationRound: 3,
		Unsettled: []Donation{
			Donation{Pid: post1, Amount: NewInt(3), Impact: NewInt(2)},
			Donation{Pid: post2, Amount: NewInt(4), Impact: NewInt(7)},
		}}
	u2 := &userMeta{
		Consumption:       NewInt(0),
		Hold:              NewInt(0),
		Reputation:        NewInt(456),
		LastSettledRound:  4,
		LastDonationRound: 5,
		Unsettled: []Donation{
			Donation{Pid: post2, Amount: NewInt(6), Impact: NewInt(11)},
		}}
	store.SetUserMeta(user1, u2)
	store.SetUserMeta(user1, u1)
	store.SetUserMeta(user2, u2)
	defer func() {
		suite.Equal(u1, store.GetUserMeta(user1))
		suite.Equal(u2, store.GetUserMeta(user2))
	}()

	round1 := &roundMeta{
		Result:  []Pid{"123", "4ed6", "xzz"},
		SumIF:   NewInt(33333),
		StartAt: 324324,
		TopN:    nil,
	}
	round2 := &roundMeta{
		Result:  []Pid{"xzz"},
		SumIF:   NewInt(234134),
		StartAt: 342,
		TopN: []PostIFPair{
			PostIFPair{
				Pid:   post1,
				SumIF: NewInt(234235311),
			},
		},
	}
	store.SetRoundMeta(3, round2)
	store.SetRoundMeta(3, round1)
	store.SetRoundMeta(4, round2)
	defer func() {
		suite.Equal(round1, store.GetRoundMeta(3))
		suite.Equal(round2, store.GetRoundMeta(4))
	}()

	rp1 := &roundPostMeta{
		SumIF: NewInt(342),
	}
	rp2 := &roundPostMeta{
		SumIF: NewInt(666),
	}

	store.SetRoundPostMeta(123, post1, rp2)
	store.SetRoundPostMeta(123, post1, rp1)
	store.SetRoundPostMeta(342, post1, rp2)
	defer func() {
		suite.Equal(rp1, store.GetRoundPostMeta(123, post1))
		suite.Equal(rp2, store.GetRoundPostMeta(342, post1))
	}()

	store.SetGameMeta(&gameMeta{CurrentRound: 33})
	store.SetGameMeta(&gameMeta{CurrentRound: 443})
	defer func() {
		suite.Equal(RoundId(443), store.GetCurrentRound())
		suite.Equal(&gameMeta{CurrentRound: 443}, store.GetGameMeta())
	}()
}

func (suite *StoreTestSuite) TestStoreImportExporter() {
	store := NewReputationStore(suite.mockDB, DefaultInitialReputation)

	var user1 Uid = "test"
	var user2 Uid = "test2"
	var post1 Pid = "post1"
	var post2 Pid = "post2"

	u1 := &userMeta{
		Reputation:        NewInt(123),
		LastSettledRound:  3,
		LastDonationRound: 3,
		Unsettled: []Donation{
			Donation{Pid: post1, Amount: NewInt(3), Impact: NewInt(2)},
			Donation{Pid: post2, Amount: NewInt(4), Impact: NewInt(7)},
		}}
	u2 := &userMeta{
		Reputation:        NewInt(456),
		LastSettledRound:  4,
		LastDonationRound: 5,
		Unsettled: []Donation{
			Donation{Pid: post2, Amount: NewInt(6), Impact: NewInt(11)},
		}}
	store.SetUserMeta(user1, u1)
	store.SetUserMeta(user2, u2)

	// export data
	data := store.Export()
	db2 := internal.NewMockStore()
	store2 := NewReputationStore(db2, DefaultInitialReputation)
	store2.Import(data)
	suite.Equal(u1.Reputation, store2.GetUserMeta(user1).Reputation)
	suite.Equal(u2.Reputation, store2.GetUserMeta(user2).Reputation)
	suite.Equal(RoundId(0), store2.GetUserMeta(user1).LastDonationRound)
	suite.Equal(RoundId(0), store2.GetUserMeta(user2).LastDonationRound)
	suite.Equal(RoundId(0), store2.GetUserMeta(user1).LastSettledRound)
	suite.Equal(RoundId(0), store2.GetUserMeta(user2).LastSettledRound)
	suite.Empty(store2.GetUserMeta(user1).Unsettled)
	suite.Empty(store2.GetUserMeta(user2).Unsettled)
}

func (suite *StoreTestSuite) TestStoreImportExporterFromUpgrade1() {
	store := NewReputationStore(suite.mockDB, DefaultInitialReputation)

	var user1 Uid = "test"
	u1 := &userMeta{
		Reputation:        NewInt(100000),
		LastSettledRound:  3,
		LastDonationRound: 3,
		Unsettled:         []Donation{}}
	store.SetUserMeta(user1, u1)

	// export data
	data := store.Export()
	data.Reputations[0].IsMiniDollar = false

	db2 := internal.NewMockStore()
	store2 := NewReputationStore(db2, DefaultInitialReputation)
	store2.Import(data)
	u1.Reputation.Mul(NewInt(12))
	suite.Equal(u1.Reputation, store2.GetUserMeta(user1).Reputation)
}
