package rates

import (
	"gorm.io/gorm"

	"github.com/skobelina/currency_converter/utils/currencies"
	errors "github.com/skobelina/currency_converter/utils/errors"
)

type Service interface {
	Get() (*float64, error)
}

type service struct {
	repo *gorm.DB
}

func NewService(repo *gorm.DB) Service {
	return &service{repo}
}

func (s *service) Get() (*float64, error) {
	rate, err := currencies.GetCurrencyRates()
	if err != nil {
		return nil, errors.NewInternalServerErrorf("cannot get currency rates: %v", err)
	}
	return &rate, nil
}
