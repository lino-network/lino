package recorder

import (
	"github.com/lino-network/lino/recorder/dbtestutil"
	dRepo "github.com/lino-network/lino/recorder/donation/repository"
	iRepo "github.com/lino-network/lino/recorder/inflation/repository"
	pRepo "github.com/lino-network/lino/recorder/postReward/repository"
	sRepo "github.com/lino-network/lino/recorder/stakeStat/repository"
)

type Recorder struct {
	DonationRepo         dRepo.DonationRepository
	InflationRepo        iRepo.InflationRepository
	StakeStatRepository  sRepo.StakeStatRepository
	PostRewardRepository pRepo.PostRewardRepository
}

func NewRecorder() Recorder {
	db, _ := dbtestutil.NewDBConn()
	donationRepo, _ := dbtestutil.NewDonationDB(db)
	inflationRepo, _ := dbtestutil.NewInflationDB(db)
	stakeStatRepo, _ := dbtestutil.NewStakeStatDB(db)
	postRewardStatRepo, _ := dbtestutil.NewPostRewardDB(db)

	dbtestutil.NewDonationDB(db)
	return Recorder{
		DonationRepo:         donationRepo,
		InflationRepo:        inflationRepo,
		StakeStatRepository:  stakeStatRepo,
		PostRewardRepository: postRewardStatRepo,
	}
}
