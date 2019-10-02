package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/post/types"
)

var (
	PostSubStore = []byte{0x00} // SubStore for all post info
)

func GetAuthorPrefix(author linotypes.AccountKey) []byte {
	return append(PostSubStore, author...)
}

// GetPostInfoKey - "post info substore" + "permlink"
func GetPostInfoKey(permlink linotypes.Permlink) []byte {
	return append(PostSubStore, permlink...)
}

// PostStorage - post storage
type PostStorage struct {
	key sdk.StoreKey
	cdc *wire.Codec
}

// NewPostStorage - returns a new PostStorage that
// uses codec to (binary) encode and decode concrete Post
func NewPostStorage(key sdk.StoreKey) PostStorage {
	cdc := wire.New()
	wire.RegisterCrypto(cdc)

	return PostStorage{
		key: key,
		cdc: cdc,
	}
}

// DoesPostExist - check if a post exists in KVStore or not
func (ps PostStorage) HasPost(ctx sdk.Context, permlink linotypes.Permlink) bool {
	store := ctx.KVStore(ps.key)
	return store.Has(GetPostInfoKey(permlink))
}

// GetPostInfo - get post info from KVStore
func (ps PostStorage) GetPost(ctx sdk.Context, permlink linotypes.Permlink) (*Post, sdk.Error) {
	store := ctx.KVStore(ps.key)
	key := GetPostInfoKey(permlink)
	infoByte := store.Get(key)
	if infoByte == nil {
		return nil, types.ErrPostNotFound(permlink)
	}
	postInfo := new(Post)
	ps.cdc.MustUnmarshalBinaryLengthPrefixed(infoByte, postInfo)
	return postInfo, nil
}

// SetPostInfo - set post info to KVStore
func (ps PostStorage) SetPost(ctx sdk.Context, postInfo *Post) {
	store := ctx.KVStore(ps.key)
	infoByte := ps.cdc.MustMarshalBinaryLengthPrefixed(*postInfo)
	store.Set(GetPostInfoKey(linotypes.GetPermlink(postInfo.Author, postInfo.PostID)), infoByte)
}

// Post cannot be deleted in the store. you can mark it as deleted.
// // SetPostInfo - set post info to KVStore
// func (ps PostStorage) DeletePost(ctx sdk.Context, permlink linotypes.Permlink) {
// 	store := ctx.KVStore(ps.key)
// 	store.Delete(GetPostInfoKey(permlink))
// }

func (ps PostStorage) StoreList(ctx sdk.Context) utils.StoreList {
	store := ctx.KVStore(ps.key)
	return utils.StoreList{
		string(PostSubStore): utils.SubStore{
			Store:      store,
			Prefix:     PostSubStore,
			ValCreator: func() interface{} { return new(Post) },
			Decoder:    ps.cdc.MustUnmarshalBinaryLengthPrefixed,
		},
	}
}
