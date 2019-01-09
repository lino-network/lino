package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/reward"
)

type RewardRepository interface {
	Add(detail *reward.Reward) errors.Error
	Get(username string) (*reward.Reward, errors.Error)
	IsEnable() bool
}
