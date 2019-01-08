package recorder

import (
	"database/sql"
	"fmt"
	"log"

	bhRepo "github.com/lino-network/lino/recorder/balancehistory/repository"
	dRepo "github.com/lino-network/lino/recorder/donation/repository"
	iRepo "github.com/lino-network/lino/recorder/inflation/repository"
	postRepo "github.com/lino-network/lino/recorder/post/repository"
	pRepo "github.com/lino-network/lino/recorder/postreward/repository"
	rRepo "github.com/lino-network/lino/recorder/reward/repository"
	stakeRepo "github.com/lino-network/lino/recorder/stake/repository"
	sRepo "github.com/lino-network/lino/recorder/stakestat/repository"
	tRepo "github.com/lino-network/lino/recorder/topcontent/repository"
	uRepo "github.com/lino-network/lino/recorder/user/repository"
)

type Recorder struct {
	NewVersionOnly           bool
	DonationRepo             dRepo.DonationRepository
	InflationRepo            iRepo.InflationRepository
	StakeStatRepository      sRepo.StakeStatRepository
	PostRewardRepository     pRepo.PostRewardRepository
	StakeRepository          stakeRepo.StakeRepository
	TopContentRepository     tRepo.TopContentRepository
	UserRepository           uRepo.UserRepository
	BalanceHistoryRepository bhRepo.BalanceHistoryRepository
	RewardRepository         rRepo.RewardRepository
	PostRepository           postRepo.PostRepository
}

var conn *sql.DB

// NewDBConn returns a new sql db connection
func NewDBConn(dbUsername, dbPassword, dbHost, dbPort, dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		dbUsername,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	conn, err := sql.Open("mysql", dsn)
	return conn, err
}

func NewRecorder() Recorder {
	configPath = fs.String("config.file", "", "Config file path.")
	configs, e := NewConfig(*configPath)
	if e != nil {
		log.Fatal().Err(e).Msg("Failed to create config.")
	}
	db, err := NewDBConn(configs.DBUsername(), configs.DBPassword(), configs.DBHost(), configs.DBPort(), configs.DBName())
	if err != nil {
		panic(err)
	}
	donationRepo, err := dRepo.NewDonationDB(db)
	if err != nil {
		panic(err)
	}
	inflationRepo, err := iRepo.NewinflationDB(db)
	if err != nil {
		panic(err)
	}
	stakeStatRepo, err := sRepo.NewStakeStatDB(db)
	if err != nil {
		panic(err)
	}
	postRewardStatRepo, err := pRepo.NewPostRewardDB(db)
	if err != nil {
		panic(err)
	}
	stakeRepo, err := stakeRepo.NewStakeDB(db)
	if err != nil {
		panic(err)
	}
	topContentRepo, err := tRepo.NewTopContentDB(db)
	if err != nil {
		panic(err)
	}
	balanceHistoryRepo, err := bhRepo.NewBalanceHistoryDB(db)
	if err != nil {
		panic(err)
	}
	postRepo, err := postRepo.NewPostDB(db)
	if err != nil {
		panic(err)
	}
	rewardRepo, err := rRepo.NewRewardDB(db)
	if err != nil {
		panic(err)
	}
	userRepo, err := uRepo.NewUserDB(db)
	if err != nil {
		panic(err)
	}
	return Recorder{
		DonationRepo:             donationRepo,
		InflationRepo:            inflationRepo,
		StakeStatRepository:      stakeStatRepo,
		PostRewardRepository:     postRewardStatRepo,
		StakeRepository:          stakeRepo,
		TopContentRepository:     topContentRepo,
		UserRepository:           userRepo,
		BalanceHistoryRepository: balanceHistoryRepo,
		RewardRepository:         rewardRepo,
		PostRepository:           postRepo,
	}
}
