package types

import (
	"fmt"

	"github.com/tendermint/go-wire"
)

type Post struct {
	Title    string `json:"denom"`
	Author   string `json:"author"`
	Sequence int  `json:"seq"`
	Content  string `json:"content"`
}

func (post Post) String() string {
	return fmt.Sprintf("title:%v, author:%v, seq:%v, content:%v",
					   post.Title, post.Author, post.Sequence, post.Content)
}

// Post id is computed by the author and sequence.
// TODO: change to a better algorithm
func (post Post) PostID() []byte {
	return wire.BinaryBytes(post.Author+"#"+string(post.Sequence))
}

func GetPost(store KVStore, pid []byte) *Post {
	data := store.Get(pid)
	if len(data) == 0 {
		return nil
	}
	var post *Post
	err := wire.ReadBinaryBytes(data, &post)
	if err != nil {
		panic(fmt.Sprintf("Error reading Post %X error: %v",
			data, err.Error()))
	}
	return post
}

func SetPost(store KVStore, pid []byte, post *Post) {
	postBytes := wire.BinaryBytes(post)
	store.Set(pid, postBytes)
}
