package post

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// CreateMsg contains information to create a post
type CreateMsg struct {
	Info CreateInfo `json:"createInfo"`
}

// CreateInfo is used to construct CreateMsg
type CreateInfo struct {
	Permlink       string            `json:"permlink"`
	Title          string            `json:"title"`
	Content        string            `json:"content"`
	Author         types.AccountKey  `json:"author"`
	ParentAuthor   types.AccountKey  `json:"parentAuthor"`
	ParentPermlink string            `json:"parentPermlink"`
	SourceAuthor   types.AccountKey  `json:"sourceAuthor"`
	SourcePermlink string            `json:"sourcePermlink"`
	Links          types.LinkMapping `json:"links"`
}

// NewCreateMsg constructs a post msg
func NewCreateMsg(createInfo CreateInfo) CreateMsg {
	return CreateMsg{Info: createInfo}
}

// Type implements sdk.Msg
func (msg CreateMsg) Type() string { return "post/create" }

// ValidateBasic implements sdk.Msg
func (msg CreateMsg) ValidateBasic() sdk.Error {
	// Ensure permlink exists
	if len(msg.Info.Permlink) == 0 {
		return ErrPostCreateNoPermlink()
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
	return []sdk.Address{sdk.Address(msg.Info.Author)}
}

func (msg CreateMsg) String() string {
	return fmt.Sprintf("Post.CreateMsg{Info:%v}", msg.Info)
}
