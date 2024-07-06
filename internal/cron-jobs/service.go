package cronjobs

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/internal/mails"
	"github.com/skobelina/currency_converter/internal/rates"

	"github.com/skobelina/currency_converter/internal/subscribers"

	"gorm.io/gorm"
)

type CronJobService struct {
	repo        *gorm.DB
	mail        mails.MailServiceInterface
	rates       rates.RateServiceInterface
	subscribers subscribers.SubscriberServiceInterface
}

func NewService(repo *gorm.DB, mail mails.MailServiceInterface, rates rates.RateServiceInterface, subscribers subscribers.SubscriberServiceInterface) *CronJobService {
	return &CronJobService{
		repo:        repo,
		mail:        mail,
		rates:       rates,
		subscribers: subscribers,
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
