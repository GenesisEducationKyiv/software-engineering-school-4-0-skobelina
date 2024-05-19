package cronjobs

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/domains/mails"
	"github.com/skobelina/currency_converter/domains/rates"

	"github.com/skobelina/currency_converter/domains/subscribers"
	"github.com/skobelina/currency_converter/repo"

	"gorm.io/gorm"
)

type Service interface {
	NotificationExchangeRates() error
	Close() error
}

type service struct {
	repo        *gorm.DB
	mail        mails.Service
	rates       rates.Service
	subscribers subscribers.Service
}

func NewService(config *CronJobConfig) Service {
	repo, err := repo.Connect(config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	rates := rates.NewService(repo)
	subscribers := subscribers.NewService(repo)
	return &service{
		repo,
		config.MailService,
		rates,
		subscribers,
	}
}

func (s *service) Close() error {
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
