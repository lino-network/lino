package model

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	postInfoSubStore           = []byte{0x00} // SubStore for all post info
	postMetaSubStore           = []byte{0x01} // SubStore for all post mata info
	postLikeSubStore           = []byte{0x02} // SubStore for all like to post
	postReportOrUpvoteSubStore = []byte{0x03} // SubStore for all like to post
	postCommentSubStore        = []byte{0x04} // SubStore for all comments
	postViewsSubStore          = []byte{0x05} // SubStore for all views
	postDonationsSubStore      = []byte{0x06} // SubStore for all donations
)

type PostStorage struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewPostStorage returns a new PostStorage that
// uses go-wire to (binary) encode and decode concrete Post
func NewPostStorage(key sdk.StoreKey) PostStorage {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)

	return PostStorage{
		key: key,
		cdc: cdc,
	}
}

func (ps PostStorage) DoesPostExist(ctx sdk.Context, permlink types.Permlink) bool {
	store := ctx.KVStore(ps.key)
	return store.Has(GetPostInfoKey(permlink))
}

func (ps PostStorage) GetPostInfo(ctx sdk.Context, permlink types.Permlink) (*PostInfo, sdk.Error) {
	store := ctx.KVStore(ps.key)
	infoByte := store.Get(GetPostInfoKey(permlink))
	if infoByte == nil {
		return nil, ErrPostNotFound(GetPostInfoKey(permlink))
	}
	postInfo := new(PostInfo)
	if err := ps.cdc.UnmarshalJSON(infoByte, postInfo); err != nil {
		return nil, ErrPostUnmarshalError(err)
	}
	return postInfo, nil
}

func (ps PostStorage) SetPostInfo(ctx sdk.Context, postInfo *PostInfo) sdk.Error {
	store := ctx.KVStore(ps.key)
	infoByte, err := ps.cdc.MarshalJSON(*postInfo)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostInfoKey(types.GetPermlink(postInfo.Author, postInfo.PostID)), infoByte)
	return nil
}

func (ps PostStorage) GetPostMeta(ctx sdk.Context, permlink types.Permlink) (*PostMeta, sdk.Error) {
	store := ctx.KVStore(ps.key)
	metaBytes := store.Get(GetPostMetaKey(permlink))
	if metaBytes == nil {
		return nil, ErrPostMetaNotFound(GetPostMetaKey(permlink))
	}
	postMeta := new(PostMeta)
	if unmarshalErr := ps.cdc.UnmarshalJSON(metaBytes, postMeta); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postMeta, nil
}

func (ps PostStorage) SetPostMeta(ctx sdk.Context, permlink types.Permlink, postMeta *PostMeta) sdk.Error {
	store := ctx.KVStore(ps.key)
	metaBytes, err := ps.cdc.MarshalJSON(*postMeta)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostMetaKey(permlink), metaBytes)
	return nil
}

func (ps PostStorage) GetPostLike(
	ctx sdk.Context, permlink types.Permlink, likeUser types.AccountKey) (*Like, sdk.Error) {
	store := ctx.KVStore(ps.key)
	likeBytes := store.Get(GetPostLikeKey(permlink, likeUser))
	if likeBytes == nil {
		return nil, ErrPostLikeNotFound(GetPostLikeKey(permlink, likeUser))
	}
	postLike := new(Like)
	if unmarshalErr := ps.cdc.UnmarshalJSON(likeBytes, postLike); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postLike, nil
}

func (ps PostStorage) SetPostLike(ctx sdk.Context, permlink types.Permlink, postLike *Like) sdk.Error {
	store := ctx.KVStore(ps.key)
	likeByte, err := ps.cdc.MarshalJSON(*postLike)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostLikeKey(permlink, postLike.Username), likeByte)
	return nil
}

func (ps PostStorage) GetPostReportOrUpvote(
	ctx sdk.Context, permlink types.Permlink, user types.AccountKey) (*ReportOrUpvote, sdk.Error) {
	store := ctx.KVStore(ps.key)
	reportOrUpvoteBytes := store.Get(GetPostReportOrUpvoteKey(permlink, user))
	if reportOrUpvoteBytes == nil {
		return nil, ErrPostReportOrUpvoteNotFound(GetPostReportOrUpvoteKey(permlink, user))
	}
	reportOrUpvote := new(ReportOrUpvote)
	if unmarshalErr := ps.cdc.UnmarshalJSON(reportOrUpvoteBytes, reportOrUpvote); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return reportOrUpvote, nil
}

func (ps PostStorage) SetPostReportOrUpvote(
	ctx sdk.Context, permlink types.Permlink, reportOrUpvote *ReportOrUpvote) sdk.Error {
	store := ctx.KVStore(ps.key)
	reportOrUpvoteByte, err := ps.cdc.MarshalJSON(*reportOrUpvote)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostReportOrUpvoteKey(permlink, reportOrUpvote.Username), reportOrUpvoteByte)
	return nil
}

func (ps PostStorage) GetPostComment(
	ctx sdk.Context, permlink types.Permlink, commentPermlink types.Permlink) (*Comment, sdk.Error) {
	store := ctx.KVStore(ps.key)
	commentBytes := store.Get(GetPostCommentKey(permlink, commentPermlink))
	if commentBytes == nil {
		return nil, ErrPostCommentNotFound(GetPostCommentKey(permlink, commentPermlink))
	}
	postComment := new(Comment)
	if unmarshalErr := ps.cdc.UnmarshalJSON(commentBytes, postComment); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postComment, nil
}

func (ps PostStorage) SetPostComment(
	ctx sdk.Context, permlink types.Permlink, postComment *Comment) sdk.Error {
	store := ctx.KVStore(ps.key)
	postCommentByte, err := ps.cdc.MarshalJSON(*postComment)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(
		GetPostCommentKey(permlink, types.GetPermlink(postComment.Author, postComment.PostID)),
		postCommentByte)
	return nil
}

func (ps PostStorage) GetPostView(
	ctx sdk.Context, permlink types.Permlink, viewUser types.AccountKey) (*View, sdk.Error) {
	store := ctx.KVStore(ps.key)
	viewBytes := store.Get(GetPostViewKey(permlink, viewUser))
	if viewBytes == nil {
		return nil, ErrPostViewNotFound(GetPostViewKey(permlink, viewUser))
	}
	postView := new(View)
	if unmarshalErr := ps.cdc.UnmarshalJSON(viewBytes, postView); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postView, nil
}

func (ps PostStorage) SetPostView(ctx sdk.Context, permlink types.Permlink, postView *View) sdk.Error {
	store := ctx.KVStore(ps.key)
	postViewByte, err := ps.cdc.MarshalJSON(*postView)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostViewKey(permlink, postView.Username), postViewByte)
	return nil
}

func (ps PostStorage) GetPostDonations(
	ctx sdk.Context, permlink types.Permlink, donateUser types.AccountKey) (*Donations, sdk.Error) {
	store := ctx.KVStore(ps.key)
	donateBytes := store.Get(GetPostDonationKey(permlink, donateUser))
	if donateBytes == nil {
		return nil, ErrPostDonationNotFound(GetPostDonationKey(permlink, donateUser))
	}
	postDonations := new(Donations)
	if unmarshalErr := ps.cdc.UnmarshalJSON(donateBytes, postDonations); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postDonations, nil
}

func (ps PostStorage) SetPostDonations(
	ctx sdk.Context, permlink types.Permlink, postDonations *Donations) sdk.Error {
	store := ctx.KVStore(ps.key)
	postDonationsByte, err := ps.cdc.MarshalJSON(*postDonations)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostDonationKey(permlink, postDonations.Username), postDonationsByte)
	return nil
}

func GetPostInfoKey(permlink types.Permlink) []byte {
	return append(postInfoSubStore, permlink...)
}

func GetPostMetaKey(permlink types.Permlink) []byte {
	return append(postMetaSubStore, permlink...)
}

// PostLikePrefix format is LikeSubStore / PostKey
// which can be used to access all likes belong to this post
func getPostLikePrefix(permlink types.Permlink) []byte {
	return append(append(postLikeSubStore, permlink...), types.KeySeparator...)
}

func GetPostLikeKey(permlink types.Permlink, likeUser types.AccountKey) []byte {
	return append(getPostLikePrefix(permlink), likeUser...)
}

// PostReportPrefix format is ReportSubStore / PostKey
// which can be used to access all reports belong to this post
func getPostReportOrUpvotePrefix(permlink types.Permlink) []byte {
	return append(append(postReportOrUpvoteSubStore, permlink...), types.KeySeparator...)
}

func GetPostReportOrUpvoteKey(permlink types.Permlink, user types.AccountKey) []byte {
	return append(getPostReportOrUpvotePrefix(permlink), user...)
}

// PostViewPrefix format is ViewSubStore / permlink
// which can be used to access all views belong to this post
func getPostViewPrefix(permlink types.Permlink) []byte {
	return append(append(postViewsSubStore, permlink...), types.KeySeparator...)
}

func GetPostViewKey(permlink types.Permlink, viewUser types.AccountKey) []byte {
	return append(getPostViewPrefix(permlink), viewUser...)
}

// PostCommentPrefix format is CommentSubStore / permlink
// which can be used to access all comments belong to this post
func getPostCommentPrefix(permlink types.Permlink) []byte {
	return append(append(postCommentSubStore, permlink...), types.KeySeparator...)
}

func GetPostCommentKey(permlink types.Permlink, commentPermlink types.Permlink) []byte {
	return append(getPostCommentPrefix(permlink), commentPermlink...)
}

// PostDonationPrefix format is DonationSubStore / permlink
// which can be used to access all donations belong to this post
func getPostDonationsPrefix(permlink types.Permlink) []byte {
	return append(append(postDonationsSubStore, permlink...), types.KeySeparator...)
}

func GetPostDonationKey(permlink types.Permlink, donateUser types.AccountKey) []byte {
	return append(getPostDonationsPrefix(permlink), donateUser...)
}
