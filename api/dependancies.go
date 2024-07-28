package api

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/configs"
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
	logrus.Info("Loading configuration")
	config, err := configs.LoadConfig(".env")
	if err != nil {
		logrus.Fatal("Failed to load config: ", err)
	}

	logrus.Info("Connecting to the database")
	repo, err := repo.Connect(config.DatabaseURL)
	if err != nil {
		logrus.Fatal("Failed to connect to database: ", err)
	}

	logrus.Info("Migrating database schema")
	if err := repo.AutoMigrate(&subscribers.Subscriber{}); err != nil {
		logrus.Fatalf("Failed to migrate database: %v", err)
	}

	logrus.Info("Initializing repositories and services")
	subscriberRepo := subscribers.NewRepository(repo)
	ratesService := rates.NewService(repo, config)

	logrus.Info("Connecting to RabbitMQ")
	rabbitMQ, err := queue.NewRabbitMQ(config.RabbitMQURL, "events")
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	logrus.Info("Initializing Saga orchestrator")
	saga := subscribers.NewSaga(rabbitMQ)

	logrus.Info("Initializing Subscriber service")
	subscribersService := subscribers.NewService(subscriberRepo, rabbitMQ, saga)

	logrus.Info("Preloading currency rates")
	currencyRates, err := ratesService.Get()
	if err != nil {
		logrus.Warnf("Cannot preload currency rates: %v", err)
	}

	logrus.Info("Dependencies initialized successfully")

	return &dependencies{
		Repo:          repo,
		Rates:         ratesService,
		Subscribers:   subscribersService,
		RabbitMQ:      rabbitMQ,
		currencyRates: currencyRates,
	}
}
