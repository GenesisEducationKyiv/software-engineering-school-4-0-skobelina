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
	repo        *gorm.DB
	mail        mails.MailService
	rates       rates.RateServiceInterface
	subscribers subscribers.SubscriberServiceInterface
}

func NewService(config *CronJobConfig) *CronJobService {
	repo, err := repo.Connect(config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	rates := rates.NewService(repo)
	subscribers := subscribers.NewService(repo)
	return &CronJobService{
		repo,
		*config.MailService,
		rates,
		subscribers,
	}
}

func (s *CronJobService) Close() error {
	db, err := s.repo.DB()
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
