package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"
)

var (
	postKeyPrefix          = []byte("post/")
	postMetaKeyPrefix      = []byte("meta/")
	postLikesKeyPrefix     = []byte("likes/")
	postCommentsKeyPrefix  = []byte("comments/")
	postViewsKeyPrefix     = []byte("views/")
	postDonationsKeyPrefix = []byte("donations/")
)

// TODO(Lino) Register cdc here.
// temporary use old wire.
// this will help marshal and unmarshal interface type.
const msgTypePost = 0x1
const msgTypePostMeta = 0x2
const msgTypePostLike = 0x3
const msgTypePostComments = 0x4
const msgTypePostDonations = 0x5

var _ = oldwire.RegisterInterface(
	struct{ PostInterface }{},
	oldwire.ConcreteType{PostInfo{}, msgTypePost},
	oldwire.ConcreteType{PostMeta{}, msgTypePostMeta},
	oldwire.ConcreteType{PostLikes{}, msgTypePostLike},
	oldwire.ConcreteType{PostComments{}, msgTypePostComments},
	oldwire.ConcreteType{PostDonations{}, msgTypePostDonations},
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
	val, err := pm.get(ctx, PostInfoKey(postKey), ErrPostNotFound)
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
	return pm.set(ctx, PostInfoKey(GetPostKey(postInfo.Author, postInfo.PostID)), postInfo)
}

func (pm PostManager) GetPostMeta(ctx sdk.Context, postKey PostKey) (*PostMeta, sdk.Error) {
	val, err := pm.get(ctx, PostMetaKey(postKey), ErrPostMetaNotFound)
	if err != nil {
		return nil, err
	}
	postMeta := &PostMeta{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postMeta); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postMeta, nil
}

func (pm PostManager) SetPostMeta(ctx sdk.Context, postKey PostKey, postMeta *PostMeta) sdk.Error {
	return pm.set(ctx, PostMetaKey(postKey), postMeta)
}

func (pm PostManager) GetPostLikes(ctx sdk.Context, postKey PostKey) (*PostLikes, sdk.Error) {
	val, err := pm.get(ctx, PostLikesKey(postKey), ErrPostLikesNotFound)
	if err != nil {
		return nil, err
	}
	postLikes := &PostLikes{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postLikes); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postLikes, nil
}

func (pm PostManager) SetPostLikes(ctx sdk.Context, postKey PostKey, postLikes *PostLikes) sdk.Error {
	return pm.set(ctx, PostLikesKey(postKey), postLikes)
}

func (pm PostManager) GetPostComments(ctx sdk.Context, postKey PostKey) (*PostComments, sdk.Error) {
	val, err := pm.get(ctx, PostCommentsKey(postKey), ErrPostCommentsNotFound)
	if err != nil {
		return nil, err
	}
	postComments := &PostComments{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postComments); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postComments, nil
}

func (pm PostManager) SetPostComments(ctx sdk.Context, postKey PostKey, postComments *PostComments) sdk.Error {
	return pm.set(ctx, PostCommentsKey(postKey), postComments)
}

func (pm PostManager) GetPostViews(ctx sdk.Context, postKey PostKey) (*PostViews, sdk.Error) {
	val, err := pm.get(ctx, PostViewsKey(postKey), ErrPostViewsNotFound)
	if err != nil {
		return nil, err
	}
	postViews := &PostViews{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postViews); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postViews, nil
}

func (pm PostManager) SetPostViews(ctx sdk.Context, postKey PostKey, postViews *PostViews) sdk.Error {
	return pm.set(ctx, PostViewsKey(postKey), postViews)
}

func (pm PostManager) GetPostDonations(ctx sdk.Context, postKey PostKey) (*PostDonations, sdk.Error) {
	val, err := pm.get(ctx, PostDonationKey(postKey), ErrPostDonationsNotFound)
	if err != nil {
		return nil, err
	}
	postDonations := &PostDonations{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postDonations); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postDonations, nil
}

func (pm PostManager) SetPostDonations(ctx sdk.Context, postKey PostKey, postDonations *PostDonations) sdk.Error {
	return pm.set(ctx, PostDonationKey(postKey), postDonations)
}

func PostInfoKey(postKey PostKey) []byte {
	return append(postKeyPrefix, postKey...)
}

func PostMetaKey(postKey PostKey) []byte {
	return append(postMetaKeyPrefix, postKey...)
}

func PostLikesKey(postKey PostKey) []byte {
	return append(postLikesKeyPrefix, postKey...)
}

func PostViewsKey(postKey PostKey) []byte {
	return append(postViewsKeyPrefix, postKey...)
}

func PostCommentsKey(postKey PostKey) []byte {
	return append(postCommentsKeyPrefix, postKey...)
}

func PostDonationKey(postKey PostKey) []byte {
	return append(postDonationsKeyPrefix, postKey...)
}
