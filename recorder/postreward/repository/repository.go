package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/postreward"
)

type PostRewardRepository interface {
	Get(permlink string) (*postreward.PostReward, errors.Error)
	Add(info *postreward.PostReward) errors.Error
}
