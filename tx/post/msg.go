package post

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

// CreatePostMsg contains information to create a post
type CreatePostMsg struct {
	types.Post
}

// NewCreatePostMsg constructs a post msg
func NewCreatePostMsg(post types.Post) CreatePostMsg {
	return CreatePostMsg{Post: post}
}

// Type implements sdk.Msg
func (msg CreatePostMsg) Type() string { return "post" } // TODO change to "post/create", wait for base app udpate

// ValidateBasic implements sdk.Msg
func (msg CreatePostMsg) ValidateBasic() sdk.Error {
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
func (msg CreatePostMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg CreatePostMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners implements Msg.
func (msg CreatePostMsg) GetSigners() []sdk.Address {
	return []sdk.Address{sdk.Address(msg.Author)}
}

func (msg CreatePostMsg) String() string {
	return fmt.Sprintf("Post.CreatePostMsg{post:%v}", msg.Post)
}
