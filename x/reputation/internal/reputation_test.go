package internal

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tendermint/libs/db"
)

// A mock store for testing.
type mockStore struct {
	store map[string]([]byte)
}

func newMockStore() Store {
	return &mockStore{store: make(map[string]([]byte))}
}

func (s *mockStore) Set(key []byte, val []byte) {
	// fmt.Printf("---> write %v\n", key)
	s.store[string(key)] = val
}

func (s *mockStore) Get(key []byte) []byte {
	// fmt.Printf("---> read %v\n", key)
	if el, ok := s.store[string(key)]; ok {
		return el
	} else {
		return nil
	}
}

func (s *mockStore) Has(key []byte) bool {
	_, ok := s.store[string(key)]
	return ok
}

func (s *mockStore) Delete(key []byte) {
	delete(s.store, string(key))
}

func (s *mockStore) Iterator(start, end []byte) db.Iterator {
	panic("not implemented")
}

var _ Store = &mockStore{}

func TestMockStore(t *testing.T) {
	var store Store = newMockStore()
	store.Set([]byte("a"), []byte("1"))
	if e := store.Get([]byte("a")); e != nil && string(e) == "1" {
		t.Log("store ok")
	} else {
		t.Log(e)
		t.Error("store failed")
	}
}

// use this if you want to test internal reputationImpl method.
func NewTestReputationImpl(s ReputationStore) *ReputationImpl {
	return &ReputationImpl{store: s}
}

func TestFirstBlock1(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewReputation(store)
	newBlockTime := int64(0)
	rep.Update(0)
	rid, startAt := rep.GetCurrentRound()
	assert.Equal(int64(1), rid)
	assert.Equal(newBlockTime, startAt)
	rep.IncFreeScore("me", big.NewInt(100))
	assert.Equal(rep.GetReputation("me"), big.NewInt(OneLinoCoin+100))
}

func TestFirstBlock2(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewReputation(store)
	newBlockTime := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
	rep.Update(newBlockTime.Unix())
	rid, startAt := rep.GetCurrentRound()
	assert.Equal(int64(2), rid)
	assert.Equal(newBlockTime.Unix(), startAt)

	// next round starts after RoundDuration, assuming RoundDuration is 25
	// need to update this test when you change it.
	nextBlockTime := time.Date(1995, time.February, 6, 12, 11, 0, 0, time.UTC)
	rep.Update(nextBlockTime.Unix())
	rid, startAt = rep.GetCurrentRound()
	assert.Equal(int64(3), rid)
	assert.Equal(nextBlockTime.Unix(), startAt)
}

func TestIncFreeScore(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewReputation(store)

	rep.IncFreeScore("user1", big.NewInt(OneLinoCoin))
	assert.Equal(big.NewInt(2*OneLinoCoin), rep.GetReputation("user1"))
}

func TestNumKeysCanBuy(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)

	nCanBuy := rep.numKeysCanBuy(big.NewInt(0), big.NewInt(1))
	assert.Equal(big.NewInt(0), nCanBuy)

	// 0.01 lino
	nCanBuy = rep.numKeysCanBuy(big.NewInt(0), big.NewInt(0.01*OneLinoCoin))
	assert.Equal(big.NewInt(1), nCanBuy)

	// 1 lino
	nCanBuy = rep.numKeysCanBuy(big.NewInt(0), big.NewInt(1*OneLinoCoin))
	assert.Equal(big.NewInt(95), nCanBuy)

	// 1000 lino
	nCanBuy = rep.numKeysCanBuy(big.NewInt(0), big.NewInt(1000*OneLinoCoin))
	assert.Equal(big.NewInt(13177), nCanBuy)

	// 3 times
	a1 := rep.numKeysCanBuy(big.NewInt(0), big.NewInt(1234*OneLinoCoin))
	a2 := rep.numKeysCanBuy(a1, big.NewInt(1777*OneLinoCoin))
	a3 := rep.numKeysCanBuy(bigIntAdd(a1, a2), big.NewInt(8163*OneLinoCoin))
	assert.Equal(big.NewInt(14742), a1)
	assert.Equal(big.NewInt(8818), a2)
	assert.Equal(big.NewInt(22724), a3)
}

func BenchmarkNumKeysCanBuy(b *testing.B) {
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	for n := 0; n < b.N; n++ {
		rep.numKeysCanBuy(big.NewInt(0), big.NewInt(100*OneLinoCoin))
	}
}

func TestBuyKey(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	user1 := "user1"
	post1 := "post1"

	rep.buyKey(0, user1, post1, big.NewInt(OneLinoCoin))
	assert.Equal(big.NewInt(95), rep.store.GetRoundNumKeysHas(0, user1, post1))
	assert.Equal(big.NewInt(95), rep.store.GetRoundNumKeysSold(0, post1))
}

func TestReportAt(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	user1 := "user1"
	post1 := "post1"

	rst := rep.ReportAt(user1, post1)
	assert.Equal(big.NewInt(-OneLinoCoin), rep.GetSumRep(post1))
	assert.Equal(big.NewInt(-OneLinoCoin), rst)
}

func TestDonationNoLessThanInit(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	t1 := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
	t3 := time.Date(1995, time.February, 7, 11, 11, 0, 0, time.UTC)
	user1 := "user1"
	post1 := "post1"
	rep.Update(t1.Unix())

	// round 2 start
	rep.DonateAt(user1, post1, big.NewInt(1000))
	assert.Equal(rep.GetReputation(user1), big.NewInt(InitialCustomerScore))
	rep.Update(t3.Unix())
	assert.Equal(big.NewInt(InitialCustomerScore), rep.GetReputation(user1))
}

func TestDonationZeroStake(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	t1 := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
	t2 := time.Date(1995, time.February, 7, 11, 11, 0, 0, time.UTC)
	user1 := "user1"
	post1 := "post1"
	rep.Update(t1.Unix())

	// round 2 start
	rep.DonateAt(user1, post1, big.NewInt(0))
	assert.Equal(big.NewInt(InitialCustomerScore), rep.GetReputation(user1))
	rep.Update(t2.Unix())
	assert.Equal(big.NewInt(InitialCustomerScore), rep.GetReputation(user1))
}

func TestDonationReturnDp1(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	user1 := "user1"
	post1 := "post1"
	post2 := "post2"

	// round 2 start
	dp1 := rep.DonateAt(user1, post1, big.NewInt(OneLinoCoin))
	dp2 := rep.DonateAt(user1, post1, big.NewInt(OneLinoCoin))
	dp3 := rep.DonateAt(user1, post2, big.NewInt(OneLinoCoin))
	assert.Equal(big.NewInt(OneLinoCoin), dp1)
	assert.Equal(big.NewInt(0), dp2)
	assert.Equal(big.NewInt(OneLinoCoin), dp3)
}

func TestDonationReturnDp2(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	t1 := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
	t2 := time.Date(1995, time.February, 7, 11, 11, 0, 0, time.UTC)
	user1 := "user1"
	user2 := "user2"
	post1 := "post1"
	post2 := "post2"

	// round 2
	rep.Update(t1.Unix())
	rep.DonateAt(user1, post1, big.NewInt(OneLinoCoin))
	rep.DonateAt(user2, post2, big.NewInt(100*OneLinoCoin)) // make user1's donation useless.

	// round3
	rep.Update(t2.Unix())
	dp := rep.DonateAt(user1, post1, big.NewInt(OneLinoCoin))
	assert.Equal(BigIntZero, dp)
}

func TestDonationBasic(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	t1 := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
	t3 := time.Date(1995, time.February, 6, 12, 11, 0, 0, time.UTC)
	user1 := "user1"
	post1 := "post1"

	// round 2
	rep.Update(t1.Unix())
	rep.DonateAt(user1, post1, big.NewInt(100*OneLinoCoin))
	assert.Equal(big.NewInt(100*OneLinoCoin), rep.store.GetRoundPostSumStake(2, post1))
	assert.Equal(rep.GetReputation(user1), big.NewInt(InitialCustomerScore))
	assert.Equal(big.NewInt(OneLinoCoin), rep.store.GetRoundSumDp(2)) // bounded by this user's dp

	// round 3
	rep.Update(t3.Unix())
	// (1 * 9 + 100) / 10
	assert.Equal(big.NewInt(1090000), rep.GetReputation(user1))
	assert.Equal(big.NewInt(OneLinoCoin), rep.GetSumRep(post1))
}

// customer score is correct after multiple rounds.
func TestDonationCase1(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	t1 := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
	t3 := time.Date(1995, time.February, 6, 12, 11, 0, 0, time.UTC)
	t4 := time.Date(1995, time.February, 7, 13, 11, 1, 0, time.UTC)
	user1 := "user1"
	post1 := "post1"
	post2 := "post2"

	// round 2
	rep.Update(t1.Unix())
	rep.DonateAt(user1, post1, big.NewInt(100*OneLinoCoin))
	assert.Equal(big.NewInt(100*OneLinoCoin), rep.store.GetRoundPostSumStake(2, post1))
	assert.Equal(rep.GetReputation(user1), big.NewInt(InitialCustomerScore))
	assert.Equal(big.NewInt(OneLinoCoin), rep.store.GetRoundSumDp(2)) // bounded by this user's dp

	// round 3
	rep.Update(t3.Unix())
	// (1 * 9 + 100) / 10
	assert.Equal(big.NewInt(1090000), rep.GetReputation(user1))
	assert.Equal(big.NewInt(OneLinoCoin), rep.GetSumRep(post1))
	rep.DonateAt(user1, post1, big.NewInt(1*OneLinoCoin)) // does not count
	rep.DonateAt(user1, post2, big.NewInt(900*OneLinoCoin))
	rep.Update(t4.Unix())
	// (10.9 * 9 + 900) / 10
	assert.Equal(big.NewInt(9981000), rep.GetReputation(user1))
	assert.Equal([]Pid{post2}, rep.store.GetRoundResult(3))
	// round 4
}

// multiple user split stake correct.
func TestDonationCase2(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()
	rep := NewTestReputationImpl(store)
	t1 := time.Date(1995, time.February, 5, 11, 11, 0, 0, time.UTC)
	t3 := time.Date(1995, time.February, 6, 12, 11, 0, 0, time.UTC)
	t4 := time.Date(1995, time.February, 7, 13, 11, 1, 0, time.UTC)
	user1 := "user1"
	user2 := "user2"
	user3 := "user3"
	post1 := "post1"
	post2 := "post2"

	// round 2
	rep.Update(t1.Unix())
	dp1 := rep.DonateAt(user1, post1, big.NewInt(100*OneLinoCoin))
	dp2 := rep.DonateAt(user2, post2, big.NewInt(1000*OneLinoCoin))
	dp3 := rep.DonateAt(user3, post2, big.NewInt(1000*OneLinoCoin))
	assert.Equal(big.NewInt(OneLinoCoin), dp1)
	assert.Equal(big.NewInt(OneLinoCoin), dp2)
	assert.Equal(big.NewInt(OneLinoCoin), dp3)
	assert.Equal(big.NewInt(100*OneLinoCoin), rep.store.GetRoundPostSumStake(2, post1))
	assert.Equal(rep.GetReputation(user1), big.NewInt(InitialCustomerScore))
	assert.Equal(big.NewInt(3*OneLinoCoin), rep.store.GetRoundSumDp(2)) // bounded by this user's dp

	// post1, dp, 1
	// post2, dp, 2
	// round 3
	rep.Update(t3.Unix())
	assert.Equal([]Pid{post2, post1}, rep.store.GetRoundResult(2))
	assert.Equal(big.NewInt(1090000), rep.GetReputation(user1))
	assert.Equal(big.NewInt(13943027), rep.GetReputation(user2))
	assert.Equal(big.NewInt(6236972), rep.GetReputation(user3))
	assert.Equal(big.NewInt(OneLinoCoin), rep.GetSumRep(post1))
	assert.Equal(big.NewInt(2*OneLinoCoin), rep.GetSumRep(post2))

	// user1: 10.9
	// user2: 139.43027
	// user3: 62.36972
	dp1 = rep.DonateAt(user2, post2, big.NewInt(200*OneLinoCoin))
	dp2 = rep.DonateAt(user1, post1, big.NewInt(400*OneLinoCoin))
	// does not count because rep used up.
	dp3 = rep.DonateAt(user1, post1, big.NewInt(900*OneLinoCoin))
	dp4 := rep.DonateAt(user3, post1, big.NewInt(500*OneLinoCoin))
	assert.Equal(big.NewInt(13943027-OneLinoCoin), dp1)
	assert.Equal(big.NewInt(1090000-OneLinoCoin), dp2)
	assert.Equal(BigIntZero, dp3)
	assert.Equal(big.NewInt(6236972), dp4)

	// round 4
	rep.Update(t4.Unix())
	assert.Equal([]Pid{post2, post1}, rep.store.GetRoundResult(3))
	assert.Equal(big.NewInt(16136841), rep.GetReputation(user1))
	assert.Equal(big.NewInt(14548724), rep.GetReputation(user2))
	assert.Equal(big.NewInt(8457432), rep.GetReputation(user3))
}
