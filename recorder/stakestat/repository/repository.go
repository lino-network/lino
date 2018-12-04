package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/stakestat"
)

type StakeStatRepository interface {
	Get(timestamp int64) (*stakestat.StakeStat, errors.Error)
	Add(info *stakestat.StakeStat) errors.Error
}
