package recorder

import (
	"github.com/lino-network/lino/recorder/dbtestutil"
	dRepo "github.com/lino-network/lino/recorder/donation/repository"
	iRepo "github.com/lino-network/lino/recorder/inflation/repository"
)

type Recorder struct {
	DonationRepo  dRepo.DonationRepository
	InflationRepo iRepo.InflationRepository
}

func NewRecorder() Recorder {
	db, _ := dbtestutil.NewDBConn()
	donationRepo, _ := dbtestutil.NewDonationDB(db)
	inflationRepo, _ := dbtestutil.NewInflationDB(db)
	dbtestutil.NewDonationDB(db)
	return Recorder{
		DonationRepo:  donationRepo,
		InflationRepo: inflationRepo,
	}
}
