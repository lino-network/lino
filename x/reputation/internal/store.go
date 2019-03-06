package internal

import (
	// XXX(yumin): hmm, seems that golang can not recursively check assignability
	// of interfaces. so we have to import an Iterator type here, because
	// the exact same one we defined is not assignable.
	db "github.com/tendermint/tendermint/libs/db"

	"encoding/binary"
	"math/big"
	"sort"
	"os"
	"strings"
)

// Store - store.
type Store interface {
	Set(key []byte, val []byte)
	Get(key []byte) []byte
	Has(key []byte) bool
	Delete(key []byte)
	// Iterator over a domain of keys in ascending order. End is exclusive.
	// Start must be less than end, or the Iterator is invalid.
	// Iterator must be closed by caller.
	// To iterate over entire domain, use store.Iterator(nil, nil)
	// CONTRACT: No writes may happen within a domain while an iterator exists over it.
	Iterator(start, end []byte) db.Iterator
}

// ReputationStore - xx
// a simple wrapper around kv-store. It does not handle any reputation computation logic,
// except for init value for some field, like customer score.
// This interface should only be used by the reputation system implementation.
// This store needs to be merkelized to support fast rollback.
type ReputationStore interface {
	// TODO(yumin): these two are all in memory, which is extremely bad if state is large.
	// should changed to io.Write/Reader.
	// Export all state to deterministic bytes
	Export() *UserReputationTable
	ExportToFile(file string)

	// Import state from bytes
	Import(tb *UserReputationTable)

	// Note that, this value may not be the exact customer score of the user
	// due to there might be unsettled keys remaining.
	// Also, because that user can get free reputation when they lockdown coins, this value
	// is solely the user's customer score.
	GetCustomerScore(u Uid) Rep
	SetCustomerScore(u Uid, r Rep)

	// user get free score by coin lockdown.
	GetFreeScore(u Uid) Rep
	SetFreeScore(u Uid, r Rep)

	// till which round has user settled, i.e. his customer score is accurate, initialized as 0.
	GetUserLastSettled(u Uid) RoundId
	SetUserLastSettled(u Uid, r RoundId)

	// The last round that user has any donation.
	// contract is that since any donation will force user to settle any previous customer score,
	// so if there would be any unsettled customer score, when lastDonationRound > LastSettled,
	// they are all in the lastDonationRound.
	GetUserLastDonationRound(u Uid) RoundId
	SetUserLastDonationRound(u Uid, r RoundId)

	// total reputation of a post which is the sum of
	// donors' reputation subtracted by reporters' reputation.
	GetSumRep(p Pid) Rep
	SetSumRep(p Pid, rep Rep)

	// @returns: how much @p u has donated to @p p.
	GetUserDonatedOn(u Uid, p Pid) Dp
	SetUserDonatedOn(u Uid, p Pid, dp Dp)

	// if has not donated before, return 0, otherwise reputation at that time.
	GetUserLastDonation(u Uid, p Pid) Rep
	SetUserLastDonation(u Uid, p Pid, rep Rep)

	// if has not reported before, return 0, otherwise reputation at that time.
	GetUserLastReport(u Uid, p Pid) Rep
	SetUserLastReport(u Uid, p Pid, rep Rep)

	// the final result of round, should be set right after round ends.
	GetRoundResult(r RoundId) []Pid
	SetRoundResult(r RoundId, rst []Pid)

	// sum of donation power of all donations happened during the round.
	GetRoundSumDp(r RoundId) Dp
	SetRoundSumDp(r RoundId, dp Dp)

	// top n posts. Ordering is maintained by this store, update upon SetSumDp is called.
	// NOTE: not the final result.
	GetRoundTopNPosts(r RoundId) []PostDpPair

	GetRoundStartAt(round RoundId) Time

	/// -----------  In this round  -------------
	// TODO(yumin): store them together to avoid second read?
	// RoundId is the current round, starts from 1.
	GetCurrentRound() RoundId
	// write (RoundId + 1, t) into db and update current round
	StartNewRound(t Time)

	// total donation power received of a @p post.
	GetRoundPostSumDp(r RoundId, p Pid) Dp
	// update TopN
	SetRoundPostSumDp(r RoundId, p Pid, dp Dp)

	// total stake received of a @p post.
	GetRoundPostSumStake(r RoundId, p Pid) Stake
	SetRoundPostSumStake(r RoundId, p Pid, s Stake)

	// @returns: how many keys has been sold for @p
	GetRoundNumKeysSold(r RoundId, p Pid) bigInt
	SetRoundNumKeysSold(r RoundId, p Pid, numKeys bigInt)

	// @returns: how many keys this user has on this post. return nil if have zero.
	GetRoundNumKeysHas(r RoundId, u Uid, p Pid) bigInt
	SetRoundNumKeysHas(r RoundId, u Uid, p Pid, numKeys bigInt)
}

// This store implementation does not have state. It is just a wrapper of read/write of
// data necessary for reputation system. Also, it takes the problem, that the underlying
// kv store may be an iavl, so the number of keys will have a large impact on performance
// into account, by trying to minimize the number of keys.
// we choosed json as serializer because:
// IT IS deterministic for struct.
// Though gob is good in following aspects, but it's not deterministic for big.Int.
// 1. official package, dependency-free.
// 2. much faster(~4 times), much smaller(..)
// 3. reflection-based, do not need a protocol compiler.
// 4. math/big supported.
// according to benchmarks:
// https://github.com/alecthomas/go_serialization_benchmarks

var (
	KeySeparator               = []byte("/")
	repUserMetaPrefix          = []byte{0x00}
	repPostMetaPrefix          = []byte{0x01}
	repUserPostMetaPrefix      = []byte{0x02}
	repRoundMetaPrefix         = []byte{0x03}
	repRoundPostMetaPrefix     = []byte{0x04}
	repRoundUserPostMetaPrefix = []byte{0x05}
	repGameMetaPrefix          = []byte{0x06}
)

type userMeta struct {
	CustomerScore     Rep
	FreeScore         Rep
	LastSettled       RoundId
	LastDonationRound RoundId
}

type postMeta struct {
	SumRep Rep
}

type userPostMeta struct {
	Donated         Dp
	LastDonationRep Rep
	LastReportRep   Rep
}

type roundMeta struct {
	Result  []Pid
	SumDp   Dp
	StartAt Time
	TopN    []PostDpPair
}

type roundPostMeta struct {
	SumDp       Dp
	SumStake    Stake
	NumKeysSold bigInt
}

type roundUserPostMeta struct {
	NumKeysHas bigInt
}

type gameMeta struct {
	CurrentRound RoundId
}

func getUserMetaKey(u Uid) []byte {
	return append(repUserMetaPrefix, []byte(u)...)
}

func getPostMetaKey(p Pid) []byte {
	return append(repPostMetaPrefix, []byte(p)...)
}

func getRoundMetaKey(r RoundId) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(r))
	return append(repRoundMetaPrefix, buf...)
}

func getUserPostMetaKey(u Uid, p Pid) []byte {
	prefix := append(getUserMetaKey(u), KeySeparator...)
	return append(prefix, []byte(p)...)
}

func getRoundPostMetaKey(r RoundId, p Pid) []byte {
	prefix := append(getRoundMetaKey(r), KeySeparator...)
	return append(prefix, []byte(p)...)
}

func getRoundUserPostMetaKey(r RoundId, u Uid, p Pid) []byte {
	roundPrefix := append(getRoundMetaKey(r), KeySeparator...)
	roundUserPrefix := append(roundPrefix, []byte(u)...)
	roundUserPrefix = append(roundUserPrefix, KeySeparator...)
	return append(roundUserPrefix, []byte(p)...)
}

func getGameKey() []byte {
	return repGameMetaPrefix
}

// The only state is the number of bestContentIndex
type reputationStoreImpl struct {
	store             Store
	BestContentIndexN int
}

func NewReputationStoreDefaultN(s Store) ReputationStore {
	return &reputationStoreImpl{store: s, BestContentIndexN: DefaultBestContentIndexN}
}

func NewReputationStore(s Store, n int) ReputationStore {
	return &reputationStoreImpl{store: s, BestContentIndexN: n}
}

func (impl reputationStoreImpl) Export() *UserReputationTable {
	rst := &UserReputationTable{}
	itr := impl.store.Iterator(repUserMetaPrefix, PrefixEndBytes(repUserMetaPrefix))
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		uid := Uid(itr.Key()[1:])
		if strings.Contains(uid, string(KeySeparator)) {
			continue
		}
		v := impl.getUserMeta(uid)
		rst.reputations = append(rst.reputations, UserReputation{
			Username: uid,
			CustomerScore: v.CustomerScore,
			FreeScore: v.FreeScore,
		})
	}
	return rst
}

func (impl reputationStoreImpl) ExportToFile(file string) {
	rst := impl.Export()
	f, err := os.Create(file)
	if err != nil {
		panic("failed to create account")
	}
	defer f.Close()
	jsonbytes, err := cdc.MarshalJSON(rst)
	f.Write(jsonbytes)
	if err != nil {
		panic("failed to marshal json for " + file + " due to " + err.Error())
	}
	f.Sync()
}


func (impl reputationStoreImpl) Import(tb *UserReputationTable) {
	for _, v := range tb.reputations {
		impl.setUserMeta(v.Username, &userMeta{
			FreeScore: v.FreeScore,
			CustomerScore: v.CustomerScore,
		})
	}
}

// TODO(yumin): a cache can help to make it faster.
// Follow https://github.com/golang/go/wiki/CodeReviewComments#declaring-empty-slices
// we prefer to use nil as empty slice in all following get/set.
// Another reason is that, gob does not differentiate nil and [], and it treats them all as
// nil.
func (impl reputationStoreImpl) getUserMeta(u Uid) *userMeta {
	buf := impl.store.Get(getUserMetaKey(u))
	rst := decodeUserMeta(buf)
	if rst == nil {
		return &userMeta{
			CustomerScore:     big.NewInt(InitialCustomerScore),
			FreeScore:         big.NewInt(0),
			LastSettled:       0,
			LastDonationRound: 0,
		}
	}
	return rst
}

func (impl reputationStoreImpl) setUserMeta(u Uid, data *userMeta) {
	if data != nil {
		impl.store.Set(getUserMetaKey(u), encodeUserMeta(data))
	}
}

func (impl reputationStoreImpl) getPostMeta(p Pid) *postMeta {
	buf := impl.store.Get(getPostMetaKey(p))
	rst := decodePostMeta(buf)
	if rst == nil {
		return &postMeta{
			SumRep: big.NewInt(0),
		}
	}
	return rst
}

func (impl reputationStoreImpl) setPostMeta(p Pid, post *postMeta) {
	if post != nil {
		impl.store.Set(getPostMetaKey(p), encodePostMeta(post))
	}
}

func (impl reputationStoreImpl) getRoundMeta(r RoundId) *roundMeta {
	buf := impl.store.Get(getRoundMetaKey(r))
	rst := decodeRoundMeta(buf)
	if rst == nil {
		return &roundMeta{
			Result:  nil,
			SumDp:   big.NewInt(0),
			StartAt: 0,
			TopN:    nil,
		}
	}
	return rst
}

func (impl reputationStoreImpl) setRoundMeta(r RoundId, dt *roundMeta) {
	if dt != nil {
		impl.store.Set(getRoundMetaKey(r), encodeRoundMeta(dt))
	}
}

func (impl reputationStoreImpl) getUserPostMeta(u Uid, p Pid) *userPostMeta {
	buf := impl.store.Get(getUserPostMetaKey(u, p))
	rst := decodeUserPostMeta(buf)
	if rst == nil {
		return &userPostMeta{
			Donated:         big.NewInt(0),
			LastDonationRep: big.NewInt(0),
			LastReportRep:   big.NewInt(0),
		}
	}
	return rst
}

func (impl reputationStoreImpl) setUserPostMeta(u Uid, p Pid, dt *userPostMeta) {
	if dt != nil {
		impl.store.Set(getUserPostMetaKey(u, p), encodeUserPostMeta(dt))
	}
}

func (impl reputationStoreImpl) getRoundPostMeta(r RoundId, p Pid) *roundPostMeta {
	buf := impl.store.Get(getRoundPostMetaKey(r, p))
	rst := decodeRoundPostMeta(buf)
	if rst == nil {
		return &roundPostMeta{
			SumDp:       big.NewInt(0),
			SumStake:    big.NewInt(0),
			NumKeysSold: big.NewInt(0),
		}
	}
	return rst
}

func (impl reputationStoreImpl) setRoundPostMeta(r RoundId, p Pid, dt *roundPostMeta) {
	if dt != nil {
		impl.store.Set(getRoundPostMetaKey(r, p), encodeRoundPostMeta(dt))
	}
}

func (impl reputationStoreImpl) getRoundUserPostMeta(r RoundId, u Uid, p Pid) *roundUserPostMeta {
	buf := impl.store.Get(getRoundUserPostMetaKey(r, u, p))
	rst := decodeRoundUserPostMeta(buf)
	if rst == nil {
		return &roundUserPostMeta{
			NumKeysHas: big.NewInt(0),
		}
	}
	return rst
}

func (impl reputationStoreImpl) setRoundUserPostMeta(r RoundId, u Uid, p Pid, dt *roundUserPostMeta) {
	if dt != nil {
		impl.store.Set(getRoundUserPostMetaKey(r, u, p), encodeRoundUserPostMeta(dt))
	}
}

// game starts with 1
func (impl reputationStoreImpl) getGameMeta() *gameMeta {
	buf := impl.store.Get(getGameKey())
	rst := decodeGameMeta(buf)
	if rst == nil {
		return &gameMeta{
			CurrentRound: 1,
		}
	}
	return rst
}

func (impl reputationStoreImpl) setGameMeta(dt *gameMeta) {
	if dt != nil {
		impl.store.Set(getGameKey(), encodeGameMeta(dt))
	}
}

//  --------------          user meta          ------------------
func (impl reputationStoreImpl) GetCustomerScore(u Uid) Rep {
	return impl.getUserMeta(u).CustomerScore
}

func (impl reputationStoreImpl) SetCustomerScore(u Uid, score Rep) {
	user := impl.getUserMeta(u)
	user.CustomerScore = score
	impl.setUserMeta(u, user)
}

func (impl reputationStoreImpl) GetFreeScore(u Uid) Rep {
	return impl.getUserMeta(u).FreeScore
}

func (impl reputationStoreImpl) SetFreeScore(u Uid, score Rep) {
	user := impl.getUserMeta(u)
	user.FreeScore = score
	impl.setUserMeta(u, user)
}

func (impl reputationStoreImpl) GetUserLastSettled(u Uid) RoundId {
	return impl.getUserMeta(u).LastSettled
}

func (impl reputationStoreImpl) SetUserLastSettled(u Uid, roundId RoundId) {
	user := impl.getUserMeta(u)
	user.LastSettled = roundId
	impl.setUserMeta(u, user)
}

func (impl reputationStoreImpl) GetUserLastDonationRound(u Uid) RoundId {
	return impl.getUserMeta(u).LastDonationRound
}

func (impl reputationStoreImpl) SetUserLastDonationRound(u Uid, roundId RoundId) {
	user := impl.getUserMeta(u)
	user.LastDonationRound = roundId
	impl.setUserMeta(u, user)
}

//  --------------          post meta          ------------------
func (impl reputationStoreImpl) GetSumRep(p Pid) Rep {
	rst := impl.getPostMeta(p)
	return rst.SumRep
}

// This function update the topN posts.
func (impl reputationStoreImpl) SetSumRep(p Pid, r Rep) {
	rst := impl.getPostMeta(p)
	rst.SumRep = r
	impl.setPostMeta(p, rst)
}

//  --------------     user post meta        ------------------

func (impl reputationStoreImpl) GetUserDonatedOn(u Uid, p Pid) Dp {
	rst := impl.getUserPostMeta(u, p)
	return rst.Donated
}

func (impl reputationStoreImpl) SetUserDonatedOn(u Uid, p Pid, dp Dp) {
	rst := impl.getUserPostMeta(u, p)
	rst.Donated = dp
	impl.setUserPostMeta(u, p, rst)
}

func (impl reputationStoreImpl) GetUserLastDonation(u Uid, p Pid) Rep {
	rst := impl.getUserPostMeta(u, p)
	return rst.LastDonationRep

}

func (impl reputationStoreImpl) SetUserLastDonation(u Uid, p Pid, rep Rep) {
	rst := impl.getUserPostMeta(u, p)
	rst.LastDonationRep = rep
	impl.setUserPostMeta(u, p, rst)

}

func (impl reputationStoreImpl) GetUserLastReport(u Uid, p Pid) Rep {
	rst := impl.getUserPostMeta(u, p)
	return rst.LastReportRep

}

func (impl reputationStoreImpl) SetUserLastReport(u Uid, p Pid, rep Rep) {
	rst := impl.getUserPostMeta(u, p)
	rst.LastReportRep = rep
	impl.setUserPostMeta(u, p, rst)
}

//  --------------     round meta        ------------------
func (impl reputationStoreImpl) GetRoundResult(r RoundId) []Pid {
	rst := impl.getRoundMeta(r)
	return rst.Result
}

func (impl reputationStoreImpl) SetRoundResult(r RoundId, bests []Pid) {
	rst := impl.getRoundMeta(r)
	rst.Result = bests
	impl.setRoundMeta(r, rst)
}

func (impl reputationStoreImpl) GetRoundSumDp(r RoundId) Dp {
	rst := impl.getRoundMeta(r)
	return rst.SumDp
}

func (impl reputationStoreImpl) SetRoundSumDp(r RoundId, dp Dp) {
	rst := impl.getRoundMeta(r)
	rst.SumDp = dp
	impl.setRoundMeta(r, rst)
}

func (impl reputationStoreImpl) GetRoundTopNPosts(r RoundId) []PostDpPair {
	rst := impl.getRoundMeta(r)
	return rst.TopN
}

func (impl reputationStoreImpl) GetRoundStartAt(r RoundId) Time {
	rst := impl.getRoundMeta(r)
	return rst.StartAt
}

func (impl reputationStoreImpl) GetCurrentRound() RoundId {
	rst := impl.getGameMeta()
	return rst.CurrentRound
}

func (impl reputationStoreImpl) StartNewRound(t Time) {
	rst := impl.getGameMeta()
	newRoundId := rst.CurrentRound + 1

	newRoundMeta := &roundMeta{
		Result:  nil,
		SumDp:   big.NewInt(0),
		StartAt: t,
		TopN:    nil,
	}
	impl.setRoundMeta(newRoundId, newRoundMeta)

	rst.CurrentRound = newRoundId
	impl.setGameMeta(rst)
}

//  ----------------     round post        ------------------
func (impl reputationStoreImpl) GetRoundPostSumDp(r RoundId, p Pid) Dp {
	rst := impl.getRoundPostMeta(r, p)
	return rst.SumDp
}

// XXX(yumin): this updates the topN.
func (impl reputationStoreImpl) SetRoundPostSumDp(r RoundId, p Pid, dp Dp) {
	// update sumDp
	rst := impl.getRoundPostMeta(r, p)
	rst.SumDp = dp
	impl.setRoundPostMeta(r, p, rst)

	// update topN
	roundMeta := impl.getRoundMeta(r)
	topN := roundMeta.TopN
	alreadyInTop := false
	for i, candidate := range topN {
		if candidate.Pid == p {
			alreadyInTop = true
			topN[i].SumDp = dp
		}
	}
	if !alreadyInTop {
		topN = append(topN, PostDpPair{Pid: p, SumDp: dp})
	}

	sort.SliceStable(topN, func(i, j int) bool {
		return topN[i].SumDp.Cmp(topN[j].SumDp) > 0
	})

	if len(topN) > impl.BestContentIndexN {
		topN = topN[:impl.BestContentIndexN]
	}
	roundMeta.TopN = topN
	impl.setRoundMeta(r, roundMeta)
}

func (impl reputationStoreImpl) GetRoundPostSumStake(r RoundId, p Pid) Stake {
	rst := impl.getRoundPostMeta(r, p)
	return rst.SumStake
}

func (impl reputationStoreImpl) SetRoundPostSumStake(r RoundId, p Pid, s Stake) {
	rst := impl.getRoundPostMeta(r, p)
	rst.SumStake = s
	impl.setRoundPostMeta(r, p, rst)
}

func (impl reputationStoreImpl) GetRoundNumKeysSold(r RoundId, p Pid) bigInt {
	rst := impl.getRoundPostMeta(r, p)
	return rst.NumKeysSold
}

func (impl reputationStoreImpl) SetRoundNumKeysSold(r RoundId, p Pid, numKeys bigInt) {
	rst := impl.getRoundPostMeta(r, p)
	rst.NumKeysSold = numKeys
	impl.setRoundPostMeta(r, p, rst)
}

//  ----------------     round user post        ------------------
func (impl reputationStoreImpl) GetRoundNumKeysHas(r RoundId, u Uid, p Pid) bigInt {
	rst := impl.getRoundUserPostMeta(r, u, p)
	return rst.NumKeysHas
}

func (impl reputationStoreImpl) SetRoundNumKeysHas(r RoundId, u Uid, p Pid, numKeys bigInt) {
	rst := impl.getRoundUserPostMeta(r, u, p)
	rst.NumKeysHas = numKeys
	impl.setRoundUserPostMeta(r, u, p, rst)
}

var _ ReputationStore = &reputationStoreImpl{}
