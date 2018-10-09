package recorder

import (
	"github.com/lino-network/lino/recorder/dbtestutil"
	dRepo "github.com/lino-network/lino/recorder/donation/repository"
	iRepo "github.com/lino-network/lino/recorder/inflation/repository"
	pRepo "github.com/lino-network/lino/recorder/postreward/repository"
	stakeRepo "github.com/lino-network/lino/recorder/stake/repository"
	sRepo "github.com/lino-network/lino/recorder/stakestat/repository"
	tRepo "github.com/lino-network/lino/recorder/topcontent/repository"
)

type Recorder struct {
	DonationRepo         dRepo.DonationRepository
	InflationRepo        iRepo.InflationRepository
	StakeStatRepository  sRepo.StakeStatRepository
	PostRewardRepository pRepo.PostRewardRepository
	StakeRepository      stakeRepo.StakeRepository
	TopContentRepository tRepo.TopContentRepository
}

func NewRecorder() Recorder {
	db, _ := dbtestutil.NewDBConn()
	donationRepo, _ := dbtestutil.NewDonationDB(db)
	inflationRepo, _ := dbtestutil.NewInflationDB(db)
	stakeStatRepo, _ := dbtestutil.NewStakeStatDB(db)
	postRewardStatRepo, _ := dbtestutil.NewPostRewardDB(db)
	stakeRepo, _ := dbtestutil.NewStakeDB(db)
	topContentRepo, _ := dbtestutil.NewTopContentDB(db)

	dbtestutil.NewDonationDB(db)
	return Recorder{
		DonationRepo:         donationRepo,
		InflationRepo:        inflationRepo,
		StakeStatRepository:  stakeStatRepo,
		PostRewardRepository: postRewardStatRepo,
		StakeRepository:      stakeRepo,
		TopContentRepository: topContentRepo,
	}
}
