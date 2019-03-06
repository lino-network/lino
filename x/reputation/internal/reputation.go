package internal

import (
	"math/big"
)

// This package is not thread-safe.

type Reputation interface {
	DonateAt(u Uid, p Pid, s Stake) Dp
	ReportAt(u Uid, p Pid) Rep
	// user's freescore += @p r, NOTE: unit is COIN.
	IncFreeScore(u Uid, r Rep)
	Update(t Time) // called every endblocker.
	// current reputation of the user.
	GetReputation(u Uid) Rep
	GetSumRep(p Pid) Rep
	GetCurrentRound() (RoundId, Time) // current round and its start time.

	// ExportImporter
	ExportToFile(file string)
	// Import(dt []byte) error
}

type ReputationImpl struct {
	store ReputationStore
}

func NewReputation(s ReputationStore) Reputation {
	return &ReputationImpl{store: s}
}

// ExportToFile - implementing ExporteImporter
func (rep ReputationImpl) ExportToFile(f string) {
	rep.store.ExportToFile(f)
}

// Import - implementing ExporteImporter
// func (rep ReputationImpl) Import(dt []byte) error {
// 	rep.store.Import(dt)
// }

func (rep ReputationImpl) GetReputation(u Uid) Rep {
	customerScore := rep.GetSettledCustomerScore(u)
	freeScore := rep.store.GetFreeScore(u)
	return bigIntAdd(customerScore, freeScore)
}

func (rep ReputationImpl) GetSumRep(p Pid) Rep {
	return rep.store.GetSumRep(p)
}

func (rep ReputationImpl) GetSettledCustomerScore(u Uid) Rep {
	current := rep.store.GetCurrentRound()
	lastSettled := rep.store.GetUserLastSettled(u)
	lastDonated := rep.store.GetUserLastDonationRound(u)
	customerScore := rep.store.GetCustomerScore(u)

	// if last donated round is not settled, that round has ended.
	if lastSettled < lastDonated && lastDonated < current {
		unsettledScore := big.NewInt(0)
		bestContents := rep.store.GetRoundResult(lastDonated)
		for _, pid := range bestContents {
			numKeysHas := rep.store.GetRoundNumKeysHas(lastDonated, u, pid)
			if numKeysHas != nil {
				numKeysSold := rep.store.GetRoundNumKeysSold(lastDonated, pid)
				totalStake := rep.store.GetRoundPostSumStake(lastDonated, pid)
				if numKeysSold.Cmp(BigIntZero) > 0 {
					unsettledScore.Add(unsettledScore,
						bigIntDiv(bigIntMul(numKeysHas, totalStake), numKeysSold))
				}
			}
		}
		newScore :=
			bigIntDiv(
				bigIntAdd(
					bigIntMul(
						customerScore,
						big.NewInt(SampleWindowSize-1)),
					unsettledScore),
				big.NewInt(SampleWindowSize))

		customerScore = bigIntMax(newScore,
			bigIntDiv(bigIntMul(customerScore, big.NewInt(DecayFactor)), big.NewInt(100)))
		customerScore = bigIntMax(customerScore, big.NewInt(InitialCustomerScore))
		rep.store.SetUserLastSettled(u, lastDonated) // last donated round is settled.
		rep.store.SetCustomerScore(u, customerScore)
	}
	return customerScore
}

// currently all donations count, though we said that only the first 40 donation counts in the spec.
func (rep ReputationImpl) DonateAt(u Uid, p Pid, s Stake) Dp {
	if len(u) == 0 {
		panic("Uid must be longer than 0")
	}
	if len(p) == 0 {
		panic("Pid must be longer than 0")
	}
	var uRep Rep = rep.GetReputation(u)
	var usedDp Dp = rep.store.GetUserDonatedOn(u, p)
	var availableDp Dp = bigIntMax(bigIntSub(uRep, usedDp), BigIntZero)
	var dp Dp = bigIntMin(availableDp, s)
	// update user's donation to p.
	rep.store.SetUserDonatedOn(u, p, bigIntAdd(dp, usedDp))
	rep.incPostSumRep(u, p, uRep)
	// update this round-related information.
	var current RoundId = rep.store.GetCurrentRound()
	rep.store.SetUserLastDonationRound(u, current)
	rep.buyKey(current, u, p, s)
	rep.incPostSumDp(current, p, dp)
	rep.incRoundSumDp(current, dp)
	rep.incPostSumStake(current, p, s)
	return dp
}

// buy keys of post @p p for user @p u using @p dp.
func (rep ReputationImpl) buyKey(roundId RoundId, u Uid, p Pid, stake Stake) {
	numKeysSold := rep.store.GetRoundNumKeysSold(roundId, p)
	numKeysBuy := rep.numKeysCanBuy(numKeysSold, stake)
	if numKeysBuy.Cmp(BigIntZero) > 0 {
		atHand := rep.store.GetRoundNumKeysHas(roundId, u, p)
		rep.store.SetRoundNumKeysHas(roundId, u, p, atHand.Add(atHand, numKeysBuy))
		rep.store.SetRoundNumKeysSold(roundId, p, numKeysSold.Add(numKeysSold, numKeysBuy))
	}
}

// returns the number of keys that uses @p stake can buy.
// bsearch: ~16000ns/op
// sqrt: ~1400ns/op
/////////////////////////////////////////////////////////////////////////
// sqrt version is much faster, though it is not as accurate as bsearch.
// sqrt version
// func (rep ReputationImpl) numKeysCanBuy(numKeysSold *big.Int, stake Stake) (numKeysCanBuy *big.Int) {
// 	paraC := big.NewInt(KeyPriceC)
// 	// current price = C + n * K, when K == 1, it becomes C + n.
// 	currentPrice := bigIntAdd(paraC, numKeysSold)

// 	// number of keys = sqrt(c * (c - 2) + 2*s) - c
// 	temp1 := bigIntMul(currentPrice, bigIntSub(currentPrice, big.NewInt(2))) // c * (c - 2)
// 	temp2 := bigIntAdd(temp1, bigIntLsh(stake, 1)) // c * (c - 2) + 2*s

// 	if temp2.Cmp(BigIntZero) >= 0 {
// 		temp3 := temp2.Sqrt(temp2) // sqrt(c * (c - 2) + 2*s)
// 		rst := bigIntSub(temp3, currentPrice) // sqrt(c * (c - 2) + 2*s) - c
//      rst.Add(rst, big.NewInt(1))
// 		numKeysCanBuy = bigIntMax(rst, big.NewInt(0))
// 	} else {
// 		numKeysCanBuy = big.NewInt(0)
// 	}
// 	return
// }
func (rep ReputationImpl) numKeysCanBuy(numKeysSold *big.Int, stake Stake) (numKeysCanBuy *big.Int) {
	paraC := big.NewInt(KeyPriceC)
	// current price = C + n * K, when K == 1, it becomes C + n.
	currentPrice := bigIntAdd(paraC, numKeysSold)
	// binary search on the largest n that
	// (n+1)*c + (k + k*n)*n/2 <= stake, which is the sum of price of [0..N] keys
	// since we assume that k is always 1, it becomes:
	// (n+1)*c + (1 + 1*n)*n/2 <= stake
	// use alternative form to speed it up
	// (n + 1) * (2*c + n) <= 2*stake

	// Edge case.
	if stake.Cmp(currentPrice) < 0 {
		return big.NewInt(0)
	}

	twoC := bigIntLsh(currentPrice, 1)
	twoS := bigIntLsh(stake, 1)
	// (n + 1) * (2*c + n)
	eval := func(n int64) *big.Int {
		return bigIntMul(bigIntAdd(big.NewInt(1), big.NewInt(n)),
			bigIntAdd(twoC, big.NewInt(n)))
	}

	beg, end := int64(0), int64(MaxNumKeysEachTime+1)
	// optimization on common cases, ~8500 ns/op
	commonValues := []int64{OneLinoCoin, 100 * OneLinoCoin, 10000 * OneLinoCoin}
	for _, v := range commonValues {
		if eval(v).Cmp(twoS) > 0 {
			end = v
			break
		} else {
			beg = v
		}
	}

	for beg+1 < end {
		mid := (beg + end) >> 1
		val := eval(mid)
		if val.Cmp(twoS) <= 0 {
			beg = mid
		} else {
			end = mid
		}
	}

	// +1 because in [0..N], there are n+1 keys.
	return big.NewInt(beg + 1)
}

// increase the post sum dp.
func (rep ReputationImpl) incPostSumDp(roundId RoundId, p Pid, dp Dp) {
	sumDp := rep.store.GetRoundPostSumDp(roundId, p)
	sumDp.Add(sumDp, dp)
	rep.store.SetRoundPostSumDp(roundId, p, sumDp)
}

// increase the post sum dp.
func (rep ReputationImpl) incPostSumStake(roundId RoundId, p Pid, s Stake) {
	sumStake := rep.store.GetRoundPostSumStake(roundId, p)
	sumStake.Add(sumStake, s)
	rep.store.SetRoundPostSumStake(roundId, p, sumStake)
}

// increase the total dp value of this post. Note that, if an user has donated before,
// and if his reputation has changed since then, the delta will be counted.
func (rep ReputationImpl) incPostSumRep(u Uid, p Pid, newRep Rep) {
	sumRep := rep.store.GetSumRep(p)
	oldRep := rep.store.GetUserLastDonation(u, p)
	var delta Rep
	delta = bigIntSub(newRep, oldRep)
	rep.store.SetSumRep(p, sumRep.Add(sumRep, delta))
	rep.store.SetUserLastDonation(u, p, newRep)
}

// increasing the total dp value of this round.
func (rep ReputationImpl) incRoundSumDp(roundId RoundId, dp Dp) {
	sumDp := rep.store.GetRoundSumDp(roundId)
	rep.store.SetRoundSumDp(roundId, sumDp.Add(sumDp, dp))
}

// similar to incPostSumRep, delta is included as well.
func (rep ReputationImpl) ReportAt(u Uid, p Pid) Rep {
	sumRep := rep.store.GetSumRep(p)
	newRep := rep.GetReputation(u) // new user rep
	var delta Rep
	oldRep := rep.store.GetUserLastReport(u, p)
	delta = bigIntSub(newRep, oldRep)
	sumRep.Sub(sumRep, delta)
	rep.store.SetSumRep(p, sumRep)
	rep.store.SetUserLastReport(u, p, newRep)
	return sumRep
}

func (rep ReputationImpl) GetCurrentRound() (RoundId, Time) {
	rid := rep.store.GetCurrentRound()
	startAt := rep.store.GetRoundStartAt(rid)
	return rid, startAt
}

func (rep ReputationImpl) IncFreeScore(u Uid, score Rep) {
	freescore := rep.store.GetFreeScore(u)
	freescore.Add(freescore, score)
	rep.store.SetFreeScore(u, freescore)
}

func (rep ReputationImpl) Update(t Time) {
	round := rep.store.GetCurrentRound()
	startAt := rep.store.GetRoundStartAt(round)
	if rep.moreThan(t, startAt, RoundDuration) {
		// process all information of this round
		// Find out top N.
		topN := rep.store.GetRoundTopNPosts(round)
		sumDpInRound := rep.store.GetRoundSumDp(round)
		// XXX(yumin): taking a floor of 80%.
		dpBound := bigIntDiv(bigIntMul(sumDpInRound, big.NewInt(8)), big.NewInt(10))
		dpCovered := big.NewInt(0)
		var rst []Pid
		for _, pidDp := range topN {
			pid := pidDp.Pid
			postDp := rep.store.GetRoundPostSumDp(round, pid)
			dpCovered.Add(dpCovered, postDp)
			// add to best content index only if it is not reported.
			sumRep := rep.store.GetSumRep(pid)
			if postDp.Cmp(BigIntZero) <= 0 {
				continue
			}
			if sumRep.Cmp(BigIntZero) > 0 {
				rst = append(rst, pid)
			}

			if dpCovered.Cmp(dpBound) >= 0 {
				break
			}
		}
		rep.store.SetRoundResult(round, rst)
		// start a new round
		rep.store.StartNewRound(t)
	}
}

func (rep ReputationImpl) moreThan(cur, startAt Time, hours int64) bool {
	return cur-startAt >= hours*3600
}

// a bunch of helper functions that takes two bitInt and returns
// a newly allocated int
func bigIntAdd(a, b *big.Int) *big.Int {
	rst := big.NewInt(0)
	return rst.Add(a, b)
}

func bigIntSub(a, b *big.Int) *big.Int {
	rst := big.NewInt(0)
	return rst.Sub(a, b)
}

func bigIntMul(a, b *big.Int) *big.Int {
	rst := big.NewInt(0)
	return rst.Mul(a, b)
}

func bigIntDiv(a, b *big.Int) *big.Int {
	rst := big.NewInt(0)
	return rst.Div(a, b)
}

func bigIntLsh(a *big.Int, n uint) *big.Int {
	rst := big.NewInt(0)
	return rst.Lsh(a, n)
}

func bigIntMin(a, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return a
	} else {
		return b
	}
}

func bigIntMax(a, b *big.Int) *big.Int {
	if a.Cmp(b) > 0 {
		return a
	} else {
		return b
	}
}

var _ Reputation = &ReputationImpl{}
