package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/stake"
)

type StakeRepository interface {
	Get(username string) (*stake.Stake, errors.Error)
	Add(info *stake.Stake) errors.Error
}
