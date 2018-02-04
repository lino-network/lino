package types

import (
	"fmt"
	"github.com/tendermint/go-wire"
	"bytes"
)

type Like struct {
	From     []byte           `json:"from"`      // address
	To       []byte           `json:"to"`        // post_id
}
type LikeId []byte

func (this *Like) EqualsTo(other *Like) bool {
	return bytes.Equal(this.From, other.From) && bytes.Equal(this.From, other.From)
}

type LikeSummary struct {
	Likes    []LikeId
}
// const LIKE_SUMMARY_KEY [...]byte("LikeSummaryKey")
const (
	LIKE_SUMMARY_KEY = "LIKE_SUMMARY_KEY"
)

func (this *Like) String() string {
	return fmt.Sprintf("Like{%v, %v}", this.From, this.To)
}

func (this Like) ToLikeID() LikeId {
	return wire.BinaryBytes("like/" + string(this.From) + "#" + string(this.To))
}

// TODO(djj): disallow invalid Like in CheckTx, instead of leave it as no-op.

func GetLikesByPostId(store KVStore, post_id []byte) []Like {
	summary := readLikeSummary(store)
	var likes []Like
	for _, like_id := range summary.Likes {
		if like := readLike(store, like_id); bytes.Equal(like.To, post_id) {
			likes = append(likes, *like)
		}
	}
	return likes
}

func AddLike(store KVStore, like Like) {
	if !doesPostExist(store, like.To) {
		return
	}
	summary := readLikeSummary(store)
	if likeExist(store, like, summary) {
		return
	}
	// insert like to db
	like_id := insertLike(store, &like)

	// update summary
	summary.Likes = append(summary.Likes, like_id)
	updateSummary(store, summary)
}

func RemoveLike(store KVStore, like Like) {
	if !doesPostExist(store, like.To) {
		return
	}
	summary := readLikeSummary(store)
	if !likeExist(store, like, summary) {
		return
	}
	// remove id to save some space? lol
	like_id := like.ToLikeID()
	store.Set(like_id, []byte{})

	// update summary
	removeLikeFromSummary(summary, like_id)
	updateSummary(store, summary)
}

func removeLikeFromSummary(summary *LikeSummary, like_id LikeId) {
	for i, v := range summary.Likes {
		if (bytes.Equal(v, like_id)) {
			summary.Likes = append(summary.Likes[:i], summary.Likes[i + 1:]...)
			return
		}
	}
	panic("Removing a Like that does not exist in LikeSummary")
}

func updateSummary(store KVStore, summary *LikeSummary) {
	bytes := wire.BinaryBytes(summary)
	store.Set([]byte(LIKE_SUMMARY_KEY), bytes);
}

func insertLike(store KVStore, like *Like) LikeId {
	bytes := wire.BinaryBytes(like)
	like_id := like.ToLikeID()
	store.Set(like_id, bytes)
	return like_id
}

// func likeExist(store KVStore, to_insert Like) bool {
// 	summary := readLikeSummary(store)
// 	return likeExist(store, to_insert, summary)
// }

func likeExist(store KVStore, to_insert Like, summary *LikeSummary) bool {
	for _, like_id := range summary.Likes {
		like := readLike(store, like_id)
		if like.EqualsTo(&to_insert) {
			return true
		}
	}
	return false
}

func readLike(store KVStore, like_id LikeId) *Like {
	data := store.Get(like_id)
	if len(data) == 0 {
		return nil
	}
	var like *Like
	err := wire.ReadBinaryBytes(data, &like)
	if err != nil {
		panic("Calling ReadLike using an invalid like_id")
	}
	return like
}

func readLikeSummary(store KVStore) *LikeSummary {
	data := store.Get([]byte(LIKE_SUMMARY_KEY))
	if len(data) == 0 {
		return &LikeSummary{}
	}
	var summary *LikeSummary
	err := wire.ReadBinaryBytes(data, &summary)
	if err != nil {
		panic("ReadLikeSummary is corrupted.")
	}
	return summary
}

func doesPostExist(store KVStore, post_id []byte) bool {
	return store.Get(post_id) != nil
}