package cronjobs

import (
	"github.com/skobelina/currency_converter/internal/rates"
	"github.com/skobelina/currency_converter/pkg/queue"
)

type CronJobService struct {
	rates    rates.RateServiceInterface
	rabbitMQ *queue.RabbitMQ
}

func NewService(rates rates.RateServiceInterface, rabbitMQ *queue.RabbitMQ) *CronJobService {
	return &CronJobService{
		rates:    rates,
		rabbitMQ: rabbitMQ,
	}
}
