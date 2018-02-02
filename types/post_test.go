package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-wire"
)

func TestPostID(t *testing.T) {
	post := Post{
		Author: "test",
		Sequence: 1,
	}
	assert.Equal(
		t, post.PostID(),
		    wire.BinaryBytes(
		        post.Author + "#" + string(post.Sequence)))

}
