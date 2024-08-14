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

type ProviderCurrencyBeacon struct {
	BaseHandler
}

type CurrencyBeaconResponse struct {
	Meta     MetaData                   `json:"meta"`
	Response CurrencyBeaconDataResponse `json:"response"`
}

type MetaData struct {
	Code       int    `json:"code"`
	Disclaimer string `json:"disclaimer"`
}

type CurrencyBeaconDataResponse struct {
	Date  string             `json:"date"`
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func (p *ProviderCurrencyBeacon) Handle(config *configs.Config) (float64, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RequestDuration.WithLabelValues("GET", config.AppCurrencyBeaconURL).Observe(duration.Seconds())
	}()

	var myClient = &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, config.AppCurrencyBeaconURL+"?api_key="+config.AppCurrencyBeaconKey+"&base=USD&symbols=UAH", http.NoBody)
	if err != nil {
		logrus.Errorf("ProviderCurrencyBeacon - Failed to create request: %v", err)
		return p.BaseHandler.Handle(config)
	}
	resp, err := myClient.Do(req)
	if err != nil {
		logrus.Errorf("ProviderCurrencyBeacon - Failed to perform request: %v", err)
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
		logrus.Errorf("ProviderCurrencyBeacon - Failed to read response body: %v", err)
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}
	logrus.Infof("ProviderCurrencyBeacon - Response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("ProviderCurrencyBeacon - Unexpected status code: %d", resp.StatusCode)
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}

	var data CurrencyBeaconResponse
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.Errorf("ProviderCurrencyBeacon - Failed to unmarshal response: %v", err)
		logrus.Errorf("ProviderCurrencyBeacon - Response body: %s", string(body))
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}

	if data.Meta.Code != 200 {
		logrus.Errorf("ProviderCurrencyBeacon - API returned error code: %d", data.Meta.Code)
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}

	uahRate, uahExists := data.Response.Rates[constants.CurrencyUAH]
	if !uahExists {
		logrus.Errorf("ProviderCurrencyBeacon - UAH rate not found")
		metrics.RequestCount.WithLabelValues("GET", "error").Inc()
		return p.BaseHandler.Handle(config)
	}

	logrus.Infof("ProviderCurrencyBeacon - Successfully fetched UAH rate: %f", uahRate)
	metrics.RequestCount.WithLabelValues("GET", "success").Inc()
	return uahRate, nil
}
