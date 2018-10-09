package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/topcontent"
)

type TopContentRepository interface {
	Get(timestamp int64) (*topcontent.TopContent, errors.Error)
	Add(info *topcontent.TopContent) errors.Error
}
