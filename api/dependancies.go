package api

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/domains/mails"
	"github.com/skobelina/currency_converter/domains/rates"
	"github.com/skobelina/currency_converter/domains/subscribers"
	"github.com/skobelina/currency_converter/repo"
	"github.com/skobelina/currency_converter/utils/currencies"
)

type dependancies struct {
	rates         rates.Service
	subscribers   subscribers.Service
	mailService   mails.Service
	currencyRates float64
}

func registerDependencies() *dependancies {
	repo, err := repo.Connect(databaseURL)
	if err != nil {
		panic(err)
	}
	if err := repo.AutoMigrate(&subscribers.Subscriber{}); err != nil {
		logrus.Infof("failed to migrate database: %v", err)
	}
	rates := rates.NewService(repo)
	subscribers := subscribers.NewService(repo)
	mailService := mails.NewService(mails.DefaultMailSendAddress, mails.DefaultMailHost)
	currencyRates, err := currencies.GetCurrencyRates()
	if err != nil {
		logrus.Infof("cannot preload currency rates: %v\n", err)
	}
	return &dependancies{
		rates:         rates,
		subscribers:   subscribers,
		mailService:   mailService,
		currencyRates: currencyRates,
	}
}
