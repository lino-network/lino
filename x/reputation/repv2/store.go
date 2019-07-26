package repv2

import (
	"math/big"
	"strconv"

	db "github.com/tendermint/tendermint/libs/db"
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

type UserIterator = func(user Uid) bool // return true to break.

// ReputationStore - store reputation values.
// a simple wrapper around kv-store. It does not handle any reputation computation logic,
// except for init value for some field, like customer score.
// This interface should only be used by the reputation system implementation.
type ReputationStore interface {
	// TODO(yumin): these two are all in memory, which is extremely bad if state is large.
	// should changed to io.Write/Reader.
	// Export all state to deterministic bytes
	Export() *UserReputationTable

	// Import state from bytes
	Import(tb *UserReputationTable)

	// Iterator over usernames
	IterateUsers(UserIterator)

	GetUserMeta(u Uid) *userMeta
	SetUserMeta(u Uid, data *userMeta)

	SetRoundMeta(r RoundId, dt *roundMeta)
	GetRoundMeta(r RoundId) *roundMeta

	// total donation power received of a @p post in @p r round.
	GetRoundPostMeta(r RoundId, p Pid) *roundPostMeta
	SetRoundPostMeta(r RoundId, p Pid, dt *roundPostMeta)
	DelRoundPostMeta(r RoundId, p Pid)

	/// -----------  In this round  -------------
	// RoundId is the current round, starts from 1.
	GetCurrentRound() RoundId

	// Global data.
	GetGameMeta() *gameMeta
	SetGameMeta(dt *gameMeta)
}

// This store implementation does not have state. It is just a wrapper of read/write of
// data necessary for reputation system. Also, it takes the problem, that the underlying
// kv store may be an iavl, into account so the number of keys will have a large
// impact on performance into account, by trying to minimize the number of keys.
// we choosed amino json as serializer because it is deterministic.

var (
	KeySeparator           = byte('/')
	repUserMetaPrefix      = []byte{0x00}
	repRoundMetaPrefix     = []byte{0x01}
	repRoundPostMetaPrefix = []byte{0x02}
	repGameMetaPrefix      = []byte{0x03}
)

// LastSettled: till which round has user settled, i.e.
//     his customer score is accurate, initialized as 0.
// LastDonationRound:
//     The last round that user has any donation.
//     contract is that since any donation will force user to settle any previous customer score,
//     so if there would be any unsettled customer score, when lastDonationRound > LastSettled,
//     they are all in the lastDonationRound.
type userMeta struct {
	Consumption       Rep        `json:"cs"`
	Hold              Rep        `json:"hold"`
	Reputation        Rep        `json:"rep"`
	LastSettledRound  RoundId    `json:"ls"`
	LastDonationRound RoundId    `json:"ldr"`
	Unsettled         []Donation `json:"ust"`
}

// Result: the final result of round, should be set right after round ends.
// SumDp: sum of donation power of all donations happened during the round.
// top n posts. Ordering is maintained by caller. Note that this is not the final result.
type roundMeta struct {
	Result  []Pid        `json:"result"`
	SumIF   IF           `json:"sum_if"`
	StartAt Time         `json:"start_at"`
	TopN    []PostIFPair `json:"top_n"`
}

type roundPostMeta struct {
	SumIF IF `json:"dp"`
}

type gameMeta struct {
	CurrentRound RoundId `json:"current_round"`
}

func getUserMetaKey(u Uid) []byte {
	return append(repUserMetaPrefix, []byte(u)...)
}

func getRoundMetaKey(r RoundId) []byte {
	return append(repRoundMetaPrefix, strconv.FormatInt(r, 36)...)
}

func getRoundPostMetaKey(r RoundId, p Pid) []byte {
	prefix := append(repRoundPostMetaPrefix, strconv.FormatInt(r, 36)...)
	return append(append(prefix, KeySeparator), []byte(p)...)
}

func getGameKey() []byte {
	return repGameMetaPrefix
}

// no state.
type reputationStoreImpl struct {
	store          Store
	initReputation int64
}

func NewReputationStore(s Store, initRep int64) ReputationStore {
	return &reputationStoreImpl{store: s, initReputation: initRep}
}

func (impl reputationStoreImpl) IterateUsers(cb UserIterator) {
	itr := impl.store.Iterator(repUserMetaPrefix, PrefixEndBytes(repUserMetaPrefix))
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		uid := Uid(itr.Key()[1:])
		if cb(uid) {
			break
		}
	}
}

func (impl reputationStoreImpl) Export() *UserReputationTable {
	rst := &UserReputationTable{}
	impl.IterateUsers(func(uid Uid) bool {
		v := impl.GetUserMeta(uid)
		rst.Reputations = append(rst.Reputations, UserReputation{
			Username:      uid,
			CustomerScore: v.Reputation,
			FreeScore:     big.NewInt(0),
		})

		return false
	})
	return rst
}

// backward compatible to v1.
func (impl reputationStoreImpl) Import(tb *UserReputationTable) {
	for _, v := range tb.Reputations {
		impl.SetUserMeta(v.Username, &userMeta{
			Reputation: bigIntAdd(v.FreeScore, v.CustomerScore),
		})
	}
}

func (impl reputationStoreImpl) GetUserMeta(u Uid) *userMeta {
	buf := impl.store.Get(getUserMetaKey(u))
	rst := decodeUserMeta(buf)
	if rst == nil {
		return &userMeta{
			Consumption:       big.NewInt(impl.initReputation),
			Hold:              big.NewInt(0),
			Reputation:        big.NewInt(impl.initReputation),
			LastSettledRound:  0,
			LastDonationRound: 0,
			Unsettled:         nil,
		}
	}
	return rst
}

func (impl reputationStoreImpl) SetUserMeta(u Uid, data *userMeta) {
	if data != nil {
		impl.store.Set(getUserMetaKey(u), encodeUserMeta(data))
	}
}

func (impl reputationStoreImpl) GetRoundMeta(r RoundId) *roundMeta {
	buf := impl.store.Get(getRoundMetaKey(r))
	rst := decodeRoundMeta(buf)
	if rst == nil {
		return &roundMeta{
			Result:  nil,
			SumIF:   big.NewInt(0),
			StartAt: 0,
			TopN:    nil,
		}
	}
	return rst
}

func (impl reputationStoreImpl) SetRoundMeta(r RoundId, dt *roundMeta) {
	if dt != nil {
		impl.store.Set(getRoundMetaKey(r), encodeRoundMeta(dt))
	}
}

func (impl reputationStoreImpl) GetRoundPostMeta(r RoundId, p Pid) *roundPostMeta {
	buf := impl.store.Get(getRoundPostMetaKey(r, p))
	rst := decodeRoundPostMeta(buf)
	if rst == nil {
		return &roundPostMeta{
			SumIF: big.NewInt(0),
		}
	}
	return rst
}

func (impl reputationStoreImpl) SetRoundPostMeta(r RoundId, p Pid, dt *roundPostMeta) {
	if dt != nil {
		impl.store.Set(getRoundPostMetaKey(r, p), encodeRoundPostMeta(dt))
	}
}

func (impl reputationStoreImpl) DelRoundPostMeta(r RoundId, p Pid) {
	impl.store.Delete(getRoundPostMetaKey(r, p))
}

// game starts with 1
func (impl reputationStoreImpl) GetGameMeta() *gameMeta {
	buf := impl.store.Get(getGameKey())
	rst := decodeGameMeta(buf)
	if rst == nil {
		return &gameMeta{
			CurrentRound: 1,
		}
	}
	return rst
}

func (impl reputationStoreImpl) SetGameMeta(dt *gameMeta) {
	if dt != nil {
		impl.store.Set(getGameKey(), encodeGameMeta(dt))
	}
}

func (impl reputationStoreImpl) GetCurrentRound() RoundId {
	rst := impl.GetGameMeta()
	return rst.CurrentRound
}

var _ ReputationStore = &reputationStoreImpl{}
