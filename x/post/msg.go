package post

import (
	"fmt"
	"unicode/utf8"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = CreatePostMsg{}
var _ types.Msg = UpdatePostMsg{}
var _ types.Msg = DeletePostMsg{}
var _ types.Msg = DonateMsg{}
var _ types.Msg = ReportOrUpvoteMsg{}
var _ types.Msg = ViewMsg{}

// CreatePostMsg contains information to create a post
type CreatePostMsg struct {
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

// UpdatePostMsg - update post
type UpdatePostMsg struct {
	Author                  types.AccountKey       `json:"author"`
	PostID                  string                 `json:"post_id"`
	Title                   string                 `json:"title"`
	Content                 string                 `json:"content"`
	Links                   []types.IDToURLMapping `json:"links"`
	RedistributionSplitRate string                 `json:"redistribution_split_rate"`
}

// DeletePostMsg - sent from a user to a post
type DeletePostMsg struct {
	Author types.AccountKey `json:"author"`
	PostID string           `json:"post_id"`
}

// DonateMsg - sent from a user to a post
type DonateMsg struct {
	Username types.AccountKey `json:"username"`
	Amount   types.LNO        `json:"amount"`
	Author   types.AccountKey `json:"author"`
	PostID   string           `json:"post_id"`
	FromApp  types.AccountKey `json:"from_app"`
	Memo     string           `json:"memo"`
}

// ViewMsg - sent from a user to a post
type ViewMsg struct {
	Username types.AccountKey `json:"username"`
	Author   types.AccountKey `json:"author"`
	PostID   string           `json:"post_id"`
}

// ReportOrUpvoteMsg - sent from a user to a post
type ReportOrUpvoteMsg struct {
	Username types.AccountKey `json:"username"`
	Author   types.AccountKey `json:"author"`
	PostID   string           `json:"post_id"`
	IsReport bool             `json:"is_report"`
}

// NewCreatePostMsg - constructs a post msg
func NewCreatePostMsg(
	author, postID, title, content, parentAuthor, parentPostID,
	sourceAuthor, sourcePostID, redistributionSplitRate string,
	links []types.IDToURLMapping) CreatePostMsg {
	return CreatePostMsg{
		Author:       types.AccountKey(author),
		PostID:       postID,
		Title:        title,
		Content:      content,
		ParentAuthor: types.AccountKey(parentAuthor),
		ParentPostID: parentPostID,
		SourceAuthor: types.AccountKey(sourceAuthor),
		SourcePostID: sourcePostID,
		Links:        links,
		RedistributionSplitRate: redistributionSplitRate,
	}
}

// NewUpdatePostMsg - constructs a UpdatePost msg
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

func NewDeletePostMsg(author, postID string) DeletePostMsg {
	return DeletePostMsg{
		Author: types.AccountKey(author),
		PostID: postID,
	}
}

// NewViewMsg - constructs a view msg
func NewViewMsg(user, author string, postID string) ViewMsg {
	return ViewMsg{
		Username: types.AccountKey(user),
		Author:   types.AccountKey(author),
		PostID:   postID,
	}
}

// NewDonateMsg - constructs a donate msg
func NewDonateMsg(
	user string, amount types.LNO, author string,
	postID string, fromApp string, memo string) DonateMsg {
	return DonateMsg{
		Username: types.AccountKey(user),
		Amount:   amount,
		Author:   types.AccountKey(author),
		PostID:   postID,
		FromApp:  types.AccountKey(fromApp),
		Memo:     memo,
	}
}

// NewReportOrUpvoteMsg - constructs a ReportOrUpvote msg
func NewReportOrUpvoteMsg(
	user, author, postID string, isReport bool) ReportOrUpvoteMsg {

	return ReportOrUpvoteMsg{
		Username: types.AccountKey(user),
		Author:   types.AccountKey(author),
		PostID:   postID,
		IsReport: isReport,
	}
}

// Type - implements sdk.Msg
func (msg CreatePostMsg) Type() string { return types.PostRouterName }

// Type - implements sdk.Msg
func (msg UpdatePostMsg) Type() string { return types.PostRouterName }

// Type - implements sdk.Msg
func (msg DeletePostMsg) Type() string { return types.PostRouterName }

// Type - implements sdk.Msg
func (msg DonateMsg) Type() string { return types.PostRouterName }

// Type - implements sdk.Msg
func (msg ReportOrUpvoteMsg) Type() string { return types.PostRouterName }

// Type - implements sdk.Msg
func (msg ViewMsg) Type() string { return types.PostRouterName }

// ValidateBasic - implements sdk.Msg
func (msg CreatePostMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.PostID) == 0 {
		return ErrNoPostID()
	}
	if len(msg.PostID) > types.MaximumLengthOfPostID {
		return ErrPostIDTooLong()
	}
	if len(msg.Author) == 0 {
		return ErrNoAuthor()
	}
	if (len(msg.ParentAuthor) > 0 || len(msg.ParentPostID) > 0) &&
		(len(msg.SourceAuthor) > 0 || len(msg.SourcePostID) > 0) {
		return ErrCommentAndRepostConflict()
	}
	if utf8.RuneCountInString(msg.Title) > types.MaxPostTitleLength {
		return ErrPostTitleExceedMaxLength()
	}
	if utf8.RuneCountInString(msg.Content) > types.MaxPostContentLength {
		return ErrPostContentExceedMaxLength()
	}
	if len(msg.RedistributionSplitRate) > types.MaximumSdkRatLength {
		return ErrRedistributionSplitRateLengthTooLong()
	}

	if len(msg.Links) > types.MaximumNumOfLinks {
		return ErrTooManyURL()
	}

	for _, link := range msg.Links {
		if len(link.Identifier) > types.MaximumLinkIdentifier {
			return ErrIdentifierLengthTooLong()
		}
		if len(link.URL) > types.MaximumLinkURL {
			return ErrURLLengthTooLong()
		}
	}

	splitRate, err := sdk.NewRatFromDecimal(msg.RedistributionSplitRate, types.NewRatFromDecimalPrecision)
	if err != nil {
		return err
	}
	if splitRate.LT(sdk.ZeroRat()) || splitRate.GT(sdk.OneRat()) {
		return ErrInvalidPostRedistributionSplitRate()
	}
	return nil
}

// ValidateBasic - implements sdk.Msg
func (msg UpdatePostMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.PostID) == 0 {
		return ErrNoPostID()
	}
	if len(msg.Author) == 0 {
		return ErrNoAuthor()
	}
	if utf8.RuneCountInString(msg.Title) > types.MaxPostTitleLength {
		return ErrPostTitleExceedMaxLength()
	}
	if utf8.RuneCountInString(msg.Content) > types.MaxPostContentLength {
		return ErrPostContentExceedMaxLength()
	}
	if len(msg.RedistributionSplitRate) > types.MaximumSdkRatLength {
		return ErrRedistributionSplitRateLengthTooLong()
	}

	for _, link := range msg.Links {
		if len(link.Identifier) > types.MaximumLinkIdentifier {
			return ErrIdentifierLengthTooLong()
		}
		if len(link.URL) > types.MaximumLinkURL {
			return ErrURLLengthTooLong()
		}
	}

	splitRate, err := sdk.NewRatFromDecimal(msg.RedistributionSplitRate, types.NewRatFromDecimalPrecision)
	if err != nil {
		return err
	}

	if splitRate.LT(sdk.ZeroRat()) || splitRate.GT(sdk.OneRat()) {
		return ErrInvalidPostRedistributionSplitRate()
	}
	return nil
}

// ValidateBasic - implements sdk.Msg
func (msg DeletePostMsg) ValidateBasic() sdk.Error {
	if len(msg.PostID) == 0 {
		return ErrNoPostID()
	}
	if len(msg.Author) == 0 {
		return ErrNoAuthor()
	}

	return nil
}

// ValidateBasic - implements sdk.Msg
func (msg DonateMsg) ValidateBasic() sdk.Error {
	// Ensure permlink  exists
	if len(msg.Username) == 0 {
		return ErrNoUsername()
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrInvalidTarget()
	}

	_, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err
	}

	if utf8.RuneCountInString(msg.Memo) > types.MaximumMemoLength {
		return ErrInvalidMemo()
	}
	return nil
}

// ValidateBasic - implements sdk.Msg
func (msg ReportOrUpvoteMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) == 0 {
		return ErrNoUsername()
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrInvalidTarget()
	}
	return nil
}

// ValidateBasic - implements sdk.Msg
func (msg ViewMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) == 0 {
		return ErrNoUsername()
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrInvalidTarget()
	}
	return nil
}

// Get implements sdk.Msg; should not be called
func (msg CreatePostMsg) GetPermission() types.Permission {
	return types.AppPermission
}
func (msg UpdatePostMsg) GetPermission() types.Permission {
	return types.AppPermission
}
func (msg DeletePostMsg) GetPermission() types.Permission {
	return types.AppPermission
}
func (msg DonateMsg) GetPermission() types.Permission {
	return types.PreAuthorizationPermission
}
func (msg ReportOrUpvoteMsg) GetPermission() types.Permission {
	return types.AppPermission
}
func (msg ViewMsg) GetPermission() types.Permission {
	return types.AppPermission
}

// GetSignBytes - implements sdk.Msg
func (msg CreatePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSignBytes - implements sdk.Msg
func (msg UpdatePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSignBytes - implements sdk.Msg
func (msg DeletePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSignBytes - implements sdk.Msg
func (msg DonateMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSignBytes - implements sdk.Msg
func (msg ReportOrUpvoteMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSignBytes - implements sdk.Msg
func (msg ViewMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

func getSignBytes(msg sdk.Msg) []byte {
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners - implements sdk.Msg.
func (msg CreatePostMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Author)}
}

// GetSigners - implements sdk.Msg.
func (msg UpdatePostMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Author)}
}

// GetSigners - implements sdk.Msg.
func (msg DeletePostMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Author)}
}

// GetSigners - implements sdk.Msg.
func (msg DonateMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetSigners - implements sdk.Msg.
func (msg ReportOrUpvoteMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// GetSigners - implements sdk.Msg.
func (msg ViewMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

// String implements Stringer
func (msg CreatePostMsg) String() string {
	return fmt.Sprintf("Post.CreatePostMsg{author:%v, postID:%v, title:%v, content:%v, parentAuthor:%v,"+
		"parentPostID:%v, sourceAuthor:%v, sourcePostID:%v,links:%v, redistribution split rate:%v}",
		msg.Author, msg.PostID, msg.Title, msg.Content, msg.ParentAuthor, msg.ParentPostID, msg.SourceAuthor, msg.SourcePostID,
		msg.Links, msg.RedistributionSplitRate)
}

func (msg UpdatePostMsg) String() string {
	return fmt.Sprintf("Post.UpdatePostMsg{author:%v, postID:%v, title:%v, content:%v, links:%v, redistribution split rate:%v}",
		msg.Author, msg.PostID, msg.Title, msg.Content,
		msg.Links, msg.RedistributionSplitRate)
}

func (msg DeletePostMsg) String() string {
	return fmt.Sprintf("Post.DeletePostMsg{author:%v, postID:%v}", msg.Author, msg.PostID)
}

func (msg DonateMsg) String() string {
	return fmt.Sprintf(
		"Post.DonateMsg{donation from: %v, amount: %v, post author:%v, post id: %v}",
		msg.Username, msg.Amount, msg.Author, msg.PostID)
}

func (msg ReportOrUpvoteMsg) String() string {
	return fmt.Sprintf(
		"Post.ReportOrUpvoteMsg{from: %v, post author:%v, post id: %v}",
		msg.Username, msg.Author, msg.PostID)
}

func (msg ViewMsg) String() string {
	return fmt.Sprintf(
		"Post.ViewMsg{from: %v, post author:%v, post id: %v}",
		msg.Username, msg.Author, msg.PostID)
}

// GetConsumeAmount - implements types.Msg.
func (msg CreatePostMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// GetConsumeAmount - implements types.Msg.
func (msg UpdatePostMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// GetConsumeAmount - implements types.Msg.
func (msg DeletePostMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// GetConsumeAmount - implements types.Msg.
func (msg DonateMsg) GetConsumeAmount() types.Coin {
	coin, _ := types.LinoToCoin(msg.Amount)
	return coin
}

// GetConsumeAmount - implements types.Msg.
func (msg ReportOrUpvoteMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// GetConsumeAmount - implements types.Msg.
func (msg ViewMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}
