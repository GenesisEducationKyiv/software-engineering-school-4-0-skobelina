package rates

import (
	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/configs"
	"github.com/skobelina/currency_converter/infrastructure/currencies"
	"github.com/skobelina/currency_converter/pkg/utils/serializer"
	"gorm.io/gorm"
)

type RateService struct {
	repo    *gorm.DB
	handler currencies.CurrencyHandler
	config  *configs.Config
}

func NewService(repo *gorm.DB, config *configs.Config) *RateService {
	providerExchangeRates := &currencies.ProviderExchangeRates{}
	providerCurrencyBeacon := &currencies.ProviderCurrencyBeacon{}
	providerExchangeRates.SetNext(providerCurrencyBeacon)
	return &RateService{repo: repo, handler: providerExchangeRates, config: config}
}

func (s *RateService) Get() (*float64, error) {
	rate, err := s.handler.Handle(s.config)
	if err != nil {
		logrus.Errorf("RateService - Error fetching rate: %v", err)
		return nil, serializer.NewInternalServerErrorf("all providers failed: %v", err)
	}
	logrus.Info("RateService - Successfully fetched rate")
	return &rate, nil
}
