package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	oldwire "github.com/tendermint/go-wire"
)

var (
	postInfoSubStore      = []byte{0x00} // SubStore for all post info
	postMetaSubStore      = []byte{0x01} // SubStore for all post mata info
	postLikeSubStore      = []byte{0x02} // SubStore for all like to post
	postCommentSubStore   = []byte{0x03} // SubStore for all comments
	postViewsSubStore     = []byte{0x04} // SubStore for all views
	postDonationsSubStore = []byte{0x05} // SubStore for all donations
)

// TODO(Lino) Register cdc here.
// temporary use old wire.
// this will help marshal and unmarshal interface type.
const msgTypePost = 0x1
const msgTypePostMeta = 0x2
const msgTypePostLike = 0x3
const msgTypePostReport = 0x4
const msgTypePostView = 0x5
const msgTypePostComment = 0x6
const msgTypePostDonations = 0x7

var _ = oldwire.RegisterInterface(
	struct{ PostInterface }{},
	oldwire.ConcreteType{PostInfo{}, msgTypePost},
	oldwire.ConcreteType{PostMeta{}, msgTypePostMeta},
	oldwire.ConcreteType{Like{}, msgTypePostLike},
	oldwire.ConcreteType{Report{}, msgTypePostReport},
	oldwire.ConcreteType{View{}, msgTypePostView},
	oldwire.ConcreteType{Comment{}, msgTypePostComment},
	oldwire.ConcreteType{Donation{}, msgTypePostDonations},
)

type PostManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewPostManager returns a new PostManager that
// uses go-wire to (binary) encode and decode concrete Post
func NewPostMananger(key sdk.StoreKey) PostManager {
	cdc := wire.NewCodec()

	return PostManager{
		key: key,
		cdc: cdc,
	}
}

func (pm PostManager) get(ctx sdk.Context, key []byte, errFunc NotFoundErrFunc) ([]byte, sdk.Error) {
	store := ctx.KVStore(pm.key)
	val := store.Get(key)
	if val == nil {
		return nil, errFunc(key)
	}
	return val, nil
}

func (pm PostManager) set(ctx sdk.Context, key []byte, postStruct PostInterface) sdk.Error {
	store := ctx.KVStore(pm.key)
	val, err := oldwire.MarshalJSON(postStruct)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(key, val)
	return nil
}

func (pm PostManager) GetPostInfo(ctx sdk.Context, postKey PostKey) (*PostInfo, sdk.Error) {
	val, err := pm.get(ctx, GetPostInfoKey(postKey), ErrPostNotFound)
	if err != nil {
		return nil, err
	}
	postInfo := new(PostInfo)
	if err := oldwire.UnmarshalJSON(val, postInfo); err != nil {
		return nil, ErrPostUnmarshalError(err)
	}
	return postInfo, nil
}

func (pm PostManager) SetPostInfo(ctx sdk.Context, postInfo *PostInfo) sdk.Error {
	return pm.set(ctx, GetPostInfoKey(GetPostKey(postInfo.Author, postInfo.PostID)), postInfo)
}

func (pm PostManager) GetPostMeta(ctx sdk.Context, postKey PostKey) (*PostMeta, sdk.Error) {
	val, err := pm.get(ctx, GetPostMetaKey(postKey), ErrPostMetaNotFound)
	if err != nil {
		return nil, err
	}
	postMeta := new(PostMeta)
	if unmarshalErr := oldwire.UnmarshalJSON(val, postMeta); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postMeta, nil
}

func (pm PostManager) SetPostMeta(ctx sdk.Context, postKey PostKey, postMeta *PostMeta) sdk.Error {
	return pm.set(ctx, GetPostMetaKey(postKey), postMeta)
}

func (pm PostManager) GetPostLike(ctx sdk.Context, postKey PostKey, likeUser acc.AccountKey) (*Like, sdk.Error) {
	val, err := pm.get(ctx, GetPostLikeKey(postKey, likeUser), ErrPostLikeNotFound)
	if err != nil {
		return nil, err
	}
	postLike := new(Like)
	if unmarshalErr := oldwire.UnmarshalJSON(val, postLike); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postLike, nil
}

func (pm PostManager) SetPostLike(ctx sdk.Context, postKey PostKey, postLike *Like) sdk.Error {
	return pm.set(ctx, GetPostLikeKey(postKey, postLike.Username), postLike)
}

func (pm PostManager) GetPostComment(ctx sdk.Context, postKey PostKey, commentPostKey PostKey) (*Comment, sdk.Error) {
	val, err := pm.get(ctx, GetPostCommentKey(postKey, commentPostKey), ErrPostCommentNotFound)
	if err != nil {
		return nil, err
	}
	postComment := new(Comment)
	if unmarshalErr := oldwire.UnmarshalJSON(val, postComment); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postComment, nil
}

func (pm PostManager) SetPostComment(ctx sdk.Context, postKey PostKey, postComment *Comment) sdk.Error {
	return pm.set(ctx, GetPostCommentKey(postKey, GetPostKey(postComment.Author, postComment.PostID)), postComment)
}

func (pm PostManager) GetPostView(ctx sdk.Context, postKey PostKey, viewUser acc.AccountKey) (*View, sdk.Error) {
	val, err := pm.get(ctx, GetPostViewKey(postKey, viewUser), ErrPostViewNotFound)
	if err != nil {
		return nil, err
	}
	postView := new(View)
	if unmarshalErr := oldwire.UnmarshalJSON(val, postView); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postView, nil
}

func (pm PostManager) SetPostView(ctx sdk.Context, postKey PostKey, postView *View) sdk.Error {
	return pm.set(ctx, GetPostViewKey(postKey, postView.Username), postView)
}

func (pm PostManager) GetPostDonation(ctx sdk.Context, postKey PostKey, donateUser acc.AccountKey) (*Donation, sdk.Error) {
	val, err := pm.get(ctx, GetPostDonationKey(postKey, donateUser), ErrPostDonationNotFound)
	if err != nil {
		return nil, err
	}
	postDonation := new(Donation)
	if unmarshalErr := oldwire.UnmarshalJSON(val, postDonation); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postDonation, nil
}

func (pm PostManager) SetPostDonation(ctx sdk.Context, postKey PostKey, postDonation *Donation) sdk.Error {
	return pm.set(ctx, GetPostDonationKey(postKey, postDonation.Username), postDonation)
}

func GetPostInfoKey(postKey PostKey) []byte {
	return append([]byte(postSubStore), postKey...)
}

func GetPostMetaKey(postKey PostKey) []byte {
	return append([]byte(postMetaSubStore), postKey...)
}

// PostLikePrefix format is LikeSubStore / PostKey
// which can be used to access all likes belong to this post
func GetPostLikePrefix(postKey PostKey) []byte {
	return append(append([]byte(postLikeSubStore), postKey...), types.KeySeparator...)
}

func GetPostLikeKey(postKey PostKey, likeUser acc.AccountKey) []byte {
	return append(GetPostLikePrefix(postKey), likeUser...)
}

// PostViewPrefix format is ViewSubStore / PostKey
// which can be used to access all views belong to this post
func GetPostViewPrefix(postKey PostKey) []byte {
	return append(append([]byte(postViewsSubStore), postKey...), types.KeySeparator...)
}

func GetPostViewKey(postKey PostKey, viewUser acc.AccountKey) []byte {
	return append(GetPostViewPrefix(postKey), viewUser...)
}

// PostCommentPrefix format is CommentSubStore / PostKey
// which can be used to access all comments belong to this post
func GetPostCommentPrefix(postKey PostKey) []byte {
	return append(append([]byte(postCommentSubStore), postKey...), types.KeySeparator...)
}

func GetPostCommentKey(postKey PostKey, commentPostKey PostKey) []byte {
	return append(GetPostCommentPrefix(postKey), commentPostKey...)
}

// PostDonationPrefix format is DonationSubStore / PostKey
// which can be used to access all donations belong to this post
func GetPostDonationPrefix(postKey PostKey) []byte {
	return append(append([]byte(postDonationsSubStore), postKey...), types.KeySeparator...)
}

func GetPostDonationKey(postKey PostKey, donateUser acc.AccountKey) []byte {
	return append(GetPostDonationPrefix(postKey), donateUser...)
}
