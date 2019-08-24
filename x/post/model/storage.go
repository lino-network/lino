package model

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/types"
)

var (
	postSubStore = []byte{0x00} // SubStore for all post info
)

func GetAuthorPrefix(author linotypes.AccountKey) []byte {
	return append(postSubStore, author...)
}

// GetPostInfoKey - "post info substore" + "permlink"
func GetPostInfoKey(permlink linotypes.Permlink) []byte {
	return append(postSubStore, permlink...)
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

// // Export post storage state.
func (ps PostStorage) Export(ctx sdk.Context) *PostTablesIR {
	panic("post export unimplemented")
	// tables := &PostTables{}
	// store := ctx.KVStore(ps.key)
	// // export table.Posts
	// func() {
	// 	itr := sdk.KVStorePrefixIterator(store, postSubStore)
	// 	defer itr.Close()
	// 	for ; itr.Valid(); itr.Next() {
	// 		k := itr.Key()
	// 		permlink := linotypes.Permlink(k[1:])
	// 		info, err := ps.GetPost(ctx, permlink)
	// 		if err != nil {
	// 			panic("failed to read post info: " + err.Error())
	// 		}
	// 		row := PostRow{
	// 			Permlink: permlink,
	// 			Info:     *info,
	// 			Meta:     *meta,
	// 		}
	// 		tables.Posts = append(tables.Posts, row)
	// 	}
	// }()
	// return tables
}

// Import from tablesIR.
func (ps PostStorage) Import(ctx sdk.Context, tb *PostTablesIR) error {
	// upgrade2 has simplied the post structure to just one post.
	for _, v := range tb.Posts {
		ps.SetPost(ctx, &Post{
			PostID:    v.Info.PostID,
			Title:     v.Info.Title,
			Content:   v.Info.Content,
			Author:    v.Info.Author,
			CreatedBy: v.Info.Author,
			CreatedAt: v.Meta.CreatedAt,
			UpdatedAt: v.Meta.LastUpdatedAt,
			IsDeleted: v.Meta.IsDeleted,
		})
	}
	return nil
}
