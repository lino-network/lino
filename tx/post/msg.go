package post

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

var _ sdk.Msg = CreatePostMsg{}
var _ sdk.Msg = LikeMsg{}
var _ sdk.Msg = DonateMsg{}
var _ sdk.Msg = ReportOrUpvoteMsg{}

// PostCreateParams can also use to publish comment(with parent) or repost(with source)
type PostCreateParams struct {
	PostID                  string                 `json:"post_id"`
	Title                   string                 `json:"title"`
	Content                 string                 `json:"content"`
	Author                  types.AccountKey       `json:"author"`
	ParentAuthor            types.AccountKey       `json:"parent_author"`
	ParentPostID            string                 `json:"parent_postID"`
	SourceAuthor            types.AccountKey       `json:"source_author"`
	SourcePostID            string                 `json:"source_postID"`
	Links                   []types.IDToURLMapping `json:"links"`
	RedistributionSplitRate sdk.Rat                `json:"redistribution_split_rate"`
}

// CreatePostMsg contains information to create a post
type CreatePostMsg struct {
	PostCreateParams
}

// LikeMsg sent from a user to a post
type LikeMsg struct {
	Username types.AccountKey
	Weight   int64
	Author   types.AccountKey
	PostID   string
}

// DonateMsg sent from a user to a post
type DonateMsg struct {
	Username types.AccountKey
	Amount   types.LNO
	Author   types.AccountKey
	PostID   string
}

// ReportOrUpvoteMsg sent from a user to a post
type ReportOrUpvoteMsg struct {
	Username types.AccountKey
	Author   types.AccountKey
	PostID   string
	IsReport bool
	IsRevoke bool
}

// NewCreatePostMsg constructs a post msg
func NewCreatePostMsg(postCreateParams PostCreateParams) CreatePostMsg {
	return CreatePostMsg{PostCreateParams: postCreateParams}
}

// NewLikeMsg constructs a like msg
func NewLikeMsg(user types.AccountKey, weight int64, author types.AccountKey, postID string) LikeMsg {
	return LikeMsg{
		Username: user,
		Weight:   weight,
		Author:   author,
		PostID:   postID,
	}
}

// NewDonateMsg constructs a donate msg
func NewDonateMsg(user types.AccountKey, amount types.LNO, author types.AccountKey, postID string) DonateMsg {
	return DonateMsg{
		Username: user,
		Amount:   amount,
		Author:   author,
		PostID:   postID,
	}
}

// NewReportOrUpvoteMsg constructs a report msg
func NewReportOrUpvoteMsg(user types.AccountKey, author types.AccountKey, postID string, isReport bool, isRevoke bool) ReportOrUpvoteMsg {
	return ReportOrUpvoteMsg{
		Username: user,
		Author:   author,
		PostID:   postID,
		IsReport: isReport,
		IsRevoke: isRevoke,
	}
}

// Type implements sdk.Msg
func (msg CreatePostMsg) Type() string     { return types.PostRouterName } // TODO change to "post/create", wait for base app udpate
func (msg LikeMsg) Type() string           { return types.PostRouterName } // TODO change to "post/create", wait for base app udpate
func (msg DonateMsg) Type() string         { return types.PostRouterName } // TODO change to "post/create", wait for base app udpate
func (msg ReportOrUpvoteMsg) Type() string { return types.PostRouterName } // TODO change to "post/create", wait for base app udpate

// ValidateBasic implements sdk.Msg
func (msg CreatePostMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.PostID) == 0 {
		return ErrPostCreateNoPostID()
	}
	if len(msg.Author) == 0 {
		return ErrPostCreateNoAuthor()
	}
	if (len(msg.ParentAuthor) > 0 || len(msg.ParentPostID) > 0) &&
		(len(msg.SourceAuthor) > 0 || len(msg.SourcePostID) > 0) {
		return ErrCommentAndRepostError()
	}
	if len(msg.Title) > types.MaxPostTitleLength {
		return ErrPostTitleExceedMaxLength()
	}
	if len(msg.Content) > types.MaxPostContentLength {
		return ErrPostContentExceedMaxLength()
	}
	return nil
}

func (msg LikeMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.Username) == 0 {
		return ErrPostLikeNoUsername()
	}
	if msg.Weight > types.MaxLikeWeight ||
		msg.Weight < types.MinLikeWeight {
		return ErrPostLikeWeightOverflow(msg.Weight)
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrPostLikeInvalidTarget()
	}
	return nil
}

func (msg DonateMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.Username) == 0 {
		return ErrPostDonateNoUsername()
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrPostDonateInvalidTarget()
	}

	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

// ValidateBasic implements sdk.Msg
func (msg ReportOrUpvoteMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) == 0 {
		return ErrPostReportOrUpvoteNoUsername()
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrPostReportOrUpvoteInvalidTarget()
	}
	return nil
}

// Get implements sdk.Msg; should not be called
func (msg CreatePostMsg) Get(key interface{}) (value interface{}) {
	return nil
}
func (msg LikeMsg) Get(key interface{}) (value interface{}) {
	return nil
}
func (msg DonateMsg) Get(key interface{}) (value interface{}) {
	return nil
}
func (msg ReportOrUpvoteMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg CreatePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}
func (msg LikeMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

func (msg DonateMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

func (msg ReportOrUpvoteMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

func getSignBytes(msg sdk.Msg) []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners implements sdk.Msg.
func (msg CreatePostMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Author)}
}
func (msg LikeMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}
func (msg DonateMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}
func (msg ReportOrUpvoteMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

// String implements Stringer
func (msg CreatePostMsg) String() string {
	return fmt.Sprintf("Post.CreatePostMsg{postInfo:%v}", msg.PostCreateParams)
}
func (msg LikeMsg) String() string {
	return fmt.Sprintf("Post.LikeMsg{like from: %v, weight: %v, post auther:%v, post id: %v}", msg.Username, msg.Weight, msg.Author, msg.PostID)
}
func (msg DonateMsg) String() string {
	return fmt.Sprintf("Post.DonateMsg{donation from: %v, amount: %v, post auther:%v, post id: %v}", msg.Username, msg.Amount, msg.Author, msg.PostID)
}
func (msg ReportOrUpvoteMsg) String() string {
	return fmt.Sprintf("Post.ReportOrUpvoteMsg{from: %v, post auther:%v, post id: %v}", msg.Username, msg.Author, msg.PostID)
}
