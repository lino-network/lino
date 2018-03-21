package post

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

// CreateMsg contains information to create a post
type CreateMsg struct {
	CreateInfo
}

type CreateInfo struct {
	PostID       string         `json:"post_id"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	Author       acc.AccountKey `json:"author"`
	ParentAuthor acc.AccountKey `json:"parent_author"`
	ParentPostID string         `json:"parent_post_id"`
	SourceAuthor acc.AccountKey `json:"source_author"`
	SourcePostID string         `json:"source_post_id"`
	Links        IDToURLMapping `json:"links"`
}

// NewCreateMsg constructs a post msg
func NewCreateMsg(createInfo CreateInfo) CreateMsg {
	return CreateMsg{CreateInfo: createInfo}
}

// Type implements sdk.Msg
func (msg CreateMsg) Type() string { return "post" } // TODO change to "post/create", wait for base app udpate

// ValidateBasic implements sdk.Msg
func (msg CreateMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.PostID) == 0 {
		return ErrPostCreateNoPostID()
	}
	if len(msg.Author) == 0 {
		return ErrPostCreateNoAuthor()
	}
	if len(msg.Title) > types.MaxPostTitleLength {
		return ErrPostTitleExceedMaxLength()
	}
	if len(msg.Content) > types.MaxPostContentLength {
		return ErrPostContentExceedMaxLength()
	}
	return nil
}

// Get implements sdk.Msg; should not be called
func (msg CreateMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg CreateMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners implements Msg.
func (msg CreateMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Author)}
}

func (msg CreateMsg) String() string {
	return fmt.Sprintf("Post.CreateMsg{Info:%v}", msg.CreateInfo)
}
