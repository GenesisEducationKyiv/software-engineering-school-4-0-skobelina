package cronjobs

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/domains/mails"
	"github.com/skobelina/currency_converter/domains/rates"
	"github.com/skobelina/currency_converter/repo"

	"github.com/skobelina/currency_converter/domains/subscribers"

	"gorm.io/gorm"
)

type CronJobService struct {
	Repo        *gorm.DB
	Mail        mails.MailServiceInterface
	Rates       rates.RateServiceInterface
	Subscribers subscribers.SubscriberServiceInterface
}

func NewService(config *CronJobConfig) *CronJobService {
	repo, err := repo.Connect(config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	rates := rates.NewService(repo)
	subscribers := subscribers.NewService(repo)
	mails := mails.NewService(mails.DefaultMailSendAddress, mails.DefaultMailHost)
	return &CronJobService{
		repo,
		mails,
		rates,
		subscribers,
	}
}

func (s *CronJobService) Close() error {
	db, err := s.Repo.DB()
	if err != nil {
		logrus.Errorf("CronJobs: Close: %v", err)
		return err
	}
	err = db.Close()
	if err != nil {
		logrus.Errorf("CronJobs: Close: %v", err)
		return err
	}
	return nil
}
