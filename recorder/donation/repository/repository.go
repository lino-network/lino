package repository

import (
	"github.com/lino-network/lino/recorder/donation"
	errors "github.com/lino-network/lino/recorder/errors"
)

type DonationRepository interface {
	Get(username string) (*donation.Donation, errors.Error)
	Add(info *donation.Donation) errors.Error
}
