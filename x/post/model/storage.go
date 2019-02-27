package model

import (
	"strings"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

var (
	postInfoSubStore           = []byte{0x00} // SubStore for all post info
	postMetaSubStore           = []byte{0x01} // SubStore for all post mata info
	postReportOrUpvoteSubStore = []byte{0x02} // SubStore for all report or upvote to post
	postCommentSubStore        = []byte{0x03} // SubStore for all comments
	postViewsSubStore          = []byte{0x04} // SubStore for all views
	// XXX(yukai): deprecated.
	// postDonationsSubStore      = []byte{0x05} // SubStore for all donations
)

// PostStorage - post storage
type PostStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
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
func (ps PostStorage) DoesPostExist(ctx sdk.Context, permlink types.Permlink) bool {
	store := ctx.KVStore(ps.key)
	return store.Has(GetPostInfoKey(permlink))
}

// GetPostInfo - get post info from KVStore
func (ps PostStorage) GetPostInfo(ctx sdk.Context, permlink types.Permlink) (*PostInfo, sdk.Error) {
	store := ctx.KVStore(ps.key)
	infoByte := store.Get(GetPostInfoKey(permlink))
	if infoByte == nil {
		return nil, ErrPostNotFound(GetPostInfoKey(permlink))
	}
	postInfo := new(PostInfo)
	if err := ps.cdc.UnmarshalBinaryBare(infoByte, postInfo); err != nil {
		return nil, ErrFailedToUnmarshalPostInfo(err)
	}
	return postInfo, nil
}

// SetPostInfo - set post info to KVStore
func (ps PostStorage) SetPostInfo(ctx sdk.Context, postInfo *PostInfo) sdk.Error {
	store := ctx.KVStore(ps.key)
	infoByte, err := ps.cdc.MarshalBinaryBare(*postInfo)
	if err != nil {
		return ErrFailedToMarshalPostInfo(err)
	}
	store.Set(GetPostInfoKey(types.GetPermlink(postInfo.Author, postInfo.PostID)), infoByte)
	return nil
}

// GetPostMeta - get post meta from KVStore
func (ps PostStorage) GetPostMeta(ctx sdk.Context, permlink types.Permlink) (*PostMeta, sdk.Error) {
	store := ctx.KVStore(ps.key)
	metaBytes := store.Get(GetPostMetaKey(permlink))
	if metaBytes == nil {
		return nil, ErrPostMetaNotFound(GetPostMetaKey(permlink))
	}
	postMeta := new(PostMeta)
	if unmarshalErr := ps.cdc.UnmarshalBinaryBare(metaBytes, postMeta); unmarshalErr != nil {
		return nil, ErrFailedToUnmarshalPostMeta(unmarshalErr)
	}
	return postMeta, nil
}

// SetPostMeta - set post meta to KVStore
func (ps PostStorage) SetPostMeta(ctx sdk.Context, permlink types.Permlink, postMeta *PostMeta) sdk.Error {
	store := ctx.KVStore(ps.key)
	metaBytes, err := ps.cdc.MarshalBinaryBare(*postMeta)
	if err != nil {
		return ErrFailedToMarshalPostMeta(err)
	}
	store.Set(GetPostMetaKey(permlink), metaBytes)
	return nil
}

// GetPostReportOrUpvote - get report or upvote from KVStore
func (ps PostStorage) GetPostReportOrUpvote(
	ctx sdk.Context, permlink types.Permlink, user types.AccountKey) (*ReportOrUpvote, sdk.Error) {
	store := ctx.KVStore(ps.key)
	reportOrUpvoteBytes := store.Get(getPostReportOrUpvoteKey(permlink, user))
	if reportOrUpvoteBytes == nil {
		return nil, ErrPostReportOrUpvoteNotFound(getPostReportOrUpvoteKey(permlink, user))
	}
	reportOrUpvote := new(ReportOrUpvote)
	if unmarshalErr := ps.cdc.UnmarshalBinaryBare(reportOrUpvoteBytes, reportOrUpvote); unmarshalErr != nil {
		return nil, ErrFailedToUnmarshalPostReportOrUpvote(unmarshalErr)
	}
	return reportOrUpvote, nil
}

// SetPostReportOrUpvote - set report or upvote to KVStore
func (ps PostStorage) SetPostReportOrUpvote(
	ctx sdk.Context, permlink types.Permlink, reportOrUpvote *ReportOrUpvote) sdk.Error {
	store := ctx.KVStore(ps.key)
	reportOrUpvoteByte, err := ps.cdc.MarshalBinaryBare(*reportOrUpvote)
	if err != nil {
		return ErrFailedToMarshalPostReportOrUpvote(err)
	}
	store.Set(getPostReportOrUpvoteKey(permlink, reportOrUpvote.Username), reportOrUpvoteByte)
	return nil
}

// GetPostComment - get post comment from KVStore
func (ps PostStorage) GetPostComment(
	ctx sdk.Context, permlink types.Permlink, commentPermlink types.Permlink) (*Comment, sdk.Error) {
	store := ctx.KVStore(ps.key)
	commentBytes := store.Get(getPostCommentKey(permlink, commentPermlink))
	if commentBytes == nil {
		return nil, ErrPostCommentNotFound(getPostCommentKey(permlink, commentPermlink))
	}
	postComment := new(Comment)
	if unmarshalErr := ps.cdc.UnmarshalBinaryBare(commentBytes, postComment); unmarshalErr != nil {
		return nil, ErrFailedToUnmarshalPostComment(unmarshalErr)
	}
	return postComment, nil
}

// SetPostComment - set post comment to KVStore
func (ps PostStorage) SetPostComment(
	ctx sdk.Context, permlink types.Permlink, postComment *Comment) sdk.Error {
	store := ctx.KVStore(ps.key)
	postCommentByte, err := ps.cdc.MarshalBinaryBare(*postComment)
	if err != nil {
		return ErrFailedToMarshalPostComment(err)
	}
	store.Set(
		getPostCommentKey(permlink, types.GetPermlink(postComment.Author, postComment.PostID)),
		postCommentByte)
	return nil
}

// GetPostView - get post view from KVStore
func (ps PostStorage) GetPostView(
	ctx sdk.Context, permlink types.Permlink, viewUser types.AccountKey) (*View, sdk.Error) {
	store := ctx.KVStore(ps.key)
	viewBytes := store.Get(getPostViewKey(permlink, viewUser))
	if viewBytes == nil {
		return nil, ErrPostViewNotFound(getPostViewKey(permlink, viewUser))
	}
	postView := new(View)
	if unmarshalErr := ps.cdc.UnmarshalBinaryBare(viewBytes, postView); unmarshalErr != nil {
		return nil, ErrFailedToUnmarshalPostView(unmarshalErr)
	}
	return postView, nil
}

// SetPostView - set post view to KVStore
func (ps PostStorage) SetPostView(ctx sdk.Context, permlink types.Permlink, postView *View) sdk.Error {
	store := ctx.KVStore(ps.key)
	postViewByte, err := ps.cdc.MarshalBinaryBare(*postView)
	if err != nil {
		return ErrFailedToMarshalPostView(err)
	}
	store.Set(getPostViewKey(permlink, postView.Username), postViewByte)
	return nil
}

// Export post storage state.
func (ps PostStorage) Export(ctx sdk.Context) *PostTables {
	tables := &PostTables{}
	store := ctx.KVStore(ps.key)
	// export table.Posts
	func() {
		itr := sdk.KVStorePrefixIterator(store, postInfoSubStore)
		defer itr.Close()
		for ; itr.Valid(); itr.Next() {
			k := itr.Key()
			permlink := types.Permlink(k[1:])
			info, err := ps.GetPostInfo(ctx, permlink)
			if err != nil {
				panic("failed to read post info: " + err.Error())
			}
			meta, err := ps.GetPostMeta(ctx, permlink)
			if err != nil {
				panic("failed to read post meta: " + err.Error())
			}
			row := PostRow{
				Permlink: permlink,
				Info:     *info,
				Meta:     *meta,
			}
			tables.Posts = append(tables.Posts, row)
		}
	}()
	// export tables.PostUser
	func() {
		itr := sdk.KVStorePrefixIterator(store, postReportOrUpvoteSubStore)
		defer itr.Close()
		for ; itr.Valid(); itr.Next() {
			k := itr.Key()
			permlinkAccount := string(k[1:])
			strs := strings.Split(permlinkAccount, types.KeySeparator)
			if len(strs) != 2 {
				panic("failed to split out permlink account: " + permlinkAccount)
			}
			permlink, username := types.Permlink(strs[0]), types.AccountKey(strs[1])
			ru, err := ps.GetPostReportOrUpvote(ctx, permlink, username)
			if err != nil {
				panic("failed to get report or upvote: " + err.Error())
			}
			row := PostUserRow{
				Permlink:       permlink,
				User:           username,
				ReportOrUpvote: *ru,
			}
			tables.PostUsers = append(tables.PostUsers, row)
		}
	}()
	return tables
}

// Import from tablesIR.
func (ps PostStorage) Import(ctx sdk.Context, tb *PostTablesIR) {
	check := func(e error) {
		if e != nil {
			panic("[ps] Failed to import: " + e.Error())
		}
	}
	// import table.developers
	for _, v := range tb.Posts {
		err := ps.SetPostInfo(ctx, &v.Info)
		check(err)
		err = ps.SetPostMeta(ctx, v.Permlink, &PostMeta{
			CreatedAt:               v.Meta.CreatedAt,
			LastUpdatedAt:           v.Meta.LastUpdatedAt,
			LastActivityAt:          v.Meta.LastActivityAt,
			AllowReplies:            v.Meta.AllowReplies,
			IsDeleted:               v.Meta.IsDeleted,
			TotalDonateCount:        v.Meta.TotalDonateCount,
			TotalReportCoinDay:      v.Meta.TotalReportCoinDay,
			TotalUpvoteCoinDay:      v.Meta.TotalUpvoteCoinDay,
			TotalViewCount:          v.Meta.TotalViewCount,
			TotalReward:             v.Meta.TotalReward,
			RedistributionSplitRate: sdk.MustNewDecFromStr(v.Meta.RedistributionSplitRate),
		})
		check(err)
	}
	// import PostUsers
	for _, v := range tb.PostUsers {
		err := ps.SetPostReportOrUpvote(ctx, v.Permlink, &v.ReportOrUpvote)
		check(err)
	}
}

// GetPostInfoPrefix - "post info substore" + "author"
func GetPostInfoPrefix(author types.AccountKey) []byte {
	return append(postInfoSubStore, author...)
}

// GetPostInfoKey - "post info substore" + "permlink"
func GetPostInfoKey(permlink types.Permlink) []byte {
	return append(postInfoSubStore, permlink...)
}

// GetPostMetaKey - "post meta substore" + "permlink"
func GetPostMetaKey(permlink types.Permlink) []byte {
	return append(postMetaSubStore, permlink...)
}

// getPostReportOrUpvotePrefix - "post report or upvote substore" + "permlink"
// which can be used to access all reports belong to this post
func getPostReportOrUpvotePrefix(permlink types.Permlink) []byte {
	return append(append(postReportOrUpvoteSubStore, permlink...), types.KeySeparator...)
}

// getPostReportOrUpvotePrefix - "post report or upvote substore" + "permlink" + "user"
func getPostReportOrUpvoteKey(permlink types.Permlink, user types.AccountKey) []byte {
	return append(getPostReportOrUpvotePrefix(permlink), user...)
}

// getPostViewPrefix - "post view substore" + "permlink"
// which can be used to access all views belong to this post
func getPostViewPrefix(permlink types.Permlink) []byte {
	return append(append(postViewsSubStore, permlink...), types.KeySeparator...)
}

// getPostViewKey - "post view substore" + "permlink" + "user"
func getPostViewKey(permlink types.Permlink, viewUser types.AccountKey) []byte {
	return append(getPostViewPrefix(permlink), viewUser...)
}

// PostCommentPrefix - "comment substore" + "permlink"
// which can be used to access all comments belong to this post
func getPostCommentPrefix(permlink types.Permlink) []byte {
	return append(append(postCommentSubStore, permlink...), types.KeySeparator...)
}

// PostCommentPrefix - "comment substore" + "permlink" + "comment permlink"
func getPostCommentKey(permlink types.Permlink, commentPermlink types.Permlink) []byte {
	return append(getPostCommentPrefix(permlink), commentPermlink...)
}
