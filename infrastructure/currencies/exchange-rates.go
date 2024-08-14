package currencies

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/configs"
	"github.com/skobelina/currency_converter/internal/constants"
	"github.com/skobelina/currency_converter/pkg/metrics"
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

func (p *ProviderExchangeRates) Handle(config *configs.Config) (float64, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RequestDuration.WithLabelValues("GET", config.AppCurrencyBeaconURL).Observe(duration.Seconds())
	}()

	var myClient = &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, config.AppCurrencyBeaconURL, http.NoBody)
	if err != nil {
		logrus.Errorf("ProviderExchangeRates - Failed to create request: %v", err)
		return p.BaseHandler.Handle(config)
	}
	req.Header.Add("apikey", config.AppCurrencyExchangeKey)
	resp, err := myClient.Do(req)
	if err != nil {
		logrus.Errorf("ProviderExchangeRates - Failed to perform request: %v", err)
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Errorf("Error closing response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("ProviderExchangeRates - Failed to read response body: %v", err)
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("ProviderExchangeRates - Unexpected status code: %d", resp.StatusCode)
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}

	var data CurrencyData
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.Errorf("ProviderExchangeRates - Failed to unmarshal response: %v", err)
		logrus.Errorf("ProviderExchangeRates - Response body: %s", string(body))
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}

	usdRate, usdExists := data.Rates[constants.CurrencyUSD]
	uahRate, uahExists := data.Rates[constants.CurrencyUAH]
	if !usdExists || !uahExists {
		logrus.Errorf("ProviderExchangeRates - USD or UAH rate not found")
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}

	logrus.Infof("ProviderExchangeRates - Successfully fetched exchange rates: USD: %f, UAH: %f", usdRate, uahRate)
	metrics.RequestCount.WithLabelValues("GET", "success").Inc()
	return uahRate / usdRate, nil
}
