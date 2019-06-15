package repv2

import (
	"math/big"
	"testing"
	"time"

	"github.com/lino-network/lino/x/reputation/repv2/internal"
	"github.com/stretchr/testify/suite"
)

type ReputationTestSuite struct {
	suite.Suite
	store                ReputationStore
	rep                  ReputationImpl
	roundDurationSeconds int64
	bestN                int
	userMaxN             int
	time                 time.Time
}

func TestReputationTestSuite(t *testing.T) {
	suite.Run(t, &ReputationTestSuite{})
}

func (suite *ReputationTestSuite) SetupTest() {
	suite.roundDurationSeconds = 25 * 3600
	suite.bestN = 30
	suite.userMaxN = 10
	suite.store = NewReputationStore(internal.NewMockStore())
	suite.rep = NewReputation(
		suite.store, suite.bestN, suite.userMaxN, suite.roundDurationSeconds).(ReputationImpl)
	suite.time = time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
}

func (suite *ReputationTestSuite) MoveToNewRound() {
	suite.time = suite.time.Add(time.Duration(suite.roundDurationSeconds) * time.Second)
	suite.rep.Update(suite.time.Unix())
}

// for bigInt(*big.Int), need this on comparing with zero value.
func (suite *ReputationTestSuite) EqualZero(a bigInt, args ...interface{}) {
	suite.Equal(0, a.Cmp(big.NewInt(0)), "%d is not bigInt zero", a.Int64())
}

func (suite *ReputationTestSuite) TestFirstBlock1() {
	rep := suite.rep
	newBlockTime := int64(0)
	rep.Update(0)
	rid, startAt := rep.GetCurrentRound()
	suite.Equal(int64(1), rid)
	suite.Equal(newBlockTime, startAt)
	suite.Equal(rep.GetReputation("me"), big.NewInt(InitialReputation))
}

func (suite *ReputationTestSuite) TestFirstBlock2() {
	rep := suite.rep
	newBlockTime := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
	rep.Update(newBlockTime.Unix())
	rid, startAt := rep.GetCurrentRound()
	suite.Equal(int64(2), rid)
	suite.Equal(newBlockTime.Unix(), startAt)

	nextBlockTime := time.Date(1995, time.February, 6, 12, 11, 0, 0, time.UTC)
	suite.Equal(newBlockTime.Add(time.Duration(suite.roundDurationSeconds)*time.Second), nextBlockTime)
	rep.Update(nextBlockTime.Unix())
	rid, startAt = rep.GetCurrentRound()
	suite.Equal(int64(3), rid)
	suite.Equal(nextBlockTime.Unix(), startAt)
}

func (suite *ReputationTestSuite) TestIncFreeScore() {
	rep := suite.rep
	rep.IncFreeScore("user1", big.NewInt(3000))
	suite.Equal(big.NewInt(3000+InitialReputation), rep.GetReputation("user1"))
}

func (suite *ReputationTestSuite) TestExtractConsumptionInfo() {
	newset := func(ids []Pid) map[Pid]bool {
		rst := make(map[Pid]bool)
		for _, id := range ids {
			rst[id] = true
		}
		return rst
	}
	cases := []struct {
		user     *userMeta
		seedSet  map[Pid]bool
		expected consumptionInfo
	}{
		{
			&userMeta{
				Unsettled: []Donation{
					{
						Pid:    "a",
						Amount: big.NewInt(333),
						Impact: big.NewInt(333),
					},
					{
						Pid:    "qf93f",
						Amount: big.NewInt(999),
						Impact: big.NewInt(222),
					},
					{
						Pid:    "b",
						Amount: big.NewInt(55555),
						Impact: big.NewInt(1),
					},
				},
			},
			newset([]Pid{"a", "b", "c", "ddd"}),
			consumptionInfo{
				seed:    big.NewInt(55888),
				other:   big.NewInt(999),
				seedIF:  big.NewInt(334),
				otherIF: big.NewInt(222),
			},
		},
		{
			&userMeta{
				Unsettled: []Donation{
					{
						Pid:    "not",
						Amount: big.NewInt(1),
						Impact: big.NewInt(1),
					},
				},
			},
			newset([]Pid{"a", "b", "c", "ddd"}),
			consumptionInfo{
				seed:    big.NewInt(0),
				other:   big.NewInt(1),
				seedIF:  big.NewInt(0),
				otherIF: big.NewInt(1),
			},
		},
	}

	for i, v := range cases {
		suite.Equal(v.expected, suite.rep.extractConsumptionInfo(v.user, v.seedSet), "case: %d", i)
	}
}

func (suite *ReputationTestSuite) TestGetSeedSet() {
	rep := suite.rep
	user1 := "user1"
	user2 := "user2"
	post1 := "post1"
	post2 := "post2"
	rep.IncFreeScore(user1, big.NewInt(10000))
	rep.IncFreeScore(user2, big.NewInt(10000))

	rep.DonateAt(user1, post1, big.NewInt(100000))
	rep.DonateAt(user1, post1, big.NewInt(1000))
	rep.DonateAt(user2, post1, big.NewInt(100000))
	rep.DonateAt(user2, post2, big.NewInt(1000))
	suite.MoveToNewRound()

	set := rep.getSeedSet(1)
	suite.Equal(1, len(set))
	suite.True(set[post1])
	suite.False(set[post2])

	rep.DonateAt(user1, post1, big.NewInt(10000))
	rep.DonateAt(user2, post2, big.NewInt(10000))
	suite.MoveToNewRound()
	set = rep.getSeedSet(2)
	suite.Equal(2, len(set))
	suite.True(set[post1])
	suite.True(set[post2])
	suite.False(set["other"])
}

func (suite *ReputationTestSuite) TestReputationDecraseToZero() {
	repData := reputationData{
		consumption: big.NewInt(1),
		hold:        big.NewInt(0),
		reputation:  big.NewInt(1),
	}
	consumptions := consumptionInfo{
		seed:    big.NewInt(0),
		other:   big.NewInt(1),
		seedIF:  big.NewInt(0),
		otherIF: big.NewInt(1),
	}

	rep := suite.rep
	newrep := rep.calcReputation(repData, consumptions)
	suite.EqualZero(newrep.reputation)
}

func (suite *ReputationTestSuite) TestCalcReputation() {
	cases := []struct {
		repData     reputationData
		consumption consumptionInfo
		repeat      int
		expected    reputationData
	}{
		{
			// growth curve
			reputationData{
				consumption: big.NewInt(1),
				hold:        big.NewInt(0),
				reputation:  big.NewInt(1),
			},
			consumptionInfo{
				seed:    big.NewInt(10000 * 100000),
				other:   big.NewInt(333 * 100000),
				seedIF:  big.NewInt(1),
				otherIF: big.NewInt(0),
			},
			5,
			reputationData{
				consumption: big.NewInt(409510000),
				hold:        big.NewInt(32804999),
				reputation:  big.NewInt(81460010),
			},
		},
		{
			// decrease curve
			reputationData{
				consumption: big.NewInt(1000 * 100000),
				hold:        big.NewInt(10 * 100000),
				reputation:  big.NewInt(900 * 100000),
			},
			consumptionInfo{
				seed:    big.NewInt(10 * 100000),
				other:   big.NewInt(600 * 100000),
				seedIF:  big.NewInt(10 * 100000),
				otherIF: big.NewInt(600 * 100000),
			},
			13,
			reputationData{
				consumption: big.NewInt(37860000),
				hold:        big.NewInt(10 *100000),
				reputation:  big.NewInt(27860000),
			},
		},
	}
	rep := suite.rep
	for i, c := range cases {
		data := c.repData
		for j := 0; j < c.repeat; j++ {
			data = rep.calcReputation(data, c.consumption)
		}
		suite.Equal(c.expected, data, "case: %d", i)
	}
}

// reputation values are greater than zero even if data does not make sense.
func (suite *ReputationTestSuite) TestReputationGTEZero() {
	repData := reputationData{
		consumption: big.NewInt(0),
		hold:        big.NewInt(0),
		reputation:  big.NewInt(0),
	}
	consumptions := consumptionInfo{
		seed:    big.NewInt(1000),
		other:   big.NewInt(10000),
		seedIF:  big.NewInt(0),
		otherIF: big.NewInt(3333),
	}

	rep := suite.rep
	newrep := rep.calcReputation(repData, consumptions)
	suite.True(newrep.reputation.Cmp(big.NewInt(0)) >= 0)
	suite.True(newrep.hold.Cmp(big.NewInt(0)) >= 0)
	suite.True(newrep.hold.Cmp(big.NewInt(0)) >= 0)
}

func (suite *ReputationTestSuite) TestReputationMigrate() {
	rep := suite.rep
	user1 := "user1"
	user2 := "user2"
	post1 := "post1"
	suite.True(rep.RequireMigrate(user1))
	suite.True(rep.RequireMigrate(user2))
	suite.MoveToNewRound()
	suite.True(rep.RequireMigrate(user1))
	suite.MoveToNewRound()
	suite.MoveToNewRound()
	suite.True(rep.RequireMigrate(user1))
	suite.MoveToNewRound()
	suite.True(rep.RequireMigrate(user1))
	rep.DonateAt(user1, post1, big.NewInt(1000))
	suite.False(rep.RequireMigrate(user1))
	suite.MoveToNewRound()
	suite.False(rep.RequireMigrate(user1))
	suite.True(rep.RequireMigrate(user2))
	rep.MigrateFromV1(user2, big.NewInt(333))
	suite.False(rep.RequireMigrate(user2))

	// double migration should be ignored.
	rep.MigrateFromV1(user2, big.NewInt(9999))

	suite.Equal(big.NewInt(333), rep.GetReputation(user2))
}

func (suite *ReputationTestSuite) TestAppendDonation() {
	rep := NewReputation(NewReputationStore(internal.NewMockStore()), 100000, 2,
		suite.roundDurationSeconds).(ReputationImpl)
	user := &userMeta{
		Reputation: big.NewInt(100),
		Unsettled:  []Donation{},
	}
	cases := []struct {
		post           Pid
		amount         LinoCoin
		expectedImpact IF
		expected       *userMeta
	}{
		{
			"p1", big.NewInt(33), big.NewInt(33),
			&userMeta{
				Reputation: big.NewInt(100),
				Unsettled: []Donation{
					Donation{Pid: "p1", Amount: big.NewInt(33), Impact: big.NewInt(33)},
				},
			},
		},
		{
			"p2", big.NewInt(77), big.NewInt(67),
			&userMeta{
				Reputation: big.NewInt(100),
				Unsettled: []Donation{
					Donation{Pid: "p1", Amount: big.NewInt(33), Impact: big.NewInt(33)},
					Donation{Pid: "p2", Amount: big.NewInt(77), Impact: big.NewInt(67)},
				},
			},
		},
		{
			"p3", big.NewInt(100), big.NewInt(0),
			&userMeta{
				Reputation: big.NewInt(100),
				Unsettled: []Donation{
					Donation{Pid: "p1", Amount: big.NewInt(33), Impact: big.NewInt(33)},
					Donation{Pid: "p2", Amount: big.NewInt(77), Impact: big.NewInt(67)},
				},
			},
		},
		{
			"p1", big.NewInt(100), big.NewInt(0),
			&userMeta{
				Reputation: big.NewInt(100),
				Unsettled: []Donation{
					Donation{Pid: "p1", Amount: big.NewInt(133), Impact: big.NewInt(33)},
					Donation{Pid: "p2", Amount: big.NewInt(77), Impact: big.NewInt(67)},
				},
			},
		},
		{
			"p2", big.NewInt(1000), big.NewInt(0),
			&userMeta{
				Reputation: big.NewInt(100),
				Unsettled: []Donation{
					Donation{Pid: "p1", Amount: big.NewInt(133), Impact: big.NewInt(33)},
					Donation{Pid: "p2", Amount: big.NewInt(1077), Impact: big.NewInt(67)},
				},
			},
		},
	}

	for i, c := range cases {
		impact := rep.appendDonation(user, c.post, c.amount)
		suite.Equal(c.expectedImpact, impact, "case: %d", i)
		suite.Equal(c.expected, user, "case: %d", i)
	}
}

// func (suite *ReputationTestSuite) TestComputeReputation() {
// 	cases := []struct{
// 		u *userMeta
// 		r RoundId
// 	}{
// 		{
// 			&userMeta{

// 			}, 3,
// 		},
// 	}
// }

func (suite *ReputationTestSuite) TestDonationReturnDp1() {
	rep := suite.rep
	user1 := "user1"
	post1 := "post1"
	post2 := "post2"

	dp1 := rep.DonateAt(user1, post1, big.NewInt(InitialReputation))
	dp2 := rep.DonateAt(user1, post1, big.NewInt(InitialReputation))
	dp3 := rep.DonateAt(user1, post2, big.NewInt(InitialReputation))
	suite.Equal(big.NewInt(InitialReputation), dp1)
	suite.Equal(big.NewInt(0), dp2)
	suite.Equal(big.NewInt(0), dp3)
}

func (suite *ReputationTestSuite) TestDonationReturnDp2() {
	rep := suite.rep
	user1 := "user1"
	user2 := "user2"
	post1 := "post1"
	post2 := "post2"

	dp1 := rep.DonateAt(user1, post1, big.NewInt(100))
	dp2 := rep.DonateAt(user1, post2, big.NewInt(100))
	dpu2 := rep.DonateAt(user2, post1, big.NewInt(100))
	suite.Equal(big.NewInt(InitialReputation), dp1)
	suite.Equal(big.NewInt(0), dp2)
	suite.Equal(big.NewInt(InitialReputation), dpu2)

	suite.MoveToNewRound()

	// round 2
	dp3 := rep.DonateAt(user1, post2, big.NewInt(3))
	dp4 := rep.DonateAt(user1, post1, big.NewInt(4))
	dp5 := rep.DonateAt(user1, post1, big.NewInt(5))
	dpu2 = rep.DonateAt(user2, post2, big.NewInt(17))
	suite.Equal(big.NewInt(3), dp3)
	suite.Equal(big.NewInt(4), dp4)
	suite.Equal(big.NewInt(2), dp5)
	suite.Equal(big.NewInt(9), dpu2)
}

// func TestDonationBasic(t *testing.T) {
// 	assert := assert.New(t)
// 	store := newReputationStoreOnMock()
// 	rep := NewTestReputationImpl(store)
// 	t1 := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
// 	t3 := time.Date(1995, time.February, 6, 12, 11, 0, 0, time.UTC)
// 	user1 := "user1"
// 	post1 := "post1"

// 	// round 2
// 	rep.Update(t1.Unix())
// 	rep.DonateAt(user1, post1, big.NewInt(100*OneLinoCoin))
// 	assert.Equal(big.NewInt(100*OneLinoCoin), rep.store.GetRoundPostSumStake(2, post1))
// 	assert.Equal(rep.GetReputation(user1), big.NewInt(InitialCustomerScore))
// 	assert.Equal(big.NewInt(OneLinoCoin), rep.store.GetRoundSumDp(2)) // bounded by this user's dp

// 	// round 3
// 	rep.Update(t3.Unix())
// 	// (1 * 9 + 100) / 10
// 	assert.Equal(big.NewInt(1090000), rep.GetReputation(user1))
// 	assert.Equal(big.NewInt(OneLinoCoin), rep.GetSumRep(post1))
// }

// // customer score is correct after multiple rounds.
// func TestDonationCase1(t *testing.T) {
// 	assert := assert.New(t)
// 	store := newReputationStoreOnMock()
// 	rep := NewTestReputationImpl(store)
// 	t1 := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
// 	t3 := time.Date(1995, time.February, 6, 12, 11, 0, 0, time.UTC)
// 	t4 := time.Date(1995, time.February, 7, 13, 11, 1, 0, time.UTC)
// 	user1 := "user1"
// 	post1 := "post1"
// 	post2 := "post2"

// 	// round 2
// 	rep.Update(t1.Unix())
// 	rep.DonateAt(user1, post1, big.NewInt(100*OneLinoCoin))
// 	assert.Equal(big.NewInt(100*OneLinoCoin), rep.store.GetRoundPostSumStake(2, post1))
// 	assert.Equal(rep.GetReputation(user1), big.NewInt(InitialCustomerScore))
// 	assert.Equal(big.NewInt(OneLinoCoin), rep.store.GetRoundSumDp(2)) // bounded by this user's dp

// 	// round 3
// 	rep.Update(t3.Unix())
// 	// (1 * 9 + 100) / 10
// 	assert.Equal(big.NewInt(1090000), rep.GetReputation(user1))
// 	assert.Equal(big.NewInt(OneLinoCoin), rep.GetSumRep(post1))
// 	rep.DonateAt(user1, post1, big.NewInt(1*OneLinoCoin)) // does not count
// 	rep.DonateAt(user1, post2, big.NewInt(900*OneLinoCoin))
// 	rep.Update(t4.Unix())
// 	// (10.9 * 9 + 900) / 10
// 	assert.Equal(big.NewInt(9981000), rep.GetReputation(user1))
// 	assert.Equal([]Pid{post2}, rep.store.GetRoundResult(3))
// 	// round 4
// }

// // multiple user split stake correct.
// func TestDonationCase2(t *testing.T) {
// 	assert := assert.New(t)
// 	store := newReputationStoreOnMock()
// 	rep := NewTestReputationImpl(store)
// 	t1 := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
// 	t3 := time.Date(1995, time.February, 6, 12, 11, 0, 0, time.UTC)
// 	t4 := time.Date(1995, time.February, 7, 13, 11, 1, 0, time.UTC)
// 	user1 := "user1"
// 	user2 := "user2"
// 	user3 := "user3"
// 	post1 := "post1"
// 	post2 := "post2"

// 	// round 2
// 	rep.Update(t1.Unix())
// 	dp1 := rep.DonateAt(user1, post1, big.NewInt(100*OneLinoCoin))
// 	dp2 := rep.DonateAt(user2, post2, big.NewInt(1000*OneLinoCoin))
// 	dp3 := rep.DonateAt(user3, post2, big.NewInt(1000*OneLinoCoin))
// 	assert.Equal(big.NewInt(OneLinoCoin), dp1)
// 	assert.Equal(big.NewInt(OneLinoCoin), dp2)
// 	assert.Equal(big.NewInt(OneLinoCoin), dp3)
// 	assert.Equal(big.NewInt(100*OneLinoCoin), rep.store.GetRoundPostSumStake(2, post1))
// 	assert.Equal(rep.GetReputation(user1), big.NewInt(InitialCustomerScore))
// 	assert.Equal(big.NewInt(3*OneLinoCoin), rep.store.GetRoundSumDp(2)) // bounded by this user's dp

// 	// post1, dp, 1
// 	// post2, dp, 2
// 	// round 3
// 	rep.Update(t3.Unix())
// 	assert.Equal([]Pid{post2, post1}, rep.store.GetRoundResult(2))
// 	assert.Equal(big.NewInt(1090000), rep.GetReputation(user1))
// 	assert.Equal(big.NewInt(13943027), rep.GetReputation(user2))
// 	assert.Equal(big.NewInt(6236972), rep.GetReputation(user3))
// 	assert.Equal(big.NewInt(OneLinoCoin), rep.GetSumRep(post1))
// 	assert.Equal(big.NewInt(2*OneLinoCoin), rep.GetSumRep(post2))

// 	// user1: 10.9
// 	// user2: 139.43027
// 	// user3: 62.36972
// 	dp1 = rep.DonateAt(user2, post2, big.NewInt(200*OneLinoCoin))
// 	dp2 = rep.DonateAt(user1, post1, big.NewInt(400*OneLinoCoin))
// 	// does not count because rep used up.
// 	dp3 = rep.DonateAt(user1, post1, big.NewInt(900*OneLinoCoin))
// 	dp4 := rep.DonateAt(user3, post1, big.NewInt(500*OneLinoCoin))
// 	assert.Equal(big.NewInt(13943027-OneLinoCoin), dp1)
// 	assert.Equal(big.NewInt(1090000-OneLinoCoin), dp2)
// 	assert.Equal(BigIntZero, dp3)
// 	assert.Equal(big.NewInt(6236972), dp4)

// 	// round 4
// 	rep.Update(t4.Unix())
// 	assert.Equal([]Pid{post2, post1}, rep.store.GetRoundResult(3))
// 	assert.Equal(big.NewInt(16136841), rep.GetReputation(user1))
// 	assert.Equal(big.NewInt(14548724), rep.GetReputation(user2))
// 	assert.Equal(big.NewInt(8457432), rep.GetReputation(user3))
// }
// func TestStartNewRound(t *testing.T) {
// 	assert := assert.New(t)
// 	store := newReputationStoreOnMock()

// 	assert.Equal(RoundId(1), store.GetCurrentRound())
// 	store.StartNewRound(222)
// 	assert.Equal(RoundId(2), store.GetCurrentRound())
// 	assert.Equal(int64(222), store.GetRoundStartAt(2))
// 	assert.Empty(store.GetRoundTopNPosts(2))
// 	assert.Empty(store.GetRoundResult(2))
// 	assert.Equal(big.NewInt(0), store.GetRoundSumDp(1))
// 	assert.Equal(big.NewInt(0), store.GetRoundSumDp(2))
// }

// func TestTopN(t *testing.T) {
// 	assert := assert.New(t)
// 	post1 := "bla"
// 	post2 := "zzz"
// 	store := newReputationStoreOnMock()

// 	// test sorting
// 	store.StartNewRound(222)
// 	store.SetRoundPostSumDp(2, post1, big.NewInt(100))
// 	assert.Equal([]PostDpPair{{post1, big.NewInt(100)}}, store.GetRoundTopNPosts(2))
// 	store.SetRoundPostSumDp(2, post2, big.NewInt(300))
// 	assert.Equal([]PostDpPair{{post2, big.NewInt(300)}, {post1, big.NewInt(100)}}, store.GetRoundTopNPosts(2))
// 	store.SetRoundPostSumDp(2, post1, big.NewInt(1000))
// 	assert.Equal([]PostDpPair{{post1, big.NewInt(1000)}, {post2, big.NewInt(300)}}, store.GetRoundTopNPosts(2))

// 	for i := 1; i <= DefaultBestContentIndexN; i++ {
// 		store.SetRoundPostSumDp(2, "p"+string(i), big.NewInt(int64(i)))
// 	}
// 	for i := DefaultBestContentIndexN; i >= 0; i-- {
// 		store.SetRoundPostSumDp(2, "pp"+string(i), big.NewInt(int64(i)))
// 	}
// 	topN := store.GetRoundTopNPosts(2)

// 	// at most N.
// 	assert.Equal(DefaultBestContentIndexN, len(topN))
// 	// decreasing order
// 	for i, v := range topN {
// 		if i > 0 {
// 			assert.Truef(v.SumDp.Cmp(topN[i-1].SumDp) <= 0, "%+v, %+v", v, topN[i-1])
// 		}
// 	}
// }

func (suite *ReputationTestSuite) TestBigIntEMA() {
	cases := []struct {
		prev     bigInt
		new      bigInt
		w        int64
		expected bigInt
	}{
		{
			prev:     big.NewInt(333),
			new:      big.NewInt(333),
			w:        10,
			expected: big.NewInt(333),
		},
		{
			prev:     big.NewInt(0),
			new:      big.NewInt(10),
			w:        10,
			expected: big.NewInt(1),
		},
		{
			prev:     big.NewInt(10),
			new:      big.NewInt(110),
			w:        10,
			expected: big.NewInt(20),
		},
		{
			prev:     big.NewInt(4),
			new:      big.NewInt(77),
			w:        7,
			expected: big.NewInt(14),
		},
	}

	for i, v := range cases {
		suite.Equal(v.expected, bigIntEMA(v.prev, v.new, v.w), "case: %d", i)
	}
}

func (suite *ReputationTestSuite) TestIntDivFrac() {
	cases := []struct {
		v        bigInt
		num      int64
		denum    int64
		expected bigInt
	}{
		{
			v:        big.NewInt(80),
			num:      8,
			denum:    10,
			expected: big.NewInt(100),
		},
		{
			v:        big.NewInt(100),
			num:      1,
			denum:    3,
			expected: big.NewInt(300),
		},
		{
			v:        big.NewInt(77),
			num:      11,
			denum:    7,
			expected: big.NewInt(49),
		},
	}

	for i, v := range cases {
		suite.Equal(v.expected, bigIntDivFrac(v.v, v.num, v.denum), "case: %d", i)
	}
}

func (suite *ReputationTestSuite) TestIntMulFrac() {
	cases := []struct {
		v        bigInt
		num      int64
		denum    int64
		expected bigInt
	}{
		{
			v:        big.NewInt(80),
			num:      8,
			denum:    10,
			expected: big.NewInt(64),
		},
		{
			v:        big.NewInt(100),
			num:      1,
			denum:    3,
			expected: big.NewInt(33),
		},
		{
			v:        big.NewInt(77),
			num:      11,
			denum:    7,
			expected: big.NewInt(121),
		},
	}

	for i, v := range cases {
		suite.Equal(v.expected, bigIntMulFrac(v.v, v.num, v.denum), "case: %d", i)
	}
}

func (suite *ReputationTestSuite) TestBubbleUp() {
	cases := []struct {
		posts    []PostIFPair
		pos      int
		expected []PostIFPair
	}{
		{
			posts:    nil,
			pos:      -1,
			expected: nil,
		},
		{
			posts:    []PostIFPair{{"1", big.NewInt(3)}},
			pos:      0,
			expected: []PostIFPair{{"1", big.NewInt(3)}},
		},
		{
			posts:    []PostIFPair{{"1", big.NewInt(3)}, {"2", big.NewInt(0)}},
			pos:      0,
			expected: []PostIFPair{{"1", big.NewInt(3)}, {"2", big.NewInt(0)}},
		},
		{
			posts:    []PostIFPair{{"1", big.NewInt(3)}, {"2", big.NewInt(4)}},
			pos:      1,
			expected: []PostIFPair{{"2", big.NewInt(4)}, {"1", big.NewInt(3)}},
		},
		{
			posts: []PostIFPair{
				{"1", big.NewInt(9)},
				{"3", big.NewInt(8)},
				{"5", big.NewInt(7)},
				{"2", big.NewInt(6)},
				{"8", big.NewInt(100)},
				{"0", big.NewInt(5)},
				{"11", big.NewInt(4)},
			},
			pos: 4,
			expected: []PostIFPair{
				{"8", big.NewInt(100)},
				{"1", big.NewInt(9)},
				{"3", big.NewInt(8)},
				{"5", big.NewInt(7)},
				{"2", big.NewInt(6)},
				{"0", big.NewInt(5)},
				{"11", big.NewInt(4)},
			},
		},
		{
			posts: []PostIFPair{
				{"1", big.NewInt(9)},
				{"3", big.NewInt(8)},
				{"5", big.NewInt(7)},
				{"2", big.NewInt(6)},
				{"0", big.NewInt(5)},
				{"11", big.NewInt(4)},
				{"8", big.NewInt(100)},
			},
			pos: 6,
			expected: []PostIFPair{
				{"8", big.NewInt(100)},
				{"1", big.NewInt(9)},
				{"3", big.NewInt(8)},
				{"5", big.NewInt(7)},
				{"2", big.NewInt(6)},
				{"0", big.NewInt(5)},
				{"11", big.NewInt(4)},
			},
		},
	}
	for i, v := range cases {
		bubbleUp(v.posts, v.pos)
		suite.Equal(v.expected, v.posts, "case: %d", i)
	}
}
