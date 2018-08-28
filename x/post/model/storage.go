package model

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	postInfoSubStore           = []byte{0x00} // SubStore for all post info
	postMetaSubStore           = []byte{0x01} // SubStore for all post mata info
	postReportOrUpvoteSubStore = []byte{0x02} // SubStore for all report or upvote to post
	postCommentSubStore        = []byte{0x03} // SubStore for all comments
	postViewsSubStore          = []byte{0x04} // SubStore for all views
	postDonationsSubStore      = []byte{0x05} // SubStore for all donations
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
	cdc := wire.NewCodec()
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
	if err := ps.cdc.UnmarshalJSON(infoByte, postInfo); err != nil {
		return nil, ErrFailedToUnmarshalPostInfo(err)
	}
	return postInfo, nil
}

// SetPostInfo - set post info to KVStore
func (ps PostStorage) SetPostInfo(ctx sdk.Context, postInfo *PostInfo) sdk.Error {
	store := ctx.KVStore(ps.key)
	infoByte, err := ps.cdc.MarshalJSON(*postInfo)
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
	if unmarshalErr := ps.cdc.UnmarshalJSON(metaBytes, postMeta); unmarshalErr != nil {
		return nil, ErrFailedToUnmarshalPostMeta(unmarshalErr)
	}
	return postMeta, nil
}

// SetPostMeta - set post meta to KVStore
func (ps PostStorage) SetPostMeta(ctx sdk.Context, permlink types.Permlink, postMeta *PostMeta) sdk.Error {
	store := ctx.KVStore(ps.key)
	metaBytes, err := ps.cdc.MarshalJSON(*postMeta)
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
	if unmarshalErr := ps.cdc.UnmarshalJSON(reportOrUpvoteBytes, reportOrUpvote); unmarshalErr != nil {
		return nil, ErrFailedToUnmarshalPostReportOrUpvote(unmarshalErr)
	}
	return reportOrUpvote, nil
}

// SetPostReportOrUpvote - set report or upvote to KVStore
func (ps PostStorage) SetPostReportOrUpvote(
	ctx sdk.Context, permlink types.Permlink, reportOrUpvote *ReportOrUpvote) sdk.Error {
	store := ctx.KVStore(ps.key)
	reportOrUpvoteByte, err := ps.cdc.MarshalJSON(*reportOrUpvote)
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
	if unmarshalErr := ps.cdc.UnmarshalJSON(commentBytes, postComment); unmarshalErr != nil {
		return nil, ErrFailedToUnmarshalPostComment(unmarshalErr)
	}
	return postComment, nil
}

// SetPostComment - set post comment to KVStore
func (ps PostStorage) SetPostComment(
	ctx sdk.Context, permlink types.Permlink, postComment *Comment) sdk.Error {
	store := ctx.KVStore(ps.key)
	postCommentByte, err := ps.cdc.MarshalJSON(*postComment)
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
	if unmarshalErr := ps.cdc.UnmarshalJSON(viewBytes, postView); unmarshalErr != nil {
		return nil, ErrFailedToUnmarshalPostView(unmarshalErr)
	}
	return postView, nil
}

// SetPostView - set post view to KVStore
func (ps PostStorage) SetPostView(ctx sdk.Context, permlink types.Permlink, postView *View) sdk.Error {
	store := ctx.KVStore(ps.key)
	postViewByte, err := ps.cdc.MarshalJSON(*postView)
	if err != nil {
		return ErrFailedToMarshalPostView(err)
	}
	store.Set(getPostViewKey(permlink, postView.Username), postViewByte)
	return nil
}

// GetPostDonations - get post donations from KVStore
func (ps PostStorage) GetPostDonations(
	ctx sdk.Context, permlink types.Permlink, donateUser types.AccountKey) (*Donations, sdk.Error) {
	store := ctx.KVStore(ps.key)
	donateBytes := store.Get(getPostDonationKey(permlink, donateUser))
	if donateBytes == nil {
		return nil, ErrPostDonationNotFound(getPostDonationKey(permlink, donateUser))
	}
	postDonations := new(Donations)
	if unmarshalErr := ps.cdc.UnmarshalJSON(donateBytes, postDonations); unmarshalErr != nil {
		return nil, ErrFailedToUnmarshalPostDonations(unmarshalErr)
	}
	return postDonations, nil
}

// SetPostDonations - set post donations to KVStore
func (ps PostStorage) SetPostDonations(
	ctx sdk.Context, permlink types.Permlink, postDonations *Donations) sdk.Error {
	store := ctx.KVStore(ps.key)
	postDonationsByte, err := ps.cdc.MarshalJSON(*postDonations)
	if err != nil {
		return ErrFailedToMarshalPostDonations(err)
	}
	store.Set(getPostDonationKey(permlink, postDonations.Username), postDonationsByte)
	return nil
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

// PostCommentPrefix - "donation substore" + "permlink"
// which can be used to access all donations belong to this post
func getPostDonationsPrefix(permlink types.Permlink) []byte {
	return append(append(postDonationsSubStore, permlink...), types.KeySeparator...)
}

// getPostDonationKey - "donation substore" + "permlink" + "donator"
func getPostDonationKey(permlink types.Permlink, donateUser types.AccountKey) []byte {
	return append(getPostDonationsPrefix(permlink), donateUser...)
}
