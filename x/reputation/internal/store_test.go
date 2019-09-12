package internal

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newReputationStoreOnMock() ReputationStore {
	return NewReputationStoreDefaultN(newMockStore())
}

func TestInitValues(t *testing.T) {
	assert := assert.New(t)
	user1 := "test"
	post1 := "bla"
	store := newReputationStoreOnMock()

	assert.Equal(RoundId(1), store.GetCurrentRound())
	assert.Equal(big.NewInt(InitialCustomerScore), store.GetCustomerScore(user1))
	assert.Equal(big.NewInt(0), store.GetFreeScore(user1))
	assert.Equal(RoundId(0), store.GetUserLastSettled(user1))
	assert.Equal(RoundId(0), store.GetUserLastDonationRound(user1))
	assert.Equal(big.NewInt(0), store.GetSumRep(post1))
	assert.Equal(big.NewInt(0), store.GetUserDonatedOn(user1, post1))
	assert.Equal(big.NewInt(0), store.GetUserLastDonation(user1, post1))
	assert.Equal(big.NewInt(0), store.GetUserLastReport(user1, post1))
	assert.Empty(store.GetRoundResult(0))
	assert.Empty(store.GetRoundResult(1))
	assert.Equal(big.NewInt(0), store.GetRoundSumDp(0))
	assert.Equal(big.NewInt(0), store.GetRoundSumDp(1))
	assert.Empty(store.GetRoundTopNPosts(0))
	assert.Empty(store.GetRoundStartAt(0))
	assert.Equal(big.NewInt(0), store.GetRoundPostSumStake(0, post1))
	assert.Equal(big.NewInt(0), store.GetRoundNumKeysSold(0, post1))
	assert.Equal(big.NewInt(0), store.GetRoundNumKeysHas(0, user1, post1))
}

func TestStartNewRound(t *testing.T) {
	assert := assert.New(t)
	store := newReputationStoreOnMock()

	assert.Equal(RoundId(1), store.GetCurrentRound())
	store.StartNewRound(222)
	assert.Equal(RoundId(2), store.GetCurrentRound())
	assert.Equal(int64(222), store.GetRoundStartAt(2))
	assert.Empty(store.GetRoundTopNPosts(2))
	assert.Empty(store.GetRoundResult(2))
	assert.Equal(big.NewInt(0), store.GetRoundSumDp(1))
	assert.Equal(big.NewInt(0), store.GetRoundSumDp(2))
}

func TestTopN(t *testing.T) {
	assert := assert.New(t)
	post1 := "bla"
	post2 := "zzz"
	store := newReputationStoreOnMock()

	// test sorting
	store.StartNewRound(222)
	store.SetRoundPostSumDp(2, post1, big.NewInt(100))
	assert.Equal([]PostDpPair{{post1, big.NewInt(100)}}, store.GetRoundTopNPosts(2))
	store.SetRoundPostSumDp(2, post2, big.NewInt(300))
	assert.Equal([]PostDpPair{{post2, big.NewInt(300)}, {post1, big.NewInt(100)}}, store.GetRoundTopNPosts(2))
	store.SetRoundPostSumDp(2, post1, big.NewInt(1000))
	assert.Equal([]PostDpPair{{post1, big.NewInt(1000)}, {post2, big.NewInt(300)}}, store.GetRoundTopNPosts(2))

	for i := 1; i <= DefaultBestContentIndexN; i++ {
		store.SetRoundPostSumDp(2, "p"+string(i), big.NewInt(int64(i)))
	}
	for i := DefaultBestContentIndexN; i >= 0; i-- {
		store.SetRoundPostSumDp(2, "pp"+string(i), big.NewInt(int64(i)))
	}
	topN := store.GetRoundTopNPosts(2)

	// at most N.
	assert.Equal(DefaultBestContentIndexN, len(topN))
	// decreasing order
	for i, v := range topN {
		if i > 0 {
			assert.Truef(v.SumDp.Cmp(topN[i-1].SumDp) <= 0, "%+v, %+v", v, topN[i-1])
		}
	}
}

// This is intended to be a large test to cover possible wrong prefix write bugs.
// For example, if the prefix is wrongly used.
func TestStoreGetSet(t *testing.T) {
	assert := assert.New(t)
	user1 := "test"
	user2 := "test2"
	post1 := "post1"
	post2 := "post2"
	store := newReputationStoreOnMock()

	store.SetCustomerScore(user1, big.NewInt(101))
	store.SetCustomerScore(user2, big.NewInt(102))
	defer func() { assert.Equal(store.GetCustomerScore(user1), big.NewInt(101)) }()
	defer func() { assert.Equal(store.GetCustomerScore(user2), big.NewInt(102)) }()

	store.SetFreeScore(user1, big.NewInt(111))
	store.SetFreeScore(user2, big.NewInt(222))
	defer func() { assert.Equal(store.GetFreeScore(user1), big.NewInt(111)) }()
	defer func() { assert.Equal(store.GetFreeScore(user2), big.NewInt(222)) }()

	store.SetUserLastSettled(user1, 888)
	store.SetUserLastSettled(user2, 999)
	defer func() { assert.Equal(store.GetUserLastSettled(user1), int64(888)) }()
	defer func() { assert.Equal(store.GetUserLastSettled(user2), int64(999)) }()

	store.SetUserLastDonationRound(user1, 201)
	store.SetUserLastDonationRound(user2, 202)
	defer func() { assert.Equal(store.GetUserLastDonationRound(user1), int64(201)) }()
	defer func() { assert.Equal(store.GetUserLastDonationRound(user2), int64(202)) }()

	store.SetSumRep(post1, big.NewInt(1000))
	store.SetSumRep(post2, big.NewInt(2000))
	defer func() { assert.Equal(store.GetSumRep(post1), big.NewInt(1000)) }()
	defer func() { assert.Equal(store.GetSumRep(post2), big.NewInt(2000)) }()

	store.SetUserDonatedOn(user1, post1, big.NewInt(77))
	store.SetUserDonatedOn(user1, post2, big.NewInt(88))
	store.SetUserDonatedOn(user2, post1, big.NewInt(99))
	defer func() { assert.Equal(big.NewInt(77), store.GetUserDonatedOn(user1, post1)) }()
	defer func() { assert.Equal(big.NewInt(88), store.GetUserDonatedOn(user1, post2)) }()
	defer func() { assert.Equal(big.NewInt(99), store.GetUserDonatedOn(user2, post1)) }()

	store.SetUserLastDonation(user1, post1, big.NewInt(222))
	store.SetUserLastDonation(user1, post2, big.NewInt(233))
	store.SetUserLastDonation(user2, post1, big.NewInt(244))
	defer func() { assert.Equal(big.NewInt(222), store.GetUserLastDonation(user1, post1)) }()
	defer func() { assert.Equal(big.NewInt(233), store.GetUserLastDonation(user1, post2)) }()
	defer func() { assert.Equal(big.NewInt(244), store.GetUserLastDonation(user2, post1)) }()

	store.SetUserLastReport(user1, post1, big.NewInt(1122))
	store.SetUserLastReport(user1, post2, big.NewInt(1133))
	store.SetUserLastReport(user2, post1, big.NewInt(1144))
	defer func() { assert.Equal(big.NewInt(1122), store.GetUserLastReport(user1, post1)) }()
	defer func() { assert.Equal(big.NewInt(1133), store.GetUserLastReport(user1, post2)) }()
	defer func() { assert.Equal(big.NewInt(1144), store.GetUserLastReport(user2, post1)) }()

	store.SetRoundSumDp(2, big.NewInt(1022))
	store.SetRoundSumDp(2, big.NewInt(1033))
	store.SetRoundSumDp(3, big.NewInt(1044))
	defer func() { assert.Equal(big.NewInt(1033), store.GetRoundSumDp(2)) }()
	defer func() { assert.Equal(big.NewInt(1044), store.GetRoundSumDp(3)) }()

	store.SetRoundPostSumDp(2, post1, big.NewInt(100))
	store.SetRoundPostSumDp(2, post1, big.NewInt(199))
	store.SetRoundPostSumDp(2, post2, big.NewInt(200))
	defer func() { assert.Equal(big.NewInt(199), store.GetRoundPostSumDp(2, post1)) }()
	defer func() { assert.Equal(big.NewInt(200), store.GetRoundPostSumDp(2, post2)) }()

	store.SetRoundPostSumStake(2, post1, big.NewInt(700))
	store.SetRoundPostSumStake(2, post1, big.NewInt(701))
	store.SetRoundPostSumStake(2, post2, big.NewInt(722))
	defer func() { assert.Equal(big.NewInt(701), store.GetRoundPostSumStake(2, post1)) }()
	defer func() { assert.Equal(big.NewInt(722), store.GetRoundPostSumStake(2, post2)) }()

	store.SetRoundNumKeysHas(2, user1, post1, big.NewInt(77))
	store.SetRoundNumKeysHas(2, user1, post1, big.NewInt(233))
	store.SetRoundNumKeysHas(2, user1, post2, big.NewInt(444))
	store.SetRoundNumKeysHas(2, user2, post1, big.NewInt(666))
	defer func() { assert.Equal(big.NewInt(233), store.GetRoundNumKeysHas(2, user1, post1)) }()
	defer func() { assert.Equal(big.NewInt(444), store.GetRoundNumKeysHas(2, user1, post2)) }()
	defer func() { assert.Equal(big.NewInt(666), store.GetRoundNumKeysHas(2, user2, post1)) }()

	store.SetRoundNumKeysSold(2, post1, big.NewInt(1))
	store.SetRoundNumKeysSold(2, post1, big.NewInt(1000))
	store.SetRoundNumKeysSold(2, post2, big.NewInt(2000))
	defer func() { assert.Equal(big.NewInt(1000), store.GetRoundNumKeysSold(2, post1)) }()
	defer func() { assert.Equal(big.NewInt(2000), store.GetRoundNumKeysSold(2, post2)) }()
}
