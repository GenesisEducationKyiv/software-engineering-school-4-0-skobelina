package api

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/internal/rates"
	"github.com/skobelina/currency_converter/internal/subscribers"
	"github.com/skobelina/currency_converter/pkg/queue"
	"github.com/skobelina/currency_converter/pkg/repo"
	"gorm.io/gorm"
)

type dependencies struct {
	Repo          *gorm.DB
	Rates         *rates.RateService
	Subscribers   *subscribers.SubscriberService
	RabbitMQ      *queue.RabbitMQ
	currencyRates *float64
}

func registerDependencies() *dependencies {
	repo, err := repo.Connect(databaseURL)
	if err != nil {
		panic(err)
	}
	if err := repo.AutoMigrate(&subscribers.Subscriber{}); err != nil {
		logrus.Infof("failed to migrate database: %v", err)
	}
	subscriberRepo := subscribers.NewRepository(repo)
	rates := rates.NewService(repo)
	subscribers := subscribers.NewService(subscriberRepo)
	rabbitMQ, err := queue.NewRabbitMQ(rabbitMQURL, "events")
	if err != nil {
		logrus.Infof("failed to connect to RabbitMQ: %v", err)
	}
	currencyRates, err := rates.Get()
	if err != nil {
		logrus.Infof("cannot preload currency rates: %v\n", err)
	}
	return &dependencies{
		Repo:          repo,
		Rates:         rates,
		Subscribers:   subscribers,
		RabbitMQ:      rabbitMQ,
		currencyRates: currencyRates,
	}
}
