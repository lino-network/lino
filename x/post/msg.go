package post

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = CreatePostMsg{}
var _ types.Msg = UpdatePostMsg{}
var _ types.Msg = DeletePostMsg{}
var _ types.Msg = LikeMsg{}
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
	Username       types.AccountKey `json:"username"`
	Amount         types.LNO        `json:"amount"`
	Author         types.AccountKey `json:"author"`
	PostID         string           `json:"post_id"`
	FromApp        types.AccountKey `json:"from_app"`
	Memo           string           `json:"memo"`
	IsMicroPayment bool             `json:"is_micropayment"`
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
func NewCreatePostMsg(
	author, postID, title, content, parentAuthor, parentPostID,
	sourceAuthor, sourcePostID, redistributionSplitRate string,
	links []types.IDToURLMapping) CreatePostMsg {
	return CreatePostMsg{
		Author:       types.AccountKey(author),
		PostID:       postID,
		Title:        title,
		Content:      content,
		SourceAuthor: types.AccountKey(sourceAuthor),
		SourcePostID: sourcePostID,
		Links:        links,
		RedistributionSplitRate: redistributionSplitRate,
	}
}

// NewUpdatePostMsg constructs a UpdatePost msg
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

// NewLikeMsg constructs a view msg
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
	postID string, fromApp string, memo string, isMicropayment bool) DonateMsg {
	return DonateMsg{
		Username:       types.AccountKey(user),
		Amount:         amount,
		Author:         types.AccountKey(author),
		PostID:         postID,
		FromApp:        types.AccountKey(fromApp),
		Memo:           memo,
		IsMicroPayment: isMicropayment,
	}
}

// NewReportOrUpvoteMsg constructs a ReportOrUpvote msg
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
func (msg UpdatePostMsg) Type() string     { return types.PostRouterName }
func (msg DeletePostMsg) Type() string     { return types.PostRouterName }
func (msg LikeMsg) Type() string           { return types.PostRouterName }
func (msg DonateMsg) Type() string         { return types.PostRouterName }
func (msg ReportOrUpvoteMsg) Type() string { return types.PostRouterName }
func (msg ViewMsg) Type() string           { return types.PostRouterName }

// ValidateBasic implements sdk.Msg
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
	if len(msg.Title) > types.MaxPostTitleLength {
		return ErrPostTitleExceedMaxLength()
	}
	if len(msg.Content) > types.MaxPostContentLength {
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

func (msg DeletePostMsg) ValidateBasic() sdk.Error {
	if len(msg.PostID) == 0 {
		return ErrNoPostID()
	}
	if len(msg.Author) == 0 {
		return ErrNoAuthor()
	}

	return nil
}

func (msg LikeMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.Username) == 0 {
		return ErrNoUsername()
	}
	if msg.Weight > types.MaxLikeWeight ||
		msg.Weight < types.MinLikeWeight {
		return ErrPostLikeWeightOverflow(msg.Weight)
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrInvalidTarget()
	}
	return nil
}

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

	if len(msg.Memo) > types.MaximumMemoLength {
		return ErrInvalidMemo()
	}
	return nil
}

// ValidateBasic implements sdk.Msg
func (msg ReportOrUpvoteMsg) ValidateBasic() sdk.Error {
	if len(msg.Username) == 0 {
		return ErrNoUsername()
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrInvalidTarget()
	}
	return nil
}

// ValidateBasic implements sdk.Msg
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
	return types.PostPermission
}
func (msg UpdatePostMsg) GetPermission() types.Permission {
	return types.PostPermission
}
func (msg DeletePostMsg) GetPermission() types.Permission {
	return types.PostPermission
}
func (msg LikeMsg) GetPermission() types.Permission {
	return types.PostPermission
}
func (msg DonateMsg) GetPermission() types.Permission {
	if msg.IsMicroPayment {
		return types.MicropaymentPermission
	}
	return types.TransactionPermission
}
func (msg ReportOrUpvoteMsg) GetPermission() types.Permission {
	return types.PostPermission
}
func (msg ViewMsg) GetPermission() types.Permission {
	return types.PostPermission
}

// GetSignBytes implements sdk.Msg
func (msg CreatePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

func (msg UpdatePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

func (msg DeletePostMsg) GetSignBytes() []byte {
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
	b, err := msgCdc.MarshalJSON(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners implements sdk.Msg.
func (msg CreatePostMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Author)}
}
func (msg UpdatePostMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Author)}
}
func (msg DeletePostMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Author)}
}
func (msg LikeMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}
func (msg DonateMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}
func (msg ReportOrUpvoteMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}
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
