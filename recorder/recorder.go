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
	db, err := dbtestutil.NewDBConn()
	if err != nil {
		panic(err)
	}
	donationRepo, err := dbtestutil.NewDonationDB(db)
	if err != nil {
		panic(err)
	}
	inflationRepo, err := dbtestutil.NewInflationDB(db)
	if err != nil {
		panic(err)
	}
	stakeStatRepo, err := dbtestutil.NewStakeStatDB(db)
	if err != nil {
		panic(err)
	}
	postRewardStatRepo, err := dbtestutil.NewPostRewardDB(db)
	if err != nil {
		panic(err)
	}
	stakeRepo, err := dbtestutil.NewStakeDB(db)
	if err != nil {
		panic(err)
	}
	topContentRepo, err := dbtestutil.NewTopContentDB(db)
	if err != nil {
		panic(err)
	}
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
