package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/post"
)

type PostRepository interface {
	Get(author string) (*post.Post, errors.Error)
	Add(info *post.Post) errors.Error
	SetReward(author, postID string, amount string) errors.Error
	IsEnable() bool
}
