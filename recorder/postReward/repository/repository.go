package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/postReward"
)

type PostRewardRepository interface {
	Get(timestamp int64) (*postReward.PostReward, errors.Error)
	Add(info *postReward.PostReward) errors.Error
}
