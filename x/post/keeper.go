package post

//go:generate mockery -name PostKeeper

import (
	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/manager"
	"github.com/lino-network/lino/x/post/model"
	"github.com/lino-network/lino/x/post/types"
)

type PostKeeper interface {
	DoesPostExist(ctx sdk.Context, permlink linotypes.Permlink) bool
	GetPost(ctx sdk.Context, permlink linotypes.Permlink) (model.Post, sdk.Error)
	CreatePost(ctx sdk.Context, author linotypes.AccountKey, postID string, createdBy linotypes.AccountKey, content string, title string) sdk.Error
	UpdatePost(ctx sdk.Context, author linotypes.AccountKey, postID, title, content string) sdk.Error
	DeletePost(ctx sdk.Context, permlink linotypes.Permlink) sdk.Error
	LinoDonate(ctx sdk.Context, from linotypes.AccountKey, amount linotypes.Coin, author linotypes.AccountKey, postID string, app linotypes.AccountKey) sdk.Error
	IDADonate(ctx sdk.Context, from linotypes.AccountKey, n linotypes.MiniIDA, author linotypes.AccountKey, postID string, app, signer linotypes.AccountKey) sdk.Error
	ExecRewardEvent(ctx sdk.Context, reward types.RewardEvent) sdk.Error

	// querier
	GetComsumptionWindow(ctx sdk.Context) linotypes.MiniDollar

	ImportFromFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
	ExportToFile(ctx sdk.Context, cdc *codec.Codec, filepath string) error
}

var _ PostKeeper = manager.PostManager{}
