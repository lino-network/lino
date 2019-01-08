package repository

import (
	"github.com/lino-network/lino/recorder/balancehistory"
	errors "github.com/lino-network/lino/recorder/errors"
)

type BalanceHistoryRepository interface {
	Add(detail *balancehistory.BalanceHistory) errors.Error
	Get(username string) (*balancehistory.BalanceHistory, errors.Error)
}
