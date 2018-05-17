package post

import (
	"encoding/json"
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = CreatePostMsg{}
var _ sdk.Msg = LikeMsg{}
var _ sdk.Msg = DonateMsg{}
var _ sdk.Msg = ReportOrUpvoteMsg{}
var _ sdk.Msg = ViewMsg{}
var _ sdk.Msg = UpdatePostMsg{}

// PostCreateParams can also use to publish comment(with parent) or repost(with source)
type PostCreateParams struct {
	Author                  types.AccountKey       `json:"author"`
	PostID                  string                 `json:"post_id"`
	Title                   string                 `json:"title"`
	Content                 string                 `json:"content"`
	ParentAuthor            types.AccountKey       `json:"parent_author"`
	ParentPostID            string                 `json:"parent_postID"`
	SourceAuthor            types.AccountKey       `json:"source_author"`
	SourcePostID            string                 `json:"source_postID"`
	Links                   []types.IDToURLMapping `json:"links"`
	RedistributionSplitRate string                 `json:"redistribution_split_rate"`
}

// CreatePostMsg contains information to create a post
type CreatePostMsg struct {
	PostCreateParams
}

type UpdatePostMsg struct {
	Author                  types.AccountKey       `json:"author"`
	PostID                  string                 `json:"post_id"`
	Title                   string                 `json:"title"`
	Content                 string                 `json:"content"`
	Links                   []types.IDToURLMapping `json:"links"`
	RedistributionSplitRate string                 `json:"redistribution_split_rate"`
}

type DeletePostMsg struct {
	Author types.AccountKey `json:"author"`
	PostID string           `json:"post_id"`
	Titile string           `json:"title"`
}

// LikeMsg sent from a user to a post
type LikeMsg struct {
	Username types.AccountKey `json:"username"`
	Weight   int64            `json:"weight"`
	Author   types.AccountKey `json:"author"`
	PostID   string           `json:"post_id"`
}

// DonateMsg sent from a user to a post
type DonateMsg struct {
	Username     types.AccountKey `json:"username"`
	Amount       types.LNO        `json:"amount"`
	Author       types.AccountKey `json:"author"`
	PostID       string           `json:"post_id"`
	FromApp      types.AccountKey `json:"from_app"`
	FromChecking bool             `json:"from_checking"`
	Memo         string           `json:"memo"`
}

// ViewMsg sent from a user to a post
type ViewMsg struct {
	Username types.AccountKey `json:"username"`
	Author   types.AccountKey `json:"author"`
	PostID   string           `json:"post_id"`
}

// ReportOrUpvoteMsg sent from a user to a post
type ReportOrUpvoteMsg struct {
	Username types.AccountKey `json:"username"`
	Author   types.AccountKey `json:"author"`
	PostID   string           `json:"post_id"`
	IsReport bool             `json:"is_report"`
}

// NewCreatePostMsg constructs a post msg
func NewCreatePostMsg(postCreateParams PostCreateParams) CreatePostMsg {
	return CreatePostMsg{PostCreateParams: postCreateParams}
}

// NewUpdatePostMsg constructs a like msg
func NewUpdatePostMsg(
	author, postID, title, content string,
	links []types.IDToURLMapping, redistributionSplitRate string) UpdatePostMsg {
	return UpdatePostMsg{
		Author:  types.AccountKey(author),
		PostID:  postID,
		Title:   title,
		Content: content,
		Links:   links,
		RedistributionSplitRate: redistributionSplitRate,
	}
}

// NewLikeMsg constructs a like msg
func NewLikeMsg(
	user string, weight int64, author, postID string) LikeMsg {

	return LikeMsg{
		Username: types.AccountKey(user),
		Weight:   weight,
		Author:   types.AccountKey(author),
		PostID:   postID,
	}
}

// NewLikeMsg constructs a like msg
func NewViewMsg(user, author string, postID string) ViewMsg {

	return ViewMsg{
		Username: types.AccountKey(user),
		Author:   types.AccountKey(author),
		PostID:   postID,
	}
}

// NewDonateMsg constructs a donate msg
func NewDonateMsg(
	user string, amount types.LNO, author string,
	postID string, fromApp string, fromChecking bool, memo string) DonateMsg {

	return DonateMsg{
		Username:     types.AccountKey(user),
		Amount:       amount,
		Author:       types.AccountKey(author),
		PostID:       postID,
		FromApp:      types.AccountKey(fromApp),
		FromChecking: fromChecking,
		Memo:         memo,
	}
}

// NewReportOrUpvoteMsg constructs a report msg
func NewReportOrUpvoteMsg(
	user, author, postID string, isReport bool) ReportOrUpvoteMsg {

	return ReportOrUpvoteMsg{
		Username: types.AccountKey(user),
		Author:   types.AccountKey(author),
		PostID:   postID,
		IsReport: isReport,
	}
}

// Type implements sdk.Msg
func (msg CreatePostMsg) Type() string     { return types.PostRouterName }
func (msg LikeMsg) Type() string           { return types.PostRouterName }
func (msg DonateMsg) Type() string         { return types.PostRouterName }
func (msg ReportOrUpvoteMsg) Type() string { return types.PostRouterName }
func (msg ViewMsg) Type() string           { return types.PostRouterName }
func (msg UpdatePostMsg) Type() string     { return types.PostRouterName }

// ValidateBasic implements sdk.Msg
func (msg CreatePostMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.PostID) == 0 {
		return ErrNoPostID()
	}
	if len(msg.Author) == 0 {
		return ErrNoAuthor()
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

	splitRate, err := sdk.NewRatFromDecimal(msg.RedistributionSplitRate)
	if err != nil {
		return ErrPostRedistributionSplitRate()
	}

	if splitRate.LT(sdk.ZeroRat) || splitRate.GT(sdk.OneRat) {
		return ErrPostRedistributionSplitRate()
	}
	return nil
}

// ValidateBasic implements sdk.Msg
func (msg UpdatePostMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.PostID) == 0 {
		return ErrNoPostID()
	}
	if len(msg.Author) == 0 {
		return ErrNoAuthor()
	}
	if len(msg.Title) > types.MaxPostTitleLength {
		return ErrPostTitleExceedMaxLength()
	}
	if len(msg.Content) > types.MaxPostContentLength {
		return ErrPostContentExceedMaxLength()
	}

	splitRate, err := sdk.NewRatFromDecimal(msg.RedistributionSplitRate)
	if err != nil {
		return ErrPostRedistributionSplitRate()
	}

	if splitRate.LT(sdk.ZeroRat) || splitRate.GT(sdk.OneRat) {
		return ErrPostRedistributionSplitRate()
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

	if len(msg.Memo) > types.MaximumMemoLength {
		return ErrInvalidMemo()
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

// ValidateBasic implements sdk.Msg
func (msg ViewMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) == 0 {
		return ErrPostViewNoUsername()
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrPostViewInvalidTarget()
	}
	return nil
}

// Get implements sdk.Msg; should not be called
func (msg CreatePostMsg) Get(key interface{}) (value interface{}) {
	return nil
}
func (msg UpdatePostMsg) Get(key interface{}) (value interface{}) {
	return nil
}
func (msg LikeMsg) Get(key interface{}) (value interface{}) {
	return nil
}
func (msg DonateMsg) Get(key interface{}) (value interface{}) {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	if keyStr == types.PermissionLevel {
		if msg.FromChecking {
			return types.PostPermission
		}
		return types.TransactionPermission
	}
	return nil
}
func (msg ReportOrUpvoteMsg) Get(key interface{}) (value interface{}) {
	return nil
}
func (msg ViewMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg CreatePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

func (msg UpdatePostMsg) GetSignBytes() []byte {
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

func (msg ViewMsg) GetSignBytes() []byte {
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
func (msg UpdatePostMsg) GetSigners() []sdk.Address {
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
func (msg ViewMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Username)}
}

// String implements Stringer
func (msg CreatePostMsg) String() string {
	return fmt.Sprintf("Post.CreatePostMsg{postInfo:%v}", msg.PostCreateParams)
}

func (msg UpdatePostMsg) String() string {
	return fmt.Sprintf("Post.UpdatePostMsg{author:%v, postID:%v, title:%v, content:%v, links:%v, redistribution split rate:%v}",
		msg.Author, msg.PostID, msg.Title, msg.Content,
		msg.Links, msg.RedistributionSplitRate)
}

func (msg LikeMsg) String() string {
	return fmt.Sprintf(
		"Post.LikeMsg{like from: %v, weight: %v, post auther:%v, post id: %v}",
		msg.Username, msg.Weight, msg.Author, msg.PostID)
}
func (msg DonateMsg) String() string {
	return fmt.Sprintf(
		"Post.DonateMsg{donation from: %v, amount: %v, post auther:%v, post id: %v}",
		msg.Username, msg.Amount, msg.Author, msg.PostID)
}
func (msg ReportOrUpvoteMsg) String() string {
	return fmt.Sprintf(
		"Post.ReportOrUpvoteMsg{from: %v, post auther:%v, post id: %v}",
		msg.Username, msg.Author, msg.PostID)
}
func (msg ViewMsg) String() string {
	return fmt.Sprintf(
		"Post.ViewMsg{from: %v, post auther:%v, post id: %v}",
		msg.Username, msg.Author, msg.PostID)
}
