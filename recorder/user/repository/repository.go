package repository

import (
	errors "github.com/lino-network/lino/recorder/errors"
	"github.com/lino-network/lino/recorder/user"
)

type UserRepository interface {
	Add(user *user.User) errors.Error
	Get(username string) (*user.User, errors.Error)
	IncreaseSequenceNumber(username string) errors.Error
	UpdateSequenceNumber(username string, seq uint64) errors.Error
	UpdatePubKey(username, resetPubKey, TxPubKey, appPubKey string) errors.Error
	UpdateBalance(username string, balance string) errors.Error
	IsEnable() bool
}
