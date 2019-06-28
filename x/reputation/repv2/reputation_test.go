package repv2

import (
	"fmt"
	"io/ioutil"
	"math/big"
	// "math/rand"
	"os"
	"path/filepath"
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
	suite.store = NewReputationStore(internal.NewMockStore(), DefaultInitialReputation)
	suite.rep = NewReputation(
		suite.store, suite.bestN, suite.userMaxN,
		DefaultRoundDurationSeconds, DefaultSampleWindowSize, DefaultDecayFactor,
	).(ReputationImpl)
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

func (suite *ReputationTestSuite) TestUpdateTime() {
	suite.MoveToNewRound()
	t := suite.time
	rep := suite.rep
	r, rt := rep.GetCurrentRound()
	suite.Equal(int64(2), r)
	suite.Equal(t.Unix(), rt)

	t2 := t.Add(time.Duration(10*3600) * time.Second)
	rep.Update(t2.Unix())
	r, rt = rep.GetCurrentRound()
	suite.Equal(int64(2), r)
	suite.Equal(t.Unix(), rt)

	t3 := t2.Add(time.Duration(16*3600) * time.Second)
	rep.Update(t3.Unix())
	r, rt = rep.GetCurrentRound()
	suite.Equal(int64(3), r)
	suite.Equal(t3.Unix(), rt)

	t4 := t3.Add(time.Duration(10*3600) * time.Second)
	rep.Update(t4.Unix())
	r, rt = rep.GetCurrentRound()
	suite.Equal(int64(3), r)
	suite.Equal(t3.Unix(), rt)
}

func (suite *ReputationTestSuite) TestExportImportFile() {
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

	rep.DonateAt(user1, post1, big.NewInt(3333))
	rep.DonateAt(user1, post1, big.NewInt(4444))
	rep.DonateAt(user2, post1, big.NewInt(5555))
	rep.DonateAt(user2, post2, big.NewInt(1324))
	suite.MoveToNewRound()
	suite.MoveToNewRound()

	suite.Require().Equal(big.NewInt(10920), suite.rep.GetReputation("user1"))
	suite.Require().Equal(big.NewInt(10910), suite.rep.GetReputation("user2"))

	dir, err := ioutil.TempDir("", "example")
	suite.Require().Nil(err)
	defer os.RemoveAll(dir) // clean up

	tmpfn := filepath.Join(dir, "tmpfile")
	suite.MoveToNewRound()
	rep.ExportToFile(tmpfn)

	imported := NewReputation(
		NewReputationStore(internal.NewMockStore(), DefaultInitialReputation),
		suite.bestN, suite.userMaxN,
		DefaultRoundDurationSeconds, DefaultSampleWindowSize, DefaultDecayFactor,
	).(ReputationImpl)

	imported.ImportFromFile(tmpfn)
	suite.Equal(suite.rep.GetReputation("user1"), imported.GetReputation("user1"))
	suite.Equal(suite.rep.GetReputation("user2"), imported.GetReputation("user2"))
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

func (suite *ReputationTestSuite) TestComputeReputation() {
	rep := suite.rep
	suite.EqualZero(rep.computeReputation(big.NewInt(0), big.NewInt(1000)))
	suite.Equal(big.NewInt(969), rep.computeReputation(big.NewInt(999), big.NewInt(3)))
}

func (suite *ReputationTestSuite) TestComputeNewRepData() {
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
				hold:        big.NewInt(10 * 100000),
				reputation:  big.NewInt(27860000),
			},
		},
	}
	rep := suite.rep
	for i, c := range cases {
		data := c.repData
		for j := 0; j < c.repeat; j++ {
			data = rep.computeNewRepData(data, c.consumption)
		}
		suite.Equal(c.expected, data, "case: %d", i)
	}
}

func (suite *ReputationTestSuite) TestComputeNewRepDataDecreaseToZero() {
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
	newrep := rep.computeNewRepData(repData, consumptions)
	suite.EqualZero(newrep.reputation)
}

// reputation values are greater than zero even if data does not make sense.
func (suite *ReputationTestSuite) TestComputeNewRepDataGTEZero() {
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
	newrep := rep.computeNewRepData(repData, consumptions)
	suite.True(newrep.reputation.Cmp(big.NewInt(0)) >= 0)
	suite.True(newrep.hold.Cmp(big.NewInt(0)) >= 0)
	suite.True(newrep.hold.Cmp(big.NewInt(0)) >= 0)
}

func (suite *ReputationTestSuite) TestDonateAtGrow1() {
	rep := suite.rep
	for i := 0; i <= 50; i++ {
		rep.DonateAt("user1", "post1", big.NewInt(100*100000))
		suite.MoveToNewRound()
	}
	suite.Equal(big.NewInt(9690802), rep.GetReputation("user1"))
}

func (suite *ReputationTestSuite) TestDonateAtGrowAndDown() {
	rep := suite.rep
	for i := 0; i <= 60; i++ {
		rep.DonateAt("user1", "post1", big.NewInt(80*100000))
		rep.DonateAt("user1", "post2", big.NewInt(20*100000))
		suite.MoveToNewRound()
	}
	suite.Equal(big.NewInt(9270170), rep.GetReputation("user1"))

	rep.IncFreeScore("majority", big.NewInt(1000000*100000))
	for i := 0; i <= 1; i++ {
		rep.DonateAt("user1", "trash", big.NewInt(1*100000))
		rep.DonateAt("majority", "good", big.NewInt(1000000*100000))
		suite.MoveToNewRound()
	}
	suite.Equal(big.NewInt(9254170), rep.GetReputation("user1"))

	for i := 0; i <= 60; i++ {
		// rep.DonateAt("user1", "good", big.NewInt(50 * 100000))
		rep.DonateAt("user1", "trash", big.NewInt(100*100000))
		rep.DonateAt("majority", "good", big.NewInt(1000000*100000))
		suite.MoveToNewRound()
	}
	suite.Equal(big.NewInt(57205), rep.GetReputation("user1"))
}

func (suite *ReputationTestSuite) TestUpdateReputationDonateAt() {
	rep := suite.rep

	// panics
	suite.Panics(func() { rep.DonateAt("", "123", big.NewInt(11)) })
	suite.Panics(func() { rep.DonateAt("u31", "", big.NewInt(11)) })
	suite.Panics(func() { rep.DonateAt("", "", big.NewInt(11)) })

	donations := []struct {
		from   Uid
		to     Pid
		amount int64
	}{
		{"user1", "post1", 10000},
		{"user2", "post1", 3},
		{"user3", "post1", 600},
		{"user4", "post1", 999},
		{"user5", "post1", 1},
		{"user6", "post1", 2},
		{"user7", "post2", 7777},
		{"user8", "post2", 2},
		{"user9", "post2", 2},
		{"user10", "post2", 100},
		{"user11", "post2", 1000000},
	}
	cases := []struct {
		user     string
		expected *userMeta
	}{
		{
			"user1",
			&userMeta{
				Consumption:       big.NewInt(1000),
				Hold:              big.NewInt(99),
				Reputation:        big.NewInt(10),
				LastDonationRound: 1,
				LastSettledRound:  1,
			},
		},
		{
			"user3",
			&userMeta{
				Consumption:       big.NewInt(60),
				Hold:              big.NewInt(5),
				Reputation:        big.NewInt(10),
				LastDonationRound: 1,
				LastSettledRound:  1,
			},
		},
		{
			"user7",
			&userMeta{
				Consumption:       big.NewInt(778),
				Hold:              big.NewInt(77),
				Reputation:        big.NewInt(8),
				LastDonationRound: 1,
				LastSettledRound:  1,
			},
		},
		{
			"user11",
			&userMeta{
				Consumption:       big.NewInt(100000),
				Hold:              big.NewInt(9999),
				Reputation:        big.NewInt(10),
				LastDonationRound: 1,
				LastSettledRound:  1,
			},
		},
	}
	for _, donation := range donations {
		rep.DonateAt(donation.from, donation.to, big.NewInt(donation.amount))
	}
	suite.MoveToNewRound()
	for i, v := range cases {
		user := rep.store.GetUserMeta(v.user)
		rep.updateReputation(user, 2)
		suite.Equal(v.expected, user, "case: %d", i)
	}
}

func (suite *ReputationTestSuite) TestAppendDonation() {
	rep := NewReputation(
		NewReputationStore(internal.NewMockStore(), DefaultInitialReputation),
		100000, 2,
		DefaultRoundDurationSeconds, DefaultSampleWindowSize, DefaultDecayFactor).(ReputationImpl)
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

func (suite *ReputationTestSuite) TestIncRoundPostSumImpactAndUpdate() {
	rep := suite.rep
	for i := 0; i < 1000; i++ {
		post := fmt.Sprintf("post%d", i)
		rep.incRoundPostSumImpact(1, post, big.NewInt(int64(1)))
	}
	for i := 999; i >= 0; i-- {
		post := fmt.Sprintf("post%d", 999-i)
		rep.incRoundPostSumImpact(1, post, big.NewInt(int64(3)))
	}
	for i := 0; i < 1000; i++ {
		post := fmt.Sprintf("post%d", i)
		if i%173 == 0 {
			rep.incRoundPostSumImpact(1, post, big.NewInt(10000))
		}
	}

	for i := 0; i < 1000; i++ {
		post := fmt.Sprintf("post%d", i)
		meta := rep.store.GetRoundPostMeta(1, post)
		if i%173 == 0 {
			suite.Equal(big.NewInt(10004), meta.SumIF)
		} else {
			suite.Equal(big.NewInt(4), meta.SumIF)
		}
	}

	round := rep.store.GetRoundMeta(1)
	suite.Nil(round.Result)
	suite.Equal(big.NewInt(4*1000+10000*6), round.SumIF)
	suite.Equal(int64(0), round.StartAt)
	suite.Equal(30, len(round.TopN))
	for _, v := range round.TopN {
		id := -1
		fmt.Sscanf(v.Pid, "post%d", &id)
		if id%173 == 0 {
			suite.Equal(big.NewInt(10004), v.SumIF)
		} else {
			suite.Equal(big.NewInt(4), v.SumIF)
		}
	}

	suite.MoveToNewRound()
	roundFinal := rep.store.GetRoundMeta(1)
	suite.Equal(6, len(roundFinal.Result))
	for _, v := range round.Result {
		id := -1
		fmt.Sscanf(v, "post%d", &id)
		suite.True(id%173 == 0)
	}
}

func (suite *ReputationTestSuite) TestFirstBlock1() {
	rep := suite.rep
	newBlockTime := int64(0)
	rep.Update(0)
	rid, startAt := rep.GetCurrentRound()
	suite.Equal(int64(1), rid)
	suite.Equal(newBlockTime, startAt)
	suite.Equal(rep.GetReputation("me"), big.NewInt(DefaultInitialReputation))
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
	suite.Equal(big.NewInt(3000+DefaultInitialReputation), rep.GetReputation("user1"))
}

func (suite *ReputationTestSuite) TestDonationReturnDp1() {
	rep := suite.rep
	user1 := "user1"
	post1 := "post1"
	post2 := "post2"

	dp1 := rep.DonateAt(user1, post1, big.NewInt(DefaultInitialReputation))
	dp2 := rep.DonateAt(user1, post1, big.NewInt(DefaultInitialReputation))
	dp3 := rep.DonateAt(user1, post2, big.NewInt(DefaultInitialReputation))
	suite.Equal(big.NewInt(DefaultInitialReputation), dp1)
	suite.Equal(big.NewInt(0), dp2)
	suite.Equal(big.NewInt(0), dp3)
}

func (suite *ReputationTestSuite) TestDonationReturnDp2() {
	rep := suite.rep
	user1 := "user1"
	user2 := "user2"
	post1 := "post1"
	post2 := "post2"

	dp1 := rep.DonateAt(user1, post1, big.NewInt(10000))
	dp2 := rep.DonateAt(user1, post2, big.NewInt(10000))
	dpu2 := rep.DonateAt(user2, post1, big.NewInt(10000))
	suite.Equal(big.NewInt(DefaultInitialReputation), dp1)
	suite.Equal(big.NewInt(0), dp2)
	suite.Equal(big.NewInt(DefaultInitialReputation), dpu2)

	suite.MoveToNewRound()

	// round 2
	dp3 := rep.DonateAt(user1, post2, big.NewInt(3))
	dp4 := rep.DonateAt(user1, post1, big.NewInt(4))
	dp5 := rep.DonateAt(user1, post1, big.NewInt(5))
	dpu2 = rep.DonateAt(user2, post2, big.NewInt(17))
	suite.Equal(big.NewInt(3), dp3)
	suite.Equal(big.NewInt(4), dp4)
	suite.Equal(big.NewInt(3), dp5)
	suite.Equal(big.NewInt(10), dpu2)
}

func (suite *ReputationTestSuite) TestBigIntEMA() {
	suite.Panics(func() { bigIntEMA(big.NewInt(1000), big.NewInt(333), 0) })
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
	suite.Panics(func() { bigIntDivFrac(big.NewInt(1000), 1, 0) })
	suite.Panics(func() { bigIntDivFrac(big.NewInt(1000), 0, 1) })
	suite.Panics(func() { bigIntDivFrac(big.NewInt(1000), 0, 0) })
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
	suite.Panics(func() { bigIntMulFrac(big.NewInt(1000), 1, 0) })
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

// func (suite *ReputationTestSuite) simPostZipf(nposts uint64) *rand.Zipf {
// 	// zipf posts, with s = 2, v = 50. number of seed: 193 if nposts = 10000
// 	zipf := rand.NewZipf(rand.New(rand.NewSource(121212)), 2, 50, uint64(nposts))
// 	return zipf
// 	// print distribution.
// 	// count := make(map[uint64]int)
// 	// for i := 0; i < nposts; i++ {
// 	// 	v := zipf.Uint64()
// 	// 	count[v]++
// 	// }
// 	// probs := make([]float64, nposts)
// 	// for k, v := range count {
// 	// 	probs[k] = float64(v) * 100 / float64(nposts)
// 	// }
// 	// total := float64(0.0)
// 	// for i, v := range probs {
// 	// 	total += v
// 	// 	if total >= 80 {
// 	// 		fmt.Printf("80: %d\n", i)
// 	// 		break
// 	// 	}
// 	// }
// 	// fmt.Println(probs)
// }

// simulations
// func (suite *ReputationTestSuite) TestSimulation() {
// 	rep := NewReputation(
// 		NewReputationStore(internal.NewMockStore(), DefaultInitialReputation),
// 		200, 30, DefaultRoundDurationSeconds, DefaultSampleWindowSize, DefaultDecayFactor)
// 	suite.rep = rep.(ReputationImpl)
// 	// zipf posts.
// 	nPosts := uint64(1000)
// 	zipf := suite.simPostZipf(nPosts)
// 	nUsers := int(5 * nPosts)
// 	toUID := func(i int) Uid {
// 		return fmt.Sprintf("user%d", i)
// 	}
// 	toPID := func(i uint64) Pid {
// 		return fmt.Sprintf("post%d", i)
// 	}

// 	for j := 0; j < 3; j++ {
// 		for i := 0; i < nUsers; i++ {
// 			rep.DonateAt(toUID(i), toPID(zipf.Uint64()), big.NewInt(10*100000))
// 		}
// 		suite.MoveToNewRound()
// 		fmt.Println(j)
// 	}

// 	for i := 0; i < nUsers; i++ {
// 		fmt.Println(rep.GetReputation(toUID(i)))
// 	}

// }

// benchmarks
func BenchmarkDonateAt1(b *testing.B) {
	suite := ReputationTestSuite{}
	suite.SetupTest()
	for n := 0; n < b.N; n++ {
		suite.rep.DonateAt("user1", "post2", big.NewInt(100*100000))
	}
}

func BenchmarkDonateAtWorstCase(b *testing.B) {
	suite := ReputationTestSuite{}
	suite.SetupTest()

	posts := make([]Pid, b.N)
	for i := 0; i < b.N; i++ {
		posts[i] = fmt.Sprintf("post%d", i)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		suite.rep.DonateAt("user1", posts[n], big.NewInt(10000*100000))
	}
}
