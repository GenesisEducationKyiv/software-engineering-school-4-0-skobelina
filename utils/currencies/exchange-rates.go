package currencies

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/constants"
)

type ProviderExchangeRates struct {
	BaseHandler
}

type CurrencyData struct {
	Base string
	Date string
	Motd struct {
		Msg     string
		Url     string
		Success bool
	}
	Rates map[string]float64
}

var (
	currencyExchangeApi = os.Getenv("APP_CURRENCY_EXCHANGE_URL")
	apiKey              = os.Getenv("APP_CURRENCY_EXCHANGE_KEY")
)

func (p *ProviderExchangeRates) Handle() (float64, error) {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", currencyExchangeApi, nil)
	if err != nil {
		logrus.Errorf("ProviderExchangeRates - Failed to create request: %v", err)
		return p.BaseHandler.Handle()
	}
	req.Header.Add("apikey", apiKey)
	resp, err := myClient.Do(req)
	if err != nil {
		logrus.Errorf("ProviderExchangeRates - Failed to perform request: %v", err)
		return p.BaseHandler.Handle()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("ProviderExchangeRates - Failed to read response body: %v", err)
		return p.BaseHandler.Handle()
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("ProviderExchangeRates - Unexpected status code: %d", resp.StatusCode)
		return p.BaseHandler.Handle()
	}

	var data CurrencyData
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.Errorf("ProviderExchangeRates - Failed to unmarshal response: %v", err)
		logrus.Errorf("ProviderExchangeRates - Response body: %s", string(body))
		return p.BaseHandler.Handle()
	}

	usdRate, usdExists := data.Rates[constants.CurrencyUSD]
	uahRate, uahExists := data.Rates[constants.CurrencyUAH]
	if !usdExists || !uahExists {
		logrus.Errorf("ProviderExchangeRates - USD or UAH rate not found")
		return p.BaseHandler.Handle()
	}
	return uahRate / usdRate, nil
}
