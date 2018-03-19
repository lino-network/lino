package post

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
)

var (
	postKeyPrefix          = []byte("post/")
	postMetaKeyPrefix      = []byte("meta/")
	postLikesKeyPrefix     = []byte("likes/")
	postCommentsKeyPrefix  = []byte("comments/")
	postViewsKeyPrefix     = []byte("views/")
	postDonationsKeyPrefix = []byte("donations/")
)

// Implements types.PostManager
type postManager struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewPostManager returns a new types.PostManager that
// uses go-wire to (binary) encode and decode concrete types.Post
func NewPostMananger(key sdk.StoreKey) postManager {
	cdc := wire.NewCodec()

	return postManager{
		key: key,
		cdc: cdc,
	}
}

func (pm postManager) get(ctx sdk.Context, postKey types.PostKey, errFunc NotFoundErrFunc, prefix []byte) ([]byte, sdk.Error) {
	store := ctx.KVStore(pm.key)
	val := store.Get(append(prefix, postKey...))
	if val == nil {
		return nil, errFunc(postKey)
	}
	return val, nil
}

func (pm postManager) set(ctx sdk.Context, postKey types.PostKey, postStruct interface{}, prefix []byte) sdk.Error {
	store := ctx.KVStore(pm.key)
	val, err := json.Marshal(postStruct)
	if err != nil {
		return ErrPostMarshalError(err)
	}
	store.Set(append(prefix, postKey...), val)
	return nil
}

func (pm postManager) GetPost(ctx sdk.Context, postKey types.PostKey) (*types.Post, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostNotFound, postKeyPrefix)
	if err != nil {
		return nil, err
	}
	post := &types.Post{}
	unmarshalErr := json.Unmarshal(val, post)
	if unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return post, nil
}

func (pm postManager) SetPost(ctx sdk.Context, post *types.Post) sdk.Error {
	return pm.set(ctx, post.Key, post, postKeyPrefix)
}

func (pm postManager) GetPostMeta(ctx sdk.Context, postKey types.PostKey) (*types.PostMeta, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostMetaNotFound, postMetaKeyPrefix)
	if err != nil {
		return nil, err
	}
	postMeta := &types.PostMeta{}
	unmarshalErr := json.Unmarshal(val, postMeta)
	if unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postMeta, nil
}

func (pm postManager) SetPostMeta(ctx sdk.Context, postKey types.PostKey, postMeta *types.PostMeta) sdk.Error {
	return pm.set(ctx, postKey, postMeta, postMetaKeyPrefix)
}

func (pm postManager) GetPostLikes(ctx sdk.Context, postKey types.PostKey) (*types.PostLikes, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostLikesNotFound, postLikesKeyPrefix)
	if err != nil {
		return nil, err
	}
	postLikes := &types.PostLikes{}
	unmarshalErr := json.Unmarshal(val, postLikes)
	if unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postLikes, nil
}

func (pm postManager) SetPostLikes(ctx sdk.Context, postKey types.PostKey, postLikes *types.PostLikes) sdk.Error {
	return pm.set(ctx, postKey, postLikes, postLikesKeyPrefix)
}

func (pm postManager) GetPostComments(ctx sdk.Context, postKey types.PostKey) (*types.PostComments, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostCommentsNotFound, postCommentsKeyPrefix)
	if err != nil {
		return nil, err
	}
	postComments := &types.PostComments{}
	unmarshalErr := json.Unmarshal(val, postComments)
	if unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postComments, nil
}

func (pm postManager) SetPostComments(ctx sdk.Context, postKey types.PostKey, postComments *types.PostComments) sdk.Error {
	return pm.set(ctx, postKey, postComments, postCommentsKeyPrefix)
}

func (pm postManager) GetPostViews(ctx sdk.Context, postKey types.PostKey) (*types.PostViews, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostViewsNotFound, postViewsKeyPrefix)
	if err != nil {
		return nil, err
	}
	postViews := &types.PostViews{}
	unmarshalErr := json.Unmarshal(val, postViews)
	if unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postViews, nil
}

func (pm postManager) SetPostViews(ctx sdk.Context, postKey types.PostKey, postViews *types.PostViews) sdk.Error {
	return pm.set(ctx, postKey, postViews, postViewsKeyPrefix)
}

func (pm postManager) GetPostDonations(ctx sdk.Context, postKey types.PostKey) (*types.PostDonations, sdk.Error) {
	val, err := pm.get(ctx, postKey, ErrPostDonationsNotFound, postDonationsKeyPrefix)
	if err != nil {
		return nil, err
	}
	postDonations := &types.PostDonations{}
	unmarshalErr := json.Unmarshal(val, postDonations)
	if unmarshalErr != nil {
		return nil, ErrPostUnmarshalError(unmarshalErr)
	}
	return postDonations, nil
}

func (pm postManager) SetPostDonations(ctx sdk.Context, postKey types.PostKey, postDonations *types.PostDonations) sdk.Error {
	return pm.set(ctx, postKey, postDonations, postDonationsKeyPrefix)
}

// Check postManager implements PostManager interface
var _ types.PostManager = postManager{}
