package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/inflation"
)

type InflationRepository interface {
	Get(timestamp int64) (*inflation.Inflation, errors.Error)
	Add(info *inflation.Inflation) errors.Error
}
