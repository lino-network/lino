package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/stakeStat"
)

type StakeStatRepository interface {
	Get(timestamp int64) (*stakeStat.StakeStat, errors.Error)
	Add(info *stakeStat.StakeStat) errors.Error
}
