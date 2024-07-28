package cronjobs

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/internal/rates"
	"github.com/skobelina/currency_converter/pkg/queue"

	"github.com/skobelina/currency_converter/internal/subscribers"

	"gorm.io/gorm"
)

type CronJobService struct {
	repo        *gorm.DB
	rates       rates.RateServiceInterface
	subscribers subscribers.SubscriberServiceInterface
	rabbitMQ    *queue.RabbitMQ
}

func NewService(repo *gorm.DB, rates rates.RateServiceInterface, subscribers subscribers.SubscriberServiceInterface, rabbitMQ *queue.RabbitMQ) *CronJobService {
	return &CronJobService{
		repo:        repo,
		rates:       rates,
		subscribers: subscribers,
		rabbitMQ:    rabbitMQ,
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
