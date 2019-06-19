package repv2

import (
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
)

// This package is not thread-safe.

type Reputation interface {
	// import reputation score from V1, foreach user @p u, this should only
	// be called once(later calls will be ignored).
	MigrateFromV1(u Uid, rep Rep)

	// requires importing reputation score from V1 if
	// no donation has been recorded for this user in V2.
	RequireMigrate(u Uid) bool

	// Donate to post @p p wit @p s coins.
	// Note that if migrate is required, it must be done before donate.
	DonateAt(u Uid, p Pid, s LinoCoin) IF

	// user's freescore += @p r, NOTE: unit is COIN.
	IncFreeScore(u Uid, r Rep)

	// needs to be called every endblocker.
	Update(t Time)

	// current reputation of the user.
	GetReputation(u Uid) Rep

	// Round 0 is an invalidated round
	// Round 1 is a short round that will last for only one block, because round-1's
	// start time is set to 0.
	GetCurrentRound() (RoundId, Time) // current round and its start time.

	// ExportImporter
	ExportToFile(file string)
	ImportFromFile(file string)
}

type ReputationImpl struct {
	store                ReputationStore
	BestN                int
	UserMaxN             int
	RoundDurationSeconds int64
	SampleWindowSize     int64
	DecayFactor          int64
}

var _ Reputation = ReputationImpl{}

func NewReputation(s ReputationStore, BestN int, UserMaxN int, RoundDurationSeconds, SampleWindowSize, DecayFactor int64) Reputation {
	return ReputationImpl{
		store:                s,
		BestN:                BestN,
		UserMaxN:             UserMaxN,
		RoundDurationSeconds: RoundDurationSeconds,
		SampleWindowSize:     SampleWindowSize,
		DecayFactor:          DecayFactor,
	}
}

// ExportToFile - implementing ExporteImporter
func (rep ReputationImpl) ExportToFile(file string) {
	// before calling store's export, update reputation.
	rep.store.IterateUsers(func(u Uid) bool {
		rep.GetReputation(u)
		return false
	})
	rst := rep.store.Export()
	f, err := os.Create(file)
	if err != nil {
		panic("failed to create account")
	}
	defer f.Close()
	jsonbytes, err := cdc.MarshalJSON(rst)
	if err != nil {
		panic("failed to marshal json for " + file + " due to " + err.Error())
	}
	f.Write(jsonbytes)
	f.Sync()
}

// ImportFromFile - implementing ExporteImporter
func (rep ReputationImpl) ImportFromFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic("failed to open " + err.Error())
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic("failed to readall: " + err.Error())
	}
	dt := &UserReputationTable{}
	err = cdc.UnmarshalJSON(bytes, dt)
	if err != nil {
		panic("failed to unmarshal: " + err.Error())
	}
	rep.store.Import(dt)
}

// Import reputation score from V1
// The previous reputation score becomes consumption in V2.
// User's last donation round will be reset to current round to
// indicate the this user has migrated.
func (rep ReputationImpl) MigrateFromV1(u Uid, prevRep Rep) {
	if !rep.RequireMigrate(u) {
		return
	}
	user := rep.store.GetUserMeta(u)
	defer func() {
		rep.store.SetUserMeta(u, user)
	}()
	user.LastDonationRound = rep.store.GetCurrentRound()
	user.Consumption = prevRep
	user.Consumption = bigIntMax(user.Consumption, big.NewInt(0))
	user.Reputation = rep.computeReputation(user.Consumption, user.Hold)
}

// any donation or import will make user's lastDonationRound
// to be updated to newest.
func (rep ReputationImpl) RequireMigrate(u Uid) bool {
	last := rep.store.GetUserMeta(u).LastDonationRound
	if last > 0 {
		return false
	}
	return true
}

// internal struct as a summary of a user's consumption in a round.
type consumptionInfo struct {
	seed    LinoCoin
	other   LinoCoin
	seedIF  IF
	otherIF IF
}

// extract user' consumptionInfo from userMeta. The @p seedSet is the result
// of the user's the last donation, and also the round when donations in user.Unsettled
// happened.
func (rep ReputationImpl) extractConsumptionInfo(user *userMeta, seedSet map[Pid]bool) consumptionInfo {
	seed := big.NewInt(0)
	other := big.NewInt(0)
	seedIF := big.NewInt(0)
	otherIF := big.NewInt(0)
	for _, pd := range user.Unsettled {
		amount := pd.Amount
		pid := pd.Pid
		impact := pd.Impact
		if _, isTop := seedSet[pid]; isTop {
			seed.Add(seed, amount)
			seedIF.Add(seedIF, impact)
		} else {
			other.Add(other, amount)
			otherIF.Add(otherIF, impact)
		}
	}
	return consumptionInfo{
		seed:    seed,
		other:   other,
		seedIF:  seedIF,
		otherIF: otherIF,
	}
}

// return the seed set of the @p id round.
func (rep ReputationImpl) getSeedSet(id RoundId) map[Pid]bool {
	result := rep.store.GetRoundMeta(id)
	tops := make(map[Pid]bool)
	for _, p := range result.Result {
		tops[p] = true
	}
	return tops
}

// compute reputation from @p consumption and @p hold.
func (rep ReputationImpl) computeReputation(consumption LinoCoin, hold Rep) Rep {
	return bigIntMax(
		bigIntSub(consumption, bigIntMul(hold, big.NewInt(rep.SampleWindowSize))),
		big.NewInt(0),
	)
}

// internal struct for reputation data of a user.
type reputationData struct {
	consumption LinoCoin
	hold        Rep
	reputation  Rep
}

// will mutate @p user to new reputation score based on @p consumptions
// Return a new reputationData that is computed from (a user's) @p repData with
// new @p consumptions in a round.
func (rep ReputationImpl) computeNewRepData(repData reputationData, consumptions consumptionInfo) reputationData {
	seed := consumptions.seed
	other := consumptions.other
	seedIF := consumptions.seedIF
	otherIF := consumptions.otherIF

	adjustedConsumption := bigIntMin(
		bigIntMax(bigIntDivFrac(seedIF, 8, 10), seed),
		bigIntAdd(seed, other),
	)
	newConsumption := repData.consumption
	if bigIntGreater(adjustedConsumption, repData.consumption) {
		newConsumption = bigIntEMA(repData.consumption, adjustedConsumption, rep.SampleWindowSize)
	}
	otherLimit := bigIntDiv(bigIntAdd(seedIF, otherIF), big.NewInt(5)) // * 20%
	if bigIntGreater(otherIF, otherLimit) {
		newConsumption = bigIntSub(newConsumption,
			bigIntMax(big.NewInt(1),
				bigIntMulFrac(bigIntSub(otherIF, otherLimit), rep.DecayFactor, 100)))
		newConsumption = bigIntMax(big.NewInt(0), newConsumption)
	}

	if bigIntGTE(newConsumption, repData.consumption) {
		delta := bigIntSub(newConsumption, repData.consumption)
		repData.hold = bigIntMax(big.NewInt(1), bigIntEMA(repData.hold, delta, rep.SampleWindowSize))
	}

	repData.consumption = newConsumption
	repData.reputation = rep.computeReputation(repData.consumption, repData.hold)
	return repData
}

// update @p user with information of @p current round.
func (rep ReputationImpl) updateReputation(user *userMeta, current RoundId) {
	// needs to update user's reputation only when the last settled
	// round is less than the last donation round, *and* current round
	// is newer than last donation round(i.e. the last donation round has ended).
	// if above conditions do not hold, skip update.
	if !(user.LastSettledRound < user.LastDonationRound && user.LastDonationRound < current) {
		return
	}

	seedset := rep.getSeedSet(user.LastDonationRound)
	consumptions := rep.extractConsumptionInfo(user, seedset)
	newrep := rep.computeNewRepData(reputationData{
		consumption: user.Consumption,
		hold:        user.Hold,
		reputation:  user.Reputation,
	}, consumptions)

	user.Consumption = newrep.consumption
	user.Hold = newrep.hold
	user.Reputation = newrep.reputation
	user.LastSettledRound = user.LastDonationRound
	user.Unsettled = nil
}

// return the reputation of @p u.
func (rep ReputationImpl) GetReputation(u Uid) Rep {
	user := rep.store.GetUserMeta(u)
	defer func() {
		rep.store.SetUserMeta(u, user)
	}()
	current := rep.store.GetCurrentRound()
	rep.updateReputation(user, current)
	return user.Reputation
}

// Record @p u has donated to @p p with @p amount LinoCoin.
// Only the first UserMaxN posts will be counted and have impact.
// The invarience is that user will have only *one* round of donations unsettled,
// either because that round is current round(not ended yet),
// or the user has never donated after that round.
// So when a user donates, we first update user's reputation, then
// we add this donation to records.
func (rep ReputationImpl) DonateAt(u Uid, p Pid, amount LinoCoin) IF {
	if len(u) == 0 {
		panic("Length of Uid must be longer than 0")
	}
	if len(p) == 0 {
		panic("Length of Pid must be longer than 0")
	}
	var current RoundId = rep.store.GetCurrentRound()
	user := rep.store.GetUserMeta(u)
	defer func() {
		rep.store.SetUserMeta(u, user)
	}()
	rep.updateReputation(user, current)
	user.LastDonationRound = current
	impact := rep.appendDonation(user, p, amount)
	rep.incRoundPostSumImpact(current, p, impact)
	return impact
}

// appendDonation: append a new donation to user's unsettled list, return the impact
// factor of this donation.
// contract: before calling this, user's reputation needs to
//           be updated by calling updateReputation.
func (rep ReputationImpl) appendDonation(user *userMeta, post Pid, amount LinoCoin) IF {
	reputation := user.Reputation
	pos := -1
	used := big.NewInt(0)
	for i, v := range user.Unsettled {
		used.Add(used, v.Impact)
		if v.Pid == post {
			pos = i
		}
	}
	if pos == -1 && len(user.Unsettled) >= rep.UserMaxN {
		return big.NewInt(0)
	}
	var available IF = bigIntMax(bigIntSub(reputation, used), big.NewInt(0))
	var impact IF = bigIntMin(available, amount)
	if pos != -1 {
		user.Unsettled[pos].Amount.Add(user.Unsettled[pos].Amount, amount)
		user.Unsettled[pos].Impact.Add(user.Unsettled[pos].Impact, impact)
	} else {
		user.Unsettled = append(user.Unsettled, Donation{
			Pid:    post,
			Amount: amount,
			Impact: impact,
		})
	}
	return impact
}

// increase the sum of impact factors of @p post by @p dp, in @p round
// It also maintains an order of posts of the round by bubbling up the rank
// of post on impact factor increasing.
func (rep ReputationImpl) incRoundPostSumImpact(round RoundId, p Pid, dp IF) {
	roundPost := rep.store.GetRoundPostMeta(round, p)
	defer func() {
		rep.store.SetRoundPostMeta(round, p, roundPost)
	}()
	roundMeta := rep.store.GetRoundMeta(round)
	defer func() {
		rep.store.SetRoundMeta(round, roundMeta)
	}()

	roundMeta.SumIF.Add(roundMeta.SumIF, dp)
	newSumIF := bigIntAdd(roundPost.SumIF, dp)
	roundPost.SumIF = newSumIF

	pos := -1
	for i, v := range roundMeta.TopN {
		if v.Pid == p {
			roundMeta.TopN[i].SumIF = newSumIF
			pos = i
			break
		}
	}
	if pos == -1 {
		// first post.
		if len(roundMeta.TopN) < rep.BestN {
			roundMeta.TopN = append(roundMeta.TopN, PostIFPair{Pid: p, SumIF: newSumIF})
		} else {
			lastSumIF := roundMeta.TopN[len(roundMeta.TopN)-1].SumIF
			// lastSumIF < newSumIF
			if bigIntLess(lastSumIF, newSumIF) {
				roundMeta.TopN = append(roundMeta.TopN, PostIFPair{Pid: p, SumIF: newSumIF})
			} else {
				// do not need to do anything, as this post's sumDP
				// is less or equal to the last one.
				return
			}
		}
		pos = len(roundMeta.TopN) - 1
	}
	bubbleUp(roundMeta.TopN, pos)
	// keeping bestN only.
	if len(roundMeta.TopN) > rep.BestN {
		roundMeta.TopN = roundMeta.TopN[:rep.BestN]
	}
}

// return the current round id the the start time of the round.
func (rep ReputationImpl) GetCurrentRound() (RoundId, Time) {
	rid := rep.store.GetCurrentRound()
	return rid, rep.store.GetRoundMeta(rid).StartAt
}

// increase @p u user's reputation by @p score.
// To make added score permanent, add it on consumption, as reputation is
// only a temporory result, same in reputation migration.
func (rep ReputationImpl) IncFreeScore(u Uid, score Rep) {
	user := rep.store.GetUserMeta(u)
	defer func() {
		rep.store.SetUserMeta(u, user)
	}()
	user.Consumption.Add(user.Consumption, score)
	user.Consumption = bigIntMax(user.Consumption, big.NewInt(0))
	user.Reputation = rep.computeReputation(user.Consumption, user.Hold)
}

// On BlockEnd(@p t), select out the seed set of the current round and start
// a new round.
func (rep ReputationImpl) Update(t Time) {
	current := rep.store.GetCurrentRound()
	roundInfo := rep.store.GetRoundMeta(current)
	if t-roundInfo.StartAt >= rep.RoundDurationSeconds {
		// need to update only when it is updated.
		defer func() {
			rep.store.SetRoundMeta(current, roundInfo)
			// start a new round
			rep.StartNewRound(t)
		}()
		// process all information of this round
		// Find out top N.
		topN := roundInfo.TopN
		sumDpInRound := roundInfo.SumIF
		// XXX(yumin): taking a floor of 80%.
		dpBound := bigIntMulFrac(sumDpInRound, 8, 10)
		dpCovered := big.NewInt(0)
		var rst []Pid
		for _, pidDp := range topN {
			pid := pidDp.Pid
			postIF := pidDp.SumIF
			dpCovered.Add(dpCovered, postIF)
			rst = append(rst, pid)
			if !bigIntLess(dpCovered, dpBound) {
				break
			}
		}
		roundInfo.Result = rst
	}
}

// write (RoundId + 1, t) into db and update current round
func (rep ReputationImpl) StartNewRound(t Time) {
	gameMeta := rep.store.GetGameMeta()
	defer func() {
		rep.store.SetGameMeta(gameMeta)
	}()
	newRoundId := gameMeta.CurrentRound + 1
	gameMeta.CurrentRound = newRoundId

	newRoundMeta := &roundMeta{
		Result:  nil,
		SumIF:   big.NewInt(0),
		StartAt: t,
		TopN:    nil,
	}
	rep.store.SetRoundMeta(newRoundId, newRoundMeta)
}

// a bunch of helper functions that takes two bitInt and returns
// a newly allocated bigInt.
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

func bigIntMin(a, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return big.NewInt(0).Set(a)
	} else {
		return big.NewInt(0).Set(b)
	}
}

func bigIntMax(a, b *big.Int) *big.Int {
	if a.Cmp(b) > 0 {
		return big.NewInt(0).Set(a)
	} else {
		return big.NewInt(0).Set(b)
	}
}

func bigIntGreater(a, b *big.Int) bool {
	return a.Cmp(b) > 0
}

func bigIntGTE(a, b *big.Int) bool {
	return a.Cmp(b) >= 0
}

func bigIntLess(a, b *big.Int) bool {
	return a.Cmp(b) < 0
}

func bigIntLTE(a, b *big.Int) bool {
	return a.Cmp(b) <= 0
}

// return the exponential moving average of @p prev on having a new sample @p new
// with sample size of @p windowSize.
func bigIntEMA(prev, new *big.Int, windowSize int64) *big.Int {
	if windowSize <= 0 {
		panic("bigIntEMA illegal windowSize: " + strconv.FormatInt(windowSize, 10))
	}
	return bigIntDiv(
		bigIntAdd(new, bigIntMul(prev, big.NewInt(windowSize-1))),
		big.NewInt(windowSize))
}

// return v / (num / denom)
func bigIntDivFrac(v *big.Int, num, denom int64) *big.Int {
	if num == 0 || denom == 0 {
		panic("bigIntDivFrac zero num or denom")
	}
	return bigIntMulFrac(v, denom, num)
}

// return v * (num / denom)
func bigIntMulFrac(v *big.Int, num, denom int64) *big.Int {
	if denom == 0 {
		panic("bigIntMulFrac zero denom")
	}
	return bigIntDiv(bigIntMul(v, big.NewInt(num)), big.NewInt(denom))
}

// contract:
//     before: all inversions are related to posts[pos].
//     after:  posts are sorted by SumIF, decreasingly.
func bubbleUp(posts []PostIFPair, pos int) {
	for i := pos; i > 0; i-- {
		if bigIntLess(posts[i-1].SumIF, posts[i].SumIF) {
			posts[i], posts[i-1] = posts[i-1], posts[i]
		} else {
			break
		}
	}
}
