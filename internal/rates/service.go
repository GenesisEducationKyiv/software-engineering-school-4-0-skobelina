package rates

import (
	"github.com/skobelina/currency_converter/infrastructure/currencies"
	errors "github.com/skobelina/currency_converter/pkg/errors"
	"gorm.io/gorm"
)

type RateService struct {
	repo    *gorm.DB
	handler currencies.CurrencyHandler
}

func NewService(repo *gorm.DB) *RateService {
	providerExchangeRates := &currencies.ProviderExchangeRates{}
	providerCurrencyBeacon := &currencies.ProviderCurrencyBeacon{}
	providerExchangeRates.SetNext(providerCurrencyBeacon)
	return &RateService{repo: repo, handler: providerExchangeRates}
}

func (s *RateService) Get() (*float64, error) {
	rate, err := s.handler.Handle()
	if err != nil {
		return nil, errors.NewInternalServerErrorf("all providers failed: %v", err)
	}
	return &rate, nil
}
