package cronjobs

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/domains/mails"
	"github.com/skobelina/currency_converter/domains/rates"

	"github.com/skobelina/currency_converter/domains/subscribers"

	"gorm.io/gorm"
)

type CronJobService struct {
	Repo        *gorm.DB
	Mail        mails.MailServiceInterface
	Rates       rates.RateServiceInterface
	Subscribers subscribers.SubscriberServiceInterface
}

func NewService(repo *gorm.DB, mail mails.MailServiceInterface, rates rates.RateServiceInterface, subscribers subscribers.SubscriberServiceInterface) *CronJobService {
	return &CronJobService{
		Repo:        repo,
		Mail:        mail,
		Rates:       rates,
		Subscribers: subscribers,
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
