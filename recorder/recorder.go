package recorder

import (
	"github.com/lino-network/lino/recorder/dbtestutil"
	dRepo "github.com/lino-network/lino/recorder/donation/repository"
)

type Recorder struct {
	DonationRepo dRepo.DonationRepository
}

func NewRecorder() Recorder {
	db, _ := dbtestutil.NewDBConn()
	donationRepo, _ := dbtestutil.NewDonationDB(db)
	dbtestutil.NewDonationDB(db)
	return Recorder{
		DonationRepo: donationRepo,
	}
}
