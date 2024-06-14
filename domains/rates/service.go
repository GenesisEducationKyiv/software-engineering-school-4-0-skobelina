package rates

import (
	"gorm.io/gorm"

	"github.com/skobelina/currency_converter/utils/currencies"
	errors "github.com/skobelina/currency_converter/utils/errors"
)

type RateService struct {
	repo *gorm.DB
}

func NewService(repo *gorm.DB) *RateService {
	return &RateService{repo}
}

func (s *RateService) Get() (*float64, error) {
	rate, err := currencies.GetCurrencyRates()
	if err != nil {
		return nil, errors.NewInternalServerErrorf("cannot get currency rates: %v", err)
	}
	return &rate, nil
}
