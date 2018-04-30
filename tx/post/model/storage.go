package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
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

// TODO(Lino) Register cdc here.
// temporary use old wire.
// this will help marshal and unmarshal interface type.
const (
	msgTypePost               = "1"
	msgTypePostMeta           = "2"
	msgTypePostLike           = "3"
	msgTypePostReportOrUpvote = "4"
	msgTypePostView           = "5"
	msgTypePostComment        = "6"
	msgTypePostDonations      = "7"
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

	cdc.RegisterInterface((*PostInterface)(nil), nil)
	cdc.RegisterConcrete(PostInfo{}, msgTypePost, nil)
	cdc.RegisterConcrete(PostMeta{}, msgTypePostMeta, nil)
	cdc.RegisterConcrete(Like{}, msgTypePostLike, nil)
	cdc.RegisterConcrete(ReportOrUpvote{}, msgTypePostReportOrUpvote, nil)
	cdc.RegisterConcrete(View{}, msgTypePostView, nil)
	cdc.RegisterConcrete(Comment{}, msgTypePostComment, nil)
	cdc.RegisterConcrete(Donation{}, msgTypePostDonations, nil)
	wire.RegisterCrypto(cdc)

	return PostStorage{
		key: key,
		cdc: cdc,
	}
}

func (ps PostStorage) get(ctx sdk.Context, key []byte, errFunc NotFoundErrFunc) ([]byte, sdk.Error) {
	store := ctx.KVStore(ps.key)
	val := store.Get(key)
	if val == nil {
		return nil, errFunc(key)
	}
	return val, nil
}

func (ps PostStorage) set(ctx sdk.Context, key []byte, postStruct PostInterface) sdk.Error {
	store := ctx.KVStore(ps.key)
	val, err := ps.cdc.MarshalJSON(postStruct)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(key, val)
	return nil
}

func (ps PostStorage) GetPostInfo(ctx sdk.Context, postKey types.PostKey) (*PostInfo, sdk.Error) {
	val, err := ps.get(ctx, GetPostInfoKey(postKey), ErrPostNotFound)
	if err != nil {
		return nil, err
	}
	postInfo := new(PostInfo)
	if err := ps.cdc.UnmarshalJSON(val, postInfo); err != nil {
		return nil, ErrPostUnmarshalError(err)
	}
	return postInfo, nil
}

func (ps PostStorage) SetPostInfo(ctx sdk.Context, postInfo *PostInfo) sdk.Error {
	return ps.set(ctx, GetPostInfoKey(types.GetPostKey(postInfo.Author, postInfo.PostID)), postInfo)
}

func (ps PostStorage) GetPostMeta(ctx sdk.Context, postKey types.PostKey) (*PostMeta, sdk.Error) {
	val, err := ps.get(ctx, GetPostMetaKey(postKey), ErrPostMetaNotFound)
	if err != nil {
		return nil, err
	}
	postMeta := new(PostMeta)
	if unmarshalErr := ps.cdc.UnmarshalJSON(val, postMeta); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postMeta, nil
}

func (ps PostStorage) SetPostMeta(ctx sdk.Context, postKey types.PostKey, postMeta *PostMeta) sdk.Error {
	return ps.set(ctx, GetPostMetaKey(postKey), postMeta)
}

func (ps PostStorage) GetPostLike(ctx sdk.Context, postKey types.PostKey, likeUser types.AccountKey) (*Like, sdk.Error) {
	val, err := ps.get(ctx, GetPostLikeKey(postKey, likeUser), ErrPostLikeNotFound)
	if err != nil {
		return nil, err
	}
	postLike := new(Like)
	if unmarshalErr := ps.cdc.UnmarshalJSON(val, postLike); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postLike, nil
}

func (ps PostStorage) SetPostLike(ctx sdk.Context, postKey types.PostKey, postLike *Like) sdk.Error {
	return ps.set(ctx, GetPostLikeKey(postKey, postLike.Username), postLike)
}

func (ps PostStorage) GetPostReportOrUpvote(ctx sdk.Context, postKey types.PostKey, user types.AccountKey) (*ReportOrUpvote, sdk.Error) {
	val, err := ps.get(ctx, GetPostReportOrUpvoteKey(postKey, user), ErrPostReportOrUpvoteNotFound)
	if err != nil {
		return nil, err
	}
	reportOrUpvote := new(ReportOrUpvote)
	if unmarshalErr := ps.cdc.UnmarshalJSON(val, reportOrUpvote); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return reportOrUpvote, nil
}

func (ps PostStorage) SetPostReportOrUpvote(ctx sdk.Context, postKey types.PostKey, reportOrUpvote *ReportOrUpvote) sdk.Error {
	return ps.set(ctx, GetPostReportOrUpvoteKey(postKey, reportOrUpvote.Username), reportOrUpvote)
}

func (ps PostStorage) RemovePostReportOrUpvote(ctx sdk.Context, postKey types.PostKey, user types.AccountKey) sdk.Error {
	store := ctx.KVStore(ps.key)
	store.Delete(GetPostReportOrUpvoteKey(postKey, user))
	return nil
}

func (ps PostStorage) GetPostComment(ctx sdk.Context, postKey types.PostKey, commentPostKey types.PostKey) (*Comment, sdk.Error) {
	val, err := ps.get(ctx, GetPostCommentKey(postKey, commentPostKey), ErrPostCommentNotFound)
	if err != nil {
		return nil, err
	}
	postComment := new(Comment)
	if unmarshalErr := ps.cdc.UnmarshalJSON(val, postComment); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postComment, nil
}

func (ps PostStorage) SetPostComment(ctx sdk.Context, postKey types.PostKey, postComment *Comment) sdk.Error {
	return ps.set(ctx, GetPostCommentKey(postKey, types.GetPostKey(postComment.Author, postComment.PostID)), postComment)
}

func (ps PostStorage) GetPostView(ctx sdk.Context, postKey types.PostKey, viewUser types.AccountKey) (*View, sdk.Error) {
	val, err := ps.get(ctx, GetPostViewKey(postKey, viewUser), ErrPostViewNotFound)
	if err != nil {
		return nil, err
	}
	postView := new(View)
	if unmarshalErr := ps.cdc.UnmarshalJSON(val, postView); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postView, nil
}

func (ps PostStorage) SetPostView(ctx sdk.Context, postKey types.PostKey, postView *View) sdk.Error {
	return ps.set(ctx, GetPostViewKey(postKey, postView.Username), postView)
}

func (ps PostStorage) GetPostDonations(ctx sdk.Context, postKey types.PostKey, donateUser types.AccountKey) (*Donations, sdk.Error) {
	val, err := ps.get(ctx, GetPostDonationKey(postKey, donateUser), ErrPostDonationNotFound)
	if err != nil {
		return nil, err
	}
	postDonations := new(Donations)
	if unmarshalErr := ps.cdc.UnmarshalJSON(val, postDonations); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postDonations, nil
}

func (ps PostStorage) SetPostDonations(ctx sdk.Context, postKey types.PostKey, postDonations *Donations) sdk.Error {
	return ps.set(ctx, GetPostDonationKey(postKey, postDonations.Username), postDonations)
}

func GetPostInfoKey(postKey types.PostKey) []byte {
	return append([]byte(postInfoSubStore), postKey...)
}

func GetPostMetaKey(postKey types.PostKey) []byte {
	return append([]byte(postMetaSubStore), postKey...)
}

// PostLikePrefix format is LikeSubStore / PostKey
// which can be used to access all likes belong to this post
func GetPostLikePrefix(postKey types.PostKey) []byte {
	return append(append([]byte(postLikeSubStore), postKey...), types.KeySeparator...)
}

func GetPostLikeKey(postKey types.PostKey, likeUser types.AccountKey) []byte {
	return append(GetPostLikePrefix(postKey), likeUser...)
}

// PostReportPrefix format is ReportSubStore / PostKey
// which can be used to access all reports belong to this post
func GetPostReportOrUpvotePrefix(postKey types.PostKey) []byte {
	return append(append([]byte(postReportOrUpvoteSubStore), postKey...), types.KeySeparator...)
}

func GetPostReportOrUpvoteKey(postKey types.PostKey, user types.AccountKey) []byte {
	return append(GetPostReportOrUpvotePrefix(postKey), user...)
}

// PostViewPrefix format is ViewSubStore / PostKey
// which can be used to access all views belong to this post
func GetPostViewPrefix(postKey types.PostKey) []byte {
	return append(append([]byte(postViewsSubStore), postKey...), types.KeySeparator...)
}

func GetPostViewKey(postKey types.PostKey, viewUser types.AccountKey) []byte {
	return append(GetPostViewPrefix(postKey), viewUser...)
}

// PostCommentPrefix format is CommentSubStore / PostKey
// which can be used to access all comments belong to this post
func GetPostCommentPrefix(postKey types.PostKey) []byte {
	return append(append([]byte(postCommentSubStore), postKey...), types.KeySeparator...)
}

func GetPostCommentKey(postKey types.PostKey, commentPostKey types.PostKey) []byte {
	return append(GetPostCommentPrefix(postKey), commentPostKey...)
}

// PostDonationPrefix format is DonationSubStore / PostKey
// which can be used to access all donations belong to this post
func GetPostDonationPrefix(postKey types.PostKey) []byte {
	return append(append([]byte(postDonationsSubStore), postKey...), types.KeySeparator...)
}

func GetPostDonationKey(postKey types.PostKey, donateUser types.AccountKey) []byte {
	return append(GetPostDonationPrefix(postKey), donateUser...)
}
