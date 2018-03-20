package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
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
// this will help use to marshal and unmarshal interface type.
const msgTypePost = 0x1
const msgTypePostMeta = 0x2
const msgTypePostLike = 0x3
const msgTypePostComments = 0x4
const msgTypePostDonations = 0x5

var _ = oldwire.RegisterInterface(
	struct{ types.PostInterface }{},
	oldwire.ConcreteType{types.Post{}, msgTypePost},
	oldwire.ConcreteType{types.PostMeta{}, msgTypePostMeta},
	oldwire.ConcreteType{types.PostLikes{}, msgTypePostLike},
	oldwire.ConcreteType{types.PostComments{}, msgTypePostComments},
	oldwire.ConcreteType{types.PostDonations{}, msgTypePostDonations},
)

// Implements types.PostManager
type PostManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewPostManager returns a new types.PostManager that
// uses go-wire to (binary) encode and decode concrete types.Post
func NewPostMananger(key sdk.StoreKey) PostManager {
	cdc := wire.NewCodec()

	return PostManager{
		key: key,
		cdc: cdc,
	}
}

func (pm PostManager) get(ctx sdk.Context, postKey types.PostKey, errFunc NotFoundErrFunc, prefix []byte) ([]byte, sdk.Error) {
	store := ctx.KVStore(pm.key)
	val := store.Get(append(prefix, postKey...))
	if val == nil {
		return nil, errFunc(postKey)
	}
	return val, nil
}

func (pm PostManager) set(ctx sdk.Context, postKey types.PostKey, postStruct types.PostInterface, prefix []byte) sdk.Error {
	store := ctx.KVStore(pm.key)
	val, err := oldwire.MarshalJSON(postStruct)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(append(prefix, postKey...), val)
	return nil
}

func (pm PostManager) GetPost(ctx sdk.Context, postKey types.PostKey) (*types.Post, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostNotFound, postKeyPrefix)
	if err != nil {
		return nil, err
	}
	post := new(types.Post)
	if err := oldwire.UnmarshalJSON(val, post); err != nil {
		return nil, ErrPostUnmarshalError(err)
	}
	return post, nil
}

func (pm PostManager) SetPost(ctx sdk.Context, post *types.Post) sdk.Error {
	return pm.set(ctx, types.GetPostKey(post.Author, post.PostID), post, postKeyPrefix)
}

func (pm PostManager) GetPostMeta(ctx sdk.Context, postKey types.PostKey) (*types.PostMeta, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostMetaNotFound, postMetaKeyPrefix)
	if err != nil {
		return nil, err
	}
	postMeta := &types.PostMeta{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postMeta); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postMeta, nil
}

func (pm PostManager) SetPostMeta(ctx sdk.Context, postKey types.PostKey, postMeta *types.PostMeta) sdk.Error {
	return pm.set(ctx, postKey, postMeta, postMetaKeyPrefix)
}

func (pm PostManager) GetPostLikes(ctx sdk.Context, postKey types.PostKey) (*types.PostLikes, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostLikesNotFound, postLikesKeyPrefix)
	if err != nil {
		return nil, err
	}
	postLikes := &types.PostLikes{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postLikes); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postLikes, nil
}

func (pm PostManager) SetPostLikes(ctx sdk.Context, postKey types.PostKey, postLikes *types.PostLikes) sdk.Error {
	return pm.set(ctx, postKey, postLikes, postLikesKeyPrefix)
}

func (pm PostManager) GetPostComments(ctx sdk.Context, postKey types.PostKey) (*types.PostComments, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostCommentsNotFound, postCommentsKeyPrefix)
	if err != nil {
		return nil, err
	}
	postComments := &types.PostComments{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postComments); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postComments, nil
}

func (pm PostManager) SetPostComments(ctx sdk.Context, postKey types.PostKey, postComments *types.PostComments) sdk.Error {
	return pm.set(ctx, postKey, postComments, postCommentsKeyPrefix)
}

func (pm PostManager) GetPostViews(ctx sdk.Context, postKey types.PostKey) (*types.PostViews, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostViewsNotFound, postViewsKeyPrefix)
	if err != nil {
		return nil, err
	}
	postViews := &types.PostViews{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postViews); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postViews, nil
}

func (pm PostManager) SetPostViews(ctx sdk.Context, postKey types.PostKey, postViews *types.PostViews) sdk.Error {
	return pm.set(ctx, postKey, postViews, postViewsKeyPrefix)
}

func (pm PostManager) GetPostDonations(ctx sdk.Context, postKey types.PostKey) (*types.PostDonations, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostDonationsNotFound, postDonationsKeyPrefix)
	if err != nil {
		return nil, err
	}
	postDonations := &types.PostDonations{}
	if unmarshalErr := oldwire.UnmarshalJSON(val, postDonations); unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postDonations, nil
}

func (pm PostManager) SetPostDonations(ctx sdk.Context, postKey types.PostKey, postDonations *types.PostDonations) sdk.Error {
	return pm.set(ctx, postKey, postDonations, postDonationsKeyPrefix)
}

// Check PostManager implements PostManager interface
var _ types.PostManager = PostManager{}
