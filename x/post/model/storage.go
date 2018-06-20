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
	commentListSubStore        = []byte{0x07} // Substore for comment list
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

func (ps PostStorage) DoesPostExist(ctx sdk.Context, permLink types.PermLink) bool {
	store := ctx.KVStore(ps.key)
	return store.Has(GetPostInfoKey(permLink))
}

func (ps PostStorage) GetPostInfo(ctx sdk.Context, permLink types.PermLink) (*PostInfo, sdk.Error) {
	store := ctx.KVStore(ps.key)
	infoByte := store.Get(GetPostInfoKey(permLink))
	if infoByte == nil {
		return nil, ErrPostNotFound(GetPostInfoKey(permLink))
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
	store.Set(GetPostInfoKey(types.GetPermLink(postInfo.Author, postInfo.PostID)), infoByte)
	return nil
}

func (ps PostStorage) GetPostMeta(ctx sdk.Context, permLink types.PermLink) (*PostMeta, sdk.Error) {
	store := ctx.KVStore(ps.key)
	metaBytes := store.Get(GetPostMetaKey(permLink))
	if metaBytes == nil {
		return nil, ErrPostMetaNotFound(GetPostMetaKey(permLink))
	}
	postMeta := new(PostMeta)
	if unmarshalErr := ps.cdc.UnmarshalJSON(metaBytes, postMeta); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postMeta, nil
}

func (ps PostStorage) SetPostMeta(ctx sdk.Context, permLink types.PermLink, postMeta *PostMeta) sdk.Error {
	store := ctx.KVStore(ps.key)
	metaBytes, err := ps.cdc.MarshalJSON(*postMeta)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostMetaKey(permLink), metaBytes)
	return nil
}

func (ps PostStorage) GetPostLike(
	ctx sdk.Context, permLink types.PermLink, likeUser types.AccountKey) (*Like, sdk.Error) {
	store := ctx.KVStore(ps.key)
	likeBytes := store.Get(GetPostLikeKey(permLink, likeUser))
	if likeBytes == nil {
		return nil, ErrPostLikeNotFound(GetPostLikeKey(permLink, likeUser))
	}
	postLike := new(Like)
	if unmarshalErr := ps.cdc.UnmarshalJSON(likeBytes, postLike); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postLike, nil
}

func (ps PostStorage) SetPostLike(ctx sdk.Context, permLink types.PermLink, postLike *Like) sdk.Error {
	store := ctx.KVStore(ps.key)
	likeByte, err := ps.cdc.MarshalJSON(*postLike)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostLikeKey(permLink, postLike.Username), likeByte)
	return nil
}

func (ps PostStorage) GetPostReportOrUpvote(
	ctx sdk.Context, permLink types.PermLink, user types.AccountKey) (*ReportOrUpvote, sdk.Error) {
	store := ctx.KVStore(ps.key)
	reportOrUpvoteBytes := store.Get(GetPostReportOrUpvoteKey(permLink, user))
	if reportOrUpvoteBytes == nil {
		return nil, ErrPostReportOrUpvoteNotFound(GetPostReportOrUpvoteKey(permLink, user))
	}
	reportOrUpvote := new(ReportOrUpvote)
	if unmarshalErr := ps.cdc.UnmarshalJSON(reportOrUpvoteBytes, reportOrUpvote); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return reportOrUpvote, nil
}

func (ps PostStorage) SetPostReportOrUpvote(
	ctx sdk.Context, permLink types.PermLink, reportOrUpvote *ReportOrUpvote) sdk.Error {
	store := ctx.KVStore(ps.key)
	reportOrUpvoteByte, err := ps.cdc.MarshalJSON(*reportOrUpvote)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostReportOrUpvoteKey(permLink, reportOrUpvote.Username), reportOrUpvoteByte)
	return nil
}

func (ps PostStorage) RemovePostReportOrUpvote(
	ctx sdk.Context, permLink types.PermLink, user types.AccountKey) sdk.Error {
	store := ctx.KVStore(ps.key)
	store.Delete(GetPostReportOrUpvoteKey(permLink, user))
	return nil
}

func (ps PostStorage) GetPostComment(
	ctx sdk.Context, permLink types.PermLink, commentPermLink types.PermLink) (*Comment, sdk.Error) {
	store := ctx.KVStore(ps.key)
	commentBytes := store.Get(GetPostCommentKey(permLink, commentPermLink))
	if commentBytes == nil {
		return nil, ErrPostCommentNotFound(GetPostCommentKey(permLink, commentPermLink))
	}
	postComment := new(Comment)
	if unmarshalErr := ps.cdc.UnmarshalJSON(commentBytes, postComment); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postComment, nil
}

func (ps PostStorage) SetPostComment(
	ctx sdk.Context, permLink types.PermLink, postComment *Comment) sdk.Error {
	store := ctx.KVStore(ps.key)
	postCommentByte, err := ps.cdc.MarshalJSON(*postComment)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(
		GetPostCommentKey(permLink, types.GetPermLink(postComment.Author, postComment.PostID)),
		postCommentByte)
	return nil
}

func (ps PostStorage) GetPostView(
	ctx sdk.Context, permLink types.PermLink, viewUser types.AccountKey) (*View, sdk.Error) {
	store := ctx.KVStore(ps.key)
	viewBytes := store.Get(GetPostViewKey(permLink, viewUser))
	if viewBytes == nil {
		return nil, ErrPostViewNotFound(GetPostViewKey(permLink, viewUser))
	}
	postView := new(View)
	if unmarshalErr := ps.cdc.UnmarshalJSON(viewBytes, postView); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postView, nil
}

func (ps PostStorage) SetPostView(ctx sdk.Context, permLink types.PermLink, postView *View) sdk.Error {
	store := ctx.KVStore(ps.key)
	postViewByte, err := ps.cdc.MarshalJSON(*postView)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostViewKey(permLink, postView.Username), postViewByte)
	return nil
}

func (ps PostStorage) GetPostDonations(
	ctx sdk.Context, permLink types.PermLink, donateUser types.AccountKey) (*Donations, sdk.Error) {
	store := ctx.KVStore(ps.key)
	donateBytes := store.Get(GetPostDonationKey(permLink, donateUser))
	if donateBytes == nil {
		return nil, ErrPostDonationNotFound(GetPostDonationKey(permLink, donateUser))
	}
	postDonations := new(Donations)
	if unmarshalErr := ps.cdc.UnmarshalJSON(donateBytes, postDonations); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postDonations, nil
}

func (ps PostStorage) SetPostDonations(
	ctx sdk.Context, permLink types.PermLink, postDonations *Donations) sdk.Error {
	store := ctx.KVStore(ps.key)
	postDonationsByte, err := ps.cdc.MarshalJSON(*postDonations)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetPostDonationKey(permLink, postDonations.Username), postDonationsByte)
	return nil
}

func (ps PostStorage) GetCommentList(ctx sdk.Context, permLink types.PermLink) (*CommentList, sdk.Error) {
	store := ctx.KVStore(ps.key)
	lstByte := store.Get(GetCommentListKey(permLink))
	if lstByte == nil {
		return nil, nil
	}
	lst := new(CommentList)
	if unmarshalErr := ps.cdc.UnmarshalJSON(lstByte, lst); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return lst, nil
}

func (ps PostStorage) SetCommentList(
	ctx sdk.Context, permLink types.PermLink, lst *CommentList) sdk.Error {
	store := ctx.KVStore(ps.key)
	lstByte, err := ps.cdc.MarshalJSON(*lst)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(GetCommentListKey(permLink), lstByte)

	return nil
}

func GetPostInfoKey(permLink types.PermLink) []byte {
	return append(postInfoSubStore, permLink...)
}

func GetPostMetaKey(permLink types.PermLink) []byte {
	return append(postMetaSubStore, permLink...)
}

// PostLikePrefix format is LikeSubStore / PostKey
// which can be used to access all likes belong to this post
func getPostLikePrefix(permLink types.PermLink) []byte {
	return append(append(postLikeSubStore, permLink...), types.KeySeparator...)
}

func GetPostLikeKey(permLink types.PermLink, likeUser types.AccountKey) []byte {
	return append(getPostLikePrefix(permLink), likeUser...)
}

// PostReportPrefix format is ReportSubStore / PostKey
// which can be used to access all reports belong to this post
func getPostReportOrUpvotePrefix(permLink types.PermLink) []byte {
	return append(append(postReportOrUpvoteSubStore, permLink...), types.KeySeparator...)
}

func GetPostReportOrUpvoteKey(permLink types.PermLink, user types.AccountKey) []byte {
	return append(getPostReportOrUpvotePrefix(permLink), user...)
}

// PostViewPrefix format is ViewSubStore / PostKey
// which can be used to access all views belong to this post
func getPostViewPrefix(permLink types.PermLink) []byte {
	return append(append(postViewsSubStore, permLink...), types.KeySeparator...)
}

func GetPostViewKey(permLink types.PermLink, viewUser types.AccountKey) []byte {
	return append(getPostViewPrefix(permLink), viewUser...)
}

// PostCommentPrefix format is CommentSubStore / PostKey
// which can be used to access all comments belong to this post
func getPostCommentPrefix(permLink types.PermLink) []byte {
	return append(append(postCommentSubStore, permLink...), types.KeySeparator...)
}

func GetPostCommentKey(permLink types.PermLink, commentPostKey types.PermLink) []byte {
	return append(getPostCommentPrefix(permLink), commentPostKey...)
}

// PostDonationPrefix format is DonationSubStore / PostKey
// which can be used to access all donations belong to this post
func getPostDonationPrefix(permLink types.PermLink) []byte {
	return append(append(postDonationsSubStore, permLink...), types.KeySeparator...)
}

func GetPostDonationKey(permLink types.PermLink, donateUser types.AccountKey) []byte {
	return append(getPostDonationPrefix(permLink), donateUser...)
}

func GetCommentListKey(permLink types.PermLink) []byte {
	return append(commentListSubStore, permLink...)
}
