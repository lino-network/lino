package repv2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// This package is not thread-safe.

type Reputation interface {
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
	ExportToFile(file string) error
	ImportFromFile(file string) error

	GetUserMeta(u Uid) string
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

func NewReputation(s ReputationStore, bestN int, userMaxN int, roundDurationSeconds, sampleWindowSize, decayFactor int64) Reputation {
	return ReputationImpl{
		store:                s,
		BestN:                bestN,
		UserMaxN:             userMaxN,
		RoundDurationSeconds: roundDurationSeconds,
		SampleWindowSize:     sampleWindowSize,
		DecayFactor:          decayFactor,
	}
}

// ExportToFile - implementing ExporteImporter
func (rep ReputationImpl) ExportToFile(file string) error {
	// before calling store's export, update reputation.
	rep.store.IterateUsers(func(u Uid) bool {
		rep.GetReputation(u)
		return false
	})
	rst := rep.store.Export()
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	defer f.Close()
	jsonbytes, err := cdc.MarshalJSON(rst)
	if err != nil {
		return fmt.Errorf("failed to marshal json for " + file + " due to " + err.Error())
	}
	_, err = f.Write(jsonbytes)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}

// ImportFromFile - implementing ExporteImporter
func (rep ReputationImpl) ImportFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open " + err.Error())
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to readall: " + err.Error())
	}
	dt := &UserReputationTable{}
	err = cdc.UnmarshalJSON(bytes, dt)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: " + err.Error())
	}
	rep.store.Import(dt)
	return nil
}

// internal struct as a summary of a user's consumption in a round.
type consumptionInfo struct {
	seed    LinoCoin
	other   LinoCoin
	seedIF  IF
	otherIF IF
}


// extract user's consumptionInfo from userMeta. The @p seedSet is the result
// of the user's the last donation, and also the round when donations in user.Unsettled
// happened.
func (rep ReputationImpl) extractConsumptionInfo(user *userMeta, seedSet map[Pid]bool) consumptionInfo {
	seed := NewInt(0)
	other := NewInt(0)
	seedIF := NewInt(0)
	otherIF := NewInt(0)
	for _, pd := range user.Unsettled {
		amount := pd.Amount
		pid := pd.Pid
		impact := pd.Impact
		if _, isTop := seedSet[pid]; isTop {
			seed.Add(amount)
			seedIF.Add(impact)
		} else {
			other.Add(amount)
			otherIF.Add(impact)
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
	return IntMax(
		IntSub(consumption, IntMul(hold, NewInt(rep.SampleWindowSize))),
		NewInt(0),
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

	adjustedConsumption := IntMin(
		IntMax(IntDivFrac(seedIF, 8, 10), seed),
		IntAdd(seed, other),
	)
	newConsumption := repData.consumption
	if IntGreater(adjustedConsumption, repData.consumption) {
		newConsumption = IntEMA(repData.consumption, adjustedConsumption, rep.SampleWindowSize)
	}
	otherLimit := IntDiv(IntAdd(seedIF, otherIF), NewInt(5)) // * 20%
	if IntGreater(otherIF, otherLimit) {
		newConsumption = IntSub(newConsumption,
			IntMax(NewInt(1),
				IntMulFrac(IntSub(otherIF, otherLimit), rep.DecayFactor, 100)))
		newConsumption = IntMax(NewInt(0), newConsumption)
	}

	if IntGTE(newConsumption, repData.consumption) {
		delta := IntSub(newConsumption, repData.consumption)
		repData.hold = IntEMA(repData.hold, delta, rep.SampleWindowSize)
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

	// remove roundPostSumImpacts that are not no longer used.
	// Note: this does not guarantee that all keys will be deleted, but rather
	// set a bound on how much data one user may hold.
	for _, pd := range user.Unsettled {
		rep.deleteRoumdPostSumImpact(user.LastDonationRound, pd.Pid)
	}

	user.Consumption = newrep.consumption
	user.Hold = newrep.hold
	user.Reputation = newrep.reputation
	user.LastSettledRound = user.LastDonationRound
	user.Unsettled = nil
}


func (rep ReputationImpl) GetUserMeta(u Uid) string {
	user := rep.store.GetUserMeta(u)
	rst, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	return string(rst)
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
	used := NewInt(0)
	for i, v := range user.Unsettled {
		used.Add(v.Impact)
		if v.Pid == post {
			pos = i
		}
	}
	if pos == -1 && len(user.Unsettled) >= rep.UserMaxN {
		return NewInt(0)
	}
	var available IF = IntMax(IntSub(reputation, used), NewInt(0))
	var impact IF = IntMin(available, amount)
	if pos != -1 {
		user.Unsettled[pos].Amount.Add(amount)
		user.Unsettled[pos].Impact.Add(impact)
	} else {
		user.Unsettled = append(user.Unsettled, Donation{
			Pid:    post,
			Amount: amount,
			Impact: impact,
		})
	}
	return impact
}

func (rep ReputationImpl) deleteRoumdPostSumImpact(round RoundId, p Pid) {
	rep.store.DelRoundPostMeta(round, p)
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

	roundMeta.SumIF.Add(dp)
	newSumIF := IntAdd(roundPost.SumIF, dp)
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
			if IntLess(lastSumIF, newSumIF) {
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
	user.Consumption.Add(score)
	user.Consumption = IntMax(user.Consumption, NewInt(0))
	user.Reputation = rep.computeReputation(user.Consumption, user.Hold)
}

// On BlockEnd(@p t), select out the seed set of the current round and start
// a new round.
func (rep ReputationImpl) Update(t Time) {
	current := rep.store.GetCurrentRound()
	roundInfo := rep.store.GetRoundMeta(current)
	if int64(t-roundInfo.StartAt) >= rep.RoundDurationSeconds {
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
		dpBound := IntMulFrac(sumDpInRound, 8, 10)
		dpCovered := NewInt(0)
		var rst []Pid
		for _, pidDp := range topN {
			pid := pidDp.Pid
			postIF := pidDp.SumIF
			if !IntGreater(postIF, NewInt(0)) {
				break
			}
			dpCovered.Add(postIF)
			rst = append(rst, pid)
			if !IntLess(dpCovered, dpBound) {
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
		SumIF:   NewInt(0),
		StartAt: t,
		TopN:    nil,
	}
	rep.store.SetRoundMeta(newRoundId, newRoundMeta)
}

// return the exponential moving average of @p prev on having a new sample @p new
// with sample size of @p windowSize.
func IntEMA(prev, new Int, windowSize int64) Int {
	if windowSize <= 0 {
		panic("IntEMA illegal windowSize: " + strconv.FormatInt(windowSize, 10))
	}
	return IntDiv(
		IntAdd(new, IntMul(prev, NewInt(windowSize-1))),
		NewInt(windowSize))
}

// contract:
//     before: all inversions are related to posts[pos].
//     after:  posts are sorted by SumIF, decreasingly.
func bubbleUp(posts []PostIFPair, pos int) {
	for i := pos; i > 0; i-- {
		if IntLess(posts[i-1].SumIF, posts[i].SumIF) {
			posts[i], posts[i-1] = posts[i-1], posts[i]
		} else {
			break
		}
	}
}
