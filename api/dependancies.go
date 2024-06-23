package api

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/domains/mails"
	"github.com/skobelina/currency_converter/domains/rates"
	"github.com/skobelina/currency_converter/domains/subscribers"
	"github.com/skobelina/currency_converter/repo"
	"gorm.io/gorm"
)

type dependencies struct {
	Repo          *gorm.DB
	Rates         *rates.RateService
	Subscribers   *subscribers.SubscriberService
	MailService   *mails.MailService
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
	mailService := mails.NewService(mails.DefaultMailSendAddress, mails.DefaultMailHost)
	currencyRates, err := rates.Get()
	if err != nil {
		logrus.Infof("cannot preload currency rates: %v\n", err)
	}
	return &dependencies{
		Repo:          repo,
		Rates:         rates,
		Subscribers:   subscribers,
		MailService:   mailService,
		currencyRates: currencyRates,
	}
}
