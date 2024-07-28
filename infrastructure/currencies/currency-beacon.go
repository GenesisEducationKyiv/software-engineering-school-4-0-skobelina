package currencies

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skobelina/currency_converter/internal/constants"
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

var (
	currencyBeaconApi = os.Getenv("APP_CURRENCY_BEACON_URL")
	apiKeyBeacon      = os.Getenv("APP_CURRENCY_BEACON_KEY")
)

func (p *ProviderCurrencyBeacon) Handle() (float64, error) {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", currencyBeaconApi+"?api_key="+apiKeyBeacon+"&base=USD&symbols=UAH", nil)
	if err != nil {
		logrus.Errorf("ProviderCurrencyBeacon - Failed to create request: %v", err)
		return p.BaseHandler.Handle()
	}
	resp, err := myClient.Do(req)
	if err != nil {
		logrus.Errorf("ProviderCurrencyBeacon - Failed to perform request: %v", err)
		return p.BaseHandler.Handle()
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Errorf("Error closing response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("ProviderCurrencyBeacon - Failed to read response body: %v", err)
		return p.BaseHandler.Handle()
	}
	logrus.Infof("ProviderCurrencyBeacon - Response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("ProviderCurrencyBeacon - Unexpected status code: %d", resp.StatusCode)
		return p.BaseHandler.Handle()
	}

	var data CurrencyBeaconResponse
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.Errorf("ProviderCurrencyBeacon - Failed to unmarshal response: %v", err)
		logrus.Errorf("ProviderCurrencyBeacon - Response body: %s", string(body))
		return p.BaseHandler.Handle()
	}

	if data.Meta.Code != 200 {
		logrus.Errorf("ProviderCurrencyBeacon - API returned error code: %d", data.Meta.Code)
		return p.BaseHandler.Handle()
	}

	uahRate, uahExists := data.Response.Rates[constants.CurrencyUAH]
	if !uahExists {
		logrus.Errorf("ProviderCurrencyBeacon - UAH rate not found")
		return p.BaseHandler.Handle()
	}
	return uahRate, nil
}
