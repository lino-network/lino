package types

import (
	"fmt"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lino-network/lino/types"
)

// CreatePostMsg contains information to create a post
// required stateful validation:
// createdBy is a developer, if not author.
type CreatePostMsg struct {
	Author    types.AccountKey `json:"author"`
	PostID    string           `json:"post_id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	CreatedBy types.AccountKey `json:"created_by"`
}

var _ types.Msg = CreatePostMsg{}

// NewCreatePostMsg - constructs a post msg
func NewCreatePostMsg(author, postID, title, content, createdBy string) CreatePostMsg {
	return CreatePostMsg{
		Author:    types.AccountKey(author),
		PostID:    postID,
		Title:     title,
		Content:   content,
		CreatedBy: types.AccountKey(createdBy),
	}
}

// Route - implements sdk.Msg
func (msg CreatePostMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg CreatePostMsg) Type() string { return "CreatePostMsg" }

// GetSigners - implements sdk.Msg
func (msg CreatePostMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.CreatedBy)}
}

// GetSignBytes - implements sdk.Msg
func (msg CreatePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetPermission - implements types.Msg
func (msg CreatePostMsg) GetPermission() types.Permission {
	if msg.CreatedBy == msg.Author {
		return types.TransactionPermission
	}
	return types.AppOrAffiliatedPermission
}

// GetConsumeAmount - implements types.Msg
func (msg CreatePostMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// ValidateBasic - implements sdk.Msg
func (msg CreatePostMsg) ValidateBasic() sdk.Error {
	err := checkPostBasic(msg.PostID, msg.Author, msg.Title, msg.Content)
	if err != nil {
		return err
	}
	if len(msg.CreatedBy) == 0 {
		return ErrNoCreatedBy()
	}
	return nil
}

// String implements Stringer
func (msg CreatePostMsg) String() string {
	return fmt.Sprintf(
		"Post.CreatePostMsg{author:%v, postID:%v, title:%v, content:%v, created_by:%v}",
		msg.Author, msg.PostID, msg.Title, msg.Content, msg.CreatedBy)
}

// UpdatePostMsg - update post
type UpdatePostMsg struct {
	Author  types.AccountKey `json:"author"`
	PostID  string           `json:"post_id"`
	Title   string           `json:"title"`
	Content string           `json:"content"`
}

var _ types.Msg = UpdatePostMsg{}

// NewUpdatePostMsg - constructs a UpdatePost msg
func NewUpdatePostMsg(author, postID, title, content string) UpdatePostMsg {
	return UpdatePostMsg{
		Author:  types.AccountKey(author),
		PostID:  postID,
		Title:   title,
		Content: content,
	}
}

// Route - implements sdk.Msg
func (msg UpdatePostMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg UpdatePostMsg) Type() string { return "UpdatePostMsg" }

// GetPermission - implements types.Msg
func (msg UpdatePostMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// ValidateBasic - implements sdk.Msg
func (msg UpdatePostMsg) ValidateBasic() sdk.Error {
	err := checkPostBasic(msg.PostID, msg.Author, msg.Title, msg.Content)
	if err != nil {
		return err
	}
	return nil
}

func (msg UpdatePostMsg) String() string {
	return fmt.Sprintf("Post.UpdatePostMsg{author:%v, postID:%v, title:%v, content:%v}",
		msg.Author, msg.PostID, msg.Title, msg.Content)
}

// GetSignBytes - implements sdk.Msg
func (msg UpdatePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg UpdatePostMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Author)}
}

// GetConsumeAmount - implements types.Msg
func (msg UpdatePostMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// DeletePostMsg - sent from a user to a post
type DeletePostMsg struct {
	Author types.AccountKey `json:"author"`
	PostID string           `json:"post_id"`
}

var _ types.Msg = DeletePostMsg{}

func NewDeletePostMsg(author, postID string) DeletePostMsg {
	return DeletePostMsg{
		Author: types.AccountKey(author),
		PostID: postID,
	}
}

// Route - implements sdk.Msg
func (msg DeletePostMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg DeletePostMsg) Type() string { return "DeletePostMsg" }

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

// GetPermission - implements types.Msg
func (msg DeletePostMsg) GetPermission() types.Permission {
	return types.TransactionPermission
}

// GetSignBytes - implements sdk.Msg
func (msg DeletePostMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg DeletePostMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Author)}
}

// GetConsumeAmount - implements types.Msg
func (msg DeletePostMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

func (msg DeletePostMsg) String() string {
	return fmt.Sprintf("Post.DeletePostMsg{author:%v, postID:%v}", msg.Author, msg.PostID)
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

var _ types.Msg = DonateMsg{}

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

// Route - implements sdk.Msg
func (msg DonateMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg DonateMsg) Type() string { return "DonateMsg" }

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

// GetPermission - implements types.Msg
func (msg DonateMsg) GetPermission() types.Permission {
	return types.PreAuthorizationPermission
}

// GetSignBytes - implements sdk.Msg
func (msg DonateMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg DonateMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Username)}
}

func (msg DonateMsg) String() string {
	return fmt.Sprintf(
		"Post.DonateMsg{donation from: %v, amount: %v, post author:%v, post id: %v}",
		msg.Username, msg.Amount, msg.Author, msg.PostID)
}

// GetConsumeAmount - implements types.Msg
func (msg DonateMsg) GetConsumeAmount() types.Coin {
	coin, _ := types.LinoToCoin(msg.Amount)
	return coin
}

// DonateMsg - sent from a user to a post
type IDADonateMsg struct {
	Username types.AccountKey `json:"username"`
	App      types.AccountKey `json:"app"`
	Amount   types.IDAStr     `json:"amount"`
	Author   types.AccountKey `json:"author"`
	PostID   string           `json:"post_id"`
	Memo     string           `json:"memo"`
}

var _ types.Msg = DonateMsg{}

// NewIDADonateMsg - constructs a donate msg use In-app digital assets.
func NewIDADonateMsg(user, app string, amount types.IDAStr, author string, postID string, memo string) IDADonateMsg {
	return IDADonateMsg{
		Username: types.AccountKey(user),
		App:      types.AccountKey(app),
		Amount:   amount,
		Author:   types.AccountKey(author),
		PostID:   postID,
		Memo:     memo,
	}
}

// Route - implements sdk.Msg
func (msg IDADonateMsg) Route() string { return RouterKey }

// Type - implements sdk.Msg
func (msg IDADonateMsg) Type() string { return "IDADonateMsg" }

// ValidateBasic - implements sdk.Msg
func (msg IDADonateMsg) ValidateBasic() sdk.Error {
	// Ensure permlink  exists
	if len(msg.Username) == 0 {
		return ErrNoUsername()
	}
	if len(msg.App) == 0 {
		return ErrNoApp()
	}
	if len(msg.Author) == 0 || len(msg.PostID) == 0 {
		return ErrInvalidTarget()
	}

	_, err := msg.Amount.ToIDA()
	if err != nil {
		return err
	}

	if utf8.RuneCountInString(msg.Memo) > types.MaximumMemoLength {
		return ErrInvalidMemo()
	}
	return nil
}

// GetPermission - implements types.Msg
func (msg IDADonateMsg) GetPermission() types.Permission {
	return types.AppOrAffiliatedPermission
}

// GetSignBytes - implements sdk.Msg
func (msg IDADonateMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// GetSigners - implements sdk.Msg
func (msg IDADonateMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.App)}
}

func (msg IDADonateMsg) String() string {
	return fmt.Sprintf(
		"Post.IDADonateMsg{donation from:%v, app:%v, amount: %v, author:%v, pid:%v, memo:%v}",
		msg.Username, msg.App, msg.Amount, msg.Author, msg.PostID, msg.Memo)
}

// GetConsumeAmount - implements types.Msg
// TODO(yumin): outdated.
func (msg IDADonateMsg) GetConsumeAmount() types.Coin {
	return types.NewCoinFromInt64(0)
}

// utils
func getSignBytes(msg sdk.Msg) []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func checkPostBasic(postID string, author types.AccountKey, title, content string) sdk.Error {
	if len(postID) == 0 {
		return ErrNoPostID()
	}
	if len(postID) > types.MaximumLengthOfPostID {
		return ErrPostIDTooLong()
	}
	if len(author) == 0 {
		return ErrNoAuthor()
	}
	if utf8.RuneCountInString(title) > types.MaxPostTitleLength {
		return ErrPostTitleExceedMaxLength()
	}
	if utf8.RuneCountInString(content) > types.MaxPostContentLength {
		return ErrPostContentExceedMaxLength()
	}
	return nil
}
