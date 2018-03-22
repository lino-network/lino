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

func (pm PostManager) get(ctx sdk.Context, postKey PostKey, errFunc NotFoundErrFunc, prefix []byte) ([]byte, sdk.Error) {
	store := ctx.KVStore(pm.key)
	val := store.Get(append(prefix, postKey...))
	if val == nil {
		return nil, errFunc(postKey)
	}
	return val, nil
}

func (pm PostManager) set(ctx sdk.Context, postKey PostKey, postStruct PostInterface, prefix []byte) sdk.Error {
	store := ctx.KVStore(pm.key)
	val, err := oldwire.MarshalJSON(postStruct)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(append(prefix, postKey...), val)
	return nil
}

func (pm PostManager) GetPostInfo(ctx sdk.Context, postKey PostKey) (*PostInfo, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostNotFound, postKeyPrefix)
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
	return pm.set(ctx, GetPostKey(postInfo.Author, postInfo.PostID), postInfo, postKeyPrefix)
}

func (pm PostManager) GetPostMeta(ctx sdk.Context, postKey PostKey) (*PostMeta, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostMetaNotFound, postMetaKeyPrefix)
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
	return pm.set(ctx, postKey, postMeta, postMetaKeyPrefix)
}

func (pm PostManager) GetPostLikes(ctx sdk.Context, postKey PostKey) (*PostLikes, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostLikesNotFound, postLikesKeyPrefix)
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
	return pm.set(ctx, postKey, postLikes, postLikesKeyPrefix)
}

func (pm PostManager) GetPostComments(ctx sdk.Context, postKey PostKey) (*PostComments, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostCommentsNotFound, postCommentsKeyPrefix)
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
	return pm.set(ctx, postKey, postComments, postCommentsKeyPrefix)
}

func (pm PostManager) GetPostViews(ctx sdk.Context, postKey PostKey) (*PostViews, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostViewsNotFound, postViewsKeyPrefix)
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
	return pm.set(ctx, postKey, postViews, postViewsKeyPrefix)
}

func (pm PostManager) GetPostDonations(ctx sdk.Context, postKey PostKey) (*PostDonations, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostDonationsNotFound, postDonationsKeyPrefix)
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
	return pm.set(ctx, postKey, postDonations, postDonationsKeyPrefix)
}
